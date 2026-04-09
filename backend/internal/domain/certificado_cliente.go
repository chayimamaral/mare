package domain

// CertificadoClienteResumo expõe metadados do A1 do cliente (sem blobs nem senha).
type CertificadoClienteResumo struct {
	TipoCertificado string `json:"tipo_certificado"`
	NomeCertificado string `json:"nome_certificado"`
	EmitidoPara     string `json:"emitido_para"`
	EmitidoPor      string `json:"emitido_por"`
	CNPJ            string `json:"cnpj"`
	ValidadeDe      string `json:"validade_de"`
	ValidadeAte     string `json:"validade_ate"`
}
