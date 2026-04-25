package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/nfeprovider"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
)

func nfeTpAmbSefaz(ambiente string) int {
	a := strings.ToLower(strings.TrimSpace(ambiente))
	if a == "homologacao" || a == "hom" || a == "trial" {
		return 2
	}
	return 1
}

func homologacaoFromAmbiente(ambiente string) bool {
	return nfeTpAmbSefaz(ambiente) == 2
}

const maxManifestacaoRetornoXML = 120_000

func clipManifestacaoXML(s string) string {
	if len(s) <= maxManifestacaoRetornoXML {
		return s
	}
	return s[:maxManifestacaoRetornoXML] + "…"
}

// ManifestarDestinatario envia manifestação (tpEvento 2102xx) ao SVRS RecepcaoEvento 4.00 (UF SC / demais que usam SVRS).
func (s *NFESerproService) ManifestarDestinatario(
	ctx context.Context,
	schemaName, tenantID, chaveNFe, tpEvento, cnpjDest, ambiente, xJust string,
	nSeqEvento int,
	simular bool,
) (domain.NFEManifestacaoDest, error) {
	if err := validateNFEChave(chaveNFe); err != nil {
		return domain.NFEManifestacaoDest{}, err
	}
	if strings.TrimSpace(tenantID) == "" {
		return domain.NFEManifestacaoDest{}, fmt.Errorf("tenant nao encontrado no contexto")
	}
	cnpjDest = onlyDigitsSync(cnpjDest)
	if len(cnpjDest) != 14 && len(cnpjDest) != 11 {
		return domain.NFEManifestacaoDest{}, fmt.Errorf("cnpj/cpf destinatario invalido")
	}
	tpAmb := nfeTpAmbSefaz(ambiente)
	homolog := homologacaoFromAmbiente(ambiente)

	if simular {
		fakeXML := `<?xml version="1.0" encoding="UTF-8"?><retEnvEvento versao="1.00" xmlns="http://www.portalfiscal.inf.br/nfe"><cStat>128</cStat><xMotivo>Lote processado (simulado)</xMotivo><retEvento versao="1.00"><infEvento><cStat>135</cStat><xMotivo>Evento registrado e vinculado a NF-e (simulado)</xMotivo><nProt>000000000000000</nProt></infEvento></retEvento></retEnvEvento>`
		return s.repo.InsertManifestacaoDest(ctx, schemaName, repository.NFEManifestacaoInsert{
			ChaveNFe:      strings.TrimSpace(chaveNFe),
			TpEvento:      strings.TrimSpace(tpEvento),
			CNPJDest:      cnpjDest,
			CStatLote:     128,
			XMotivoLote:   "Lote processado (simulado)",
			CStatEvento:   135,
			XMotivoEvento: "Evento registrado e vinculado a NF-e (simulado)",
			NProt:         "000000000000000",
			RetornoXML:    fakeXML,
		})
	}

	if s.certSvc == nil || !s.certSvc.Configurado() {
		return domain.NFEManifestacaoDest{}, fmt.Errorf("certificado A1 nao configurado para mTLS")
	}
	material, err := s.certSvc.MaterialEmMemoria(ctx, tenantID)
	if err != nil {
		return domain.NFEManifestacaoDest{}, fmt.Errorf("material de certificado: %w", err)
	}
	defer material.Zero()

	client, tlsCert, err := nfeprovider.NewHTTPClientWithPFX(material.PFX, material.Senha, 120*time.Second)
	if err != nil {
		return domain.NFEManifestacaoDest{}, err
	}

	res, err := nfeprovider.EnvManifestacaoDestinatarioSVRS(ctx, client, homolog, tlsCert, nfeprovider.ManifestacaoDestinatarioRequest{
		ChaveNFe:    chaveNFe,
		TpEvento:    tpEvento,
		CNPJCPFDest: cnpjDest,
		TpAmb:       tpAmb,
		XJust:       xJust,
		NSeqEvento:  nSeqEvento,
	})
	if err != nil {
		return domain.NFEManifestacaoDest{}, err
	}

	return s.repo.InsertManifestacaoDest(ctx, schemaName, repository.NFEManifestacaoInsert{
		ChaveNFe:      strings.TrimSpace(chaveNFe),
		TpEvento:      strings.TrimSpace(tpEvento),
		CNPJDest:      cnpjDest,
		CStatLote:     res.CStatLote,
		XMotivoLote:   res.XMotivoLote,
		CStatEvento:   res.CStatEvento,
		XMotivoEvento: res.XMotivoEvento,
		NProt:         res.NProt,
		RetornoXML:    clipManifestacaoXML(res.RetEnvEventoXML),
	})
}

// ListManifestacaoDest histórico por chave (mais recentes primeiro).
func (s *NFESerproService) ListManifestacaoDest(ctx context.Context, schemaName, chaveNFe string) (domain.NFEManifestacaoListResponse, error) {
	if err := validateNFEChave(chaveNFe); err != nil {
		return domain.NFEManifestacaoListResponse{}, err
	}
	items, total, err := s.repo.ListManifestacaoByChave(ctx, schemaName, chaveNFe, 50)
	if err != nil {
		return domain.NFEManifestacaoListResponse{}, err
	}
	return domain.NFEManifestacaoListResponse{Items: items, TotalRecords: total}, nil
}
