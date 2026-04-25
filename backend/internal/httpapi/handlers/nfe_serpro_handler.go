package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chayimamaral/vecontab/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecontab/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
	"github.com/chayimamaral/vecontab/backend/internal/service"
)

type NFESerproHandler struct {
	svc *service.NFESerproService
}

func NewNFESerproHandler(svc *service.NFESerproService) *NFESerproHandler {
	return &NFESerproHandler{svc: svc}
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

func (h *NFESerproHandler) Consultar(w http.ResponseWriter, r *http.Request) {
	var req nfeConsultaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	schema := middleware.TenantSchema(r.Context())
	out, err := h.svc.ConsultarNFe(r.Context(), schema, req.Ambiente, req.ChaveNFe, req.RequestTag, req.Assinar)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

// SincronizarProvider aciona o provider pattern (EF-920) com estado de NSU por tenant/schema.
func (h *NFESerproHandler) SincronizarProvider(w http.ResponseWriter, r *http.Request) {
	var req nfeSyncProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
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
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

// ManifestarDestinatario POST /serpro/nfe/manifestar-destinatario — Recepção de Evento SVRS (manifestação 2102xx).
func (h *NFESerproHandler) ManifestarDestinatario(w http.ResponseWriter, r *http.Request) {
	var req nfeManifestarDestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
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
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

// ListManifestacaoDest GET /serpro/nfe/manifestacao — histórico por chave.
func (h *NFESerproHandler) ListManifestacaoDest(w http.ResponseWriter, r *http.Request) {
	chave := strings.TrimSpace(r.URL.Query().Get("chave"))
	schema := middleware.TenantSchema(r.Context())
	out, err := h.svc.ListManifestacaoDest(r.Context(), schema, chave)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
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
		render.WriteError(w, http.StatusBadRequest, err.Error())
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
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func (h *NFESerproHandler) GetDocumento(w http.ResponseWriter, r *http.Request) {
	chave := strings.TrimSpace(r.URL.Query().Get("chave"))
	schema := middleware.TenantSchema(r.Context())
	out, err := h.svc.BuscarDocumento(r.Context(), schema, chave)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func (h *NFESerproHandler) ExportarXML(w http.ResponseWriter, r *http.Request) {
	chave := strings.TrimSpace(r.URL.Query().Get("chave"))
	schema := middleware.TenantSchema(r.Context())
	xmlPayload, err := h.svc.ExportarXML(r.Context(), schema, chave)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
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
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
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
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	src := strings.TrimSpace(req.Retorno)
	if src == "" {
		src = strings.TrimSpace(req.XML)
	}
	if src == "" {
		render.WriteError(w, http.StatusBadRequest, "informe retorno ou xml")
		return
	}
	xmlPayload, err := service.DanfeXMLFromConsultaRetorno(src)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	html, err := h.svc.GerarDanfeHTMLFromXML(r.Context(), xmlPayload)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

func (h *NFESerproHandler) PushNotificacao(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, "erro ao ler payload")
		return
	}
	headers := map[string]string{}
	for k, vals := range r.Header {
		headers[k] = strings.Join(vals, ",")
	}
	if err := h.svc.RegistrarPushNotificacao(r.Context(), rawBody, headers); err != nil {
		render.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
