package nfeprovider

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

const (
	scURLProducao    = "https://satnfe.sef.sc.gov.br/ws/distribuicao/nfedownloadV2.asmx"
	scURLHomologacao = "https://hom.satnfe.sef.sc.gov.br/ws/distribuicao/nfedownloadV2.asmx"
)

type SCProvider struct {
	url      string
	verAplic string
	client   *http.Client
}

func NewSCProvider(homologacao bool, verAplic string) *SCProvider {
	url := scURLProducao
	if homologacao {
		url = scURLHomologacao
	}
	verAplic = strings.TrimSpace(verAplic)
	if verAplic == "" {
		verAplic = "vecontab-ef920"
	}
	return &SCProvider{url: url, verAplic: verAplic}
}

func (p *SCProvider) ConfigurarCertificado(pfx []byte, password string) error {
	priv, cert, err := pkcs12.Decode(pfx, password)
	if err != nil {
		return fmt.Errorf("pkcs12 decode: %w", err)
	}
	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  priv,
		Leaf:        cert,
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{tlsCert},
		},
	}
	p.client = &http.Client{
		Transport: tr,
		Timeout:   100 * time.Second,
	}
	return nil
}

type soapEnvelope struct {
	XMLName xml.Name `xml:"soap12:Envelope"`
	Soap12  string   `xml:"xmlns:soap12,attr"`
	XSI     string   `xml:"xmlns:xsi,attr,omitempty"`
	XSD     string   `xml:"xmlns:xsd,attr,omitempty"`
	Body    soapBody `xml:"soap12:Body"`
}

type soapBody struct {
	NFEDownloadContab nfeDownloadContabReq `xml:"nfeDownloadContab"`
}

type nfeDownloadContabReq struct {
	XMLName   xml.Name     `xml:"nfeDownloadContab"`
	XMLNS     string       `xml:"xmlns,attr"`
	DistNFeSC distNFeSCReq `xml:"distNFeSC"`
}

type distNFeSCReq struct {
	XMLName  xml.Name   `xml:"distNFeSC"`
	XMLNS    string     `xml:"xmlns,attr"`
	Versao   string     `xml:"versao,attr"`
	TpAmb    int        `xml:"tpAmb"`
	VerAplic string     `xml:"verAplic"`
	CUF      int        `xml:"cUF"`
	CNPJ     string     `xml:"CNPJ,omitempty"`
	CPF      string     `xml:"CPF,omitempty"`
	SolRel   *solRelReq `xml:"solRel,omitempty"`
}

type solRelReq struct {
	IndXML   int    `xml:"indXML"`
	IndAtor  int    `xml:"indAtor"`
	UltNuNSU string `xml:"ultNuNSU"`
}

type retDistNFeSC struct {
	XMLName      xml.Name `xml:"retDistNFeSC"`
	Versao       string   `xml:"versao,attr"`
	CStat        int      `xml:"cStat"`
	XMotivo      string   `xml:"xMotivo"`
	UltNuNSURet  string   `xml:"ultNuNSURet"`
	QtDFeRet     int      `xml:"qtDfeRet"`
	LoteDistComp string   `xml:"loteDistComp"`
}

type loteDistNFeSC struct {
	XMLName xml.Name        `xml:"loteDistNFeSC"`
	Versao  string          `xml:"versao,attr"`
	Itens   []distNFeSCItem `xml:"distNFeSC"`
}

type distNFeSCItem struct {
	NSU      string `xml:"NSU,attr"`
	ChAcesso string `xml:"chAcesso,attr"`
	InnerXML string `xml:",innerxml"`
}

func (p *SCProvider) SincronizarDocumentos(ctx context.Context, cnpj string, ultNSU string) (ResultadoSincronizacao, error) {
	if p.client == nil {
		return ResultadoSincronizacao{}, fmt.Errorf("certificado nao configurado para provider SC")
	}
	cnpj = onlyDigits(cnpj)
	if len(cnpj) != 14 && len(cnpj) != 11 {
		return ResultadoSincronizacao{}, fmt.Errorf("cnpj/cpf invalido para sincronizacao SC")
	}
	cursor := normalizeNSU(ultNSU)
	all := make([]DocumentoFiscal, 0, 64)
	lastCStat := 0
	lastMotivo := ""

	for i := 0; i < 200; i++ {
		ret, err := p.callDist(ctx, cnpj, cursor)
		if err != nil {
			return ResultadoSincronizacao{}, err
		}
		lastCStat = ret.CStat
		lastMotivo = strings.TrimSpace(ret.XMotivo)
		if higherNSU(cursor, ret.UltNuNSURet) {
			cursor = normalizeNSU(ret.UltNuNSURet)
		}

		switch ret.CStat {
		case 117:
			return ResultadoSincronizacao{Documentos: all, NovoMaxNSU: cursor, CStat: ret.CStat, XMotivo: ret.XMotivo}, nil
		case 110:
			return ResultadoSincronizacao{Documentos: all, NovoMaxNSU: cursor, CStat: ret.CStat, XMotivo: ret.XMotivo}, nil
		case 118, 138:
			decoded, err := decodeLote(ret.LoteDistComp)
			if err != nil {
				return ResultadoSincronizacao{}, err
			}
			items, err := parseLote(decoded)
			if err != nil {
				return ResultadoSincronizacao{}, err
			}
			all = append(all, items...)

			// Continua imediatamente quando vier lote cheio (50 docs) e houve avanço de NSU.
			if ret.QtDFeRet >= 50 {
				continue
			}
			return ResultadoSincronizacao{Documentos: all, NovoMaxNSU: cursor, CStat: ret.CStat, XMotivo: ret.XMotivo}, nil
		default:
			return ResultadoSincronizacao{}, fmt.Errorf("provider SC rejeitou requisicao: cStat=%d xMotivo=%s", ret.CStat, strings.TrimSpace(ret.XMotivo))
		}
	}
	return ResultadoSincronizacao{Documentos: all, NovoMaxNSU: cursor, CStat: lastCStat, XMotivo: lastMotivo}, nil
}

func (p *SCProvider) callDist(ctx context.Context, cnpj, ultNSU string) (retDistNFeSC, error) {
	env := soapEnvelope{
		Soap12: "http://www.w3.org/2003/05/soap-envelope",
		XSI:    "http://www.w3.org/2001/XMLSchema-instance",
		XSD:    "http://www.w3.org/2001/XMLSchema",
		Body: soapBody{
			NFEDownloadContab: nfeDownloadContabReq{
				XMLNS: "http://www.satnfe.sef.sc.gov.br/ws/distribuicao-v2",
				DistNFeSC: distNFeSCReq{
					XMLNS:    "http://www.satnfe.sef.sc.gov.br/ws/distribuicao-v2",
					Versao:   "2.02",
					TpAmb:    1,
					VerAplic: p.verAplic,
					CUF:      42,
					CNPJ:     cnpj,
					SolRel: &solRelReq{
						IndXML:   1,
						IndAtor:  3,
						UltNuNSU: normalizeNSU(ultNSU),
					},
				},
			},
		},
	}
	body, err := xml.Marshal(env)
	if err != nil {
		return retDistNFeSC{}, fmt.Errorf("marshal soap sc: %w", err)
	}
	payload := append([]byte(xml.Header), body...)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.url, bytes.NewReader(payload))
	if err != nil {
		return retDistNFeSC{}, err
	}
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("Accept", "application/soap+xml, text/xml, application/xml")

	resp, err := p.client.Do(req)
	if err != nil {
		return retDistNFeSC{}, fmt.Errorf("request provider SC: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 30<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return retDistNFeSC{}, fmt.Errorf("provider SC status %d: %s", resp.StatusCode, clip(raw, 1600))
	}
	retXML, err := extractXMLNode(raw, "retDistNFeSC")
	if err != nil {
		return retDistNFeSC{}, err
	}
	var ret retDistNFeSC
	if err := xml.Unmarshal(retXML, &ret); err != nil {
		return retDistNFeSC{}, fmt.Errorf("parse retDistNFeSC: %w", err)
	}
	return ret, nil
}

func decodeLote(loteB64 string) ([]byte, error) {
	if strings.TrimSpace(loteB64) == "" {
		return nil, nil
	}
	clean := strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' || r == ' ' {
			return -1
		}
		return r
	}, loteB64)
	z, err := base64.StdEncoding.DecodeString(clean)
	if err != nil {
		return nil, fmt.Errorf("base64 loteDistComp: %w", err)
	}
	gz, err := gzip.NewReader(bytes.NewReader(z))
	if err != nil {
		return nil, fmt.Errorf("gzip loteDistComp: %w", err)
	}
	defer gz.Close()
	out, err := io.ReadAll(io.LimitReader(gz, 60<<20))
	if err != nil {
		return nil, fmt.Errorf("read gzip loteDistComp: %w", err)
	}
	return out, nil
}

func parseLote(raw []byte) ([]DocumentoFiscal, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var lote loteDistNFeSC
	if err := xml.Unmarshal(raw, &lote); err != nil {
		return nil, fmt.Errorf("parse loteDistNFeSC: %w", err)
	}
	out := make([]DocumentoFiscal, 0, len(lote.Itens))
	now := time.Now().UTC()
	for _, it := range lote.Itens {
		tipo := guessTipo(strings.TrimSpace(it.InnerXML))
		xmlRaw := strings.TrimSpace(it.InnerXML)
		out = append(out, DocumentoFiscal{
			NSU:         strings.TrimSpace(it.NSU),
			ChaveAcesso: strings.TrimSpace(it.ChAcesso),
			Tipo:        tipo,
			XML:         xmlRaw,
			RecebidoEm:  now,
		})
	}
	return out, nil
}

func guessTipo(inner string) string {
	s := strings.ToLower(inner)
	if strings.Contains(s, "<nfeproc") {
		return "nfeProc"
	}
	if strings.Contains(s, "<proceventonfe") {
		return "procEventoNFe"
	}
	if inner == "" {
		return "chave"
	}
	return "documento"
}

func extractXMLNode(raw []byte, tag string) ([]byte, error) {
	s := string(raw)
	start := strings.Index(strings.ToLower(s), "<"+strings.ToLower(tag))
	if start < 0 {
		return nil, fmt.Errorf("tag %s nao encontrada no SOAP", tag)
	}
	endTag := "</" + strings.ToLower(tag) + ">"
	end := strings.Index(strings.ToLower(s[start:]), endTag)
	if end < 0 {
		return nil, fmt.Errorf("fecho da tag %s nao encontrado no SOAP", tag)
	}
	end += start + len(endTag)
	return []byte(s[start:end]), nil
}

func normalizeNSU(v string) string {
	v = onlyDigits(v)
	if v == "" {
		return "0"
	}
	return strings.TrimLeft(v, "0")
}

func higherNSU(current, candidate string) bool {
	a := normalizeNSU(current)
	b := normalizeNSU(candidate)
	if b == "" || b == "0" {
		return false
	}
	if a == "" {
		a = "0"
	}
	ai, errA := strconv.ParseUint(a, 10, 64)
	bi, errB := strconv.ParseUint(b, 10, 64)
	if errA != nil || errB != nil {
		return len(b) > len(a) || (len(b) == len(a) && b > a)
	}
	return bi > ai
}

func onlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func clip(b []byte, max int) string {
	s := strings.TrimSpace(string(b))
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
