package handlers

import (
	"net/http"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/repository"
)

type GlobalMonitorHandler struct {
	audit *repository.VecxAuditRepository
}

func NewGlobalMonitorHandler(audit *repository.VecxAuditRepository) *GlobalMonitorHandler {
	return &GlobalMonitorHandler{audit: audit}
}

// ListActiveSessions lista sessoes ativas apenas a partir do VECX_AUDIT (EF-929).
func (h *GlobalMonitorHandler) ListActiveSessions(w http.ResponseWriter, r *http.Request) {
	rows, err := h.audit.ListActiveSessions(r.Context())
	if err != nil {
		render.WriteError(w, http.StatusServiceUnavailable, "Nao foi possivel consultar o banco de auditoria global")
		return
	}
	render.WriteJSON(w, http.StatusOK, rows)
}
