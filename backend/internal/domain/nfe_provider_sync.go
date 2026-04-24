package domain

import "time"

type NFESyncEstado struct {
	ID                  string     `json:"id"`
	Provider            string     `json:"provider"`
	UF                  string     `json:"uf"`
	CNPJ                string     `json:"cnpj"`
	UltimoNSU           string     `json:"ultimo_nsu"`
	UltimoCStat         int        `json:"ultimo_cstat,omitempty"`
	UltimoMotivo        string     `json:"ultimo_motivo,omitempty"`
	UltimaVerificacao   *time.Time `json:"ultima_verificacao,omitempty"`
	ProximaConsultaApos *time.Time `json:"proxima_consulta_apos,omitempty"`
}

type NFESincronizacaoResultado struct {
	Provider         string `json:"provider"`
	UF               string `json:"uf"`
	CNPJ             string `json:"cnpj"`
	AnteriorNSU      string `json:"anterior_nsu"`
	NovoNSU          string `json:"novo_nsu"`
	TotalRecebidos   int    `json:"total_recebidos"`
	TotalPersistidos int    `json:"total_persistidos"`
	CStat            int    `json:"cstat"`
	XMotivo          string `json:"x_motivo,omitempty"`
}

type NFESyncEstadoListResponse struct {
	Items        []NFESyncEstado `json:"items"`
	TotalRecords int64           `json:"totalRecords"`
}
