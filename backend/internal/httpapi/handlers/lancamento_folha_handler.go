package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/service"
)

type LancamentoFolhaHandler struct {
	service *service.LancamentoFolhaService
}

type lancamentoFolhaEnvelope struct {
	Params struct {
		ID              string  `json:"id"`
		ClienteID       string  `json:"cliente_id"`
		Competencia     string  `json:"competencia"`
		ValorFolha      float64 `json:"valor_folha"`
		ValorFaturamento float64 `json:"valor_faturamento"`
		Observacoes     string  `json:"observacoes"`
	} `json:"params"`
}

func NewLancamentoFolhaHandler(svc *service.LancamentoFolhaService) *LancamentoFolhaHandler {
	return &LancamentoFolhaHandler{service: svc}
}

// ListTree retorna a arvore de clientes com lancamentos mensais.
func (h *LancamentoFolhaHandler) ListTree(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	if strings.TrimSpace(tenantID) == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}

	response, err := h.service.ListTree(r.Context(), tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

// Create cria um novo lancamento de folha.
func (h *LancamentoFolhaHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	if strings.TrimSpace(tenantID) == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}

	var payload lancamentoFolhaEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if strings.TrimSpace(payload.Params.ClienteID) == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o cliente")
		return
	}

	if strings.TrimSpace(payload.Params.Competencia) == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar a competencia (MM/AAAA)")
		return
	}

	competencia, err := parseCompetencia(payload.Params.Competencia)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, "Competencia invalida! Use MM/AAAA")
		return
	}

	response, err := h.service.Create(r.Context(), tenantID, payload.Params.ClienteID, competencia, payload.Params.ValorFolha, payload.Params.ValorFaturamento, payload.Params.Observacoes)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

// Update atualiza um lancamento existente.
func (h *LancamentoFolhaHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	if strings.TrimSpace(tenantID) == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}

	var payload lancamentoFolhaEnvelope
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if strings.TrimSpace(payload.Params.ID) == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o id do lancamento")
		return
	}

	competencia, err := parseCompetencia(payload.Params.Competencia)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, "Competencia invalida! Use MM/AAAA")
		return
	}

	response, err := h.service.Update(r.Context(), payload.Params.ID, tenantID, competencia, payload.Params.ValorFolha, payload.Params.ValorFaturamento, payload.Params.Observacoes)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

// Delete remove um lancamento.
func (h *LancamentoFolhaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	if strings.TrimSpace(tenantID) == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		var payload lancamentoFolhaEnvelope
		if err := json.NewDecoder(r.Body).Decode(&payload); err == nil {
			id = strings.TrimSpace(payload.Params.ID)
		}
	}

	if id == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o id do lancamento")
		return
	}

	if err := h.service.Delete(r.Context(), id, tenantID); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, map[string]string{"message": "Lancamento removido com sucesso"})
}

// GetByID retorna um lancamento pelo ID.
func (h *LancamentoFolhaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	if strings.TrimSpace(tenantID) == "" {
		render.WriteError(w, http.StatusUnauthorized, "tenant nao identificado")
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		render.WriteError(w, http.StatusBadRequest, "Favor informar o id do lancamento")
		return
	}

	response, err := h.service.GetByID(r.Context(), id, tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

// parseCompetencia converte "MM/AAAA" para time.Time (primeiro dia do mes).
func parseCompetencia(value string) (time.Time, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return time.Time{}, parseCompetenciaErr()
	}

	mes := strings.TrimSpace(parts[0])
	ano := strings.TrimSpace(parts[1])

	return time.Parse("01/2006", mes+"/"+ano)
}

func parseCompetenciaErr() error {
	return &parseError{"formato invalido, use MM/AAAA"}
}

type parseError struct {
	msg string
}

func (e *parseError) Error() string {
	return e.msg
}
