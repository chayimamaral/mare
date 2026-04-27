package domain

type NFEValidacaoRegra struct {
	ID          string `json:"id"`
	Etapa       string `json:"etapa"`
	CodigoRegra string `json:"codigo_regra"`
	Titulo      string `json:"titulo"`
	Descricao   string `json:"descricao"`
	Severidade  string `json:"severidade"`
	Ordem       int    `json:"ordem"`
	Ativo       bool   `json:"ativo"`
}

type NFECodigoErro struct {
	ID             string `json:"id"`
	Origem         string `json:"origem"`
	EtapaValidacao string `json:"etapa_validacao"`
	Codigo         string `json:"codigo"`
	Mensagem       string `json:"mensagem"`
	DescricaoTec   string `json:"descricao_tecnica"`
	AcaoSugerida   string `json:"acao_sugerida"`
	HTTPStatus     int    `json:"http_status,omitempty"`
	Ativo          bool   `json:"ativo"`
}
