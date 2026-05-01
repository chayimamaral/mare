package nfepayload

import "testing"

func TestValorTotalFromJSON_nfeProc(t *testing.T) {
	raw := []byte(`{
		"nfeProc": {
			"NFe": {
				"infNFe": {
					"total": { "ICMSTot": { "vNF": "1500.50" } }
				}
			}
		}
	}`)
	v := ValorTotalFromJSON(raw)
	if v == nil || *v != 1500.50 {
		t.Fatalf("got %+v", v)
	}
}

func TestValorTotalFromJSON_invalid(t *testing.T) {
	if v := ValorTotalFromJSON([]byte(`{}`)); v != nil {
		t.Fatalf("expected nil, got %v", *v)
	}
}
