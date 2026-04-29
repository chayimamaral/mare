package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GrupoPassosListParams struct {
	First     int
	Rows      int
	SortField string
	SortOrder int
	Descricao string
}

type GrupoPassosRepository struct {
	pool *pgxpool.Pool
}

func NewGrupoPassosRepository(pool *pgxpool.Pool) *GrupoPassosRepository {
	return &GrupoPassosRepository{pool: pool}
}

func (r *GrupoPassosRepository) List(ctx context.Context, params GrupoPassosListParams) ([]domain.GrupoPassosListItem, int64, error) {
	whereParts := []string{"g.ativo = true"}
	args := []any{}
	argIndex := 1

	if strings.TrimSpace(params.Descricao) != "" {
		whereParts = append(whereParts, fmt.Sprintf("g.descricao ILIKE $%d", argIndex))
		args = append(args, "%"+strings.TrimSpace(params.Descricao)+"%")
		argIndex++
	}

	orderBy := "g.descricao ASC"
	switch params.SortField {
	case "municipio":
		if params.SortOrder == -1 {
			orderBy = "m.nome ASC"
		} else {
			orderBy = "m.nome DESC"
		}
	case "descricao":
		if params.SortOrder == -1 {
			orderBy = "g.descricao ASC"
		} else {
			orderBy = "g.descricao DESC"
		}
	case "tipoempresa":
		if params.SortOrder == -1 {
			orderBy = "t.descricao ASC"
		} else {
			orderBy = "t.descricao DESC"
		}
	}

	query := fmt.Sprintf(`
		SELECT
			g.id,
			g.descricao,
			g.municipio_id,
			g.tipoempresa_id,
			m.id,
			m.nome,
			e.sigla,
			t.id,
			t.descricao
		FROM grupopassos g
		JOIN public.municipio m ON m.id = g.municipio_id
		JOIN public.tipoempresa t ON t.id = g.tipoempresa_id
		JOIN public.estado e ON e.id = m.ufid
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`, strings.Join(whereParts, " AND "), orderBy, argIndex, argIndex+1)
	args = append(args, params.Rows, params.First)

	rows, err := dbQuery(ctx, r.pool, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list grupopassos: %w", err)
	}
	defer rows.Close()

	grupos := make([]domain.GrupoPassosListItem, 0)
	for rows.Next() {
		var id, descricao, municipioID, tipoEmpresaID, mid, mnome, sigla, tid, tdesc string
		if err := rows.Scan(&id, &descricao, &municipioID, &tipoEmpresaID, &mid, &mnome, &sigla, &tid, &tdesc); err != nil {
			return nil, 0, fmt.Errorf("scan grupopassos: %w", err)
		}

		grupos = append(grupos, domain.GrupoPassosListItem{
			ID:            id,
			Descricao:     descricao,
			MunicipioID:   municipioID,
			TipoEmpresaID: tipoEmpresaID,
			Municipio: domain.GrupoPassosMunicipio{
				ID:   mid,
				Nome: mnome + " / " + sigla,
			},
			TipoEmpresa: domain.GrupoPassosTipoEmpresa{
				ID:        tid,
				Descricao: tdesc,
			},
		})
	}

	countQuery := fmt.Sprintf("SELECT count(*) FROM grupopassos g WHERE %s", strings.Join(whereParts, " AND "))
	var total int64
	if err := dbQueryRow(ctx, r.pool, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count grupopassos: %w", err)
	}

	return grupos, total, nil
}

func (r *GrupoPassosRepository) Create(ctx context.Context, descricao, municipioID, tipoEmpresaID string) ([]domain.GrupoPassosMutationItem, int64, error) {
	const query = `
		INSERT INTO grupopassos (descricao, municipio_id, tipoempresa_id)
		VALUES ($1, $2, $3)
		RETURNING id, descricao, municipio_id, tipoempresa_id, ativo`

	rows, err := dbQuery(ctx, r.pool, query, descricao, municipioID, tipoEmpresaID)
	if err != nil {
		return nil, 0, fmt.Errorf("create grupopassos: %w", err)
	}
	defer rows.Close()

	grupos := make([]domain.GrupoPassosMutationItem, 0)
	for rows.Next() {
		var id, d, m, t string
		var ativo bool
		if err := rows.Scan(&id, &d, &m, &t, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan created grupopassos: %w", err)
		}
		grupos = append(grupos, domain.GrupoPassosMutationItem{ID: id, Descricao: d, MunicipioID: m, TipoEmpresaID: t, Ativo: ativo})
	}

	var total int64
	if err := dbQueryRow(ctx, r.pool, `SELECT count(*) FROM grupopassos WHERE ativo = true`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count grupopassos: %w", err)
	}

	return grupos, total, nil
}

func (r *GrupoPassosRepository) Update(ctx context.Context, id, descricao, municipioID, tipoEmpresaID string) ([]domain.GrupoPassosMutationItem, int64, error) {
	const query = `
		UPDATE grupopassos
		SET descricao = $1, municipio_id = $2, tipoempresa_id = $3
		WHERE id = $4
		RETURNING id, descricao, municipio_id, tipoempresa_id, ativo`

	rows, err := dbQuery(ctx, r.pool, query, descricao, municipioID, tipoEmpresaID, id)
	if err != nil {
		return nil, 0, fmt.Errorf("update grupopassos: %w", err)
	}
	defer rows.Close()

	grupos := make([]domain.GrupoPassosMutationItem, 0)
	for rows.Next() {
		var gid, d, m, t string
		var ativo bool
		if err := rows.Scan(&gid, &d, &m, &t, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan updated grupopassos: %w", err)
		}
		grupos = append(grupos, domain.GrupoPassosMutationItem{ID: gid, Descricao: d, MunicipioID: m, TipoEmpresaID: t, Ativo: ativo})
	}

	var total int64
	if err := dbQueryRow(ctx, r.pool, `SELECT count(*) FROM grupopassos WHERE ativo = true`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count grupopassos: %w", err)
	}

	return grupos, total, nil
}

func (r *GrupoPassosRepository) Delete(ctx context.Context, id string) ([]domain.GrupoPassosMutationItem, int64, error) {
	const query = `
		UPDATE grupopassos
		SET ativo = false
		WHERE id = $1
		RETURNING id, descricao, municipio_id, tipoempresa_id, ativo`

	rows, err := dbQuery(ctx, r.pool, query, id)
	if err != nil {
		return nil, 0, fmt.Errorf("delete grupopassos: %w", err)
	}
	defer rows.Close()

	grupos := make([]domain.GrupoPassosMutationItem, 0)
	for rows.Next() {
		var gid, d, m, t string
		var ativo bool
		if err := rows.Scan(&gid, &d, &m, &t, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan deleted grupopassos: %w", err)
		}
		grupos = append(grupos, domain.GrupoPassosMutationItem{ID: gid, Descricao: d, MunicipioID: m, TipoEmpresaID: t, Ativo: ativo})
	}

	var total int64
	if err := dbQueryRow(ctx, r.pool, `SELECT count(*) FROM grupopassos WHERE ativo = true`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count grupopassos: %w", err)
	}

	return grupos, total, nil
}

func (r *GrupoPassosRepository) GetByID(ctx context.Context, id string) ([]domain.GrupoPassosMutationItem, int64, error) {
	const query = `
		SELECT id, descricao, municipio_id, tipoempresa_id, ativo
		FROM grupopassos
		WHERE id = $1`

	rows, err := dbQuery(ctx, r.pool, query, id)
	if err != nil {
		return nil, 0, fmt.Errorf("get grupopassos by id: %w", err)
	}
	defer rows.Close()

	grupos := make([]domain.GrupoPassosMutationItem, 0)
	for rows.Next() {
		var gid, d, m, t string
		var ativo bool
		if err := rows.Scan(&gid, &d, &m, &t, &ativo); err != nil {
			return nil, 0, fmt.Errorf("scan grupopassos by id: %w", err)
		}
		grupos = append(grupos, domain.GrupoPassosMutationItem{ID: gid, Descricao: d, MunicipioID: m, TipoEmpresaID: t, Ativo: ativo})
	}

	var total int64
	if err := dbQueryRow(ctx, r.pool, `SELECT count(*) FROM grupopassos WHERE ativo = true`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count grupopassos: %w", err)
	}

	return grupos, total, nil
}
