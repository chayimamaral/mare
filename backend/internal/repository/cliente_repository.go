package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClienteRepository acesso unificado a clientes PF/PJ: tabelas cliente/empresa no schema do tenant (search_path).
// O ID exposto na API permanece empresa.id (agenda/compromissos).
// Validação PJ (rotina/CNAEs) e PF (documento) fica na camada de serviço.
type ClienteRepository struct {
	pool *pgxpool.Pool
}

func NewClienteRepository(pool *pgxpool.Pool) *ClienteRepository {
	return &ClienteRepository{pool: pool}
}

// GetByID carrega um cliente ativo do tenant (id = empresa.id). Documento: c.documento ou fallback em clientes_dados.cnpj.
func (r *ClienteRepository) GetByID(ctx context.Context, tenantID, id string) (domain.Cliente, error) {
	const q = `
		SELECT
			e.id,
			c.tenant_id,
			COALESCE(NULLIF(BTRIM(c.tipo_pessoa::text), ''), 'PJ')::text,
			c.nome,
			COALESCE(NULLIF(BTRIM(c.documento), ''), NULLIF(BTRIM(ed.cnpj::text), '')),
			COALESCE(c.municipio_id, ed.municipio_id)::text,
			c.cnaes,
			COALESCE(c.bairro, ''),
			e.iniciado,
			e.ativo
		FROM empresa e
		INNER JOIN cliente c ON c.id = e.cliente_id
		LEFT JOIN clientes_dados ed ON ed.cliente_id = c.id
		WHERE e.id = $1 AND e.tenant_id = $2 AND e.ativo = true AND c.ativo = true`

	var c domain.Cliente
	var doc, munID sql.NullString
	if err := dbQueryRow(ctx, r.pool, q, id, tenantID).Scan(
		&c.ID,
		&c.TenantID,
		&c.TipoPessoa,
		&c.Nome,
		&doc,
		&munID,
		&c.Cnaes,
		&c.Bairro,
		&c.Iniciado,
		&c.Ativo,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.Cliente{}, fmt.Errorf("cliente nao encontrado")
		}
		return domain.Cliente{}, fmt.Errorf("get cliente: %w", err)
	}
	if doc.Valid {
		c.Documento = doc.String
	}
	if munID.Valid {
		s := munID.String
		c.MunicipioID = &s
	}
	return c, nil
}

// ListByTenant lista clientes ativos do escritório (paginação pode ser acrescentada depois).
func (r *ClienteRepository) ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]domain.Cliente, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	const q = `
		SELECT
			e.id,
			c.tenant_id,
			COALESCE(NULLIF(BTRIM(c.tipo_pessoa::text), ''), 'PJ')::text,
			c.nome,
			COALESCE(NULLIF(BTRIM(c.documento), ''), NULLIF(BTRIM(ed.cnpj::text), '')),
			COALESCE(c.municipio_id, ed.municipio_id)::text,
			c.cnaes,
			COALESCE(c.bairro, ''),
			e.iniciado,
			e.ativo
		FROM empresa e
		INNER JOIN cliente c ON c.id = e.cliente_id
		LEFT JOIN clientes_dados ed ON ed.cliente_id = c.id
		WHERE e.tenant_id = $1 AND e.ativo = true AND c.ativo = true
		ORDER BY c.nome ASC
		LIMIT $2 OFFSET $3`

	rows, err := dbQuery(ctx, r.pool, q, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list clientes: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Cliente, 0)
	for rows.Next() {
		var c domain.Cliente
		var doc, munID sql.NullString
		if err := rows.Scan(
			&c.ID,
			&c.TenantID,
			&c.TipoPessoa,
			&c.Nome,
			&doc,
			&munID,
			&c.Cnaes,
			&c.Bairro,
			&c.Iniciado,
			&c.Ativo,
		); err != nil {
			return nil, fmt.Errorf("scan cliente: %w", err)
		}
		if doc.Valid {
			c.Documento = doc.String
		}
		if munID.Valid {
			s := munID.String
			c.MunicipioID = &s
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
