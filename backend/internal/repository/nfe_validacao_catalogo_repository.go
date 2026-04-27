package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NFEValidacaoCatalogoRepository struct {
	pool *pgxpool.Pool
}

func NewNFEValidacaoCatalogoRepository(pool *pgxpool.Pool) *NFEValidacaoCatalogoRepository {
	return &NFEValidacaoCatalogoRepository{pool: pool}
}

func (r *NFEValidacaoCatalogoRepository) ListRegrasAtivasPorEtapa(ctx context.Context, etapa string) ([]domain.NFEValidacaoRegra, error) {
	const q = `
		SELECT id::text, etapa, codigo_regra, titulo, descricao, severidade, ordem, ativo
		FROM public.nfe_validacao_regra
		WHERE ativo = true
		  AND ($1 = '' OR etapa = $1)
		ORDER BY etapa, ordem, codigo_regra`
	rows, err := dbQuery(ctx, r.pool, q, strings.TrimSpace(etapa))
	if err != nil {
		return nil, fmt.Errorf("list nfe_validacao_regra: %w", err)
	}
	defer rows.Close()

	out := make([]domain.NFEValidacaoRegra, 0, 16)
	for rows.Next() {
		var it domain.NFEValidacaoRegra
		if err := rows.Scan(&it.ID, &it.Etapa, &it.CodigoRegra, &it.Titulo, &it.Descricao, &it.Severidade, &it.Ordem, &it.Ativo); err != nil {
			return nil, fmt.Errorf("scan nfe_validacao_regra: %w", err)
		}
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iter nfe_validacao_regra: %w", err)
	}
	return out, nil
}

func (r *NFEValidacaoCatalogoRepository) GetCodigoErro(ctx context.Context, origem, codigo string) (*domain.NFECodigoErro, error) {
	const q = `
		SELECT id::text, origem, etapa_validacao, codigo, mensagem, descricao_tecnica, acao_sugerida, COALESCE(http_status, 0), ativo
		FROM public.nfe_codigo_erro
		WHERE ativo = true
		  AND origem = $1
		  AND codigo = $2
		LIMIT 1`
	row := dbQueryRow(ctx, r.pool, q, strings.TrimSpace(origem), strings.TrimSpace(codigo))
	var out domain.NFECodigoErro
	if err := row.Scan(
		&out.ID,
		&out.Origem,
		&out.EtapaValidacao,
		&out.Codigo,
		&out.Mensagem,
		&out.DescricaoTec,
		&out.AcaoSugerida,
		&out.HTTPStatus,
		&out.Ativo,
	); err != nil {
		return nil, err
	}
	return &out, nil
}
