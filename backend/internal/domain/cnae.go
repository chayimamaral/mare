package domain

type CnaeRecord struct {
	ID          string `json:"id"`
	Secao       string `json:"secao"`
	Divisao     string `json:"divisao"`
	Grupo       string `json:"grupo"`
	Classe      string `json:"classe"`
	Subclasse   string `json:"subclasse"`
	Denominacao string `json:"denominacao"`
	Ativo       bool   `json:"ativo"`
}

type CnaeLiteItem struct {
	ID          string `json:"id"`
	Denominacao string `json:"denominacao"`
	Subclasse   string `json:"subclasse"`
}

type CnaeValidateItem struct {
	ID string `json:"id"`
}

// CnaeIbgeResolve descreve uma subclasse encontrada nas tabelas ibge_cnae_* (catálogo oficial).
type CnaeIbgeResolve struct {
	Found       bool   `json:"found"`
	Secao       string `json:"secao"`
	Divisao     string `json:"divisao"`
	Grupo       string `json:"grupo"`
	Classe      string `json:"classe"`
	Denominacao string `json:"denominacao"`
	Subclasse   string `json:"subclasse"`
}
