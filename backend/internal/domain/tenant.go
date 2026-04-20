package domain

type TenantEntity struct {
	ID         string `json:"id"`
	Nome       string `json:"nome"`
	Contato    string `json:"contato"`
	Active     bool   `json:"active"`
	Plano      string `json:"plano"`
	SchemaName string `json:"schema_name,omitempty"`
}

// TenantListRow agrega dados de tenant_dados por schema para listagem (SUPER).
type TenantListRow struct {
	ID          string `json:"id"`
	Nome        string `json:"nome"`
	Contato     string `json:"contato"`
	Active      bool   `json:"active"`
	Plano       string `json:"plano"`
	CNPJ        string `json:"cnpj"`
	RazaoSocial string `json:"razaosocial"`
	Fantasia    string `json:"fantasia"`
}
