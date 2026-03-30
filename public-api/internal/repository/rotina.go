package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Rotina struct {
	pool *pgxpool.Pool
}

func NewRotina(pool *pgxpool.Pool) *Rotina {
	return &Rotina{pool: pool}
}

// ListByMunicipioTipoJSON retorna rotinas ativas com rotinaitens em uma única query (json_agg / json_build_object).
//
// tipoEmpresaID: UUID ou os literais "null" / "none" para rotinas com tipo_empresa_id nulo.
func (r *Rotina) ListByMunicipioTipoJSON(ctx context.Context, municipioID, tipoEmpresaID string) (json.RawMessage, error) {
	tipo := strings.TrimSpace(tipoEmpresaID)

	const q = `
WITH params AS (
  SELECT 
    $1::uuid AS mun_id,
    lower(trim($2)) AS tipo_str
)
SELECT COALESCE(
  (
    SELECT json_agg(doc ORDER BY doc->>'descricao')
    FROM (
      SELECT json_build_object(
        'id', r.id,
        'descricao', r.descricao,
        'municipio_id', r.municipio_id,
        'tipo_empresa_id', r.tipo_empresa_id,
        'rotinaitens', COALESCE(
          (
            SELECT json_agg(
              json_build_object(
                'id', p.id,
                'descricao', p.descricao,
                'tempoestimado', p.tempoestimado,
                'ordem', ri.ordem,
                'link', COALESCE(l.link, '')
              )
              ORDER BY ri.ordem ASC NULLS LAST
            )
            FROM public.rotinaitens ri
            INNER JOIN public.passos p ON p.id = ri.passo_id
            LEFT JOIN public.linkpassos l ON l.passo_id = p.id
            WHERE ri.rotina_id = r.id
          ),
          '[]'::json
        )
      ) AS doc
      FROM public.rotinas r, params p
      WHERE r.ativo = true
        AND r.municipio_id = p.mun_id
        AND (
          (p.tipo_str IN ('null', 'none', '') AND r.tipo_empresa_id IS NULL)
          OR
          (p.tipo_str NOT IN ('null', 'none', '') AND r.tipo_empresa_id = p.tipo_str::uuid)
        )
    ) sub
  ),
  '[]'::json
)
`

	var raw []byte
	if err := r.pool.QueryRow(ctx, q, municipioID, tipo).Scan(&raw); err != nil {
		return nil, fmt.Errorf("list rotinas json: %w", err)
	}
	return json.RawMessage(raw), nil
}
