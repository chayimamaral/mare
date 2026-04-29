package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type CaixaPostalHandler struct {
	svc *service.CaixaPostalService
}

func NewCaixaPostalHandler(svc *service.CaixaPostalService) *CaixaPostalHandler {
	return &CaixaPostalHandler{svc: svc}
}

func (h *CaixaPostalHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	schema := middleware.TenantSchema(r.Context())
	tenantID := middleware.TenantID(r.Context())
	count, err := h.svc.Count(r.Context(), schema, tenantID)
	if err != nil {
		render.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	render.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *CaixaPostalHandler) List(w http.ResponseWriter, r *http.Request) {
	schema := middleware.TenantSchema(r.Context())
	tenantID := middleware.TenantID(r.Context())
	msg, err := h.svc.ListMensagens(r.Context(), schema, tenantID)
	if err != nil {
		render.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	render.WriteJSON(w, http.StatusOK, msg)
}

func (h *CaixaPostalHandler) Send(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantID string `json:"tenant_id"` // vazio se global ou user normal
		IsGlobal bool   `json:"is_global"`
		Titulo   string `json:"titulo"`
		Conteudo string `json:"conteudo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.WriteError(w, http.StatusBadRequest, "Corpo inválido")
		return
	}

	userID := middleware.UserID(r.Context())
	userName := middleware.UserName(r.Context())
	requesterRole := middleware.Role(r.Context())
	currentSchema := middleware.TenantSchema(r.Context())

	err := h.svc.Enviar(r.Context(), body.TenantID, body.IsGlobal, body.Titulo, body.Conteudo, userID, userName, requesterRole, currentSchema)
	if err != nil {
		render.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	render.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *CaixaPostalHandler) Read(w http.ResponseWriter, r *http.Request) {
	msgID := chi.URLParam(r, "id")
	if msgID == "" {
		render.WriteError(w, http.StatusBadRequest, "ID da mensagem não fornecido")
		return
	}

	schema := middleware.TenantSchema(r.Context())
	userID := middleware.UserID(r.Context())

	err := h.svc.MarkAsRead(r.Context(), schema, msgID, userID)
	if err != nil {
		render.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	render.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
