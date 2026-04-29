package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/crypto/certseal"
	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/chayimamaral/vecx/backend/internal/repository"
	"golang.org/x/crypto/pkcs12"
)

var (
	// ErrChaveCriptografiaNaoConfigurada quando VECONTAB_CERT_CRYPTO_KEY_HEX está vazia ou inválida.
	ErrChaveCriptografiaNaoConfigurada = errors.New("criptografia de certificado nao configurada (VECONTAB_CERT_CRYPTO_KEY_HEX)")
)

// CertificadoService cifra PFX e senha antes do Postgres; decifra só em memória para uso em TLS/SERPRO.
type CertificadoService struct {
	repo        *repository.CertificadoRepository
	repoCliente *repository.CertificadoClienteRepository
	key         []byte
}

// NewCertificadoService aceita keyHex vazia (serviço inoperante até configurar).
// repoCliente pode ser nil (upload por cliente ficará indisponível).
func NewCertificadoService(repo *repository.CertificadoRepository, repoCliente *repository.CertificadoClienteRepository, keyHex string) (*CertificadoService, error) {
	keyHex = strings.TrimSpace(keyHex)
	if keyHex == "" {
		return &CertificadoService{repo: repo, repoCliente: repoCliente, key: nil}, nil
	}
	k, err := certseal.ParseKeyHex(keyHex)
	if err != nil {
		return nil, fmt.Errorf("VECONTAB_CERT_CRYPTO_KEY_HEX: %w", err)
	}
	return &CertificadoService{repo: repo, repoCliente: repoCliente, key: k}, nil
}

// Configurado indica se a chave AES-256 está disponível.
func (s *CertificadoService) Configurado() bool {
	return s != nil && len(s.key) == certseal.KeySize
}

// ClienteRepoConfigurado indica se persistência por cliente está disponível.
func (s *CertificadoService) ClienteRepoConfigurado() bool {
	return s != nil && s.repoCliente != nil
}

// UpsertPFX valida o PKCS#12, extrai validade (e metadados básicos) e persiste cifrado por tenant.
func (s *CertificadoService) UpsertPFX(ctx context.Context, tenantID string, pfx []byte, senhaPlana, cnpjHint, titularHint string) error {
	if s == nil || s.repo == nil {
		return fmt.Errorf("repositorio nao configurado")
	}
	if !s.Configurado() {
		return ErrChaveCriptografiaNaoConfigurada
	}
	tid := strings.TrimSpace(tenantID)
	if tid == "" {
		return fmt.Errorf("tenant obrigatorio")
	}
	senhaPlana = strings.TrimSpace(senhaPlana)
	if len(pfx) == 0 || senhaPlana == "" {
		return fmt.Errorf("pfx e senha obrigatorios")
	}

	priv, leaf, err := pkcs12.Decode(pfx, senhaPlana)
	if err != nil {
		return fmt.Errorf("pfx invalido ou senha incorreta: %w", err)
	}
	_ = priv

	pfxSealed, err := certseal.Seal(s.key, pfx)
	if err != nil {
		return fmt.Errorf("cifrar pfx: %w", err)
	}
	senhaSealed, err := certseal.Seal(s.key, []byte(senhaPlana))
	if err != nil {
		return fmt.Errorf("cifrar senha: %w", err)
	}

	cnpj := strings.TrimSpace(cnpjHint)
	nome := strings.TrimSpace(titularHint)
	if nome == "" {
		nome = strings.TrimSpace(leaf.Subject.CommonName)
	}

	row := &domain.Certificado{
		Tenant:       tid,
		PFXCifrado:   pfxSealed,
		SenhaCifrada: senhaSealed,
		CNPJ:         cnpj,
		TitularNome:  nome,
		EmitidoPor:   strings.TrimSpace(leaf.Issuer.CommonName),
		ValidadeDe:   leaf.NotBefore,
		ValidadeAte:  leaf.NotAfter,
		Ativo:        true,
	}
	return s.repo.UpsertAtivo(ctx, row)
}

// UpsertPFXCliente valida PKCS#12 e persiste em certificado_cliente no schema do tenant (1:1 com cliente_id).
func (s *CertificadoService) UpsertPFXCliente(ctx context.Context, tenantID, empresaID string, pfx []byte, senhaPlana, cnpjHint, titularHint string) error {
	if s == nil || s.repoCliente == nil {
		return fmt.Errorf("repositorio certificado_cliente nao configurado")
	}
	if !s.Configurado() {
		return ErrChaveCriptografiaNaoConfigurada
	}
	tid := strings.TrimSpace(tenantID)
	eid := strings.TrimSpace(empresaID)
	if tid == "" || eid == "" {
		return fmt.Errorf("tenant e empresa obrigatorios")
	}
	senhaPlana = strings.TrimSpace(senhaPlana)
	if len(pfx) == 0 || senhaPlana == "" {
		return fmt.Errorf("pfx e senha obrigatorios")
	}

	priv, leaf, err := pkcs12.Decode(pfx, senhaPlana)
	if err != nil {
		return fmt.Errorf("pfx invalido ou senha incorreta: %w", err)
	}
	_ = priv

	pfxSealed, err := certseal.Seal(s.key, pfx)
	if err != nil {
		return fmt.Errorf("cifrar pfx: %w", err)
	}
	senhaSealed, err := certseal.Seal(s.key, []byte(senhaPlana))
	if err != nil {
		return fmt.Errorf("cifrar senha: %w", err)
	}

	cnpj := strings.TrimSpace(cnpjHint)
	nome := strings.TrimSpace(titularHint)
	if nome == "" {
		nome = strings.TrimSpace(leaf.Subject.CommonName)
	}

	clienteID, err := s.repoCliente.ClienteIDEmpresaTenant(ctx, eid, tid)
	if err != nil {
		return err
	}

	return s.repoCliente.UpsertAtivo(ctx, clienteID, pfxSealed, senhaSealed, cnpj, nome,
		strings.TrimSpace(leaf.Issuer.CommonName), leaf.NotBefore, leaf.NotAfter)
}

// ResumoCertificadoCliente retorna metadados do certificado do cliente (por empresa_id da API).
func (s *CertificadoService) ResumoCertificadoCliente(ctx context.Context, tenantID, empresaID string) (*domain.CertificadoClienteResumo, error) {
	if s == nil || s.repoCliente == nil {
		return nil, fmt.Errorf("repositorio certificado_cliente nao configurado")
	}
	tid := strings.TrimSpace(tenantID)
	eid := strings.TrimSpace(empresaID)
	if tid == "" || eid == "" {
		return nil, fmt.Errorf("tenant e empresa obrigatorios")
	}
	clienteID, err := s.repoCliente.ClienteIDEmpresaTenant(ctx, eid, tid)
	if err != nil {
		return nil, err
	}
	row, err := s.repoCliente.GetResumoAtivo(ctx, clienteID)
	if err != nil || row == nil {
		return nil, err
	}
	return &domain.CertificadoClienteResumo{
		TipoCertificado: "A1",
		NomeCertificado: firstNonEmptyCert(row.EmitidoPor, row.TitularNome),
		EmitidoPara:     strings.TrimSpace(row.TitularNome),
		EmitidoPor:      strings.TrimSpace(row.TitularNome),
		CNPJ:            strings.TrimSpace(row.CNPJ),
		ValidadeDe:      row.ValidadeDe.Format("2006-01-02"),
		ValidadeAte:     row.ValidadeAte.Format("2006-01-02"),
	}, nil
}

func firstNonEmptyCert(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// MaterialEmMemoria decifra PFX e senha; o chamador deve chamar CertificadoMaterial.Zero() após uso.
func (s *CertificadoService) MaterialEmMemoria(ctx context.Context, tenantID string) (*domain.CertificadoMaterial, error) {
	if !s.Configurado() {
		return nil, ErrChaveCriptografiaNaoConfigurada
	}
	row, err := s.repo.GetAtivoPorTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	pfx, err := certseal.Open(s.key, row.PFXCifrado)
	if err != nil {
		return nil, fmt.Errorf("decifrar pfx: %w", err)
	}
	pwdBin, err := certseal.Open(s.key, row.SenhaCifrada)
	if err != nil {
		clearBytes(pfx)
		return nil, fmt.Errorf("decifrar senha: %w", err)
	}
	senha := string(pwdBin)
	clearBytes(pwdBin)

	out := &domain.CertificadoMaterial{
		PFX:         append([]byte(nil), pfx...),
		Senha:       senha,
		CNPJ:        row.CNPJ,
		Nome:        row.TitularNome,
		ValidadeAte: row.ValidadeAte,
	}
	clearBytes(pfx)
	return out, nil
}

// ResumoPorTenant retorna metadados persistidos do certificado ativo do tenant.
func (s *CertificadoService) ResumoPorTenant(ctx context.Context, tenantID string) (*domain.Certificado, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("repositorio nao configurado")
	}
	return s.repo.GetResumoAtivoPorTenant(ctx, tenantID)
}

// TLSClientCertificate monta tls.Certificate para mTLS (ex.: SERPRO) a partir do material decifrado.
func (s *CertificadoService) TLSClientCertificate(ctx context.Context, tenantID, empresaID string) (tls.Certificate, func(), error) {
	mat, err := s.MaterialEmMemoria(ctx, tenantID)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	tcert, err := tlsCertificateFromPFX(mat.PFX, mat.Senha)
	cleanup := func() { mat.Zero() }
	if err != nil {
		cleanup()
		return tls.Certificate{}, nil, err
	}
	return tcert, cleanup, nil
}

func tlsCertificateFromPFX(pfx []byte, password string) (tls.Certificate, error) {
	priv, cert, err := pkcs12.Decode(pfx, password)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("pkcs12 decode: %w", err)
	}
	if cert == nil {
		return tls.Certificate{}, fmt.Errorf("certificado folha ausente no pfx")
	}
	return tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  priv,
		Leaf:        cert,
	}, nil
}

func clearBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
