package nfeprovider

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
	"golang.org/x/crypto/pkcs12"
)

const (
	svrsRecepcaoEventoProducao    = "https://nfe.svrs.rs.gov.br/ws/recepcaoevento/recepcaoevento4.asmx"
	svrsRecepcaoEventoHomologacao = "https://nfe-homologacao.svrs.rs.gov.br/ws/recepcaoevento/recepcaoevento4.asmx"
	soapActionRecepcaoEvento4     = "http://www.portalfiscal.inf.br/nfe/wsdl/NFeRecepcaoEvento4/nfeRecepcaoEvento"
	nfeNamespace                  = "http://www.portalfiscal.inf.br/nfe"
	// COrgao SVRS (NF-e autorizada via ambiente virtual RS).
	cOrgaoSVRS = 91
)

// NewHTTPClientWithPFX cliente HTTPS com certificado de cliente (mTLS SVRS), para Recepção de Evento.
func NewHTTPClientWithPFX(pfx []byte, password string, timeout time.Duration) (*http.Client, tls.Certificate, error) {
	priv, cert, err := pkcs12.Decode(pfx, password)
	if err != nil {
		return nil, tls.Certificate{}, fmt.Errorf("pkcs12 decode: %w", err)
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
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	return &http.Client{Transport: tr, Timeout: timeout}, tlsCert, nil
}

// ManifestacaoDestinatarioRequest parâmetros para envio de manifestação (tpEvento 2102xx) ao SVRS.
type ManifestacaoDestinatarioRequest struct {
	ChaveNFe    string
	TpEvento    string
	CNPJCPFDest string
	TpAmb       int // 1 produção, 2 homologação
	XJust       string
	NSeqEvento  int
}

// ManifestacaoDestinatarioResult resultado após chamada ao webservice.
type ManifestacaoDestinatarioResult struct {
	CStatLote       int
	XMotivoLote     string
	CStatEvento     int
	XMotivoEvento   string
	NProt           string
	RetEnvEventoXML string
}

// EnvManifestacaoDestinatarioSVRS monta envEvento, assina infEvento (RSA-SHA256, exc-c14n), compacta e envia ao RecepcaoEvento 4.00.
func EnvManifestacaoDestinatarioSVRS(ctx context.Context, client *http.Client, homologacao bool, tlsCert tls.Certificate, req ManifestacaoDestinatarioRequest) (ManifestacaoDestinatarioResult, error) {
	if client == nil {
		return ManifestacaoDestinatarioResult{}, fmt.Errorf("http client nulo")
	}
	req.ChaveNFe = onlyDigits(req.ChaveNFe)
	req.CNPJCPFDest = onlyDigits(req.CNPJCPFDest)
	req.TpEvento = strings.TrimSpace(req.TpEvento)
	if len(req.ChaveNFe) != 44 {
		return ManifestacaoDestinatarioResult{}, fmt.Errorf("chave_nfe deve ter 44 digitos")
	}
	if req.TpAmb != 1 && req.TpAmb != 2 {
		return ManifestacaoDestinatarioResult{}, fmt.Errorf("tpAmb invalido (use 1 ou 2)")
	}
	seq := req.NSeqEvento
	if seq < 1 {
		seq = 1
	}
	if seq > 99 {
		return ManifestacaoDestinatarioResult{}, fmt.Errorf("nSeqEvento maximo 99")
	}
	if err := validateTpEventoManifest(req.TpEvento, req.XJust); err != nil {
		return ManifestacaoDestinatarioResult{}, err
	}
	if len(req.CNPJCPFDest) != 11 && len(req.CNPJCPFDest) != 14 {
		return ManifestacaoDestinatarioResult{}, fmt.Errorf("cnpj/cpf destinatario invalido")
	}

	envXML, err := buildSignedEnvEvento(tlsCert, req, seq)
	if err != nil {
		return ManifestacaoDestinatarioResult{}, err
	}
	b64, err := gzipBase64(envXML)
	if err != nil {
		return ManifestacaoDestinatarioResult{}, err
	}
	url := svrsRecepcaoEventoProducao
	if homologacao {
		url = svrsRecepcaoEventoHomologacao
	}
	soap := buildRecepcaoEventoSOAP(b64)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader([]byte(soap)))
	if err != nil {
		return ManifestacaoDestinatarioResult{}, err
	}
	httpReq.Header.Set("Content-Type", "application/soap+xml; charset=utf-8; action=\""+soapActionRecepcaoEvento4+"\"")
	httpReq.Header.Set("Accept", "application/soap+xml, text/xml, application/xml")

	resp, err := client.Do(httpReq)
	if err != nil {
		return ManifestacaoDestinatarioResult{}, fmt.Errorf("recepcao evento svrs: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 20<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ManifestacaoDestinatarioResult{}, fmt.Errorf("recepcao evento svrs status %d: %s", resp.StatusCode, clip(raw, 1200))
	}
	retXML, err := decodeRecepcaoEventoResultBody(raw)
	if err != nil {
		return ManifestacaoDestinatarioResult{}, err
	}
	return parseRetEnvEvento(retXML)
}

func validateTpEventoManifest(tp, xjust string) error {
	switch tp {
	case "210200", "210210":
		return nil
	case "210220", "210240":
		t := strings.TrimSpace(xjust)
		if len(t) < 15 {
			return fmt.Errorf("tpEvento %s exige x_just com no minimo 15 caracteres", tp)
		}
		return nil
	default:
		return fmt.Errorf("tpEvento invalido (use 210200, 210210, 210220 ou 210240)")
	}
}

func descEventoManifest(tp string) (string, error) {
	switch tp {
	case "210200":
		return "Confirmacao da Operacao", nil
	case "210210":
		return "Ciencia da Operacao", nil
	case "210220":
		return "Desconhecimento da Operacao", nil
	case "210240":
		return "Operacao nao Realizada", nil
	default:
		return "", fmt.Errorf("tpEvento desconhecido")
	}
}

func buildSignedEnvEvento(tlsCert tls.Certificate, req ManifestacaoDestinatarioRequest, seq int) ([]byte, error) {
	desc, err := descEventoManifest(req.TpEvento)
	if err != nil {
		return nil, err
	}
	idVal := "ID" + req.TpEvento + req.ChaveNFe + fmt.Sprintf("%02d", seq)

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		loc = time.UTC
	}
	dh := time.Now().In(loc).Format("2006-01-02T15:04:05-07:00")

	doc := etree.NewDocument()
	env := doc.CreateElement("envEvento")
	env.CreateAttr("xmlns", nfeNamespace)
	env.CreateAttr("versao", "1.00")
	idLote := env.CreateElement("idLote")
	idLote.SetText("1")

	evento := env.CreateElement("evento")
	evento.CreateAttr("versao", "1.00")
	inf := evento.CreateElement("infEvento")
	inf.CreateAttr("Id", idVal)
	inf.CreateAttr("xmlns", nfeNamespace)

	inf.CreateElement("cOrgao").SetText(strconv.Itoa(cOrgaoSVRS))
	inf.CreateElement("tpAmb").SetText(strconv.Itoa(req.TpAmb))
	if len(req.CNPJCPFDest) == 14 {
		inf.CreateElement("CNPJ").SetText(req.CNPJCPFDest)
	} else {
		inf.CreateElement("CPF").SetText(req.CNPJCPFDest)
	}
	inf.CreateElement("chNFe").SetText(req.ChaveNFe)
	inf.CreateElement("dhEvento").SetText(dh)
	inf.CreateElement("tpEvento").SetText(req.TpEvento)
	inf.CreateElement("nSeqEvento").SetText(strconv.Itoa(seq))
	inf.CreateElement("verEvento").SetText("1.00")

	det := inf.CreateElement("detEvento")
	det.CreateAttr("versao", "1.00")
	det.CreateElement("descEvento").SetText(desc)
	if req.TpEvento == "210220" || req.TpEvento == "210240" {
		det.CreateElement("xJust").SetText(strings.TrimSpace(req.XJust))
	}

	keyStore := dsig.TLSCertKeyStore(tlsCert)
	ctx := dsig.NewDefaultSigningContext(keyStore)
	ctx.IdAttribute = "Id"
	ctx.Canonicalizer = dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("")

	sig, err := ctx.ConstructSignature(inf, true)
	if err != nil {
		return nil, fmt.Errorf("assinar envEvento: %w", err)
	}
	evento.AddChild(sig)

	return doc.WriteToBytes()
}

func gzipBase64(xmlBytes []byte) (string, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(xmlBytes); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func buildRecepcaoEventoSOAP(nfeDadosMsgB64 string) string {
	esc := html.EscapeString(nfeDadosMsgB64)
	return `<?xml version="1.0" encoding="UTF-8"?>` +
		`<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">` +
		`<soap12:Body>` +
		`<nfeRecepcaoEvento xmlns="http://www.portalfiscal.inf.br/nfe/wsdl/NFeRecepcaoEvento4">` +
		`<nfeDadosMsg>` + esc + `</nfeDadosMsg>` +
		`</nfeRecepcaoEvento>` +
		`</soap12:Body>` +
		`</soap12:Envelope>`
}

func decodeRecepcaoEventoResultBody(soap []byte) ([]byte, error) {
	s := string(soap)
	// Corpo útil costuma estar em <nfeRecepcaoEventoResult> (várias capitalizações).
	low := strings.ToLower(s)
	idx := strings.Index(low, "<nfeRecepcaoEventoResult")
	if idx < 0 {
		idx = strings.Index(low, "<nferecepcaoeventoresult")
	}
	if idx < 0 {
		return nil, fmt.Errorf("tag nfeRecepcaoEventoResult ausente na resposta SVRS")
	}
	start := strings.Index(s[idx:], ">")
	if start < 0 {
		return nil, fmt.Errorf("resultado SOAP malformado")
	}
	start += idx + 1
	endTag := strings.Index(strings.ToLower(s[start:]), "</nfeRecepcaoEventoResult")
	if endTag < 0 {
		endTag = strings.Index(strings.ToLower(s[start:]), "</nferecepcaoeventoresult")
	}
	if endTag < 0 {
		return nil, fmt.Errorf("fechamento nfeRecepcaoEventoResult ausente")
	}
	inner := strings.TrimSpace(s[start : start+endTag])
	inner = html.UnescapeString(inner)
	if inner == "" {
		return nil, fmt.Errorf("nfeRecepcaoEventoResult vazio")
	}
	// XML direto
	if strings.HasPrefix(inner, "<") {
		return []byte(inner), nil
	}
	// Base64 (+ opcional gzip)
	raw, err := base64.StdEncoding.DecodeString(strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' || r == ' ' {
			return -1
		}
		return r
	}, inner))
	if err != nil {
		return nil, fmt.Errorf("base64 resultado recepcao: %w", err)
	}
	if len(raw) >= 2 && raw[0] == 0x1f && raw[1] == 0x8b {
		gz, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("gzip resultado recepcao: %w", err)
		}
		out, err := io.ReadAll(io.LimitReader(gz, 10<<20))
		_ = gz.Close()
		if err != nil {
			return nil, err
		}
		return out, nil
	}
	return raw, nil
}

type retEnvEventoXML struct {
	XMLName   xml.Name        `xml:"retEnvEvento"`
	CStat     int             `xml:"cStat"`
	XMotivo   string          `xml:"xMotivo"`
	RetEvento []retEventoItem `xml:"retEvento"`
}

type retEventoItem struct {
	InfEvento infEventoRetXML `xml:"infEvento"`
}

type infEventoRetXML struct {
	CStat   int    `xml:"cStat"`
	XMotivo string `xml:"xMotivo"`
	NProt   string `xml:"nProt"`
}

func parseRetEnvEvento(raw []byte) (ManifestacaoDestinatarioResult, error) {
	var root retEnvEventoXML
	if err := xml.Unmarshal(raw, &root); err != nil {
		return ManifestacaoDestinatarioResult{}, fmt.Errorf("parse retEnvEvento: %w", err)
	}
	out := ManifestacaoDestinatarioResult{
		CStatLote:       root.CStat,
		XMotivoLote:     strings.TrimSpace(root.XMotivo),
		RetEnvEventoXML: string(raw),
	}
	if len(root.RetEvento) > 0 {
		ev := root.RetEvento[0].InfEvento
		out.CStatEvento = ev.CStat
		out.XMotivoEvento = strings.TrimSpace(ev.XMotivo)
		out.NProt = strings.TrimSpace(ev.NProt)
	}
	return out, nil
}
