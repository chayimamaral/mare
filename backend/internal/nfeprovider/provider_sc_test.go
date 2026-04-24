package nfeprovider

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"testing"
)

func TestDecodeLoteAndParse(t *testing.T) {
	xml := `<loteDistNFeSC versao="2.00"><distNFeSC NSU="1" chAcesso="35240312345678000199550010000006371000000000"><nfeProc><NFe><infNFe/></NFe></nfeProc></distNFeSC></loteDistNFeSC>`
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write([]byte(xml)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	decoded, err := decodeLote(b64)
	if err != nil {
		t.Fatal(err)
	}
	items, err := parseLote(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d", len(items))
	}
	if items[0].Tipo != "nfeProc" {
		t.Fatalf("tipo = %s", items[0].Tipo)
	}
	if items[0].NSU != "1" {
		t.Fatalf("nsu = %s", items[0].NSU)
	}
}
