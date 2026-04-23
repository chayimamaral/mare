package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecontab/backend/internal/httpapi/render"
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
