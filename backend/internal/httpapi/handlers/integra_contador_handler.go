package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/httpapi/middleware"
	"github.com/chayimamaral/vecx/backend/internal/httpapi/render"
	"github.com/chayimamaral/vecx/backend/internal/service"
)

type IntegraContadorHandler struct {
	svc *service.IntegraContadorService
}

func NewIntegraContadorHandler(svc *service.IntegraContadorService) *IntegraContadorHandler {
	return &IntegraContadorHandler{svc: svc}
}

type integraCallRequest struct {
	Ambiente                  string                      `json:"ambiente"`
	Operacao                  string                      `json:"operacao"`
	Payload                   service.IntegraDadosEntrada `json:"payload"`
	AccessToken               string                      `json:"access_token"`
	JWTToken                  string                      `json:"jwt_token"`
	AutenticarProcuradorToken string                      `json:"autenticar_procurador_token"`
}

func (h *IntegraContadorHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantID(r.Context())
	out, err := h.svc.Authenticate(r.Context(), tenantID)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func (h *IntegraContadorHandler) Call(w http.ResponseWriter, r *http.Request) {
	var req integraCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	amb := strings.ToLower(strings.TrimSpace(req.Ambiente))
	if amb == "" {
		amb = string(service.IntegraAmbienteTrial)
	}
	out, err := h.svc.Call(r.Context(), service.IntegraCallInput{
		TenantID:                  middleware.TenantID(r.Context()),
		Ambiente:                  service.IntegraContadorAmbiente(amb),
		Operacao:                  req.Operacao,
		Payload:                   req.Payload,
		AccessToken:               req.AccessToken,
		JWTToken:                  req.JWTToken,
		AutenticarProcuradorToken: req.AutenticarProcuradorToken,
	})
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}

func (h *IntegraContadorHandler) PGMEIGerarDAS(w http.ResponseWriter, r *http.Request) {
	h.callPGMEI(w, r, "emitir", "GERARDASPDF21")
}

func (h *IntegraContadorHandler) PGMEIGerarDASCodBarras(w http.ResponseWriter, r *http.Request) {
	h.callPGMEI(w, r, "emitir", "GERARDASCODBARRA22")
}

func (h *IntegraContadorHandler) PGMEIAtualizarBeneficio(w http.ResponseWriter, r *http.Request) {
	h.callPGMEI(w, r, "emitir", "ATUBENEFICIO23")
}

func (h *IntegraContadorHandler) PGMEIConsultarDividaAtiva(w http.ResponseWriter, r *http.Request) {
	h.callPGMEI(w, r, "consultar", "DIVIDAATIVA24")
}

type pgmeiRequest struct {
	Ambiente                  string                       `json:"ambiente"`
	Contratante               service.IntegraIdentificacao `json:"contratante"`
	AutorPedidoDados          service.IntegraIdentificacao `json:"autorPedidoDados"`
	Contribuinte              service.IntegraIdentificacao `json:"contribuinte"`
	Dados                     json.RawMessage              `json:"dados"`
	AccessToken               string                       `json:"access_token"`
	JWTToken                  string                       `json:"jwt_token"`
	AutenticarProcuradorToken string                       `json:"autenticar_procurador_token"`
}

func (h *IntegraContadorHandler) callPGMEI(w http.ResponseWriter, r *http.Request, operacao, idServico string) {
	var req pgmeiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.WriteError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	rawDados := strings.TrimSpace(string(req.Dados))
	if rawDados == "" || rawDados == "null" {
		render.WriteError(w, http.StatusBadRequest, "dados obrigatorio")
		return
	}
	dadosEscaped, err := json.Marshal(rawDados)
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, "dados invalido")
		return
	}
	payload := service.IntegraDadosEntrada{
		Contratante:  req.Contratante,
		AutorPedido:  req.AutorPedidoDados,
		Contribuinte: req.Contribuinte,
		PedidoDados: service.IntegraPedidoDados{
			IDSistema:     "PGMEI",
			IDServico:     idServico,
			VersaoSistema: "1.0",
			Dados:         string(dadosEscaped),
		},
	}
	amb := strings.ToLower(strings.TrimSpace(req.Ambiente))
	if amb == "" {
		amb = string(service.IntegraAmbienteTrial)
	}

	out, err := h.svc.Call(r.Context(), service.IntegraCallInput{
		TenantID:                  middleware.TenantID(r.Context()),
		Ambiente:                  service.IntegraContadorAmbiente(amb),
		Operacao:                  operacao,
		Payload:                   payload,
		AccessToken:               req.AccessToken,
		JWTToken:                  req.JWTToken,
		AutenticarProcuradorToken: req.AutenticarProcuradorToken,
	})
	if err != nil {
		render.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	render.WriteJSON(w, http.StatusOK, out)
}
