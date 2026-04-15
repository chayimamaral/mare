package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/repository"
)

type SerproServicoEnquadramentoService struct {
	repo *repository.SerproServicoEnquadramentoRepository
}

type SerproServicoEnquadramentoSaveInput struct {
	EnquadramentoID    string   `json:"enquadramento_id"`
	RegimeTributarioID string   `json:"regime_tributario_id"`
	ServicosIDs        []string `json:"servicos_ids"`
}

func NewSerproServicoEnquadramentoService(repo *repository.SerproServicoEnquadramentoRepository) *SerproServicoEnquadramentoService {
	return &SerproServicoEnquadramentoService{repo: repo}
}

func (s *SerproServicoEnquadramentoService) ListServicosIDs(ctx context.Context, enquadramentoID, regimeTributarioID string) ([]string, error) {
	eid := strings.TrimSpace(enquadramentoID)
	rid := strings.TrimSpace(regimeTributarioID)
	if eid == "" || rid == "" {
		return []string{}, nil
	}
	return s.repo.ListServicosIDs(ctx, eid, rid)
}

func (s *SerproServicoEnquadramentoService) SaveServicosIDs(ctx context.Context, in SerproServicoEnquadramentoSaveInput) error {
	eid := strings.TrimSpace(in.EnquadramentoID)
	rid := strings.TrimSpace(in.RegimeTributarioID)
	if eid == "" {
		return fmt.Errorf("enquadramento_id obrigatorio")
	}
	if rid == "" {
		return fmt.Errorf("regime_tributario_id obrigatorio")
	}

	seen := make(map[string]struct{}, len(in.ServicosIDs))
	ids := make([]string, 0, len(in.ServicosIDs))
	for _, raw := range in.ServicosIDs {
		id := strings.TrimSpace(raw)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	return s.repo.SaveServicosIDs(ctx, eid, rid, ids)
}
