package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
)

type MonitorOperacaoService struct {
	repo *repository.MonitorOperacaoRepository
}

func NewMonitorOperacaoService(repo *repository.MonitorOperacaoRepository) *MonitorOperacaoService {
	return &MonitorOperacaoService{repo: repo}
}

type MonitorOperacaoListResponse struct {
	Itens []domain.MonitorOperacaoItem `json:"itens"`
	Total int64                        `json:"total"`
}

func (s *MonitorOperacaoService) Registrar(ctx context.Context, in repository.MonitorOperacaoInsert) error {
	o := strings.TrimSpace(in.Origem)
	t := strings.TrimSpace(in.Tipo)
	st := strings.TrimSpace(in.Status)
	if o == "" || t == "" || st == "" {
		return fmt.Errorf("origem, tipo e status sao obrigatorios")
	}
	if strings.TrimSpace(in.TenantID) == "" {
		return fmt.Errorf("tenant_id obrigatorio")
	}
	return s.repo.Insert(ctx, in)
}

func (s *MonitorOperacaoService) ListPage(ctx context.Context, viewerRole, viewerTenantID string, limit, offset int) (MonitorOperacaoListResponse, error) {
	role := strings.TrimSpace(strings.ToUpper(viewerRole))
	if role != "SUPER" && role != "ADMIN" {
		return MonitorOperacaoListResponse{}, fmt.Errorf("perfil nao autorizado")
	}
	if role == "ADMIN" && strings.TrimSpace(viewerTenantID) == "" {
		return MonitorOperacaoListResponse{}, fmt.Errorf("tenant nao identificado")
	}
	total, err := s.repo.CountList(ctx, role, viewerTenantID)
	if err != nil {
		return MonitorOperacaoListResponse{}, err
	}
	items, err := s.repo.ListPage(ctx, role, viewerTenantID, limit, offset)
	if err != nil {
		return MonitorOperacaoListResponse{}, err
	}
	return MonitorOperacaoListResponse{Itens: items, Total: total}, nil
}
