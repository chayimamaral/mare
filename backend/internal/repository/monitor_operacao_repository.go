package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MonitorOperacaoRepository struct {
	pool *pgxpool.Pool
}

type MonitorOperacaoListFilter struct {
	ClienteNome string
	Status      string
	DataDeISO   string
	DataAteISO  string
}

type monitorSchemaTarget struct {
	TenantID   string
	TenantNome string
	SchemaName string
}

func NewMonitorOperacaoRepository(pool *pgxpool.Pool) *MonitorOperacaoRepository {
	return &MonitorOperacaoRepository{pool: pool}
}

type MonitorOperacaoInsert struct {
	TenantID string
	UserID   *string
	Origem   string
	Tipo     string
	Status   string
	Mensagem *string
	Detalhe  map[string]any
}

func (r *MonitorOperacaoRepository) Insert(ctx context.Context, row MonitorOperacaoInsert) (string, error) {
	tid := strings.TrimSpace(row.TenantID)
	if tid == "" {
		return "", fmt.Errorf("tenant_id obrigatorio")
	}
	targetSchema, err := r.resolveTargetSchema(ctx, tid)
	if err != nil {
		return "", err
	}
	var detJSON []byte
	if row.Detalhe != nil {
		var err error
		detJSON, err = json.Marshal(row.Detalhe)
		if err != nil {
			return "", fmt.Errorf("marshal detalhe monitor_operacao: %w", err)
		}
	}

	const q = `
INSERT INTO monitor_operacao (tenant_id, user_id, origem, tipo, status, mensagem, detalhe)
VALUES ($1::uuid, $2, $3, $4, $5, $6, $7::jsonb)
RETURNING id::text
`
	var id string
	err = withTenantSchemaContextByName(ctx, r.pool, targetSchema, func(inner context.Context) error {
		return dbQueryRow(inner, r.pool, q,
			tid,
			row.UserID,
			strings.TrimSpace(row.Origem),
			strings.TrimSpace(row.Tipo),
			strings.TrimSpace(row.Status),
			row.Mensagem,
			detJSON,
		).Scan(&id)
	})
	if err != nil {
		return "", fmt.Errorf("insert monitor_operacao: %w", err)
	}
	return id, nil
}

func (r *MonitorOperacaoRepository) InsertCompromissosRefs(ctx context.Context, monitorID string, compromissoIDs []string) error {
	mid := strings.TrimSpace(monitorID)
	if mid == "" || len(compromissoIDs) == 0 {
		return nil
	}
	targetSchema, err := r.findMonitorSchemaByID(ctx, mid)
	if err != nil {
		return err
	}
	const q = `
INSERT INTO monitor_operacao_compromisso (monitor_operacao_id, empresa_compromisso_id)
VALUES ($1::uuid, $2::uuid)
ON CONFLICT DO NOTHING`
	if err := withTenantSchemaContextByName(ctx, r.pool, targetSchema, func(inner context.Context) error {
		for _, cid := range compromissoIDs {
			cid = strings.TrimSpace(cid)
			if cid == "" {
				continue
			}
			if _, err := dbExec(inner, r.pool, q, mid, cid); err != nil {
				return fmt.Errorf("insert monitor_operacao_compromisso: %w", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *MonitorOperacaoRepository) CountList(ctx context.Context, viewerRole, viewerTenantID string, f MonitorOperacaoListFilter) (int64, error) {
	items, err := r.loadMatchedItems(ctx, viewerRole, viewerTenantID, f)
	if err != nil {
		return 0, err
	}
	return int64(len(items)), nil
}

func (r *MonitorOperacaoRepository) ListPage(ctx context.Context, viewerRole, viewerTenantID string, limit, offset int, f MonitorOperacaoListFilter) ([]domain.MonitorOperacaoItem, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}

	items, err := r.loadMatchedItems(ctx, viewerRole, viewerTenantID, f)
	if err != nil {
		return nil, err
	}
	if offset >= len(items) {
		return []domain.MonitorOperacaoItem{}, nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end], nil
}

func (r *MonitorOperacaoRepository) loadMatchedItems(ctx context.Context, viewerRole, viewerTenantID string, f MonitorOperacaoListFilter) ([]domain.MonitorOperacaoItem, error) {
	targets, err := r.listMonitorTargets(ctx, viewerRole, viewerTenantID)
	if err != nil {
		return nil, err
	}
	all := make([]domain.MonitorOperacaoItem, 0)
	for _, target := range targets {
		items, err := r.loadItemsFromSchema(ctx, target, f)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
	}
	sort.SliceStable(all, func(i, j int) bool {
		if all[i].CriadoEm.Equal(all[j].CriadoEm) {
			return all[i].ID > all[j].ID
		}
		return all[i].CriadoEm.After(all[j].CriadoEm)
	})
	return all, nil
}

func (r *MonitorOperacaoRepository) listMonitorTargets(ctx context.Context, viewerRole, viewerTenantID string) ([]monitorSchemaTarget, error) {
	role := strings.TrimSpace(strings.ToUpper(viewerRole))
	if role == "ADMIN" {
		schemaName, err := r.resolveTargetSchema(ctx, viewerTenantID)
		if err != nil {
			return nil, err
		}
		var tenantNome string
		_ = dbQueryRow(ctx, r.pool, `SELECT COALESCE(nome, '') FROM public.tenant WHERE id = $1::uuid`, viewerTenantID).Scan(&tenantNome)
		return []monitorSchemaTarget{{TenantID: viewerTenantID, TenantNome: tenantNome, SchemaName: schemaName}}, nil
	}

	rows, err := dbQuery(ctx, r.pool, `
		SELECT t.id::text, COALESCE(t.nome, ''), tsc.schema_name
		FROM public.tenant t
		JOIN public.tenant_schema_catalog tsc ON tsc.tenant_id = t.id
		WHERE NULLIF(BTRIM(COALESCE(tsc.schema_name, '')), '') IS NOT NULL
		ORDER BY COALESCE(t.nome, '')`)
	if err != nil {
		return nil, fmt.Errorf("list monitor targets: %w", err)
	}
	defer rows.Close()

	targets := make([]monitorSchemaTarget, 0)
	for rows.Next() {
		var target monitorSchemaTarget
		if err := rows.Scan(&target.TenantID, &target.TenantNome, &target.SchemaName); err != nil {
			return nil, fmt.Errorf("scan monitor target: %w", err)
		}
		targets = append(targets, target)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return targets, nil
}

func (r *MonitorOperacaoRepository) loadItemsFromSchema(ctx context.Context, target monitorSchemaTarget, f MonitorOperacaoListFilter) ([]domain.MonitorOperacaoItem, error) {
	items := make([]domain.MonitorOperacaoItem, 0)
	err := withTenantSchemaContextByName(ctx, r.pool, target.SchemaName, func(inner context.Context) error {
		args := make([]any, 0, 5)
		q := `
SELECT
  mo.id::text,
  mo.tenant_id::text,
  NULLIF(BTRIM(CASE WHEN mo.tenant_id = $1::uuid THEN 'Plataforma' ELSE COALESCE(t.nome, '') END), '') AS tenant_nome,
  c.nome,
  mo.user_id::text,
  mo.origem,
  mo.tipo,
  mo.status,
  mo.mensagem,
  mo.detalhe,
  mo.criado_em
FROM monitor_operacao mo
LEFT JOIN public.tenant t ON t.id = mo.tenant_id
LEFT JOIN empresa e
  ON e.id = CASE
    WHEN COALESCE(mo.detalhe->>'empresa_id', '') ~* '^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
      THEN (mo.detalhe->>'empresa_id')::uuid
    ELSE NULL
  END
LEFT JOIN cliente c ON c.id = e.cliente_id
WHERE 1=1`
		args = append(args, domain.MonitorOperacaoTenantPlataformaID)
		if clienteNome := strings.TrimSpace(f.ClienteNome); clienteNome != "" {
			q += fmt.Sprintf(" AND c.nome ILIKE $%d", len(args)+1)
			args = append(args, "%"+clienteNome+"%")
		}
		if status := strings.TrimSpace(strings.ToUpper(f.Status)); status != "" {
			q += fmt.Sprintf(" AND UPPER(TRIM(mo.status)) = $%d", len(args)+1)
			args = append(args, status)
		}
		if dataDe := strings.TrimSpace(f.DataDeISO); dataDe != "" {
			q += fmt.Sprintf(" AND mo.criado_em::date >= $%d::date", len(args)+1)
			args = append(args, dataDe)
		}
		if dataAte := strings.TrimSpace(f.DataAteISO); dataAte != "" {
			q += fmt.Sprintf(" AND mo.criado_em::date <= $%d::date", len(args)+1)
			args = append(args, dataAte)
		}
		q += ` ORDER BY mo.criado_em DESC`

		rows, err := dbQuery(inner, r.pool, q, args...)
		if err != nil {
			return fmt.Errorf("list monitor_operacao: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var it domain.MonitorOperacaoItem
			var tenantNome, clienteNome, userID, mensagem *string
			var detBytes []byte
			if err := rows.Scan(
				&it.ID,
				&it.TenantID,
				&tenantNome,
				&clienteNome,
				&userID,
				&it.Origem,
				&it.Tipo,
				&it.Status,
				&mensagem,
				&detBytes,
				&it.CriadoEm,
			); err != nil {
				return fmt.Errorf("scan monitor_operacao: %w", err)
			}
			it.TenantNome = tenantNome
			it.ClienteNome = clienteNome
			it.UserID = userID
			it.Mensagem = mensagem
			if len(detBytes) > 0 {
				var m map[string]any
				if err := json.Unmarshal(detBytes, &m); err == nil {
					it.Detalhe = m
				}
			}
			items = append(items, it)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return r.attachCompromissosForSchema(inner, items)
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MonitorOperacaoRepository) attachCompromissosForSchema(ctx context.Context, itens []domain.MonitorOperacaoItem) error {
	if len(itens) == 0 {
		return nil
	}
	monitorIDs := make([]string, 0, len(itens))
	index := make(map[string]int, len(itens))
	for i, it := range itens {
		monitorIDs = append(monitorIDs, it.ID)
		index[it.ID] = i
	}
	const q = `
SELECT
  rel.monitor_operacao_id::text,
  ec.id::text,
  ec.empresa_id::text,
  c.nome,
  ec.descricao,
  ec.competencia::date::text,
  ec.vencimento::date::text,
  ec.status,
  ec.valor
FROM monitor_operacao_compromisso rel
INNER JOIN empresa_compromissos ec ON ec.id = rel.empresa_compromisso_id
LEFT JOIN empresa e ON e.id = ec.empresa_id
LEFT JOIN cliente c ON c.id = e.cliente_id
WHERE rel.monitor_operacao_id = ANY($1::uuid[])
ORDER BY rel.criado_em ASC, ec.vencimento ASC, ec.descricao ASC`
	rows, err := dbQuery(ctx, r.pool, q, monitorIDs)
	if err != nil {
		return fmt.Errorf("list monitor_operacao_compromisso: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var monitorID string
		var child domain.MonitorOperacaoCompromissoItem
		var clienteNome *string
		var valor *float64
		if err := rows.Scan(
			&monitorID,
			&child.CompromissoID,
			&child.EmpresaID,
			&clienteNome,
			&child.Descricao,
			&child.Competencia,
			&child.Vencimento,
			&child.Status,
			&valor,
		); err != nil {
			return fmt.Errorf("scan monitor_operacao_compromisso: %w", err)
		}
		child.ClienteNome = clienteNome
		child.Valor = valor
		if idx, ok := index[monitorID]; ok {
			itens[idx].Compromissos = append(itens[idx].Compromissos, child)
		}
	}
	return rows.Err()
}

func (r *MonitorOperacaoRepository) resolveTargetSchema(ctx context.Context, tenantID string) (string, error) {
	if strings.TrimSpace(tenantID) == domain.MonitorOperacaoTenantPlataformaID {
		return r.getSuperSchema(ctx)
	}
	return repositoryResolveTenantSchema(ctx, r.pool, tenantID)
}

func (r *MonitorOperacaoRepository) findMonitorSchemaByID(ctx context.Context, monitorID string) (string, error) {
	rows, err := dbQuery(ctx, r.pool, `SELECT DISTINCT schema_name FROM public.tenant_schema_catalog WHERE NULLIF(BTRIM(schema_name), '') IS NOT NULL`)
	if err != nil {
		return "", fmt.Errorf("list schemas monitor_operacao: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return "", err
		}
		var found bool
		if err := withTenantSchemaContextByName(ctx, r.pool, schemaName, func(inner context.Context) error {
			return dbQueryRow(inner, r.pool, `SELECT EXISTS(SELECT 1 FROM monitor_operacao WHERE id = $1::uuid)`, monitorID).Scan(&found)
		}); err != nil {
			return "", err
		}
		if found {
			return schemaName, nil
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("monitor_operacao nao encontrado")
}

func (r *MonitorOperacaoRepository) getSuperSchema(ctx context.Context) (string, error) {
	const q = `
		SELECT tsc.schema_name
		FROM public.usuario u
		JOIN public.tenant_schema_catalog tsc ON tsc.tenant_id = u.tenantid
		WHERE UPPER(TRIM(COALESCE(u.role::text, ''))) = 'SUPER'
		  AND NULLIF(BTRIM(COALESCE(tsc.schema_name, '')), '') IS NOT NULL
		LIMIT 1`
	var schemaName string
	if err := dbQueryRow(ctx, r.pool, q).Scan(&schemaName); err != nil {
		return "", fmt.Errorf("schema SUPER nao encontrado: %w", err)
	}
	return schemaName, nil
}
