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

type EmpresaTipoEmpresaRef struct {
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
}

type EmpresaListItem struct {
	ID                  string                `json:"id"`
	Nome                string                `json:"nome"`
	TipoPessoa          string                `json:"tipo_pessoa"`
	Documento           string                `json:"documento"`
	Municipio           EmpresaRef            `json:"municipio"`
	Rotina              EmpresaRotinaRef      `json:"rotina"`
	RotinaPF            EmpresaRotinaPFRef    `json:"rotina_pf"`
	TipoEmpresa         EmpresaTipoEmpresaRef `json:"tipo_empresa"`
	Cnaes               any                   `json:"cnaes"`
	Bairro              string                `json:"bairro"`
	Iniciado            bool                  `json:"iniciado"`
	PassosConcluidos    bool                  `json:"passos_concluidos"`
	CompromissosGerados bool                  `json:"compromissos_gerados"`
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
