package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LocalAgentClient struct {
	baseURL      string
	sharedSecret string
	client       *http.Client
}

type LocalAgentCertificate struct {
	ID         string   `json:"id"`
	Label      string   `json:"label"`
	Subject    string   `json:"subject"`
	SerialHex  string   `json:"serial_hex"`
	SlotID     uint     `json:"slot_id"`
	TokenLabel string   `json:"token_label"`
	TaxIDs     []string `json:"tax_ids,omitempty"`
}

type LocalAgentSignRequest struct {
	HashSHA256Base64 string `json:"hash_sha256_base64"`
	CertificateID    string `json:"certificate_id,omitempty"`
	PIN              string `json:"pin,omitempty"`
	DocumentID       string `json:"document_id,omitempty"`
	TaxID            string `json:"tax_id,omitempty"`
	Procuracao       bool   `json:"procuracao"`
	SignerTaxID      string `json:"signer_tax_id,omitempty"`
}

type LocalAgentSignResponse struct {
	Algorithm        string `json:"algorithm"`
	SignatureBase64  string `json:"signature_base64"`
	SelectedCertID   string `json:"selected_cert_id"`
	SelectedCertName string `json:"selected_cert_name"`
}

func NewLocalAgentClient(baseURL string, sharedSecret string, timeout time.Duration) *LocalAgentClient {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		base = "http://127.0.0.1:9999"
	}
	if timeout <= 0 {
		timeout = 8 * time.Second
	}
	return &LocalAgentClient{
		baseURL:      base,
		sharedSecret: strings.TrimSpace(sharedSecret),
		client:       &http.Client{Timeout: timeout},
	}
}

func (c *LocalAgentClient) ListCertificates(ctx context.Context) ([]LocalAgentCertificate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/certificates", nil)
	if err != nil {
		return nil, err
	}
	c.applySecret(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nao foi possivel contactar o vecx-agent em %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(body))
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("vecx-agent recusou o segredo (HTTP 401); confira LOCAL_AGENT_SHARED_SECRET no backend e AGENT_SHARED_SECRET no agente")
		}
		var errPayload struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(body, &errPayload) == nil && errPayload.Error != "" {
			msg = errPayload.Error
		}
		return nil, fmt.Errorf("vecx-agent respondeu HTTP %d: %s", resp.StatusCode, msg)
	}
	var wrapped struct {
		Items []LocalAgentCertificate `json:"items"`
	}
	if err := json.Unmarshal(body, &wrapped); err == nil && bytes.Contains(body, []byte(`"items"`)) {
		return wrapped.Items, nil
	}
	var out []LocalAgentCertificate
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("resposta invalida do agente local: %w", err)
	}
	return out, nil
}

func (c *LocalAgentClient) SignHash(ctx context.Context, in LocalAgentSignRequest) (LocalAgentSignResponse, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return LocalAgentSignResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/sign", bytes.NewReader(b))
	if err != nil {
		return LocalAgentSignResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.applySecret(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return LocalAgentSignResponse{}, fmt.Errorf("nao foi possivel contactar o vecx-agent em %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return LocalAgentSignResponse{}, fmt.Errorf("agente local retornou status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var out LocalAgentSignResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return LocalAgentSignResponse{}, fmt.Errorf("resposta invalida do agente local: %w", err)
	}
	return out, nil
}

func (c *LocalAgentClient) applySecret(req *http.Request) {
	if req == nil {
		return
	}
	if c.sharedSecret == "" {
		return
	}
	req.Header.Set("X-Local-Agent-Secret", c.sharedSecret)
}
