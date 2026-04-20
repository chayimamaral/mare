package domain

import "time"

type CaixaPostalMensagem struct {
	ID                string     `json:"id"`
	RemetenteID       *string    `json:"remetente_id,omitempty"`
	RemetenteTenantID *string    `json:"remetente_tenantid,omitempty"`
	RemetenteNome     string     `json:"remetente_nome"`
	Tipo              string     `json:"tipo"` // INBOX ou OUTBOX
	IsGlobal          bool       `json:"is_global"`
	Titulo            string     `json:"titulo"`
	Conteudo          string     `json:"conteudo"`
	Lida              bool       `json:"lida"`
	LidaPor           *string    `json:"lida_por,omitempty"` // ID do Usuario que leu primeiro (no tenant)
	LidaEm            *time.Time `json:"lida_em,omitempty"`
	CriadoEm          time.Time  `json:"criado_em"`
}
