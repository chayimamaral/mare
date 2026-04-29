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

type SalarioMinimoService struct {
	repo *repository.SalarioMinimoRepository
}

type SalarioMinimoListResponse struct {
	Salarios []domain.SalarioMinimoNacional `json:"salarios"`
}

func NewSalarioMinimoService(repo *repository.SalarioMinimoRepository) *SalarioMinimoService {
	return &SalarioMinimoService{repo: repo}
}

func (s *SalarioMinimoService) List(ctx context.Context) (SalarioMinimoListResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return SalarioMinimoListResponse{}, err
	}
	return SalarioMinimoListResponse{Salarios: items}, nil
}

func (s *SalarioMinimoService) Create(ctx context.Context, ano int, valor float64) (domain.SalarioMinimoNacional, error) {
	if err := validaSalarioMinimo(ano, valor); err != nil {
		return domain.SalarioMinimoNacional{}, err
	}
	item, err := s.repo.Create(ctx, ano, valor)
	if err != nil {
		return domain.SalarioMinimoNacional{}, mapSalarioMinimoErr(err)
	}
	return item, nil
}

func (s *SalarioMinimoService) Update(ctx context.Context, id string, ano int, valor float64) (domain.SalarioMinimoNacional, error) {
	if strings.TrimSpace(id) == "" {
		return domain.SalarioMinimoNacional{}, fmt.Errorf("id obrigatorio")
	}
	if err := validaSalarioMinimo(ano, valor); err != nil {
		return domain.SalarioMinimoNacional{}, err
	}
	item, err := s.repo.Update(ctx, strings.TrimSpace(id), ano, valor)
	if err != nil {
		return domain.SalarioMinimoNacional{}, mapSalarioMinimoErr(err)
	}
	return item, nil
}

func (s *SalarioMinimoService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id obrigatorio")
	}
	return s.repo.Delete(ctx, strings.TrimSpace(id))
}

func validaSalarioMinimo(ano int, valor float64) error {
	if ano < 1994 {
		return fmt.Errorf("ano invalido")
	}
	if valor <= 0 {
		return fmt.Errorf("valor deve ser maior que zero")
	}
	return nil
}

func mapSalarioMinimoErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		if strings.Contains(pgErr.ConstraintName, "salario_minimo_nacional_ano") {
			return fmt.Errorf("ja existe salario minimo cadastrado para este ano")
		}
	}
	return err
}
