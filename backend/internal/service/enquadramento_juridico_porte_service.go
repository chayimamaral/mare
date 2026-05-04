package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/chayimamaral/vecx/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

type EnquadramentoJuridicoPorteService struct {
	repo *repository.EnquadramentoJuridicoPorteRepository
}

type EnquadramentoJuridicoPorteListResponse struct {
	Items []domain.EnquadramentoJuridicoPorte `json:"items"`
}

func NewEnquadramentoJuridicoPorteService(repo *repository.EnquadramentoJuridicoPorteRepository) *EnquadramentoJuridicoPorteService {
	return &EnquadramentoJuridicoPorteService{repo: repo}
}

func (s *EnquadramentoJuridicoPorteService) List(ctx context.Context, anoVigencia *int) (EnquadramentoJuridicoPorteListResponse, error) {
	items, err := s.repo.List(ctx, anoVigencia)
	if err != nil {
		return EnquadramentoJuridicoPorteListResponse{}, err
	}
	return EnquadramentoJuridicoPorteListResponse{Items: items}, nil
}

func (s *EnquadramentoJuridicoPorteService) Create(ctx context.Context, sigla, descricao string, limiteInicial float64, limiteFinal *float64, anoVigencia int) (domain.EnquadramentoJuridicoPorte, error) {
	if err := validateEnquadramentoJuridicoPorte(sigla, descricao, limiteInicial, limiteFinal, anoVigencia); err != nil {
		return domain.EnquadramentoJuridicoPorte{}, err
	}
	item, err := s.repo.Create(ctx, strings.TrimSpace(sigla), strings.TrimSpace(descricao), limiteInicial, limiteFinal, anoVigencia)
	if err != nil {
		return domain.EnquadramentoJuridicoPorte{}, mapEnquadramentoJuridicoPorteErr(err)
	}
	return item, nil
}

func (s *EnquadramentoJuridicoPorteService) Update(ctx context.Context, id, sigla, descricao string, limiteInicial float64, limiteFinal *float64, anoVigencia int, ativo bool) (domain.EnquadramentoJuridicoPorte, error) {
	if strings.TrimSpace(id) == "" {
		return domain.EnquadramentoJuridicoPorte{}, fmt.Errorf("id obrigatorio")
	}
	if err := validateEnquadramentoJuridicoPorte(sigla, descricao, limiteInicial, limiteFinal, anoVigencia); err != nil {
		return domain.EnquadramentoJuridicoPorte{}, err
	}
	item, err := s.repo.Update(ctx, strings.TrimSpace(id), strings.TrimSpace(sigla), strings.TrimSpace(descricao), limiteInicial, limiteFinal, anoVigencia, ativo)
	if err != nil {
		return domain.EnquadramentoJuridicoPorte{}, mapEnquadramentoJuridicoPorteErr(err)
	}
	return item, nil
}

func (s *EnquadramentoJuridicoPorteService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id obrigatorio")
	}
	return s.repo.Delete(ctx, strings.TrimSpace(id))
}

func validateEnquadramentoJuridicoPorte(sigla, descricao string, limiteInicial float64, limiteFinal *float64, anoVigencia int) error {
	if strings.TrimSpace(sigla) == "" {
		return fmt.Errorf("sigla obrigatoria")
	}
	if strings.TrimSpace(descricao) == "" {
		return fmt.Errorf("descricao obrigatoria")
	}
	if anoVigencia < 1995 || anoVigencia > 2100 {
		return fmt.Errorf("ano de vigencia invalido")
	}
	if limiteFinal != nil && *limiteFinal < limiteInicial {
		return fmt.Errorf("limite final nao pode ser menor que limite inicial")
	}
	return nil
}

func mapEnquadramentoJuridicoPorteErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		if strings.Contains(pgErr.ConstraintName, "enquadramento_juridico_porte_sigla_ano") {
			return fmt.Errorf("ja existe registro com esta sigla para o mesmo ano de vigencia")
		}
	}
	return err
}
