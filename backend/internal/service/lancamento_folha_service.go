package service

import (
	"context"
	"time"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/chayimamaral/vecx/backend/internal/repository"
)

type LancamentoFolhaService struct {
	repo *repository.LancamentoFolhaRepository
}

type LancamentoFolhaTreeResponse struct {
	Tree []domain.LancamentoFolhaTreeNode `json:"tree"`
}

type LancamentoFolhaMutationResponse struct {
	Lancamento domain.LancamentoFolha `json:"lancamento"`
}

func NewLancamentoFolhaService(repo *repository.LancamentoFolhaRepository) *LancamentoFolhaService {
	return &LancamentoFolhaService{repo: repo}
}

func (s *LancamentoFolhaService) ListTree(ctx context.Context, tenantID string) (LancamentoFolhaTreeResponse, error) {
	tree, err := s.repo.ListClientesWithLancamentos(ctx, tenantID)
	if err != nil {
		return LancamentoFolhaTreeResponse{}, err
	}
	return LancamentoFolhaTreeResponse{Tree: tree}, nil
}

func (s *LancamentoFolhaService) Create(ctx context.Context, tenantID, clienteID string, competencia time.Time, valorFolha, valorFaturamento float64, observacoes string) (LancamentoFolhaMutationResponse, error) {
	lf, err := s.repo.Create(ctx, tenantID, clienteID, competencia, valorFolha, valorFaturamento, observacoes)
	if err != nil {
		return LancamentoFolhaMutationResponse{}, err
	}
	return LancamentoFolhaMutationResponse{Lancamento: lf}, nil
}

func (s *LancamentoFolhaService) Update(ctx context.Context, id, tenantID string, competencia time.Time, valorFolha, valorFaturamento float64, observacoes string) (LancamentoFolhaMutationResponse, error) {
	lf, err := s.repo.Update(ctx, id, tenantID, competencia, valorFolha, valorFaturamento, observacoes)
	if err != nil {
		return LancamentoFolhaMutationResponse{}, err
	}
	return LancamentoFolhaMutationResponse{Lancamento: lf}, nil
}

func (s *LancamentoFolhaService) Delete(ctx context.Context, id, tenantID string) error {
	return s.repo.Delete(ctx, id, tenantID)
}

func (s *LancamentoFolhaService) GetByID(ctx context.Context, id, tenantID string) (LancamentoFolhaMutationResponse, error) {
	lf, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return LancamentoFolhaMutationResponse{}, err
	}
	return LancamentoFolhaMutationResponse{Lancamento: lf}, nil
}
