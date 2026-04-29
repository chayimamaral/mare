package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/service"
)

type CertificadoClienteHandler struct {
	cert *service.CertificadoService
}

func NewCertificadoClienteHandler(cert *service.CertificadoService) *CertificadoClienteHandler {
	return &CertificadoClienteHandler{cert: cert}
}

func (h *CertificadoClienteHandler) Get(w http.ResponseWriter, r *http.Request) {
	empresaID := strings.TrimSpace(r.URL.Query().Get("empresa_id"))
	if empresaID == "" {
		render.WriteError(w, http.StatusBadRequest, "empresa_id obrigatorio")
		return
	}
	if h.cert == nil || !h.cert.ClienteRepoConfigurado() {
		render.WriteJSON(w, http.StatusOK, map[string]any{"certificado": map[string]any{}})
		return
	}
	tenantID := middleware.TenantID(r.Context())
	resumo, err := h.cert.ResumoCertificadoCliente(r.Context(), tenantID, empresaID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if resumo == nil {
		render.WriteJSON(w, http.StatusOK, map[string]any{"certificado": map[string]any{}})
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"certificado": resumo})
}

func (h *CertificadoClienteHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if h.cert == nil || !h.cert.ClienteRepoConfigurado() || !h.cert.Configurado() {
		render.WriteError(w, http.StatusBadRequest, "servico de certificado nao configurado")
		return
	}
	if err := r.ParseMultipartForm(16 << 20); err != nil {
		render.WriteError(w, http.StatusBadRequest, "multipart invalido")
		return
	}
	empresaID := strings.TrimSpace(r.FormValue("empresa_id"))
	if empresaID == "" {
		render.WriteError(w, http.StatusBadRequest, "empresa_id obrigatorio")
		return
	}
	senha := r.FormValue("senha_certificado")
	cnpj := strings.TrimSpace(r.FormValue("cnpj"))
	titular := strings.TrimSpace(r.FormValue("titular_nome"))

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
	if err := h.cert.UpsertPFXCliente(r.Context(), tenantID, empresaID, pfxBytes, senha, cnpj, titular); err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{"success": true})
}
