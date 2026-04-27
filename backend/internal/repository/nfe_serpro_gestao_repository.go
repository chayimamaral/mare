package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
)

type NFEGestaoListParams struct {
	First            int
	Rows             int
	SortField        string
	SortOrder        int
	TipoArquivo      string
	ChaveNFe         string
	EmissaoIni       *time.Time
	EmissaoFim       *time.Time
	CNPJEmitente     string
	CNPJDestinatario string
}

func tenantNFEGestaoTable(schemaName string) (string, error) {
	s, err := normalizeSchemaForNFE(schemaName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`"%s"."nfe_gestao"`, s), nil
}

func (r *NFESerproRepository) UpsertGestao(ctx context.Context, schemaName string, g domain.NFEGestao) (domain.NFEGestao, error) {
	tbl, err := tenantNFEGestaoTable(schemaName)
	if err != nil {
		return domain.NFEGestao{}, err
	}

	var emiss any
	if g.DataEmissao != nil {
		emiss = g.DataEmissao.UTC().Format("2006-01-02")
	} else {
		emiss = nil
	}
	var valor any
	if g.ValorTotal != nil {
		valor = *g.ValorTotal
	} else {
		valor = nil
	}

	q := fmt.Sprintf(`
		INSERT INTO %s (
		    chave_nfe, tipo_arquivo, numero_nfe, razao_social_emitente, cnpj_emitente,
		    data_emissao, cnpj_destinatario, valor_total, data_download
		)
		VALUES ($1, $2, $3, $4, $5, $6::date, $7, $8::numeric, $9)
		ON CONFLICT (chave_nfe) DO UPDATE SET
		    tipo_arquivo = EXCLUDED.tipo_arquivo,
		    numero_nfe = EXCLUDED.numero_nfe,
		    razao_social_emitente = EXCLUDED.razao_social_emitente,
		    cnpj_emitente = EXCLUDED.cnpj_emitente,
		    data_emissao = EXCLUDED.data_emissao,
		    cnpj_destinatario = EXCLUDED.cnpj_destinatario,
		    valor_total = EXCLUDED.valor_total,
		    data_download = EXCLUDED.data_download,
		    updatedat = CURRENT_TIMESTAMP
		RETURNING id, chave_nfe, tipo_arquivo,
		          COALESCE(numero_nfe, ''), COALESCE(razao_social_emitente, ''),
		          COALESCE(cnpj_emitente, ''), data_emissao, COALESCE(cnpj_destinatario, ''),
		          valor_total, data_download
	`, tbl)

	var out domain.NFEGestao
	var dataEm sql.NullTime
	var val sql.NullFloat64
	err = dbQueryRow(
		ctx,
		r.pool,
		q,
		g.ChaveNFe,
		g.TipoArquivo,
		nullStrPtr(g.NumeroNFe),
		nullStrPtr(g.RazaoSocialEmitente),
		nullStrPtr(g.CNPJEmitente),
		emiss,
		nullStrPtr(g.CNPJDestinatario),
		valor,
		g.DataDownload,
	).Scan(
		&out.ID,
		&out.ChaveNFe,
		&out.TipoArquivo,
		&out.NumeroNFe,
		&out.RazaoSocialEmitente,
		&out.CNPJEmitente,
		&dataEm,
		&out.CNPJDestinatario,
		&val,
		&out.DataDownload,
	)
	if err != nil {
		return domain.NFEGestao{}, fmt.Errorf("upsert nfe_gestao: %w", err)
	}
	if dataEm.Valid {
		t := dataEm.Time.UTC()
		out.DataEmissao = &t
	}
	if val.Valid {
		x := val.Float64
		out.ValorTotal = &x
	}
	return out, nil
}

func nullStrPtr(s string) any {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return s
}

func (r *NFESerproRepository) ListGestao(ctx context.Context, schemaName string, p NFEGestaoListParams) ([]domain.NFEGestao, int64, error) {
	tbl, err := tenantNFEGestaoTable(schemaName)
	if err != nil {
		return nil, 0, err
	}

	where := []string{"TRUE"}
	args := []any{}
	ai := 1

	if t := strings.TrimSpace(p.TipoArquivo); t != "" {
		where = append(where, fmt.Sprintf("tipo_arquivo = $%d", ai))
		args = append(args, t)
		ai++
	}
	if d := onlyDigitsRepo(p.ChaveNFe); d != "" {
		where = append(where, fmt.Sprintf("regexp_replace(COALESCE(chave_nfe, ''), '[^0-9]', '', 'g') LIKE $%d", ai))
		args = append(args, "%"+d+"%")
		ai++
	}
	if p.EmissaoIni != nil {
		where = append(where, fmt.Sprintf("data_emissao >= $%d::date", ai))
		args = append(args, p.EmissaoIni.UTC().Format("2006-01-02"))
		ai++
	}
	if p.EmissaoFim != nil {
		where = append(where, fmt.Sprintf("data_emissao <= $%d::date", ai))
		args = append(args, p.EmissaoFim.UTC().Format("2006-01-02"))
		ai++
	}
	if d := onlyDigitsRepo(p.CNPJEmitente); d != "" {
		where = append(where, fmt.Sprintf("regexp_replace(COALESCE(cnpj_emitente, ''), '[^0-9]', '', 'g') LIKE $%d", ai))
		args = append(args, "%"+d+"%")
		ai++
	}
	if d := onlyDigitsRepo(p.CNPJDestinatario); d != "" {
		where = append(where, fmt.Sprintf("regexp_replace(COALESCE(cnpj_destinatario, ''), '[^0-9]', '', 'g') LIKE $%d", ai))
		args = append(args, "%"+d+"%")
		ai++
	}

	whereSQL := strings.Join(where, " AND ")

	orderBy := "data_download DESC"
	switch p.SortField {
	case "chave_nfe":
		orderBy = orderCol("chave_nfe", p.SortOrder)
	case "tipo_arquivo":
		orderBy = orderCol("tipo_arquivo", p.SortOrder)
	case "numero_nfe":
		orderBy = orderCol("numero_nfe", p.SortOrder)
	case "razao_social_emitente":
		orderBy = orderCol("razao_social_emitente", p.SortOrder)
	case "cnpj_emitente":
		orderBy = orderCol("cnpj_emitente", p.SortOrder)
	case "data_emissao":
		orderBy = orderCol("data_emissao", p.SortOrder)
	case "cnpj_destinatario":
		orderBy = orderCol("cnpj_destinatario", p.SortOrder)
	case "valor_total":
		orderBy = orderCol("valor_total", p.SortOrder)
	case "data_download":
		orderBy = orderCol("data_download", p.SortOrder)
	}

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, tbl, whereSQL)
	var total int64
	if err := dbQueryRow(ctx, r.pool, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count nfe_gestao: %w", err)
	}

	rows := p.Rows
	if rows <= 0 {
		rows = 25
	}
	if rows > 3000 {
		rows = 3000
	}
	first := p.First
	if first < 0 {
		first = 0
	}

	listQ := fmt.Sprintf(`
		SELECT id, chave_nfe, tipo_arquivo,
		       COALESCE(numero_nfe, ''), COALESCE(razao_social_emitente, ''),
		       COALESCE(cnpj_emitente, ''), data_emissao, COALESCE(cnpj_destinatario, ''),
		       valor_total, data_download
		FROM %s
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, tbl, whereSQL, orderBy, ai, ai+1)
	args = append(args, rows, first)

	gr, err := dbQuery(ctx, r.pool, listQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list nfe_gestao: %w", err)
	}
	defer gr.Close()

	var out []domain.NFEGestao
	for gr.Next() {
		var row domain.NFEGestao
		var dataEm sql.NullTime
		var val sql.NullFloat64
		if err := gr.Scan(
			&row.ID,
			&row.ChaveNFe,
			&row.TipoArquivo,
			&row.NumeroNFe,
			&row.RazaoSocialEmitente,
			&row.CNPJEmitente,
			&dataEm,
			&row.CNPJDestinatario,
			&val,
			&row.DataDownload,
		); err != nil {
			return nil, 0, fmt.Errorf("scan nfe_gestao: %w", err)
		}
		if dataEm.Valid {
			t := dataEm.Time.UTC()
			row.DataEmissao = &t
		}
		if val.Valid {
			x := val.Float64
			row.ValorTotal = &x
		}
		out = append(out, row)
	}
	return out, total, nil
}

func onlyDigitsRepo(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func orderCol(col string, sortOrder int) string {
	dir := "ASC"
	if sortOrder < 0 {
		dir = "DESC"
	}
	return col + " " + dir + " NULLS LAST"
}
