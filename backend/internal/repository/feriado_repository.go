package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeriadoListParams struct {
	First       int
	Rows        int
	Descricao   string
	HolidayCode string
}

type FeriadoUpsertInput struct {
	ID          string
	Descricao   string
	Data        string
	HolidayCode string
	MunicipioID string
	EstadoID    string
}

type FeriadoRepository struct {
	pool *pgxpool.Pool
}

func NewFeriadoRepository(pool *pgxpool.Pool) *FeriadoRepository {
	return &FeriadoRepository{pool: pool}
}

func (r *FeriadoRepository) List(ctx context.Context, params FeriadoListParams) ([]domain.FeriadoListItem, int64, error) {
	holiday := strings.TrimSpace(params.HolidayCode)
	if holiday == "" {
		holiday = "VARIAVEL"
	}

	whereParts := []string{"f.ativo = true", "f.feriado = $1"}
	args := []any{holiday}
	argIndex := 2

	if strings.TrimSpace(params.Descricao) != "" {
		whereParts = append(whereParts, fmt.Sprintf("f.descricao ILIKE $%d", argIndex))
		args = append(args, "%"+strings.TrimSpace(params.Descricao)+"%")
		argIndex++
	}

	join := ""
	selectExtra := ""
	if holiday == "MUNICIPAL" {
		join = `
			LEFT JOIN public.feriado_municipal fm ON fm.feriado_id = f.id
			LEFT JOIN public.municipio m ON m.id = fm.municipio_id`
		selectExtra = ", m.id, m.nome"
	}
	if holiday == "ESTADUAL" {
		join = `
			LEFT JOIN public.feriado_estadual fe ON fe.feriado_id = f.id
			LEFT JOIN public.estado e ON e.id = fe.uf_id`
		selectExtra = ", e.id, e.nome"
	}

	query := fmt.Sprintf(`
		SELECT f.id, f.descricao, f.data, f.feriado%s
		FROM public.feriados f
		%s
		WHERE %s
		ORDER BY TO_DATE(f.data, 'DD/MM') ASC
		LIMIT $%d OFFSET $%d`, selectExtra, join, strings.Join(whereParts, " AND "), argIndex, argIndex+1)
	args = append(args, params.Rows, params.First)

	rows, err := dbQuery(ctx, r.pool, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list feriados: %w", err)
	}
	defer rows.Close()

	feriados := make([]domain.FeriadoListItem, 0)
	for rows.Next() {
		if holiday == "MUNICIPAL" {
			var id, descricao, data, tipo, mid, mnome string
			if err := rows.Scan(&id, &descricao, &data, &tipo, &mid, &mnome); err != nil {
				return nil, 0, fmt.Errorf("scan feriado municipal: %w", err)
			}
			feriados = append(feriados, domain.FeriadoListItem{
				ID:        id,
				Descricao: descricao,
				Data:      data,
				Feriado:   tipo,
				Municipio: &domain.FeriadoRef{ID: mid, Nome: mnome},
			})
			continue
		}

		if holiday == "ESTADUAL" {
			var id, descricao, data, tipo, eid, enome string
			if err := rows.Scan(&id, &descricao, &data, &tipo, &eid, &enome); err != nil {
				return nil, 0, fmt.Errorf("scan feriado estadual: %w", err)
			}
			feriados = append(feriados, domain.FeriadoListItem{
				ID:        id,
				Descricao: descricao,
				Data:      data,
				Feriado:   tipo,
				Estado:    &domain.FeriadoRef{ID: eid, Nome: enome},
			})
			continue
		}

		var id, descricao, data, tipo string
		if err := rows.Scan(&id, &descricao, &data, &tipo); err != nil {
			return nil, 0, fmt.Errorf("scan feriado: %w", err)
		}
		feriados = append(feriados, domain.FeriadoListItem{ID: id, Descricao: descricao, Data: data, Feriado: tipo})
	}

	countQuery := fmt.Sprintf("SELECT count(*) FROM public.feriados f WHERE %s", strings.Join(whereParts, " AND "))
	var total int64
	if err := dbQueryRow(ctx, r.pool, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count feriados: %w", err)
	}

	return feriados, total, nil
}

func (r *FeriadoRepository) Create(ctx context.Context, input FeriadoUpsertInput) ([]domain.FeriadoMutationItem, int64, error) {
	const existsQuery = `SELECT count(*) FROM public.feriados WHERE descricao = $1`
	var count int64
	if err := dbQueryRow(ctx, r.pool, existsQuery, input.Descricao).Scan(&count); err != nil {
		return nil, 0, fmt.Errorf("check feriado exists: %w", err)
	}
	if count > 0 {
		return nil, 0, fmt.Errorf("Feriado ja cadastrado")
	}

	const query = `
		INSERT INTO public.feriados (descricao, data, feriado)
		VALUES ($1, $2, $3)
		RETURNING id, descricao, data, feriado, ativo`

	rows, err := dbQuery(ctx, r.pool, query, input.Descricao, input.Data, input.HolidayCode)
	if err != nil {
		return nil, 0, fmt.Errorf("create feriado: %w", err)
	}
	defer rows.Close()

	feriados := make([]domain.FeriadoMutationItem, 0)
	var createdID string
	for rows.Next() {
		var id, descricao, data, tipo string
		var ativo bool
		if err := rows.Scan(&id, &descricao, &data, &tipo, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan created feriado: %w", err)
		}
		createdID = id
		feriados = append(feriados, domain.FeriadoMutationItem{ID: id, Descricao: descricao, Data: data, Feriado: tipo, Ativo: ativo})
	}

	if input.HolidayCode == "MUNICIPAL" && strings.TrimSpace(input.MunicipioID) != "" {
		_, _ = dbExec(ctx, r.pool, `INSERT INTO public.feriado_municipal (feriado_id, municipio_id) VALUES ($1, $2)`, createdID, input.MunicipioID)
	}
	if input.HolidayCode == "ESTADUAL" && strings.TrimSpace(input.EstadoID) != "" {
		_, _ = dbExec(ctx, r.pool, `INSERT INTO public.feriado_estadual (feriado_id, uf_id) VALUES ($1, $2)`, createdID, input.EstadoID)
	}

	return feriados, int64(len(feriados)), nil
}

func (r *FeriadoRepository) Update(ctx context.Context, input FeriadoUpsertInput) ([]domain.FeriadoMutationItem, int64, error) {
	const query = `
		UPDATE public.feriados
		SET descricao = $1, data = $2
		WHERE id = $3
		RETURNING id, descricao, data, feriado, ativo`

	rows, err := dbQuery(ctx, r.pool, query, input.Descricao, input.Data, input.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("update feriado: %w", err)
	}
	defer rows.Close()

	feriados := make([]domain.FeriadoMutationItem, 0)
	for rows.Next() {
		var id, descricao, data, tipo string
		var ativo bool
		if err := rows.Scan(&id, &descricao, &data, &tipo, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan updated feriado: %w", err)
		}
		feriados = append(feriados, domain.FeriadoMutationItem{ID: id, Descricao: descricao, Data: data, Feriado: tipo, Ativo: ativo})
	}

	if strings.TrimSpace(input.MunicipioID) != "" {
		_, _ = dbExec(ctx, r.pool, `UPDATE public.feriado_municipal SET municipio_id = $1 WHERE feriado_id = $2`, input.MunicipioID, input.ID)
	}
	if strings.TrimSpace(input.EstadoID) != "" {
		_, _ = dbExec(ctx, r.pool, `UPDATE public.feriado_estadual SET uf_id = $1 WHERE feriado_id = $2`, input.EstadoID, input.ID)
	}

	return feriados, int64(len(feriados)), nil
}

func (r *FeriadoRepository) Delete(ctx context.Context, id string) ([]domain.FeriadoMutationItem, int64, error) {
	const query = `
		UPDATE public.feriados
		SET ativo = false
		WHERE id = $1
		RETURNING id, descricao, data, feriado, ativo`

	rows, err := dbQuery(ctx, r.pool, query, id)
	if err != nil {
		return nil, 0, fmt.Errorf("delete feriado: %w", err)
	}
	defer rows.Close()

	feriados := make([]domain.FeriadoMutationItem, 0)
	for rows.Next() {
		var fid, descricao, data, tipo string
		var ativo bool
		if err := rows.Scan(&fid, &descricao, &data, &tipo, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan deleted feriado: %w", err)
		}
		feriados = append(feriados, domain.FeriadoMutationItem{ID: fid, Descricao: descricao, Data: data, Feriado: tipo, Ativo: ativo})
	}

	return feriados, int64(len(feriados)), nil
}
