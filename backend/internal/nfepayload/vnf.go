// Package nfepayload extrai totais de JSON de NF-e (layouts nfeProc / NFe) sem depender de service/repository.
package nfepayload

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ValorTotalFromJSON devolve vNF em total.ICMSTot quando o JSON segue o layout usual (SERPRO / SEFAZ).
func ValorTotalFromJSON(raw []byte) *float64 {
	if len(raw) == 0 {
		return nil
	}
	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil
	}
	inf := findInfNFeMap(root)
	if inf == nil {
		return nil
	}
	tot, ok := inf["total"].(map[string]any)
	if !ok {
		return nil
	}
	icms, ok := tot["ICMSTot"].(map[string]any)
	if !ok {
		return nil
	}
	return parseDecimalPtr(stringFromAny(icms["vNF"]))
}

func findInfNFeMap(root map[string]any) map[string]any {
	if proc, ok := root["nfeProc"].(map[string]any); ok {
		if nfe, ok := proc["NFe"].(map[string]any); ok {
			if inf, ok := nfe["infNFe"].(map[string]any); ok {
				return inf
			}
		}
	}
	if nfe, ok := root["NFe"].(map[string]any); ok {
		if inf, ok := nfe["infNFe"].(map[string]any); ok {
			return inf
		}
	}
	return nil
}

func stringFromAny(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x)
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case json.Number:
		return strings.TrimSpace(x.String())
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", v))
	}
}

func parseDecimalPtr(s string) *float64 {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", "."))
	if s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}
