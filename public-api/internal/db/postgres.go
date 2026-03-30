package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/chayimamaral/vecontab/public-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse pg config: %w", err)
	}

	if cfg.SSLRootCertPath != "" && !cfg.SSLInsecure {
		certBytes, certErr := os.ReadFile(cfg.SSLRootCertPath)
		if certErr == nil {
			rootCAs := x509.NewCertPool()
			rootCAs.AppendCertsFromPEM(certBytes)
			poolConfig.ConnConfig.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				RootCAs:    rootCAs,
			}
		}
	} else {
		poolConfig.ConnConfig.TLSConfig = nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pg pool: %w", err)
		log.Printf("Erro ao criar pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping pg: %w", err)
	}

	return pool, nil
}
