package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/domain"
)

func tenantNFEManifestacaoTable(schemaName string) (string, error) {
	s, err := normalizeSchemaForNFE(schemaName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`"%s"."nfe_manifestacao_dest"`, s), nil
}

type NFEManifestacaoInsert struct {
	ChaveNFe      string
	TpEvento      string
	CNPJDest      string
	CStatLote     int
	XMotivoLote   string
	CStatEvento   int
	XMotivoEvento string
	NProt         string
	RetornoXML    string
}

func (r *NFESerproRepository) InsertManifestacaoDest(ctx context.Context, schemaName string, in NFEManifestacaoInsert) (domain.NFEManifestacaoDest, error) {
	tbl, err := tenantNFEManifestacaoTable(schemaName)
	if err != nil {
		return domain.NFEManifestacaoDest{}, err
	}
	q := fmt.Sprintf(`
		INSERT INTO %s (
		    chave_nfe, tp_evento, cnpj_dest, cstat_lote, x_motivo_lote,
		    cstat_evento, x_motivo_evento, n_prot, retorno_xml
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, chave_nfe, tp_evento, COALESCE(cnpj_dest, ''),
		          COALESCE(cstat_lote, 0), COALESCE(x_motivo_lote, ''),
		          COALESCE(cstat_evento, 0), COALESCE(x_motivo_evento, ''),
		          COALESCE(n_prot, ''), COALESCE(retorno_xml, ''), criado_em
	`, tbl)
	var out domain.NFEManifestacaoDest
	err = dbQueryRow(
		ctx,
		r.pool,
		q,
		strings.TrimSpace(in.ChaveNFe),
		strings.TrimSpace(in.TpEvento),
		onlyDigitsRepo(in.CNPJDest),
		in.CStatLote,
		strings.TrimSpace(in.XMotivoLote),
		in.CStatEvento,
		strings.TrimSpace(in.XMotivoEvento),
		strings.TrimSpace(in.NProt),
		in.RetornoXML,
	).Scan(
		&out.ID,
		&out.ChaveNFe,
		&out.TpEvento,
		&out.CNPJDest,
		&out.CStatLote,
		&out.XMotivoLote,
		&out.CStatEvento,
		&out.XMotivoEvento,
		&out.NProt,
		&out.RetornoXML,
		&out.CriadoEm,
	)
	if err != nil {
		return domain.NFEManifestacaoDest{}, fmt.Errorf("insert nfe_manifestacao_dest: %w", err)
	}
	return out, nil
}

func (r *NFESerproRepository) ListManifestacaoByChave(ctx context.Context, schemaName, chaveNFe string, limit int) ([]domain.NFEManifestacaoDest, int64, error) {
	tbl, err := tenantNFEManifestacaoTable(schemaName)
	if err != nil {
		return nil, 0, err
	}
	ch := strings.TrimSpace(chaveNFe)
	if len(ch) != 44 {
		return nil, 0, fmt.Errorf("chave_nfe invalida")
	}
	if limit <= 0 {
		limit = 30
	}
	if limit > 200 {
		limit = 200
	}
	const countQ = `SELECT COUNT(*) FROM %s WHERE chave_nfe = $1`
	var total int64
	if err := dbQueryRow(ctx, r.pool, fmt.Sprintf(countQ, tbl), ch).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count nfe_manifestacao_dest: %w", err)
	}
	listQ := fmt.Sprintf(`
		SELECT id, chave_nfe, tp_evento, COALESCE(cnpj_dest, ''),
		       COALESCE(cstat_lote, 0), COALESCE(x_motivo_lote, ''),
		       COALESCE(cstat_evento, 0), COALESCE(x_motivo_evento, ''),
		       COALESCE(n_prot, ''), COALESCE(retorno_xml, ''), criado_em
		FROM %s
		WHERE chave_nfe = $1
		ORDER BY criado_em DESC
		LIMIT $2
	`, tbl)
	gr, err := dbQuery(ctx, r.pool, listQ, ch, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("list nfe_manifestacao_dest: %w", err)
	}
	defer gr.Close()
	out := make([]domain.NFEManifestacaoDest, 0, limit)
	for gr.Next() {
		var row domain.NFEManifestacaoDest
		if err := gr.Scan(
			&row.ID,
			&row.ChaveNFe,
			&row.TpEvento,
			&row.CNPJDest,
			&row.CStatLote,
			&row.XMotivoLote,
			&row.CStatEvento,
			&row.XMotivoEvento,
			&row.NProt,
			&row.RetornoXML,
			&row.CriadoEm,
		); err != nil {
			return nil, 0, fmt.Errorf("scan nfe_manifestacao_dest: %w", err)
		}
		out = append(out, row)
	}
	return out, total, gr.Err()
}
