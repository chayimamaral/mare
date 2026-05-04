package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LancamentoFolhaRepository struct {
	pool *pgxpool.Pool
}

func NewLancamentoFolhaRepository(pool *pgxpool.Pool) *LancamentoFolhaRepository {
	return &LancamentoFolhaRepository{pool: pool}
}

// ListClientesWithLancamentos retorna a arvore de clientes com seus lancamentos.
// Nivel 1: clientes ativos do tenant.
// Nivel 2: lancamentos mensais ordenados pelo mais recente.
func (r *LancamentoFolhaRepository) ListClientesWithLancamentos(ctx context.Context, tenantID string) ([]domain.LancamentoFolhaTreeNode, error) {
	clientesQ := `
		SELECT
			e.id,
			c.nome,
			COALESCE(NULLIF(BTRIM(c.documento), ''), NULLIF(BTRIM(ed.cnpj::text), ''))
		FROM empresa e
		INNER JOIN cliente c ON c.id = e.cliente_id
		LEFT JOIN clientes_dados ed ON ed.cliente_id = c.id
		WHERE e.tenant_id = $1 AND e.ativo = true AND c.ativo = true
		ORDER BY c.nome ASC`

	cRows, err := dbQuery(ctx, r.pool, clientesQ, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list clientes lancamentos: %w", err)
	}
	defer cRows.Close()

	type clienteInfo struct {
		id        string
		nome      string
		documento string
	}

	var clientes []clienteInfo
	for cRows.Next() {
		var ci clienteInfo
		var doc sql.NullString
		if err := cRows.Scan(&ci.id, &ci.nome, &doc); err != nil {
			return nil, fmt.Errorf("scan cliente lancamento: %w", err)
		}
		ci.documento = doc.String
		clientes = append(clientes, ci)
	}

	if err := cRows.Err(); err != nil {
		return nil, err
	}

	lancamentosQ := `
		SELECT id, cliente_id, competencia, valor_folha, valor_faturamento, observacoes
		FROM lancamentos_folha
		WHERE tenant_id = $1
		ORDER BY competencia DESC`

	lRows, err := dbQuery(ctx, r.pool, lancamentosQ, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list lancamentos: %w", err)
	}
	defer lRows.Close()

	lancamentosByCliente := make(map[string][]domain.LancamentoFolhaTreeNode)
	for lRows.Next() {
		var id, clienteID, observacoes string
		var competencia time.Time
		var valorFolha, valorFaturamento float64

		if err := lRows.Scan(&id, &clienteID, &competencia, &valorFolha, &valorFaturamento, &observacoes); err != nil {
			return nil, fmt.Errorf("scan lancamento: %w", err)
		}

		compStr := competencia.Format("01/2006")
		childNode := domain.LancamentoFolhaTreeNode{
			Key:  "l-" + id,
			Leaf: true,
			Data: domain.LancamentoFolhaLancamentoNode{
				Tipo:             "lancamento",
				ID:               id,
				Competencia:      compStr,
				ValorFolha:       valorFolha,
				ValorFaturamento: valorFaturamento,
				Observacoes:      observacoes,
			},
		}
		lancamentosByCliente[clienteID] = append(lancamentosByCliente[clienteID], childNode)
	}

	if err := lRows.Err(); err != nil {
		return nil, err
	}

	tree := make([]domain.LancamentoFolhaTreeNode, 0, len(clientes))
	for _, cli := range clientes {
		lancs := lancamentosByCliente[cli.id]
		if lancs == nil {
			lancs = []domain.LancamentoFolhaTreeNode{}
		}

		var totalFolha, totalFaturamento float64
		for _, l := range lancs {
			data := l.Data.(domain.LancamentoFolhaLancamentoNode)
			totalFolha += data.ValorFolha
			totalFaturamento += data.ValorFaturamento
		}

		tree = append(tree, domain.LancamentoFolhaTreeNode{
			Key:  "c-" + cli.id,
			Leaf: len(lancs) == 0,
			Data: domain.LancamentoFolhaClienteNode{
				Tipo:            "cliente",
				ClienteID:       cli.id,
				Nome:            cli.nome,
				Documento:       cli.documento,
				TotalFolha:      totalFolha,
				TotalFaturamento: totalFaturamento,
			},
			Children: lancs,
		})
	}

	return tree, nil
}

// Create insere um lancamento com validacao de duplicidade por cliente+competencia.
func (r *LancamentoFolhaRepository) Create(ctx context.Context, tenantID, clienteID string, competencia time.Time, valorFolha, valorFaturamento float64, observacoes string) (domain.LancamentoFolha, error) {
	const insertQ = `
		INSERT INTO lancamentos_folha (tenant_id, cliente_id, competencia, valor_folha, valor_faturamento, observacoes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, tenant_id, cliente_id, competencia, valor_folha, valor_faturamento, observacoes, created_at, updated_at`

	var lf domain.LancamentoFolha
	var comp time.Time
	err := dbQueryRow(ctx, r.pool, insertQ, tenantID, clienteID, competencia, valorFolha, valorFaturamento, observacoes).Scan(
		&lf.ID, &lf.TenantID, &lf.ClienteID, &comp,
		&lf.ValorFolha, &lf.ValorFaturamento, &lf.Observacoes,
		&lf.CreatedAt, &lf.UpdatedAt,
	)
	if err != nil {
		return domain.LancamentoFolha{}, fmt.Errorf("create lancamento folha: %w", err)
	}

	lf.Competencia = comp.Format("01/2006")
	return lf, nil
}

// Update atualiza um lancamento existente.
func (r *LancamentoFolhaRepository) Update(ctx context.Context, id, tenantID string, competencia time.Time, valorFolha, valorFaturamento float64, observacoes string) (domain.LancamentoFolha, error) {
	const updateQ = `
		UPDATE lancamentos_folha
		SET competencia = $1, valor_folha = $2, valor_faturamento = $3, observacoes = $4, updated_at = now()
		WHERE id = $5 AND tenant_id = $6
		RETURNING id, tenant_id, cliente_id, competencia, valor_folha, valor_faturamento, observacoes, created_at, updated_at`

	var lf domain.LancamentoFolha
	var comp time.Time
	err := dbQueryRow(ctx, r.pool, updateQ, competencia, valorFolha, valorFaturamento, observacoes, id, tenantID).Scan(
		&lf.ID, &lf.TenantID, &lf.ClienteID, &comp,
		&lf.ValorFolha, &lf.ValorFaturamento, &lf.Observacoes,
		&lf.CreatedAt, &lf.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.LancamentoFolha{}, fmt.Errorf("lancamento nao encontrado")
		}
		return domain.LancamentoFolha{}, fmt.Errorf("update lancamento folha: %w", err)
	}

	lf.Competencia = comp.Format("01/2006")
	return lf, nil
}

// Delete remove um lancamento.
func (r *LancamentoFolhaRepository) Delete(ctx context.Context, id, tenantID string) error {
	const deleteQ = `DELETE FROM lancamentos_folha WHERE id = $1 AND tenant_id = $2`

	tag, err := dbExec(ctx, r.pool, deleteQ, id, tenantID)
	if err != nil {
		return fmt.Errorf("delete lancamento folha: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("lancamento nao encontrado")
	}
	return nil
}

// GetByID retorna um lancamento pelo ID.
func (r *LancamentoFolhaRepository) GetByID(ctx context.Context, id, tenantID string) (domain.LancamentoFolha, error) {
	const q = `
		SELECT id, tenant_id, cliente_id, competencia, valor_folha, valor_faturamento, observacoes, created_at, updated_at
		FROM lancamentos_folha
		WHERE id = $1 AND tenant_id = $2`

	var lf domain.LancamentoFolha
	var comp time.Time
	err := dbQueryRow(ctx, r.pool, q, id, tenantID).Scan(
		&lf.ID, &lf.TenantID, &lf.ClienteID, &comp,
		&lf.ValorFolha, &lf.ValorFaturamento, &lf.Observacoes,
		&lf.CreatedAt, &lf.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.LancamentoFolha{}, fmt.Errorf("lancamento nao encontrado")
		}
		return domain.LancamentoFolha{}, fmt.Errorf("get lancamento folha: %w", err)
	}

	lf.Competencia = comp.Format("01/2006")
	return lf, nil
}
