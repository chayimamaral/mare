package service

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
)

type nfeInnerXML struct {
	Inner string `xml:",innerxml"`
}

type nfeIPIXMLElement struct {
	ClEnq string `xml:"clEnq"`
	CEnq  string `xml:"cEnq"`
	IPINT struct {
		CST string `xml:"CST"`
	} `xml:"IPINT"`
	IPITrib struct {
		CST      string `xml:"CST"`
		VBC      string `xml:"vBC"`
		PIPI     string `xml:"pIPI"`
		VIPI     string `xml:"vIPI"`
		CNPJProd string `xml:"CNPJProd"`
		CSelo    string `xml:"cSelo"`
		QSelo    string `xml:"qSelo"`
		QUnid    string `xml:"qUnid"`
		VUnid    string `xml:"vUnid"`
	} `xml:"IPITrib"`
}

type nfeInfXML struct {
	ID  string `xml:"Id,attr"`
	Ide struct {
		NatOp    string `xml:"natOp"`
		Mod      string `xml:"mod"`
		Serie    string `xml:"serie"`
		NNF      string `xml:"nNF"`
		DhEmi    string `xml:"dhEmi"`
		DhSaiEnt string `xml:"dhSaiEnt"`
		TpAmb    string `xml:"tpAmb"`
		TpEmi    string `xml:"tpEmi"`
		TpImp    string `xml:"tpImp"`
		TpNF     string `xml:"tpNF"`
		IDest    string `xml:"idDest"`
		IndPres  string `xml:"indPres"`
		ProcEmi  string `xml:"procEmi"`
		VerProc  string `xml:"verProc"`
		FinNFe   string `xml:"finNFe"`
		CMunFG   string `xml:"cMunFG"`
		IndFinal string `xml:"indFinal"`
	} `xml:"ide"`
	Emit struct {
		XNome string `xml:"xNome"`
		XFant string `xml:"xFant"`
		CNPJ  string `xml:"CNPJ"`
		CPF   string `xml:"CPF"`
		IE    string `xml:"IE"`
		IEST  string `xml:"IEST"`
		IM    string `xml:"IM"`
		CNAE  string `xml:"CNAE"`
		CRT   string `xml:"CRT"`
		Ender struct {
			XLgr    string `xml:"xLgr"`
			Nro     string `xml:"nro"`
			XCpl    string `xml:"xCpl"`
			XBairro string `xml:"xBairro"`
			CMun    string `xml:"cMun"`
			XMun    string `xml:"xMun"`
			UF      string `xml:"UF"`
			CEP     string `xml:"CEP"`
			CPais   string `xml:"cPais"`
			XPais   string `xml:"xPais"`
			Fone    string `xml:"fone"`
		} `xml:"enderEmit"`
	} `xml:"emit"`
	Dest struct {
		XNome     string `xml:"xNome"`
		CNPJ      string `xml:"CNPJ"`
		CPF       string `xml:"CPF"`
		IE        string `xml:"IE"`
		IM        string `xml:"IM"`
		IndIEDest string `xml:"indIEDest"`
		ISUF      string `xml:"ISUF"`
		Email     string `xml:"email"`
		Ender struct {
			XLgr    string `xml:"xLgr"`
			Nro     string `xml:"nro"`
			XCpl    string `xml:"xCpl"`
			XBairro string `xml:"xBairro"`
			CMun    string `xml:"cMun"`
			XMun    string `xml:"xMun"`
			UF      string `xml:"UF"`
			CEP     string `xml:"CEP"`
			CPais   string `xml:"cPais"`
			XPais   string `xml:"xPais"`
			Fone    string `xml:"fone"`
		} `xml:"enderDest"`
	} `xml:"dest"`
	Det []struct {
		Prod struct {
			CProd    string `xml:"cProd"`
			XProd    string `xml:"xProd"`
			NCM      string `xml:"NCM"`
			EXTIPI   string `xml:"EXTIPI"`
			CFOP     string `xml:"CFOP"`
			CEAN     string `xml:"cEAN"`
			UCom     string `xml:"uCom"`
			QCom     string `xml:"qCom"`
			VUnCom   string `xml:"vUnCom"`
			VProd    string `xml:"vProd"`
			VDesc    string `xml:"vDesc"`
			VFrete   string `xml:"vFrete"`
			VSeg     string `xml:"vSeg"`
			VOutro   string `xml:"vOutro"`
			IndTot   string `xml:"indTot"`
			CEANTrib string `xml:"cEANTrib"`
			UTrib    string `xml:"uTrib"`
			QTrib    string `xml:"qTrib"`
			VUnTrib  string `xml:"vUnTrib"`
			XPed     string `xml:"xPed"`
			NItemPed string `xml:"nItemPed"`
			NFCI     string `xml:"nFCI"`
		} `xml:"prod"`
		Imposto struct {
			VTotTrib string           `xml:"vTotTrib"`
			ICMS     nfeInnerXML      `xml:"ICMS"`
			IPI      nfeIPIXMLElement `xml:"IPI"`
			PIS      nfeInnerXML      `xml:"PIS"`
			COFINS   nfeInnerXML      `xml:"COFINS"`
		} `xml:"imposto"`
	} `xml:"det"`
	Total struct {
		ICMSTot struct {
			VBC     string `xml:"vBC"`
			VICMS   string `xml:"vICMS"`
			VICMSDeson string `xml:"vICMSDeson"`
			VBCST   string `xml:"vBCST"`
			VST     string `xml:"vST"`
			VII     string `xml:"vII"`
			VIPI    string `xml:"vIPI"`
			VPIS    string `xml:"vPIS"`
			VCOFINS string `xml:"vCOFINS"`
			VProd   string `xml:"vProd"`
			VFrete  string `xml:"vFrete"`
			VSeg    string `xml:"vSeg"`
			VDesc   string `xml:"vDesc"`
			VOutro  string `xml:"vOutro"`
			VTotTrib string `xml:"vTotTrib"`
			VNF     string `xml:"vNF"`
		} `xml:"ICMSTot"`
	} `xml:"total"`
	Transp struct {
		ModFrete   string `xml:"modFrete"`
		Transporta struct {
			XNome string `xml:"xNome"`
			CNPJ  string `xml:"CNPJ"`
			CPF   string `xml:"CPF"`
			IE    string `xml:"IE"`
			XEnder string `xml:"xEnder"`
			XMun  string `xml:"xMun"`
			UF    string `xml:"UF"`
		} `xml:"transporta"`
		VeicTransp struct {
			Placa string `xml:"placa"`
			UF    string `xml:"UF"`
			RNTC  string `xml:"RNTC"`
		} `xml:"veicTransp"`
		Vol []struct {
			QVol string `xml:"qVol"`
			Esp  string `xml:"esp"`
			Marca string `xml:"marca"`
			NVol string `xml:"nVol"`
			PesoL string `xml:"pesoL"`
			PesoB string `xml:"pesoB"`
		} `xml:"vol"`
	} `xml:"transp"`
	InfAdic struct {
		InfCpl    string `xml:"infCpl"`
		InfAdFisco string `xml:"infAdFisco"`
	} `xml:"infAdic"`
	Pag struct {
		DetPag []nfeDetPagXML `xml:"detPag"`
	} `xml:"pag"`
	Cobr struct {
		Fat struct {
			NFat string `xml:"nFat"`
			VOrig string `xml:"vOrig"`
			VDesc string `xml:"vDesc"`
			VLiq string `xml:"vLiq"`
		} `xml:"fat"`
		Dup []struct {
			NDup string `xml:"nDup"`
			DVenc string `xml:"dVenc"`
			VDup string `xml:"vDup"`
		} `xml:"dup"`
	} `xml:"cobr"`
}

type nfeProcXML struct {
	NFe struct {
		InfNFe nfeInfXML `xml:"infNFe"`
	} `xml:"NFe"`
	ProtNFe struct {
		InfProt struct {
			NProt   string `xml:"nProt"`
			CStat   string `xml:"cStat"`
			XMotivo string `xml:"xMotivo"`
		} `xml:"infProt"`
	} `xml:"protNFe"`
}

type procNFeRootXML struct {
	NFe struct {
		InfNFe nfeInfXML `xml:"infNFe"`
	} `xml:"NFe"`
	ProtNFe struct {
		InfProt struct {
			NProt   string `xml:"nProt"`
			CStat   string `xml:"cStat"`
			XMotivo string `xml:"xMotivo"`
		} `xml:"infProt"`
	} `xml:"protNFe"`
}

type nfeRootXML struct {
	InfNFe  nfeInfXML `xml:"infNFe"`
	ProtNFe struct {
		InfProt struct {
			NProt   string `xml:"nProt"`
			CStat   string `xml:"cStat"`
			XMotivo string `xml:"xMotivo"`
		} `xml:"infProt"`
	} `xml:"protNFe"`
}

type nfeInfProtXML struct {
	NProt   string `xml:"nProt"`
	CStat   string `xml:"cStat"`
	XMotivo string `xml:"xMotivo"`
	DhRecbto string `xml:"dhRecbto"`
}

type nfeDetPagXML struct {
	TPag string `xml:"tPag"`
	XPag string `xml:"xPag"`
	VPag string `xml:"vPag"`
	Card struct {
		CNPJ string `xml:"CNPJ"`
		TBand string `xml:"tBand"`
		CAut string `xml:"cAut"`
	} `xml:"card"`
}

func (s *NFESerproService) BuildDanfeView(ctx context.Context, schemaName, chaveNFe string) (domain.NFEDANFEView, error) {
	doc, err := s.BuscarDocumento(ctx, schemaName, chaveNFe)
	if err != nil {
		return domain.NFEDANFEView{}, err
	}
	xmlStr := strings.TrimSpace(doc.PayloadXML)
	if xmlStr == "" && len(doc.PayloadJSON) > 0 {
		xmlFromJSON, convErr := DanfeXMLFromConsultaRetorno(string(doc.PayloadJSON))
		if convErr == nil {
			xmlStr = strings.TrimSpace(xmlFromJSON)
		}
	}
	if xmlStr == "" {
		return domain.NFEDANFEView{}, fmt.Errorf("documento sem payload_xml/payload_json para danfe")
	}

	inf, prot, err := parseNFEXMLAnyEnvelope(xmlStr)
	if err != nil {
		// Fallback: alguns documentos de distribuição não trazem infNFe no payload_xml,
		// mas possuem estrutura utilizável dentro do payload_json (trial/docZip).
		if len(doc.PayloadJSON) > 0 {
			xmlFromJSON, convErr := DanfeXMLFromConsultaRetorno(string(doc.PayloadJSON))
			if convErr == nil {
				inf2, prot2, err2 := parseNFEXMLAnyEnvelope(strings.TrimSpace(xmlFromJSON))
				if err2 == nil {
					inf = inf2
					prot = prot2
					err = nil
				}
			}
		}
		if err != nil {
			// Última tentativa: envelope completo do documento persistido.
			rawDoc, mErr := json.Marshal(doc)
			if mErr == nil {
				xmlFromDoc, convErr := DanfeXMLFromConsultaRetorno(string(rawDoc))
				if convErr == nil {
					inf2, prot2, err2 := parseNFEXMLAnyEnvelope(strings.TrimSpace(xmlFromDoc))
					if err2 == nil {
						inf = inf2
						prot = prot2
						err = nil
					}
				}
			}
		}
		if err != nil {
			return domain.NFEDANFEView{}, err
		}
	}

	chave := strings.TrimPrefix(strings.TrimSpace(inf.ID), "NFe")
	if chave == "" {
		chave = strings.TrimSpace(chaveNFe)
	}

	crtCod := strings.TrimSpace(inf.Emit.CRT)
	view := domain.NFEDANFEView{
		Identificacao: domain.NFEDANFEIdentificacao{
			Chave:             chave,
			Modelo:            strings.TrimSpace(inf.Ide.Mod),
			Serie:             strings.TrimSpace(inf.Ide.Serie),
			Numero:            strings.TrimSpace(inf.Ide.NNF),
			EmissaoEm:         strings.TrimSpace(inf.Ide.DhEmi),
			SaidaEntradaEm:    strings.TrimSpace(inf.Ide.DhSaiEnt),
			Protocolo:         strings.TrimSpace(prot.NProt),
			CodigoStatus:      strings.TrimSpace(prot.CStat),
			DataAutorizacao:   strings.TrimSpace(prot.DhRecbto),
			EventoDescricao:   eventDescriptionByCStat(strings.TrimSpace(prot.CStat), strings.TrimSpace(prot.XMotivo)),
			Ambiente:          ambienteAutorizacaoLabel(strings.TrimSpace(inf.Ide.TpAmb)),
			Situacao:          statusAtualByCStat(strings.TrimSpace(prot.CStat), strings.TrimSpace(prot.XMotivo)),
			NaturezaOp:        strings.TrimSpace(inf.Ide.NatOp),
			TipoOperacao:      tipoOperacaoLabel(strings.TrimSpace(inf.Ide.TpNF)),
			DestinoOperacao:   destinoOperacaoLabel(strings.TrimSpace(inf.Ide.IDest)),
			ConsumidorFinal:   consumidorFinalLabel(strings.TrimSpace(inf.Ide.IndFinal)),
			PresencaComprador: presencaCompradorLabel(strings.TrimSpace(inf.Ide.IndPres)),
			ProcessoEmissao:   processoEmissaoLabel(strings.TrimSpace(inf.Ide.ProcEmi)),
			VersaoProcesso:    strings.TrimSpace(inf.Ide.VerProc),
			TipoEmissao:       tipoEmissaoLabel(strings.TrimSpace(inf.Ide.TpEmi)),
			Finalidade:        finalidadeLabel(strings.TrimSpace(inf.Ide.FinNFe)),
			FormaPagamento:    firstPayment(inf.Pag.DetPag),
			DigestValue:       extractDigestValue(xmlStr),
			DataInclusaoBD:    formatDataInclusaoBD(doc.RecebidoEm),
		},
		Emitente: domain.NFEDANFEPessoa{
			Nome:                 strings.TrimSpace(inf.Emit.XNome),
			NomeFantasia:         strings.TrimSpace(inf.Emit.XFant),
			CNPJCPF:              firstNotEmpty(inf.Emit.CNPJ, inf.Emit.CPF),
			IE:                   strings.TrimSpace(inf.Emit.IE),
			IEST:                 strings.TrimSpace(inf.Emit.IEST),
			IM:                   strings.TrimSpace(inf.Emit.IM),
			CNAE:                 strings.TrimSpace(inf.Emit.CNAE),
			CRT:                  crtCod,
			CRTDescricao:         crtCodNome(crtCod),
			Logradouro:           strings.TrimSpace(inf.Emit.Ender.XLgr),
			Numero:               strings.TrimSpace(inf.Emit.Ender.Nro),
			Bairro:               strings.TrimSpace(inf.Emit.Ender.XBairro),
			Municipio:            strings.TrimSpace(inf.Emit.Ender.XMun),
			MunicipioCodigo:      strings.TrimSpace(inf.Emit.Ender.CMun),
			MunicipioCodNome:     formatCodNome(inf.Emit.Ender.CMun, inf.Emit.Ender.XMun),
			UF:                   strings.TrimSpace(inf.Emit.Ender.UF),
			CEP:                  strings.TrimSpace(inf.Emit.Ender.CEP),
			PaisCodigo:           strings.TrimSpace(inf.Emit.Ender.CPais),
			PaisNome:             strings.TrimSpace(inf.Emit.Ender.XPais),
			PaisCodNome:          formatCodNome(inf.Emit.Ender.CPais, inf.Emit.Ender.XPais),
			Telefone:             strings.TrimSpace(inf.Emit.Ender.Fone),
			CodMunFatoGerador:    strings.TrimSpace(inf.Ide.CMunFG),
			EnderecoCompleto:     joinLinhaEndereco(inf.Emit.Ender.XLgr, inf.Emit.Ender.Nro, inf.Emit.Ender.XCpl),
		},
		Destinatario: domain.NFEDANFEPessoa{
			Nome:                 strings.TrimSpace(inf.Dest.XNome),
			CNPJCPF:              firstNotEmpty(inf.Dest.CNPJ, inf.Dest.CPF),
			IE:                   strings.TrimSpace(inf.Dest.IE),
			IndicadorIEDest:      strings.TrimSpace(inf.Dest.IndIEDest),
			IndicadorIEDescricao: indIEDestLabel(strings.TrimSpace(inf.Dest.IndIEDest)),
			IM:                   strings.TrimSpace(inf.Dest.IM),
			ISUF:                 strings.TrimSpace(inf.Dest.ISUF),
			Email:                strings.ToLower(strings.TrimSpace(inf.Dest.Email)),
			Logradouro:           strings.TrimSpace(inf.Dest.Ender.XLgr),
			Numero:               strings.TrimSpace(inf.Dest.Ender.Nro),
			Bairro:               strings.TrimSpace(inf.Dest.Ender.XBairro),
			Municipio:            strings.TrimSpace(inf.Dest.Ender.XMun),
			MunicipioCodigo:      strings.TrimSpace(inf.Dest.Ender.CMun),
			MunicipioCodNome:     formatCodNome(inf.Dest.Ender.CMun, inf.Dest.Ender.XMun),
			UF:                   strings.TrimSpace(inf.Dest.Ender.UF),
			CEP:                  strings.TrimSpace(inf.Dest.Ender.CEP),
			PaisCodigo:           strings.TrimSpace(inf.Dest.Ender.CPais),
			PaisNome:             strings.TrimSpace(inf.Dest.Ender.XPais),
			PaisCodNome:          formatCodNome(inf.Dest.Ender.CPais, inf.Dest.Ender.XPais),
			Telefone:             strings.TrimSpace(inf.Dest.Ender.Fone),
			EnderecoCompleto:     joinLinhaEndereco(inf.Dest.Ender.XLgr, inf.Dest.Ender.Nro, inf.Dest.Ender.XCpl),
		},
		Totais: domain.NFEDANFETotais{
			BaseICMS:   strings.TrimSpace(inf.Total.ICMSTot.VBC),
			ValorICMS:  strings.TrimSpace(inf.Total.ICMSTot.VICMS),
			ValorICMSDeson: strings.TrimSpace(inf.Total.ICMSTot.VICMSDeson),
			BaseICMSST: strings.TrimSpace(inf.Total.ICMSTot.VBCST),
			ValorST:    strings.TrimSpace(inf.Total.ICMSTot.VST),
			ValorII:    strings.TrimSpace(inf.Total.ICMSTot.VII),
			ValorIPI:   strings.TrimSpace(inf.Total.ICMSTot.VIPI),
			ValorPIS:   strings.TrimSpace(inf.Total.ICMSTot.VPIS),
			ValorCOF:   strings.TrimSpace(inf.Total.ICMSTot.VCOFINS),
			ValorProd:  strings.TrimSpace(inf.Total.ICMSTot.VProd),
			ValorFrete: strings.TrimSpace(inf.Total.ICMSTot.VFrete),
			ValorSeg:   strings.TrimSpace(inf.Total.ICMSTot.VSeg),
			ValorDesc:  strings.TrimSpace(inf.Total.ICMSTot.VDesc),
			ValorOutro: strings.TrimSpace(inf.Total.ICMSTot.VOutro),
			ValorTotTrib: strings.TrimSpace(inf.Total.ICMSTot.VTotTrib),
			ValorNF:    strings.TrimSpace(inf.Total.ICMSTot.VNF),
		},
		Transporte: domain.NFEDANFETransporte{
			Modalidade:   modFreteLabelCompleto(strings.TrimSpace(inf.Transp.ModFrete)),
			Transportado: strings.TrimSpace(inf.Transp.Transporta.XNome),
			CNPJCPF:      firstNotEmpty(inf.Transp.Transporta.CNPJ, inf.Transp.Transporta.CPF),
			IE:           strings.TrimSpace(inf.Transp.Transporta.IE),
			Endereco:     strings.TrimSpace(inf.Transp.Transporta.XEnder),
			Municipio:    strings.TrimSpace(inf.Transp.Transporta.XMun),
			Placa:        strings.TrimSpace(inf.Transp.VeicTransp.Placa),
			UF:           firstNotEmpty(inf.Transp.VeicTransp.UF, inf.Transp.Transporta.UF),
			RNTC:         strings.TrimSpace(inf.Transp.VeicTransp.RNTC),
			QuantidadeVol: sumVolumes(inf.Transp.Vol),
			Volumes:      []domain.NFEDANFEVolume{},
		},
		Cobranca: domain.NFEDANFECobranca{
			NumeroFatura:  strings.TrimSpace(inf.Cobr.Fat.NFat),
			ValorOriginal: strings.TrimSpace(inf.Cobr.Fat.VOrig),
			ValorDesconto: strings.TrimSpace(inf.Cobr.Fat.VDesc),
			ValorLiquido:  strings.TrimSpace(inf.Cobr.Fat.VLiq),
			Duplicatas:    []domain.NFEDANFEDuplicata{},
			Pagamentos:    []domain.NFEDANFEPagamento{},
		},
		Adicionais: domain.NFEDANFEAdicionais{
			TpImp:                    tpImpLabel(strings.TrimSpace(inf.Ide.TpImp)),
			InformacoesComplementares: strings.TrimSpace(inf.InfAdic.InfCpl),
			InformacoesFisco:          strings.TrimSpace(inf.InfAdic.InfAdFisco),
		},
	}
	for _, d := range inf.Cobr.Dup {
		view.Cobranca.Duplicatas = append(view.Cobranca.Duplicatas, domain.NFEDANFEDuplicata{
			Numero:     strings.TrimSpace(d.NDup),
			Vencimento: strings.TrimSpace(d.DVenc),
			Valor:      strings.TrimSpace(d.VDup),
		})
	}
	for _, p := range inf.Pag.DetPag {
		view.Cobranca.Pagamentos = append(view.Cobranca.Pagamentos, domain.NFEDANFEPagamento{
			Forma:       firstNotEmpty(strings.TrimSpace(p.XPag), formaPagamentoLabel(strings.TrimSpace(p.TPag))),
			Valor:       strings.TrimSpace(p.VPag),
			CNPJCred:    strings.TrimSpace(p.Card.CNPJ),
			Bandeira:    bandeiraCartaoLabel(strings.TrimSpace(p.Card.TBand)),
			Autorizacao: strings.TrimSpace(p.Card.CAut),
		})
	}
	for _, v := range inf.Transp.Vol {
		view.Transporte.Volumes = append(view.Transporte.Volumes, domain.NFEDANFEVolume{
			Quantidade: strings.TrimSpace(v.QVol),
			Especie:    strings.TrimSpace(v.Esp),
			Marca:      strings.TrimSpace(v.Marca),
			Numero:     strings.TrimSpace(v.NVol),
			PesoLiquido: strings.TrimSpace(v.PesoL),
			PesoBruto:  strings.TrimSpace(v.PesoB),
		})
	}
	view.Itens = make([]domain.NFEDANFEItem, 0, len(inf.Det))
	for _, d := range inf.Det {
		vbcIcms, vicms, picms := extractICMSValores(d.Imposto.ICMS.Inner)
		ipiDet := buildItemIPI(d.Imposto.IPI)
		vipi := strings.TrimSpace(ipiDet.VIPI)
		if vipi == "" {
			vipi = strings.TrimSpace(d.Imposto.IPI.IPITrib.VIPI)
		}
		pisCST := pisCofinsCSTDesc(extractFirstCST(d.Imposto.PIS.Inner), pisCSTLabel)
		cofCST := pisCofinsCSTDesc(extractFirstCST(d.Imposto.COFINS.Inner), cofinsCSTLabel)
		view.Itens = append(view.Itens, domain.NFEDANFEItem{
			Codigo:             strings.TrimSpace(d.Prod.CProd),
			Descricao:          strings.TrimSpace(d.Prod.XProd),
			NCM:                strings.TrimSpace(d.Prod.NCM),
			EXTIPI:             strings.TrimSpace(d.Prod.EXTIPI),
			CFOP:               strings.TrimSpace(d.Prod.CFOP),
			CEAN:               strings.TrimSpace(d.Prod.CEAN),
			Unidade:            strings.TrimSpace(d.Prod.UCom),
			Quantidade:         strings.TrimSpace(d.Prod.QCom),
			ValorUnit:          strings.TrimSpace(d.Prod.VUnCom),
			ValorTotal:         strings.TrimSpace(d.Prod.VProd),
			ValorDesc:          strings.TrimSpace(d.Prod.VDesc),
			ValorFrete:         strings.TrimSpace(d.Prod.VFrete),
			ValorSeg:           strings.TrimSpace(d.Prod.VSeg),
			ValorOutro:         strings.TrimSpace(d.Prod.VOutro),
			IndicadorTotal:     strings.TrimSpace(d.Prod.IndTot),
			IndicadorTotalDesc: indTotLabel(strings.TrimSpace(d.Prod.IndTot)),
			CEANTrib:           strings.TrimSpace(d.Prod.CEANTrib),
			UTrib:              strings.TrimSpace(d.Prod.UTrib),
			QTrib:              strings.TrimSpace(d.Prod.QTrib),
			VUnTrib:            strings.TrimSpace(d.Prod.VUnTrib),
			ValorTotTrib:       strings.TrimSpace(d.Imposto.VTotTrib),
			XPed:               strings.TrimSpace(d.Prod.XPed),
			NItemPed:           strings.TrimSpace(d.Prod.NItemPed),
			NFCI:               strings.TrimSpace(d.Prod.NFCI),
			BaseICMS:           vbcIcms,
			ValorICMS:          vicms,
			ValorIPI:           vipi,
			AliquotaICM:        picms,
			AliquotaIPI:        strings.TrimSpace(d.Imposto.IPI.IPITrib.PIPI),
			ICMS:               parseICMSDetalhe(d.Imposto.ICMS.Inner),
			IPI:                ipiDet,
			PISCST:             pisCST,
			COFINSCST:          cofCST,
		})
	}
	return view, nil
}

func extractDigestValue(xmlStr string) string {
	re := regexp.MustCompile(`(?i)<DigestValue>\s*([^<]+?)\s*</DigestValue>`)
	if m := re.FindStringSubmatch(xmlStr); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func formatDataInclusaoBD(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func formatCodNome(cod, nome string) string {
	cod = strings.TrimSpace(cod)
	nome = strings.TrimSpace(nome)
	if cod != "" && nome != "" {
		return cod + " - " + nome
	}
	if cod != "" {
		return cod
	}
	return nome
}

func joinLinhaEndereco(xLgr, nro, xCpl string) string {
	xLgr = strings.TrimSpace(xLgr)
	nro = strings.TrimSpace(nro)
	xCpl = strings.TrimSpace(xCpl)
	if xLgr == "" && nro == "" && xCpl == "" {
		return ""
	}
	s := xLgr
	if nro != "" {
		if s != "" {
			s += ", "
		}
		s += nro
	}
	if xCpl != "" {
		if s != "" {
			s += " — "
		}
		s += xCpl
	}
	return s
}

func ambienteAutorizacaoLabel(v string) string {
	switch strings.TrimSpace(v) {
	case "1":
		return "1 - Produção"
	case "2":
		return "2 - Homologação"
	default:
		return strings.TrimSpace(v)
	}
}

func tipoEmissaoLabel(v string) string {
	switch strings.TrimSpace(v) {
	case "1":
		return "1 - Emissão normal"
	case "2":
		return "2 - Contingência FS"
	case "3":
		return "3 - Contingência SCAN"
	case "4":
		return "4 - Contingência DPEC"
	case "5":
		return "5 - Contingência FS-DA"
	case "6":
		return "6 - Contingência SVC-AN"
	case "7":
		return "7 - Contingência SVC-RS"
	case "9":
		return "9 - Contingência off-line"
	default:
		return strings.TrimSpace(v)
	}
}

func tpImpLabel(v string) string {
	switch strings.TrimSpace(v) {
	case "1":
		return "1 - Retrato"
	case "2":
		return "2 - Paisagem"
	case "3":
		return "3 - Simplificado"
	case "4":
		return "4 - DANFE NFC-e"
	case "5":
		return "5 - DANFE NFC-e em mensagem eletrônica"
	default:
		if strings.TrimSpace(v) == "" {
			return ""
		}
		return strings.TrimSpace(v)
	}
}

func crtLabel(v string) string {
	switch strings.TrimSpace(v) {
	case "1":
		return "Simples Nacional"
	case "2":
		return "Simples Nacional (excesso de sublimite de receita bruta)"
	case "3":
		return "Regime Normal"
	case "4":
		return "MEI (Simples Nacional)"
	default:
		return ""
	}
}

func crtCodNome(crt string) string {
	crt = strings.TrimSpace(crt)
	if crt == "" {
		return ""
	}
	l := crtLabel(crt)
	if l == "" {
		return crt
	}
	return crt + " - " + l
}

func indIEDestLabel(v string) string {
	switch strings.TrimSpace(v) {
	case "1":
		return "1 - Contribuinte ICMS (informar IE do destinatário)"
	case "2":
		return "2 - Contribuinte isento de Inscrição no cadastro de Contribuintes"
	case "9":
		return "9 - Não Contribuinte"
	default:
		return strings.TrimSpace(v)
	}
}

func modFreteLabelCompleto(v string) string {
	switch strings.TrimSpace(v) {
	case "0":
		return "0 - Contratação do Frete por conta do Remetente (CIF)"
	case "1":
		return "1 - Contratação do Frete por conta do Destinatário (FOB)"
	case "2":
		return "2 - Contratação do Frete por conta de Terceiros"
	case "3":
		return "3 - Transporte Próprio por conta do Remetente"
	case "4":
		return "4 - Transporte Próprio por conta do Destinatário"
	case "9":
		return "9 - Sem Ocorrência de Transporte"
	default:
		return strings.TrimSpace(v)
	}
}

func indTotLabel(v string) string {
	switch strings.TrimSpace(v) {
	case "0":
		return "0 - Valor do item (vProd) não compõe o valor total da NF-e"
	case "1":
		return "1 - Valor do item (vProd) compõe o valor total da NF-e"
	default:
		return strings.TrimSpace(v)
	}
}

func extractFirstCST(inner string) string {
	if strings.TrimSpace(inner) == "" {
		return ""
	}
	re := regexp.MustCompile(`<CST>([0-9]{2})</CST>`)
	if m := re.FindStringSubmatch(inner); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	re2 := regexp.MustCompile(`<CST>([0-9]{3})</CST>`)
	if m := re2.FindStringSubmatch(inner); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func extractICMSValores(inner string) (vbc, vicms, picms string) {
	if strings.TrimSpace(inner) == "" {
		return "", "", ""
	}
	reBC := regexp.MustCompile(`<vBC>([^<]+)</vBC>`)
	reICMS := regexp.MustCompile(`<vICMS>([^<]+)</vICMS>`)
	reP := regexp.MustCompile(`<pICMS>([^<]+)</pICMS>`)
	if m := reBC.FindStringSubmatch(inner); len(m) > 1 {
		vbc = strings.TrimSpace(m[1])
	}
	if m := reICMS.FindStringSubmatch(inner); len(m) > 1 {
		vicms = strings.TrimSpace(m[1])
	}
	if m := reP.FindStringSubmatch(inner); len(m) > 1 {
		picms = strings.TrimSpace(m[1])
	}
	return vbc, vicms, picms
}

func parseICMSDetalhe(inner string) domain.NFEDANFEItemICMS {
	out := domain.NFEDANFEItemICMS{}
	if strings.TrimSpace(inner) == "" {
		return out
	}
	reOrig := regexp.MustCompile(`<orig>([0-9])</orig>`)
	if m := reOrig.FindStringSubmatch(inner); len(m) > 1 {
		out.Origem = origemMercadoriaLabel(m[1])
	}
	reCSOSN := regexp.MustCompile(`<CSOSN>([0-9]{3})</CSOSN>`)
	if m := reCSOSN.FindStringSubmatch(inner); len(m) > 1 {
		out.Tributacao = csosnLabel(m[1])
		return out
	}
	reCST := regexp.MustCompile(`<CST>([0-9]{2,3})</CST>`)
	if m := reCST.FindStringSubmatch(inner); len(m) > 1 {
		out.Tributacao = cstICMSLabel(m[1])
	}
	return out
}

func origemMercadoriaLabel(o string) string {
	switch strings.TrimSpace(o) {
	case "0":
		return "0 - Nacional"
	case "1":
		return "1 - Estrangeira — Importação direta"
	case "2":
		return "2 - Estrangeira — Adquirida no mercado interno"
	case "3":
		return "3 - Nacional com mais de 40% de conteúdo estrangeiro"
	case "4":
		return "4 - Nacional com processos produtivos básicos"
	case "5":
		return "5 - Nacional com menos de 40% de conteúdo estrangeiro"
	case "6":
		return "6 - Estrangeira — Importação direta sem similar nacional"
	case "7":
		return "7 - Estrangeira — Adquirida mercado interno sem similar nacional"
	case "8":
		return "8 - Nacional com conteúdo de importação superior a 70%"
	default:
		return strings.TrimSpace(o)
	}
}

func csosnLabel(c string) string {
	switch strings.TrimSpace(c) {
	case "101":
		return "101 - Tributada pelo Simples Nacional com permissão de crédito"
	case "102":
		return "102 - Tributada pelo Simples Nacional sem permissão de crédito"
	case "103":
		return "103 - Isenção do ICMS no Simples Nacional para faixa de receita bruta"
	case "201":
		return "201 - Tributada pelo Simples Nacional com permissão de crédito e cobrança do ICMS por ST"
	case "202":
		return "202 - Tributada pelo Simples Nacional sem permissão de crédito e com cobrança do ICMS por ST"
	case "203":
		return "203 - Isenção do ICMS no Simples Nacional para faixa de receita bruta e com cobrança do ICMS por ST"
	case "300":
		return "300 - Imune"
	case "400":
		return "400 - Não tributada pelo Simples Nacional"
	case "500":
		return "500 - ICMS cobrado anteriormente por ST ou por antecipação"
	default:
		if strings.TrimSpace(c) == "" {
			return ""
		}
		return c + " - CSOSN (ver manual)"
	}
}

func cstICMSLabel(cst string) string {
	cst = strings.TrimSpace(cst)
	switch cst {
	case "00":
		return "00 - Tributada integralmente"
	case "10":
		return "10 - Tributada e com cobrança do ICMS por ST"
	case "20":
		return "20 - Com redução de base de cálculo"
	case "30":
		return "30 - Isenta ou não tributada e com cobrança do ICMS por ST"
	case "40":
		return "40 - Isenta"
	case "41":
		return "41 - Não tributada"
	case "50":
		return "50 - Suspensão"
	case "51":
		return "51 - Diferimento"
	case "60":
		return "60 - ICMS cobrado anteriormente por ST"
	case "70":
		return "70 - Com redução de base de cálculo e cobrança do ICMS por ST"
	case "90":
		return "90 - Outras"
	default:
		if cst == "" {
			return ""
		}
		return cst + " - CST ICMS (ver manual)"
	}
}

func ipiCSTLabel(cst string) string {
	cst = strings.TrimSpace(cst)
	switch cst {
	case "00":
		return "00 - Entrada com recuperação de crédito"
	case "01":
		return "01 - Entrada tributável com alíquota zero"
	case "02":
		return "02 - Entrada isenta"
	case "03":
		return "03 - Entrada não-tributada"
	case "04":
		return "04 - Entrada imune"
	case "05":
		return "05 - Entrada com suspensão"
	case "49":
		return "49 - Outras entradas"
	case "50":
		return "50 - Saída tributada"
	case "51":
		return "51 - Saída tributável com alíquota zero"
	case "52":
		return "52 - Saída isenta"
	case "53":
		return "53 - Saída não-tributada"
	case "54":
		return "54 - Saída imune"
	case "55":
		return "55 - Saída com suspensão"
	case "99":
		return "99 - Outras saídas"
	default:
		if cst == "" {
			return ""
		}
		return cst + " - CST IPI (ver manual)"
	}
}

func pisCofinsCSTDesc(cod string, descFn func(string) string) string {
	cod = strings.TrimSpace(cod)
	if cod == "" {
		return ""
	}
	d := descFn(cod)
	if d == "" {
		return cod
	}
	return d
}

func pisCSTLabel(c string) string {
	c = strings.TrimSpace(c)
	switch c {
	case "01":
		return "01 - Operação Tributável com Alíquota Básica"
	case "02":
		return "02 - Operação Tributável com Alíquota Diferenciada"
	case "03":
		return "03 - Operação Tributável com Alíquota por Unidade de Medida de Produto"
	case "04":
		return "04 - Operação Tributável monofásica — Revenda a Alíquota Zero"
	case "05":
		return "05 - Operação Tributável por Substituição Tributária"
	case "06":
		return "06 - Operação Tributável a Alíquota Zero"
	case "07":
		return "07 - Operação Isenta da Contribuição"
	case "08":
		return "08 - Operação sem Incidência da Contribuição"
	case "09":
		return "09 - Operação com Suspensão da Contribuição"
	case "49":
		return "49 - Outras Operações de Saída"
	case "50":
		return "50 - Operação com Direito a Crédito — Vinculada Exclusivamente a Receita Tributada no Mercado Interno"
	case "51":
		return "51 - Operação com Direito a Crédito — Vinculada Exclusivamente a Receita Não Tributada no Mercado Interno"
	case "52":
		return "52 - Operação com Direito a Crédito — Vinculada Exclusivamente a Receita de Exportação"
	case "53":
		return "53 - Operação com Direito a Crédito — Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno"
	case "54":
		return "54 - Operação com Direito a Crédito — Vinculada a Receitas Tributadas no Mercado Interno e de Exportação"
	case "55":
		return "55 - Operação com Direito a Crédito — Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação"
	case "56":
		return "56 - Operação com Direito a Crédito — Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno e de Exportação"
	case "60":
		return "60 - Crédito Presumido — Operação de Aquisição Vinculada Exclusivamente a Receita Tributada no Mercado Interno"
	case "61":
		return "61 - Crédito Presumido — Operação de Aquisição Vinculada Exclusivamente a Receita Não-Tributada no Mercado Interno"
	case "62":
		return "62 - Crédito Presumido — Operação de Aquisição Vinculada Exclusivamente a Receita de Exportação"
	case "63":
		return "63 - Crédito Presumido — Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno"
	case "64":
		return "64 - Crédito Presumido — Operação de Aquisição Vinculada a Receitas Tributadas no Mercado Interno e de Exportação"
	case "65":
		return "65 - Crédito Presumido — Operação de Aquisição Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação"
	case "66":
		return "66 - Crédito Presumido — Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno e de Exportação"
	case "67":
		return "67 - Crédito Presumido — Outras Operações"
	case "70":
		return "70 - Operação de Aquisição sem Direito a Crédito"
	case "71":
		return "71 - Operação de Aquisição com Isenção"
	case "72":
		return "72 - Operação de Aquisição com Suspensão"
	case "73":
		return "73 - Operação de Aquisição a Alíquota Zero"
	case "74":
		return "74 - Operação de Aquisição sem Incidência da Contribuição"
	case "75":
		return "75 - Operação de Aquisição por Substituição Tributária"
	case "98":
		return "98 - Outras Operações de Entrada"
	case "99":
		return "99 - Outras Operações"
	default:
		if c == "" {
			return ""
		}
		return c + " - CST PIS/COFINS (ver manual)"
	}
}

func cofinsCSTLabel(c string) string {
	// CST COFINS usa a mesma tabela de descrições resumidas do PIS em muitos casos
	return pisCSTLabel(c)
}

func buildItemIPI(ipi nfeIPIXMLElement) domain.NFEDANFEItemIPI {
	out := domain.NFEDANFEItemIPI{
		ClEnq:    strings.TrimSpace(ipi.ClEnq),
		CEnq:     strings.TrimSpace(ipi.CEnq),
		CSelo:    strings.TrimSpace(ipi.IPITrib.CSelo),
		CNPJProd: strings.TrimSpace(ipi.IPITrib.CNPJProd),
		QSelo:    strings.TrimSpace(ipi.IPITrib.QSelo),
		QUnid:    strings.TrimSpace(ipi.IPITrib.QUnid),
		VUnid:    strings.TrimSpace(ipi.IPITrib.VUnid),
		VBC:      strings.TrimSpace(ipi.IPITrib.VBC),
		PIPI:     strings.TrimSpace(ipi.IPITrib.PIPI),
		VIPI:     strings.TrimSpace(ipi.IPITrib.VIPI),
	}
	cst := strings.TrimSpace(ipi.IPITrib.CST)
	if cst == "" {
		cst = strings.TrimSpace(ipi.IPINT.CST)
	}
	if cst != "" {
		out.CST = ipiCSTLabel(cst)
	}
	return out
}

func firstPayment(det []nfeDetPagXML) string {
	if len(det) == 0 {
		return ""
	}
	if strings.TrimSpace(det[0].XPag) != "" {
		return strings.TrimSpace(det[0].XPag)
	}
	return formaPagamentoLabel(strings.TrimSpace(det[0].TPag))
}

func sumVolumes(vols []struct {
	QVol  string `xml:"qVol"`
	Esp   string `xml:"esp"`
	Marca string `xml:"marca"`
	NVol  string `xml:"nVol"`
	PesoL string `xml:"pesoL"`
	PesoB string `xml:"pesoB"`
}) string {
	total := 0
	has := false
	for _, v := range vols {
		q := strings.TrimSpace(v.QVol)
		if q == "" {
			continue
		}
		has = true
		var n int
		_, _ = fmt.Sscanf(q, "%d", &n)
		total += n
	}
	if !has {
		return ""
	}
	return fmt.Sprintf("%d", total)
}

func destinoOperacaoLabel(v string) string {
	switch v {
	case "1":
		return "1 - Operação Interna"
	case "2":
		return "2 - Operação Interestadual"
	case "3":
		return "3 - Operação com Exterior"
	default:
		return v
	}
}

func consumidorFinalLabel(v string) string {
	switch v {
	case "0":
		return "0 - Normal"
	case "1":
		return "1 - Consumidor final"
	default:
		return v
	}
}

func presencaCompradorLabel(v string) string {
	switch v {
	case "0":
		return "0 - Não se aplica"
	case "1":
		return "1 - Operação presencial"
	case "2":
		return "2 - Operação pela internet"
	case "3":
		return "3 - Operação não presencial (teleatendimento)"
	case "4":
		return "4 - NFC-e com entrega a domicilio"
	case "9":
		return "9 - Operação não presencial (outros)"
	default:
		return v
	}
}

func processoEmissaoLabel(v string) string {
	switch v {
	case "0":
		return "0 - com aplicativo do Contribuinte"
	case "1":
		return "1 - avulsa pelo Fisco"
	case "2":
		return "2 - avulsa pelo Contribuinte com Certificado Digital"
	case "3":
		return "3 - com aplicativo fornecido pelo Fisco"
	default:
		return v
	}
}

func finalidadeLabel(v string) string {
	switch v {
	case "1":
		return "1 - Normal"
	case "2":
		return "2 - complementar"
	case "3":
		return "3 - de Ajuste"
	case "4":
		return "4 - devolução de mercadoria"
	default:
		return v
	}
}

func tipoOperacaoLabel(v string) string {
	switch v {
	case "0":
		return "0 - Entrada"
	case "1":
		return "1 - Saída"
	default:
		return v
	}
}

func formaPagamentoLabel(v string) string {
	switch v {
	case "01":
		return "01 - Dinheiro"
	case "02":
		return "02 - Cheque"
	case "03":
		return "03 - Cartão de Crédito"
	case "04":
		return "04 - Cartão de Débito"
	case "05":
		return "05 - Crédito Loja"
	case "15":
		return "15 - Boleto Bancário"
	case "90":
		return "90 - Sem pagamento"
	case "99":
		return "99 - Outros"
	default:
		return v
	}
}

func bandeiraCartaoLabel(v string) string {
	switch v {
	case "01":
		return "01 - Visa"
	case "02":
		return "02 - Mastercard"
	case "03":
		return "03 - American Express"
	case "04":
		return "04 - Sorocred"
	case "99":
		return "99 - Outros"
	default:
		return v
	}
}

func eventDescriptionByCStat(cStat, motivo string) string {
	switch cStat {
	case "100":
		return "Autorização de Uso (Cód.: 110100)"
	case "150":
		return "Autorização de Uso Fora de Prazo (Cód.: 110100)"
	case "301":
		return "Denegação de Uso - Situação do emitente (Cód.: 110101)"
	case "302":
		return "Denegação de Uso - Situação do destinatário (Cód.: 110101)"
	default:
		if cStat == "" && motivo == "" {
			return ""
		}
		return strings.TrimSpace(cStat + " - " + motivo)
	}
}

func statusAtualByCStat(cStat, motivo string) string {
	switch cStat {
	case "100", "150":
		return "AUTORIZADA"
	case "301", "302":
		return "DENEGADA"
	default:
		if strings.TrimSpace(motivo) != "" {
			return strings.TrimSpace(motivo)
		}
		return cStat
	}
}

func parseNFEXMLAnyEnvelope(xmlStr string) (nfeInfXML, nfeInfProtXML, error) {
	var inf nfeInfXML
	var prot nfeInfProtXML

	dec := xml.NewDecoder(bytes.NewReader([]byte(xmlStr)))
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		switch start.Name.Local {
		case "infNFe":
			if inf.ID == "" {
				if err := dec.DecodeElement(&inf, &start); err != nil {
					return nfeInfXML{}, nfeInfProtXML{}, fmt.Errorf("xml nf-e invalido para danfe view (infNFe): %w", err)
				}
			}
		case "infProt":
			if prot.NProt == "" && prot.XMotivo == "" {
				if err := dec.DecodeElement(&prot, &start); err != nil {
					return nfeInfXML{}, nfeInfProtXML{}, fmt.Errorf("xml nf-e invalido para danfe view (infProt): %w", err)
				}
			}
		}
	}
	if strings.TrimSpace(inf.ID) == "" && strings.TrimSpace(inf.Ide.NNF) == "" {
		return nfeInfXML{}, nfeInfProtXML{}, fmt.Errorf("xml nf-e invalido para danfe view (infNFe ausente no envelope)")
	}
	return inf, prot, nil
}

func firstNotEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return strings.TrimSpace(a)
	}
	return strings.TrimSpace(b)
}
