package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/repository"
	"github.com/chayimamaral/vecx/backend/internal/service"
)

type NFESerproHandler struct {
	svc            *service.NFESerproService
	validacaoCat   *service.NFEValidacaoCatalogoService
	certificadoSvc *service.CertificadoService
}

func NewNFESerproHandler(svc *service.NFESerproService, validacaoCat *service.NFEValidacaoCatalogoService, certificadoSvc *service.CertificadoService) *NFESerproHandler {
	return &NFESerproHandler{svc: svc, validacaoCat: validacaoCat, certificadoSvc: certificadoSvc}
}

type nfeConsultaRequest struct {
	Ambiente   string `json:"ambiente"`
	ChaveNFe   string `json:"chave_nfe"`
	RequestTag string `json:"request_tag"`
	Assinar    bool   `json:"assinar"`
}

type nfeSyncProviderRequest struct {
	Provider string `json:"provider"`
	UF       string `json:"uf"`
	CNPJ     string `json:"cnpj"`
	Ambiente string `json:"ambiente"`
	Simular  bool   `json:"simular"`
}

type nfeManifestarDestRequest struct {
	ChaveNFe   string `json:"chave_nfe"`
	TpEvento   string `json:"tp_evento"`
	CNPJDest   string `json:"cnpj_destinatario"`
	Ambiente   string `json:"ambiente"`
	XJust      string `json:"x_just"`
	NSeqEvento int    `json:"n_seq_evento"`
	Simular    bool   `json:"simular"`
}

func hasOnlyDigits(v string) bool {
	s := strings.TrimSpace(v)
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func (h *NFESerproHandler) Consultar(w http.ResponseWriter, r *http.Request) {
	var req nfeConsultaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeNFEError(w, r, err, "mensagem_inicial_webservice")
		return
	}
	ch := strings.TrimSpace(req.ChaveNFe)
	if len(ch) != 44 || !hasOnlyDigits(ch) {
		h.writeNFEError(w, r, errNFEDadosInvalidos("chave_nfe deve ter 44 digitos numericos"), "area_dados")
		return
	}
	schema := middleware.TenantSchema(r.Context())
	tenantID := middleware.TenantID(r.Context())
	// Certificado de transmissao so e obrigatorio quando a chamada pede assinatura explicita.
	// Em cenarios de consulta trial/exemplo (assinar=false), nao bloqueia.
	if h.certificadoSvc != nil && req.Assinar {
		if _, err := h.certificadoSvc.ResumoPorTenant(r.Context(), tenantID); err != nil {
			h.writeNFEError(w, r, err, "certificado_transmissao")
			return
		}
	}
	out, err := h.svc.ConsultarNFe(r.Context(), schema, req.Ambiente, req.ChaveNFe, req.RequestTag, req.Assinar)
	if err != nil {
		h.writeNFEError(w, r, err, "mensagem_inicial_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

// SincronizarProvider aciona o provider pattern (EF-920) com estado de NSU por tenant/schema.
func (h *NFESerproHandler) SincronizarProvider(w http.ResponseWriter, r *http.Request) {
	var req nfeSyncProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeNFEError(w, r, err, "mensagem_inicial_webservice")
		return
	}
	cnpj := strings.TrimSpace(req.CNPJ)
	if cnpj != "" && !hasOnlyDigits(cnpj) {
		h.writeNFEError(w, r, errNFEDadosInvalidos("cnpj/cpf deve conter apenas digitos"), "area_dados")
		return
	}
	schema := middleware.TenantSchema(r.Context())
	tenantID := middleware.TenantID(r.Context())
	out, err := h.svc.SincronizarPorProvider(
		r.Context(),
		schema,
		tenantID,
		req.Provider,
		req.UF,
		req.CNPJ,
		req.Ambiente,
		req.Simular,
	)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

// ManifestarDestinatario POST /serpro/nfe/manifestar-destinatario — Recepção de Evento SVRS (manifestação 2102xx).
func (h *NFESerproHandler) ManifestarDestinatario(w http.ResponseWriter, r *http.Request) {
	var req nfeManifestarDestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeNFEError(w, r, err, "mensagem_inicial_webservice")
		return
	}
	ch := strings.TrimSpace(req.ChaveNFe)
	if len(ch) != 44 || !hasOnlyDigits(ch) {
		h.writeNFEError(w, r, errNFEDadosInvalidos("chave_nfe deve ter 44 digitos numericos"), "area_dados")
		return
	}
	schema := middleware.TenantSchema(r.Context())
	tenantID := middleware.TenantID(r.Context())
	if strings.TrimSpace(req.Ambiente) == "" {
		req.Ambiente = "producao"
	}
	out, err := h.svc.ManifestarDestinatario(
		r.Context(),
		schema,
		tenantID,
		req.ChaveNFe,
		req.TpEvento,
		req.CNPJDest,
		req.Ambiente,
		req.XJust,
		req.NSeqEvento,
		req.Simular,
	)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func (h *NFESerproHandler) ListRegrasValidacao(w http.ResponseWriter, r *http.Request) {
	etapa := strings.TrimSpace(r.URL.Query().Get("etapa"))
	if h.validacaoCat == nil {
		render.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}, "totalRecords": 0})
		return
	}
	items, err := h.validacaoCat.ListRegrasAtivasPorEtapa(r.Context(), etapa)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{
		"items":        items,
		"totalRecords": len(items),
	})
}

func (h *NFESerproHandler) writeNFEError(w http.ResponseWriter, r *http.Request, err error, etapa string) {
	if h.validacaoCat == nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	d := h.validacaoCat.ResolverErro(r.Context(), "webservice", etapa, err)
	render.WriteJSON(w, http.StatusBadRequest, map[string]any{
		"error":           d.Mensagem,
		"codigo":          d.Codigo,
		"mensagem":        d.Mensagem,
		"etapa_validacao": d.EtapaValidacao,
		"acao_sugerida":   d.AcaoSugerida,
		"origem":          d.Origem,
	})
}

type errNFEDadosInvalidos string

func (e errNFEDadosInvalidos) Error() string {
	return string(e)
}

// ListManifestacaoDest GET /serpro/nfe/manifestacao — histórico por chave.
func (h *NFESerproHandler) ListManifestacaoDest(w http.ResponseWriter, r *http.Request) {
	chave := strings.TrimSpace(r.URL.Query().Get("chave"))
	schema := middleware.TenantSchema(r.Context())
	out, err := h.svc.ListManifestacaoDest(r.Context(), schema, chave)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

// ListSyncEstado GET /serpro/nfe/sync-estado - checkpoint de sincronizacao por provider/UF/CNPJ.
func (h *NFESerproHandler) ListSyncEstado(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	schema := middleware.TenantSchema(r.Context())
	p := repository.NFESyncEstadoListParams{
		First:    parseIntNFEGestao(q.Get("first"), 0),
		Rows:     parseIntNFEGestao(q.Get("rows"), 50),
		Provider: strings.TrimSpace(q.Get("provider")),
		UF:       strings.TrimSpace(q.Get("uf")),
		CNPJ:     strings.TrimSpace(q.Get("cnpj")),
	}
	out, err := h.svc.ListSyncEstado(r.Context(), schema, p)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func parseIntNFEGestao(s string, def int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

// ListGestao GET /serpro/nfe/gestao — listagem paginada com filtros (EF-919).
func (h *NFESerproHandler) ListGestao(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	schema := middleware.TenantSchema(r.Context())
	p := repository.NFEGestaoListParams{
		First:            parseIntNFEGestao(q.Get("first"), 0),
		Rows:             parseIntNFEGestao(q.Get("rows"), 25),
		SortField:        strings.TrimSpace(q.Get("sortField")),
		SortOrder:        parseIntNFEGestao(q.Get("sortOrder"), -1),
		TipoArquivo:      strings.TrimSpace(q.Get("tipo_arquivo")),
		ChaveNFe:         strings.TrimSpace(q.Get("chave_nfe")),
		CNPJEmitente:     q.Get("cnpj_emitente"),
		CNPJDestinatario: q.Get("cnpj_destinatario"),
	}
	if v := strings.TrimSpace(q.Get("emissao_ini")); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			p.EmissaoIni = &t
		}
	}
	if v := strings.TrimSpace(q.Get("emissao_fim")); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			p.EmissaoFim = &t
		}
	}
	if p.SortField == "" {
		p.SortField = "data_download"
		p.SortOrder = -1
	}
	out, err := h.svc.ListGestao(r.Context(), schema, p)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func (h *NFESerproHandler) GetDocumento(w http.ResponseWriter, r *http.Request) {
	chave := strings.TrimSpace(r.URL.Query().Get("chave"))
	schema := middleware.TenantSchema(r.Context())
	out, err := h.svc.BuscarDocumento(r.Context(), schema, chave)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func (h *NFESerproHandler) ExportarXML(w http.ResponseWriter, r *http.Request) {
	chave := strings.TrimSpace(r.URL.Query().Get("chave"))
	schema := middleware.TenantSchema(r.Context())
	xmlPayload, err := h.svc.ExportarXML(r.Context(), schema, chave)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xmlPayload))
}

func (h *NFESerproHandler) ExportarDanfeHTML(w http.ResponseWriter, r *http.Request) {
	chave := strings.TrimSpace(r.URL.Query().Get("chave"))
	schema := middleware.TenantSchema(r.Context())
	html, err := h.svc.ExportarDanfeHTML(r.Context(), schema, chave)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

func (h *NFESerproHandler) GetDanfeJSON(w http.ResponseWriter, r *http.Request) {
	chave := strings.TrimSpace(r.URL.Query().Get("chave"))
	if len(chave) != 44 || !hasOnlyDigits(chave) {
		h.writeNFEError(w, r, errNFEDadosInvalidos("chave_nfe deve ter 44 digitos numericos"), "area_dados")
		return
	}
	schema := middleware.TenantSchema(r.Context())
	view, err := h.svc.BuildDanfeView(r.Context(), schema, chave)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, view)
}

type danfeHTMLFromBodyRequest struct {
	XML     string `json:"xml"`
	Retorno string `json:"retorno"`
}

// GerarDanfeHTMLFromXMLBody gera DANFE a partir do quadro Retorno (JSON trial com payload_json, XML SEFAZ, ou campo xml legado).
func (h *NFESerproHandler) GerarDanfeHTMLFromXMLBody(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)
	var req danfeHTMLFromBodyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeNFEError(w, r, err, "mensagem_inicial_webservice")
		return
	}
	src := strings.TrimSpace(req.Retorno)
	if src == "" {
		src = strings.TrimSpace(req.XML)
	}
	if src == "" {
		h.writeNFEError(w, r, errNFEDadosInvalidos("informe retorno ou xml"), "area_dados")
		return
	}
	xmlPayload, err := service.DanfeXMLFromConsultaRetorno(src)
	if err != nil {
		h.writeNFEError(w, r, err, "area_dados")
		return
	}
	html, err := h.svc.GerarDanfeHTMLFromXML(r.Context(), xmlPayload)
	if err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

func (h *NFESerproHandler) PushNotificacao(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		h.writeNFEError(w, r, err, "mensagem_inicial_webservice")
		return
	}
	headers := map[string]string{}
	for k, vals := range r.Header {
		headers[k] = strings.Join(vals, ",")
	}
	if err := h.svc.RegistrarPushNotificacao(r.Context(), rawBody, headers); err != nil {
		h.writeNFEError(w, r, err, "regras_negocio_webservice")
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
