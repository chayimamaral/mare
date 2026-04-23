package service

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// JSON envelope como retornado pela API (payload_json com nfeProc completo).
const trialEnvelopeJSON = `{
  "id": "cb6c7d17-cf37-4db0-ac59-224d12c20900",
  "chave_nfe": "35170608530528000184550000000154301000771561",
  "payload_json": {
    "nfeProc": {
      "NFe": {
        "infNFe": {
          "Id": "NFe35170608530528000184550000000154301000771561",
          "det": [
            {
              "prod": {
                "NCM": 48025610,
                "CFOP": 5102,
                "qCom": 5,
                "uCom": "RS",
                "cProd": 346,
                "vProd": 745,
                "xProd": "SULFITE A4",
                "vUnCom": 149
              },
              "nItem": 1
            }
          ],
          "ide": {
            "cUF": 35,
            "mod": 55,
            "serie": 0,
            "nNF": 15430,
            "dhEmi": "2017-06-05T08:31:06-03:00",
            "tpAmb": 1
          },
          "dest": {
            "CPF": 67879577696,
            "xNome": "MARIA FICTICIA"
          },
          "emit": {
            "CNPJ": "56776378000136",
            "xNome": "COMERCIO DE TESTE LTDA"
          },
          "total": {
            "ICMSTot": { "vNF": 745, "vProd": 745 }
          }
        }
      },
      "versao": 3.1,
      "protNFe": {
        "infProt": {
          "cStat": 100,
          "chNFe": "35170608530528000184550000000154301000771561",
          "dhRecbto": "2017-06-05T08:31:06-03:00"
        }
      }
    }
  }
}`

func TestDanfeXMLFromConsultaRetorno_documentEnvelope(t *testing.T) {
	x, err := DanfeXMLFromConsultaRetorno(trialEnvelopeJSON)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(x, "portalfiscal.inf.br/nfe") {
		t.Fatal("namespace missing")
	}
	if !strings.Contains(x, "<nfeProc") || !strings.Contains(x, "<protNFe") {
		t.Fatalf("unexpected xml: %s", x[:min(400, len(x))])
	}
}

func TestDanfeXML_envelopeSaxonProducesBody(t *testing.T) {
	repo := filepath.Join("..", "..", "..")
	jar := filepath.Join(repo, "docs", "arquitetura", "SaxonHE12-9J", "saxon-he-12.9.jar")
	xsl := filepath.Join(repo, "frontend", "public", "svrs-nfe-xslt", "_Visualizacao_Internet.xsl")
	if _, err := os.Stat(jar); err != nil {
		t.Skip("Saxon JAR not present:", jar)
	}
	if _, err := os.Stat(xsl); err != nil {
		t.Skip("XSLT dir not present:", xsl)
	}
	if _, err := exec.LookPath("java"); err != nil {
		t.Skip("java not in PATH")
	}
	x, err := DanfeXMLFromConsultaRetorno(trialEnvelopeJSON)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("java", "-jar", jar, "-s:-", "-xsl:"+xsl)
	cmd.Stdin = strings.NewReader(x)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("saxon: %v\nstderr=%s", err, string(ee.Stderr))
		}
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "<body") || !strings.Contains(html, "Consulta da NF-e") {
		t.Fatalf("unexpected html head: %s", html[:min(600, len(html))])
	}
}
