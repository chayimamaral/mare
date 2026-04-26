package service

import (
	"context"
	"errors"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
)

type RepresentanteService struct {
	repo *repository.RepresentanteRepository
}

func NewRepresentanteService(repo *repository.RepresentanteRepository) *RepresentanteService {
	return &RepresentanteService{repo: repo}
}

func (s *RepresentanteService) List(ctx context.Context) ([]domain.Representante, error) {
	return s.repo.List(ctx)
}

func (s *RepresentanteService) Create(ctx context.Context, nome, emailContato string) (domain.Representante, error) {
	if strings.TrimSpace(nome) == "" {
		return domain.Representante{}, errors.New("Nome obrigatorio")
	}
	return s.repo.Create(ctx, nome, emailContato)
}

func (s *RepresentanteService) Update(ctx context.Context, id, nome, emailContato string, ativo bool) (domain.Representante, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(nome) == "" {
		return domain.Representante{}, errors.New("Dados incompletos")
	}
	return s.repo.Update(ctx, id, nome, emailContato, ativo)
}

func (s *RepresentanteService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *RepresentanteService) ListModulos(ctx context.Context) ([]domain.ModuloPlataforma, error) {
	return s.repo.ListModulos(ctx)
}

func (s *RepresentanteService) GetMatriz(ctx context.Context, representanteID string) ([]domain.MatrizAcessoItem, error) {
	if strings.TrimSpace(representanteID) == "" {
		return nil, errors.New("representante_id obrigatorio")
	}
	return s.repo.GetMatriz(ctx, representanteID)
}

func (s *RepresentanteService) ReplaceMatriz(ctx context.Context, representanteID string, entries []domain.MatrizAcessoItem) error {
	if strings.TrimSpace(representanteID) == "" {
		return errors.New("representante_id obrigatorio")
	}
	return s.repo.ReplaceMatriz(ctx, representanteID, entries)
}
