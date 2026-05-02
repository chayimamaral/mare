package certlayout

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/chayimamaral/vecx/agente-local/internal/certtax"
)

const (
	SubdirClientes = "cert_clientes"
	SubdirContador = "cert_contador"
)

// ResolveClientePFX retorna caminho para {root}/cert_clientes/{taxID}.pfx
func ResolveClientePFX(root, taxID string) (string, error) {
	root = strings.TrimSpace(root)
	taxID = certtax.NormalizeTaxID(taxID)
	if root == "" {
		return "", errors.New("pasta raiz de certificados nao configurada")
	}
	if len(taxID) != 11 && len(taxID) != 14 {
		return "", errors.New("tax_id deve ser CPF (11) ou CNPJ (14) digitos")
	}
	name := taxID + ".pfx"
	p := filepath.Join(root, SubdirClientes, name)
	return p, nil
}

// ResolveContadorPFX escolhe o .pfx do contador em cert_contador (único arquivo ou contador.pfx).
func ResolveContadorPFX(root string) (string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return "", errors.New("pasta raiz de certificados nao configurada")
	}
	dir := filepath.Join(root, SubdirContador)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("pasta %s inexistente", SubdirContador)
		}
		return "", fmt.Errorf("ler %s: %w", dir, err)
	}
	var pfxFiles []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if strings.EqualFold(filepath.Ext(n), ".pfx") {
			pfxFiles = append(pfxFiles, filepath.Join(dir, n))
		}
	}
	if len(pfxFiles) == 0 {
		return "", fmt.Errorf("nenhum .pfx em %s", SubdirContador)
	}
	if len(pfxFiles) == 1 {
		return pfxFiles[0], nil
	}
	var contador []string
	for _, p := range pfxFiles {
		base := filepath.Base(p)
		if strings.EqualFold(base, "contador.pfx") {
			contador = append(contador, p)
		}
	}
	if len(contador) == 1 {
		return contador[0], nil
	}
	return "", fmt.Errorf("cert_contador: varios .pfx; renomeie um para contador.pfx ou mantenha apenas um arquivo")
}
