package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CertificadoClienteRepository struct {
	pool *pgxpool.Pool
}

func NewCertificadoClienteRepository(pool *pgxpool.Pool) *CertificadoClienteRepository {
	return &CertificadoClienteRepository{pool: pool}
}

// ClienteIDEmpresaTenant resolve o cliente vinculado à empresa no tenant (empresa.id exposto na API).
func (r *CertificadoClienteRepository) ClienteIDEmpresaTenant(ctx context.Context, empresaID, tenantID string) (string, error) {
	eid := strings.TrimSpace(empresaID)
	tid := strings.TrimSpace(tenantID)
	if eid == "" || tid == "" {
		return "", fmt.Errorf("empresa e tenant obrigatorios")
	}
	const q = `
		SELECT c.id::text
		FROM public.empresa e
		INNER JOIN public.cliente c ON c.id = e.cliente_id AND c.ativo = true
		WHERE e.id = $1 AND e.tenant_id = $2 AND e.ativo = true`
	var cid string
	if err := r.pool.QueryRow(ctx, q, eid, tid).Scan(&cid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("empresa nao encontrada neste tenant")
		}
		return "", fmt.Errorf("resolver cliente: %w", err)
	}
	return cid, nil
}

type CertificadoClienteResumoRow struct {
	CNPJ        string
	TitularNome string
	EmitidoPor  string
	ValidadeDe  time.Time
	ValidadeAte time.Time
}

// UpsertAtivo grava ou substitui o certificado ativo do cliente (PK = cliente_id).
func (r *CertificadoClienteRepository) UpsertAtivo(
	ctx context.Context,
	clienteID string,
	pfxCifrado, senhaCifrada []byte,
	cnpj, titularNome, emitidoPor string,
	validadeDe, validadeAte time.Time,
) error {
	cid := strings.TrimSpace(clienteID)
	if cid == "" {
		return fmt.Errorf("cliente_id obrigatorio")
	}
	if len(pfxCifrado) == 0 || len(senhaCifrada) == 0 {
		return fmt.Errorf("blobs cifrados obrigatorios")
	}
	const q = `
		INSERT INTO public.certificado_cliente (
			cliente_id, pfx_cifrado, senha_cifrada, cnpj, titular_nome, emitido_por,
			validade_de, validade_ate, ativo, atualizado_em
		)
		VALUES ($1, $2, $3, NULLIF(TRIM($4), ''), NULLIF(TRIM($5), ''), NULLIF(TRIM($6), ''), $7, $8, true, NOW())
		ON CONFLICT (cliente_id) DO UPDATE SET
			pfx_cifrado = EXCLUDED.pfx_cifrado,
			senha_cifrada = EXCLUDED.senha_cifrada,
			cnpj = EXCLUDED.cnpj,
			titular_nome = EXCLUDED.titular_nome,
			emitido_por = EXCLUDED.emitido_por,
			validade_de = EXCLUDED.validade_de,
			validade_ate = EXCLUDED.validade_ate,
			ativo = true,
			atualizado_em = NOW()`
	_, err := r.pool.Exec(ctx, q, cid, pfxCifrado, senhaCifrada, cnpj, titularNome, emitidoPor, validadeDe, validadeAte)
	if err != nil {
		return fmt.Errorf("gravar certificado_cliente: %w", err)
	}
	return nil
}

// GetResumoAtivo retorna metadados sem blobs cifrados.
func (r *CertificadoClienteRepository) GetResumoAtivo(ctx context.Context, clienteID string) (*CertificadoClienteResumoRow, error) {
	cid := strings.TrimSpace(clienteID)
	if cid == "" {
		return nil, fmt.Errorf("cliente_id obrigatorio")
	}
	const q = `
		SELECT COALESCE(cnpj, ''), COALESCE(titular_nome, ''), COALESCE(emitido_por, ''),
			validade_de, validade_ate
		FROM public.certificado_cliente
		WHERE cliente_id = $1 AND ativo = true`
	row := r.pool.QueryRow(ctx, q, cid)
	var out CertificadoClienteResumoRow
	if err := row.Scan(&out.CNPJ, &out.TitularNome, &out.EmitidoPor, &out.ValidadeDe, &out.ValidadeAte); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("ler certificado_cliente: %w", err)
	}
	return &out, nil
}
