package nfeprovider

import (
	"strings"
	"testing"
)

func TestParseRetEnvEvento(t *testing.T) {
	raw := `<?xml version="1.0" encoding="UTF-8"?><retEnvEvento xmlns="http://www.portalfiscal.inf.br/nfe" versao="1.00">
  <cStat>128</cStat><xMotivo>Lote processado</xMotivo>
  <retEvento versao="1.00"><infEvento><cStat>135</cStat><xMotivo>Evento registrado</xMotivo><nProt>391123456789012</nProt></infEvento></retEvento>
</retEnvEvento>`
	got, err := parseRetEnvEvento([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if got.CStatLote != 128 || !strings.Contains(got.XMotivoLote, "Lote") {
		t.Fatalf("lote: %+v", got)
	}
	if got.CStatEvento != 135 || got.NProt != "391123456789012" {
		t.Fatalf("evento: %+v", got)
	}
}
