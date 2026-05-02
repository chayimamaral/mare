package pkcs11

import (
	"context"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/chayimamaral/vecx/agente-local/internal/certtax"
	"github.com/chayimamaral/vecx/agente-local/internal/domain"
	miekgpkcs11 "github.com/miekg/pkcs11"
)

type Provider struct {
	libraryPath string
}

func NewProvider(linuxLibrary, windowsLibrary string) *Provider {
	library := linuxLibrary
	if runtime.GOOS == "windows" {
		library = windowsLibrary
	}
	return &Provider{libraryPath: strings.TrimSpace(library)}
}

func (p *Provider) ListCertificates(ctx context.Context) ([]domain.Certificate, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	mod, slots, err := p.openWithSlots()
	if err != nil {
		return nil, err
	}
	defer mod.Destroy()
	defer mod.Finalize()

	certs := make([]domain.Certificate, 0, len(slots))
	for _, slotID := range slots {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		tokenInfo, err := mod.GetTokenInfo(slotID)
		if err != nil {
			continue
		}
		session, err := mod.OpenSession(slotID, miekgpkcs11.CKF_SERIAL_SESSION|miekgpkcs11.CKF_RW_SESSION)
		if err != nil {
			continue
		}
		_ = mod.FindObjectsInit(session, []*miekgpkcs11.Attribute{
			miekgpkcs11.NewAttribute(miekgpkcs11.CKA_CLASS, miekgpkcs11.CKO_CERTIFICATE),
		})
		objects, _, _ := mod.FindObjects(session, 128)
		_ = mod.FindObjectsFinal(session)

		tokenLabel := strings.TrimSpace(tokenInfo.Label)
		if tokenLabel == "" {
			tokenLabel = fmt.Sprintf("Token slot %d", slotID)
		}
		for _, obj := range objects {
			attrs, err := mod.GetAttributeValue(session, obj, []*miekgpkcs11.Attribute{
				miekgpkcs11.NewAttribute(miekgpkcs11.CKA_ID, nil),
				miekgpkcs11.NewAttribute(miekgpkcs11.CKA_LABEL, nil),
				miekgpkcs11.NewAttribute(miekgpkcs11.CKA_VALUE, nil),
			})
			if err != nil {
				continue
			}
			certID := attributeValue(attrs, miekgpkcs11.CKA_ID)
			label := strings.TrimSpace(string(attributeValue(attrs, miekgpkcs11.CKA_LABEL)))
			der := attributeValue(attrs, miekgpkcs11.CKA_VALUE)
			subject := ""
			serialHex := ""
			taxIDs := []string(nil)
			if len(der) > 0 {
				if parsed, err := x509.ParseCertificate(der); err == nil {
					subject = parsed.Subject.String()
					serialHex = strings.ToUpper(parsed.SerialNumber.Text(16))
					if label == "" {
						label = parsed.Subject.CommonName
					}
					taxIDs = certtax.TaxIDsFromCertificate(parsed)
				}
			}
			if label == "" {
				label = fmt.Sprintf("Certificado slot %d", slotID)
			}
			certs = append(certs, domain.Certificate{
				ID:         buildCertificateID(slotID, certID),
				Label:      label,
				Subject:    subject,
				SerialHex:  serialHex,
				SlotID:     slotID,
				TokenLabel: tokenLabel,
				TaxIDs:     taxIDs,
			})
		}
		_ = mod.CloseSession(session)
	}
	return certs, nil
}

func (p *Provider) SignSHA256(ctx context.Context, certificateID string, hashSHA256 []byte, pin string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if len(hashSHA256) != 32 {
		return nil, errors.New("hash SHA-256 invalido")
	}
	if strings.TrimSpace(certificateID) == "" {
		return nil, errors.New("certificate_id obrigatorio")
	}

	slotID, certObjID, err := parseCertificateID(certificateID)
	if err != nil {
		return nil, err
	}

	mod, _, err := p.openWithSlots()
	if err != nil {
		return nil, err
	}
	defer mod.Destroy()
	defer mod.Finalize()

	session, err := mod.OpenSession(slotID, miekgpkcs11.CKF_SERIAL_SESSION|miekgpkcs11.CKF_RW_SESSION)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir sessao PKCS#11: %w", err)
	}
	defer mod.CloseSession(session)

	if strings.TrimSpace(pin) != "" {
		if err := mod.Login(session, miekgpkcs11.CKU_USER, pin); err != nil && err.Error() != "CKR_USER_ALREADY_LOGGED_IN" {
			return nil, fmt.Errorf("falha no login PIN do token: %w", err)
		}
		defer mod.Logout(session)
	}

	privateKey, err := p.findPrivateKeyByCertID(mod, session, certObjID)
	if err != nil {
		return nil, err
	}

	mechanisms := []*miekgpkcs11.Mechanism{
		miekgpkcs11.NewMechanism(miekgpkcs11.CKM_SHA256_RSA_PKCS, nil),
	}
	if err := mod.SignInit(session, mechanisms, privateKey); err == nil {
		signature, err := mod.Sign(session, hashSHA256)
		if err == nil {
			return signature, nil
		}
	}

	digestInfo := append(sha256DigestInfoPrefix(), hashSHA256...)
	if err := mod.SignInit(session, []*miekgpkcs11.Mechanism{
		miekgpkcs11.NewMechanism(miekgpkcs11.CKM_RSA_PKCS, nil),
	}, privateKey); err != nil {
		return nil, fmt.Errorf("falha ao inicializar assinatura RSA no token: %w", err)
	}
	signature, err := mod.Sign(session, digestInfo)
	if err != nil {
		return nil, fmt.Errorf("falha ao assinar com PKCS#11: %w", err)
	}
	return signature, nil
}

func (p *Provider) openWithSlots() (*miekgpkcs11.Ctx, []uint, error) {
	if p.libraryPath == "" {
		return nil, nil, errors.New("biblioteca PKCS#11 nao configurada")
	}
	mod := miekgpkcs11.New(p.libraryPath)
	if mod == nil {
		return nil, nil, fmt.Errorf("falha ao carregar biblioteca PKCS#11: %s", p.libraryPath)
	}
	if err := mod.Initialize(); err != nil {
		return nil, nil, fmt.Errorf("erro ao inicializar PKCS#11 (%s): %w", p.libraryPath, err)
	}
	slots, err := mod.GetSlotList(true)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao listar slots PKCS#11: %w", err)
	}
	if len(slots) == 0 {
		return nil, nil, errors.New("nenhum slot/token PKCS#11 disponivel")
	}
	return mod, slots, nil
}

func (p *Provider) findPrivateKeyByCertID(mod *miekgpkcs11.Ctx, session miekgpkcs11.SessionHandle, certID []byte) (miekgpkcs11.ObjectHandle, error) {
	if err := mod.FindObjectsInit(session, []*miekgpkcs11.Attribute{
		miekgpkcs11.NewAttribute(miekgpkcs11.CKA_CLASS, miekgpkcs11.CKO_PRIVATE_KEY),
		miekgpkcs11.NewAttribute(miekgpkcs11.CKA_ID, certID),
	}); err != nil {
		return 0, fmt.Errorf("erro ao buscar chave privada no token: %w", err)
	}
	objs, _, err := mod.FindObjects(session, 1)
	_ = mod.FindObjectsFinal(session)
	if err != nil {
		return 0, fmt.Errorf("erro ao buscar chave privada no token: %w", err)
	}
	if len(objs) == 0 {
		return 0, errors.New("chave privada nao encontrada para o certificado selecionado")
	}
	return objs[0], nil
}

func buildCertificateID(slotID uint, certObjID []byte) string {
	return fmt.Sprintf("slot-%d-id-%s", slotID, strings.ToLower(hex.EncodeToString(certObjID)))
}

func parseCertificateID(value string) (uint, []byte, error) {
	parts := strings.Split(strings.TrimSpace(value), "-id-")
	if len(parts) != 2 {
		return 0, nil, errors.New("certificate_id invalido")
	}
	slotPart := strings.TrimPrefix(parts[0], "slot-")
	slotParsed, err := strconv.ParseUint(slotPart, 10, 32)
	if err != nil {
		return 0, nil, errors.New("certificate_id invalido")
	}
	certID, err := hex.DecodeString(parts[1])
	if err != nil || len(certID) == 0 {
		return 0, nil, errors.New("certificate_id invalido")
	}
	return uint(slotParsed), certID, nil
}

func attributeValue(attrs []*miekgpkcs11.Attribute, typ uint) []byte {
	for _, a := range attrs {
		if a.Type == typ {
			return a.Value
		}
	}
	return nil
}

func sha256DigestInfoPrefix() []byte {
	return []byte{
		0x30, 0x31, 0x30, 0x0d,
		0x06, 0x09, 0x60, 0x86,
		0x48, 0x01, 0x65, 0x03,
		0x04, 0x02, 0x01, 0x05,
		0x00, 0x04, 0x20,
	}
}
