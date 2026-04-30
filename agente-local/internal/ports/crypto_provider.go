package ports

import (
	"context"

	"github.com/chayimamaral/vecx/agente-local/internal/domain"
)

type CryptoProvider interface {
	ListCertificates(ctx context.Context) ([]domain.Certificate, error)
	SignSHA256(ctx context.Context, certificateID string, hashSHA256 []byte, pin string) ([]byte, error)
}
