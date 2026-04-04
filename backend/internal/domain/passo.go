package domain

type Passo struct {
	ID          string `json:"id"`
	Descricao   string `json:"descricao"`
	Tempo       int    `json:"tempoestimado"`
	TipoPasso   string `json:"tipopasso"`
	MunicipioID string `json:"municipio_id"`
	Link        string `json:"link,omitempty"`
}

type PassoMunicipioRef struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type PassoListItem struct {
	ID          string            `json:"id"`
	Descricao   string            `json:"descricao"`
	Tempo       int               `json:"tempoestimado"`
	TipoPasso   string            `json:"tipopasso"`
	Link        string            `json:"link"`
	MunicipioID string            `json:"municipio_id"`
	Municipio   PassoMunicipioRef `json:"municipio"`
}

type PassoMutationItem struct {
	ID          string `json:"id"`
	Descricao   string `json:"descricao"`
	Tempo       int    `json:"tempoestimado"`
	TipoPasso   string `json:"tipopasso"`
	MunicipioID string `json:"municipio_id"`
	Active      bool   `json:"active"`
}

type PassoDetailItem struct {
	ID          string `json:"id"`
	Descricao   string `json:"descricao"`
	Tempo       int    `json:"tempoestimado"`
	TipoPasso   string `json:"tipopasso"`
	MunicipioID string `json:"municipio_id"`
	Link        string `json:"link"`
}

type PassoCidadeItem struct {
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
	Tempo     int    `json:"tempoestimado"`
	TipoPasso string `json:"tipopasso"`
	RotinaID  string `json:"rotina_id"`
	Ordem     any    `json:"ordem"`
	Link      string `json:"link"`
}
