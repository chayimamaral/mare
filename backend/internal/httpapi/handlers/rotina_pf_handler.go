package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecontab/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
	"github.com/chayimamaral/vecontab/backend/internal/service"
)

type RotinaPFHandler struct {
	svc *service.RotinaPFService
}

func NewRotinaPFHandler(svc *service.RotinaPFService) *RotinaPFHandler {
	return &RotinaPFHandler{svc: svc}
}

type rotinaPFEnvelope struct {
	Params struct {
		ID          string `json:"id"`
		Nome        string `json:"nome"`
		Categoria   string `json:"categoria"`
		Descricao   string `json:"descricao"`
		Ativo       *bool  `json:"ativo"`
		RotinaPFID  string `json:"rotina_pf_id"`
		ItemID      string `json:"item_id"`
		Ordem       int    `json:"ordem"`
		PassoID     string `json:"passo_id"`
		TempoEstimado int  `json:"tempo_estimado"`
	} `json:"params"`
}

func rotinaPFBoolOr(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}

func (h *RotinaPFHandler) ListLite(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(middleware.TenantID(r.Context()))
	if tenantID == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}
	resp, err := h.svc.ListLite(r.Context(), tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, resp)
}

func (h *RotinaPFHandler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(middleware.TenantID(r.Context()))
	if tenantID == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}
	first, rows := rotinaPFListPaging(r.URL.Query().Get("first"), r.URL.Query().Get("rows"))
	params := repository.RotinaPFListParams{
		First:     first,
		Rows:      rows,
		SortField: r.URL.Query().Get("sortField"),
		SortOrder: parseIntRotinaPF(r.URL.Query().Get("sortOrder"), 1),
		Nome:      parseNomeFilterRotinaPF(r.URL.Query().Get("filters")),
		TenantID:  tenantID,
	}
	resp, err := h.svc.ListAdmin(r.Context(), params)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, resp)
}

func (h *RotinaPFHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(middleware.TenantID(r.Context()))
	if tenantID == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}
	var payload rotinaPFEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	if strings.TrimSpace(payload.Params.Nome) == "" {
		render.WriteError(w, http.StatusBadRequest, "Nome obrigatorio")
		return
	}
	if strings.TrimSpace(payload.Params.Categoria) == "" {
		render.WriteError(w, http.StatusBadRequest, "Categoria obrigatoria (MENSALISTA, SAZONAL_IRPF, AVULSO)")
		return
	}
	resp, err := h.svc.Create(r.Context(), repository.RotinaPFUpsertInput{
		TenantID:  tenantID,
		Nome:      payload.Params.Nome,
		Categoria: payload.Params.Categoria,
		Descricao: payload.Params.Descricao,
		Ativo:     rotinaPFBoolOr(payload.Params.Ativo, true),
	})
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, resp)
}

func (h *RotinaPFHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(middleware.TenantID(r.Context()))
	if tenantID == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}
	var payload rotinaPFEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	if strings.TrimSpace(payload.Params.ID) == "" {
		render.WriteError(w, http.StatusBadRequest, "Id obrigatorio")
		return
	}
	if strings.TrimSpace(payload.Params.Nome) == "" {
		render.WriteError(w, http.StatusBadRequest, "Nome obrigatorio")
		return
	}
	if strings.TrimSpace(payload.Params.Categoria) == "" {
		render.WriteError(w, http.StatusBadRequest, "Categoria obrigatoria")
		return
	}
	resp, err := h.svc.Update(r.Context(), repository.RotinaPFUpsertInput{
		ID:        strings.TrimSpace(payload.Params.ID),
		TenantID:  tenantID,
		Nome:      payload.Params.Nome,
		Categoria: payload.Params.Categoria,
		Descricao: payload.Params.Descricao,
		Ativo:     rotinaPFBoolOr(payload.Params.Ativo, true),
	})
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, resp)
}

func (h *RotinaPFHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(middleware.TenantID(r.Context()))
	if tenantID == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}
	var payload rotinaPFEnvelope
	_ = json.NewDecoder(r.Body).Decode(&payload)
	id := strings.TrimSpace(payload.Params.ID)
	if id == "" {
		id = strings.TrimSpace(r.URL.Query().Get("id"))
	}
	if id == "" {
		render.WriteError(w, http.StatusBadRequest, "Id obrigatorio")
		return
	}
	resp, err := h.svc.SoftDelete(r.Context(), id, tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, resp)
}

func (h *RotinaPFHandler) ListItens(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(middleware.TenantID(r.Context()))
	if tenantID == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}
	rid := strings.TrimSpace(r.URL.Query().Get("rotina_pf_id"))
	if rid == "" {
		render.WriteError(w, http.StatusBadRequest, "rotina_pf_id obrigatorio")
		return
	}
	resp, err := h.svc.ListItens(r.Context(), rid, tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, resp)
}

func (h *RotinaPFHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(middleware.TenantID(r.Context()))
	if tenantID == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}
	var payload rotinaPFEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	rid := strings.TrimSpace(payload.Params.RotinaPFID)
	if rid == "" {
		render.WriteError(w, http.StatusBadRequest, "rotina_pf_id obrigatorio")
		return
	}
	desc := strings.TrimSpace(payload.Params.Descricao)
	passo := strings.TrimSpace(payload.Params.PassoID)
	if passo == "" && desc == "" {
		render.WriteError(w, http.StatusBadRequest, "Informe passo_id ou descricao do item")
		return
	}
	resp, err := h.svc.CreateItem(r.Context(), repository.RotinaPFItemUpsertInput{
		RotinaPFID:    rid,
		TenantID:      tenantID,
		Ordem:         payload.Params.Ordem,
		PassoID:       passo,
		Descricao:     payload.Params.Descricao,
		TempoEstimado: payload.Params.TempoEstimado,
	})
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, resp)
}

func (h *RotinaPFHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(middleware.TenantID(r.Context()))
	if tenantID == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}
	var payload rotinaPFEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	itemID := strings.TrimSpace(payload.Params.ItemID)
	if itemID == "" {
		itemID = strings.TrimSpace(payload.Params.ID)
	}
	if itemID == "" {
		render.WriteError(w, http.StatusBadRequest, "item_id obrigatorio")
		return
	}
	rid := strings.TrimSpace(payload.Params.RotinaPFID)
	if rid == "" {
		render.WriteError(w, http.StatusBadRequest, "rotina_pf_id obrigatorio")
		return
	}
	passo := strings.TrimSpace(payload.Params.PassoID)
	desc := strings.TrimSpace(payload.Params.Descricao)
	if passo == "" && desc == "" {
		render.WriteError(w, http.StatusBadRequest, "Informe passo_id ou descricao do item")
		return
	}
	resp, err := h.svc.UpdateItem(r.Context(), repository.RotinaPFItemUpsertInput{
		ID:            itemID,
		RotinaPFID:    rid,
		TenantID:      tenantID,
		Ordem:         payload.Params.Ordem,
		PassoID:       passo,
		Descricao:     payload.Params.Descricao,
		TempoEstimado: payload.Params.TempoEstimado,
	})
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, resp)
}

func (h *RotinaPFHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimSpace(middleware.TenantID(r.Context()))
	if tenantID == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}
	var payload rotinaPFEnvelope
	_ = json.NewDecoder(r.Body).Decode(&payload)
	itemID := strings.TrimSpace(payload.Params.ItemID)
	if itemID == "" {
		itemID = strings.TrimSpace(payload.Params.ID)
	}
	if itemID == "" {
		itemID = strings.TrimSpace(r.URL.Query().Get("item_id"))
	}
	rid := strings.TrimSpace(payload.Params.RotinaPFID)
	if rid == "" {
		rid = strings.TrimSpace(r.URL.Query().Get("rotina_pf_id"))
	}
	if itemID == "" || rid == "" {
		render.WriteError(w, http.StatusBadRequest, "item_id e rotina_pf_id obrigatorios")
		return
	}
	resp, err := h.svc.DeleteItem(r.Context(), itemID, rid, tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, resp)
}

func rotinaPFListPaging(firstRaw, rowsRaw string) (first int, rows int) {
	first = parseIntRotinaPF(firstRaw, 0)
	rows = parseIntRotinaPF(rowsRaw, 25)
	if rows <= 0 {
		rows = 25
	}
	if rows > 200 {
		rows = 200
	}
	if first < 0 {
		first = 0
	}
	return first, rows
}

func parseIntRotinaPF(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return parsed
}

func parseNomeFilterRotinaPF(raw string) string {
	type filtersPayload struct {
		Nome struct {
			Value string `json:"value"`
		} `json:"nome"`
	}
	if strings.TrimSpace(raw) == "" {
		return ""
	}
	var payload filtersPayload
	if err := json.Unmarshal([]byte(raw), &payload); err == nil {
		return payload.Nome.Value
	}
	return ""
}
