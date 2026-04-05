package domain

// CategoriaRotinaPF classifica o template de rotina PF para automação (agenda, integrações futuras).
type CategoriaRotinaPF string

const (
	CategoriaRotinaPFMensalista  CategoriaRotinaPF = "MENSALISTA"
	CategoriaRotinaPFSazonalIRPF CategoriaRotinaPF = "SAZONAL_IRPF"
	CategoriaRotinaPFAvulso      CategoriaRotinaPF = "AVULSO"
)

// RotinaPF é o template de passos/obrigações federais ou sazonais para clientes PF (escopo tenant).
type RotinaPF struct {
	ID        string            `json:"id"`
	TenantID  string            `json:"tenant_id"`
	Nome      string            `json:"nome"`
	Categoria CategoriaRotinaPF `json:"categoria"`
	Descricao string            `json:"descricao,omitempty"`
	Ativo     bool              `json:"ativo"`
}

// RotinaPFLiteItem atende dropdowns (listagem leve por tenant).
type RotinaPFLiteItem struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Categoria string `json:"categoria"`
}

// RotinaPFListRow é uma linha da grade administrativa (todos os status ativo).
type RotinaPFListRow struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Categoria string `json:"categoria"`
	Descricao string `json:"descricao"`
	Ativo     bool   `json:"ativo"`
	CriadoEm  string `json:"criado_em"`
	ItemCount int64  `json:"item_count"`
}

// RotinaPFItemRow é um passo da rotina PF (com rótulo do passo global opcional).
type RotinaPFItemRow struct {
	ID             string `json:"id"`
	RotinaPFID     string `json:"rotina_pf_id"`
	Ordem          int    `json:"ordem"`
	PassoID        string `json:"passo_id"`
	PassoDescricao string `json:"passo_descricao"`
	Descricao      string `json:"descricao"`
	TempoEstimado  int    `json:"tempo_estimado"`
}
