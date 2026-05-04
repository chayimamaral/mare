package domain

// EnquadramentoJuridicoPorte classifica porte da empresa por faixa de faturamento anual (public).
type EnquadramentoJuridicoPorte struct {
	ID            string   `json:"id"`
	Sigla         string   `json:"sigla"`
	Descricao     string   `json:"descricao"`
	LimiteInicial float64  `json:"limite_inicial"`
	LimiteFinal   *float64 `json:"limite_final"`
	AnoVigencia   int      `json:"ano_vigencia"`
	Ativo         bool     `json:"ativo"`
}
