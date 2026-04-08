package domain

type CatalogoServico struct {
	ID                  string `json:"id"`
	Secao               string `json:"secao"`
	Sequencial          int    `json:"sequencial"`
	Codigo              string `json:"codigo"`
	IDSistema           string `json:"id_sistema"`
	IDServico           string `json:"id_servico"`
	SituacaoImplantacao string `json:"situacao_implantacao"`
	DataImplantacao     string `json:"data_implantacao,omitempty"`
	Tipo                string `json:"tipo"`
	Descricao           string `json:"descricao"`
	Ativo               bool   `json:"ativo"`
}
