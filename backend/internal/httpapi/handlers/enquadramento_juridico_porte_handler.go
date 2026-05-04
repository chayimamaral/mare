package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/service"
)

type EnquadramentoJuridicoPorteHandler struct {
	service *service.EnquadramentoJuridicoPorteService
}

type enquadramentoJuridicoPorteEnvelope struct {
	Params struct {
		ID            string   `json:"id"`
		Sigla         string   `json:"sigla"`
		Descricao     string   `json:"descricao"`
		LimiteInicial float64  `json:"limite_inicial"`
		LimiteFinal   *float64 `json:"limite_final"`
		AnoVigencia   int      `json:"ano_vigencia"`
		Ativo         *bool    `json:"ativo"`
	} `json:"params"`
}

func NewEnquadramentoJuridicoPorteHandler(svc *service.EnquadramentoJuridicoPorteService) *EnquadramentoJuridicoPorteHandler {
	return &EnquadramentoJuridicoPorteHandler{service: svc}
}

func (h *EnquadramentoJuridicoPorteHandler) List(w http.ResponseWriter, r *http.Request) {
	var ano *int
	if raw := strings.TrimSpace(r.URL.Query().Get("ano_vigencia")); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 1995 || v > 2100 {
			render.WriteError(w, http.StatusBadRequest, "ano_vigencia invalido")
			return
		}
		ano = &v
	}
	out, err := h.service.List(r.Context(), ano)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func (h *EnquadramentoJuridicoPorteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload enquadramentoJuridicoPorteEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	item, err := h.service.Create(r.Context(), payload.Params.Sigla, payload.Params.Descricao, payload.Params.LimiteInicial, payload.Params.LimiteFinal, payload.Params.AnoVigencia)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"item": item})
}

func (h *EnquadramentoJuridicoPorteHandler) Update(w http.ResponseWriter, r *http.Request) {
	var payload enquadramentoJuridicoPorteEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	ativo := true
	if payload.Params.Ativo != nil {
		ativo = *payload.Params.Ativo
	}
	item, err := h.service.Update(r.Context(), payload.Params.ID, payload.Params.Sigla, payload.Params.Descricao, payload.Params.LimiteInicial, payload.Params.LimiteFinal, payload.Params.AnoVigencia, ativo)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"item": item})
}

func (h *EnquadramentoJuridicoPorteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		var payload enquadramentoJuridicoPorteEnvelope
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
