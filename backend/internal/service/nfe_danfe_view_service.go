package service

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
)

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
		CNPJ  string `xml:"CNPJ"`
		CPF   string `xml:"CPF"`
		IE    string `xml:"IE"`
		Ender struct {
			XLgr    string `xml:"xLgr"`
			Nro     string `xml:"nro"`
			XBairro string `xml:"xBairro"`
			XMun    string `xml:"xMun"`
			UF      string `xml:"UF"`
			CEP     string `xml:"CEP"`
		} `xml:"enderEmit"`
	} `xml:"emit"`
	Dest struct {
		XNome string `xml:"xNome"`
		CNPJ  string `xml:"CNPJ"`
		CPF   string `xml:"CPF"`
		IE    string `xml:"IE"`
		IndIEDest string `xml:"indIEDest"`
		Ender struct {
			XLgr    string `xml:"xLgr"`
			Nro     string `xml:"nro"`
			XBairro string `xml:"xBairro"`
			XMun    string `xml:"xMun"`
			UF      string `xml:"UF"`
			CEP     string `xml:"CEP"`
		} `xml:"enderDest"`
	} `xml:"dest"`
	Det []struct {
		Prod struct {
			CProd  string `xml:"cProd"`
			XProd  string `xml:"xProd"`
			NCM    string `xml:"NCM"`
			EXTIPI string `xml:"EXTIPI"`
			CFOP   string `xml:"CFOP"`
			CEAN   string `xml:"cEAN"`
			UCom   string `xml:"uCom"`
			QCom   string `xml:"qCom"`
			VUnCom string `xml:"vUnCom"`
			VProd  string `xml:"vProd"`
			VDesc  string `xml:"vDesc"`
			VFrete string `xml:"vFrete"`
			VSeg   string `xml:"vSeg"`
			VOutro string `xml:"vOutro"`
			IndTot string `xml:"indTot"`
			CEANTrib string `xml:"cEANTrib"`
			UTrib  string `xml:"uTrib"`
			QTrib  string `xml:"qTrib"`
			VUnTrib string `xml:"vUnTrib"`
		} `xml:"prod"`
		Imposto struct {
			VTotTrib string `xml:"vTotTrib"`
			ICMS struct {
				Any struct {
					VBC   string `xml:"vBC"`
					VICMS string `xml:"vICMS"`
					PICMS string `xml:"pICMS"`
				} `xml:",any"`
			} `xml:"ICMS"`
			IPI struct {
				IPITrib struct {
					VIPI string `xml:"vIPI"`
					PIPI string `xml:"pIPI"`
				} `xml:"IPITrib"`
			} `xml:"IPI"`
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
			Ambiente:          strings.TrimSpace(inf.Ide.TpAmb),
			Situacao:          statusAtualByCStat(strings.TrimSpace(prot.CStat), strings.TrimSpace(prot.XMotivo)),
			NaturezaOp:        strings.TrimSpace(inf.Ide.NatOp),
			TipoOperacao:      tipoOperacaoLabel(strings.TrimSpace(inf.Ide.TpNF)),
			DestinoOperacao:   destinoOperacaoLabel(strings.TrimSpace(inf.Ide.IDest)),
			ConsumidorFinal:   consumidorFinalLabel(strings.TrimSpace(inf.Ide.IndFinal)),
			PresencaComprador: presencaCompradorLabel(strings.TrimSpace(inf.Ide.IndPres)),
			ProcessoEmissao:   processoEmissaoLabel(strings.TrimSpace(inf.Ide.ProcEmi)),
			VersaoProcesso:    strings.TrimSpace(inf.Ide.VerProc),
			Finalidade:        finalidadeLabel(strings.TrimSpace(inf.Ide.FinNFe)),
			FormaPagamento:    firstPayment(inf.Pag.DetPag),
		},
		Emitente: domain.NFEDANFEPessoa{
			Nome:       strings.TrimSpace(inf.Emit.XNome),
			CNPJCPF:    firstNotEmpty(inf.Emit.CNPJ, inf.Emit.CPF),
			IE:         strings.TrimSpace(inf.Emit.IE),
			Logradouro: strings.TrimSpace(inf.Emit.Ender.XLgr),
			Numero:     strings.TrimSpace(inf.Emit.Ender.Nro),
			Bairro:     strings.TrimSpace(inf.Emit.Ender.XBairro),
			Municipio:  strings.TrimSpace(inf.Emit.Ender.XMun),
			UF:         strings.TrimSpace(inf.Emit.Ender.UF),
			CEP:        strings.TrimSpace(inf.Emit.Ender.CEP),
		},
		Destinatario: domain.NFEDANFEPessoa{
			Nome:       strings.TrimSpace(inf.Dest.XNome),
			CNPJCPF:    firstNotEmpty(inf.Dest.CNPJ, inf.Dest.CPF),
			IE:         strings.TrimSpace(inf.Dest.IE),
			IndicadorIEDest: strings.TrimSpace(inf.Dest.IndIEDest),
			Logradouro: strings.TrimSpace(inf.Dest.Ender.XLgr),
			Numero:     strings.TrimSpace(inf.Dest.Ender.Nro),
			Bairro:     strings.TrimSpace(inf.Dest.Ender.XBairro),
			Municipio:  strings.TrimSpace(inf.Dest.Ender.XMun),
			UF:         strings.TrimSpace(inf.Dest.Ender.UF),
			CEP:        strings.TrimSpace(inf.Dest.Ender.CEP),
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
			Modalidade:   strings.TrimSpace(inf.Transp.ModFrete),
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
		view.Itens = append(view.Itens, domain.NFEDANFEItem{
			Codigo:      strings.TrimSpace(d.Prod.CProd),
			Descricao:   strings.TrimSpace(d.Prod.XProd),
			NCM:         strings.TrimSpace(d.Prod.NCM),
			EXTIPI:      strings.TrimSpace(d.Prod.EXTIPI),
			CFOP:        strings.TrimSpace(d.Prod.CFOP),
			CEAN:        strings.TrimSpace(d.Prod.CEAN),
			Unidade:     strings.TrimSpace(d.Prod.UCom),
			Quantidade:  strings.TrimSpace(d.Prod.QCom),
			ValorUnit:   strings.TrimSpace(d.Prod.VUnCom),
			ValorTotal:  strings.TrimSpace(d.Prod.VProd),
			ValorDesc:   strings.TrimSpace(d.Prod.VDesc),
			ValorFrete:  strings.TrimSpace(d.Prod.VFrete),
			ValorSeg:    strings.TrimSpace(d.Prod.VSeg),
			ValorOutro:  strings.TrimSpace(d.Prod.VOutro),
			IndicadorTotal: strings.TrimSpace(d.Prod.IndTot),
			CEANTrib:    strings.TrimSpace(d.Prod.CEANTrib),
			UTrib:       strings.TrimSpace(d.Prod.UTrib),
			QTrib:       strings.TrimSpace(d.Prod.QTrib),
			VUnTrib:     strings.TrimSpace(d.Prod.VUnTrib),
			ValorTotTrib: strings.TrimSpace(d.Imposto.VTotTrib),
			BaseICMS:    strings.TrimSpace(d.Imposto.ICMS.Any.VBC),
			ValorICMS:   strings.TrimSpace(d.Imposto.ICMS.Any.VICMS),
			ValorIPI:    strings.TrimSpace(d.Imposto.IPI.IPITrib.VIPI),
			AliquotaICM: strings.TrimSpace(d.Imposto.ICMS.Any.PICMS),
			AliquotaIPI: strings.TrimSpace(d.Imposto.IPI.IPITrib.PIPI),
		})
	}
	return view, nil
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
