package service

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
	"github.com/jackc/pgx/v5"
)

const (
	consultaNFEBaseTrial    = "https://gateway.apiserpro.serpro.gov.br/consulta-nfe-df-trial/api"
	consultaNFEBaseProducao = "https://gateway.apiserpro.serpro.gov.br/consulta-nfe-df/api"
)

type NFESerproService struct {
	repo       *repository.NFESerproRepository
	serproAuth *SerproService
}

func NewNFESerproService(repo *repository.NFESerproRepository, serproAuth *SerproService) *NFESerproService {
	return &NFESerproService{repo: repo, serproAuth: serproAuth}
}

var nfeChaveRegex = regexp.MustCompile(`^\d{44}$`)

func normalizeStaticBearer(raw string) string {
	t := strings.TrimSpace(raw)
	if t == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(t), "bearer ") {
		t = strings.TrimSpace(t[7:])
	}
	return t
}

func validateNFEChave(chave string) error {
	c := strings.TrimSpace(chave)
	if !nfeChaveRegex.MatchString(c) {
		return fmt.Errorf("chave_nfe deve ter 44 digitos numericos")
	}
	return nil
}

func normalizeNFEAmbienteIntegra(raw string) string {
	a := strings.ToLower(strings.TrimSpace(raw))
	if a == string(IntegraAmbienteProducao) {
		return string(IntegraAmbienteProducao)
	}
	return string(IntegraAmbienteTrial)
}

// resolveConsultaNFEAPIBase: se SERPRO_NFE_API_BASE_URL estiver definida, ela prevalece (ambiente no body e ignorado).
// Caso contrario, usa o gateway SERPRO conforme trial ou producao (mesmo padrao do Integra Contador).
func (s *NFESerproService) resolveConsultaNFEAPIBase(ambiente string) string {
	if s.serproAuth == nil {
		return ""
	}
	if b := strings.TrimSpace(s.serproAuth.cfg.SerproNFEAPIBaseURL); b != "" {
		return strings.TrimSuffix(b, "/")
	}
	switch normalizeNFEAmbienteIntegra(ambiente) {
	case string(IntegraAmbienteProducao):
		return strings.TrimSuffix(consultaNFEBaseProducao, "/")
	default:
		return strings.TrimSuffix(consultaNFEBaseTrial, "/")
	}
}

func (s *NFESerproService) ConsultarNFe(ctx context.Context, schemaName, ambiente, chaveNFe, requestTag string, assinar bool) (domain.NFEDocumento, error) {
	if s.serproAuth == nil {
		return domain.NFEDocumento{}, fmt.Errorf("servico SERPRO nao configurado")
	}
	if err := validateNFEChave(chaveNFe); err != nil {
		return domain.NFEDocumento{}, err
	}

	chave := strings.TrimSpace(chaveNFe)
	if cached, err := s.repo.GetDocumentoByChave(ctx, schemaName, chave); err == nil {
		cached.JaBaixada = true
		return cached, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return domain.NFEDocumento{}, err
	}

	ambNorm := normalizeNFEAmbienteIntegra(ambiente)
	base := s.resolveConsultaNFEAPIBase(ambNorm)
	if base == "" {
		return domain.NFEDocumento{}, fmt.Errorf("base URL consulta NF-e nao configurada")
	}

	var token string
	if static := normalizeStaticBearer(s.serproAuth.cfg.SerproNFEBearerToken); static != "" {
		token = static
	} else {
		var err error
		token, err = s.serproAuth.ObterBearerToken(ctx)
		if err != nil {
			return domain.NFEDocumento{}, err
		}
	}

	u := base + "/v1/nfe/" + strings.TrimSpace(chaveNFe)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return domain.NFEDocumento{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if strings.TrimSpace(requestTag) != "" {
		req.Header.Set("x-request-tag", strings.TrimSpace(requestTag))
	}
	if assinar {
		req.Header.Set("x-signature", "1")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.NFEDocumento{}, fmt.Errorf("consulta nfe serpro: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.NFEDocumento{}, fmt.Errorf("consulta nfe serpro status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	jsonPayload, err := normalizeJSON(body)
	if err != nil {
		return domain.NFEDocumento{}, fmt.Errorf("resposta json invalida: %w", err)
	}
	payloadXML, err := jsonToXML(jsonPayload)
	if err != nil {
		return domain.NFEDocumento{}, fmt.Errorf("erro ao converter json para xml: %w", err)
	}

	doc := domain.NFEDocumento{
		ChaveNFe:          strings.TrimSpace(chaveNFe),
		Ambiente:          ambNorm,
		Origem:            "CONSULTA_NFE_SERPRO",
		PayloadJSON:       jsonPayload,
		PayloadXML:        payloadXML,
		ContentTypeOrigem: resp.Header.Get("Content-Type"),
		RequestTag:        strings.TrimSpace(requestTag),
		StatusHTTP:        resp.StatusCode,
	}
	fillEventoMetaFromJSON(&doc)
	out, err := s.repo.UpsertDocumento(ctx, schemaName, doc)
	if err != nil {
		return domain.NFEDocumento{}, err
	}
	return out, nil
}

func (s *NFESerproService) BuscarDocumento(ctx context.Context, schemaName, chaveNFe string) (domain.NFEDocumento, error) {
	if err := validateNFEChave(chaveNFe); err != nil {
		return domain.NFEDocumento{}, err
	}
	return s.repo.GetDocumentoByChave(ctx, schemaName, strings.TrimSpace(chaveNFe))
}

func (s *NFESerproService) ExportarXML(ctx context.Context, schemaName, chaveNFe string) (string, error) {
	doc, err := s.BuscarDocumento(ctx, schemaName, chaveNFe)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(doc.PayloadXML) != "" {
		return doc.PayloadXML, nil
	}
	if len(doc.PayloadJSON) == 0 {
		return "<nfe/>", nil
	}
	return jsonToXML(doc.PayloadJSON)
}

func (s *NFESerproService) RegistrarPushNotificacao(ctx context.Context, rawBody []byte, headers map[string]string) error {
	var payload map[string]any
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		payload = map[string]any{"raw_body": string(rawBody)}
	}
	bPayload, _ := json.Marshal(payload)
	bHeaders, _ := json.Marshal(headers)
	notif := domain.NFEPushNotificacao{
		Payload: json.RawMessage(bPayload),
		Headers: json.RawMessage(bHeaders),
	}
	if v, ok := payload["chaveNFe"]; ok {
		notif.ChaveNFe = strings.TrimSpace(fmt.Sprintf("%v", v))
	}
	if v, ok := payload["dataHoraEnvio"]; ok {
		ts := strings.TrimSpace(fmt.Sprintf("%v", v))
		if ts != "" {
			if parsed, err := parseTimeFlexible(ts); err == nil {
				notif.DataHoraEnvio = &parsed
			}
		}
	}
	return s.repo.SavePushNotificacao(ctx, notif)
}

func normalizeJSON(raw []byte) (json.RawMessage, error) {
	var anyPayload any
	if err := json.Unmarshal(raw, &anyPayload); err != nil {
		return nil, err
	}
	b, err := json.Marshal(anyPayload)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

func fillEventoMetaFromJSON(doc *domain.NFEDocumento) {
	if doc == nil || len(doc.PayloadJSON) == 0 {
		return
	}
	var obj map[string]any
	if err := json.Unmarshal(doc.PayloadJSON, &obj); err != nil {
		return
	}
	if nfeProc, ok := obj["nfeProc"].(map[string]any); ok {
		if protNFe, ok := nfeProc["protNFe"].(map[string]any); ok {
			if infProt, ok := protNFe["infProt"].(map[string]any); ok {
				if cStat, ok := infProt["cStat"]; ok {
					doc.EventoCodigo = strings.TrimSpace(fmt.Sprintf("%v", cStat))
				}
				if motivo, ok := infProt["xMotivo"]; ok {
					doc.EventoDescricao = strings.TrimSpace(fmt.Sprintf("%v", motivo))
				}
			}
		}
	}
}

func parseTimeFlexible(v string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	trimmed := strings.TrimSpace(v)
	for _, layout := range layouts {
		if t, err := time.Parse(layout, trimmed); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato de data/hora invalido")
}

func jsonToXML(raw json.RawMessage) (string, error) {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", err
	}
	var b bytes.Buffer
	b.WriteString(xml.Header)
	b.WriteString("<nfe>")
	if err := writeXMLValue(&b, "documento", v); err != nil {
		return "", err
	}
	b.WriteString("</nfe>")
	return b.String(), nil
}

func writeXMLValue(buf *bytes.Buffer, key string, v any) error {
	tag := sanitizeXMLTag(key)
	switch val := v.(type) {
	case map[string]any:
		buf.WriteString("<" + tag + ">")
		for k, child := range val {
			if err := writeXMLValue(buf, k, child); err != nil {
				return err
			}
		}
		buf.WriteString("</" + tag + ">")
	case []any:
		for i, item := range val {
			itemTag := tag + "_item_" + strconv.Itoa(i)
			if err := writeXMLValue(buf, itemTag, item); err != nil {
				return err
			}
		}
	default:
		buf.WriteString("<" + tag + ">")
		if err := xml.EscapeText(buf, []byte(fmt.Sprintf("%v", val))); err != nil {
			return err
		}
		buf.WriteString("</" + tag + ">")
	}
	return nil
}

func sanitizeXMLTag(tag string) string {
	t := strings.TrimSpace(strings.ToLower(tag))
	if t == "" {
		return "campo"
	}
	var out strings.Builder
	for _, r := range t {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			out.WriteRune(r)
			continue
		}
		out.WriteRune('_')
	}
	s := out.String()
	if s == "" {
		return "campo"
	}
	if s[0] >= '0' && s[0] <= '9' {
		return "f_" + s
	}
	return s
}
