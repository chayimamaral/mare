package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/chayimamaral/vecx/agente-local/internal/config"
	"github.com/chayimamaral/vecx/agente-local/internal/httpserver"
	"github.com/chayimamaral/vecx/agente-local/internal/provider/pkcs11"
	"github.com/chayimamaral/vecx/agente-local/internal/usecase"
)

func main() {
	cfg := config.Load()

	provider := pkcs11.NewProvider(cfg.PKCS11LibraryLinux, cfg.PKCS11LibraryWindow)
	signUC := usecase.NewSignUseCase(provider)
	handler := httpserver.NewHandler(signUC, nil)
	server := httpserver.NewServer(cfg.HTTPAddr, cfg.AllowedOrigins, cfg.SharedSecret, handler)

	log.Printf("agente local iniciado em http://%s", cfg.HTTPAddr)
	if len(cfg.AllowedOrigins) > 0 {
		log.Printf("cors liberado para: %v", cfg.AllowedOrigins)
	}
	if cfg.SharedSecret != "" {
		log.Printf("autenticacao local habilitada via X-Local-Agent-Secret")
	}

	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-stopCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("erro ao encerrar servidor: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("erro no servidor local: %v", err)
	}
}
