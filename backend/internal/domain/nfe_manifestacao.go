package domain

import "time"

// NFEManifestacaoDest histórico de manifestação do destinatário (Recepção de Evento SVRS, tpEvento 2102xx).
type NFEManifestacaoDest struct {
	ID            string    `json:"id"`
	ChaveNFe      string    `json:"chave_nfe"`
	TpEvento      string    `json:"tp_evento"`
	CNPJDest      string    `json:"cnpj_dest"`
	CStatLote     int       `json:"cstat_lote"`
	XMotivoLote   string    `json:"x_motivo_lote,omitempty"`
	CStatEvento   int       `json:"cstat_evento"`
	XMotivoEvento string    `json:"x_motivo_evento,omitempty"`
	NProt         string    `json:"n_prot,omitempty"`
	RetornoXML    string    `json:"retorno_xml,omitempty"`
	CriadoEm      time.Time `json:"criado_em"`
}

type NFEManifestacaoListResponse struct {
	Items        []NFEManifestacaoDest `json:"items"`
	TotalRecords int64                 `json:"totalRecords"`
}
