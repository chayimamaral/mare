package config

import (
	"os"
	"strings"
)

type Config struct {
	HTTPAddr            string
	AllowedOrigins      []string
	PKCS11LibraryLinux  string
	PKCS11LibraryWindow string
	SharedSecret        string
}

func Load() Config {
	return Config{
		HTTPAddr:            envOrDefault("AGENT_HTTP_ADDR", "127.0.0.1:9999"),
		AllowedOrigins:      parseCSV(strings.TrimSpace(os.Getenv("AGENT_ALLOWED_ORIGINS"))),
		PKCS11LibraryLinux:  envOrDefault("PKCS11_LIBRARY_LINUX", "/usr/lib64/libeToken.so"),
		PKCS11LibraryWindow: envOrDefault("PKCS11_LIBRARY_WINDOWS", `C:\Windows\System32\eTPKCS11.dll`),
		SharedSecret:        strings.TrimSpace(os.Getenv("AGENT_SHARED_SECRET")),
	}
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
