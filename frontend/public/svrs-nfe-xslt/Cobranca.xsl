<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:output method="html"/>
  <xsl:include href="Utils.xsl"/>
  <xsl:include href="_Estilos_Geral.xsl"/>
  <xsl:include href="_Scripts_Geral.xsl"/>
  <xsl:template match="/" name="Cobranca">
    <div id="Cobranca"  class="GeralXslt">

      <xsl:variable name = "cobr" select = "//n:cobr"/>
      <xsl:if test="$cobr != ''">      
        <fieldset>
        <legend class="titulo-aba">Dados de Cobrança</legend>
        <xsl:variable name = "fatura" select = "//n:cobr/n:fat"/>
        <xsl:if test = "$fatura != ''">
          <fieldset>
            <legend>Fatura</legend>
            <table class="box">
              <tr class="col-3">
                <xsl:for-each select="//n:cobr/n:fat/n:nFat">
                  <td>
                    <label>Número</label>
                    <span>
                      <xsl:value-of select = "text()"/>
                    </span>
                  </td>
                </xsl:for-each>
                <td>
                  <label>Valor Original</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:cobr/n:fat/n:vOrig"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td>
                  <label>Valor do Desconto</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:cobr/n:fat/n:vDesc"/>
                    </xsl:call-template> 
                  </span>
                </td>
              </tr>
              <tr>
                <td>
                  <label>Valor Líquido</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:cobr/n:fat/n:vLiq"/>
                    </xsl:call-template>
                  </span>
                </td>
              </tr>
            </table>
          </fieldset>
          <br />
        </xsl:if>
        <xsl:variable name = "duplicata" select = "//n:cobr/n:dup"/>
        <xsl:if test = "$duplicata != ''">
          <fieldset>
            <legend>Duplicatas</legend>
            <table class="box">
              <tr class="col-3">
                <td>
                  <label>Número</label>
                </td>
                <td>
                  <label>Vencimento</label>
                </td>
                <td>
                  <label>Valor</label>
                </td>
              </tr>
              <xsl:for-each select = "//n:cobr/n:dup">
                <tr class="col-3">
                  <td>
                    <span>
                      <xsl:value-of select = "n:nDup"/>
                    </span>
                  </td>
                  <td>
                    <span>
                      <xsl:call-template name="formatDate">
                        <xsl:with-param name="date" select="n:dVenc"/>
                      </xsl:call-template>
                    </span>
                  </td>
                  <td>
                    <span>
                      <xsl:call-template name="format2Casas">
                        <xsl:with-param name="num" select="n:vDup"/>
                      </xsl:call-template> 
                    </span>
                  </td>
                </tr>
              </xsl:for-each>
            </table>
          </fieldset>
          <br />
        </xsl:if>
      </fieldset>
      </xsl:if>

      <fieldset>
        <xsl:variable name = "pagamento" select = "//n:pag"/>
        <xsl:if test = "$pagamento != ''">
          <legend class="titulo-aba">Formas de Pagamento</legend>
          <fieldset>
            <xsl:variable name = "tPag" select = "//n:pag/n:tPag"/>
            <table class="box">
              <tr class="col-5">
                <td>
                  <label>Forma de Pagamento</label>
                </td>
                <td>
                  <label>Valor do Pagamento</label>
                </td>
                <td>
                  <label>CNPJ da Credenciadora </label>
                </td>
                <td>
                  <label>Bandeira da operadora</label>
                </td>
                <td>
                  <label> Número de autorização </label>
                </td>
              </tr>
            <xsl:for-each select="//n:infNFe/n:pag">
              <tr>
                <td>
                  <span> 
                    <xsl:choose>
                      <xsl:when test="n:tPag='01'">
                        1 - Dinheiro
                      </xsl:when>
                      <xsl:when test="n:tPag='02'">
                        2 - Cheque
                      </xsl:when>
                      <xsl:when test="n:tPag='03'">
                        3 - Cartão de Crédito
                      </xsl:when>
                      <xsl:when test="n:tPag='04'">
                        4 - Cartão de Débito
                      </xsl:when>
                      <xsl:when test="n:tPag='05'">
                        5 - Crédito Loja
                      </xsl:when>
                      <xsl:when test="n:tPag='10'">
                        10 - Vale Alimentação
                      </xsl:when>
                      <xsl:when test="n:tPag='11'">
                        11 - Vale Refeição
                      </xsl:when>
                      <xsl:when test="n:tPag='12'">
                        12 - Vale Presente
                      </xsl:when>
                      <xsl:when test="n:tPag='13'">
                        13 - Vale Combustível
                      </xsl:when>
                      <xsl:when test="n:tPag='99'">
                        99 - Outros
                      </xsl:when>
                      <xsl:otherwise>
                        <xsl:value-of select="n:tPag"/>
                      </xsl:otherwise>
                    </xsl:choose> 
                  </span>
                </td>
                <td>
                  <span>
                    <xsl:variable name = "vPag" select = "n:vPag"/>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="$vPag"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td>
                  <span>
                    <xsl:variable name = "cnpj" select = "n:card/n:CNPJ"/>
                    <xsl:value-of select="$cnpj"/>
                  </span>
                </td>
                <td>
                  <span>
                    <xsl:variable name = "tBand" select = "n:card/n:tBand"/>
                    <xsl:choose>
                      <xsl:when test="$tBand = 1">
                        01 - Visa
                      </xsl:when>
                      <xsl:when test="$tBand = 2">
                        02 - Mastercard
                      </xsl:when>
                      <xsl:when test="$tBand = 3">
                        03 - American Express
                      </xsl:when>
                      <xsl:when test="$tBand = 4">
                        04 - Sorocred
                      </xsl:when>
                      <xsl:when test="$tBand = 99">
                        99 - Outros
                      </xsl:when> 
                      <xsl:otherwise>
                        <xsl:value-of select="$tBand"/>
                      </xsl:otherwise>
                    </xsl:choose> 
                  </span>
                </td>
                <td>
                  <span>
                    <xsl:variable name = "cAut" select = "n:card/n:cAut"/>
                    <xsl:value-of select="$cAut"/>
                  </span>
                </td>
              </tr>
            </xsl:for-each>
          </table> 
        </fieldset> 
        </xsl:if> 
      </fieldset> 
    </div>
  </xsl:template>
</xsl:stylesheet>
