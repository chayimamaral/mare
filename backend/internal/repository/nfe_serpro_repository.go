package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NFESerproRepository struct {
	pool *pgxpool.Pool
}

func NewNFESerproRepository(pool *pgxpool.Pool) *NFESerproRepository {
	return &NFESerproRepository{pool: pool}
}

var schemaNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]{2,62}$`)

func normalizeSchemaForNFE(schemaName string) (string, error) {
	s := strings.TrimSpace(strings.ToLower(schemaName))
	if !schemaNameRegex.MatchString(s) {
		return "", fmt.Errorf("schema invalido")
	}
	return s, nil
}

func tenantNFETable(schemaName string) (string, error) {
	s, err := normalizeSchemaForNFE(schemaName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`"%s"."nfe_documento"`, s), nil
}

func (r *NFESerproRepository) UpsertDocumento(ctx context.Context, schemaName string, doc domain.NFEDocumento) (domain.NFEDocumento, error) {
	tableName, err := tenantNFETable(schemaName)
	if err != nil {
		return domain.NFEDocumento{}, err
	}

	q := fmt.Sprintf(`
		INSERT INTO %s
		    (chave_nfe, ambiente, evento_codigo, evento_descricao, origem, payload_json, payload_xml, content_type_origem, request_tag, status_http)
		VALUES
		    ($1, $2, $3, $4, $5, $6::jsonb, $7, $8, $9, $10)
		ON CONFLICT (chave_nfe) DO UPDATE
		SET ambiente = EXCLUDED.ambiente,
		    evento_codigo = EXCLUDED.evento_codigo,
		    evento_descricao = EXCLUDED.evento_descricao,
		    origem = EXCLUDED.origem,
		    payload_json = EXCLUDED.payload_json,
		    payload_xml = EXCLUDED.payload_xml,
		    content_type_origem = EXCLUDED.content_type_origem,
		    request_tag = EXCLUDED.request_tag,
		    status_http = EXCLUDED.status_http,
		    recebido_em = now(),
		    updatedat = CURRENT_TIMESTAMP
		RETURNING id, chave_nfe, COALESCE(ambiente, ''), COALESCE(evento_codigo, ''), COALESCE(evento_descricao, ''),
		          origem, payload_json, COALESCE(payload_xml, ''), COALESCE(content_type_origem, ''), COALESCE(request_tag, ''),
		          COALESCE(status_http, 0), recebido_em
	`, tableName)

	var out domain.NFEDocumento
	var payload []byte
	err = dbQueryRow(
		ctx,
		r.pool,
		q,
		doc.ChaveNFe,
		doc.Ambiente,
		doc.EventoCodigo,
		doc.EventoDescricao,
		doc.Origem,
		doc.PayloadJSON,
		doc.PayloadXML,
		doc.ContentTypeOrigem,
		doc.RequestTag,
		doc.StatusHTTP,
	).Scan(
		&out.ID,
		&out.ChaveNFe,
		&out.Ambiente,
		&out.EventoCodigo,
		&out.EventoDescricao,
		&out.Origem,
		&payload,
		&out.PayloadXML,
		&out.ContentTypeOrigem,
		&out.RequestTag,
		&out.StatusHTTP,
		&out.RecebidoEm,
	)
	if err != nil {
		return domain.NFEDocumento{}, fmt.Errorf("upsert nfe_documento: %w", err)
	}
	if len(payload) == 0 {
		out.PayloadJSON = json.RawMessage(`{}`)
	} else {
		out.PayloadJSON = json.RawMessage(payload)
	}
	return out, nil
}

func (r *NFESerproRepository) GetDocumentoByChave(ctx context.Context, schemaName, chave string) (domain.NFEDocumento, error) {
	tableName, err := tenantNFETable(schemaName)
	if err != nil {
		return domain.NFEDocumento{}, err
	}
	q := fmt.Sprintf(`
		SELECT id, chave_nfe, COALESCE(ambiente, ''), COALESCE(evento_codigo, ''), COALESCE(evento_descricao, ''),
		       origem, payload_json, COALESCE(payload_xml, ''), COALESCE(content_type_origem, ''), COALESCE(request_tag, ''),
		       COALESCE(status_http, 0), recebido_em
		FROM %s
		WHERE chave_nfe = $1
	`, tableName)
	var out domain.NFEDocumento
	var payload []byte
	err = dbQueryRow(ctx, r.pool, q, chave).Scan(
		&out.ID,
		&out.ChaveNFe,
		&out.Ambiente,
		&out.EventoCodigo,
		&out.EventoDescricao,
		&out.Origem,
		&payload,
		&out.PayloadXML,
		&out.ContentTypeOrigem,
		&out.RequestTag,
		&out.StatusHTTP,
		&out.RecebidoEm,
	)
	if err != nil {
		return domain.NFEDocumento{}, fmt.Errorf("get nfe_documento por chave: %w", err)
	}
	if len(payload) == 0 {
		out.PayloadJSON = json.RawMessage(`{}`)
	} else {
		out.PayloadJSON = json.RawMessage(payload)
	}
	return out, nil
}

func (r *NFESerproRepository) SavePushNotificacao(ctx context.Context, notif domain.NFEPushNotificacao) error {
	const q = `
		INSERT INTO public.nfe_serpro_push_notificacao
		    (chave_nfe, data_hora_envio, payload, headers)
		VALUES
		    ($1, $2, $3::jsonb, $4::jsonb)
	`
	_, err := dbExec(ctx, r.pool, q, notif.ChaveNFe, notif.DataHoraEnvio, notif.Payload, notif.Headers)
	if err != nil {
		return fmt.Errorf("save nfe_serpro_push_notificacao: %w", err)
	}
	return nil
}
