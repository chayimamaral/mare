package nfeprovider

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) ConfigurarCertificado(_ []byte, _ string) error {
	return nil
}

func (p *MockProvider) SincronizarDocumentos(_ context.Context, cnpj string, ultNSU string) (ResultadoSincronizacao, error) {
	cnpj = onlyDigits(cnpj)
	if len(cnpj) != 14 && len(cnpj) != 11 {
		return ResultadoSincronizacao{}, fmt.Errorf("cnpj/cpf invalido para provider mock")
	}
	now := time.Now().UTC().Format(time.RFC3339)
	nsu := "1"
	if strings.TrimSpace(ultNSU) != "" && strings.TrimSpace(ultNSU) != "0" {
		return ResultadoSincronizacao{
			Documentos: nil,
			NovoMaxNSU: ultNSU,
			CStat:      117,
			XMotivo:    "Nenhum DF-e localizado (mock)",
		}, nil
	}
	chave := "35240312345678000199550010000006371000000000"
	docXML := `<nfeProc><NFe><infNFe><ide><mod>55</mod><nNF>637</nNF><dhEmi>` + now + `</dhEmi></ide><emit><CNPJ>` + cnpj + `</CNPJ><xNome>Mock Emitente</xNome></emit><dest><CNPJ>11222333000181</CNPJ></dest><total><ICMSTot><vNF>1500.50</vNF></ICMSTot></total></infNFe></NFe></nfeProc>`
	return ResultadoSincronizacao{
		Documentos: []DocumentoFiscal{
			{
				NSU:         nsu,
				ChaveAcesso: chave,
				Tipo:        "nfeProc",
				XML:         docXML,
				RecebidoEm:  time.Now().UTC(),
			},
		},
		NovoMaxNSU: nsu,
		CStat:      118,
		XMotivo:    "DF-e localizados (mock)",
	}, nil
}
