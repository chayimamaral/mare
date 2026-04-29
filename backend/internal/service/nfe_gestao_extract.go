package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chayimamaral/vecx/backend/internal/domain"
)

// BuildNFEGestao monta o registro de gestão a partir da chave e do JSON SERPRO / nfeProc.
func BuildNFEGestao(chaveNFe string, payload json.RawMessage, dataDownload time.Time) domain.NFEGestao {
	chave := strings.TrimSpace(chaveNFe)
	out := domain.NFEGestao{
		ChaveNFe:     chave,
		DataDownload: dataDownload,
		TipoArquivo:  tipoArquivoFromChave(chave),
	}
	if len(chave) == 44 && strings.Trim(chave, "0123456789") == "" {
		out.CNPJEmitente = chave[6:20]
		out.NumeroNFe = strings.TrimLeft(chave[25:34], "0")
		if out.NumeroNFe == "" {
			out.NumeroNFe = "0"
		}
	}

	if len(payload) == 0 {
		return out
	}

	var root map[string]any
	if err := json.Unmarshal(payload, &root); err != nil {
		return out
	}

	inf := findInfNFeMap(root)
	if inf == nil {
		return out
	}

	if ide, ok := inf["ide"].(map[string]any); ok {
		if v := stringFromAny(ide["nNF"]); v != "" {
			out.NumeroNFe = v
		}
		if t := parseDataEmissaoNFe(ide["dhEmi"], ide["dEmi"]); t != nil {
			out.DataEmissao = t
		}
	}

	if emit, ok := inf["emit"].(map[string]any); ok {
		if v := stringFromAny(emit["xNome"]); v != "" {
			out.RazaoSocialEmitente = v
		}
		if v := normalizeCNPJCPF(stringFromAny(emit["CNPJ"])); v != "" {
			out.CNPJEmitente = v
		} else if v := normalizeCNPJCPF(stringFromAny(emit["CPF"])); v != "" {
			out.CNPJEmitente = v
		}
	}

	if dest, ok := inf["dest"].(map[string]any); ok {
		if v := normalizeCNPJCPF(stringFromAny(dest["CNPJ"])); v != "" {
			out.CNPJDestinatario = v
		} else if v := normalizeCNPJCPF(stringFromAny(dest["CPF"])); v != "" {
			out.CNPJDestinatario = v
		}
	}

	if tot, ok := inf["total"].(map[string]any); ok {
		if icms, ok := tot["ICMSTot"].(map[string]any); ok {
			if x := parseDecimalPtr(stringFromAny(icms["vNF"])); x != nil {
				out.ValorTotal = x
			}
		}
	}

	tipo := detectTipoArquivoFromInf(inf, chave)
	if tipo != "" {
		out.TipoArquivo = tipo
	}

	return out
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

func normalizeCNPJCPF(s string) string {
	return onlyDigitsNFe(s)
}

func onlyDigitsNFe(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

func padModelo2Digits(s string) string {
	d := onlyDigitsNFe(s)
	if len(d) >= 2 {
		return d[len(d)-2:]
	}
	if len(d) == 1 {
		return "0" + d
	}
	return ""
}

func parseDataEmissaoNFe(dhEmi, dEmi any) *time.Time {
	for _, raw := range []any{dhEmi, dEmi} {
		s := stringFromAny(raw)
		if s == "" {
			continue
		}
		layouts := []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02T15:04:05-07:00",
			"2006-01-02T15:04:05",
			"2006-01-02",
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, s); err == nil {
				tt := t.UTC()
				return &tt
			}
		}
	}
	return nil
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

func tipoArquivoFromChave(chave string) string {
	if len(chave) != 44 || strings.Trim(chave, "0123456789") != "" {
		return "Outro"
	}
	return labelModeloNF(chave[20:22])
}

func labelModeloNF(mod string) string {
	switch mod {
	case "55":
		return "NF-e"
	case "65":
		return "NFC-e"
	case "57":
		return "CT-e"
	case "67":
		return "CT-e OS"
	case "07":
		return "NFSe (modelo 07)"
	case "08":
		return "NFSe (modelo 08)"
	default:
		return "Modelo " + mod
	}
}

func detectTipoArquivoFromInf(inf map[string]any, chave string) string {
	mod := padModelo2Digits(stringFromAny(inf["mod"]))
	if mod != "" {
		return labelModeloNF(mod)
	}
	// NFS-e nacional / outros layouts futuros
	if _, ok := inf["DPS"]; ok {
		return "NFS-e Nacional"
	}
	if _, ok := inf["infNFSe"]; ok {
		return "NFS-e Nacional"
	}
	return tipoArquivoFromChave(chave)
}
