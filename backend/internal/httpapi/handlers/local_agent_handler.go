package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/service"
)

type LocalAgentHandler struct {
	client *service.LocalAgentClient
}

func NewLocalAgentHandler(client *service.LocalAgentClient) *LocalAgentHandler {
	return &LocalAgentHandler{client: client}
}

func (h *LocalAgentHandler) Certificates(w http.ResponseWriter, r *http.Request) {
	out, err := h.client.ListCertificates(r.Context())
	if err != nil {
		render.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{
		"items":        out,
		"totalRecords": len(out),
	})
}

func (h *LocalAgentHandler) SignHash(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RawText       string `json:"raw_text"`
		CertificateID string `json:"certificate_id"`
		PIN           string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	raw := strings.TrimSpace(req.RawText)
	if raw == "" {
		render.WriteError(w, http.StatusBadRequest, "raw_text obrigatorio")
		return
	}
	sum := sha256.Sum256([]byte(raw))
	hashB64 := base64.StdEncoding.EncodeToString(sum[:])
	out, err := h.client.SignHash(r.Context(), service.LocalAgentSignRequest{
		HashSHA256Base64: hashB64,
		CertificateID:    strings.TrimSpace(req.CertificateID),
		PIN:              req.PIN,
	})
	if err != nil {
		render.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, map[string]any{
		"hash_sha256_base64": hashB64,
		"signature":          out,
	})
}
