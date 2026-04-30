package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/chayimamaral/vecx/agente-local/internal/domain"
	"github.com/chayimamaral/vecx/agente-local/internal/ports"
)

type SignUseCase struct {
	provider ports.CryptoProvider
}

func NewSignUseCase(provider ports.CryptoProvider) *SignUseCase {
	return &SignUseCase{provider: provider}
}

func (uc *SignUseCase) ListCertificates(ctx context.Context) ([]domain.Certificate, error) {
	return uc.provider.ListCertificates(ctx)
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

	certs, err := uc.provider.ListCertificates(ctx)
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

	signature, err := uc.provider.SignSHA256(ctx, selected, hashBytes, input.PIN)
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
