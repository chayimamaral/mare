package domain

type FeriadoRef struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type FeriadoListItem struct {
	ID        string      `json:"id"`
	Descricao string      `json:"descricao"`
	Data      string      `json:"data"`
	Feriado   string      `json:"feriado"`
	Municipio *FeriadoRef `json:"municipio,omitempty"`
	Estado    *FeriadoRef `json:"estado,omitempty"`
}

type FeriadoMutationItem struct {
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
	Data      string `json:"data"`
	Feriado   string `json:"feriado"`
	Ativo     bool   `json:"ativo"`
}
