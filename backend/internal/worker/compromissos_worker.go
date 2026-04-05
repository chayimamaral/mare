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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type CompromissosWorker struct {
	pool    *pgxpool.Pool
	cfg     config.Config
	cron    *cron.Cron
	running int32
	monitor *repository.MonitorOperacaoRepository
}

const compromissosAdvisoryLockKey int64 = 9482217701

func NewCompromissosWorker(pool *pgxpool.Pool, cfg config.Config, monitor *repository.MonitorOperacaoRepository) (*CompromissosWorker, error) {
	loc, err := time.LoadLocation(cfg.CompromissosWorkerTimezone)
	if err != nil {
		return nil, fmt.Errorf("timezone inválida do worker: %w", err)
	}

	c := cron.New(
		cron.WithLocation(loc),
		cron.WithParser(cron.NewParser(cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	w := &CompromissosWorker{
		pool:    pool,
		cfg:     cfg,
		cron:    c,
		monitor: monitor,
	}

	if _, err := c.AddFunc(cfg.CompromissosWorkerCron, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := w.runOnce(ctx); err != nil {
			log.Printf("worker compromissos: erro na execução agendada: %v", err)
		}
	}); err != nil {
		return nil, fmt.Errorf("cron inválido do worker: %w", err)
	}

	return w, nil
}

func (w *CompromissosWorker) Start(ctx context.Context) {
	log.Printf("worker compromissos: habilitado cron=%q tz=%q", w.cfg.CompromissosWorkerCron, w.cfg.CompromissosWorkerTimezone)
	w.cron.Start()

	if w.cfg.CompromissosWorkerRunOnStartup {
		go func() {
			startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			if err := w.runOnce(startupCtx); err != nil {
				log.Printf("worker compromissos: erro no run_on_startup: %v", err)
			}
		}()
	}

	<-ctx.Done()
	stopCtx := w.cron.Stop()
	select {
	case <-stopCtx.Done():
	case <-time.After(10 * time.Second):
	}
	log.Printf("worker compromissos: finalizado")
}

func (w *CompromissosWorker) runOnce(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&w.running, 0, 1) {
		return nil
	}
	defer atomic.StoreInt32(&w.running, 0)

	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx worker: %w", err)
	}
	defer tx.Rollback(ctx)

	var gotLock bool
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock($1)`, compromissosAdvisoryLockKey).Scan(&gotLock); err != nil {
		return fmt.Errorf("obter advisory lock: %w", err)
	}
	if !gotLock {
		log.Printf("worker compromissos: execução já em andamento em outra réplica")
		return nil
	}

	refDate := time.Now().AddDate(0, 1, 0)
	refMonth := time.Date(refDate.Year(), refDate.Month(), 1, 0, 0, 0, 0, refDate.Location())

	var totalInseridos int
	if err := tx.QueryRow(
		ctx,
		`SELECT public.gerar_compromissos_geral($1::date)`,
		refMonth.Format("2006-01-02"),
	).Scan(&totalInseridos); err != nil {
		w.recordMonitorOperacao(context.Background(), domain.MonitorOperacaoStatusErro, err.Error(), map[string]any{
			"competencia": refMonth.Format("2006-01-02"),
			"fase":        "gerar_compromissos_geral",
		})
		return fmt.Errorf("executar gerar_compromissos_geral: %w", err)
	}

	log.Printf("worker compromissos: competencia=%s inseridos=%d", refMonth.Format("2006-01-02"), totalInseridos)
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx worker: %w", err)
	}

	w.recordMonitorOperacao(context.Background(), domain.MonitorOperacaoStatusSucesso, fmt.Sprintf("inseridos=%d", totalInseridos), map[string]any{
		"competencia": refMonth.Format("2006-01-02"),
		"inseridos":   totalInseridos,
	})
	return nil
}

func (w *CompromissosWorker) recordMonitorOperacao(ctx context.Context, status, msg string, det map[string]any) {
	if w.monitor == nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	m := msg
	_ = w.monitor.Insert(ctx, repository.MonitorOperacaoInsert{
		TenantID: domain.MonitorOperacaoTenantPlataformaID,
		UserID:   nil,
		Origem:   domain.MonitorOperacaoOrigemAutomatico,
		Tipo:     domain.MonitorOperacaoTipoWorkerCompromissosMensal,
		Status:   status,
		Mensagem: &m,
		Detalhe:  det,
	})
}
