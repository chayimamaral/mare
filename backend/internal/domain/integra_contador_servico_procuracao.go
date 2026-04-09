package domain

type IntegraContadorServicoProcuracao struct {
	IDSistema     string `json:"id_sistema"`
	IDServico     string `json:"id_servico"`
	CodProcuracao string `json:"cod_procuracao"`
	NomeServico   string `json:"nome_servico"`
}
