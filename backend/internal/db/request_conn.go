package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type connCtxKey string

const tenantConnKey connCtxKey = "tenantConn"

func ContextWithConn(ctx context.Context, conn *pgxpool.Conn) context.Context {
	return context.WithValue(ctx, tenantConnKey, conn)
}

func ConnFromContext(ctx context.Context) *pgxpool.Conn {
	c, _ := ctx.Value(tenantConnKey).(*pgxpool.Conn)
	return c
}

func AcquireTenantConn(ctx context.Context, pool *pgxpool.Pool, tenantSchema string) (*pgxpool.Conn, error) {
	tenantSchema = strings.TrimSpace(tenantSchema)
	if tenantSchema == "" {
		return nil, fmt.Errorf("tenant schema obrigatorio")
	}
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire conn: %w", err)
	}
	if _, err := conn.Exec(ctx, fmt.Sprintf("SET search_path TO %s, public", quoteIdent(tenantSchema))); err != nil {
		conn.Release()
		return nil, fmt.Errorf("set search_path: %w", err)
	}
	return conn, nil
}
