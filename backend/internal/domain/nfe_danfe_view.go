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
	TipoEmissao       string `json:"tipo_emissao"`
	Finalidade        string `json:"finalidade"`
	FormaPagamento    string `json:"forma_pagamento"`
	DigestValue       string `json:"digest_value"`
	DataInclusaoBD    string `json:"data_inclusao_bd"`
}

type NFEDANFEPessoa struct {
	Nome                 string `json:"nome"`
	NomeFantasia         string `json:"nome_fantasia,omitempty"`
	CNPJCPF              string `json:"cnpj_cpf"`
	IE                   string `json:"ie"`
	IEST                 string `json:"ie_substituto,omitempty"`
	IM                   string `json:"im,omitempty"`
	CNAE                 string `json:"cnae,omitempty"`
	CRT                  string `json:"crt,omitempty"`
	CRTDescricao         string `json:"crt_descricao,omitempty"`
	IndicadorIEDest      string `json:"indicador_ie_dest,omitempty"`
	IndicadorIEDescricao string `json:"indicador_ie_descricao,omitempty"`
	Logradouro           string `json:"logradouro"`
	Numero               string `json:"numero"`
	Bairro               string `json:"bairro"`
	Municipio            string `json:"municipio"`
	MunicipioCodigo      string `json:"municipio_codigo,omitempty"`
	MunicipioCodNome     string `json:"municipio_cod_nome,omitempty"`
	UF                   string `json:"uf"`
	CEP                  string `json:"cep"`
	PaisCodigo           string `json:"pais_codigo,omitempty"`
	PaisNome             string `json:"pais_nome,omitempty"`
	PaisCodNome          string `json:"pais_cod_nome,omitempty"`
	Telefone             string `json:"telefone,omitempty"`
	Email                string `json:"email,omitempty"`
	ISUF                 string `json:"isuf,omitempty"`
	CodMunFatoGerador    string `json:"cod_mun_fato_gerador_icms,omitempty"`
	EnderecoCompleto     string `json:"endereco_completo,omitempty"`
}

type NFEDANFEItemICMS struct {
	Origem     string `json:"origem"`
	Tributacao string `json:"tributacao"`
}

type NFEDANFEItemIPI struct {
	ClEnq      string `json:"cl_enq"`
	CEnq       string `json:"c_enq"`
	CSelo      string `json:"c_selo"`
	CNPJProd   string `json:"cnpj_prod"`
	QSelo      string `json:"q_selo"`
	CST        string `json:"cst"`
	QUnid      string `json:"q_unid"`
	VUnid      string `json:"v_unid"`
	VIPI       string `json:"v_ipi"`
	VBC        string `json:"v_bc"`
	PIPI       string `json:"p_ipi"`
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
	IndicadorTotal     string `json:"indicador_total_nf"`
	IndicadorTotalDesc string `json:"indicador_total_desc,omitempty"`
	CEAN        string `json:"cean"`
	CEANTrib    string `json:"cean_trib"`
	UTrib       string `json:"u_trib"`
	QTrib       string `json:"q_trib"`
	VUnTrib     string `json:"v_un_trib"`
	ValorTotTrib string `json:"valor_total_tributos"`
	XPed        string `json:"x_ped,omitempty"`
	NItemPed    string `json:"n_item_ped,omitempty"`
	NFCI        string `json:"n_fci,omitempty"`
	BaseICMS    string `json:"base_icms"`
	ValorICMS   string `json:"valor_icms"`
	ValorIPI    string `json:"valor_ipi"`
	AliquotaICM string `json:"aliquota_icms"`
	AliquotaIPI string `json:"aliquota_ipi"`
	ICMS        NFEDANFEItemICMS `json:"icms"`
	IPI         NFEDANFEItemIPI  `json:"ipi"`
	PISCST      string           `json:"pis_cst,omitempty"`
	COFINSCST   string           `json:"cofins_cst,omitempty"`
}

type NFEDANFETotais struct {
	BaseICMS       string `json:"base_icms"`
	ValorICMS      string `json:"valor_icms"`
	ValorICMSDeson string `json:"valor_icms_desonerado"`
	BaseICMSST     string `json:"base_icms_st"`
	ValorST        string `json:"valor_st"`
	ValorII        string `json:"valor_ii"`
	ValorIPI       string `json:"valor_ipi"`
	ValorPIS       string `json:"valor_pis"`
	ValorCOF       string `json:"valor_cofins"`
	ValorProd      string `json:"valor_produtos"`
	ValorFrete     string `json:"valor_frete"`
	ValorSeg       string `json:"valor_seguro"`
	ValorDesc      string `json:"valor_desconto"`
	ValorOutro     string `json:"valor_outros"`
	ValorTotTrib   string `json:"valor_total_tributos"`
	ValorNF        string `json:"valor_nota"`
}

type NFEDANFETransporte struct {
	Modalidade    string `json:"modalidade"`
	Transportado  string `json:"transportador"`
	CNPJCPF       string `json:"cnpj_cpf"`
	IE            string `json:"ie"`
	Endereco      string `json:"endereco"`
	Municipio     string `json:"municipio"`
	Placa         string `json:"placa"`
	UF            string `json:"uf"`
	RNTC          string `json:"rntc"`
	QuantidadeVol string `json:"quantidade_volumes"`
	Volumes       []NFEDANFEVolume `json:"volumes"`
}

type NFEDANFEAdicionais struct {
	TpImp                    string `json:"tp_imp,omitempty"`
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
	Numero     string `json:"numero"`
	Vencimento string `json:"vencimento"`
	Valor      string `json:"valor"`
}

type NFEDANFEPagamento struct {
	Forma         string `json:"forma"`
	Valor         string `json:"valor"`
	CNPJCred      string `json:"cnpj_credenciadora"`
	Bandeira      string `json:"bandeira"`
	Autorizacao   string `json:"autorizacao"`
}

type NFEDANFEVolume struct {
	Quantidade  string `json:"quantidade"`
	Especie     string `json:"especie"`
	Marca       string `json:"marca"`
	Numero      string `json:"numero"`
	PesoLiquido string `json:"peso_liquido"`
	PesoBruto   string `json:"peso_bruto"`
}
