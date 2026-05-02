package certtax

import (
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"strings"
)

// TaxIDsFromCertificate extrai candidatos a CNPJ (14) ou CPF (11) do certificado ICP-Brasil.
// Valores retornados contêm apenas dígitos.
func TaxIDsFromCertificate(cert *x509.Certificate) []string {
	if cert == nil {
		return nil
	}
	seen := map[string]struct{}{}
	var out []string
	add := func(raw string) {
		d := onlyDigits(raw)
		if len(d) != 11 && len(d) != 14 {
			return
		}
		if _, ok := seen[d]; ok {
			return
		}
		seen[d] = struct{}{}
		out = append(out, d)
	}

	for _, name := range cert.Subject.Names {
		add(fmt.Sprint(name.Value))
	}
	add(cert.Subject.CommonName)
	add(cert.Subject.SerialNumber)

	for _, ext := range cert.Extensions {
		if len(ext.Value) == 0 {
			continue
		}
		for _, s := range digitRunsFromBytes(ext.Value) {
			add(s)
		}
	}

	for _, s := range digitRunsFromBytes(cert.RawSubject) {
		add(s)
	}

	for _, s := range otherNamesFromSAN(cert) {
		add(s)
	}
	return out
}

func digitRunsFromBytes(b []byte) []string {
	var runs []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() >= 11 {
			runs = append(runs, cur.String())
		}
		cur.Reset()
	}
	for _, c := range b {
		if c >= '0' && c <= '9' {
			cur.WriteByte(c)
			continue
		}
		flush()
	}
	flush()
	return runs
}

func onlyDigits(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// NormalizeTaxID mantém apenas dígitos (CNPJ/CPF para pastas e comparação).
func NormalizeTaxID(s string) string {
	return onlyDigits(strings.TrimSpace(s))
}

func otherNamesFromSAN(cert *x509.Certificate) []string {
	if cert == nil {
		return nil
	}
	sanOID := asn1.ObjectIdentifier{2, 5, 29, 17}
	for _, ext := range cert.Extensions {
		if !ext.Id.Equal(sanOID) {
			continue
		}
		var seq asn1.RawValue
		if _, err := asn1.Unmarshal(ext.Value, &seq); err != nil {
			continue
		}
		rest := seq.Bytes
		var out []string
		for len(rest) > 0 {
			var raw asn1.RawValue
			tr, err := asn1.Unmarshal(rest, &raw)
			if err != nil {
				break
			}
			rest = tr
			if raw.Class == 2 && raw.Tag == 0 {
				var inner asn1.RawValue
				if _, err := asn1.Unmarshal(raw.Bytes, &inner); err == nil {
					out = append(out, string(inner.Bytes))
				}
			}
		}
		return out
	}
	return nil
}
