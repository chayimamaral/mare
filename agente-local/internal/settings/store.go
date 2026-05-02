package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const settingsFile = "settings.json"

// AgentSettings persistido localmente (EF-937).
type AgentSettings struct {
	CertRootDir string `json:"cert_root_dir"`
	PreferA3    bool   `json:"prefer_a3"`
}

// Store lê/grava JSON em UserConfigDir/vecx-agent/settings.json.
type Store struct {
	path string
	mu   sync.RWMutex
}

// DefaultStore cria o diretório de configuração se necessário.
func DefaultStore() (*Store, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(base, "vecx-agent")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return &Store{path: filepath.Join(dir, settingsFile)}, nil
}

// Load retorna configuração; arquivo ausente ⇒ valores zero.
func (s *Store) Load() AgentSettings {
	if s == nil {
		return AgentSettings{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	raw, err := os.ReadFile(s.path)
	if err != nil || len(raw) == 0 {
		return AgentSettings{}
	}
	var out AgentSettings
	if json.Unmarshal(raw, &out) != nil {
		return AgentSettings{}
	}
	out.CertRootDir = cleanRootDir(out.CertRootDir)
	return out
}

// Save grava atomically (write temp + rename).
func (s *Store) Save(in AgentSettings) error {
	if s == nil {
		return nil
	}
	in.CertRootDir = cleanRootDir(in.CertRootDir)
	s.mu.Lock()
	defer s.mu.Unlock()
	raw, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, ".vecx-agent-settings-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(raw); err != nil {
		tmp.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, s.path)
}

func cleanRootDir(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return filepath.Clean(s)
}
