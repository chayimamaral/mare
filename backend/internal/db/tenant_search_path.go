package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ResolveTenantSchema retorna o schema_name do tenant no catálogo central.
func ResolveTenantSchema(ctx context.Context, pool *pgxpool.Pool, tenantID string) (string, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return "", fmt.Errorf("tenant_id obrigatorio")
	}
	var schemaName string
	err := pool.QueryRow(
		ctx,
		`SELECT schema_name FROM public.tenant_schema_catalog WHERE tenant_id = $1::uuid`,
		tenantID,
	).Scan(&schemaName)
	if err != nil {
		return "", fmt.Errorf("resolve tenant schema: %w", err)
	}
	return strings.TrimSpace(schemaName), nil
}

// WithTenantConnSearchPath executa fn em uma conexao com search_path do tenant + public.
// Uso incremental em repositorios/servicos durante o refactoring schema-per-tenant.
func WithTenantConnSearchPath(
	ctx context.Context,
	pool *pgxpool.Pool,
	tenantSchema string,
	fn func(context.Context, *pgxpool.Conn) error,
) error {
	tenantSchema = strings.TrimSpace(tenantSchema)
	if tenantSchema == "" {
		return fmt.Errorf("tenant schema obrigatorio")
	}
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn tenant schema: %w", err)
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, fmt.Sprintf("SET search_path TO %s, public", quoteIdent(tenantSchema))); err != nil {
		return fmt.Errorf("set search_path: %w", err)
	}
	return fn(ctx, conn)
}

func quoteIdent(ident string) string {
	return `"` + strings.ReplaceAll(ident, `"`, `""`) + `"`
}
