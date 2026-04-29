package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/service"
)

type SerproServicoEnquadramentoHandler struct {
	service *service.SerproServicoEnquadramentoService
}

type serproServicoEnquadramentoEnvelope struct {
	Params service.SerproServicoEnquadramentoSaveInput `json:"params"`
}

func NewSerproServicoEnquadramentoHandler(s *service.SerproServicoEnquadramentoService) *SerproServicoEnquadramentoHandler {
	return &SerproServicoEnquadramentoHandler{service: s}
}

func (h *SerproServicoEnquadramentoHandler) List(w http.ResponseWriter, r *http.Request) {
	enquadramentoID := strings.TrimSpace(r.URL.Query().Get("enquadramento_id"))
	regimeID := strings.TrimSpace(r.URL.Query().Get("regime_tributario_id"))

	ids, err := h.service.ListServicosIDs(r.Context(), enquadramentoID, regimeID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{
		"servicos_ids": ids,
	})
}

func (h *SerproServicoEnquadramentoHandler) Save(w http.ResponseWriter, r *http.Request) {
	var payload serproServicoEnquadramentoEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	if err := h.service.SaveServicosIDs(r.Context(), payload.Params); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}
