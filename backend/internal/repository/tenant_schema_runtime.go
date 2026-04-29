package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var repositoryTenantSchemaRegex = regexp.MustCompile(`^[a-z][a-z0-9_]{2,62}$`)

func normalizeRepositoryTenantSchema(schemaName string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(schemaName))
	if !repositoryTenantSchemaRegex.MatchString(normalized) {
		return "", fmt.Errorf("tenant schema invalido")
	}
	return normalized, nil
}

func quoteRepositoryIdent(ident string) string {
	return `"` + strings.ReplaceAll(ident, `"`, `""`) + `"`
}

func setTxTenantSearchPath(ctx context.Context, tx pgx.Tx, schemaName string) error {
	normalized, err := normalizeRepositoryTenantSchema(schemaName)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL search_path TO %s, public", quoteRepositoryIdent(normalized))); err != nil {
		return fmt.Errorf("set local search_path: %w", err)
	}
	return nil
}

func withTenantSchemaContext(ctx context.Context, pool *pgxpool.Pool, tenantID string, fn func(context.Context) error) error {
	schemaName, err := repositoryResolveTenantSchema(ctx, pool, tenantID)
	if err != nil {
		return err
	}
	return withTenantSchemaContextByName(ctx, pool, schemaName, fn)
}

func withTenantSchemaContextByName(ctx context.Context, pool *pgxpool.Pool, schemaName string, fn func(context.Context) error) error {
	normalized, err := normalizeRepositoryTenantSchema(schemaName)
	if err != nil {
		return err
	}
	return db.WithTenantConnSearchPath(ctx, pool, normalized, func(inner context.Context, conn *pgxpool.Conn) error {
		return fn(db.ContextWithConn(inner, conn))
	})
}

func repositoryResolveTenantSchema(ctx context.Context, pool *pgxpool.Pool, tenantID string) (string, error) {
	return db.ResolveTenantSchema(ctx, pool, tenantID)
}
