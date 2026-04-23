<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <xsl:include href="Utils.xsl"/>
  <xsl:include href="_Estilos_Geral.xsl"/>
  <xsl:include href="_Scripts_Geral.xsl"/>
  <xsl:include href="Template_Impostos.xsl"/>
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:template match="/" name="RegPass_Passagens">

    <div id="Passagens">
      <fieldset>
        <legend class="titulo-aba">Dados das Passagens</legend>
        <div>
          <xsl:for-each select="//n:Passagens/n:Passagem">
            <fieldset>
              <legend>
                Passagem <xsl:value-of select="position()"/>
              </legend>
              <table>
                <tr>
                  <td>
                    <label>Posto Fiscal</label>
                  </td>
                  <td width="100">
                    <label>UF Passagem</label>
                  </td>
                  <td width="250">
                    <label>Data Hora Passagem</label>
                  </td>
                </tr>
                <tr>
                  <td>
                    <span>
                      <xsl:value-of select="n:cUnidFiscal"/> - <xsl:value-of select="n:xUnidFiscal"/>
                    </span>
                  </td>
                  <td>
                    <span>
                      <xsl:value-of select="n:cUF"/> - <xsl:call-template name="siglaUF">
                        <xsl:with-param name="uf" select="n:cUF"/>
                      </xsl:call-template>
                    </span>
                  </td>
                  <td>
                    <span>
                      <xsl:call-template name="formatDateTime">
                        <xsl:with-param name="dateTime" select="n:dPass"/>
                        <xsl:with-param name="include_as" select="1"/>
                      </xsl:call-template>
                    </span>
                  </td>
                </tr>
              </table>

              <xsl:variable name="rodo" select="count(n:mdTransp/n:rodoviario)"/>
              <xsl:variable name="outro" select="n:mdTransp/n:outro/n:cMod"/>
              <table>
                <tr>
                  <td width="150">
                    <label>Modal</label>
                  </td>
                  <xsl:choose>
                    <xsl:when test="$rodo = 1">
                      <td>
                        <label>Veículo</label>
                      </td>
                      <td>
                        <label>Carreta</label>
                      </td>
                      <td>
                        <label>Carreta 2</label>
                      </td>
                    </xsl:when>
                    <xsl:when test="$outro != ''">
                      <td>
                        <label>Informações do Transporte</label>
                      </td>                      
                    </xsl:when>
                  </xsl:choose>
                </tr>
                <tr>
                  <td>
                    <span>
                      <xsl:choose>
                        <xsl:when test="$rodo = 1">Rodoviário</xsl:when>
                        <xsl:when test="$outro = 'F'">Ferroviário</xsl:when>
                        <xsl:when test="$outro = 'D'">Dutoviário</xsl:when>
                        <xsl:when test="$outro = 'AE'">Aéreo</xsl:when>
                        <xsl:when test="$outro = 'AQ'">Aquaviário</xsl:when>
                        <xsl:when test="$outro = 'O'">Outro</xsl:when>
                      </xsl:choose>
                    </span>
                  </td>

                  <xsl:choose>
                    <xsl:when test="$rodo = 1">
                      <td>
                        <span>
                          <xsl:value-of select="n:mdTransp/n:rodoviario/n:pVeic"/> / <xsl:call-template name="siglaUF">
                            <xsl:with-param name="uf" select="n:mdTransp/n:rodoviario/n:cUFVeic"/>
                          </xsl:call-template>
                        </span>
                      </td>
                      <td>
                        <span>
                          <xsl:choose>
                            <xsl:when test="n:mdTransp/n:rodoviario/n:pCarreta != ''">
                              <xsl:value-of select="n:mdTransp/n:rodoviario/n:pCarreta"/> / <xsl:call-template name="siglaUF">
                                <xsl:with-param name="uf" select="n:mdTransp/n:rodoviario/n:cUFCarreta"/>
                              </xsl:call-template>
                            </xsl:when>
                            <xsl:otherwise>-</xsl:otherwise>
                          </xsl:choose>
                        </span>
                      </td>
                      <td>
                        <span>
                          <xsl:choose>
                            <xsl:when test="n:mdTransp/n:rodoviario/n:pCarreta2 != ''">
                              <xsl:value-of select="n:mdTransp/n:rodoviario/n:pCarreta2"/> / <xsl:call-template name="siglaUF">
                                <xsl:with-param name="uf" select="n:mdTransp/n:rodoviario/n:cUFCarreta2"/>
                              </xsl:call-template>
                            </xsl:when>
                            <xsl:otherwise>-</xsl:otherwise>
                          </xsl:choose>
                        </span>
                      </td>
                    </xsl:when>
                    <xsl:when test="$outro != ''">
                      <td>
                        <span>
                          <xsl:if test="n:mdTransp/n:outro/n:xIdent=''">-</xsl:if>
                          <xsl:value-of select="n:mdTransp/n:outro/n:xIdent"/>
                        </span>
                      </td>
                    </xsl:when>
                  </xsl:choose>

                </tr>
              </table>

              <table>
                <tr>
                  <td width="150">
                    <label>Tipo Transmissão</label>
                  </td>
                  <td>
                    <label>Sentido da Passagem</label>
                  </td>
                </tr>
                <tr>
                  <td>
                    <span>
                      <xsl:variable name="tTrasm" select="n:tTrasm"/>
                      <xsl:choose>
                        <xsl:when test="$tTrasm = 'N'">Normal</xsl:when>
                        <xsl:when test="$tTrasm = 'A'">Atrasada</xsl:when>
                        <xsl:otherwise>-</xsl:otherwise>
                      </xsl:choose>
                    </span>
                  </td>
                  <td>
                    <span>
                      <xsl:variable name="tSentido" select="n:tSentido"/>
                      <xsl:choose>
                        <xsl:when test="$tSentido = 'E'">Entrada na UF </xsl:when>
                        <xsl:when test="$tSentido = 'S'">Saída da UF</xsl:when>
                        <xsl:when test="$tSentido = 'I'">Indeterminado</xsl:when>
                        <xsl:otherwise>-</xsl:otherwise>
                      </xsl:choose>
                    </span>
                  </td>                  
                </tr>
              </table>

              <table>
                <tr>
                  <td>
                    <label>Funcionário</label>
                  </td>
                </tr>
                <tr>
                  <td>
                    <span>
                      <xsl:call-template name="formatCpf">
                        <xsl:with-param name="cpf" select="n:cpfFunc"/>
                      </xsl:call-template> - <xsl:value-of select="n:xFunc"/>
                    </span>
                  </td>
                </tr>
              </table>
              
              <table>
                <tr>
                  <td>
                    <label>Observação</label>
                  </td>
                  <td width="200">
                    <label>Retorno</label>
                  </td>
                </tr>
                <tr>
                  <td>
                    <span>
                      <xsl:if test="n:xObs=''">-</xsl:if>
                      <xsl:value-of select="n:xObs"/>
                    </span>
                  </td>
                  <td>
                    <span>
                      <xsl:variable name="tIndRet" select="n:tIndRet"/>
                      <xsl:choose>
                        <xsl:when test="$tIndRet = 'D'">Devolução</xsl:when>
                        <xsl:when test="$tIndRet = 'R'">Retorno</xsl:when>
                        <xsl:otherwise>-</xsl:otherwise>
                      </xsl:choose>
                    </span>
                  </td>
                </tr>
              </table>

            </fieldset>
          </xsl:for-each>
        </div>
      </fieldset>
    </div>

  </xsl:template>
</xsl:stylesheet>
