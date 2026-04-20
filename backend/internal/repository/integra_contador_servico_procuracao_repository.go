package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IntegraContadorServicoProcuracaoRepository struct {
	pool *pgxpool.Pool
}

func NewIntegraContadorServicoProcuracaoRepository(pool *pgxpool.Pool) *IntegraContadorServicoProcuracaoRepository {
	return &IntegraContadorServicoProcuracaoRepository{pool: pool}
}

func (r *IntegraContadorServicoProcuracaoRepository) List(ctx context.Context, idSistema, idServico string) ([]domain.IntegraContadorServicoProcuracao, error) {
	args := []any{}
	conds := make([]string, 0, 2)
	if strings.TrimSpace(idSistema) != "" {
		args = append(args, strings.TrimSpace(idSistema))
		conds = append(conds, fmt.Sprintf("id_sistema = $%d", len(args)))
	}
	if strings.TrimSpace(idServico) != "" {
		args = append(args, strings.TrimSpace(idServico))
		conds = append(conds, fmt.Sprintf("id_servico = $%d", len(args)))
	}
	where := "true"
	if len(conds) > 0 {
		where = strings.Join(conds, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT COALESCE(id_sistema, ''), COALESCE(id_servico, ''), COALESCE(cod_procuracao, ''), COALESCE(nome_servico, '')
		  FROM public.integra_contador_servicos
		 WHERE %s
		 ORDER BY id_sistema, id_servico`, where)

	rows, err := dbQuery(ctx, r.pool, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list integra_contador_servicos: %w", err)
	}
	defer rows.Close()

	out := make([]domain.IntegraContadorServicoProcuracao, 0)
	for rows.Next() {
		var item domain.IntegraContadorServicoProcuracao
		if err := rows.Scan(&item.IDSistema, &item.IDServico, &item.CodProcuracao, &item.NomeServico); err != nil {
			return nil, fmt.Errorf("scan integra_contador_servicos: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows integra_contador_servicos: %w", err)
	}
	return out, nil
}
