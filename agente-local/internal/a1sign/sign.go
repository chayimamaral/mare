package a1sign

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/pkcs12"
)

// SignSHA256PKCS1v15 abre um .pfx (PKCS#12), valida SHA-256 e devolve assinatura RSA PKCS#1 v1.5
// (compatível com CKM_SHA256_RSA_PKCS no token).
func SignSHA256PKCS1v15(pfxPath, password string, hashSHA256 []byte) ([]byte, *x509.Certificate, error) {
	if len(hashSHA256) != 32 {
		return nil, nil, errors.New("hash SHA-256 deve ter 32 bytes")
	}
	raw, err := os.ReadFile(pfxPath)
	if err != nil {
		return nil, nil, fmt.Errorf("ler pfx: %w", err)
	}
	priv, cert, err := pkcs12.Decode(raw, password)
	if err != nil {
		return nil, nil, fmt.Errorf("pfx invalido ou senha incorreta: %w", err)
	}
	rsaKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		if _, ecdsaOK := priv.(*ecdsa.PrivateKey); ecdsaOK {
			return nil, cert, errors.New("certificado A1 ECDSA nao suportado no agente (use RSA)")
		}
		return nil, cert, errors.New("chave privada A1 deve ser RSA")
	}
	// hashSHA256 já é o digest SHA-256 do payload (32 bytes).
	sig, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hashSHA256)
	if err != nil {
		return nil, cert, fmt.Errorf("assinar rsa: %w", err)
	}
	return sig, cert, nil
}
