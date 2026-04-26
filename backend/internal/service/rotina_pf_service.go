package service

import (
	"context"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
)

type RotinaPFService struct {
	repo *repository.RotinaPFRepository
}

type RotinaPFListResponse struct {
	RotinasPF    []domain.RotinaPFLiteItem `json:"rotinas_pf"`
	TotalRecords int64                     `json:"totalRecords"`
}

type RotinaPFAdminListResponse struct {
	RotinasPF    []domain.RotinaPFListRow `json:"rotinas_pf"`
	TotalRecords int64                    `json:"totalRecords"`
}

type RotinaPFMutationResponse struct {
	RotinasPF    []domain.RotinaPFListRow `json:"rotinas_pf"`
	TotalRecords int64                    `json:"totalRecords"`
}

type RotinaPFItensResponse struct {
	Itens        []domain.RotinaPFItemRow `json:"itens"`
	TotalRecords int64                    `json:"totalRecords"`
}

func NewRotinaPFService(repo *repository.RotinaPFRepository) *RotinaPFService {
	return &RotinaPFService{repo: repo}
}

func (s *RotinaPFService) ListLite(ctx context.Context, tenantID string) (RotinaPFListResponse, error) {
	items, total, err := s.repo.ListLite(ctx, tenantID)
	if err != nil {
		return RotinaPFListResponse{}, err
	}
	return RotinaPFListResponse{RotinasPF: items, TotalRecords: total}, nil
}

func (s *RotinaPFService) ListAdmin(ctx context.Context, params repository.RotinaPFListParams) (RotinaPFAdminListResponse, error) {
	rows, total, err := s.repo.List(ctx, params)
	if err != nil {
		return RotinaPFAdminListResponse{}, err
	}
	return RotinaPFAdminListResponse{RotinasPF: rows, TotalRecords: total}, nil
}

func (s *RotinaPFService) Create(ctx context.Context, in repository.RotinaPFUpsertInput) (RotinaPFMutationResponse, error) {
	rows, total, err := s.repo.Create(ctx, in)
	if err != nil {
		return RotinaPFMutationResponse{}, err
	}
	return RotinaPFMutationResponse{RotinasPF: rows, TotalRecords: total}, nil
}

func (s *RotinaPFService) Update(ctx context.Context, in repository.RotinaPFUpsertInput) (RotinaPFMutationResponse, error) {
	rows, total, err := s.repo.Update(ctx, in)
	if err != nil {
		return RotinaPFMutationResponse{}, err
	}
	return RotinaPFMutationResponse{RotinasPF: rows, TotalRecords: total}, nil
}

func (s *RotinaPFService) SoftDelete(ctx context.Context, id, tenantID string) (RotinaPFMutationResponse, error) {
	rows, total, err := s.repo.SoftDelete(ctx, id, tenantID)
	if err != nil {
		return RotinaPFMutationResponse{}, err
	}
	return RotinaPFMutationResponse{RotinasPF: rows, TotalRecords: total}, nil
}

func (s *RotinaPFService) ListItens(ctx context.Context, rotinaPFID, tenantID string) (RotinaPFItensResponse, error) {
	itens, total, err := s.repo.ListItens(ctx, rotinaPFID, tenantID)
	if err != nil {
		return RotinaPFItensResponse{}, err
	}
	return RotinaPFItensResponse{Itens: itens, TotalRecords: total}, nil
}

func (s *RotinaPFService) CreateItem(ctx context.Context, in repository.RotinaPFItemUpsertInput) (RotinaPFItensResponse, error) {
	itens, total, err := s.repo.CreateItem(ctx, in)
	if err != nil {
		return RotinaPFItensResponse{}, err
	}
	return RotinaPFItensResponse{Itens: itens, TotalRecords: total}, nil
}

func (s *RotinaPFService) UpdateItem(ctx context.Context, in repository.RotinaPFItemUpsertInput) (RotinaPFItensResponse, error) {
	itens, total, err := s.repo.UpdateItem(ctx, in)
	if err != nil {
		return RotinaPFItensResponse{}, err
	}
	return RotinaPFItensResponse{Itens: itens, TotalRecords: total}, nil
}

func (s *RotinaPFService) DeleteItem(ctx context.Context, itemID, rotinaPFID, tenantID string) (RotinaPFItensResponse, error) {
	itens, total, err := s.repo.DeleteItem(ctx, itemID, rotinaPFID, tenantID)
	if err != nil {
		return RotinaPFItensResponse{}, err
	}
	return RotinaPFItensResponse{Itens: itens, TotalRecords: total}, nil
}
