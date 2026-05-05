package domain

type DadosComplementaresRecord struct {
	Tenantid           string  `json:"tenantid"`
	CNPJ               *string `json:"cnpj,omitempty"`
	CEP                *string `json:"cep,omitempty"`
	Endereco           *string `json:"endereco,omitempty"`
	Bairro             *string `json:"bairro,omitempty"`
	Cidade             *string `json:"cidade,omitempty"`
	Estado             *string `json:"estado,omitempty"`
	Telefone           *string `json:"telefone,omitempty"`
	Email              *string `json:"email,omitempty"`
	IE                 *string `json:"ie,omitempty"`
	IM                 *string `json:"im,omitempty"`
	RazaoSocial        *string `json:"razaosocial,omitempty"`
	Fantasia           *string `json:"fantasia,omitempty"`
	Observacoes        *string `json:"observacoes,omitempty"`
	EnviarResumoMensal bool    `json:"enviar_resumo_mensal"`
}

type RegistroUserRecord struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	TenantID     string `json:"tenantid"`
	TenantSchema string `json:"tenant_schema,omitempty"`
	Active       bool   `json:"active"`
}
