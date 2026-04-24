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
}

type NFeProvider interface {
	ConfigurarCertificado(pfx []byte, password string) error
	SincronizarDocumentos(ctx context.Context, cnpj string, ultNSU string) (ResultadoSincronizacao, error)
}
