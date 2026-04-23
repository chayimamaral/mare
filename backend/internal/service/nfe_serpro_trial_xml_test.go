package service

import (
	"strings"
	"testing"
)

func TestDanfeXMLFromConsultaRetorno_NFeRoot(t *testing.T) {
	const j = `{"NFe":{"infNFe":{"Id":"NFe35200152764578000158550010000000011000000010","ide":{"cUF":35,"mod":"55","serie":1,"nNF":1,"dhEmi":"2024-01-15T10:00:00-03:00"}}},"protNFe":{"infProt":{"chNFe":"35200152764578000158550010000000011000000010"}}}`
	x, err := DanfeXMLFromConsultaRetorno(j)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(x, "http://www.portalfiscal.inf.br/nfe") {
		t.Fatalf("missing namespace: %s", x[:min(200, len(x))])
	}
	if !strings.Contains(x, `Id="NFe35200152764578000158550010000000011000000010"`) {
		t.Fatal("missing infNFe Id attribute")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestSaxonStderrLooksFatal(t *testing.T) {
	if saxonStderrLooksFatal("") {
		t.Fatal("empty stderr should not be fatal")
	}
	if !saxonStderrLooksFatal("Error at xsl:template on line 1") {
		t.Fatal("Saxon error line should be fatal")
	}
	if !saxonStderrLooksFatal("Type error at line 5 in template") {
		t.Fatal("Saxon type error line should be fatal")
	}
	if !saxonStderrLooksFatal("XPath failure XPST0008 in expression") {
		t.Fatal("XPath error code should be fatal")
	}
	if saxonStderrLooksFatal("Warning at xsl:stylesheet ... SXWN9019 ... included more than once") {
		t.Fatal("SVRS duplicate-import warning should not be fatal")
	}
	if saxonStderrLooksFatal("Warning at xsl:stylesheet\n  may lead to errors or unexpected behavior") {
		t.Fatal("boilerplate 'errors or unexpected' should not be fatal")
	}
	if !saxonStderrIsOnlySXWN9019Noise("Warning at xsl:stylesheet\n  SXWN9019 duplicate include") {
		t.Fatal("SXWN9019-only blob should be noise")
	}
	if saxonStderrIsOnlySXWN9019Noise("Error at line 1") {
		t.Fatal("Error at should not be classified as noise-only")
	}
}
