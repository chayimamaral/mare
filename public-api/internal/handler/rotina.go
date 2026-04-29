package handler

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/chayimamaral/vecx/public-api/internal/repository"
	"github.com/gin-gonic/gin"
)

var uuidPathSegment = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type Rotina struct {
	repo *repository.Rotina
}

func NewRotina(repo *repository.Rotina) *Rotina {
	return &Rotina{repo: repo}
}

func (h *Rotina) GetPorMunicipioTipo(c *gin.Context) {
	municipioID := strings.TrimSpace(c.Param("municipio_id"))
	tipoEmpresaID := strings.TrimSpace(c.Param("tipo_empresa_id"))
	if municipioID == "" || tipoEmpresaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "municipio_id e tipo_empresa_id são obrigatórios"})
		return
	}
	if !uuidPathSegment.MatchString(municipioID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "municipio_id inválido"})
		return
	}
	tl := strings.ToLower(tipoEmpresaID)
	if tl != "null" && tl != "none" && !uuidPathSegment.MatchString(tipoEmpresaID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo_empresa_id deve ser UUID ou os literais null/none (rotina sem tipo)"})
		return
	}

	payload, err := h.repo.ListByMunicipioTipoJSON(c.Request.Context(), municipioID, tipoEmpresaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao consultar rotinas"})
		_ = c.Error(err)
		log.Printf("Erro no banco: %v", err)
		return
	}

	c.Data(http.StatusOK, "application/json", payload)
}
