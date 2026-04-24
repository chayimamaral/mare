package worker

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/chayimamaral/vecontab/backend/internal/config"
	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
	"github.com/chayimamaral/vecontab/backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type NFESyncWorker struct {
	pool    *pgxpool.Pool
	cfg     config.Config
	cron    *cron.Cron
	running int32
	monitor *repository.MonitorOperacaoRepository
	svc     *service.NFESerproService
}

const nfeSyncAdvisoryLockKey int64 = 9482217703

func NewNFESyncWorker(pool *pgxpool.Pool, cfg config.Config, svc *service.NFESerproService, monitor *repository.MonitorOperacaoRepository) (*NFESyncWorker, error) {
	if svc == nil {
		return nil, fmt.Errorf("nfe sync worker: servico nfe nil")
	}
	loc, err := time.LoadLocation(cfg.NFESyncWorkerTimezone)
	if err != nil {
		return nil, fmt.Errorf("timezone inválida do worker nfe sync: %w", err)
	}

	c := cron.New(
		cron.WithLocation(loc),
		cron.WithParser(cron.NewParser(cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	w := &NFESyncWorker{
		pool:    pool,
		cfg:     cfg,
		cron:    c,
		monitor: monitor,
		svc:     svc,
	}

	if _, err := c.AddFunc(cfg.NFESyncWorkerCron, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		if err := w.runOnce(ctx); err != nil {
			log.Printf("worker nfe sync: erro na execução agendada: %v", err)
		}
	}); err != nil {
		return nil, fmt.Errorf("cron inválido do worker nfe sync: %w", err)
	}

	return w, nil
}

func (w *NFESyncWorker) Start(ctx context.Context) {
	log.Printf("worker nfe sync: habilitado cron=%q tz=%q ambiente=%q", w.cfg.NFESyncWorkerCron, w.cfg.NFESyncWorkerTimezone, w.cfg.NFESyncWorkerAmbiente)
	w.cron.Start()

	if w.cfg.NFESyncWorkerRunOnStartup {
		go func() {
			startupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			defer cancel()
			if err := w.runOnce(startupCtx); err != nil {
				log.Printf("worker nfe sync: erro no run_on_startup: %v", err)
			}
		}()
	}

	<-ctx.Done()
	stopCtx := w.cron.Stop()
	select {
	case <-stopCtx.Done():
	case <-time.After(10 * time.Second):
	}
	log.Printf("worker nfe sync: finalizado")
}

func (w *NFESyncWorker) minGap() time.Duration {
	s := w.cfg.NFESyncWorkerMinIntervalSecs
	if s < 60 {
		s = 60
	}
	return time.Duration(s) * time.Second
}

func (w *NFESyncWorker) runOnce(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&w.running, 0, 1) {
		return nil
	}
	defer atomic.StoreInt32(&w.running, 0)

	conn, err := w.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn worker nfe sync: %w", err)
	}
	defer conn.Release()

	var gotLock bool
	if err := conn.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, nfeSyncAdvisoryLockKey).Scan(&gotLock); err != nil {
		return fmt.Errorf("obter advisory lock nfe sync: %w", err)
	}
	if !gotLock {
		log.Printf("worker nfe sync: execução já em andamento em outra réplica")
		return nil
	}
	defer func() {
		_, _ = conn.Exec(context.Background(), `SELECT pg_advisory_unlock($1)`, nfeSyncAdvisoryLockKey)
	}()

	trows, err := w.pool.Query(ctx, `
		SELECT tenant_id::text, schema_name
		FROM public.tenant_schema_catalog
		WHERE NULLIF(BTRIM(schema_name), '') IS NOT NULL
		ORDER BY tenant_id`)
	if err != nil {
		return fmt.Errorf("listar tenants worker nfe sync: %w", err)
	}
	defer trows.Close()

	type tenantRow struct {
		TenantID   string
		SchemaName string
	}
	var tenants []tenantRow
	for trows.Next() {
		var t tenantRow
		if err := trows.Scan(&t.TenantID, &t.SchemaName); err != nil {
			return fmt.Errorf("scan tenant_schema_catalog nfe sync: %w", err)
		}
		tenants = append(tenants, t)
	}
	if err := trows.Err(); err != nil {
		return err
	}

	repo := repository.NewNFESerproRepository(w.pool)
	now := time.Now().UTC()
	gap := w.minGap()
	var totalOK, totalErr int
	for _, t := range tenants {
		due, qerr := repo.ListSyncEstadosDue(ctx, t.SchemaName, now, gap, 80)
		if qerr != nil {
			log.Printf("worker nfe sync: list due tenant=%s schema=%s: %v", t.TenantID, t.SchemaName, qerr)
			totalErr++
			continue
		}
		for _, st := range due {
			_, serr := w.svc.SincronizarPorProvider(ctx, t.SchemaName, t.TenantID, st.Provider, st.UF, st.CNPJ, w.cfg.NFESyncWorkerAmbiente, false)
			if serr != nil {
				totalErr++
				log.Printf("worker nfe sync: sync tenant=%s provider=%s uf=%s cnpj=%s: %v", t.TenantID, st.Provider, st.UF, st.CNPJ, serr)
				continue
			}
			totalOK++
		}
	}

	log.Printf("worker nfe sync: tenants=%d sincronizacoes_ok=%d erros=%d", len(tenants), totalOK, totalErr)

	status := domain.MonitorOperacaoStatusSucesso
	msg := fmt.Sprintf("ok=%d erros=%d tenants=%d", totalOK, totalErr, len(tenants))
	if totalErr > 0 {
		status = domain.MonitorOperacaoStatusErro
	}
	w.recordMonitor(ctx, status, msg, map[string]any{
		"sincronizacoes_ok": totalOK,
		"erros":             totalErr,
		"tenants_catalogo":  len(tenants),
	})
	return nil
}

func (w *NFESyncWorker) recordMonitor(ctx context.Context, status, msg string, det map[string]any) {
	if w.monitor == nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	m := msg
	_, err := w.monitor.Insert(ctx, repository.MonitorOperacaoInsert{
		TenantID: domain.MonitorOperacaoTenantPlataformaID,
		UserID:   nil,
		Origem:   domain.MonitorOperacaoOrigemAutomatico,
		Tipo:     domain.MonitorOperacaoTipoWorkerNFESyncProvider,
		Status:   status,
		Mensagem: &m,
		Detalhe:  det,
	})
	if err != nil {
		return
	}
}
