package domain

type Representante struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	EmailContato string `json:"email_contato,omitempty"`
	Ativo        bool   `json:"ativo"`
}

type ModuloPlataforma struct {
	ID    string `json:"id"`
	Slug  string `json:"slug"`
	Nome  string `json:"nome"`
	Ordem int    `json:"ordem"`
}

type MatrizAcessoItem struct {
	ModuloID   string `json:"modulo_id"`
	Slug       string `json:"slug"`
	Habilitado bool   `json:"habilitado"`
}
