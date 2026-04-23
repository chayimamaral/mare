package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Runtime                        string
	Port                           string
	DatabaseURL                    string
	JWTSecret                      string
	SSLRootCertPath                string
	SSLInsecure                    bool
	CompromissosWorkerEnabled      bool
	CompromissosWorkerCron         string
	CompromissosWorkerRunOnStartup bool
	CompromissosWorkerTimezone     string

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
	NFeSaxonJAR  string
	NFeJavaPath  string
	NFeXSLTDir   string
}

func Load() (Config, error) {
	// Tenta carregar o .env do diretório atual
	_ = godotenv.Load()
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

	cfg := Config{
		Runtime:                        runtime,
		Port:                           getEnv("PORT", defaultPort),
		DatabaseURL:                    os.Getenv("PG_URL"),
		JWTSecret:                      os.Getenv("JWT_SECRET"),
		SSLRootCertPath:                getEnv("PG_SSL_ROOT_CERT", "/home/camaral/.postgres/ca.crt"),
		SSLInsecure:                    getEnv("PG_SSL_INSECURE", "true") == "true",
		CompromissosWorkerEnabled:      getEnv("COMPROMISSOS_WORKER_ENABLED", "false") == "true",
		CompromissosWorkerCron:         getEnv("COMPROMISSOS_WORKER_CRON", "0 5 1 * *"),
		CompromissosWorkerRunOnStartup: getEnv("COMPROMISSOS_WORKER_RUN_ON_STARTUP", "false") == "true",
		CompromissosWorkerTimezone:     getEnv("COMPROMISSOS_WORKER_TIMEZONE", "America/Sao_Paulo"),
		CertCryptoKeyHex:               os.Getenv("VECONTAB_CERT_CRYPTO_KEY_HEX"),
		SerproOAuthTokenURL:            getEnv("SERPRO_OAUTH_TOKEN_URL", ""),
		SerproClientID:                 os.Getenv("SERPRO_CLIENT_ID"),
		SerproClientSecret:             os.Getenv("SERPRO_CLIENT_SECRET"),
		SerproAPIBaseURL:               getEnv("SERPRO_API_BASE_URL", ""),
		SerproNFEAPIBaseURL:            getEnv("SERPRO_NFE_API_BASE_URL", getEnv("SERPRO_API_BASE_URL", "")),
		SerproNFEBearerToken:           strings.TrimSpace(os.Getenv("SERPRO_NFE_BEARER_TOKEN")),
		NFeSaxonJAR:                    strings.TrimSpace(os.Getenv("NFE_SAXON_JAR")),
		NFeJavaPath:                    getEnv("NFE_JAVA_PATH", "java"),
		NFeXSLTDir:                     strings.TrimSpace(os.Getenv("NFE_XSLT_DIR")),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("PG_URL is required")
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
