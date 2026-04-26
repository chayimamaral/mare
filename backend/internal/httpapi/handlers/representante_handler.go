package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecontab/backend/internal/service"
)

type RepresentanteHandler struct {
	service *service.RepresentanteService
}

func NewRepresentanteHandler(service *service.RepresentanteService) *RepresentanteHandler {
	return &RepresentanteHandler{service: service}
}

func (h *RepresentanteHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context())
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, items)
}

type representantePayload struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	EmailContato string `json:"email_contato"`
	Ativo        *bool  `json:"ativo"`
}

func (h *RepresentanteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload representantePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	item, err := h.service.Create(r.Context(), payload.Nome, payload.EmailContato)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, item)
}

func (h *RepresentanteHandler) Update(w http.ResponseWriter, r *http.Request) {
	var payload representantePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	ativo := true
	if payload.Ativo != nil {
		ativo = *payload.Ativo
	}
	item, err := h.service.Update(r.Context(), payload.ID, payload.Nome, payload.EmailContato, ativo)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, item)
}

func (h *RepresentanteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		render.WriteError(w, http.StatusBadRequest, "id obrigatorio")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]string{"ok": "true"})
}

func (h *RepresentanteHandler) ListModulos(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.ListModulos(r.Context())
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, items)
}

func (h *RepresentanteHandler) GetMatriz(w http.ResponseWriter, r *http.Request) {
	rid := strings.TrimSpace(r.URL.Query().Get("representante_id"))
	items, err := h.service.GetMatriz(r.Context(), rid)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, items)
}

type matrizPutPayload struct {
	RepresentanteID string                    `json:"representante_id"`
	Itens           []domain.MatrizAcessoItem `json:"itens"`
}

func (h *RepresentanteHandler) PutMatriz(w http.ResponseWriter, r *http.Request) {
	var payload matrizPutPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	if err := h.service.ReplaceMatriz(r.Context(), payload.RepresentanteID, payload.Itens); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]string{"ok": "true"})
}
