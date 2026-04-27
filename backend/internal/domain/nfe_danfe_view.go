package domain

type NFEDANFEView struct {
	Identificacao NFEDANFEIdentificacao `json:"identificacao"`
	Emitente      NFEDANFEPessoa        `json:"emitente"`
	Destinatario  NFEDANFEPessoa        `json:"destinatario"`
	Itens         []NFEDANFEItem        `json:"itens"`
	Totais        NFEDANFETotais        `json:"totais"`
	Transporte    NFEDANFETransporte    `json:"transporte"`
	Cobranca      NFEDANFECobranca      `json:"cobranca"`
	Adicionais    NFEDANFEAdicionais    `json:"adicionais"`
}

type NFEDANFEIdentificacao struct {
	Chave             string `json:"chave"`
	Modelo            string `json:"modelo"`
	Serie             string `json:"serie"`
	Numero            string `json:"numero"`
	EmissaoEm         string `json:"emissao_em"`
	SaidaEntradaEm    string `json:"saida_entrada_em"`
	Protocolo         string `json:"protocolo"`
	CodigoStatus      string `json:"codigo_status"`
	DataAutorizacao   string `json:"data_autorizacao"`
	EventoDescricao   string `json:"evento_descricao"`
	Ambiente          string `json:"ambiente"`
	Situacao          string `json:"situacao"`
	NaturezaOp        string `json:"natureza_operacao"`
	TipoOperacao      string `json:"tipo_operacao"`
	DestinoOperacao   string `json:"destino_operacao"`
	ConsumidorFinal   string `json:"consumidor_final"`
	PresencaComprador string `json:"presenca_comprador"`
	ProcessoEmissao   string `json:"processo_emissao"`
	VersaoProcesso    string `json:"versao_processo"`
	Finalidade        string `json:"finalidade"`
	FormaPagamento    string `json:"forma_pagamento"`
}

type NFEDANFEPessoa struct {
	Nome            string `json:"nome"`
	CNPJCPF         string `json:"cnpj_cpf"`
	IE              string `json:"ie"`
	IndicadorIEDest string `json:"indicador_ie_dest"`
	Logradouro      string `json:"logradouro"`
	Numero          string `json:"numero"`
	Bairro          string `json:"bairro"`
	Municipio       string `json:"municipio"`
	UF              string `json:"uf"`
	CEP             string `json:"cep"`
}

type NFEDANFEItem struct {
	Codigo      string `json:"codigo"`
	Descricao   string `json:"descricao"`
	NCM         string `json:"ncm"`
	EXTIPI      string `json:"extipi"`
	CFOP        string `json:"cfop"`
	Unidade     string `json:"unidade"`
	Quantidade  string `json:"quantidade"`
	ValorUnit   string `json:"valor_unitario"`
	ValorTotal  string `json:"valor_total"`
	ValorDesc   string `json:"valor_desconto"`
	ValorFrete  string `json:"valor_frete"`
	ValorSeg    string `json:"valor_seguro"`
	ValorOutro  string `json:"valor_outros"`
	IndicadorTotal string `json:"indicador_total_nf"`
	CEAN        string `json:"cean"`
	CEANTrib    string `json:"cean_trib"`
	UTrib       string `json:"u_trib"`
	QTrib       string `json:"q_trib"`
	VUnTrib     string `json:"v_un_trib"`
	ValorTotTrib string `json:"valor_total_tributos"`
	BaseICMS    string `json:"base_icms"`
	ValorICMS   string `json:"valor_icms"`
	ValorIPI    string `json:"valor_ipi"`
	AliquotaICM string `json:"aliquota_icms"`
	AliquotaIPI string `json:"aliquota_ipi"`
}

type NFEDANFETotais struct {
	BaseICMS   string `json:"base_icms"`
	ValorICMS  string `json:"valor_icms"`
	ValorICMSDeson string `json:"valor_icms_desonerado"`
	BaseICMSST string `json:"base_icms_st"`
	ValorST    string `json:"valor_st"`
	ValorII    string `json:"valor_ii"`
	ValorIPI   string `json:"valor_ipi"`
	ValorPIS   string `json:"valor_pis"`
	ValorCOF   string `json:"valor_cofins"`
	ValorProd  string `json:"valor_produtos"`
	ValorFrete string `json:"valor_frete"`
	ValorSeg   string `json:"valor_seguro"`
	ValorDesc  string `json:"valor_desconto"`
	ValorOutro string `json:"valor_outros"`
	ValorTotTrib string `json:"valor_total_tributos"`
	ValorNF    string `json:"valor_nota"`
}

type NFEDANFETransporte struct {
	Modalidade   string `json:"modalidade"`
	Transportado string `json:"transportador"`
	CNPJCPF      string `json:"cnpj_cpf"`
	IE           string `json:"ie"`
	Endereco     string `json:"endereco"`
	Municipio    string `json:"municipio"`
	Placa        string `json:"placa"`
	UF           string `json:"uf"`
	RNTC         string `json:"rntc"`
	QuantidadeVol string `json:"quantidade_volumes"`
	Volumes      []NFEDANFEVolume `json:"volumes"`
}

type NFEDANFEAdicionais struct {
	InformacoesComplementares string `json:"informacoes_complementares"`
	InformacoesFisco          string `json:"informacoes_fisco"`
}

type NFEDANFECobranca struct {
	NumeroFatura  string             `json:"numero_fatura"`
	ValorOriginal string             `json:"valor_original"`
	ValorDesconto string             `json:"valor_desconto"`
	ValorLiquido  string             `json:"valor_liquido"`
	Duplicatas    []NFEDANFEDuplicata `json:"duplicatas"`
	Pagamentos    []NFEDANFEPagamento `json:"pagamentos"`
}

type NFEDANFEDuplicata struct {
	Numero      string `json:"numero"`
	Vencimento  string `json:"vencimento"`
	Valor       string `json:"valor"`
}

type NFEDANFEPagamento struct {
	Forma         string `json:"forma"`
	Valor         string `json:"valor"`
	CNPJCred      string `json:"cnpj_credenciadora"`
	Bandeira      string `json:"bandeira"`
	Autorizacao   string `json:"autorizacao"`
}

type NFEDANFEVolume struct {
	Quantidade string `json:"quantidade"`
	Especie    string `json:"especie"`
	Marca      string `json:"marca"`
	Numero     string `json:"numero"`
	PesoLiquido string `json:"peso_liquido"`
	PesoBruto  string `json:"peso_bruto"`
}
