package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClienteRepository acesso unificado a clientes PF/PJ: cadastro em public.cliente,
// operação em public.empresa; o ID exposto na API permanece empresa.id (agenda/compromissos).
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
			c.rotina_id::text,
			c.rotina_pf_id::text,
			COALESCE(rpf.nome, ''),
			COALESCE(rpf.categoria, ''),
			c.cnaes,
			COALESCE(c.bairro, ''),
			e.iniciado,
			e.ativo
		FROM public.empresa e
		INNER JOIN public.cliente c ON c.id = e.cliente_id
		LEFT JOIN public.clientes_dados ed ON ed.cliente_id = c.id
		LEFT JOIN public.rotina_pf rpf ON rpf.id = c.rotina_pf_id
		WHERE e.id = $1 AND e.tenant_id = $2 AND e.ativo = true AND c.ativo = true`

	var c domain.Cliente
	var doc, munID, rotID, rpfID sql.NullString
	var rpfNome, rpfCat string
	if err := r.pool.QueryRow(ctx, q, id, tenantID).Scan(
		&c.ID,
		&c.TenantID,
		&c.TipoPessoa,
		&c.Nome,
		&doc,
		&munID,
		&rotID,
		&rpfID,
		&rpfNome,
		&rpfCat,
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
	if rotID.Valid {
		s := rotID.String
		c.RotinaID = &s
	}
	if rpfID.Valid && strings.TrimSpace(rpfID.String) != "" {
		s := rpfID.String
		c.RotinaPFID = &s
	}
	c.RotinaPFNome = rpfNome
	c.CategoriaPF = rpfCat
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
			c.rotina_id::text,
			c.rotina_pf_id::text,
			COALESCE(rpf.nome, ''),
			COALESCE(rpf.categoria, ''),
			c.cnaes,
			COALESCE(c.bairro, ''),
			e.iniciado,
			e.ativo
		FROM public.empresa e
		INNER JOIN public.cliente c ON c.id = e.cliente_id
		LEFT JOIN public.clientes_dados ed ON ed.cliente_id = c.id
		LEFT JOIN public.rotina_pf rpf ON rpf.id = c.rotina_pf_id
		WHERE e.tenant_id = $1 AND e.ativo = true AND c.ativo = true
		ORDER BY c.nome ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, q, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list clientes: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Cliente, 0)
	for rows.Next() {
		var c domain.Cliente
		var doc, munID, rotID, rpfID sql.NullString
		var rpfNome, rpfCat string
		if err := rows.Scan(
			&c.ID,
			&c.TenantID,
			&c.TipoPessoa,
			&c.Nome,
			&doc,
			&munID,
			&rotID,
			&rpfID,
			&rpfNome,
			&rpfCat,
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
		if rotID.Valid {
			s := rotID.String
			c.RotinaID = &s
		}
		if rpfID.Valid && strings.TrimSpace(rpfID.String) != "" {
			s := rpfID.String
			c.RotinaPFID = &s
		}
		c.RotinaPFNome = rpfNome
		c.CategoriaPF = rpfCat
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
