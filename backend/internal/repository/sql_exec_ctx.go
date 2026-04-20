package repository

import (
	"context"

	"github.com/chayimamaral/vecontab/backend/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func dbQuery(ctx context.Context, pool *pgxpool.Pool, sql string, args ...any) (pgx.Rows, error) {
	if conn := db.ConnFromContext(ctx); conn != nil {
		return conn.Query(ctx, sql, args...)
	}
	return pool.Query(ctx, sql, args...)
}

func dbExec(ctx context.Context, pool *pgxpool.Pool, sql string, args ...any) (pgconn.CommandTag, error) {
	if conn := db.ConnFromContext(ctx); conn != nil {
		return conn.Exec(ctx, sql, args...)
	}
	return pool.Exec(ctx, sql, args...)
}

func dbQueryRow(ctx context.Context, pool *pgxpool.Pool, sql string, args ...any) pgx.Row {
	if conn := db.ConnFromContext(ctx); conn != nil {
		return conn.QueryRow(ctx, sql, args...)
	}
	return pool.QueryRow(ctx, sql, args...)
}

func dbBegin(ctx context.Context, pool *pgxpool.Pool) (pgx.Tx, error) {
	if conn := db.ConnFromContext(ctx); conn != nil {
		return conn.Begin(ctx)
	}
	return pool.Begin(ctx)
}

func dbBeginTx(ctx context.Context, pool *pgxpool.Pool, opts pgx.TxOptions) (pgx.Tx, error) {
	if conn := db.ConnFromContext(ctx); conn != nil {
		return conn.BeginTx(ctx, opts)
	}
	return pool.BeginTx(ctx, opts)
}
