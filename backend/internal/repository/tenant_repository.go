package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TenantEntity struct {
	ID      string `json:"id"`
	Nome    string `json:"nome"`
	Contato string `json:"contato"`
	Active  bool   `json:"active"`
	Plano   string `json:"plano"`
}

type TenantRepository struct {
	pool *pgxpool.Pool
}

func NewTenantRepository(pool *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{pool: pool}
}

func (r *TenantRepository) Create(ctx context.Context, nome, contato string) (TenantEntity, error) {
	const existsQuery = `SELECT count(*) FROM public.tenant WHERE nome = $1`
	var count int64
	if err := r.pool.QueryRow(ctx, existsQuery, nome).Scan(&count); err != nil {
		return TenantEntity{}, fmt.Errorf("check tenant exists: %w", err)
	}

	if count > 0 {
		return TenantEntity{}, fmt.Errorf("Empresa ja cadastrada")
	}

	const query = `
		INSERT INTO public.tenant (nome, contato, active, plano)
		VALUES ($1, $2, $3, $4)
		RETURNING id, nome, contato, active, plano`

	var tenant TenantEntity
	if err := r.pool.QueryRow(ctx, query, nome, contato, true, "DEMO").Scan(
		&tenant.ID,
		&tenant.Nome,
		&tenant.Contato,
		&tenant.Active,
		&tenant.Plano,
	); err != nil {
		return TenantEntity{}, fmt.Errorf("create tenant: %w", err)
	}

	return tenant, nil
}

func (r *TenantRepository) Detail(ctx context.Context, id string) (TenantEntity, error) {
	const query = `SELECT id, nome, contato, active, plano FROM public.tenant WHERE id = $1`

	var tenant TenantEntity
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Nome,
		&tenant.Contato,
		&tenant.Active,
		&tenant.Plano,
	); err != nil {
		return TenantEntity{}, fmt.Errorf("detail tenant: %w", err)
	}

	return tenant, nil
}

func (r *TenantRepository) Update(ctx context.Context, id, nome, contato string, active bool) (TenantEntity, error) {
	const query = `
		UPDATE public.tenant
		SET nome = $1, active = $2, contato = $3
		WHERE id = $4
		RETURNING id, nome, contato, active, plano`

	var tenant TenantEntity
	if err := r.pool.QueryRow(ctx, query, nome, active, contato, id).Scan(
		&tenant.ID,
		&tenant.Nome,
		&tenant.Contato,
		&tenant.Active,
		&tenant.Plano,
	); err != nil {
		return TenantEntity{}, fmt.Errorf("update tenant: %w", err)
	}

	return tenant, nil
}

func (r *TenantRepository) List(ctx context.Context, role, tenantID string) ([]TenantEntity, error) {
	query := `SELECT id, nome, contato, active, plano FROM public.tenant WHERE id = $1`
	args := []any{tenantID}

	if role == "SUPER" {
		query = `SELECT id, nome, contato, active, plano FROM public.tenant`
		args = []any{}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %w", err)
	}
	defer rows.Close()

	tenants := make([]TenantEntity, 0)
	for rows.Next() {
		var tenant TenantEntity
		if err := rows.Scan(
			&tenant.ID,
			&tenant.Nome,
			&tenant.Contato,
			&tenant.Active,
			&tenant.Plano,
		); err != nil {
			return nil, fmt.Errorf("scan tenant: %w", err)
		}

		tenants = append(tenants, tenant)
	}

	return tenants, nil
}
