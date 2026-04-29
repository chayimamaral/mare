package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chayimamaral/vecx/backend/pkg/masterkey"
	"github.com/joho/godotenv"
)

// Compilacao: cd tools/encryptor && go build -o encryptor .

// encryptor <nome do arquivo .env.<cliente>

func main() {
	_ = godotenv.Load()
	// Raiz do repo (quando o comando roda em tools/encryptor); depois cwd.
	_ = godotenv.Overload(filepath.Join("..", "..", ".env.senha_compilacao"))
	_ = godotenv.Overload(filepath.Join("..", "..", "backend", ".env.senha_compilacao"))
	_ = godotenv.Overload(".env.senha_compilacao")

	mk, err := masterkey.Passphrase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "senha mestre: %v\n", err)
		os.Exit(2)
	}
	if strings.TrimSpace(mk) == "" {
		fmt.Fprintln(os.Stderr, "defina VECX_MASTER_KEY (ou VECONTAB_MASTER_KEY/SENHA_COMPILACAO) em .env.senha_compilacao")
		os.Exit(2)
	}

	envPath, err := promptEnvPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro de leitura da entrada: %v\n", err)
		os.Exit(2)
	}
	if strings.TrimSpace(envPath) == "" {
		fmt.Fprintln(os.Stderr, "nome do arquivo .env nao informado")
		os.Exit(2)
	}
	outputPath := strings.TrimSpace(envPath) + ".dist"
	if err := processEnvFile(envPath, outputPath, mk, false); err != nil {
		fmt.Fprintf(os.Stderr, "erro ao processar .env: %v\n", err)
		os.Exit(1)
	}
}

func encryptValue(plain string, passphrase string) (string, error) {
	key := sha256.Sum256([]byte(passphrase))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plain), nil)
	payload := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(payload), nil
}

func decryptValue(encryptedBase64 string, passphrase string) (string, error) {
	raw := strings.TrimSpace(encryptedBase64)
	if raw == "" {
		return "", nil
	}
	payload, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", err
	}
	key := sha256.Sum256([]byte(passphrase))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(payload) <= nonceSize {
		return "", fmt.Errorf("payload invalido")
	}
	nonce := payload[:nonceSize]
	ciphertext := payload[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func processEnvFile(envPath, outputPath, passphrase string, inplace bool) error {
	absEnvPath, err := filepath.Abs(strings.TrimSpace(envPath))
	if err != nil {
		return err
	}

	raw, err := os.ReadFile(absEnvPath)
	if err != nil {
		return err
	}

	content := string(raw)
	lines := strings.Split(content, "\n")
	changed := 0

	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}

		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		if !strings.HasPrefix(value, "ENC:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(value, "ENC:"))
		payload = unwrapQuoted(payload)
		if payload == "" {
			continue
		}

		// Se ja for criptografado valido, mantem.
		if _, decErr := decryptValue(payload, passphrase); decErr == nil {
			continue
		}

		enc, encErr := encryptValue(payload, passphrase)
		if encErr != nil {
			return fmt.Errorf("chave %s: %w", key, encErr)
		}
		lines[i] = key + "=ENC:" + enc
		changed++
	}

	newContent := strings.Join(lines, "\n")
	if !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}

	if !inplace {
		dest := strings.TrimSpace(outputPath)
		if dest == "" {
			dest = absEnvPath + ".dist"
		}
		absDest, e := filepath.Abs(dest)
		if e != nil {
			return e
		}
		if err := os.WriteFile(absDest, []byte(newContent), 0o600); err != nil {
			return err
		}
		fmt.Printf("Arquivo gerado: %s | entradas criptografadas: %d\n", absDest, changed)
		return nil
	}

	backup := absEnvPath + ".bak-" + time.Now().Format("20060102-150405")
	if err := os.WriteFile(backup, raw, 0o600); err != nil {
		return err
	}
	if err := os.WriteFile(absEnvPath, []byte(newContent), 0o600); err != nil {
		return err
	}
	fmt.Printf(".env atualizado: %s | backup: %s | entradas criptografadas: %d\n", absEnvPath, backup, changed)
	return nil
}

func promptEnvPath() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Arquivo de ambiente (.env, .env.globalbusiness, env.jpncontabilidade, etc): ")
	envPath, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(envPath), nil
}

func unwrapQuoted(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
