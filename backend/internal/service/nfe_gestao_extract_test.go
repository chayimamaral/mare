package service

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBuildNFEGestao_NFeProc(t *testing.T) {
	payload := json.RawMessage(`{
		"nfeProc": {
			"NFe": {
				"infNFe": {
					"mod": "55",
					"ide": { "nNF": "637", "dhEmi": "2024-03-15T10:30:00-03:00" },
					"emit": { "CNPJ": "12345678000199", "xNome": "Emitente Teste SA" },
					"dest": { "CNPJ": "98765432000188" },
					"total": { "ICMSTot": { "vNF": "1500.50" } }
				}
			}
		}
	}`)
	chave := "35240312345678000199550010000006371000000000"
	ts := time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC)
	g := BuildNFEGestao(chave, payload, ts)
	if g.TipoArquivo != "NF-e" {
		t.Fatalf("tipo: got %q", g.TipoArquivo)
	}
	if g.NumeroNFe != "637" {
		t.Fatalf("numero: got %q", g.NumeroNFe)
	}
	if g.RazaoSocialEmitente != "Emitente Teste SA" {
		t.Fatalf("razao: got %q", g.RazaoSocialEmitente)
	}
	if g.CNPJEmitente != "12345678000199" {
		t.Fatalf("cnpj emit: got %q", g.CNPJEmitente)
	}
	if g.CNPJDestinatario != "98765432000188" {
		t.Fatalf("cnpj dest: got %q", g.CNPJDestinatario)
	}
	if g.ValorTotal == nil || *g.ValorTotal != 1500.50 {
		t.Fatalf("valor: %+v", g.ValorTotal)
	}
	if g.DataEmissao == nil {
		t.Fatal("data emissao nil")
	}
}

func TestBuildNFEGestao_SomenteChave(t *testing.T) {
	chave := "35240312345678000199550010000006371000000000"
	ts := time.Now().UTC()
	g := BuildNFEGestao(chave, nil, ts)
	if g.TipoArquivo != "NF-e" {
		t.Fatalf("tipo: got %q", g.TipoArquivo)
	}
	if g.CNPJEmitente != "12345678000199" {
		t.Fatalf("cnpj from chave: got %q", g.CNPJEmitente)
	}
}
