package domain

type EmpresaCompromissoItem struct {
	ID                     string   `json:"id"`
	Descricao              string   `json:"descricao"`
	Valor                  *float64 `json:"valor,omitempty"`
	Vencimento             string   `json:"vencimento"`
	Observacao             string   `json:"observacao,omitempty"`
	Status                 string   `json:"status"`
	EmpresaID              string   `json:"empresa_id"`
	TipoempresaObrigacaoID string   `json:"tipoempresa_obrigacao_id"`
}

// EmpresaCompromissoAcompanhamentoItem alinha com o dashboard (mesmos nomes JSON).
type EmpresaCompromissoAcompanhamentoItem struct {
	EmpresaID      string   `json:"empresa_id"`
	EmpresaNome    string   `json:"empresa_nome"`
	CompromissoID  string   `json:"compromisso_id"`
	Descricao      string   `json:"descricao"`
	DataVencimento string   `json:"data_vencimento"`
	Status         string   `json:"status"`
	Tipo           string   `json:"tipo"`
	Classificacao  string   `json:"classificacao"`
	AgendaItemID   string   `json:"agenda_item_id"`
	ValorEstimado  *float64 `json:"valor_estimado"`
}

type EmpresaCompromissoEmpresaOption struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type EmpresaCompromissoObrigacaoOption struct {
	ID            string `json:"id"`
	Descricao     string `json:"descricao"`
	Periodicidade string `json:"periodicidade"`
}
