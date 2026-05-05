package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ResumoMensalRepository struct {
	pool *pgxpool.Pool
}

type ResumoMensalTenant struct {
	TenantID   string
	SchemaName string
}

type ResumoMensalCliente struct {
	ClienteID    string
	ClienteNome  string
	EmailContato string
	ValorTotal   float64
}

func NewResumoMensalRepository(pool *pgxpool.Pool) *ResumoMensalRepository {
	return &ResumoMensalRepository{pool: pool}
}

func (r *ResumoMensalRepository) ListTenants(ctx context.Context) ([]ResumoMensalTenant, error) {
	rows, err := dbQuery(ctx, r.pool, `
		SELECT tenant_id::text, schema_name
		FROM public.tenant_schema_catalog
		WHERE NULLIF(BTRIM(schema_name), '') IS NOT NULL
		ORDER BY tenant_id`)
	if err != nil {
		return nil, fmt.Errorf("listar tenants resumo mensal: %w", err)
	}
	defer rows.Close()

	out := make([]ResumoMensalTenant, 0)
	for rows.Next() {
		var item ResumoMensalTenant
		if err := rows.Scan(&item.TenantID, &item.SchemaName); err != nil {
			return nil, fmt.Errorf("scan tenant resumo mensal: %w", err)
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *ResumoMensalRepository) IsTenantResumoMensalAtivo(ctx context.Context, tenantID, schemaName string) (bool, error) {
	var ativo bool
	err := withTenantSchemaContextByName(ctx, r.pool, schemaName, func(inner context.Context) error {
		err := dbQueryRow(inner, r.pool, `
			SELECT COALESCE(enviar_resumo_mensal, false)
			FROM tenant_configuracoes
			WHERE tenant_id = $1::uuid
			LIMIT 1
		`, strings.TrimSpace(tenantID)).Scan(&ativo)
		if err != nil {
			ativo = false
			return nil
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("verificar flag enviar_resumo_mensal: %w", err)
	}
	return ativo, nil
}

func (r *ResumoMensalRepository) ListResumoMensalPorCliente(ctx context.Context, schemaName, tenantID string, inicio, fim time.Time) ([]ResumoMensalCliente, error) {
	out := make([]ResumoMensalCliente, 0)
	err := withTenantSchemaContextByName(ctx, r.pool, schemaName, func(inner context.Context) error {
		rows, err := dbQuery(inner, r.pool, `
			SELECT
				c.id::text AS cliente_id,
				COALESCE(c.nome, '') AS cliente_nome,
				COALESCE(cd.email_contato, '') AS email_contato,
				COALESCE(SUM(COALESCE(n.valor_total, 0)), 0)::float8 AS valor_total
			FROM cliente c
			INNER JOIN empresa e ON e.cliente_id = c.id AND e.ativo = true AND e.tenant_id = $1::uuid
			LEFT JOIN clientes_dados cd ON cd.cliente_id = c.id
			LEFT JOIN nfe_gestao n ON regexp_replace(COALESCE(n.cnpj_destinatario, ''), '[^0-9]', '', 'g') = regexp_replace(COALESCE(c.documento, ''), '[^0-9]', '', 'g')
				AND n.data_emissao >= $2::date
				AND n.data_emissao < $3::date
			WHERE c.ativo = true
			  AND NULLIF(BTRIM(COALESCE(cd.email_contato, '')), '') IS NOT NULL
			GROUP BY c.id, c.nome, cd.email_contato
			HAVING COALESCE(SUM(COALESCE(n.valor_total, 0)), 0) > 0
			ORDER BY c.nome ASC
		`, strings.TrimSpace(tenantID), inicio.Format("2006-01-02"), fim.Format("2006-01-02"))
		if err != nil {
			return fmt.Errorf("query resumo mensal por cliente: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var item ResumoMensalCliente
			if err := rows.Scan(&item.ClienteID, &item.ClienteNome, &item.EmailContato, &item.ValorTotal); err != nil {
				return fmt.Errorf("scan resumo mensal por cliente: %w", err)
			}
			out = append(out, item)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
