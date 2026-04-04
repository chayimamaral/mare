package domain

type Estado struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Sigla string `json:"sigla"`
	Ativo bool   `json:"ativo,omitempty"`
}
