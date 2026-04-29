package handlers

import (
	"net/http"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/service"
)

type IntegraContadorServicoProcuracaoHandler struct {
	svc *service.IntegraContadorServicoProcuracaoService
}

func NewIntegraContadorServicoProcuracaoHandler(svc *service.IntegraContadorServicoProcuracaoService) *IntegraContadorServicoProcuracaoHandler {
	return &IntegraContadorServicoProcuracaoHandler{svc: svc}
}

func (h *IntegraContadorServicoProcuracaoHandler) List(w http.ResponseWriter, r *http.Request) {
	idSistema := r.URL.Query().Get("id_sistema")
	idServico := r.URL.Query().Get("id_servico")
	items, err := h.svc.List(r.Context(), idSistema, idServico)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"servicos_procuracao": items})
}
