package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chayimamaral/vecontab/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecontab/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
	"github.com/chayimamaral/vecontab/backend/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
	audit   *repository.VecxAuditRepository
}

func NewAuthHandler(service *service.AuthService, audit *repository.VecxAuditRepository) *AuthHandler {
	return &AuthHandler{service: service, audit: audit}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if input.Email == "" || input.Password == "" {
		render.WriteError(w, http.StatusBadRequest, "Email e password sao obrigatorios")
		return
	}

	response, err := h.service.Login(r.Context(), input, clientIP(r))
	if err != nil {
		if errors.Is(err, service.ErrTenantNaoAutorizadoVecX) {
			render.WriteJSON(w, http.StatusForbidden, map[string]any{
				"error": service.MsgTenantVecxNaoAutorizado,
				"code":  "TENANT_INACTIVE_VECX",
			})
			return
		}
		if errors.Is(err, service.ErrAuditoriaGlobalIndisponivel) {
			render.WriteError(w, http.StatusServiceUnavailable, "Servico de auditoria global indisponivel. Nao foi possivel concluir o login.")
			return
		}
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) AssumeTenant(w http.ResponseWriter, r *http.Request) {
	var input service.AssumeTenantInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	response, err := h.service.AssumeTenant(
		r.Context(),
		middleware.UserID(r.Context()),
		middleware.Role(r.Context()),
		middleware.TenantID(r.Context()),
		middleware.RepresentativeID(r.Context()),
		clientIP(r),
		input,
	)
	if err != nil {
		if errors.Is(err, service.ErrTenantNaoAutorizadoVecX) {
			render.WriteJSON(w, http.StatusForbidden, map[string]any{
				"error": service.MsgTenantVecxNaoAutorizado,
				"code":  "TENANT_INACTIVE_VECX",
			})
			return
		}
		if errors.Is(err, service.ErrAuditoriaGlobalIndisponivel) {
			render.WriteError(w, http.StatusServiceUnavailable, "Servico de auditoria global indisponivel.")
			return
		}
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	render.WriteJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) SessionEnd(w http.ResponseWriter, r *http.Request) {
	if err := h.audit.DeactivateUserSession(r.Context(), middleware.UserID(r.Context())); err != nil {
		render.WriteError(w, http.StatusServiceUnavailable, "Auditoria global indisponivel")
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
