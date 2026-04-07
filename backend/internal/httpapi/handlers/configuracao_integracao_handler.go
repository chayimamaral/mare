package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecontab/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecontab/backend/internal/service"
)

type ConfiguracaoIntegracaoHandler struct {
	service     *service.ConfiguracaoIntegracaoService
	certService *service.CertificadoService
}

func NewConfiguracaoIntegracaoHandler(service *service.ConfiguracaoIntegracaoService, certService *service.CertificadoService) *ConfiguracaoIntegracaoHandler {
	return &ConfiguracaoIntegracaoHandler{service: service, certService: certService}
}

func (h *ConfiguracaoIntegracaoHandler) GetChavesSuper(w http.ResponseWriter, r *http.Request) {
	role := strings.ToUpper(strings.TrimSpace(middleware.Role(r.Context())))
	if role != "SUPER" {
		render.WriteError(w, http.StatusForbidden, "Somente SUPER")
		return
	}
	tenantID := middleware.TenantID(r.Context())
	resp, err := h.service.GetChavesSuper(r.Context(), tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"chaves": resp})
}

func (h *ConfiguracaoIntegracaoHandler) SaveChavesSuper(w http.ResponseWriter, r *http.Request) {
	role := strings.ToUpper(strings.TrimSpace(middleware.Role(r.Context())))
	if role != "SUPER" {
		render.WriteError(w, http.StatusForbidden, "Somente SUPER")
		return
	}
	var payload domain.ChavesSuper
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	payload.TenantID = middleware.TenantID(r.Context())
	if err := h.service.SaveChavesSuper(r.Context(), payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *ConfiguracaoIntegracaoHandler) GetTenantConfiguracoes(w http.ResponseWriter, r *http.Request) {
	role := strings.ToUpper(strings.TrimSpace(middleware.Role(r.Context())))
	if role != "ADMIN" && role != "SUPER" {
		render.WriteError(w, http.StatusForbidden, "Somente ADMIN/SUPER")
		return
	}
	tenantID := middleware.TenantID(r.Context())
	resp, err := h.service.GetTenantConfiguracoes(r.Context(), tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"configuracoes": resp})
}

func (h *ConfiguracaoIntegracaoHandler) SaveTenantConfiguracoes(w http.ResponseWriter, r *http.Request) {
	role := strings.ToUpper(strings.TrimSpace(middleware.Role(r.Context())))
	if role != "ADMIN" && role != "SUPER" {
		render.WriteError(w, http.StatusForbidden, "Somente ADMIN/SUPER")
		return
	}
	var payload domain.TenantConfiguracoes
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	payload.TenantID = middleware.TenantID(r.Context())
	if err := h.service.SaveTenantConfiguracoes(r.Context(), payload); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *ConfiguracaoIntegracaoHandler) UploadCertificadoDigital(w http.ResponseWriter, r *http.Request) {
	role := strings.ToUpper(strings.TrimSpace(middleware.Role(r.Context())))
	if role != "ADMIN" && role != "SUPER" {
		render.WriteError(w, http.StatusForbidden, "Somente ADMIN/SUPER")
		return
	}
	if h.certService == nil {
		render.WriteError(w, http.StatusBadRequest, "Servico de certificado nao configurado")
		return
	}
	if err := r.ParseMultipartForm(16 << 20); err != nil {
		render.WriteError(w, http.StatusBadRequest, "multipart invalido")
		return
	}
	empresaID := strings.TrimSpace(r.FormValue("empresa_id"))
	senha := r.FormValue("senha_certificado")
	cnpj := strings.TrimSpace(r.FormValue("cnpj"))
	titular := strings.TrimSpace(r.FormValue("titular_nome"))
	if empresaID == "" {
		render.WriteError(w, http.StatusBadRequest, "empresa_id obrigatorio")
		return
	}

	file, _, err := r.FormFile("arquivo")
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, "arquivo .pfx obrigatorio")
		return
	}
	defer file.Close()

	pfxBytes, err := io.ReadAll(io.LimitReader(file, 16<<20))
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, "falha ao ler arquivo")
		return
	}
	tenantID := middleware.TenantID(r.Context())
	if err := h.certService.UpsertPFX(r.Context(), tenantID, empresaID, pfxBytes, senha, cnpj, titular); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	mat, err := h.certService.MaterialEmMemoria(r.Context(), tenantID, empresaID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer mat.Zero()

	render.WriteJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"certificado": map[string]any{
			"nome_certificado": mat.Nome,
			"cnpj":             mat.CNPJ,
			"validade_ate":     mat.ValidadeAte.Format("2006-01-02"),
		},
	})
}
