package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FeatureMatrixRepository struct {
	pool *pgxpool.Pool
}

func NewFeatureMatrixRepository(pool *pgxpool.Pool) *FeatureMatrixRepository {
	return &FeatureMatrixRepository{pool: pool}
}

func (r *FeatureMatrixRepository) AllSlugs(ctx context.Context) ([]string, error) {
	const q = `SELECT slug FROM public.modulo_plataforma ORDER BY ordem ASC, slug ASC`
	rows, err := dbQuery(ctx, r.pool, q)
	if err != nil {
		return nil, fmt.Errorf("list modulos: %w", err)
	}
	defer rows.Close()

	out := make([]string, 0)
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("scan slug: %w", err)
		}
		out = append(out, strings.TrimSpace(s))
	}
	return out, nil
}

func (r *FeatureMatrixRepository) SlugsForRepresentante(ctx context.Context, representanteID string) ([]string, error) {
	rid := strings.TrimSpace(representanteID)
	if rid == "" {
		return nil, nil
	}
	const q = `
		SELECT m.slug
		FROM public.matriz_acesso ma
		JOIN public.modulo_plataforma m ON m.id = ma.modulo_id
		WHERE ma.representante_id = $1::uuid AND ma.habilitado = true
		ORDER BY m.ordem ASC, m.slug ASC`
	rows, err := dbQuery(ctx, r.pool, q, rid)
	if err != nil {
		return nil, fmt.Errorf("matriz representante: %w", err)
	}
	defer rows.Close()

	out := make([]string, 0)
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("scan slug: %w", err)
		}
		out = append(out, strings.TrimSpace(s))
	}
	return out, nil
}

func (r *FeatureMatrixRepository) TenantRepresentativeID(ctx context.Context, tenantID string) (string, bool, error) {
	const q = `SELECT representative_id::text FROM public.tenant WHERE id = $1::uuid`
	var raw sql.NullString
	if err := dbQueryRow(ctx, r.pool, q, tenantID).Scan(&raw); err != nil {
		return "", false, fmt.Errorf("tenant representative: %w", err)
	}
	if !raw.Valid || strings.TrimSpace(raw.String) == "" {
		return "", false, nil
	}
	return strings.TrimSpace(raw.String), true, nil
}

func (r *FeatureMatrixRepository) TenantLinkedToRepresentante(ctx context.Context, tenantID, representanteID string) (bool, error) {
	tid := strings.TrimSpace(tenantID)
	rid := strings.TrimSpace(representanteID)
	if tid == "" || rid == "" {
		return false, nil
	}
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM public.tenant t
			WHERE t.id = $1::uuid AND t.representative_id = $2::uuid
		)`
	var ok bool
	if err := dbQueryRow(ctx, r.pool, q, tid, rid).Scan(&ok); err != nil {
		return false, fmt.Errorf("vinculo tenant representante: %w", err)
	}
	return ok, nil
}

func (r *FeatureMatrixRepository) ResolveForUser(ctx context.Context, roleUpper, tenantID, representanteID string) ([]string, error) {
	roleUpper = strings.ToUpper(strings.TrimSpace(roleUpper))
	switch roleUpper {
	case "SUPER":
		return r.AllSlugs(ctx)
	case "REPRESENTANTE":
		return r.SlugsForRepresentante(ctx, representanteID)
	default:
		rid, has, err := r.TenantRepresentativeID(ctx, tenantID)
		if err != nil {
			return nil, err
		}
		if !has {
			return r.AllSlugs(ctx)
		}
		return r.SlugsForRepresentante(ctx, rid)
	}
}
