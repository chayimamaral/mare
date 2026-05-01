package handlers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chayimamaral/vecx/backend/internal/config"
	"github.com/chayimamaral/vecx/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AIHandler struct {
	cfg  config.Config
	pool *pgxpool.Pool
}

func NewAIHandler(cfg config.Config, pool *pgxpool.Pool) *AIHandler {
	return &AIHandler{cfg: cfg, pool: pool}
}

// PublicStatus expõe se o assistente está habilitado (sem autenticação).
func (h *AIHandler) PublicStatus(w http.ResponseWriter, r *http.Request) {
	render.WriteJSON(w, http.StatusOK, map[string]bool{"iaEnabled": h.cfg.IAEnabled})
}

type aiChatRequest struct {
	Message string `json:"message"`
}

type ollamaGenerateLine struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

// buildAIUserSessionBlock monta o contexto da sessão a partir do JWT (RequireAuth).
// Não coloque dados pessoais ou fiscais no Modelfile: isso seria partilhado por todos os utilizadores.
func buildAIUserSessionBlock(ctx context.Context) string {
	uid := strings.TrimSpace(middleware.UserID(ctx))
	if uid == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString("<<<CONTEXTO_SESSAO_VECX>>>\n")
	b.WriteString("Estes dados identificam APENAS o utilizador desta sessão autenticada. Não misture com outros utilizadores nem reutilize dados de outras sessões.\n")
	if n := strings.TrimSpace(middleware.UserName(ctx)); n != "" {
		fmt.Fprintf(&b, "Nome: %s\n", n)
	}
	if e := strings.TrimSpace(middleware.UserEmail(ctx)); e != "" {
		fmt.Fprintf(&b, "Email: %s\n", e)
	}
	if role := strings.TrimSpace(middleware.Role(ctx)); role != "" {
		fmt.Fprintf(&b, "Perfil: %s\n", role)
	}
	if tn := strings.TrimSpace(middleware.TenantNome(ctx)); tn != "" {
		fmt.Fprintf(&b, "Escritório (tenant): %s\n", tn)
	}
	b.WriteString("<<</CONTEXTO_SESSAO_VECX>>>\n")
	return b.String()
}

// ChatStream encaminha o prompt ao Ollama com stream e devolve SSE ao cliente.
// Dados do usuário saem do browser → backend → Ollama local; não há envio a LLMs remotas.
func (h *AIHandler) ChatStream(w http.ResponseWriter, r *http.Request) {
	if !h.cfg.IAEnabled {
		render.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Serviço Indisponível"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 512*1024)
	var req aiChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		render.WriteError(w, http.StatusBadRequest, "message obrigatorio")
		return
	}

	session := buildAIUserSessionBlock(r.Context())
	flags := repository.AIPreContextFlagsFromMessage(msg)
	dbNarrative := ""
	if h.pool != nil && (flags.NFe || flags.Imposto || flags.Restituicao) {
		dbNarrative = repository.BuildAIPreContextNarrative(r.Context(), h.pool, middleware.TenantID(r.Context()), flags)
	}

	prompt := msg
	var pre strings.Builder
	if session != "" {
		pre.WriteString(session)
		pre.WriteString("\n")
	}
	if dbNarrative != "" {
		pre.WriteString(dbNarrative)
		pre.WriteString("\n")
	}
	if pre.Len() > 0 {
		prompt = pre.String() + "\n### Pergunta do utilizador\n" + msg
	}

	payload, err := json.Marshal(map[string]any{
		"model":  h.cfg.OllamaModel,
		"prompt": prompt,
		"stream": true,
	})
	if err != nil {
		render.WriteError(w, http.StatusInternalServerError, "erro interno")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()

	ollamaReq, err := http.NewRequestWithContext(ctx, http.MethodPost, h.cfg.OllamaGenerateURL, bytes.NewReader(payload))
	if err != nil {
		render.WriteError(w, http.StatusInternalServerError, "erro interno")
		return
	}
	ollamaReq.Header.Set("Content-Type", "application/json")

	ollamaResp, err := http.DefaultClient.Do(ollamaReq)
	if err != nil {
		render.WriteError(w, http.StatusBadGateway, fmt.Sprintf("falha ao contatar Ollama: %v", err))
		return
	}
	defer ollamaResp.Body.Close()

	if ollamaResp.StatusCode != http.StatusOK {
		slurp, _ := io.ReadAll(io.LimitReader(ollamaResp.Body, 4096))
		render.WriteError(w, http.StatusBadGateway, fmt.Sprintf("Ollama: %s — %s", ollamaResp.Status, strings.TrimSpace(string(slurp))))
		return
	}

	fl, ok := w.(http.Flusher)
	if !ok {
		render.WriteError(w, http.StatusInternalServerError, "streaming nao suportado")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	writeSSE := func(obj map[string]any) bool {
		b, jerr := json.Marshal(obj)
		if jerr != nil {
			return false
		}
		if _, werr := fmt.Fprintf(w, "data: %s\n\n", b); werr != nil {
			return false
		}
		fl.Flush()
		return true
	}

	sc := bufio.NewScanner(ollamaResp.Body)
	// Linhas NDJSON do Ollama podem ser longas.
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		var row ollamaGenerateLine
		if err := json.Unmarshal(line, &row); err != nil {
			if !writeSSE(map[string]any{"error": "resposta invalida do Ollama"}) {
				return
			}
			break
		}
		if row.Error != "" {
			if !writeSSE(map[string]any{"error": row.Error}) {
				return
			}
			break
		}
		if row.Response != "" {
			if !writeSSE(map[string]any{"delta": row.Response}) {
				return
			}
		}
		if row.Done {
			_ = writeSSE(map[string]any{"done": true})
			break
		}
	}
	if err := sc.Err(); err != nil && !errors.Is(err, context.Canceled) {
		_ = writeSSE(map[string]any{"error": "leitura interrompida"})
	}
}
