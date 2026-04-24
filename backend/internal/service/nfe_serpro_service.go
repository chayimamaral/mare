package service

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/nfeprovider"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
	"github.com/jackc/pgx/v5"
)

const (
	consultaNFEBaseTrial    = "https://gateway.apiserpro.serpro.gov.br/consulta-nfe-df-trial/api"
	consultaNFEBaseProducao = "https://gateway.apiserpro.serpro.gov.br/consulta-nfe-df/api"
)

var (
	// Saxon reports errors with a line starting like "Error at" / "Type error at" (not "Warning at").
	saxonStderrErrorKind = regexp.MustCompile(`(?i)(\A|\n)\s*(Error at |Type error at |Static error at )`)
	// XPath/XSLT static error codes (avoid loose substring checks like "xtse" inside unrelated text).
	saxonStderrErrCode = regexp.MustCompile(`(?i)\b(SXXP|XTSE|XTDE|XPST|XPTY)[0-9]{3,}\b`)
)

type NFESerproService struct {
	repo       *repository.NFESerproRepository
	serproAuth *SerproService
	certSvc    *CertificadoService
}

func NewNFESerproService(repo *repository.NFESerproRepository, serproAuth *SerproService, certSvc *CertificadoService) *NFESerproService {
	return &NFESerproService{repo: repo, serproAuth: serproAuth, certSvc: certSvc}
}

var nfeChaveRegex = regexp.MustCompile(`^\d{44}$`)

func normalizeStaticBearer(raw string) string {
	t := strings.TrimSpace(raw)
	if t == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(t), "bearer ") {
		t = strings.TrimSpace(t[7:])
	}
	return t
}

func validateNFEChave(chave string) error {
	c := strings.TrimSpace(chave)
	if !nfeChaveRegex.MatchString(c) {
		return fmt.Errorf("chave_nfe deve ter 44 digitos numericos")
	}
	return nil
}

func normalizeNFEAmbienteIntegra(raw string) string {
	a := strings.ToLower(strings.TrimSpace(raw))
	if a == string(IntegraAmbienteProducao) {
		return string(IntegraAmbienteProducao)
	}
	return string(IntegraAmbienteTrial)
}

// resolveConsultaNFEAPIBase: se SERPRO_NFE_API_BASE_URL estiver definida, ela prevalece (ambiente no body e ignorado).
// Caso contrario, usa o gateway SERPRO conforme trial ou producao (mesmo padrao do Integra Contador).
func (s *NFESerproService) resolveConsultaNFEAPIBase(ambiente string) string {
	if s.serproAuth == nil {
		return ""
	}
	if b := strings.TrimSpace(s.serproAuth.cfg.SerproNFEAPIBaseURL); b != "" {
		return strings.TrimSuffix(b, "/")
	}
	switch normalizeNFEAmbienteIntegra(ambiente) {
	case string(IntegraAmbienteProducao):
		return strings.TrimSuffix(consultaNFEBaseProducao, "/")
	default:
		return strings.TrimSuffix(consultaNFEBaseTrial, "/")
	}
}

func (s *NFESerproService) syncGestaoResumo(ctx context.Context, schemaName string, doc domain.NFEDocumento) error {
	g := BuildNFEGestao(doc.ChaveNFe, doc.PayloadJSON, doc.RecebidoEm)
	_, err := s.repo.UpsertGestao(ctx, schemaName, g)
	return err
}

// ListGestao lista resumos persistidos (EF-919) com filtros e paginação.
func (s *NFESerproService) ListGestao(ctx context.Context, schemaName string, p repository.NFEGestaoListParams) (domain.NFEGestaoListResponse, error) {
	items, total, err := s.repo.ListGestao(ctx, schemaName, p)
	if err != nil {
		return domain.NFEGestaoListResponse{}, err
	}
	return domain.NFEGestaoListResponse{Items: items, TotalRecords: total}, nil
}

func (s *NFESerproService) ListSyncEstado(ctx context.Context, schemaName string, p repository.NFESyncEstadoListParams) (domain.NFESyncEstadoListResponse, error) {
	items, total, err := s.repo.ListSyncEstados(ctx, schemaName, p)
	if err != nil {
		return domain.NFESyncEstadoListResponse{}, err
	}
	return domain.NFESyncEstadoListResponse{Items: items, TotalRecords: total}, nil
}

func (s *NFESerproService) resolveProvider(providerName, uf, ambiente string, simular bool) (nfeprovider.NFeProvider, string, string, error) {
	name := strings.ToUpper(strings.TrimSpace(providerName))
	if name == "" {
		name = "SC"
	}
	if simular {
		return nfeprovider.NewMockProvider(), "MOCK", "SC", nil
	}
	ufNorm := strings.ToUpper(strings.TrimSpace(uf))
	if ufNorm == "" {
		ufNorm = "SC"
	}
	switch name {
	case "SC", "SEF_SC":
		hom := strings.EqualFold(strings.TrimSpace(ambiente), "homologacao") || strings.EqualFold(strings.TrimSpace(ambiente), "hom")
		return nfeprovider.NewSCProvider(hom, "vecontab-ef920"), "SC", "SC", nil
	case "NACIONAL":
		return nfeprovider.NewNacionalProvider(), "NACIONAL", ufNorm, nil
	default:
		return nil, "", "", fmt.Errorf("provider nfe nao suportado: %s", name)
	}
}

func proximaConsultaPorCStat(cstat int) *time.Time {
	now := time.Now().UTC()
	switch cstat {
	case 110:
		t := now.Add(1 * time.Hour)
		return &t
	case 117:
		t := now.Add(12 * time.Hour)
		return &t
	default:
		return nil
	}
}

func onlyDigitsSync(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (s *NFESerproService) SincronizarPorProvider(
	ctx context.Context,
	schemaName, tenantID, providerName, uf, cnpj, ambiente string,
	simular bool,
) (domain.NFESincronizacaoResultado, error) {
	cnpj = onlyDigitsSync(cnpj)
	if len(cnpj) != 14 && len(cnpj) != 11 {
		return domain.NFESincronizacaoResultado{}, fmt.Errorf("cnpj/cpf invalido para sincronizacao")
	}
	if strings.TrimSpace(tenantID) == "" {
		return domain.NFESincronizacaoResultado{}, fmt.Errorf("tenant nao encontrado no contexto")
	}

	provider, providerNorm, ufNorm, err := s.resolveProvider(providerName, uf, ambiente, simular)
	if err != nil {
		return domain.NFESincronizacaoResultado{}, err
	}
	if !simular && (s.certSvc == nil || !s.certSvc.Configurado()) {
		return domain.NFESincronizacaoResultado{}, fmt.Errorf("certificado A1 nao configurado para mTLS")
	}
	if !simular {
		material, err := s.certSvc.MaterialEmMemoria(ctx, tenantID)
		if err != nil {
			return domain.NFESincronizacaoResultado{}, fmt.Errorf("material de certificado: %w", err)
		}
		defer material.Zero()
		if err := provider.ConfigurarCertificado(material.PFX, material.Senha); err != nil {
			return domain.NFESincronizacaoResultado{}, fmt.Errorf("configurar certificado no provider: %w", err)
		}
	}

	estado, err := s.repo.GetSyncEstado(ctx, schemaName, providerNorm, ufNorm, cnpj)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return domain.NFESincronizacaoResultado{}, err
	}
	antNSU := "0"
	if err == nil {
		antNSU = strings.TrimSpace(estado.UltimoNSU)
	}

	res, err := provider.SincronizarDocumentos(ctx, cnpj, antNSU)
	if err != nil {
		return domain.NFESincronizacaoResultado{}, err
	}

	totalPersistidos := 0
	for _, d := range res.Documentos {
		chave := strings.TrimSpace(d.ChaveAcesso)
		if chave == "" {
			continue
		}
		doc := domain.NFEDocumento{
			ChaveNFe:          chave,
			Ambiente:          "producao",
			Origem:            "DOWNLOAD_NFE_PROVIDER_" + providerNorm,
			PayloadJSON:       json.RawMessage(`{}`),
			PayloadXML:        strings.TrimSpace(d.XML),
			ContentTypeOrigem: "application/xml",
			RequestTag:        "provider_sync",
			StatusHTTP:        200,
		}
		if strings.EqualFold(d.Tipo, "procEventoNFe") {
			doc.EventoDescricao = "Evento NF-e (provider SC)"
		}
		saved, err := s.repo.UpsertDocumento(ctx, schemaName, doc)
		if err != nil {
			return domain.NFESincronizacaoResultado{}, err
		}
		_ = s.syncGestaoResumo(ctx, schemaName, saved)
		totalPersistidos++
	}

	now := time.Now().UTC()
	proxima := proximaConsultaPorCStat(res.CStat)
	_, err = s.repo.UpsertSyncEstado(ctx, schemaName, repository.NFESyncStateUpsert{
		Provider:            providerNorm,
		UF:                  ufNorm,
		CNPJ:                cnpj,
		UltimoNSU:           res.NovoMaxNSU,
		UltimoCStat:         res.CStat,
		UltimoMotivo:        res.XMotivo,
		UltimaVerificacao:   now,
		ProximaConsultaApos: proxima,
	})
	if err != nil {
		return domain.NFESincronizacaoResultado{}, err
	}

	return domain.NFESincronizacaoResultado{
		Provider:         providerNorm,
		UF:               ufNorm,
		CNPJ:             cnpj,
		AnteriorNSU:      antNSU,
		NovoNSU:          res.NovoMaxNSU,
		TotalRecebidos:   len(res.Documentos),
		TotalPersistidos: totalPersistidos,
		CStat:            res.CStat,
		XMotivo:          strings.TrimSpace(res.XMotivo),
	}, nil
}

func (s *NFESerproService) ConsultarNFe(ctx context.Context, schemaName, ambiente, chaveNFe, requestTag string, assinar bool) (domain.NFEDocumento, error) {
	if s.serproAuth == nil {
		return domain.NFEDocumento{}, fmt.Errorf("servico SERPRO nao configurado")
	}
	if err := validateNFEChave(chaveNFe); err != nil {
		return domain.NFEDocumento{}, err
	}

	chave := strings.TrimSpace(chaveNFe)
	if cached, err := s.repo.GetDocumentoByChave(ctx, schemaName, chave); err == nil {
		cached.JaBaixada = true
		_ = s.syncGestaoResumo(ctx, schemaName, cached)
		return cached, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return domain.NFEDocumento{}, err
	}

	ambNorm := normalizeNFEAmbienteIntegra(ambiente)
	base := s.resolveConsultaNFEAPIBase(ambNorm)
	if base == "" {
		return domain.NFEDocumento{}, fmt.Errorf("base URL consulta NF-e nao configurada")
	}

	var token string
	if static := normalizeStaticBearer(s.serproAuth.cfg.SerproNFEBearerToken); static != "" {
		token = static
	} else {
		var err error
		token, err = s.serproAuth.ObterBearerToken(ctx)
		if err != nil {
			return domain.NFEDocumento{}, err
		}
	}

	u := base + "/v1/nfe/" + strings.TrimSpace(chaveNFe)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return domain.NFEDocumento{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if strings.TrimSpace(requestTag) != "" {
		req.Header.Set("x-request-tag", strings.TrimSpace(requestTag))
	}
	if assinar {
		req.Header.Set("x-signature", "1")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.NFEDocumento{}, fmt.Errorf("consulta nfe serpro: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.NFEDocumento{}, fmt.Errorf("consulta nfe serpro status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	jsonPayload, err := normalizeJSON(body)
	if err != nil {
		return domain.NFEDocumento{}, fmt.Errorf("resposta json invalida: %w", err)
	}
	payloadXML, err := jsonToXML(jsonPayload)
	if err != nil {
		return domain.NFEDocumento{}, fmt.Errorf("erro ao converter json para xml: %w", err)
	}

	doc := domain.NFEDocumento{
		ChaveNFe:          strings.TrimSpace(chaveNFe),
		Ambiente:          ambNorm,
		Origem:            "CONSULTA_NFE_SERPRO",
		PayloadJSON:       jsonPayload,
		PayloadXML:        payloadXML,
		ContentTypeOrigem: resp.Header.Get("Content-Type"),
		RequestTag:        strings.TrimSpace(requestTag),
		StatusHTTP:        resp.StatusCode,
	}
	fillEventoMetaFromJSON(&doc)
	out, err := s.repo.UpsertDocumento(ctx, schemaName, doc)
	if err != nil {
		return domain.NFEDocumento{}, err
	}
	if err := s.syncGestaoResumo(ctx, schemaName, out); err != nil {
		return domain.NFEDocumento{}, err
	}
	return out, nil
}

func (s *NFESerproService) BuscarDocumento(ctx context.Context, schemaName, chaveNFe string) (domain.NFEDocumento, error) {
	if err := validateNFEChave(chaveNFe); err != nil {
		return domain.NFEDocumento{}, err
	}
	doc, err := s.repo.GetDocumentoByChave(ctx, schemaName, strings.TrimSpace(chaveNFe))
	if err != nil {
		return doc, err
	}
	_ = s.syncGestaoResumo(ctx, schemaName, doc)
	return doc, nil
}

func (s *NFESerproService) ExportarXML(ctx context.Context, schemaName, chaveNFe string) (string, error) {
	doc, err := s.BuscarDocumento(ctx, schemaName, chaveNFe)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(doc.PayloadXML) != "" {
		return doc.PayloadXML, nil
	}
	if len(doc.PayloadJSON) == 0 {
		return "<nfe/>", nil
	}
	return jsonToXML(doc.PayloadJSON)
}

func defaultNFeXSLTDir() string {
	candidates := []string{
		"frontend/public/svrs-nfe-xslt",
		"../frontend/public/svrs-nfe-xslt",
		filepath.Join("..", "..", "frontend", "public", "svrs-nfe-xslt"),
	}
	for _, rel := range candidates {
		st, err := os.Stat(rel)
		if err != nil || !st.IsDir() {
			continue
		}
		abs, err := filepath.Abs(rel)
		if err == nil {
			return abs
		}
	}
	return ""
}

type saxonXSLTPaths struct {
	jar     string
	java    string
	xslMain string
}

func (s *NFESerproService) resolveSaxonXSLT() (saxonXSLTPaths, error) {
	var out saxonXSLTPaths
	if s.serproAuth == nil {
		return out, fmt.Errorf("servico SERPRO nao configurado")
	}
	cfg := s.serproAuth.cfg
	jar := strings.TrimSpace(cfg.NFeSaxonJAR)
	if jar == "" {
		return out, fmt.Errorf("configure NFE_SAXON_JAR com o caminho do saxon-he-*.jar (https://www.saxonica.com/download/java)")
	}
	if abs, e := filepath.Abs(jar); e == nil {
		jar = abs
	}
	if _, e := os.Stat(jar); e != nil {
		return out, fmt.Errorf("NFE_SAXON_JAR inacessivel: %w", e)
	}
	xsltDir := strings.TrimSpace(cfg.NFeXSLTDir)
	if xsltDir == "" {
		xsltDir = defaultNFeXSLTDir()
	}
	if xsltDir == "" {
		return out, fmt.Errorf("configure NFE_XSLT_DIR apontando para a pasta svrs-nfe-xslt (ex.: .../frontend/public/svrs-nfe-xslt)")
	}
	if abs, e := filepath.Abs(xsltDir); e == nil {
		xsltDir = abs
	}
	xslMain := filepath.Join(xsltDir, "_Visualizacao_Internet.xsl")
	if _, e := os.Stat(xslMain); e != nil {
		return out, fmt.Errorf("folha XSLT _Visualizacao_Internet.xsl nao encontrada em %s", xsltDir)
	}
	java := strings.TrimSpace(cfg.NFeJavaPath)
	if java == "" {
		java = "java"
	}
	out.jar = jar
	out.java = java
	out.xslMain = xslMain
	return out, nil
}

// saxonStderrIsOnlySXWN9019Noise is true when stderr looks like só avisos SXWN9019 do pacote SVRS
// (sem "Error at" / códigos XPath), para não tratar como falha.
func saxonStderrIsOnlySXWN9019Noise(stderr string) bool {
	s := strings.TrimSpace(stderr)
	if s == "" {
		return false
	}
	low := strings.ToLower(s)
	if !strings.Contains(low, "sxwn9019") {
		return false
	}
	if strings.Contains(low, "error at ") {
		return false
	}
	if strings.Contains(low, "fatal error") || strings.Contains(low, "exception in thread") {
		return false
	}
	return !saxonStderrErrCode.MatchString(stderr)
}

// saxonStderrLooksFatal returns true when Saxon stderr indicates a real failure (not only SXWN… warnings).
// SVRS XSLT triggers SXWN9019 (duplicate xsl:include) and Saxon may exit 2 while still writing HTML.
func saxonStderrLooksFatal(stderr string) bool {
	if saxonStderrIsOnlySXWN9019Noise(stderr) {
		return false
	}
	if strings.TrimSpace(stderr) == "" {
		return false
	}
	s := strings.ToLower(stderr)
	if saxonStderrErrorKind.MatchString(stderr) {
		return true
	}
	if strings.Contains(s, "fatal error") {
		return true
	}
	if strings.Contains(s, "exception in thread") {
		return true
	}
	return saxonStderrErrCode.MatchString(stderr)
}

func readDanfeHTMLFileRetries(htmlPath string) ([]byte, error) {
	const attempts = 15
	var lastErr error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			time.Sleep(40 * time.Millisecond)
		}
		b, err := os.ReadFile(htmlPath)
		if err != nil {
			lastErr = err
			continue
		}
		if strings.TrimSpace(string(b)) != "" {
			return b, nil
		}
		lastErr = fmt.Errorf("arquivo HTML vazio")
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("arquivo HTML inexistente ou vazio")
	}
	return nil, lastErr
}

func (s *NFESerproService) danfeHTMLWithSaxon(ctx context.Context, xmlStr string) (string, error) {
	sx, err := s.resolveSaxonXSLT()
	if err != nil {
		return "", err
	}
	xmlStr = strings.TrimSpace(xmlStr)
	if xmlStr == "" || xmlStr == "<nfe/>" {
		return "", fmt.Errorf("XML vazio ou sem conteudo de NF-e")
	}

	tmpDir, err := os.MkdirTemp("", "vecontab-danfe-*")
	if err != nil {
		return "", fmt.Errorf("mkdir temp: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	xmlPath := filepath.Join(tmpDir, "nfe.xml")
	if err := os.WriteFile(xmlPath, []byte(xmlStr), 0600); err != nil {
		return "", err
	}

	htmlPath := filepath.Join(tmpDir, "danfe.html")

	javaArgs := []string{"-jar", sx.jar, "-s:" + xmlPath, "-xsl:" + sx.xslMain}

	// Milhares de linhas de SXWN9019 no stderr podem encher o pipe e travar o Java
	// (stdout bloqueado). Descartamos stderr nas tentativas que buscam HTML.
	var stdout1 bytes.Buffer
	cmd1 := exec.CommandContext(ctx, sx.java, javaArgs...)
	cmd1.Stderr = io.Discard
	cmd1.Stdout = &stdout1
	_ = cmd1.Run()
	out := strings.TrimSpace(stdout1.String())
	if out != "" {
		return out, nil
	}

	cmd2 := exec.CommandContext(ctx, sx.java, append(javaArgs, "-o:"+htmlPath)...)
	cmd2.Stderr = io.Discard
	cmd2.Stdout = io.Discard
	_ = cmd2.Run()
	fileBytes, fileErr := readDanfeHTMLFileRetries(htmlPath)
	outFile := strings.TrimSpace(string(fileBytes))
	if outFile != "" {
		return outFile, nil
	}

	// Diagnóstico: uma execução com stderr capturado
	var stderr3, stdout3 bytes.Buffer
	cmd3 := exec.CommandContext(ctx, sx.java, javaArgs...)
	cmd3.Stderr = &stderr3
	cmd3.Stdout = &stdout3
	runErr3 := cmd3.Run()
	msg3 := strings.TrimSpace(stderr3.String())
	out3 := strings.TrimSpace(stdout3.String())
	if out3 != "" && !saxonStderrLooksFatal(msg3) {
		return out3, nil
	}
	if out3 != "" && saxonStderrIsOnlySXWN9019Noise(msg3) {
		return out3, nil
	}
	if runErr3 != nil {
		if msg3 != "" {
			return "", fmt.Errorf("saxon: %w (%s)", runErr3, saxonTruncateRunLog(msg3))
		}
		return "", fmt.Errorf("saxon: %w", runErr3)
	}
	if fileErr != nil {
		return "", fmt.Errorf("saxon nao gerou HTML (%v); verifique XML nfeProc e XSLT SVRS", fileErr)
	}
	return "", fmt.Errorf("saxon nao gerou HTML; o XML precisa ser NF-e autorizada (nfeProc) no namespace da SEFAZ")
}

// saxonTruncateRunLog limita texto devolvido ao cliente (stderr Saxon pode ser enorme).
func saxonTruncateRunLog(s string) string {
	const max = 8000
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

// GerarDanfeHTMLFromXML transforma um XML de NF-e ja obtido (ex.: trial no quadro Retorno) sem consultar o banco.
func (s *NFESerproService) GerarDanfeHTMLFromXML(ctx context.Context, xmlStr string) (string, error) {
	return s.danfeHTMLWithSaxon(ctx, xmlStr)
}

// ExportarDanfeHTML gera o HTML da DANFE a partir do XML persistido no tenant (chave).
func (s *NFESerproService) ExportarDanfeHTML(ctx context.Context, schemaName, chaveNFe string) (string, error) {
	xmlStr, err := s.ExportarXML(ctx, schemaName, chaveNFe)
	if err != nil {
		return "", err
	}
	xmlStr = strings.TrimSpace(xmlStr)
	if xmlStr == "" || xmlStr == "<nfe/>" {
		return "", fmt.Errorf("nao ha XML de NF-e para transformar; consulte ou busque a nota no tenant antes")
	}
	return s.danfeHTMLWithSaxon(ctx, xmlStr)
}

func (s *NFESerproService) RegistrarPushNotificacao(ctx context.Context, rawBody []byte, headers map[string]string) error {
	var payload map[string]any
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		payload = map[string]any{"raw_body": string(rawBody)}
	}
	bPayload, _ := json.Marshal(payload)
	bHeaders, _ := json.Marshal(headers)
	notif := domain.NFEPushNotificacao{
		Payload: json.RawMessage(bPayload),
		Headers: json.RawMessage(bHeaders),
	}
	if v, ok := payload["chaveNFe"]; ok {
		notif.ChaveNFe = strings.TrimSpace(fmt.Sprintf("%v", v))
	}
	if v, ok := payload["dataHoraEnvio"]; ok {
		ts := strings.TrimSpace(fmt.Sprintf("%v", v))
		if ts != "" {
			if parsed, err := parseTimeFlexible(ts); err == nil {
				notif.DataHoraEnvio = &parsed
			}
		}
	}
	return s.repo.SavePushNotificacao(ctx, notif)
}

func normalizeJSON(raw []byte) (json.RawMessage, error) {
	var anyPayload any
	if err := json.Unmarshal(raw, &anyPayload); err != nil {
		return nil, err
	}
	b, err := json.Marshal(anyPayload)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

func fillEventoMetaFromJSON(doc *domain.NFEDocumento) {
	if doc == nil || len(doc.PayloadJSON) == 0 {
		return
	}
	var obj map[string]any
	if err := json.Unmarshal(doc.PayloadJSON, &obj); err != nil {
		return
	}
	if nfeProc, ok := obj["nfeProc"].(map[string]any); ok {
		if protNFe, ok := nfeProc["protNFe"].(map[string]any); ok {
			if infProt, ok := protNFe["infProt"].(map[string]any); ok {
				if cStat, ok := infProt["cStat"]; ok {
					doc.EventoCodigo = strings.TrimSpace(fmt.Sprintf("%v", cStat))
				}
				if motivo, ok := infProt["xMotivo"]; ok {
					doc.EventoDescricao = strings.TrimSpace(fmt.Sprintf("%v", motivo))
				}
			}
		}
	}
}

func parseTimeFlexible(v string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	trimmed := strings.TrimSpace(v)
	for _, layout := range layouts {
		if t, err := time.Parse(layout, trimmed); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato de data/hora invalido")
}

func jsonToXML(raw json.RawMessage) (string, error) {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", err
	}
	var b bytes.Buffer
	b.WriteString(xml.Header)
	b.WriteString("<nfe>")
	if err := writeXMLValue(&b, "documento", v); err != nil {
		return "", err
	}
	b.WriteString("</nfe>")
	return b.String(), nil
}

func writeXMLValue(buf *bytes.Buffer, key string, v any) error {
	tag := sanitizeXMLTag(key)
	switch val := v.(type) {
	case map[string]any:
		buf.WriteString("<" + tag + ">")
		for k, child := range val {
			if err := writeXMLValue(buf, k, child); err != nil {
				return err
			}
		}
		buf.WriteString("</" + tag + ">")
	case []any:
		for i, item := range val {
			itemTag := tag + "_item_" + strconv.Itoa(i)
			if err := writeXMLValue(buf, itemTag, item); err != nil {
				return err
			}
		}
	default:
		buf.WriteString("<" + tag + ">")
		if err := xml.EscapeText(buf, []byte(fmt.Sprintf("%v", val))); err != nil {
			return err
		}
		buf.WriteString("</" + tag + ">")
	}
	return nil
}

func sanitizeXMLTag(tag string) string {
	t := strings.TrimSpace(strings.ToLower(tag))
	if t == "" {
		return "campo"
	}
	var out strings.Builder
	for _, r := range t {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			out.WriteRune(r)
			continue
		}
		out.WriteRune('_')
	}
	s := out.String()
	if s == "" {
		return "campo"
	}
	if s[0] >= '0' && s[0] <= '9' {
		return "f_" + s
	}
	return s
}
