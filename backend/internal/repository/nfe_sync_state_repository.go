package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
)

type NFESyncStateUpsert struct {
	Provider            string
	UF                  string
	CNPJ                string
	UltimoNSU           string
	UltimoCStat         int
	UltimoMotivo        string
	UltimaVerificacao   time.Time
	ProximaConsultaApos *time.Time
	UltimaQtDFeRet      int
}

type NFESyncEstadoListParams struct {
	First    int
	Rows     int
	Provider string
	UF       string
	CNPJ     string
}

func tenantNFESyncStateTable(schemaName string) (string, error) {
	s, err := normalizeSchemaForNFE(schemaName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`"%s"."nfe_sync_estado"`, s), nil
}

func (r *NFESerproRepository) GetSyncEstado(ctx context.Context, schemaName, provider, uf, cnpj string) (domain.NFESyncEstado, error) {
	tbl, err := tenantNFESyncStateTable(schemaName)
	if err != nil {
		return domain.NFESyncEstado{}, err
	}
	q := fmt.Sprintf(`
		SELECT id, provider, uf, cnpj, COALESCE(ultimo_nsu, '0'),
		       COALESCE(ultimo_cstat, 0), COALESCE(ultimo_motivo, ''),
		       ultima_verificacao, proxima_consulta_apos,
		       COALESCE(ultima_qt_dfe_ret, 0)
		FROM %s
		WHERE provider = $1 AND uf = $2 AND cnpj = $3
	`, tbl)
	var out domain.NFESyncEstado
	var ultimaVerif sql.NullTime
	var proxima sql.NullTime
	err = dbQueryRow(ctx, r.pool, q, strings.ToUpper(strings.TrimSpace(provider)), strings.ToUpper(strings.TrimSpace(uf)), onlyDigitsRepo(cnpj)).Scan(
		&out.ID,
		&out.Provider,
		&out.UF,
		&out.CNPJ,
		&out.UltimoNSU,
		&out.UltimoCStat,
		&out.UltimoMotivo,
		&ultimaVerif,
		&proxima,
		&out.UltimaQtDFeRet,
	)
	if err != nil {
		return domain.NFESyncEstado{}, err
	}
	if ultimaVerif.Valid {
		t := ultimaVerif.Time.UTC()
		out.UltimaVerificacao = &t
	}
	if proxima.Valid {
		t := proxima.Time.UTC()
		out.ProximaConsultaApos = &t
	}
	return out, nil
}

func (r *NFESerproRepository) UpsertSyncEstado(ctx context.Context, schemaName string, in NFESyncStateUpsert) (domain.NFESyncEstado, error) {
	tbl, err := tenantNFESyncStateTable(schemaName)
	if err != nil {
		return domain.NFESyncEstado{}, err
	}
	q := fmt.Sprintf(`
		INSERT INTO %s (
		    provider, uf, cnpj, ultimo_nsu, ultimo_cstat, ultimo_motivo, ultima_verificacao, proxima_consulta_apos, ultima_qt_dfe_ret
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (provider, uf, cnpj) DO UPDATE SET
		    ultimo_nsu = EXCLUDED.ultimo_nsu,
		    ultimo_cstat = EXCLUDED.ultimo_cstat,
		    ultimo_motivo = EXCLUDED.ultimo_motivo,
		    ultima_verificacao = EXCLUDED.ultima_verificacao,
		    proxima_consulta_apos = EXCLUDED.proxima_consulta_apos,
		    ultima_qt_dfe_ret = EXCLUDED.ultima_qt_dfe_ret,
		    updatedat = CURRENT_TIMESTAMP
		RETURNING id, provider, uf, cnpj, COALESCE(ultimo_nsu, '0'),
		          COALESCE(ultimo_cstat, 0), COALESCE(ultimo_motivo, ''),
		          ultima_verificacao, proxima_consulta_apos,
		          COALESCE(ultima_qt_dfe_ret, 0)
	`, tbl)
	var out domain.NFESyncEstado
	var ultimaVerif sql.NullTime
	var proxima sql.NullTime
	err = dbQueryRow(
		ctx,
		r.pool,
		q,
		strings.ToUpper(strings.TrimSpace(in.Provider)),
		strings.ToUpper(strings.TrimSpace(in.UF)),
		onlyDigitsRepo(in.CNPJ),
		normalizeNSUState(in.UltimoNSU),
		in.UltimoCStat,
		strings.TrimSpace(in.UltimoMotivo),
		in.UltimaVerificacao.UTC(),
		in.ProximaConsultaApos,
		in.UltimaQtDFeRet,
	).Scan(
		&out.ID,
		&out.Provider,
		&out.UF,
		&out.CNPJ,
		&out.UltimoNSU,
		&out.UltimoCStat,
		&out.UltimoMotivo,
		&ultimaVerif,
		&proxima,
		&out.UltimaQtDFeRet,
	)
	if err != nil {
		return domain.NFESyncEstado{}, fmt.Errorf("upsert nfe_sync_estado: %w", err)
	}
	if ultimaVerif.Valid {
		t := ultimaVerif.Time.UTC()
		out.UltimaVerificacao = &t
	}
	if proxima.Valid {
		t := proxima.Time.UTC()
		out.ProximaConsultaApos = &t
	}
	return out, nil
}

func (r *NFESerproRepository) ListSyncEstados(ctx context.Context, schemaName string, p NFESyncEstadoListParams) ([]domain.NFESyncEstado, int64, error) {
	tbl, err := tenantNFESyncStateTable(schemaName)
	if err != nil {
		return nil, 0, err
	}
	where := []string{"TRUE"}
	args := []any{}
	ai := 1
	if strings.TrimSpace(p.Provider) != "" {
		where = append(where, fmt.Sprintf("provider = $%d", ai))
		args = append(args, strings.ToUpper(strings.TrimSpace(p.Provider)))
		ai++
	}
	if strings.TrimSpace(p.UF) != "" {
		where = append(where, fmt.Sprintf("uf = $%d", ai))
		args = append(args, strings.ToUpper(strings.TrimSpace(p.UF)))
		ai++
	}
	if d := onlyDigitsRepo(p.CNPJ); d != "" {
		where = append(where, fmt.Sprintf("cnpj LIKE $%d", ai))
		args = append(args, "%"+d+"%")
		ai++
	}
	whereSQL := strings.Join(where, " AND ")
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, tbl, whereSQL)
	var total int64
	if err := dbQueryRow(ctx, r.pool, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count nfe_sync_estado: %w", err)
	}
	rows := p.Rows
	if rows <= 0 {
		rows = 50
	}
	if rows > 500 {
		rows = 500
	}
	first := p.First
	if first < 0 {
		first = 0
	}
	listQ := fmt.Sprintf(`
		SELECT id, provider, uf, cnpj, COALESCE(ultimo_nsu, '0'),
		       COALESCE(ultimo_cstat, 0), COALESCE(ultimo_motivo, ''),
		       ultima_verificacao, proxima_consulta_apos,
		       COALESCE(ultima_qt_dfe_ret, 0)
		FROM %s
		WHERE %s
		ORDER BY updatedat DESC
		LIMIT $%d OFFSET $%d
	`, tbl, whereSQL, ai, ai+1)
	args = append(args, rows, first)
	gr, err := dbQuery(ctx, r.pool, listQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list nfe_sync_estado: %w", err)
	}
	defer gr.Close()
	out := make([]domain.NFESyncEstado, 0, rows)
	for gr.Next() {
		var row domain.NFESyncEstado
		var ultima sql.NullTime
		var proxima sql.NullTime
		if err := gr.Scan(
			&row.ID,
			&row.Provider,
			&row.UF,
			&row.CNPJ,
			&row.UltimoNSU,
			&row.UltimoCStat,
			&row.UltimoMotivo,
			&ultima,
			&proxima,
			&row.UltimaQtDFeRet,
		); err != nil {
			return nil, 0, fmt.Errorf("scan nfe_sync_estado: %w", err)
		}
		if ultima.Valid {
			t := ultima.Time.UTC()
			row.UltimaVerificacao = &t
		}
		if proxima.Valid {
			t := proxima.Time.UTC()
			row.ProximaConsultaApos = &t
		}
		out = append(out, row)
	}
	return out, total, nil
}

// ListSyncEstadosDue retorna checkpoints elegíveis para nova consulta ao provider:
// respeita proxima_consulta_apos quando preenchida; quando nula, exige ultima_verificacao
// anterior a (now - minGap) ou ultima nula.
func (r *NFESerproRepository) ListSyncEstadosDue(ctx context.Context, schemaName string, now time.Time, minGap time.Duration, limit int) ([]domain.NFESyncEstado, error) {
	tbl, err := tenantNFESyncStateTable(schemaName)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	cutoff := now.Add(-minGap)
	q := fmt.Sprintf(`
		SELECT id, provider, uf, cnpj, COALESCE(ultimo_nsu, '0'),
		       COALESCE(ultimo_cstat, 0), COALESCE(ultimo_motivo, ''),
		       ultima_verificacao, proxima_consulta_apos,
		       COALESCE(ultima_qt_dfe_ret, 0)
		FROM %s
		WHERE UPPER(TRIM(provider)) NOT IN ('MOCK', 'NACIONAL')
		  AND (
		      (proxima_consulta_apos IS NOT NULL AND proxima_consulta_apos <= $1)
		      OR (proxima_consulta_apos IS NULL AND (ultima_verificacao IS NULL OR ultima_verificacao <= $2))
		  )
		ORDER BY COALESCE(proxima_consulta_apos, ultima_verificacao) ASC NULLS FIRST
		LIMIT $3
	`, tbl)
	gr, err := dbQuery(ctx, r.pool, q, now.UTC(), cutoff.UTC(), limit)
	if err != nil {
		return nil, fmt.Errorf("list nfe_sync_estado due: %w", err)
	}
	defer gr.Close()
	out := make([]domain.NFESyncEstado, 0, limit)
	for gr.Next() {
		var row domain.NFESyncEstado
		var ultima sql.NullTime
		var proxima sql.NullTime
		if err := gr.Scan(
			&row.ID,
			&row.Provider,
			&row.UF,
			&row.CNPJ,
			&row.UltimoNSU,
			&row.UltimoCStat,
			&row.UltimoMotivo,
			&ultima,
			&proxima,
			&row.UltimaQtDFeRet,
		); err != nil {
			return nil, fmt.Errorf("scan nfe_sync_estado due: %w", err)
		}
		if ultima.Valid {
			t := ultima.Time.UTC()
			row.UltimaVerificacao = &t
		}
		if proxima.Valid {
			t := proxima.Time.UTC()
			row.ProximaConsultaApos = &t
		}
		out = append(out, row)
	}
	return out, gr.Err()
}

func normalizeNSUState(v string) string {
	v = onlyDigitsRepo(v)
	if v == "" {
		return "0"
	}
	return strings.TrimLeft(v, "0")
}
