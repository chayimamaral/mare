package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecontab/backend/internal/service"
)

type CatalogoServicoHandler struct {
	service *service.CatalogoServicoService
}

type catalogoServicoEnvelope struct {
	Params service.CatalogoServicoInput `json:"params"`
}

func NewCatalogoServicoHandler(s *service.CatalogoServicoService) *CatalogoServicoHandler {
	return &CatalogoServicoHandler{service: s}
}

func (h *CatalogoServicoHandler) List(w http.ResponseWriter, r *http.Request) {
	secao := strings.TrimSpace(r.URL.Query().Get("secao"))
	items, err := h.service.List(r.Context(), secao)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{
		"servicos":     items,
		"totalRecords": len(items),
	})
}

func (h *CatalogoServicoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload catalogoServicoEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	item, err := h.service.Create(r.Context(), payload.Params)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"servico": item})
}

func (h *CatalogoServicoHandler) Update(w http.ResponseWriter, r *http.Request) {
	var payload catalogoServicoEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	item, err := h.service.Update(r.Context(), payload.Params)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"servico": item})
}

func (h *CatalogoServicoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var payload catalogoServicoEnvelope
	_ = json.NewDecoder(r.Body).Decode(&payload)
	id := strings.TrimSpace(payload.Params.ID)
	if id == "" {
		id = strings.TrimSpace(r.URL.Query().Get("id"))
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"success": true})
}
