package repository

import (
	"context"
	"fmt"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NodeRepository struct {
	pool *pgxpool.Pool
}

func NewNodeRepository(pool *pgxpool.Pool) *NodeRepository {
	return &NodeRepository{pool: pool}
}

// Nodes returns a flat list via recursive CTE (mirrors NodeService.nodes).
func (r *NodeRepository) Nodes(ctx context.Context) ([]domain.NodePasso, error) {
	const query = `
		WITH RECURSIVE tree AS (
			SELECT id, descricao, parent_id
			FROM passos
			WHERE parent_id IS NULL
			UNION ALL
			SELECT t.id, t.descricao, t.parent_id
			FROM passos t
			JOIN tree ON t.parent_id = tree.id
		)
		SELECT id, descricao, parent_id FROM tree`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("nodes query: %w", err)
	}
	defer rows.Close()

	passos := make([]domain.NodePasso, 0)
	for rows.Next() {
		var p domain.NodePasso
		if err := rows.Scan(&p.ID, &p.Descricao, &p.ParentID); err != nil {
			return nil, fmt.Errorf("nodes scan: %w", err)
		}
		passos = append(passos, p)
	}

	return passos, nil
}

// Family calls the PG stored function get_passos_nested() and returns raw JSON.
func (r *NodeRepository) Family(ctx context.Context) (any, error) {
	const query = `SELECT get_passos_nested() AS dados`

	var dados any
	if err := r.pool.QueryRow(ctx, query).Scan(&dados); err != nil {
		return nil, fmt.Errorf("family query: %w", err)
	}

	return dados, nil
}

// Recurso returns all passos as a flat list ready for tree building.
func (r *NodeRepository) Recurso(ctx context.Context) ([]domain.NodePasso, error) {
	const query = `SELECT id, descricao, parent_id FROM passos`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("recurso query: %w", err)
	}
	defer rows.Close()

	passos := make([]domain.NodePasso, 0)
	for rows.Next() {
		var p domain.NodePasso
		if err := rows.Scan(&p.ID, &p.Descricao, &p.ParentID); err != nil {
			return nil, fmt.Errorf("recurso scan: %w", err)
		}
		passos = append(passos, p)
	}

	return passos, nil
}
