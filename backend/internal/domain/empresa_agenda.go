package domain

type EmpresaAgendaItem struct {
	ID             string   `json:"id"`
	EmpresaID      string   `json:"empresa_id"`
	TemplateID     string   `json:"template_id"`
	Descricao      string   `json:"descricao"`
	DataVencimento string   `json:"data_vencimento"`
	Status         string   `json:"status"`
	ValorEstimado  *float64 `json:"valor_estimado"`
}

type EmpresaAgendaAcompanhamentoItem struct {
	EmpresaID      string   `json:"empresa_id"`
	EmpresaNome    string   `json:"empresa_nome"`
	CompromissoID  string   `json:"compromisso_id"`
	Descricao      string   `json:"descricao"`
	DataVencimento string   `json:"data_vencimento"`
	Status         string   `json:"status"`
	Tipo           string   `json:"tipo"`          // TRIBUTARIA | INFORMATIVA (template; TRIBUTO legado)
	Classificacao  string   `json:"classificacao"` // FINANCEIRO | NAO_FINANCEIRO | vazio se sem instância
	AgendaItemID   string   `json:"agenda_item_id"`
	ValorEstimado  *float64 `json:"valor_estimado"`
}
