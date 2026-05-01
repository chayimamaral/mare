package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/nfepayload"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AIPreContextFlags indica quais blocos de dados do tenant carregar antes do Ollama.
type AIPreContextFlags struct {
	NFe         bool
	Imposto     bool
	Restituicao bool
}

func normalizeAccentsAI(s string) string {
	repl := strings.NewReplacer(
		"à", "a", "á", "a", "â", "a", "ã", "a", "ä", "a",
		"è", "e", "é", "e", "ê", "e",
		"ì", "i", "í", "i", "î", "i",
		"ò", "o", "ó", "o", "ô", "o", "õ", "o",
		"ù", "u", "ú", "u", "û", "u",
		"ç", "c",
	)
	return repl.Replace(s)
}

// AIPreContextFlagsFromMessage detecta palavras-chave na pergunta (acentos normalizados de forma simples).
func AIPreContextFlagsFromMessage(msg string) AIPreContextFlags {
	n := normalizeAccentsAI(strings.TrimSpace(strings.ToLower(msg)))
	var f AIPreContextFlags
	if strings.Contains(n, "nfe") || strings.Contains(n, "nf-e") || strings.Contains(n, "nota fiscal") ||
		strings.Contains(n, "vnf") || strings.Contains(n, "icmstot") {
		f.NFe = true
	}
	if !f.NFe && (strings.Contains(n, "valor") || strings.Contains(n, "valores")) && strings.Contains(n, "nota") {
		f.NFe = true
	}
	if strings.Contains(n, "restituicao") || strings.Contains(n, "restitu") {
		f.Restituicao = true
	}
	if strings.Contains(n, "imposto") || strings.Contains(n, "tributo") || strings.Contains(n, "tributario") {
		f.Imposto = true
	}
	return f
}

// BuildAIPreContextNarrative executa consultas leves no schema do tenant (search_path via RequireAuth)
// e devolve texto para o prompt do Ollama. Falhas de SQL não interrompem o chat: omitem o bloco ou avisam.
func BuildAIPreContextNarrative(ctx context.Context, pool *pgxpool.Pool, tenantID string, flags AIPreContextFlags) string {
	tenantID = strings.TrimSpace(tenantID)
	if pool == nil || tenantID == "" || (!flags.NFe && !flags.Imposto && !flags.Restituicao) {
		return ""
	}

	var sections []string
	sections = append(sections, "<<<DADOS_VECX_BANCO (resumo do tenant atual; use apenas estes números; não extrapole)>>>")

	if flags.NFe {
		var docCount, syncCount int64
		errDoc := dbQueryRow(ctx, pool, `SELECT COUNT(*)::bigint FROM nfe_documento`).Scan(&docCount)
		errSync := dbQueryRow(ctx, pool, `SELECT COUNT(*)::bigint FROM nfe_sync_estado`).Scan(&syncCount)
		if errDoc != nil || errSync != nil {
			sections = append(sections, "NFe: não foi possível ler totais (tabela ausente ou sem permissão).")
		} else {
			var sumVNF float64
			var docsComVNF int64
			payRows, payErr := dbQuery(ctx, pool, `
				SELECT payload_json
				FROM nfe_documento
				WHERE payload_json IS NOT NULL
				  AND payload_json <> 'null'::jsonb
				  AND payload_json::text <> '{}'`)
			if payErr != nil {
				sections = append(sections, fmt.Sprintf(
					"NFe: documentos armazenados: %d; pontos de sincronização (UF/CNPJ): %d. Não foi possível ler payload_json para somar vNF.",
					docCount, syncCount))
			} else {
				func() {
					defer payRows.Close()
					for payRows.Next() {
						var raw []byte
						if scanErr := payRows.Scan(&raw); scanErr != nil {
							continue
						}
						if v := nfepayload.ValorTotalFromJSON(raw); v != nil {
							sumVNF += *v
							docsComVNF++
						}
					}
					_ = payRows.Err()
				}()
				sections = append(sections, fmt.Sprintf(
					"NFe: documentos armazenados: %d; pontos de sincronização (UF/CNPJ): %d. "+
						"Documentos com valor total vNF (ICMSTot) extraível do JSON: %d; soma desses vNF: %.2f. "+
						"Documentos sem vNF no JSON (layout diferente ou resumo incompleto) não entram na soma.",
					docCount, syncCount, docsComVNF, sumVNF))
			}
		}
	}

	if flags.Imposto {
		rows, err := dbQuery(ctx, pool, `
			SELECT COALESCE(NULLIF(TRIM(ec.status), ''), '(sem status)') AS st,
			       COUNT(*)::bigint,
			       COALESCE(SUM(ec.valor), 0)::float8
			FROM empresa_compromissos ec
			INNER JOIN empresa e ON e.id = ec.empresa_id
			INNER JOIN public.tipoempresa_obrigacao t ON t.id = ec.tipoempresa_obrigacao_id
			WHERE e.tenant_id = $1::uuid AND e.ativo = true
			  AND UPPER(TRIM(COALESCE(t.tipo_classificacao, ''))) IN ('TRIBUTARIA', 'TRIBUTO')
			GROUP BY COALESCE(NULLIF(TRIM(ec.status), ''), '(sem status)')
			ORDER BY st`,
			tenantID)
		if err != nil {
			sections = append(sections, "Impostos (compromissos tributários): não foi possível agregar dados.")
		} else {
			func() {
				defer rows.Close()
				var parts []string
				for rows.Next() {
					var st string
					var cnt int64
					var sum float64
					if scanErr := rows.Scan(&st, &cnt, &sum); scanErr != nil {
						continue
					}
					parts = append(parts, fmt.Sprintf("%s: %d item(ns), soma de valores %.2f", st, cnt, sum))
				}
				if err := rows.Err(); err != nil {
					sections = append(sections, "Impostos (compromissos tributários): erro ao ler agregados.")
					return
				}
				if len(parts) == 0 {
					sections = append(sections, "Impostos (compromissos tributários classificados como TRIBUTARIA/TRIBUTO): nenhum registro.")
				} else {
					sections = append(sections, "Impostos / compromissos tributários por status: "+strings.Join(parts, "; ")+".")
				}
			}()
		}
	}

	if flags.Restituicao {
		var cnt int64
		var sum float64
		err := dbQueryRow(ctx, pool, `
			SELECT COUNT(*)::bigint, COALESCE(SUM(ec.valor), 0)::float8
			FROM empresa_compromissos ec
			INNER JOIN empresa e ON e.id = ec.empresa_id
			WHERE e.tenant_id = $1::uuid AND e.ativo = true
			  AND (
			    UPPER(ec.descricao) LIKE '%IRPF%'
			    OR UPPER(ec.descricao) LIKE '%RESTIT%'
			  )`,
			tenantID).Scan(&cnt, &sum)
		if err != nil {
			sections = append(sections, "Restituição (heurística por descrição IRPF/restituição): não foi possível consultar.")
		} else {
			sections = append(sections, fmt.Sprintf("Restituição (heurística: compromissos cuja descrição menciona IRPF/restituição): %d registro(s); soma de valores cadastrados %.2f (valores são os do sistema, não declaração Receita).", cnt, sum))
		}
	}

	sections = append(sections, "<<</DADOS_VECX_BANCO>>>")
	return strings.Join(sections, "\n")
}
