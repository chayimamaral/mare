package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chayimamaral/vecx/backend/pkg/masterkey"
	"github.com/joho/godotenv"
)

type Config struct {
	Runtime            string
	Port               string
	CORSAllowedOrigins []string
	DatabaseURL        string
	// AuditDatabaseURL: banco global VECX_AUDIT (EF-929). Obrigatório; mesma stack TLS que PG_URL.
	AuditDatabaseURL               string
	JWTSecret                      string
	SSLRootCertPath                string
	SSLInsecure                    bool
	CompromissosWorkerEnabled      bool
	CompromissosWorkerCron         string
	CompromissosWorkerRunOnStartup bool
	CompromissosWorkerTimezone     string

	// Worker de sincronização NFe por provider (EF-920): percorre nfe_sync_estado em todos os tenants.
	NFESyncWorkerEnabled         bool
	NFESyncWorkerCron            string
	NFESyncWorkerRunOnStartup    bool
	NFESyncWorkerTimezone        string
	NFESyncWorkerAmbiente        string
	NFESyncWorkerMinIntervalSecs int

	// CertCryptoKeyHex: 64 caracteres hex (32 bytes) para AES-256-GCM de PFX/senha (issue #55).
	CertCryptoKeyHex string
	// SERPRO Integra Contador — OAuth2 cliente (credenciais de desenvolvedor); URLs conforme documentação oficial.
	SerproOAuthTokenURL string
	SerproClientID      string
	SerproClientSecret  string
	SerproAPIBaseURL    string
	SerproNFEAPIBaseURL string
	// SerproNFEBearerToken: opcional. Se preenchido, a consulta NF-e usa este Bearer direto (ex.: API de teste),
	// sem chamar OAuth2 (SERPRO_OAUTH_TOKEN_URL / CLIENT_ID / SECRET).
	SerproNFEBearerToken string

	// DANFE HTML (XSLT 2.0 SVRS): Saxon-HE + pasta com os .xsl (ex.: frontend/public/svrs-nfe-xslt).
	// https://www.saxonica.com/download/java
	NFeSaxonJAR string
	NFeJavaPath string
	NFeXSLTDir  string
}

func Load() (Config, error) {
	// Tenta carregar o .env do diretório atual
	_ = godotenv.Load()
	// Senha mestre para ENC: (nao versionar — ver env.senha_compilacao.example)
	_ = godotenv.Load(".env.senha_compilacao")
	_ = godotenv.Load("../.env.senha_compilacao")
	_ = godotenv.Load("backend/.env.senha_compilacao")
	// if err != nil {
	// 	fmt.Println("Aviso: .env não encontrado no diretório atual, tentando caminho relativo...")
	// }
	// err = godotenv.Load(".env") // Ajuste o caminho se necessário
	// if err != nil {
	// 	// Se este também falhar, talvez seja um problema real
	// 	fmt.Printf("Erro ao carregar ../../.env: %v\n", err)
	// }

	runtime := strings.ToLower(strings.TrimSpace(getEnv("VECONTAB_RUNTIME", "")))
	if runtime == "" {
		if exe, err := os.Executable(); err == nil {
			exeNorm := strings.ReplaceAll(strings.ToLower(exe), "\\", "/")
			if strings.Contains(exeNorm, "/backend/bin/") || strings.Contains(exeNorm, "/bin/vecontab-backend") {
				runtime = "binary"
			}
		}
	}
	if runtime == "" {
		runtime = "web"
	}

	defaultPort := "8080"
	if runtime == "binary" || runtime == "desktop" {
		defaultPort = "3333"
	}

	pgURL, err := decryptSensitiveEnv("PG_URL")
	if err != nil {
		return Config{}, fmt.Errorf("PG_URL invalida para descriptografia: %w", err)
	}
	auditURL, err := decryptSensitiveEnv("VECX_AUDIT")
	if err != nil {
		return Config{}, fmt.Errorf("VECX_AUDIT invalida para descriptografia: %w", err)
	}
	jwtSecret, err := decryptSensitiveEnv("JWT_SECRET")
	if err != nil {
		return Config{}, fmt.Errorf("JWT_SECRET invalida para descriptografia: %w", err)
	}
	certCryptoKeyHex, err := decryptSensitiveEnv("VECONTAB_CERT_CRYPTO_KEY_HEX")
	if err != nil {
		return Config{}, fmt.Errorf("VECONTAB_CERT_CRYPTO_KEY_HEX invalida para descriptografia: %w", err)
	}
	serproClientID, err := decryptSensitiveEnv("SERPRO_CLIENT_ID")
	if err != nil {
		return Config{}, fmt.Errorf("SERPRO_CLIENT_ID invalida para descriptografia: %w", err)
	}
	serproClientSecret, err := decryptSensitiveEnv("SERPRO_CLIENT_SECRET")
	if err != nil {
		return Config{}, fmt.Errorf("SERPRO_CLIENT_SECRET invalida para descriptografia: %w", err)
	}
	serproNFEBearerToken, err := decryptSensitiveEnv("SERPRO_NFE_BEARER_TOKEN")
	if err != nil {
		return Config{}, fmt.Errorf("SERPRO_NFE_BEARER_TOKEN invalida para descriptografia: %w", err)
	}

	cfg := Config{
		Runtime:                        runtime,
		Port:                           getEnv("PORT", defaultPort),
		CORSAllowedOrigins:             parseCSVEnv("CORS_ALLOWED_ORIGINS"),
		DatabaseURL:                    pgURL,
		AuditDatabaseURL:               auditURL,
		JWTSecret:                      jwtSecret,
		SSLRootCertPath:                getEnv("PG_SSL_ROOT_CERT", "/home/camaral/.postgres/ca.crt"),
		SSLInsecure:                    getEnv("PG_SSL_INSECURE", "true") == "true",
		CompromissosWorkerEnabled:      getEnv("COMPROMISSOS_WORKER_ENABLED", "false") == "true",
		CompromissosWorkerCron:         getEnv("COMPROMISSOS_WORKER_CRON", "0 5 1 * *"),
		CompromissosWorkerRunOnStartup: getEnv("COMPROMISSOS_WORKER_RUN_ON_STARTUP", "false") == "true",
		CompromissosWorkerTimezone:     getEnv("COMPROMISSOS_WORKER_TIMEZONE", "America/Sao_Paulo"),
		NFESyncWorkerEnabled:           getEnv("NFE_SYNC_WORKER_ENABLED", "false") == "true",
		NFESyncWorkerCron:              getEnv("NFE_SYNC_WORKER_CRON", "*/15 * * * *"),
		NFESyncWorkerRunOnStartup:      getEnv("NFE_SYNC_WORKER_RUN_ON_STARTUP", "false") == "true",
		NFESyncWorkerTimezone:          getEnv("NFE_SYNC_WORKER_TIMEZONE", "America/Sao_Paulo"),
		NFESyncWorkerAmbiente:          getEnv("NFE_SYNC_WORKER_AMBIENTE", "producao"),
		NFESyncWorkerMinIntervalSecs:   parseIntEnv("NFE_SYNC_WORKER_MIN_INTERVAL_SECS", 600),
		CertCryptoKeyHex:               certCryptoKeyHex,
		SerproOAuthTokenURL:            getEnv("SERPRO_OAUTH_TOKEN_URL", ""),
		SerproClientID:                 serproClientID,
		SerproClientSecret:             serproClientSecret,
		SerproAPIBaseURL:               getEnv("SERPRO_API_BASE_URL", ""),
		SerproNFEAPIBaseURL:            getEnv("SERPRO_NFE_API_BASE_URL", getEnv("SERPRO_API_BASE_URL", "")),
		SerproNFEBearerToken:           strings.TrimSpace(serproNFEBearerToken),
		NFeSaxonJAR:                    strings.TrimSpace(os.Getenv("NFE_SAXON_JAR")),
		NFeJavaPath:                    getEnv("NFE_JAVA_PATH", "java"),
		NFeXSLTDir:                     strings.TrimSpace(os.Getenv("NFE_XSLT_DIR")),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("PG_URL is required")
	}
	if cfg.AuditDatabaseURL == "" {
		return Config{}, fmt.Errorf("VECX_AUDIT is required")
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func parseIntEnv(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return fallback
	}
	return n
}

func parseCSVEnv(key string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(p)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}

// DecryptValue descriptografa "nonce|ciphertext" (concatenados) em Base64 usando AES-256-GCM.
// A chave AES (32 bytes) é derivada por SHA-256 da passphrase.
func DecryptValue(encryptedBase64 string, passphrase string) (string, error) {
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
		return "", errors.New("payload criptografado invalido")
	}

	nonce := payload[:nonceSize]
	ciphertext := payload[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func decryptSensitiveEnv(key string) (string, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return "", nil
	}

	mk, err := masterkey.Passphrase()
	if err != nil {
		return "", fmt.Errorf("senha mestre para %s: %w", key, err)
	}

	// Modo explicito: ENC:<base64_nonce_ciphertext>
	if strings.HasPrefix(strings.ToUpper(raw), "ENC:") {
		if mk == "" {
			return "", fmt.Errorf(
				"%s usa ENC: defina VECX_MASTER_KEY (ou VECONTAB_MASTER_KEY/SENHA_COMPILACAO) em .env.senha_compilacao",
				key,
			)
		}
		return DecryptValue(strings.TrimSpace(raw[4:]), mk)
	}

	// Compatibilidade: se vier puro em texto, mantem; se vier base64 criptografado, descriptografa.
	if mk == "" {
		return raw, nil
	}
	dec, err := DecryptValue(raw, mk)
	if err != nil {
		return raw, nil
	}
	return dec, nil
}
