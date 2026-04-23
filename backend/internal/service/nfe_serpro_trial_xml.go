package service

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const sefazNFeXMLMarker = "portalfiscal.inf.br/nfe"

// DanfeXMLFromConsultaRetorno interpreta o texto do quadro "Retorno": XML SEFAZ, ou JSON do documento
// (payload_xml oficial, ou payload_json do trial SERPRO convertido para nfeProc).
func DanfeXMLFromConsultaRetorno(retorno string) (string, error) {
	r := strings.TrimSpace(retorno)
	if r == "" {
		return "", fmt.Errorf("retorno vazio")
	}
	if strings.HasPrefix(r, "<") {
		if strings.Contains(r, sefazNFeXMLMarker) {
			return r, nil
		}
		if strings.Contains(r, "<nfe>") && strings.Contains(r, "documento") {
			return "", fmt.Errorf("este XML foi gerado a partir do JSON generico; use o JSON completo no retorno para converter em layout trial (payload_json)")
		}
		return r, nil
	}
	if !strings.HasPrefix(r, "{") {
		return "", fmt.Errorf("retorno deve ser XML ou JSON")
	}
	var wrap map[string]any
	if err := json.Unmarshal([]byte(r), &wrap); err != nil {
		return "", fmt.Errorf("JSON invalido: %w", err)
	}
	if _, ok := wrap["nfeProc"]; ok {
		return serproTrialJSONToNFeProcXML([]byte(r))
	}
	if _, ok := wrap["NFe"]; ok {
		return serproTrialJSONToNFeProcXML([]byte(r))
	}
	if px, ok := wrap["payload_xml"].(string); ok {
		px = strings.TrimSpace(px)
		if px != "" && strings.Contains(px, sefazNFeXMLMarker) {
			return px, nil
		}
	}
	pj, ok := wrap["payload_json"]
	if !ok || pj == nil {
		return "", fmt.Errorf("JSON sem payload_json (resposta trial/documento)")
	}
	inner, err := json.Marshal(pj)
	if err != nil {
		return "", err
	}
	return serproTrialJSONToNFeProcXML(inner)
}

func serproTrialJSONToNFeProcXML(raw json.RawMessage) (string, error) {
	var root any
	if err := json.Unmarshal(raw, &root); err != nil {
		return "", err
	}
	var procMap map[string]any
	switch v := root.(type) {
	case map[string]any:
		if inner, ok := v["nfeProc"].(map[string]any); ok {
			procMap = inner
		} else if _, ok := v["NFe"]; ok {
			procMap = v
		} else {
			return "", fmt.Errorf("payload_json trial: falta nfeProc ou NFe (estrutura SERPRO inesperada)")
		}
	default:
		return "", fmt.Errorf("payload_json trial: raiz deve ser objeto")
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	versao := "4.00"
	if v, ok := procMap["versao"].(string); ok && strings.TrimSpace(v) != "" {
		versao = strings.TrimSpace(v)
	} else if v, ok := procMap["versao"].(float64); ok {
		versao = formatTrialNumber(v)
	}
	buf.WriteString(`<nfeProc xmlns="http://www.portalfiscal.inf.br/nfe" versao="`)
	if err := xml.EscapeText(&buf, []byte(versao)); err != nil {
		return "", err
	}
	buf.WriteString(`">`)
	for _, k := range sortedTrialKeys(procMap) {
		if k == "versao" {
			continue
		}
		if err := writeTrialNFeXMLNode(&buf, k, procMap[k]); err != nil {
			return "", err
		}
	}
	buf.WriteString(`</nfeProc>`)
	return buf.String(), nil
}

func sortedTrialKeys(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

func writeTrialNFeXMLNode(buf *bytes.Buffer, tag string, v any) error {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case map[string]any:
		return writeTrialElementMap(buf, tag, val)
	case []any:
		for _, item := range val {
			if err := writeTrialNFeXMLNode(buf, tag, item); err != nil {
				return err
			}
		}
		return nil
	default:
		buf.WriteString("<" + tag + ">")
		if err := xml.EscapeText(buf, []byte(fmtTrialScalar(val))); err != nil {
			return err
		}
		buf.WriteString("</" + tag + ">")
		return nil
	}
}

func writeTrialElementMap(buf *bytes.Buffer, tag string, m map[string]any) error {
	buf.WriteString("<" + tag)
	if id, ok := m["Id"]; ok && isTrialScalar(id) {
		buf.WriteString(` Id="`)
		if err := xml.EscapeText(buf, []byte(fmtTrialScalar(id))); err != nil {
			return err
		}
		buf.WriteString(`"`)
	}
	buf.WriteString(">")
	for _, k := range sortedTrialKeys(m) {
		if k == "Id" {
			continue
		}
		if err := writeTrialNFeXMLNode(buf, k, m[k]); err != nil {
			return err
		}
	}
	buf.WriteString("</" + tag + ">")
	return nil
}

func isTrialScalar(v any) bool {
	switch v.(type) {
	case string, bool, float64, json.Number:
		return true
	default:
		return false
	}
}

func fmtTrialScalar(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		return formatTrialNumber(x)
	case json.Number:
		return x.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatTrialNumber(x float64) string {
	if x == float64(int64(x)) {
		return strconv.FormatInt(int64(x), 10)
	}
	return strconv.FormatFloat(x, 'f', -1, 64)
}
