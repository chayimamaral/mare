package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/repository"
	"github.com/chayimamaral/vecx/backend/internal/service"
)

type MatrizConfiguracaoTributariaHandler struct {
	service *service.MatrizConfiguracaoTributariaService
}

type matrizConfiguracaoTributariaEnvelope struct {
	Params struct {
		ID                  string  `json:"id"`
		Nome                string  `json:"nome"`
		NaturezaJuridicaID  string  `json:"natureza_juridica_id"`
		EnquadramentoPorteID string `json:"enquadramento_porte_id"`
		RegimeTributarioID  string  `json:"regime_tributario_id"`
		AliquotaBase        float64 `json:"aliquota_base"`
		PossuiFatorR        bool    `json:"possui_fator_r"`
		AliquotaFatorR      float64 `json:"aliquota_fator_r"`
		Ativo               *bool   `json:"ativo"`
	} `json:"params"`
}

func NewMatrizConfiguracaoTributariaHandler(svc *service.MatrizConfiguracaoTributariaService) *MatrizConfiguracaoTributariaHandler {
	return &MatrizConfiguracaoTributariaHandler{service: svc}
}

func (h *MatrizConfiguracaoTributariaHandler) List(w http.ResponseWriter, r *http.Request) {
	params := repository.MatrizConfiguracaoTributariaListParams{
		First:     parseIntMatriz(r.URL.Query().Get("first"), 0),
		Rows:      parseIntMatriz(r.URL.Query().Get("rows"), 25),
		SortField: r.URL.Query().Get("sortField"),
		SortOrder: parseIntMatriz(r.URL.Query().Get("sortOrder"), 1),
		Nome:      parseNomeFilterMatriz(r.URL.Query().Get("filters")),
	}

	response, err := h.service.List(r.Context(), params)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

func (h *MatrizConfiguracaoTributariaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload matrizConfiguracaoTributariaEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if strings.TrimSpace(payload.Params.Nome) == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o nome da configuracao")
		return
	}
	if strings.TrimSpace(payload.Params.NaturezaJuridicaID) == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar a natureza juridica")
		return
	}
	if strings.TrimSpace(payload.Params.EnquadramentoPorteID) == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o enquadramento/porte")
		return
	}
	if strings.TrimSpace(payload.Params.RegimeTributarioID) == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o regime tributario")
		return
	}

	response, err := h.service.Create(r.Context(), payload.Params.Nome, payload.Params.NaturezaJuridicaID, payload.Params.EnquadramentoPorteID, payload.Params.RegimeTributarioID, payload.Params.AliquotaBase, payload.Params.PossuiFatorR, payload.Params.AliquotaFatorR)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

func (h *MatrizConfiguracaoTributariaHandler) Update(w http.ResponseWriter, r *http.Request) {
	var payload matrizConfiguracaoTributariaEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if strings.TrimSpace(payload.Params.ID) == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o id")
		return
	}

	ativo := true
	if payload.Params.Ativo != nil {
		ativo = *payload.Params.Ativo
	}

	response, err := h.service.Update(r.Context(), payload.Params.ID, payload.Params.Nome, payload.Params.NaturezaJuridicaID, payload.Params.EnquadramentoPorteID, payload.Params.RegimeTributarioID, payload.Params.AliquotaBase, payload.Params.PossuiFatorR, payload.Params.AliquotaFatorR, ativo)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

func (h *MatrizConfiguracaoTributariaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		var payload matrizConfiguracaoTributariaEnvelope
		if err := json.NewDecoder(r.Body).Decode(&payload); err == nil {
			id = strings.TrimSpace(payload.Params.ID)
		}
	}

	if id == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *MatrizConfiguracaoTributariaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o id")
		return
	}

	response, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

func parseIntMatriz(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseNomeFilterMatriz(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}

	type filtersPayload struct {
		Nome struct {
			Value string `json:"value"`
		} `json:"nome"`
	}

	var payload filtersPayload
	if err := json.Unmarshal([]byte(raw), &payload); err == nil {
		return payload.Nome.Value
	}

	return ""
}
