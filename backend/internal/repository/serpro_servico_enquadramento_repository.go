package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SerproServicoEnquadramentoRepository struct {
	pool *pgxpool.Pool
}

func NewSerproServicoEnquadramentoRepository(pool *pgxpool.Pool) *SerproServicoEnquadramentoRepository {
	return &SerproServicoEnquadramentoRepository{pool: pool}
}

func (r *SerproServicoEnquadramentoRepository) ListServicosIDs(ctx context.Context, enquadramentoID, regimeTributarioID string) ([]string, error) {
	rows, err := dbQuery(ctx, r.pool, `
		SELECT serpro_servico_id::text
		FROM public.serpro_servico_enquadramento
		WHERE enquadramento_id = $1::uuid
		  AND regime_tributario_id = $2::uuid
		ORDER BY serpro_servico_id::text ASC
	`, strings.TrimSpace(enquadramentoID), strings.TrimSpace(regimeTributarioID))
	if err != nil {
		return nil, fmt.Errorf("listar matriz de conformidade fiscal: %w", err)
	}
	defer rows.Close()

	out := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan matriz de conformidade fiscal: %w", err)
		}
		out = append(out, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows matriz de conformidade fiscal: %w", err)
	}
	return out, nil
}

func (r *SerproServicoEnquadramentoRepository) SaveServicosIDs(ctx context.Context, enquadramentoID, regimeTributarioID string, servicosIDs []string) error {
	tx, err := dbBegin(ctx, r.pool)
	if err != nil {
		return fmt.Errorf("begin salvar matriz de conformidade fiscal: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		DELETE FROM public.serpro_servico_enquadramento
		WHERE enquadramento_id = $1::uuid
		  AND regime_tributario_id = $2::uuid
	`, strings.TrimSpace(enquadramentoID), strings.TrimSpace(regimeTributarioID))
	if err != nil {
		return fmt.Errorf("limpar matriz de conformidade fiscal: %w", err)
	}

	for _, raw := range servicosIDs {
		sid := strings.TrimSpace(raw)
		if sid == "" {
			continue
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO public.serpro_servico_enquadramento (serpro_servico_id, enquadramento_id, regime_tributario_id)
			VALUES ($1::uuid, $2::uuid, $3::uuid)
		`, sid, strings.TrimSpace(enquadramentoID), strings.TrimSpace(regimeTributarioID)); err != nil {
			return fmt.Errorf("inserir matriz de conformidade fiscal: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit salvar matriz de conformidade fiscal: %w", err)
	}
	return nil
}
