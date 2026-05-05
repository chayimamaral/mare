package service

import (
	"context"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/chayimamaral/vecx/backend/internal/repository"
)

type MatrizConfiguracaoTributariaService struct {
	repo *repository.MatrizConfiguracaoTributariaRepository
}

type MatrizConfiguracaoTributariaListResponse struct {
	Items        []domain.MatrizConfiguracaoTributaria `json:"items"`
	TotalRecords int64                                 `json:"totalRecords"`
}

type MatrizConfiguracaoTributariaMutationResponse struct {
	Item domain.MatrizConfiguracaoTributaria `json:"item"`
}

func NewMatrizConfiguracaoTributariaService(repo *repository.MatrizConfiguracaoTributariaRepository) *MatrizConfiguracaoTributariaService {
	return &MatrizConfiguracaoTributariaService{repo: repo}
}

func (s *MatrizConfiguracaoTributariaService) List(ctx context.Context, params repository.MatrizConfiguracaoTributariaListParams) (MatrizConfiguracaoTributariaListResponse, error) {
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return MatrizConfiguracaoTributariaListResponse{}, err
	}
	return MatrizConfiguracaoTributariaListResponse{Items: items, TotalRecords: total}, nil
}

func (s *MatrizConfiguracaoTributariaService) Create(ctx context.Context, nome, naturezaJuridicaID, enquadramentoPorteID, regimeTributarioID string, aliquotaBase float64, possuiFatorR bool, aliquotaFatorR float64, substituicaoTributaria bool) (MatrizConfiguracaoTributariaMutationResponse, error) {
	item, err := s.repo.Create(ctx, nome, naturezaJuridicaID, enquadramentoPorteID, regimeTributarioID, aliquotaBase, possuiFatorR, aliquotaFatorR, substituicaoTributaria)
	if err != nil {
		return MatrizConfiguracaoTributariaMutationResponse{}, err
	}
	return MatrizConfiguracaoTributariaMutationResponse{Item: item}, nil
}

func (s *MatrizConfiguracaoTributariaService) Update(ctx context.Context, id, nome, naturezaJuridicaID, enquadramentoPorteID, regimeTributarioID string, aliquotaBase float64, possuiFatorR bool, aliquotaFatorR float64, substituicaoTributaria bool, ativo bool) (MatrizConfiguracaoTributariaMutationResponse, error) {
	item, err := s.repo.Update(ctx, id, nome, naturezaJuridicaID, enquadramentoPorteID, regimeTributarioID, aliquotaBase, possuiFatorR, aliquotaFatorR, substituicaoTributaria, ativo)
	if err != nil {
		return MatrizConfiguracaoTributariaMutationResponse{}, err
	}
	return MatrizConfiguracaoTributariaMutationResponse{Item: item}, nil
}

func (s *MatrizConfiguracaoTributariaService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *MatrizConfiguracaoTributariaService) GetByID(ctx context.Context, id string) (MatrizConfiguracaoTributariaMutationResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return MatrizConfiguracaoTributariaMutationResponse{}, err
	}
	return MatrizConfiguracaoTributariaMutationResponse{Item: item}, nil
}
