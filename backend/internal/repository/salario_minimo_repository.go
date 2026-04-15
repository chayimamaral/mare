package repository

import (
	"context"
	"fmt"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SalarioMinimoRepository struct {
	pool *pgxpool.Pool
}

func NewSalarioMinimoRepository(pool *pgxpool.Pool) *SalarioMinimoRepository {
	return &SalarioMinimoRepository{pool: pool}
}

func (r *SalarioMinimoRepository) List(ctx context.Context) ([]domain.SalarioMinimoNacional, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, ano, valor::float8
		FROM public.salario_minimo_nacional
		ORDER BY ano DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list salario minimo nacional: %w", err)
	}
	defer rows.Close()

	out := make([]domain.SalarioMinimoNacional, 0)
	for rows.Next() {
		var rec domain.SalarioMinimoNacional
		if err := rows.Scan(&rec.ID, &rec.Ano, &rec.Valor); err != nil {
			return nil, fmt.Errorf("scan salario minimo nacional: %w", err)
		}
		out = append(out, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows salario minimo nacional: %w", err)
	}
	return out, nil
}

func (r *SalarioMinimoRepository) Create(ctx context.Context, ano int, valor float64) (domain.SalarioMinimoNacional, error) {
	var rec domain.SalarioMinimoNacional
	if err := r.pool.QueryRow(ctx, `
		INSERT INTO public.salario_minimo_nacional (ano, valor)
		VALUES ($1, $2)
		RETURNING id::text, ano, valor::float8
	`, ano, valor).Scan(&rec.ID, &rec.Ano, &rec.Valor); err != nil {
		return domain.SalarioMinimoNacional{}, fmt.Errorf("create salario minimo nacional: %w", err)
	}
	return rec, nil
}

func (r *SalarioMinimoRepository) Update(ctx context.Context, id string, ano int, valor float64) (domain.SalarioMinimoNacional, error) {
	var rec domain.SalarioMinimoNacional
	if err := r.pool.QueryRow(ctx, `
		UPDATE public.salario_minimo_nacional
		SET ano = $1, valor = $2, atualizado_em = now()
		WHERE id = $3::uuid
		RETURNING id::text, ano, valor::float8
	`, ano, valor, id).Scan(&rec.ID, &rec.Ano, &rec.Valor); err != nil {
		return domain.SalarioMinimoNacional{}, fmt.Errorf("update salario minimo nacional: %w", err)
	}
	return rec, nil
}

func (r *SalarioMinimoRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM public.salario_minimo_nacional WHERE id = $1::uuid`, id)
	if err != nil {
		return fmt.Errorf("delete salario minimo nacional: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("salario minimo nao encontrado")
	}
	return nil
}
