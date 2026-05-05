package domain

type EmpresaRef struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type EmpresaRotinaRef struct {
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
}

type EmpresaRotinaPFRef struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Categoria string `json:"categoria"`
}

// EmpresaMatrizTributariaRef referencia a matriz de configuração tributária do cliente.
type EmpresaMatrizTributariaRef struct {
	ID                     string `json:"id"`
	Nome                   string `json:"nome"`
	SubstituicaoTributaria bool   `json:"substituicao_tributaria"`
}

type EmpresaListItem struct {
	ID         string           `json:"id"`
	Nome       string           `json:"nome"`
	TipoPessoa string           `json:"tipo_pessoa"`
	Documento  string           `json:"documento"`
	IE         string           `json:"ie"`
	IM         string           `json:"im"`
	Municipio  EmpresaRef       `json:"municipio"`
	Rotina     EmpresaRotinaRef `json:"rotina"`
	RotinaPF   EmpresaRotinaPFRef `json:"rotina_pf"`
	// MatrizTributaria vincula o cliente à Matriz de Configuração Tributária,
	// substituindo os campos individuais natureza_juridica, porte e regime_tributario.
	MatrizTributaria EmpresaMatrizTributariaRef `json:"matriz_tributaria"`
	// FaturamentoAcumuladoAno soma valor_total em nfe_gestao (emitente = documento PJ) no ano calendário corrente (sessão DB).
	FaturamentoAcumuladoAno float64 `json:"faturamento_acumulado_ano"`
	Cnaes                   any     `json:"cnaes"`
	Bairro                  string  `json:"bairro"`
	Iniciado                bool    `json:"iniciado"`
	PassosConcluidos        bool    `json:"passos_concluidos"`
	CompromissosGerados     bool    `json:"compromissos_gerados"`
}

type EmpresaMutationItem struct {
	ID          string `json:"id"`
	Nome        string `json:"nome"`
	MunicipioID string `json:"municipio_id"`
	TenantID    string `json:"tenant_id"`
	RotinaID    string `json:"rotina_id"`
	RotinaPFID  string `json:"rotina_pf_id"`
	Cnaes       any    `json:"cnaes"`
	Iniciado    bool   `json:"iniciado"`
	Ativo       bool   `json:"ativo"`
}

type EmpresaProcessoItem struct {
	ID                  string `json:"id"`
	EmpresaID           string `json:"empresa_id"`
	TenantID            string `json:"tenant_id"`
	RotinaID            string `json:"rotina_id"`
	Descricao           string `json:"descricao"`
	CriadoEm            string `json:"criado_em"`
	Iniciado            bool   `json:"iniciado"`
	PassosConcluidos    bool   `json:"passos_concluidos"`
	CompromissosGerados bool   `json:"compromissos_gerados"`
	Ativo               bool   `json:"ativo"`
}
