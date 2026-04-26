package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepresentanteRepository struct {
	pool *pgxpool.Pool
}

func NewRepresentanteRepository(pool *pgxpool.Pool) *RepresentanteRepository {
	return &RepresentanteRepository{pool: pool}
}

func (r *RepresentanteRepository) List(ctx context.Context) ([]domain.Representante, error) {
	const q = `
		SELECT id::text, COALESCE(nome, ''), COALESCE(email_contato, ''), COALESCE(ativo, false)
		FROM public.representantes
		ORDER BY nome ASC`
	rows, err := dbQuery(ctx, r.pool, q)
	if err != nil {
		return nil, fmt.Errorf("list representantes: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Representante, 0)
	for rows.Next() {
		var item domain.Representante
		if err := rows.Scan(&item.ID, &item.Nome, &item.EmailContato, &item.Ativo); err != nil {
			return nil, fmt.Errorf("scan representante: %w", err)
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *RepresentanteRepository) Get(ctx context.Context, id string) (domain.Representante, error) {
	const q = `
		SELECT id::text, COALESCE(nome, ''), COALESCE(email_contato, ''), COALESCE(ativo, false)
		FROM public.representantes WHERE id = $1::uuid`
	var item domain.Representante
	if err := dbQueryRow(ctx, r.pool, q, id).Scan(&item.ID, &item.Nome, &item.EmailContato, &item.Ativo); err != nil {
		return domain.Representante{}, fmt.Errorf("get representante: %w", err)
	}
	return item, nil
}

func (r *RepresentanteRepository) Create(ctx context.Context, nome, emailContato string) (domain.Representante, error) {
	const q = `
		INSERT INTO public.representantes (nome, email_contato, ativo)
		VALUES ($1, $2, true)
		RETURNING id::text, COALESCE(nome, ''), COALESCE(email_contato, ''), COALESCE(ativo, false)`
	var item domain.Representante
	if err := dbQueryRow(ctx, r.pool, q, strings.TrimSpace(nome), strings.TrimSpace(emailContato)).Scan(
		&item.ID, &item.Nome, &item.EmailContato, &item.Ativo,
	); err != nil {
		return domain.Representante{}, fmt.Errorf("create representante: %w", err)
	}
	return item, nil
}

func (r *RepresentanteRepository) Update(ctx context.Context, id, nome, emailContato string, ativo bool) (domain.Representante, error) {
	const q = `
		UPDATE public.representantes
		SET nome = $1, email_contato = $2, ativo = $3
		WHERE id = $4::uuid
		RETURNING id::text, COALESCE(nome, ''), COALESCE(email_contato, ''), COALESCE(ativo, false)`
	var item domain.Representante
	if err := dbQueryRow(ctx, r.pool, q, strings.TrimSpace(nome), strings.TrimSpace(emailContato), ativo, id).Scan(
		&item.ID, &item.Nome, &item.EmailContato, &item.Ativo,
	); err != nil {
		return domain.Representante{}, fmt.Errorf("update representante: %w", err)
	}
	return item, nil
}

func (r *RepresentanteRepository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM public.representantes WHERE id = $1::uuid`
	tag, err := dbExec(ctx, r.pool, q, id)
	if err != nil {
		return fmt.Errorf("delete representante: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("representante nao encontrado")
	}
	return nil
}

func (r *RepresentanteRepository) ListModulos(ctx context.Context) ([]domain.ModuloPlataforma, error) {
	const q = `SELECT id::text, slug, nome, ordem FROM public.modulo_plataforma ORDER BY ordem ASC, slug ASC`
	rows, err := dbQuery(ctx, r.pool, q)
	if err != nil {
		return nil, fmt.Errorf("list modulos: %w", err)
	}
	defer rows.Close()

	out := make([]domain.ModuloPlataforma, 0)
	for rows.Next() {
		var m domain.ModuloPlataforma
		if err := rows.Scan(&m.ID, &m.Slug, &m.Nome, &m.Ordem); err != nil {
			return nil, fmt.Errorf("scan modulo: %w", err)
		}
		out = append(out, m)
	}
	return out, nil
}

func (r *RepresentanteRepository) GetMatriz(ctx context.Context, representanteID string) ([]domain.MatrizAcessoItem, error) {
	const q = `
		SELECT m.id::text, m.slug, COALESCE(ma.habilitado, false)
		FROM public.modulo_plataforma m
		LEFT JOIN public.matriz_acesso ma
		  ON ma.modulo_id = m.id AND ma.representante_id = $1::uuid
		ORDER BY m.ordem ASC, m.slug ASC`
	rows, err := dbQuery(ctx, r.pool, q, representanteID)
	if err != nil {
		return nil, fmt.Errorf("matriz: %w", err)
	}
	defer rows.Close()

	out := make([]domain.MatrizAcessoItem, 0)
	for rows.Next() {
		var item domain.MatrizAcessoItem
		if err := rows.Scan(&item.ModuloID, &item.Slug, &item.Habilitado); err != nil {
			return nil, fmt.Errorf("scan matriz: %w", err)
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *RepresentanteRepository) ReplaceMatriz(ctx context.Context, representanteID string, entries []domain.MatrizAcessoItem) error {
	tx, err := dbBeginTx(ctx, r.pool, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx matriz: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM public.matriz_acesso WHERE representante_id = $1::uuid`, representanteID); err != nil {
		return fmt.Errorf("clear matriz: %w", err)
	}

	const ins = `
		INSERT INTO public.matriz_acesso (representante_id, modulo_id, habilitado)
		VALUES ($1::uuid, $2::uuid, $3)`
	for _, e := range entries {
		if strings.TrimSpace(e.ModuloID) == "" {
			continue
		}
		if _, err := tx.Exec(ctx, ins, representanteID, e.ModuloID, e.Habilitado); err != nil {
			return fmt.Errorf("insert matriz: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit matriz: %w", err)
	}
	return nil
}
