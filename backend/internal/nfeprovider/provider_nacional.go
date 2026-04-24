package nfeprovider

import (
	"context"
	"fmt"
	"strings"
)

type NacionalProvider struct{}

func NewNacionalProvider() *NacionalProvider {
	return &NacionalProvider{}
}

func (p *NacionalProvider) ConfigurarCertificado(_ []byte, _ string) error {
	return nil
}

func (p *NacionalProvider) SincronizarDocumentos(_ context.Context, cnpj string, _ string) (ResultadoSincronizacao, error) {
	cnpj = strings.TrimSpace(cnpj)
	return ResultadoSincronizacao{}, fmt.Errorf("provider nacional ainda nao implementado: seguira NT 2014.002 e exigira manifestacao do destinatario para XML completo (cnpj=%s)", cnpj)
}
