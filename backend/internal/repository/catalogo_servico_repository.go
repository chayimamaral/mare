package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogoServicoRepository struct {
	pool *pgxpool.Pool
}

type CatalogoServicoUpsertInput struct {
	ID                  string
	Secao               string
	Sequencial          int
	Codigo              string
	IDSistema           string
	IDServico           string
	SituacaoImplantacao string
	DataImplantacao     string
	Tipo                string
	Descricao           string
}

func NewCatalogoServicoRepository(pool *pgxpool.Pool) *CatalogoServicoRepository {
	return &CatalogoServicoRepository{pool: pool}
}

func (r *CatalogoServicoRepository) List(ctx context.Context, secao string) ([]domain.CatalogoServico, error) {
	args := []any{}
	where := "ativo = true"
	if strings.TrimSpace(secao) != "" && strings.ToUpper(strings.TrimSpace(secao)) != "TODAS" {
		where += " AND secao = $1"
		args = append(args, strings.TrimSpace(secao))
	}

	query := fmt.Sprintf(`
		SELECT id, secao, sequencial, codigo, id_sistema, id_servico, situacao_implantacao,
		       COALESCE(to_char(data_implantacao, 'YYYY-MM-DD'), ''), tipo, descricao, ativo
		  FROM public.catalogo_servico_integra_contador
		 WHERE %s
		 ORDER BY secao ASC, sequencial ASC, codigo ASC`, where)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list catalogo_servico: %w", err)
	}
	defer rows.Close()

	out := make([]domain.CatalogoServico, 0)
	for rows.Next() {
		var item domain.CatalogoServico
		if err := rows.Scan(
			&item.ID,
			&item.Secao,
			&item.Sequencial,
			&item.Codigo,
			&item.IDSistema,
			&item.IDServico,
			&item.SituacaoImplantacao,
			&item.DataImplantacao,
			&item.Tipo,
			&item.Descricao,
			&item.Ativo,
		); err != nil {
			return nil, fmt.Errorf("scan catalogo_servico: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows catalogo_servico: %w", err)
	}
	return out, nil
}

func (r *CatalogoServicoRepository) Create(ctx context.Context, input CatalogoServicoUpsertInput) (domain.CatalogoServico, error) {
	const q = `
		INSERT INTO public.catalogo_servico_integra_contador
			(secao, sequencial, codigo, id_sistema, id_servico, situacao_implantacao, data_implantacao, tipo, descricao)
		VALUES
			($1,$2,$3,$4,$5,$6,NULLIF($7, '')::date,$8,$9)
		RETURNING id, secao, sequencial, codigo, id_sistema, id_servico, situacao_implantacao,
		          COALESCE(to_char(data_implantacao, 'YYYY-MM-DD'), ''), tipo, descricao, ativo`

	var item domain.CatalogoServico
	if err := r.pool.QueryRow(
		ctx,
		q,
		input.Secao,
		input.Sequencial,
		input.Codigo,
		input.IDSistema,
		input.IDServico,
		input.SituacaoImplantacao,
		input.DataImplantacao,
		input.Tipo,
		input.Descricao,
	).Scan(
		&item.ID,
		&item.Secao,
		&item.Sequencial,
		&item.Codigo,
		&item.IDSistema,
		&item.IDServico,
		&item.SituacaoImplantacao,
		&item.DataImplantacao,
		&item.Tipo,
		&item.Descricao,
		&item.Ativo,
	); err != nil {
		return domain.CatalogoServico{}, fmt.Errorf("create catalogo_servico: %w", err)
	}
	return item, nil
}

func (r *CatalogoServicoRepository) Update(ctx context.Context, input CatalogoServicoUpsertInput) (domain.CatalogoServico, error) {
	const q = `
		UPDATE public.catalogo_servico_integra_contador
		   SET secao = $1, sequencial = $2, codigo = $3, id_sistema = $4, id_servico = $5,
		       situacao_implantacao = $6, data_implantacao = NULLIF($7, '')::date, tipo = $8, descricao = $9,
		       atualizado_em = now()
		 WHERE id = $10 AND ativo = true
		RETURNING id, secao, sequencial, codigo, id_sistema, id_servico, situacao_implantacao,
		          COALESCE(to_char(data_implantacao, 'YYYY-MM-DD'), ''), tipo, descricao, ativo`

	var item domain.CatalogoServico
	if err := r.pool.QueryRow(
		ctx,
		q,
		input.Secao,
		input.Sequencial,
		input.Codigo,
		input.IDSistema,
		input.IDServico,
		input.SituacaoImplantacao,
		input.DataImplantacao,
		input.Tipo,
		input.Descricao,
		input.ID,
	).Scan(
		&item.ID,
		&item.Secao,
		&item.Sequencial,
		&item.Codigo,
		&item.IDSistema,
		&item.IDServico,
		&item.SituacaoImplantacao,
		&item.DataImplantacao,
		&item.Tipo,
		&item.Descricao,
		&item.Ativo,
	); err != nil {
		return domain.CatalogoServico{}, fmt.Errorf("update catalogo_servico: %w", err)
	}
	return item, nil
}

func (r *CatalogoServicoRepository) Delete(ctx context.Context, id string) error {
	const q = `
		UPDATE public.catalogo_servico_integra_contador
		   SET ativo = false, atualizado_em = now()
		 WHERE id = $1 AND ativo = true`
	ct, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete catalogo_servico: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("registro nao encontrado")
	}
	return nil
}
