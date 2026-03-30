package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr            string
	DatabaseURL     string
	PublicAPIKey    string
	SSLRootCertPath string
	SSLInsecure     bool
}

// Load lê .env no diretório de trabalho e variáveis de ambiente.
func Load() (Config, error) {
	_ = godotenv.Load()
	_ = godotenv.Load(".env")

	dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL é obrigatória")
	}

	key := strings.TrimSpace(os.Getenv("PUBLIC_API_KEY"))
	if key == "" {
		return Config{}, fmt.Errorf("PUBLIC_API_KEY é obrigatória")
	}

	addr := strings.TrimSpace(os.Getenv("SERVER_ADDR"))
	if addr == "" {
		addr = ":8081"
	}
	if !strings.HasPrefix(addr, ":") && !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	cfg := Config{
		Addr:            addr,
		DatabaseURL:     dbURL,
		PublicAPIKey:    key,
		SSLRootCertPath: strings.TrimSpace(getEnv("PG_SSL_ROOT_CERT", "")),
		SSLInsecure:     getEnv("PG_SSL_INSECURE", "true") == "true",
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
