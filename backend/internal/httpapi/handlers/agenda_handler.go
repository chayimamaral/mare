package handlers

import (
	"net/http"
	"strings"

	"github.com/chayimamaral/vecontab/backendgo/internal/httpapi/middleware"
	"github.com/chayimamaral/vecontab/backendgo/internal/httpapi/render"
	"github.com/chayimamaral/vecontab/backendgo/internal/service"
)

type AgendaHandler struct {
	service *service.AgendaService
}

func NewAgendaHandler(service *service.AgendaService) *AgendaHandler {
	return &AgendaHandler{service: service}
}

func (h *AgendaHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	response, err := h.service.List(r.Context(), tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

func (h *AgendaHandler) Detail(w http.ResponseWriter, r *http.Request) {
	agendaID := strings.TrimSpace(r.URL.Query().Get("agenda_id"))
	response, err := h.service.Detail(r.Context(), middleware.TenantID(r.Context()), agendaID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}
