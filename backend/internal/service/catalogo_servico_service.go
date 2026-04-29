package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/chayimamaral/vecx/backend/internal/repository"
)

type CatalogoServicoService struct {
	repo *repository.CatalogoServicoRepository
}

type CatalogoServicoInput struct {
	ID              string `json:"id"`
	Secao           string `json:"secao"`
	Sequencial      int    `json:"sequencial"`
	Codigo          string `json:"codigo"`
	IDSistema       string `json:"id_sistema"`
	IDServico       string `json:"id_servico"`
	DataImplantacao string `json:"data_implantacao"`
	Tipo            string `json:"tipo"`
	Descricao       string `json:"descricao"`
}

func NewCatalogoServicoService(repo *repository.CatalogoServicoRepository) *CatalogoServicoService {
	return &CatalogoServicoService{repo: repo}
}

func (s *CatalogoServicoService) List(ctx context.Context, secao string, incluirInativos bool) ([]domain.CatalogoServico, error) {
	return s.repo.List(ctx, secao, incluirInativos)
}

func (s *CatalogoServicoService) Create(ctx context.Context, input CatalogoServicoInput) (domain.CatalogoServico, error) {
	if err := validarCatalogoServico(input); err != nil {
		return domain.CatalogoServico{}, err
	}
	return s.repo.Create(ctx, repository.CatalogoServicoUpsertInput{
		Secao:           strings.TrimSpace(input.Secao),
		Sequencial:      input.Sequencial,
		Codigo:          strings.TrimSpace(input.Codigo),
		IDSistema:       strings.TrimSpace(input.IDSistema),
		IDServico:       strings.TrimSpace(input.IDServico),
		DataImplantacao: strings.TrimSpace(input.DataImplantacao),
		Tipo:            strings.TrimSpace(input.Tipo),
		Descricao:       strings.TrimSpace(input.Descricao),
	})
}

func (s *CatalogoServicoService) Update(ctx context.Context, input CatalogoServicoInput) (domain.CatalogoServico, error) {
	if strings.TrimSpace(input.ID) == "" {
		return domain.CatalogoServico{}, fmt.Errorf("id obrigatorio")
	}
	if err := validarCatalogoServico(input); err != nil {
		return domain.CatalogoServico{}, err
	}
	return s.repo.Update(ctx, repository.CatalogoServicoUpsertInput{
		ID:              strings.TrimSpace(input.ID),
		Secao:           strings.TrimSpace(input.Secao),
		Sequencial:      input.Sequencial,
		Codigo:          strings.TrimSpace(input.Codigo),
		IDSistema:       strings.TrimSpace(input.IDSistema),
		IDServico:       strings.TrimSpace(input.IDServico),
		DataImplantacao: strings.TrimSpace(input.DataImplantacao),
		Tipo:            strings.TrimSpace(input.Tipo),
		Descricao:       strings.TrimSpace(input.Descricao),
	})
}

func (s *CatalogoServicoService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id obrigatorio")
	}
	return s.repo.Delete(ctx, strings.TrimSpace(id))
}

func validarCatalogoServico(input CatalogoServicoInput) error {
	if strings.TrimSpace(input.Secao) == "" {
		return fmt.Errorf("secao obrigatoria")
	}
	if input.Sequencial <= 0 {
		return fmt.Errorf("sequencial deve ser maior que zero")
	}
	if strings.TrimSpace(input.Codigo) == "" {
		return fmt.Errorf("codigo obrigatorio")
	}
	if strings.TrimSpace(input.IDSistema) == "" {
		return fmt.Errorf("id_sistema obrigatorio")
	}
	if strings.TrimSpace(input.IDServico) == "" {
		return fmt.Errorf("id_servico obrigatorio")
	}
	if strings.TrimSpace(input.Tipo) == "" {
		return fmt.Errorf("tipo obrigatorio")
	}
	if strings.TrimSpace(input.Descricao) == "" {
		return fmt.Errorf("descricao obrigatoria")
	}
	return nil
}
