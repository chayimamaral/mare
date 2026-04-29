package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/domain"
	"github.com/chayimamaral/vecx/backend/internal/repository"
)

type CaixaPostalService struct {
	repo *repository.CaixaPostalRepository
}

func NewCaixaPostalService(repo *repository.CaixaPostalRepository) *CaixaPostalService {
	return &CaixaPostalService{repo: repo}
}

func (s *CaixaPostalService) ListMensagens(ctx context.Context, schemaName, tenantID string) ([]domain.CaixaPostalMensagem, error) {
	if strings.TrimSpace(schemaName) == "" {
		return nil, fmt.Errorf("schema nao definido")
	}
	if strings.TrimSpace(tenantID) == "" {
		return nil, fmt.Errorf("tenant_id nao definido")
	}
	return s.repo.List(ctx, schemaName, tenantID)
}

func (s *CaixaPostalService) Count(ctx context.Context, schemaName, tenantID string) (int, error) {
	if strings.TrimSpace(schemaName) == "" {
		return 0, fmt.Errorf("schema nao definido")
	}
	if strings.TrimSpace(tenantID) == "" {
		return 0, fmt.Errorf("tenant_id nao definido")
	}
	return s.repo.CountUnread(ctx, schemaName, tenantID)
}

func (s *CaixaPostalService) MarkAsRead(ctx context.Context, schemaName, msgID, userID string) error {
	if schemaName == "" {
		return fmt.Errorf("schema não definido")
	}
	return s.repo.MarkAsRead(ctx, schemaName, msgID, userID)
}

func (s *CaixaPostalService) Enviar(ctx context.Context, targetTenantID string, isGlobal bool, titulo, conteudo, remetenteID, remetenteNome, requesterRole, currentSchema string) error {
	if len(strings.TrimSpace(titulo)) == 0 || len(strings.TrimSpace(conteudo)) == 0 {
		return fmt.Errorf("título e conteúdo são obrigatórios")
	}

	role := strings.ToUpper(strings.TrimSpace(requesterRole))
	targetTenantID = strings.TrimSpace(targetTenantID)
	currentSchema = strings.TrimSpace(currentSchema)

	tipo := "INBOX"

	// Cenário A: Usuário enviando suporte/dúvida para VEC Sistemas (SUPER)
	if role != "SUPER" {
		superSchema, err := s.repo.GetSuperSchema(ctx)
		if err != nil {
			return err
		}

		msgOutbox := domain.CaixaPostalMensagem{
			RemetenteID:   &remetenteID,
			RemetenteNome: remetenteNome,
			Tipo:          "OUTBOX",
			IsGlobal:      false,
			Titulo:        titulo,
			Conteudo:      conteudo,
		}
		// Salva OUTBOX no proprio schema do requerente (quem enviou) pra ficar registrado.
		err = s.repo.Insert(ctx, currentSchema, msgOutbox)
		if err != nil {
			return fmt.Errorf("erro ao salvar OUTBOX: %w", err)
		}

		msgInbox := msgOutbox
		msgInbox.Tipo = "INBOX"
		// Salva INBOX na caixa postal da VEC
		return s.repo.Insert(ctx, superSchema, msgInbox)
	}

	// Cenário B: SUPER emitindo avisos do sistema (Global ou pra um Tenant Específico)
	if isGlobal && targetTenantID != "" {
		return fmt.Errorf("envio global não deve informar tenant destino")
	}
	if !isGlobal && targetTenantID == "" {
		return fmt.Errorf("tenant destino é obrigatório para envio específico")
	}
	superSchema, err := s.repo.GetSuperSchema(ctx)
	if err != nil {
		return err
	}

	msgBase := domain.CaixaPostalMensagem{
		RemetenteID:   &remetenteID,
		RemetenteNome: remetenteNome,
		Tipo:          tipo,
		IsGlobal:      isGlobal,
		Titulo:        titulo,
		Conteudo:      conteudo,
	}

	if isGlobal {
		// Repassa para TODOS os schemas ativos no banco.
		schemas, err := s.repo.GetAllActiveSchemas(ctx)
		if err != nil {
			return err
		}

		currentSchemaLower := strings.ToLower(currentSchema)
		superSchemaLower := strings.ToLower(strings.TrimSpace(superSchema))
		seen := make(map[string]struct{}, len(schemas))
		for _, schema := range schemas {
			schema = strings.TrimSpace(schema)
			if schema == "" {
				continue
			}
			schemaLower := strings.ToLower(schema)
			if schemaLower == currentSchemaLower {
				// Evita cair na própria INBOX do SUPER quando já será registrado em OUTBOX.
				continue
			}
			if superSchemaLower != "" && schemaLower == superSchemaLower {
				// Garante regra de negócio: tenant da plataforma (SUPER) não recebe INBOX global.
				continue
			}
			if _, exists := seen[schemaLower]; exists {
				continue
			}
			seen[schemaLower] = struct{}{}
			_ = s.repo.Insert(ctx, schema, msgBase)
		}

		// Guarda na caixa do Super como OUTBOX
		msgOutbox := msgBase
		msgOutbox.Tipo = "OUTBOX"
		_ = s.repo.Insert(ctx, currentSchema, msgOutbox)

		return nil
	}

	// Tenant específico
	targetSchema, err := s.repo.GetSchemaByTenantID(ctx, targetTenantID)
	if err != nil {
		return fmt.Errorf("tenant destino (%s) inválido ou sem schema: %w", targetTenantID, err)
	}
	if strings.EqualFold(strings.TrimSpace(targetSchema), strings.TrimSpace(superSchema)) {
		return fmt.Errorf("tenant do provedor (SUPER) não pode ser destinatário")
	}

	err = s.repo.Insert(ctx, targetSchema, msgBase)
	if err != nil {
		return err
	}

	msgOutbox := msgBase
	msgOutbox.Tipo = "OUTBOX"
	_ = s.repo.Insert(ctx, currentSchema, msgOutbox)

	return nil
}
