package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecontab/backend/internal/service"
)

type SalarioMinimoHandler struct {
	service *service.SalarioMinimoService
}

type salarioMinimoEnvelope struct {
	Params struct {
		ID    string  `json:"id"`
		Ano   int     `json:"ano"`
		Valor float64 `json:"valor"`
	} `json:"params"`
}

func NewSalarioMinimoHandler(svc *service.SalarioMinimoService) *SalarioMinimoHandler {
	return &SalarioMinimoHandler{service: svc}
}

func (h *SalarioMinimoHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.List(r.Context())
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func (h *SalarioMinimoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload salarioMinimoEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	item, err := h.service.Create(r.Context(), payload.Params.Ano, payload.Params.Valor)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"salario": item})
}

func (h *SalarioMinimoHandler) Update(w http.ResponseWriter, r *http.Request) {
	var payload salarioMinimoEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	item, err := h.service.Update(r.Context(), payload.Params.ID, payload.Params.Ano, payload.Params.Valor)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"salario": item})
}

func (h *SalarioMinimoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		var payload salarioMinimoEnvelope
		if err := json.NewDecoder(r.Body).Decode(&payload); err == nil {
			id = strings.TrimSpace(payload.Params.ID)
		}
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
