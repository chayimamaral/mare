package service

import (
	"context"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
)

type IntegraContadorServicoProcuracaoService struct {
	repo *repository.IntegraContadorServicoProcuracaoRepository
}

func NewIntegraContadorServicoProcuracaoService(repo *repository.IntegraContadorServicoProcuracaoRepository) *IntegraContadorServicoProcuracaoService {
	return &IntegraContadorServicoProcuracaoService{repo: repo}
}

func (s *IntegraContadorServicoProcuracaoService) List(ctx context.Context, idSistema, idServico string) ([]domain.IntegraContadorServicoProcuracao, error) {
	return s.repo.List(ctx, strings.TrimSpace(idSistema), strings.TrimSpace(idServico))
}
