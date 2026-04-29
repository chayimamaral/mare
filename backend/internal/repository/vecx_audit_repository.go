package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VecxAuditRepository acessa apenas o banco VECX_AUDIT (schema vecx_audit).
type VecxAuditRepository struct {
	pool *pgxpool.Pool
}

func NewVecxAuditRepository(pool *pgxpool.Pool) *VecxAuditRepository {
	return &VecxAuditRepository{pool: pool}
}

func (r *VecxAuditRepository) Pool() *pgxpool.Pool {
	return r.pool
}

// EnsureSchema cria schema e tabela active_sessions se ainda nao existirem (EF-929).
func (r *VecxAuditRepository) EnsureSchema(ctx context.Context) error {
	const ddl = `
CREATE SCHEMA IF NOT EXISTS vecx_audit;
CREATE TABLE IF NOT EXISTS vecx_audit.active_sessions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	user_email text NOT NULL,
	tenant_id uuid NOT NULL,
	tenant_name text NOT NULL DEFAULT '',
	tenant_cnpj text NOT NULL DEFAULT '',
	tenant_contato text NOT NULL DEFAULT '',
	logged_at timestamptz NOT NULL DEFAULT now(),
	ip_address text NOT NULL DEFAULT '',
	active boolean NOT NULL DEFAULT true,
	CONSTRAINT active_sessions_user_id_key UNIQUE (user_id)
);
CREATE INDEX IF NOT EXISTS idx_active_sessions_active_logged ON vecx_audit.active_sessions (active, logged_at DESC);
`
	if _, err := r.pool.Exec(ctx, ddl); err != nil {
		return fmt.Errorf("vecx_audit ensure schema: %w", err)
	}
	return nil
}

// UpsertActiveSession grava ou atualiza a sessao ativa do usuario (login / troca de tenant).
func (r *VecxAuditRepository) UpsertActiveSession(ctx context.Context, userID, userEmail, tenantID, tenantName, tenantCNPJ, tenantContato, ipAddress string) error {
	const q = `
INSERT INTO vecx_audit.active_sessions (
	user_id, user_email, tenant_id, tenant_name, tenant_cnpj, tenant_contato, logged_at, ip_address, active
) VALUES ($1::uuid, $2, $3::uuid, $4, $5, $6, now(), $7, true)
ON CONFLICT (user_id) DO UPDATE SET
	user_email = EXCLUDED.user_email,
	tenant_id = EXCLUDED.tenant_id,
	tenant_name = EXCLUDED.tenant_name,
	tenant_cnpj = EXCLUDED.tenant_cnpj,
	tenant_contato = EXCLUDED.tenant_contato,
	logged_at = now(),
	ip_address = EXCLUDED.ip_address,
	active = true`
	_, err := r.pool.Exec(ctx, q,
		strings.TrimSpace(userID),
		strings.TrimSpace(userEmail),
		strings.TrimSpace(tenantID),
		strings.TrimSpace(tenantName),
		strings.TrimSpace(tenantCNPJ),
		strings.TrimSpace(tenantContato),
		strings.TrimSpace(ipAddress),
	)
	if err != nil {
		return fmt.Errorf("vecx_audit upsert session: %w", err)
	}
	return nil
}

// DeactivateUserSession marca a sessao como inativa (logout).
func (r *VecxAuditRepository) DeactivateUserSession(ctx context.Context, userID string) error {
	const q = `UPDATE vecx_audit.active_sessions SET active = false WHERE user_id = $1::uuid`
	if _, err := r.pool.Exec(ctx, q, strings.TrimSpace(userID)); err != nil {
		return fmt.Errorf("vecx_audit deactivate session: %w", err)
	}
	return nil
}

// ListActiveSessions retorna sessoes marcadas ativas (monitoramento SUPER).
func (r *VecxAuditRepository) ListActiveSessions(ctx context.Context) ([]domain.ActiveSessionRow, error) {
	const q = `
SELECT id, user_id, user_email, tenant_id, tenant_name, tenant_cnpj, tenant_contato, logged_at, ip_address, active
FROM vecx_audit.active_sessions
WHERE active = true
ORDER BY logged_at DESC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("vecx_audit list active: %w", err)
	}
	defer rows.Close()

	out := make([]domain.ActiveSessionRow, 0)
	for rows.Next() {
		var row domain.ActiveSessionRow
		if err := rows.Scan(
			&row.ID,
			&row.UserID,
			&row.UserEmail,
			&row.TenantID,
			&row.TenantName,
			&row.TenantCNPJ,
			&row.TenantContato,
			&row.LoggedAt,
			&row.IPAddress,
			&row.Active,
		); err != nil {
			return nil, fmt.Errorf("vecx_audit scan: %w", err)
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

// Ping verifica conectividade do banco de auditoria.
func (r *VecxAuditRepository) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return r.pool.Ping(ctx)
}
