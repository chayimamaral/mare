package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/auth"
	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// ErrTenantNaoAutorizadoVecX: tenant inativo no cadastro; login bloqueado com aviso VECX (EF-929).
var ErrTenantNaoAutorizadoVecX = errors.New("tenant_vecx_nao_autorizado")

// MsgTenantVecxNaoAutorizado mensagem exibida ao usuario (dialog vermelho).
const MsgTenantVecxNaoAutorizado = "Esta utilizacao nao esta autorizada. Todos os dados estao salvos. Entre em contato diretamente com a VECX para regularizar a situacao."

// ErrAuditoriaGlobalIndisponivel: falha ao gravar sessao no VECX_AUDIT; login nao pode concluir (EF-929).
var ErrAuditoriaGlobalIndisponivel = errors.New("auditoria_global_indisponivel")

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Tenant   string `json:"tenant"`
}

type AuthService struct {
	users    *repository.UserRepository
	features *repository.FeatureMatrixRepository
	tokens   *auth.TokenService
	audit    *repository.VecxAuditRepository
}

type LoginResponse struct {
	ID           string        `json:"id"`
	Nome         string        `json:"nome"`
	Email        string        `json:"email"`
	TenantID     string        `json:"tenantid"`
	Token        string        `json:"token"`
	Tenant       domain.Tenant `json:"tenant"`
	Role         string        `json:"role"`
	FeatureSlugs []string      `json:"feature_slugs,omitempty"`
}

func NewAuthService(users *repository.UserRepository, features *repository.FeatureMatrixRepository, tokens *auth.TokenService, audit *repository.VecxAuditRepository) *AuthService {
	return &AuthService{users: users, features: features, tokens: tokens, audit: audit}
}

func contactoEscritorioAudit(tenantContato, telDados string) string {
	t := strings.TrimSpace(telDados)
	c := strings.TrimSpace(tenantContato)
	if t != "" && c != "" && t != c {
		return t + " | " + c
	}
	if t != "" {
		return t
	}
	return c
}

func (s *AuthService) syncAuditSession(ctx context.Context, userID, userEmail string, tenant domain.Tenant, clientIP string) error {
	cnpj, tel := s.users.LoadTenantDadosForAudit(ctx, tenant.ID)
	cont := contactoEscritorioAudit(tenant.Contato, tel)
	if err := s.audit.UpsertActiveSession(ctx, userID, userEmail, tenant.ID, tenant.Nome, cnpj, cont, clientIP); err != nil {
		return fmt.Errorf("%w: %v", ErrAuditoriaGlobalIndisponivel, err)
	}
	return nil
}

func jwtFeatureSlugs(slugs []string) []string {
	if slugs == nil {
		return []string{}
	}
	return slugs
}

func (s *AuthService) Login(ctx context.Context, input LoginInput, clientIP string) (LoginResponse, error) {
	input.Email = strings.TrimSpace(input.Email)
	user, err := s.users.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrLoginEmailNaoEncontrado) {
			return LoginResponse{}, errors.New("Email/password/empresa incorretos...")
		}
		return LoginResponse{}, fmt.Errorf(
			"Falha ao consultar usuario (verifique migrations/019_representantes_matriz_acesso.sql, nao confundir com 019_cliente_pf_pj_cliente_socios.sql). %w",
			err,
		)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		if user.Password != input.Password {
			return LoginResponse{}, errors.New("Email/password incorretos ...")
		}

		// Legacy migration path: old records may still store plaintext password.
		passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte(input.Password), 8)
		if hashErr != nil {
			return LoginResponse{}, hashErr
		}

		if updateErr := s.users.UpdatePassword(ctx, user.ID, string(passwordHash)); updateErr != nil {
			return LoginResponse{}, updateErr
		}
	}

	if !user.Tenant.Active {
		return LoginResponse{}, fmt.Errorf("%w: %s", ErrTenantNaoAutorizadoVecX, MsgTenantVecxNaoAutorizado)
	}

	roleUpper := strings.ToUpper(strings.TrimSpace(user.Role))
	if roleUpper == "REPRESENTANTE" {
		if strings.TrimSpace(user.RepresentanteID) == "" {
			return LoginResponse{}, errors.New("Usuario representante sem representante_id")
		}
		if !user.Tenant.IsVecMaster {
			ok, err := s.features.TenantLinkedToRepresentante(ctx, user.TenantID, user.RepresentanteID)
			if err != nil {
				return LoginResponse{}, err
			}
			if !ok {
				return LoginResponse{}, errors.New("Usuario representante sem vínculo valido com o tenant de login")
			}
		}
	}

	featureSlugs, err := s.features.ResolveForUser(ctx, roleUpper, user.TenantID, user.RepresentanteID)
	if err != nil {
		return LoginResponse{}, err
	}

	if err := s.syncAuditSession(ctx, user.ID, user.Email, user.Tenant, clientIP); err != nil {
		return LoginResponse{}, err
	}

	repID := strings.TrimSpace(user.RepresentanteID)
	token, err := s.tokens.Generate(auth.Claims{
		UserID:           user.ID,
		Nome:             user.Nome,
		Email:            user.Email,
		Tenant:           auth.TenantClaimsFromDomain(user.Tenant),
		Role:             user.Role,
		RepresentativeID: repID,
		FeatureSlugs:     jwtFeatureSlugs(featureSlugs),
	})
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		ID:           user.ID,
		Nome:         user.Nome,
		Email:        user.Email,
		TenantID:     user.TenantID,
		Token:        token,
		Tenant:       user.Tenant,
		Role:         user.Role,
		FeatureSlugs: featureSlugs,
	}, nil
}

type AssumeTenantInput struct {
	TenantID string `json:"tenant_id"`
}

type AssumeTenantResponse struct {
	Token        string        `json:"token"`
	Tenant       domain.Tenant `json:"tenant"`
	Role         string        `json:"role"`
	FeatureSlugs []string      `json:"feature_slugs,omitempty"`
}

func (s *AuthService) AssumeTenant(ctx context.Context, userID, role, currentTenantID, repID, clientIP string, input AssumeTenantInput) (AssumeTenantResponse, error) {
	target := strings.TrimSpace(input.TenantID)
	if target == "" {
		return AssumeTenantResponse{}, errors.New("tenant_id obrigatorio")
	}

	roleUpper := strings.ToUpper(strings.TrimSpace(role))
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return AssumeTenantResponse{}, err
	}

	switch roleUpper {
	case "SUPER":
		// SUPER pode assumir qualquer tenant ativo.
		t, err := s.users.LoadTenantForAssume(ctx, target)
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		if !t.Active {
			return AssumeTenantResponse{}, fmt.Errorf("%w: %s", ErrTenantNaoAutorizadoVecX, MsgTenantVecxNaoAutorizado)
		}
		slugs, err := s.features.ResolveForUser(ctx, roleUpper, t.ID, "")
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		if err := s.syncAuditSession(ctx, user.ID, user.Email, t, clientIP); err != nil {
			return AssumeTenantResponse{}, err
		}
		token, err := s.tokens.Generate(auth.Claims{
			UserID:       user.ID,
			Nome:         user.Nome,
			Email:        user.Email,
			Tenant:       auth.TenantClaimsFromDomain(t),
			Role:         user.Role,
			FeatureSlugs: jwtFeatureSlugs(slugs),
		})
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		return AssumeTenantResponse{Token: token, Tenant: t, Role: user.Role, FeatureSlugs: jwtFeatureSlugs(slugs)}, nil

	case "REPRESENTANTE":
		rid := strings.TrimSpace(user.RepresentanteID)
		if rid == "" {
			return AssumeTenantResponse{}, errors.New("Representante nao configurado para este usuario")
		}
		ok, err := s.features.TenantLinkedToRepresentante(ctx, target, rid)
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		if !ok {
			return AssumeTenantResponse{}, errors.New("Tenant nao vinculado a este representante")
		}
		t, err := s.users.LoadTenantForAssume(ctx, target)
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		if !t.Active {
			return AssumeTenantResponse{}, fmt.Errorf("%w: %s", ErrTenantNaoAutorizadoVecX, MsgTenantVecxNaoAutorizado)
		}
		slugs, err := s.features.ResolveForUser(ctx, roleUpper, t.ID, rid)
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		if err := s.syncAuditSession(ctx, user.ID, user.Email, t, clientIP); err != nil {
			return AssumeTenantResponse{}, err
		}
		token, err := s.tokens.Generate(auth.Claims{
			UserID:           user.ID,
			Nome:             user.Nome,
			Email:            user.Email,
			Tenant:           auth.TenantClaimsFromDomain(t),
			Role:             user.Role,
			RepresentativeID: rid,
			FeatureSlugs:     jwtFeatureSlugs(slugs),
		})
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		return AssumeTenantResponse{Token: token, Tenant: t, Role: user.Role, FeatureSlugs: jwtFeatureSlugs(slugs)}, nil

	default:
		if target != strings.TrimSpace(currentTenantID) {
			return AssumeTenantResponse{}, errors.New("Troca de tenant nao permitida para este perfil")
		}
		t, err := s.users.LoadTenantForAssume(ctx, target)
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		slugs, err := s.features.ResolveForUser(ctx, roleUpper, t.ID, user.RepresentanteID)
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		if err := s.syncAuditSession(ctx, user.ID, user.Email, t, clientIP); err != nil {
			return AssumeTenantResponse{}, err
		}
		token, err := s.tokens.Generate(auth.Claims{
			UserID:           user.ID,
			Nome:             user.Nome,
			Email:            user.Email,
			Tenant:           auth.TenantClaimsFromDomain(t),
			Role:             user.Role,
			RepresentativeID: strings.TrimSpace(user.RepresentanteID),
			FeatureSlugs:     jwtFeatureSlugs(slugs),
		})
		if err != nil {
			return AssumeTenantResponse{}, err
		}
		return AssumeTenantResponse{Token: token, Tenant: t, Role: user.Role, FeatureSlugs: jwtFeatureSlugs(slugs)}, nil
	}
}
