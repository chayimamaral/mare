package domain

import "time"

type LancamentoFolha struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	ClienteID       string    `json:"cliente_id"`
	ClienteNome     string    `json:"cliente_nome,omitempty"`
	ClienteDocumento string   `json:"cliente_documento,omitempty"`
	Competencia     string    `json:"competencia"`
	ValorFolha      float64   `json:"valor_folha"`
	ValorFaturamento float64  `json:"valor_faturamento"`
	Observacoes     string    `json:"observacoes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LancamentoFolhaTreeNode struct {
	Key      string                      `json:"key"`
	Data     any                         `json:"data"`
	Children []LancamentoFolhaTreeNode   `json:"children,omitempty"`
	Leaf     bool                        `json:"leaf"`
}

type LancamentoFolhaClienteNode struct {
	Tipo            string  `json:"tipo"`             // "cliente"
	ClienteID       string  `json:"cliente_id"`
	Nome            string  `json:"nome"`
	Documento       string  `json:"documento"`
	TotalFolha      float64 `json:"total_folha"`
	TotalFaturamento float64 `json:"total_faturamento"`
}

type LancamentoFolhaLancamentoNode struct {
	Tipo            string  `json:"tipo"`             // "lancamento"
	ID              string  `json:"id"`
	Competencia     string  `json:"competencia"`
	ValorFolha      float64 `json:"valor_folha"`
	ValorFaturamento float64 `json:"valor_faturamento"`
	Observacoes     string  `json:"observacoes"`
}
