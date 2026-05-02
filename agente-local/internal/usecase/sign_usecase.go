package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chayimamaral/vecx/agente-local/internal/a1sign"
	"github.com/chayimamaral/vecx/agente-local/internal/certlayout"
	"github.com/chayimamaral/vecx/agente-local/internal/certtax"
	"github.com/chayimamaral/vecx/agente-local/internal/domain"
	"github.com/chayimamaral/vecx/agente-local/internal/provider/pkcs11"
	"github.com/chayimamaral/vecx/agente-local/internal/settings"
)

type SignUseCase struct {
	pkcs  *pkcs11.Provider
	store *settings.Store
}

func NewSignUseCase(p *pkcs11.Provider, store *settings.Store) *SignUseCase {
	return &SignUseCase{pkcs: p, store: store}
}

func (uc *SignUseCase) ListCertificates(ctx context.Context) ([]domain.Certificate, error) {
	return uc.pkcs.ListCertificates(ctx)
}

func (uc *SignUseCase) Sign(ctx context.Context, input domain.SignInput) (domain.SignResult, error) {
	if input.HashSHA256Base64 == "" {
		return domain.SignResult{}, errors.New("hash_sha256_base64 obrigatorio")
	}

	hashBytes, err := base64.StdEncoding.DecodeString(input.HashSHA256Base64)
	if err != nil {
		return domain.SignResult{}, fmt.Errorf("hash_sha256_base64 invalido: %w", err)
	}
	if len(hashBytes) != 32 {
		return domain.SignResult{}, errors.New("hash SHA-256 deve ter 32 bytes")
	}

	taxClient := certtax.NormalizeTaxID(input.TaxID)
	structured := taxClient != "" || input.Procuracao

	if !structured {
		return uc.signLegacy(ctx, input, hashBytes)
	}

	if input.Procuracao && taxClient == "" {
		return domain.SignResult{}, errors.New("tax_id obrigatorio (CNPJ/CPF do cliente) quando procuracao=true")
	}

	st := uc.loadSettings()
	if st.PreferA3 {
		return uc.signA3Structured(ctx, input, hashBytes, taxClient)
	}

	if st.CertRootDir != "" {
		var pfxPath string
		var resErr error
		if input.Procuracao {
			pfxPath, resErr = certlayout.ResolveContadorPFX(st.CertRootDir)
		} else {
			pfxPath, resErr = certlayout.ResolveClientePFX(st.CertRootDir, taxClient)
		}
		if resErr == nil {
			if _, statErr := os.Stat(pfxPath); statErr == nil {
				return uc.signA1File(ctx, pfxPath, input.PIN, hashBytes)
			}
		}
	}

	return uc.signA3Structured(ctx, input, hashBytes, taxClient)
}

func (uc *SignUseCase) loadSettings() settings.AgentSettings {
	if uc.store == nil {
		return settings.AgentSettings{}
	}
	return uc.store.Load()
}

func (uc *SignUseCase) signA1File(ctx context.Context, pfxPath, password string, hash []byte) (domain.SignResult, error) {
	select {
	case <-ctx.Done():
		return domain.SignResult{}, ctx.Err()
	default:
	}
	sig, cert, err := a1sign.SignSHA256PKCS1v15(pfxPath, password, hash)
	if err != nil {
		return domain.SignResult{}, err
	}
	label := pfxPath
	if cert != nil {
		if cn := strings.TrimSpace(cert.Subject.CommonName); cn != "" {
			label = cn
		}
	}
	return domain.SignResult{
		Algorithm:        "RSA-SHA256",
		SignatureBase64:  base64.StdEncoding.EncodeToString(sig),
		SelectedCertID:   "a1:" + pfxPath,
		SelectedCertName: label,
	}, nil
}

func (uc *SignUseCase) signA3Structured(ctx context.Context, input domain.SignInput, hash []byte, taxClient string) (domain.SignResult, error) {
	certs, err := uc.pkcs.ListCertificates(ctx)
	if err != nil {
		return domain.SignResult{}, err
	}
	if len(certs) == 0 {
		return domain.SignResult{}, errors.New("nenhum certificado A3 encontrado no token")
	}

	if strings.TrimSpace(input.CertificateID) != "" {
		return uc.finishSign(ctx, input.CertificateID, hash, input.PIN, certs)
	}

	want := taxClient
	if input.Procuracao {
		want = certtax.NormalizeTaxID(input.SignerTaxID)
	}

	if want != "" {
		var match []domain.Certificate
		for _, c := range certs {
			for _, tid := range c.TaxIDs {
				if certtax.NormalizeTaxID(tid) == want {
					match = append(match, c)
					break
				}
			}
		}
		if len(match) == 1 {
			return uc.finishSign(ctx, match[0].ID, hash, input.PIN, certs)
		}
		if len(match) > 1 {
			return domain.SignResult{}, errors.New("varios certificados A3 correspondem ao titular; informe certificate_id")
		}
		return domain.SignResult{}, fmt.Errorf("nenhum certificado A3 corresponde ao titular %s", want)
	}

	if input.Procuracao && len(certs) == 1 {
		return uc.finishSign(ctx, certs[0].ID, hash, input.PIN, certs)
	}

	if len(certs) == 1 {
		return uc.finishSign(ctx, certs[0].ID, hash, input.PIN, certs)
	}

	return domain.SignResult{}, errors.New("varios certificados no token; informe certificate_id ou signer_tax_id (procuracao)")
}

func (uc *SignUseCase) finishSign(ctx context.Context, certID string, hash []byte, pin string, certs []domain.Certificate) (domain.SignResult, error) {
	signature, err := uc.pkcs.SignSHA256(ctx, certID, hash, pin)
	if err != nil {
		return domain.SignResult{}, err
	}
	selectedName := certID
	for _, c := range certs {
		if c.ID == certID {
			selectedName = c.Label
			break
		}
	}
	return domain.SignResult{
		Algorithm:        "RSA-SHA256",
		SignatureBase64:  base64.StdEncoding.EncodeToString(signature),
		SelectedCertID:   certID,
		SelectedCertName: selectedName,
	}, nil
}

func (uc *SignUseCase) signLegacy(ctx context.Context, input domain.SignInput, hashBytes []byte) (domain.SignResult, error) {
	certs, err := uc.pkcs.ListCertificates(ctx)
	if err != nil {
		return domain.SignResult{}, err
	}
	if len(certs) == 0 {
		return domain.SignResult{}, errors.New("nenhum certificado encontrado no token/smartcard")
	}

	selected := input.CertificateID
	if selected == "" {
		if len(certs) > 1 {
			return domain.SignResult{}, errors.New("mais de um certificado encontrado; informe certificate_id")
		}
		selected = certs[0].ID
	}

	signature, err := uc.pkcs.SignSHA256(ctx, selected, hashBytes, input.PIN)
	if err != nil {
		return domain.SignResult{}, err
	}

	selectedName := selected
	for _, c := range certs {
		if c.ID == selected {
			selectedName = c.Label
			break
		}
	}

	return domain.SignResult{
		Algorithm:        "RSA-SHA256",
		SignatureBase64:  base64.StdEncoding.EncodeToString(signature),
		SelectedCertID:   selected,
		SelectedCertName: selectedName,
	}, nil
}
