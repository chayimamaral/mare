package domain

type Cidade struct {
	ID     string `json:"id"`
	Nome   string `json:"nome"`
	Codigo string `json:"codigo"`
	UfID   string `json:"ufid"`
	Uf     UfLite `json:"uf,omitempty"`
	Ativo  bool   `json:"ativo,omitempty"`
}

type CidadeListItem struct {
	ID     string `json:"id"`
	Nome   string `json:"nome"`
	Codigo string `json:"codigo"`
	UfID   string `json:"ufid"`
	Uf     UfLite `json:"uf"`
}

type CidadeLiteItem struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type UfLite struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Sigla string `json:"sigla,omitempty"`
}
