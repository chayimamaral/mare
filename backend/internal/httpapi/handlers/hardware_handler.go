package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecontab/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecontab/backend/internal/service"
)

type HardwareHandler struct {
	svc *service.HardwareService
}

func NewHardwareHandler(svc *service.HardwareService) *HardwareHandler {
	return &HardwareHandler{svc: svc}
}

func (h *HardwareHandler) ScanLocalDevices(w http.ResponseWriter, r *http.Request) {
	devs, err := h.svc.ListLocalDevices(r.Context())
	if err != nil {
		render.WriteError(w, http.StatusInternalServerError, "falha ao escanear dispositivos locais")
		return
	}

	// Auditoria simples: identifica quem e de qual tenant realizou o scan.
	log.Printf(
		"[AUDIT] hardware_scan user_id=%s user=%s role=%s tenant_id=%s total=%d",
		middleware.UserID(r.Context()),
		middleware.UserName(r.Context()),
		strings.ToUpper(strings.TrimSpace(middleware.Role(r.Context()))),
		middleware.TenantID(r.Context()),
		len(devs),
	)

	render.WriteJSON(w, http.StatusOK, map[string]any{
		"items": devs,
	})
}

