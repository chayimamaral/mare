package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnquadramentoJuridicoPorteRepository struct {
	pool *pgxpool.Pool
}

func NewEnquadramentoJuridicoPorteRepository(pool *pgxpool.Pool) *EnquadramentoJuridicoPorteRepository {
	return &EnquadramentoJuridicoPorteRepository{pool: pool}
}

func (r *EnquadramentoJuridicoPorteRepository) List(ctx context.Context, anoVigencia *int) ([]domain.EnquadramentoJuridicoPorte, error) {
	where := []string{"ativo = true"}
	args := []any{}
	argN := 1
	if anoVigencia != nil {
		where = append(where, fmt.Sprintf("ano_vigencia = $%d", argN))
		args = append(args, *anoVigencia)
		argN++
	}
	q := fmt.Sprintf(`
		SELECT id::text, sigla, descricao, limite_inicial::float8,
		       limite_final::float8, ano_vigencia, ativo
		FROM public.enquadramento_juridico_porte
		WHERE %s
		ORDER BY ano_vigencia DESC, limite_inicial ASC`, strings.Join(where, " AND "))

	rows, err := dbQuery(ctx, r.pool, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list enquadramento_juridico_porte: %w", err)
	}
	defer rows.Close()

	out := make([]domain.EnquadramentoJuridicoPorte, 0)
	for rows.Next() {
		var rec domain.EnquadramentoJuridicoPorte
		var limF *float64
		if err := rows.Scan(&rec.ID, &rec.Sigla, &rec.Descricao, &rec.LimiteInicial, &limF, &rec.AnoVigencia, &rec.Ativo); err != nil {
			return nil, fmt.Errorf("scan enquadramento_juridico_porte: %w", err)
		}
		rec.LimiteFinal = limF
		out = append(out, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows enquadramento_juridico_porte: %w", err)
	}
	return out, nil
}

func (r *EnquadramentoJuridicoPorteRepository) Create(ctx context.Context, sigla, descricao string, limiteInicial float64, limiteFinal *float64, anoVigencia int) (domain.EnquadramentoJuridicoPorte, error) {
	const q = `
		INSERT INTO public.enquadramento_juridico_porte (sigla, descricao, limite_inicial, limite_final, ano_vigencia)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text, sigla, descricao, limite_inicial::float8, limite_final::float8, ano_vigencia, ativo`
	var rec domain.EnquadramentoJuridicoPorte
	var limF *float64
	err := dbQueryRow(ctx, r.pool, q, sigla, descricao, limiteInicial, limiteFinal, anoVigencia).Scan(
		&rec.ID, &rec.Sigla, &rec.Descricao, &rec.LimiteInicial, &limF, &rec.AnoVigencia, &rec.Ativo,
	)
	if err != nil {
		return domain.EnquadramentoJuridicoPorte{}, fmt.Errorf("create enquadramento_juridico_porte: %w", err)
	}
	rec.LimiteFinal = limF
	return rec, nil
}

func (r *EnquadramentoJuridicoPorteRepository) Update(ctx context.Context, id, sigla, descricao string, limiteInicial float64, limiteFinal *float64, anoVigencia int, ativo bool) (domain.EnquadramentoJuridicoPorte, error) {
	const q = `
		UPDATE public.enquadramento_juridico_porte
		SET sigla = $1, descricao = $2, limite_inicial = $3, limite_final = $4, ano_vigencia = $5, ativo = $6, atualizado_em = now()
		WHERE id = $7::uuid
		RETURNING id::text, sigla, descricao, limite_inicial::float8, limite_final::float8, ano_vigencia, ativo`
	var rec domain.EnquadramentoJuridicoPorte
	var limF *float64
	err := dbQueryRow(ctx, r.pool, q, sigla, descricao, limiteInicial, limiteFinal, anoVigencia, ativo, id).Scan(
		&rec.ID, &rec.Sigla, &rec.Descricao, &rec.LimiteInicial, &limF, &rec.AnoVigencia, &rec.Ativo,
	)
	if err != nil {
		return domain.EnquadramentoJuridicoPorte{}, fmt.Errorf("update enquadramento_juridico_porte: %w", err)
	}
	rec.LimiteFinal = limF
	return rec, nil
}

func (r *EnquadramentoJuridicoPorteRepository) Delete(ctx context.Context, id string) error {
	tag, err := dbExec(ctx, r.pool, `DELETE FROM public.enquadramento_juridico_porte WHERE id = $1::uuid`, id)
	if err != nil {
		return fmt.Errorf("delete enquadramento_juridico_porte: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("registro nao encontrado")
	}
	return nil
}
