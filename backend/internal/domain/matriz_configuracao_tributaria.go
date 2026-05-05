package domain

type MatrizConfiguracaoTributaria struct {
	ID                   string  `json:"id"`
	Nome                 string  `json:"nome"`
	NaturezaJuridicaID   string  `json:"natureza_juridica_id"`
	NaturezaJuridica     string  `json:"natureza_juridica,omitempty"`
	EnquadramentoPorteID string  `json:"enquadramento_porte_id"`
	EnquadramentoPorte   string  `json:"enquadramento_porte,omitempty"`
	RegimeTributarioID   string  `json:"regime_tributario_id"`
	RegimeTributario     string  `json:"regime_tributario,omitempty"`
	AliquotaBase         float64 `json:"aliquota_base"`
	PossuiFatorR         bool    `json:"possui_fator_r"`
	AliquotaFatorR       float64 `json:"aliquota_fator_r"`
	SubstituicaoTributaria bool `json:"substituicao_tributaria"`
	Ativo                bool    `json:"ativo"`
}
