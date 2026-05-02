package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecx/agente-local/internal/domain"
	"github.com/chayimamaral/vecx/agente-local/internal/usecase"
)

type Handler struct {
	signUC *usecase.SignUseCase
	onEvent func(string)
}

func NewHandler(signUC *usecase.SignUseCase, onEvent func(string)) *Handler {
	return &Handler{signUC: signUC, onEvent: onEvent}
}

func (h *Handler) emit(event string) {
	if h.onEvent == nil {
		return
	}
	h.onEvent(event)
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListCertificates(w http.ResponseWriter, r *http.Request) {
	h.emit("Conexao remota em /certificates de " + r.RemoteAddr)
	certs, err := h.signUC.ListCertificates(r.Context())
	if err != nil {
		h.emit("Falha em /certificates: " + err.Error())
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	h.emit("Consulta de certificados concluida")
	writeJSON(w, http.StatusOK, map[string]any{"items": certs})
}

func (h *Handler) Sign(w http.ResponseWriter, r *http.Request) {
	h.emit("Solicitacao de assinatura recebida em /sign de " + r.RemoteAddr)
	var in domain.SignInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.emit("Falha em /sign: corpo invalido")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "corpo invalido"})
		return
	}
	if doc := strings.TrimSpace(in.DocumentID); doc != "" {
		h.emit("Documento (EF-937): " + doc)
	}
	if tid := strings.TrimSpace(in.TaxID); tid != "" || in.Procuracao {
		h.emit("Resolucao estruturada de certificado (tax_id / procuracao)")
	}
	result, err := h.signUC.Sign(r.Context(), in)
	if err != nil {
		h.emit("Falha em /sign: " + err.Error())
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	h.emit("Assinatura digital concluida")
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
