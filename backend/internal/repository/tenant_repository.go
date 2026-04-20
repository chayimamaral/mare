package repository

import (
	"context"
	"fmt"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TenantRepository struct {
	pool *pgxpool.Pool
}

func NewTenantRepository(pool *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{pool: pool}
}

func (r *TenantRepository) Create(ctx context.Context, nome, contato, plano string) (domain.TenantEntity, error) {
	const existsQuery = `SELECT count(*) FROM public.tenant WHERE nome = $1`
	var count int64
	if err := dbQueryRow(ctx, r.pool, existsQuery, nome).Scan(&count); err != nil {
		return domain.TenantEntity{}, fmt.Errorf("check tenant exists: %w", err)
	}

	if count > 0 {
		return domain.TenantEntity{}, fmt.Errorf("Empresa ja cadastrada")
	}

	tx, err := dbBeginTx(ctx, r.pool, pgx.TxOptions{})
	if err != nil {
		return domain.TenantEntity{}, fmt.Errorf("begin tx create tenant: %w", err)
	}
	defer tx.Rollback(ctx)

	schemaSlug, err := buildUniqueSchemaSlug(ctx, tx, nome)
	if err != nil {
		return domain.TenantEntity{}, err
	}

	const query = `
		INSERT INTO public.tenant (nome, contato, active, plano)
		VALUES ($1, $2, $3, $4::public.plano)
		RETURNING id, nome, contato, active, COALESCE(plano::text, '')`

	var tenant domain.TenantEntity
	if err := tx.QueryRow(ctx, query, nome, contato, true, plano).Scan(
		&tenant.ID,
		&tenant.Nome,
		&tenant.Contato,
		&tenant.Active,
		&tenant.Plano,
	); err != nil {
		return domain.TenantEntity{}, fmt.Errorf("create tenant: %w", err)
	}

	if err := tx.QueryRow(
		ctx,
		`SELECT public.provision_tenant_schema($1::uuid, $2::text, NULL)`,
		tenant.ID,
		schemaSlug,
	).Scan(&tenant.SchemaName); err != nil {
		return domain.TenantEntity{}, fmt.Errorf("provision tenant schema: %w", err)
	}

	if err := setTxTenantSearchPath(ctx, tx, tenant.SchemaName); err != nil {
		return domain.TenantEntity{}, fmt.Errorf("set tenant search_path: %w", err)
	}

	const dadosQuery = `INSERT INTO tenant_dados (tenantid) VALUES ($1) ON CONFLICT DO NOTHING`
	if _, err := tx.Exec(ctx, dadosQuery, tenant.ID); err != nil {
		return domain.TenantEntity{}, fmt.Errorf("create tenant_dados local: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.TenantEntity{}, fmt.Errorf("commit create tenant: %w", err)
	}

	return tenant, nil
}

func (r *TenantRepository) Detail(ctx context.Context, id string) (domain.TenantEntity, error) {
	const query = `
		SELECT id,
		       COALESCE(nome, ''),
		       COALESCE(contato, ''),
		       COALESCE(active, false),
		       COALESCE(plano::text, '')
		FROM public.tenant
		WHERE id = $1::uuid`

	var tenant domain.TenantEntity
	if err := dbQueryRow(ctx, r.pool, query, id).Scan(
		&tenant.ID,
		&tenant.Nome,
		&tenant.Contato,
		&tenant.Active,
		&tenant.Plano,
	); err != nil {
		return domain.TenantEntity{}, fmt.Errorf("detail tenant: %w", err)
	}

	return tenant, nil
}

func (r *TenantRepository) Update(ctx context.Context, id, nome, contato, plano string, active bool) (domain.TenantEntity, error) {
	const query = `
		UPDATE public.tenant
		SET nome = $1,
		    active = $2,
		    contato = $3,
		    plano = CASE WHEN BTRIM($4) = '' THEN plano ELSE $4::public.plano END
		WHERE id = $5::uuid
		RETURNING id, COALESCE(nome, ''), COALESCE(contato, ''), COALESCE(active, false), COALESCE(plano::text, '')`

	var tenant domain.TenantEntity
	if err := dbQueryRow(ctx, r.pool, query, nome, active, contato, plano, id).Scan(
		&tenant.ID,
		&tenant.Nome,
		&tenant.Contato,
		&tenant.Active,
		&tenant.Plano,
	); err != nil {
		return domain.TenantEntity{}, fmt.Errorf("update tenant: %w", err)
	}

	return tenant, nil
}

func (r *TenantRepository) List(ctx context.Context, role, tenantID string) ([]domain.TenantEntity, error) {
	sqlQuery := `
		SELECT id,
		       COALESCE(nome, ''),
		       COALESCE(contato, ''),
		       COALESCE(active, false),
		       COALESCE(plano::text, '')
		FROM public.tenant
		WHERE id = $1::uuid`
	args := []any{tenantID}

	if role == "SUPER" {
		sqlQuery = `
			SELECT id,
			       COALESCE(nome, ''),
			       COALESCE(contato, ''),
			       COALESCE(active, false),
			       COALESCE(plano::text, '')
			FROM public.tenant
			WHERE NULLIF(BTRIM(COALESCE(nome, '')), '') IS NOT NULL`
		args = []any{}
	}

	rows, err := dbQuery(ctx, r.pool, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %w", err)
	}
	defer rows.Close()

	tenants := make([]domain.TenantEntity, 0)
	for rows.Next() {
		var tenant domain.TenantEntity
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

func (r *TenantRepository) ListWithDadosForSuper(ctx context.Context) ([]domain.TenantListRow, error) {
	const sqlQuery = `
		SELECT t.id,
		       COALESCE(t.nome, ''),
		       COALESCE(t.contato, ''),
		       COALESCE(t.active, false),
		       COALESCE(t.plano::text, ''),
		       COALESCE(tsc.schema_name, '')
		FROM public.tenant t
		LEFT JOIN public.tenant_schema_catalog tsc ON tsc.tenant_id = t.id
		WHERE NULLIF(BTRIM(COALESCE(t.nome, '')), '') IS NOT NULL
		ORDER BY COALESCE(t.nome, '')`

	rows, err := dbQuery(ctx, r.pool, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("list tenants with dados: %w", err)
	}
	defer rows.Close()

	out := make([]domain.TenantListRow, 0)
	for rows.Next() {
		var row domain.TenantListRow
		var schemaName string
		if err := rows.Scan(
			&row.ID,
			&row.Nome,
			&row.Contato,
			&row.Active,
			&row.Plano,
			&schemaName,
		); err != nil {
			return nil, fmt.Errorf("scan tenant list row: %w", err)
		}
		if schemaName != "" {
			if err := withTenantSchemaContext(ctx, r.pool, row.ID, func(inner context.Context) error {
				return dbQueryRow(inner, r.pool, `
					SELECT COALESCE(cnpj::text, ''), COALESCE(razaosocial::text, ''), COALESCE(fantasia::text, '')
					FROM tenant_dados
					WHERE tenantid = $1::uuid
					LIMIT 1`, row.ID).Scan(&row.CNPJ, &row.RazaoSocial, &row.Fantasia)
			}); err != nil && err != pgx.ErrNoRows {
				return nil, fmt.Errorf("load tenant_dados local: %w", err)
			}
		}
		out = append(out, row)
	}

	return out, nil
}
