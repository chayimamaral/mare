package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CaixaPostalRepository struct {
	pool *pgxpool.Pool
}

func NewCaixaPostalRepository(pool *pgxpool.Pool) *CaixaPostalRepository {
	return &CaixaPostalRepository{pool: pool}
}

var tenantSchemaRegex = regexp.MustCompile(`^[a-z][a-z0-9_]{2,62}$`)

func normalizeSchemaName(schemaName string) (string, error) {
	s := strings.TrimSpace(strings.ToLower(schemaName))
	if !tenantSchemaRegex.MatchString(s) {
		return "", fmt.Errorf("schema invalido")
	}
	return s, nil
}

func quoteIdentLocal(ident string) string {
	return `"` + strings.ReplaceAll(ident, `"`, `""`) + `"`
}

func tenantTable(schemaName, tableName string) (string, error) {
	normalized, err := normalizeSchemaName(schemaName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", quoteIdentLocal(normalized), quoteIdentLocal(tableName)), nil
}

func (r *CaixaPostalRepository) Insert(ctx context.Context, schemaName string, msg domain.CaixaPostalMensagem) error {
	tableName, err := tenantTable(schemaName, "caixa_postal_mensagens")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`
		INSERT INTO %s
		(remetente_id, remetente_nome, tipo, is_global, titulo, conteudo)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, tableName)

	_, err = dbExec(ctx, r.pool, q, msg.RemetenteID, msg.RemetenteNome, msg.Tipo, msg.IsGlobal, msg.Titulo, msg.Conteudo)
	return err
}

func (r *CaixaPostalRepository) List(ctx context.Context, schemaName, tenantID string) ([]domain.CaixaPostalMensagem, error) {
	tableName, err := tenantTable(schemaName, "caixa_postal_mensagens")
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`
		SELECT cpm.id, cpm.remetente_id, u.tenantid AS remetente_tenantid,
		       cpm.remetente_nome,
		       CASE
		           WHEN cpm.tipo = 'OUTBOX' AND COALESCE(u.tenantid::text, '') <> $1 THEN 'INBOX'
		           WHEN cpm.tipo = 'INBOX' AND COALESCE(u.tenantid::text, '') = $1 THEN 'OUTBOX'
		           ELSE cpm.tipo
		       END AS tipo,
		       cpm.is_global, cpm.titulo, cpm.conteudo,
		       cpm.lida, cpm.lida_por, cpm.lida_em, cpm.criado_em
		FROM %s cpm
		LEFT JOIN public.usuario u ON u.id = cpm.remetente_id
		ORDER BY cpm.criado_em DESC
	`, tableName)

	rows, err := dbQuery(ctx, r.pool, q, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.CaixaPostalMensagem
	for rows.Next() {
		var m domain.CaixaPostalMensagem
		err := rows.Scan(
			&m.ID, &m.RemetenteID, &m.RemetenteTenantID, &m.RemetenteNome, &m.Tipo, &m.IsGlobal,
			&m.Titulo, &m.Conteudo, &m.Lida, &m.LidaPor, &m.LidaEm, &m.CriadoEm,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (r *CaixaPostalRepository) CountUnread(ctx context.Context, schemaName, tenantID string) (int, error) {
	tableName, err := tenantTable(schemaName, "caixa_postal_mensagens")
	if err != nil {
		return 0, err
	}
	q := fmt.Sprintf(`
		SELECT count(*)
		FROM %s cpm
		LEFT JOIN public.usuario u ON u.id = cpm.remetente_id
		WHERE cpm.lida = false
		  AND (
		       CASE
		           WHEN cpm.tipo = 'OUTBOX' AND COALESCE(u.tenantid::text, '') <> $1 THEN 'INBOX'
		           WHEN cpm.tipo = 'INBOX' AND COALESCE(u.tenantid::text, '') = $1 THEN 'OUTBOX'
		           ELSE cpm.tipo
		       END
		  ) = 'INBOX'
	`, tableName)

	var count int
	err = dbQueryRow(ctx, r.pool, q, tenantID).Scan(&count)
	return count, err
}

func (r *CaixaPostalRepository) MarkAsRead(ctx context.Context, schemaName, msgID, userID string) error {
	tableName, err := tenantTable(schemaName, "caixa_postal_mensagens")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`
		UPDATE %s
		SET lida = true, lida_por = $1, lida_em = now()
		WHERE id = $2 AND lida = false
	`, tableName)

	_, err = dbExec(ctx, r.pool, q, userID, msgID)
	return err
}

func (r *CaixaPostalRepository) GetAllActiveSchemas(ctx context.Context) ([]string, error) {
	q := `SELECT DISTINCT schema_name FROM public.tenant_schema_catalog WHERE TRIM(schema_name) <> ''`
	rows, err := dbQuery(ctx, r.pool, q)
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
	err := dbQueryRow(ctx, r.pool,
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
	err := dbQueryRow(ctx, r.pool, q).Scan(&schema)
	if err != nil {
		return "", fmt.Errorf("schema VEC Sistemas (SUPER) não encontrado")
	}
	return schema, nil
}
