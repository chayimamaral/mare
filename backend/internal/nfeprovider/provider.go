package nfeprovider

import (
	"context"
	"time"
)

type DocumentoFiscal struct {
	NSU         string
	ChaveAcesso string
	Tipo        string
	XML         string
	RecebidoEm  time.Time
}

type ResultadoSincronizacao struct {
	Documentos []DocumentoFiscal
	NovoMaxNSU string
	CStat      int
	XMotivo    string
	// QtDFeRet: quantidade de DF-e na última resposta retDistNFeSC (BT SC-2021/001). Para backoff 118 vs lote cheio.
	QtDFeRet int
}

type NFeProvider interface {
	ConfigurarCertificado(pfx []byte, password string) error
	SincronizarDocumentos(ctx context.Context, cnpj string, ultNSU string) (ResultadoSincronizacao, error)
}
