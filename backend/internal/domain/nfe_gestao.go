package domain

import "time"

// NFEGestao: resumo persistido por chave para listagem (EF-919).
type NFEGestao struct {
	ID                  string     `json:"id"`
	ChaveNFe            string     `json:"chave_nfe"`
	TipoArquivo         string     `json:"tipo_arquivo"`
	NumeroNFe           string     `json:"numero_nfe"`
	RazaoSocialEmitente string     `json:"razao_social_emitente"`
	CNPJEmitente        string     `json:"cnpj_emitente"`
	DataEmissao         *time.Time `json:"data_emissao,omitempty"`
	CNPJDestinatario    string     `json:"cnpj_destinatario"`
	ValorTotal          *float64   `json:"valor_total,omitempty"`
	DataDownload        time.Time  `json:"data_download"`
}

type NFEGestaoListResponse struct {
	Items        []NFEGestao `json:"items"`
	TotalRecords int64       `json:"totalRecords"`
}
