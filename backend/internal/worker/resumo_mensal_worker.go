package worker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/chayimamaral/vecx/backend/internal/config"
	"github.com/chayimamaral/vecx/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type ResumoMensalMailer interface {
	SendResumoMensal(ctx context.Context, tenantID, clienteNome, email, assunto, corpo string, valorTotal float64) error
}

type LogResumoMensalMailer struct{}

func (LogResumoMensalMailer) SendResumoMensal(ctx context.Context, tenantID, clienteNome, email, assunto, corpo string, valorTotal float64) error {
	log.Printf("worker resumo mensal [MOCK EMAIL]: tenant=%s cliente=%s email=%s assunto=%q valor_total=%.2f corpo=%q", tenantID, clienteNome, email, assunto, valorTotal, corpo)
	return nil
}

type ResumoMensalWorker struct {
	cfg     config.Config
	pool    *pgxpool.Pool
	repo    *repository.ResumoMensalRepository
	mailer  ResumoMensalMailer
	cron    *cron.Cron
	running int32
}

const resumoMensalAdvisoryLockKey int64 = 9482217702

func NewResumoMensalWorker(pool *pgxpool.Pool, cfg config.Config, mailer ResumoMensalMailer) (*ResumoMensalWorker, error) {
	loc, err := time.LoadLocation(cfg.ResumoMensalWorkerTimezone)
	if err != nil {
		return nil, fmt.Errorf("timezone inválida do worker resumo mensal: %w", err)
	}
	if mailer == nil {
		mailer = LogResumoMensalMailer{}
	}
	c := cron.New(
		cron.WithLocation(loc),
		cron.WithParser(cron.NewParser(cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)
	w := &ResumoMensalWorker{
		cfg:    cfg,
		pool:   pool,
		repo:   repository.NewResumoMensalRepository(pool),
		mailer: mailer,
		cron:   c,
	}

	if _, err := c.AddFunc(cfg.ResumoMensalWorkerCron, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := w.runOnce(ctx); err != nil {
			log.Printf("worker resumo mensal: erro na execução agendada: %v", err)
		}
	}); err != nil {
		return nil, fmt.Errorf("cron inválido do worker resumo mensal: %w", err)
	}

	return w, nil
}

func (w *ResumoMensalWorker) Start(ctx context.Context) {
	log.Printf("worker resumo mensal: habilitado cron=%q tz=%q dia_envio=%d", w.cfg.ResumoMensalWorkerCron, w.cfg.ResumoMensalWorkerTimezone, w.cfg.ResumoDiaEnvio)
	w.cron.Start()
	if w.cfg.ResumoMensalWorkerRunOnStartup {
		go func() {
			startupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			if err := w.runOnce(startupCtx); err != nil {
				log.Printf("worker resumo mensal: erro no run_on_startup: %v", err)
			}
		}()
	}
	<-ctx.Done()
	stopCtx := w.cron.Stop()
	select {
	case <-stopCtx.Done():
	case <-time.After(10 * time.Second):
	}
	log.Printf("worker resumo mensal: finalizado")
}

func (w *ResumoMensalWorker) runOnce(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&w.running, 0, 1) {
		return nil
	}
	defer atomic.StoreInt32(&w.running, 0)

	conn, err := w.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn worker resumo mensal: %w", err)
	}
	defer conn.Release()

	var gotLock bool
	if err := conn.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, resumoMensalAdvisoryLockKey).Scan(&gotLock); err != nil {
		return fmt.Errorf("obter advisory lock resumo mensal: %w", err)
	}
	if !gotLock {
		log.Printf("worker resumo mensal: execução já em andamento em outra réplica")
		return nil
	}
	defer func() {
		_, _ = conn.Exec(context.Background(), `SELECT pg_advisory_unlock($1)`, resumoMensalAdvisoryLockKey)
	}()

	loc, _ := time.LoadLocation(w.cfg.ResumoMensalWorkerTimezone)
	if loc == nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	if w.cfg.ResumoDiaEnvio > 0 && now.Day() != w.cfg.ResumoDiaEnvio {
		log.Printf("worker resumo mensal: ignorado dia_atual=%d dia_envio=%d", now.Day(), w.cfg.ResumoDiaEnvio)
		return nil
	}

	base := time.Now().In(loc)
	inicioMesAtual := time.Date(base.Year(), base.Month(), 1, 0, 0, 0, 0, loc)
	inicioMesAnterior := inicioMesAtual.AddDate(0, -1, 0)
	fimMesAnterior := inicioMesAtual.AddDate(0, 0, -1)

	tenants, err := w.repo.ListTenants(ctx)
	if err != nil {
		return err
	}

	var totalMensagens int
	var totalTenantsAtivos int
	for _, tenant := range tenants {
		ativo, err := w.repo.IsTenantResumoMensalAtivo(ctx, tenant.TenantID, tenant.SchemaName)
		if err != nil {
			log.Printf("worker resumo mensal: tenant=%s erro ao validar flag: %v", tenant.TenantID, err)
			continue
		}
		if !ativo {
			continue
		}
		totalTenantsAtivos++

		resumos, err := w.repo.ListResumoMensalPorCliente(ctx, tenant.SchemaName, tenant.TenantID, inicioMesAnterior, inicioMesAtual)
		if err != nil {
			log.Printf("worker resumo mensal: tenant=%s erro ao consultar resumo por cliente: %v", tenant.TenantID, err)
			continue
		}
		for _, resumo := range resumos {
			assunto := fmt.Sprintf("Resumo Mensal de Notas Fiscais - %s", inicioMesAnterior.Format("01/2006"))
			corpo := fmt.Sprintf(
				"Prezado cliente, este é o resumo das suas notas fiscais recebidas no período de %s a %s. O somatório total foi de: R$ %.2f.",
				inicioMesAnterior.Format("02/01/2006"),
				fimMesAnterior.Format("02/01/2006"),
				resumo.ValorTotal,
			)
			if err := w.mailer.SendResumoMensal(ctx, tenant.TenantID, resumo.ClienteNome, strings.TrimSpace(resumo.EmailContato), assunto, corpo, resumo.ValorTotal); err != nil {
				log.Printf("worker resumo mensal: tenant=%s cliente=%s email=%s erro ao enviar resumo: %v", tenant.TenantID, resumo.ClienteNome, resumo.EmailContato, err)
				continue
			}
			totalMensagens++
		}
		log.Printf("worker resumo mensal: tenant=%s enviados=%d periodo=%s..%s", tenant.TenantID, len(resumos), inicioMesAnterior.Format("2006-01-02"), fimMesAnterior.Format("2006-01-02"))
	}

	log.Printf("worker resumo mensal: tenants_ativos=%d mensagens_enviadas=%d periodo=%s..%s", totalTenantsAtivos, totalMensagens, inicioMesAnterior.Format("2006-01-02"), fimMesAnterior.Format("2006-01-02"))
	return nil
}
