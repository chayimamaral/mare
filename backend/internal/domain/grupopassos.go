package domain

type GrupoPassosMunicipio struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type GrupoPassosTipoEmpresa struct {
	ID        string `json:"id"`
	Descricao string `json:"descricao"`
}

type GrupoPassosListItem struct {
	ID            string                 `json:"id"`
	Descricao     string                 `json:"descricao"`
	MunicipioID   string                 `json:"municipio_id"`
	TipoEmpresaID string                 `json:"tipoempresa_id"`
	Municipio     GrupoPassosMunicipio   `json:"municipio"`
	TipoEmpresa   GrupoPassosTipoEmpresa `json:"tipoempresa"`
}

type GrupoPassosMutationItem struct {
	ID            string `json:"id"`
	Descricao     string `json:"descricao"`
	MunicipioID   string `json:"municipio_id"`
	TipoEmpresaID string `json:"tipoempresa_id"`
	Ativo         bool   `json:"ativo"`
}
