package repository

import (
	"context"
	"fmt"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CaixaPostalRepository struct {
	pool *pgxpool.Pool
}

func NewCaixaPostalRepository(pool *pgxpool.Pool) *CaixaPostalRepository {
	return &CaixaPostalRepository{pool: pool}
}

func (r *CaixaPostalRepository) Insert(ctx context.Context, schemaName string, msg domain.CaixaPostalMensagem) error {
	q := fmt.Sprintf(`
		INSERT INTO %s.caixa_postal_mensagens 
		(remetente_id, remetente_nome, tipo, is_global, titulo, conteudo)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, schemaName)

	_, err := r.pool.Exec(ctx, q, msg.RemetenteID, msg.RemetenteNome, msg.Tipo, msg.IsGlobal, msg.Titulo, msg.Conteudo)
	return err
}

func (r *CaixaPostalRepository) List(ctx context.Context, schemaName string) ([]domain.CaixaPostalMensagem, error) {
	q := fmt.Sprintf(`
		SELECT id, remetente_id, remetente_nome, tipo, is_global, titulo, conteudo, lida, lida_por, lida_em, criado_em
		FROM %s.caixa_postal_mensagens
		ORDER BY criado_em DESC
	`, schemaName)

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.CaixaPostalMensagem
	for rows.Next() {
		var m domain.CaixaPostalMensagem
		err := rows.Scan(
			&m.ID, &m.RemetenteID, &m.RemetenteNome, &m.Tipo, &m.IsGlobal,
			&m.Titulo, &m.Conteudo, &m.Lida, &m.LidaPor, &m.LidaEm, &m.CriadoEm,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (r *CaixaPostalRepository) CountUnread(ctx context.Context, schemaName string) (int, error) {
	q := fmt.Sprintf(`
		SELECT count(*)
		FROM %s.caixa_postal_mensagens
		WHERE lida = false AND tipo = 'INBOX'
	`, schemaName)

	var count int
	err := r.pool.QueryRow(ctx, q).Scan(&count)
	return count, err
}

func (r *CaixaPostalRepository) MarkAsRead(ctx context.Context, schemaName, msgID, userID string) error {
	q := fmt.Sprintf(`
		UPDATE %s.caixa_postal_mensagens 
		SET lida = true, lida_por = $1, lida_em = now()
		WHERE id = $2 AND lida = false
	`, schemaName)

	_, err := r.pool.Exec(ctx, q, userID, msgID)
	return err
}

func (r *CaixaPostalRepository) GetAllActiveSchemas(ctx context.Context) ([]string, error) {
	q := `SELECT DISTINCT schema_name FROM public.tenant_schema_catalog WHERE TRIM(schema_name) <> ''`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		schemas = append(schemas, s)
	}
	return schemas, nil
}

func (r *CaixaPostalRepository) GetSchemaByTenantID(ctx context.Context, tenantID string) (string, error) {
	var schema string
	err := r.pool.QueryRow(ctx,
		`SELECT schema_name FROM public.tenant_schema_catalog WHERE tenant_id = $1::uuid`,
		tenantID,
	).Scan(&schema)
	if err != nil {
		return "", fmt.Errorf("schema do tenant (%s) não encontrado: %w", tenantID, err)
	}
	return schema, nil
}

func (r *CaixaPostalRepository) GetSuperSchema(ctx context.Context) (string, error) {
	q := `
		SELECT tsc.schema_name
		FROM public.usuario u
		JOIN public.tenant_schema_catalog tsc ON u.tenantid = tsc.tenant_id
		WHERE UPPER(TRIM(COALESCE(u.role::text, ''))) = 'SUPER' AND TRIM(tsc.schema_name) <> ''
		LIMIT 1
	`
	var schema string
	err := r.pool.QueryRow(ctx, q).Scan(&schema)
	if err != nil {
		return "", fmt.Errorf("schema VEC Sistemas (SUPER) não encontrado")
	}
	return schema, nil
}
