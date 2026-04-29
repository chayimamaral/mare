package main

import (
	"context"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chayimamaral/vecontab/backend/frontend"
	"github.com/chayimamaral/vecontab/backend/internal/config"
	"github.com/chayimamaral/vecontab/backend/internal/db"
	"github.com/chayimamaral/vecontab/backend/internal/httpapi"
	apiMiddleware "github.com/chayimamaral/vecontab/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
	"github.com/chayimamaral/vecontab/backend/internal/service"
	"github.com/chayimamaral/vecontab/backend/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()

	auditPool, err := db.NewAuditPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatalf("connect vecx audit postgres: %v", err)
	}
	defer auditPool.Close()

	auditRepo := repository.NewVecxAuditRepository(auditPool)
	if err := auditRepo.EnsureSchema(ctx); err != nil {
		log.Fatalf("vecx audit schema: %v", err)
	}

	// Compatibilidade com tokens antigos sem tenant.schema_name.
	apiMiddleware.SetTenantSchemaResolver(func(ctx context.Context, tenantID string) (string, error) {
		return db.ResolveTenantSchema(ctx, pool, tenantID)
	})
	apiMiddleware.SetTenantConnPool(pool)

	staticRoot, err := fs.Sub(frontend.FS, "out")
	if err != nil {
		log.Fatalf("frontend static: %v", err)
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.NewRouter(cfg, pool, auditPool, staticRoot),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil && (cfg.Runtime == "binary" || cfg.Runtime == "desktop") && strings.Contains(strings.ToLower(err.Error()), "address already in use") {
		fallbackPorts := []string{"3334", "3335", "3336", "3337", "3338"}
		for _, p := range fallbackPorts {
			if p == cfg.Port {
				continue
			}
			l, e := net.Listen("tcp", ":"+p)
			if e == nil {
				log.Printf("porta %s ocupada no modo %s; usando porta alternativa %s", cfg.Port, cfg.Runtime, p)
				cfg.Port = p
				server.Addr = ":" + p
				listener = l
				err = nil
				break
			}
		}
	}
	if err != nil {
		log.Fatalf("listen http: %v", err)
	}

	if cfg.CompromissosWorkerEnabled {
		monRepo := repository.NewMonitorOperacaoRepository(pool)
		w, err := worker.NewCompromissosWorker(pool, cfg, monRepo)
		if err != nil {
			log.Fatalf("init compromissos worker: %v", err)
		}
		go w.Start(ctx)
	}

	if cfg.NFESyncWorkerEnabled {
		certificadoService, _ := service.NewCertificadoService(
			repository.NewCertificadoRepository(pool),
			repository.NewCertificadoClienteRepository(pool),
			cfg.CertCryptoKeyHex,
		)
		nfeSerproRepo := repository.NewNFESerproRepository(pool)
		nfeSerproService := service.NewNFESerproService(nfeSerproRepo, service.NewSerproService(cfg, certificadoService), certificadoService)
		monRepo := repository.NewMonitorOperacaoRepository(pool)
		nw, err := worker.NewNFESyncWorker(pool, cfg, nfeSerproService, monRepo)
		if err != nil {
			log.Fatalf("init nfe sync worker: %v", err)
		}
		go nw.Start(ctx)
	}

	go func() {
		log.Printf("backendgo listening on :%s (runtime=%s)", cfg.Port, cfg.Runtime)
		errCh <- server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown: %v", err)
		}
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve http: %v", err)
		}
	}
}
