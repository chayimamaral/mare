package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
)

type TenantService struct {
	repo *repository.TenantRepository
}

type TenantCreatedResponse struct {
	TenantCreated domain.TenantEntity `json:"tenantCreated"`
}

type TenantDetailResponse struct {
	Tenant domain.TenantEntity `json:"tenant"`
}

func NewTenantService(repo *repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func normalizePlanoTenant(plano string) (string, error) {
	p := strings.ToUpper(strings.TrimSpace(plano))
	if p == "" {
		return "DEMO", nil
	}
	switch p {
	case "DEMO", "BASICO", "PRO", "PREMIUM":
		return p, nil
	default:
		return "", fmt.Errorf("Plano invalido: use DEMO, BASICO, PRO ou PREMIUM")
	}
}

func (s *TenantService) Create(ctx context.Context, nome, contato, plano, representativeID string) (TenantCreatedResponse, error) {
	p, err := normalizePlanoTenant(plano)
	if err != nil {
		return TenantCreatedResponse{}, err
	}
	tenant, err := s.repo.Create(ctx, nome, contato, p, strings.TrimSpace(representativeID))
	if err != nil {
		return TenantCreatedResponse{}, err
	}

	return TenantCreatedResponse{TenantCreated: tenant}, nil
}

func (s *TenantService) Detail(ctx context.Context, id string) (TenantDetailResponse, error) {
	tenant, err := s.repo.Detail(ctx, id)
	if err != nil {
		return TenantDetailResponse{}, err
	}

	return TenantDetailResponse{Tenant: tenant}, nil
}

func (s *TenantService) Update(ctx context.Context, id, nome, contato, plano string, active bool, representativeID *string) (domain.TenantEntity, error) {
	if representativeID != nil && strings.TrimSpace(*representativeID) != "" {
		cur, err := s.repo.Detail(ctx, id)
		if err != nil {
			return domain.TenantEntity{}, err
		}
		if cur.IsVecMaster {
			return domain.TenantEntity{}, fmt.Errorf("Tenant master da plataforma nao pode ter representante comercial")
		}
	}

	p := strings.TrimSpace(plano)
	if p != "" {
		normalized, err := normalizePlanoTenant(p)
		if err != nil {
			return domain.TenantEntity{}, err
		}
		p = normalized
	}

	tenant, err := s.repo.Update(ctx, id, nome, contato, p, active, representativeID)
	if err != nil {
		return domain.TenantEntity{}, err
	}

	return tenant, nil
}

func (s *TenantService) List(ctx context.Context, role, tenantID, representanteID string) (any, error) {
	role = strings.ToUpper(strings.TrimSpace(role))
	tenantID = strings.TrimSpace(tenantID)
	representanteID = strings.TrimSpace(representanteID)

	if role == "REPRESENTANTE" {
		if representanteID == "" {
			return []domain.TenantListRow{}, nil
		}
		return s.repo.ListForRepresentante(ctx, representanteID)
	}

	if role != "SUPER" && tenantID == "" {
		return []domain.TenantEntity{}, nil
	}

	if role == "SUPER" {
		return s.repo.ListWithDadosForSuper(ctx)
	}

	tenants, err := s.repo.List(ctx, role, tenantID)
	if err != nil {
		return nil, err
	}

	if len(tenants) == 0 {
		return []domain.TenantEntity{}, nil
	}

	return tenants[0], nil
}
