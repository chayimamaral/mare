// Package masterkey resolve a passphrase usada para ENC: nas variaveis sensiveis (AES-256-GCM).
// Nunca commite a chave real: use VECX_MASTER_KEY em .env.senha_compilacao (gitignored) ou injecao via build.
package masterkey

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	mu       sync.Mutex
	cached   string
	resolved bool
)

// EmbeddedMasterKey pode ser injetada em build com:
// -ldflags "-X 'github.com/chayimamaral/vecx/backend/pkg/masterkey.EmbeddedMasterKey=...'"
var EmbeddedMasterKey string

// Passphrase retorna a senha mestre para criptografia de .env (mesma ordem de precedencia em todo o projeto).
// Chamadas seguintes reutilizam o valor da primeira resolucao bem-sucedida.
func Passphrase() (string, error) {
	mu.Lock()
	defer mu.Unlock()
	if resolved {
		return cached, nil
	}
	p, err := resolveUnlocked()
	if err != nil {
		return "", err
	}
	cached = p
	resolved = true
	return cached, nil
}

// ResetForTests limpa o cache (apenas testes).
func ResetForTests() {
	mu.Lock()
	defer mu.Unlock()
	cached = ""
	resolved = false
}

func resolveUnlocked() (string, error) {
	if v := strings.TrimSpace(os.Getenv("VECX_MASTER_KEY")); v != "" {
		return v, nil
	}
	if v := strings.TrimSpace(os.Getenv("VECONTAB_MASTER_KEY")); v != "" {
		return v, nil
	}
	if v := strings.TrimSpace(os.Getenv("SENHA_COMPILACAO")); v != "" {
		return v, nil
	}
	if v := strings.TrimSpace(EmbeddedMasterKey); v != "" {
		return v, nil
	}

	path := strings.TrimSpace(os.Getenv("VECX_MASTER_KEY_FILE"))
	if path == "" {
		path = strings.TrimSpace(os.Getenv("VECONTAB_MASTER_KEY_FILE"))
	}
	explicitFile := path != ""
	if path == "" {
		candidates := []string{
			".env.senha_compilacao",
			filepath.Join("backend", ".env.senha_compilacao"),
			filepath.Join("..", ".env.senha_compilacao"),
		}
		for _, c := range candidates {
			if _, statErr := os.Stat(c); statErr == nil {
				path = c
				break
			}
		}
		if path == "" {
			path = ".env.senha_compilacao"
		}
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if explicitFile {
				return "", fmt.Errorf("VECX_MASTER_KEY_FILE (%s) nao encontrado", path)
			}
			return "", nil
		}
		return "", fmt.Errorf("ler %s: %w", path, err)
	}

	parsed, perr := godotenv.Parse(bytes.NewReader(b))
	if perr == nil {
		if v := strings.TrimSpace(parsed["VECX_MASTER_KEY"]); v != "" {
			return v, nil
		}
		if v := strings.TrimSpace(parsed["VECONTAB_MASTER_KEY"]); v != "" {
			return v, nil
		}
		if v := strings.TrimSpace(parsed["SENHA_COMPILACAO"]); v != "" {
			return v, nil
		}
	}

	return "", fmt.Errorf(
		"arquivo %s: defina VECX_MASTER_KEY (ou VECONTAB_MASTER_KEY/SENHA_COMPILACAO) em formato .env",
		path,
	)
}
