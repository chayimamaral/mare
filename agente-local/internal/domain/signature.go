package domain

type SignInput struct {
	HashSHA256Base64 string `json:"hash_sha256_base64"`
	CertificateID    string `json:"certificate_id,omitempty"`
	PIN              string `json:"pin,omitempty"`
}

type SignResult struct {
	Algorithm        string `json:"algorithm"`
	SignatureBase64  string `json:"signature_base64"`
	SelectedCertID   string `json:"selected_cert_id"`
	SelectedCertName string `json:"selected_cert_name"`
}
