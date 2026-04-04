package domain

type TenantEntity struct {
	ID      string `json:"id"`
	Nome    string `json:"nome"`
	Contato string `json:"contato"`
	Active  bool   `json:"active"`
	Plano   string `json:"plano"`
}
