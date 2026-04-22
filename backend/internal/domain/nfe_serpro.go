package domain

import (
	"encoding/json"
	"time"
)

type NFEDocumento struct {
	ID                string          `json:"id"`
	ChaveNFe          string          `json:"chave_nfe"`
	Ambiente          string          `json:"ambiente,omitempty"`
	EventoCodigo      string          `json:"evento_codigo,omitempty"`
	EventoDescricao   string          `json:"evento_descricao,omitempty"`
	Origem            string          `json:"origem"`
	PayloadJSON       json.RawMessage `json:"payload_json"`
	PayloadXML        string          `json:"payload_xml,omitempty"`
	ContentTypeOrigem string          `json:"content_type_origem,omitempty"`
	RequestTag        string          `json:"request_tag,omitempty"`
	StatusHTTP        int             `json:"status_http,omitempty"`
	RecebidoEm        time.Time       `json:"recebido_em"`
	// JaBaixada: true quando o retorno veio do banco do tenant (cache), sem nova chamada à SERPRO.
	JaBaixada bool `json:"ja_baixada,omitempty"`
}

type NFEPushNotificacao struct {
	ID            string          `json:"id"`
	ChaveNFe      string          `json:"chave_nfe,omitempty"`
	DataHoraEnvio *time.Time      `json:"data_hora_envio,omitempty"`
	Payload       json.RawMessage `json:"payload"`
	Headers       json.RawMessage `json:"headers"`
	RecebidoEm    time.Time       `json:"recebido_em"`
}
