package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CnaeListParams struct {
	First       int
	Rows        int
	SortField   string
	SortOrder   int
	Denominacao string
	Subclasse   string
}

type CnaeRecord struct {
	ID          string `json:"id"`
	Denominacao string `json:"denominacao"`
	Subclasse   string `json:"subclasse"`
	Ativo       bool   `json:"ativo"`
}

type CnaeLiteItem struct {
	ID          string `json:"id"`
	Denominacao string `json:"denominacao"`
}

type CnaeValidateItem struct {
	ID string `json:"id"`
}

type CnaeRepository struct {
	pool *pgxpool.Pool
}

func NewCnaeRepository(pool *pgxpool.Pool) *CnaeRepository {
	return &CnaeRepository{pool: pool}
}

func (r *CnaeRepository) List(ctx context.Context, params CnaeListParams) ([]CnaeRecord, int64, error) {
	whereParts := []string{"ativo = true"}
	args := []any{}
	argIndex := 1

	if strings.TrimSpace(params.Denominacao) != "" {
		whereParts = append(whereParts, fmt.Sprintf("denominacao ILIKE $%d", argIndex))
		args = append(args, "%"+strings.TrimSpace(params.Denominacao)+"%")
		argIndex++
	} else if strings.TrimSpace(params.Subclasse) != "" {
		whereParts = append(whereParts, fmt.Sprintf("subclasse ILIKE $%d", argIndex))
		args = append(args, "%"+strings.TrimSpace(params.Subclasse)+"%")
		argIndex++
	}

	orderBy := "subclasse ASC"
	switch params.SortField {
	case "denominacao":
		if params.SortOrder == -1 {
			orderBy = "denominacao ASC"
		} else {
			orderBy = "denominacao DESC"
		}
	case "subclasse":
		if params.SortOrder == -1 {
			orderBy = "subclasse ASC"
		} else {
			orderBy = "subclasse DESC"
		}
	}

	query := fmt.Sprintf(
		"SELECT id, denominacao, subclasse, ativo FROM public.cnae WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d",
		strings.Join(whereParts, " AND "),
		orderBy,
		argIndex,
		argIndex+1,
	)
	args = append(args, params.Rows, params.First)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list cnae: %w", err)
	}
	defer rows.Close()

	cnaes := make([]CnaeRecord, 0)
	for rows.Next() {
		var id, denominacao, subclasse string
		var ativo bool
		if err := rows.Scan(&id, &denominacao, &subclasse, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan cnae: %w", err)
		}

		cnaes = append(cnaes, CnaeRecord{ID: id, Denominacao: denominacao, Subclasse: subclasse, Ativo: ativo})
	}

	countQuery := fmt.Sprintf("SELECT count(*) FROM public.cnae WHERE %s", strings.Join(whereParts, " AND "))
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count cnae: %w", err)
	}

	return cnaes, total, nil
}

func (r *CnaeRepository) Create(ctx context.Context, denominacao, subclasse string) ([]CnaeRecord, int64, error) {
	const query = `
		INSERT INTO public.cnae (denominacao, subclasse)
		VALUES ($1, $2)
		RETURNING id, denominacao, subclasse, ativo`

	rows, err := r.pool.Query(ctx, query, denominacao, subclasse)
	if err != nil {
		return nil, 0, fmt.Errorf("create cnae: %w", err)
	}
	defer rows.Close()

	cnaes := make([]CnaeRecord, 0)
	for rows.Next() {
		var id, d, s string
		var ativo bool
		if err := rows.Scan(&id, &d, &s, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan created cnae: %w", err)
		}
		cnaes = append(cnaes, CnaeRecord{ID: id, Denominacao: d, Subclasse: s, Ativo: ativo})
	}

	var total int64
	if err := r.pool.QueryRow(ctx, `SELECT count(*) FROM public.cnae WHERE ativo = true`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count cnae: %w", err)
	}

	return cnaes, total, nil
}

func (r *CnaeRepository) Update(ctx context.Context, id, denominacao, subclasse string) ([]CnaeRecord, int64, error) {
	const query = `
		UPDATE public.cnae
		SET denominacao = $1, subclasse = $2
		WHERE id = $3
		RETURNING id, denominacao, subclasse, ativo`

	rows, err := r.pool.Query(ctx, query, denominacao, subclasse, id)
	if err != nil {
		return nil, 0, fmt.Errorf("update cnae: %w", err)
	}
	defer rows.Close()

	cnaes := make([]CnaeRecord, 0)
	for rows.Next() {
		var cid, d, s string
		var ativo bool
		if err := rows.Scan(&cid, &d, &s, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan updated cnae: %w", err)
		}
		cnaes = append(cnaes, CnaeRecord{ID: cid, Denominacao: d, Subclasse: s, Ativo: ativo})
	}

	var total int64
	if err := r.pool.QueryRow(ctx, `SELECT count(*) FROM public.cnae WHERE ativo = true`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count cnae: %w", err)
	}

	return cnaes, total, nil
}

func (r *CnaeRepository) Delete(ctx context.Context, id string) ([]CnaeRecord, int64, error) {
	const query = `
		UPDATE public.cnae
		SET ativo = false
		WHERE id = $1
		RETURNING id, denominacao, subclasse, ativo`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, 0, fmt.Errorf("delete cnae: %w", err)
	}
	defer rows.Close()

	cnaes := make([]CnaeRecord, 0)
	for rows.Next() {
		var cid, d, s string
		var ativo bool
		if err := rows.Scan(&cid, &d, &s, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan deleted cnae: %w", err)
		}
		cnaes = append(cnaes, CnaeRecord{ID: cid, Denominacao: d, Subclasse: s, Ativo: ativo})
	}

	var total int64
	if err := r.pool.QueryRow(ctx, `SELECT count(*) FROM public.cnae WHERE ativo = true`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count cnae: %w", err)
	}

	return cnaes, total, nil
}

func (r *CnaeRepository) Lite(ctx context.Context) ([]CnaeLiteItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, denominacao FROM public.cnae WHERE ativo = true ORDER BY denominacao ASC`)
	if err != nil {
		return nil, fmt.Errorf("lite cnae: %w", err)
	}
	defer rows.Close()

	cnaes := make([]CnaeLiteItem, 0)
	for rows.Next() {
		var id, d string
		if err := rows.Scan(&id, &d); err != nil {
			return nil, fmt.Errorf("scan lite cnae: %w", err)
		}
		cnaes = append(cnaes, CnaeLiteItem{ID: id, Denominacao: d})
	}

	return cnaes, nil
}

func (r *CnaeRepository) Validate(ctx context.Context, cnae string) ([]CnaeValidateItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT id FROM public.cnae WHERE ativo = true AND subclasse = $1`, cnae)
	if err != nil {
		return nil, fmt.Errorf("validate cnae: %w", err)
	}
	defer rows.Close()

	result := make([]CnaeValidateItem, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan validate cnae: %w", err)
		}
		result = append(result, CnaeValidateItem{ID: id})
	}

	return result, nil
}
