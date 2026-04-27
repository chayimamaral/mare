package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/domain"
	"github.com/chayimamaral/vecontab/backend/internal/repository"
	"github.com/jackc/pgx/v5"
)

type NFEErroDetalhado struct {
	Codigo         string `json:"codigo"`
	Mensagem       string `json:"mensagem"`
	EtapaValidacao string `json:"etapa_validacao"`
	AcaoSugerida   string `json:"acao_sugerida,omitempty"`
	Origem         string `json:"origem"`
}

type NFEValidacaoCatalogoService struct {
	repo *repository.NFEValidacaoCatalogoRepository
}

func NewNFEValidacaoCatalogoService(repo *repository.NFEValidacaoCatalogoRepository) *NFEValidacaoCatalogoService {
	return &NFEValidacaoCatalogoService{repo: repo}
}

func (s *NFEValidacaoCatalogoService) ListRegrasAtivasPorEtapa(ctx context.Context, etapa string) ([]domain.NFEValidacaoRegra, error) {
	if s == nil || s.repo == nil {
		return nil, nil
	}
	return s.repo.ListRegrasAtivasPorEtapa(ctx, etapa)
}

var nfeCodigoRegex = regexp.MustCompile(`\b\d{3,4}\b`)

func (s *NFEValidacaoCatalogoService) ResolverErro(ctx context.Context, origem, etapa string, err error) NFEErroDetalhado {
	rawMsg := "falha no processamento da NF-e"
	if err != nil {
		rawMsg = strings.TrimSpace(err.Error())
	}
	codigo := "NFE-NAO-MAPEADO"
	lowerMsg := strings.ToLower(rawMsg)
	if strings.Contains(lowerMsg, "certificado nao encontrado") {
		codigo = "CERTIFICADO_NAO_ENCONTRADO"
	}
	if strings.Contains(lowerMsg, "json") {
		codigo = "400"
	}
	if strings.Contains(lowerMsg, "chave_nfe") || strings.Contains(lowerMsg, "cnpj/cpf") || strings.Contains(lowerMsg, "informe retorno ou xml") {
		codigo = "DADOS_INVALIDOS"
	}
	if strings.Contains(lowerMsg, "consulta nfe serpro status") {
		codigo = "CSTAT_INDISPONIVEL"
	}
	if m := nfeCodigoRegex.FindString(rawMsg); m != "" {
		codigo = m
	}

	if s != nil && s.repo != nil {
		cat, catErr := s.repo.GetCodigoErro(ctx, origem, codigo)
		if catErr == nil && cat != nil {
			return NFEErroDetalhado{
				Codigo:         cat.Codigo,
				Mensagem:       cat.Mensagem,
				EtapaValidacao: cat.EtapaValidacao,
				AcaoSugerida:   cat.AcaoSugerida,
				Origem:         cat.Origem,
			}
		}
		if catErr != nil && !errors.Is(catErr, pgx.ErrNoRows) {
			return NFEErroDetalhado{
				Codigo:         "NFE-CATALOGO-INDISPONIVEL",
				Mensagem:       fmt.Sprintf("%s (%s)", rawMsg, "catalogo de erro indisponivel"),
				EtapaValidacao: etapa,
				Origem:         origem,
			}
		}
	}

	return NFEErroDetalhado{
		Codigo:         codigo,
		Mensagem:       rawMsg,
		EtapaValidacao: etapa,
		Origem:         origem,
	}
}
