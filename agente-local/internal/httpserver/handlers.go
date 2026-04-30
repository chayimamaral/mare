package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/chayimamaral/vecx/agente-local/internal/domain"
	"github.com/chayimamaral/vecx/agente-local/internal/usecase"
)

type Handler struct {
	signUC *usecase.SignUseCase
}

func NewHandler(signUC *usecase.SignUseCase) *Handler {
	return &Handler{signUC: signUC}
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListCertificates(w http.ResponseWriter, r *http.Request) {
	certs, err := h.signUC.ListCertificates(r.Context())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, certs)
}

func (h *Handler) Sign(w http.ResponseWriter, r *http.Request) {
	var in domain.SignInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "corpo invalido"})
		return
	}
	result, err := h.signUC.Sign(r.Context(), in)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
