package domain

type SignInput struct {
	HashSHA256Base64 string `json:"hash_sha256_base64"`
	CertificateID    string `json:"certificate_id,omitempty"`
	PIN              string `json:"pin,omitempty"`
	// EF-937: contexto fiscal (pasta cert_clientes / cert_contador e A3 por titular).
	DocumentID  string `json:"document_id,omitempty"`
	TaxID       string `json:"tax_id,omitempty"`
	Procuracao  bool   `json:"procuracao"`
	SignerTaxID string `json:"signer_tax_id,omitempty"` // CNPJ/CPF do contador quando procuracao + A3
}

type SignResult struct {
	Algorithm        string `json:"algorithm"`
	SignatureBase64  string `json:"signature_base64"`
	SelectedCertID   string `json:"selected_cert_id"`
	SelectedCertName string `json:"selected_cert_name"`
}
