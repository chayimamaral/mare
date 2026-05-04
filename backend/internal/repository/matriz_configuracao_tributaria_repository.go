package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatrizConfiguracaoTributariaListParams struct {
	First     int
	Rows      int
	SortField string
	SortOrder int
	Nome      string
	Ativo     *bool
}

type MatrizConfiguracaoTributariaRepository struct {
	pool *pgxpool.Pool
}

func NewMatrizConfiguracaoTributariaRepository(pool *pgxpool.Pool) *MatrizConfiguracaoTributariaRepository {
	return &MatrizConfiguracaoTributariaRepository{pool: pool}
}

func (r *MatrizConfiguracaoTributariaRepository) List(ctx context.Context, params MatrizConfiguracaoTributariaListParams) ([]domain.MatrizConfiguracaoTributaria, int64, error) {
	whereParts := []string{"1 = 1"}
	args := []any{}
	argIndex := 1

	if strings.TrimSpace(params.Nome) != "" {
		whereParts = append(whereParts, fmt.Sprintf("m.nome ILIKE $%d", argIndex))
		args = append(args, "%"+strings.TrimSpace(params.Nome)+"%")
		argIndex++
	}

	if params.Ativo != nil {
		whereParts = append(whereParts, fmt.Sprintf("m.ativo = $%d", argIndex))
		args = append(args, *params.Ativo)
		argIndex++
	}

	orderBy := "m.nome ASC"
	switch params.SortField {
	case "natureza_juridica":
		if params.SortOrder == -1 {
			orderBy = "t.descricao DESC"
		} else {
			orderBy = "t.descricao ASC"
		}
	case "enquadramento_porte":
		if params.SortOrder == -1 {
			orderBy = "e.sigla DESC"
		} else {
			orderBy = "e.sigla ASC"
		}
	case "regime_tributario":
		if params.SortOrder == -1 {
			orderBy = "r.nome DESC"
		} else {
			orderBy = "r.nome ASC"
		}
	case "aliquota_base":
		if params.SortOrder == -1 {
			orderBy = "m.aliquota_base DESC"
		} else {
			orderBy = "m.aliquota_base ASC"
		}
	case "nome":
		if params.SortOrder == -1 {
			orderBy = "m.nome DESC"
		} else {
			orderBy = "m.nome ASC"
		}
	}

	listQuery := fmt.Sprintf(`
		SELECT
			m.id,
			m.nome,
			m.natureza_juridica_id,
			t.descricao,
			m.enquadramento_porte_id,
			COALESCE(e.sigla || ' - ' || e.descricao, e.sigla),
			m.regime_tributario_id,
			r.nome,
			m.aliquota_base,
			m.possui_fator_r,
			m.aliquota_fator_r,
			m.ativo
		FROM public.matriz_configuracao_tributaria m
		JOIN public.tipoempresa t ON t.id = m.natureza_juridica_id
		JOIN public.enquadramento_juridico_porte e ON e.id = m.enquadramento_porte_id
		JOIN public.regime_tributario r ON r.id = m.regime_tributario_id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		strings.Join(whereParts, " AND "), orderBy, argIndex, argIndex+1)
	args = append(args, params.Rows, params.First)

	rows, err := dbQuery(ctx, r.pool, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list matriz: %w", err)
	}
	defer rows.Close()

	items := make([]domain.MatrizConfiguracaoTributaria, 0)
	for rows.Next() {
		var m domain.MatrizConfiguracaoTributaria
		if err := rows.Scan(
			&m.ID, &m.Nome,
			&m.NaturezaJuridicaID, &m.NaturezaJuridica,
			&m.EnquadramentoPorteID, &m.EnquadramentoPorte,
			&m.RegimeTributarioID, &m.RegimeTributario,
			&m.AliquotaBase, &m.PossuiFatorR, &m.AliquotaFatorR,
			&m.Ativo,
		); err != nil {
			return nil, 0, fmt.Errorf("scan matriz: %w", err)
		}
		items = append(items, m)
	}

	countQuery := fmt.Sprintf(`
		SELECT count(*)
		FROM public.matriz_configuracao_tributaria m
		JOIN public.tipoempresa t ON t.id = m.natureza_juridica_id
		JOIN public.enquadramento_juridico_porte e ON e.id = m.enquadramento_porte_id
		JOIN public.regime_tributario r ON r.id = m.regime_tributario_id
		WHERE %s`, strings.Join(whereParts, " AND "))

	var total int64
	if err := dbQueryRow(ctx, r.pool, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count matriz: %w", err)
	}

	return items, total, nil
}

func (r *MatrizConfiguracaoTributariaRepository) Create(ctx context.Context, nome, naturezaJuridicaID, enquadramentoPorteID, regimeTributarioID string, aliquotaBase float64, possuiFatorR bool, aliquotaFatorR float64) (domain.MatrizConfiguracaoTributaria, error) {
	const query = `
		INSERT INTO public.matriz_configuracao_tributaria
			(nome, natureza_juridica_id, enquadramento_porte_id, regime_tributario_id, aliquota_base, possui_fator_r, aliquota_fator_r)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, nome, natureza_juridica_id, enquadramento_porte_id, regime_tributario_id, aliquota_base, possui_fator_r, aliquota_fator_r, ativo`

	var m domain.MatrizConfiguracaoTributaria
	err := dbQueryRow(ctx, r.pool, query, nome, naturezaJuridicaID, enquadramentoPorteID, regimeTributarioID, aliquotaBase, possuiFatorR, aliquotaFatorR).Scan(
		&m.ID, &m.Nome, &m.NaturezaJuridicaID, &m.EnquadramentoPorteID, &m.RegimeTributarioID,
		&m.AliquotaBase, &m.PossuiFatorR, &m.AliquotaFatorR, &m.Ativo,
	)
	if err != nil {
		return domain.MatrizConfiguracaoTributaria{}, fmt.Errorf("create matriz: %w", err)
	}
	return m, nil
}

func (r *MatrizConfiguracaoTributariaRepository) Update(ctx context.Context, id, nome, naturezaJuridicaID, enquadramentoPorteID, regimeTributarioID string, aliquotaBase float64, possuiFatorR bool, aliquotaFatorR float64, ativo bool) (domain.MatrizConfiguracaoTributaria, error) {
	const query = `
		UPDATE public.matriz_configuracao_tributaria
		SET nome = $1, natureza_juridica_id = $2, enquadramento_porte_id = $3,
		    regime_tributario_id = $4, aliquota_base = $5, possui_fator_r = $6,
		    aliquota_fator_r = $7, ativo = $8, updated_at = now()
		WHERE id = $9
		RETURNING id, nome, natureza_juridica_id, enquadramento_porte_id, regime_tributario_id, aliquota_base, possui_fator_r, aliquota_fator_r, ativo`

	var m domain.MatrizConfiguracaoTributaria
	err := dbQueryRow(ctx, r.pool, query, nome, naturezaJuridicaID, enquadramentoPorteID, regimeTributarioID, aliquotaBase, possuiFatorR, aliquotaFatorR, ativo, id).Scan(
		&m.ID, &m.Nome, &m.NaturezaJuridicaID, &m.EnquadramentoPorteID, &m.RegimeTributarioID,
		&m.AliquotaBase, &m.PossuiFatorR, &m.AliquotaFatorR, &m.Ativo,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.MatrizConfiguracaoTributaria{}, fmt.Errorf("configuracao nao encontrada")
		}
		return domain.MatrizConfiguracaoTributaria{}, fmt.Errorf("update matriz: %w", err)
	}
	return m, nil
}

func (r *MatrizConfiguracaoTributariaRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM public.matriz_configuracao_tributaria WHERE id = $1`
	tag, err := dbExec(ctx, r.pool, query, id)
	if err != nil {
		return fmt.Errorf("delete matriz: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("configuracao nao encontrada")
	}
	return nil
}

func (r *MatrizConfiguracaoTributariaRepository) GetByID(ctx context.Context, id string) (domain.MatrizConfiguracaoTributaria, error) {
	const query = `
		SELECT
			m.id, m.nome,
			m.natureza_juridica_id, t.descricao,
			m.enquadramento_porte_id, COALESCE(e.sigla || ' - ' || e.descricao, e.sigla),
			m.regime_tributario_id, r.nome,
			m.aliquota_base, m.possui_fator_r, m.aliquota_fator_r, m.ativo
		FROM public.matriz_configuracao_tributaria m
		JOIN public.tipoempresa t ON t.id = m.natureza_juridica_id
		JOIN public.enquadramento_juridico_porte e ON e.id = m.enquadramento_porte_id
		JOIN public.regime_tributario r ON r.id = m.regime_tributario_id
		WHERE m.id = $1`

	var item domain.MatrizConfiguracaoTributaria
	err := dbQueryRow(ctx, r.pool, query, id).Scan(
		&item.ID, &item.Nome,
		&item.NaturezaJuridicaID, &item.NaturezaJuridica,
		&item.EnquadramentoPorteID, &item.EnquadramentoPorte,
		&item.RegimeTributarioID, &item.RegimeTributario,
		&item.AliquotaBase, &item.PossuiFatorR, &item.AliquotaFatorR,
		&item.Ativo,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.MatrizConfiguracaoTributaria{}, fmt.Errorf("configuracao nao encontrada")
		}
		return domain.MatrizConfiguracaoTributaria{}, fmt.Errorf("get matriz: %w", err)
	}
	return item, nil
}
