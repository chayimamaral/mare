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

  <xsl:template match="/" name="Totais">
    <div id="Totais" class="GeralXslt">
      <fieldset>
        <legend class="titulo-aba">Totais</legend>
        <fieldset>
          <legend>
            ICMS
          </legend>
          <table class="box">
            <tr class="col-4">
              <td>
                <label>Base de Cálculo ICMS</label>
                <span>
                  <xsl:if test="//n:total/n:ICMSTot/n:vBC != ''">
                     <xsl:value-of select = "format-number(//n:total/n:ICMSTot/n:vBC,'##.##.##0,00')"/>
                  </xsl:if>
                </span>
              </td>
              <td>
                <label>Valor do ICMS</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vICMS"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Valor do ICMS Desonerado</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vICMSDeson"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Base de Cálculo ICMS ST</label>
                <span>
                  <xsl:if test="//n:total/n:ICMSTot/n:vBCST != ''">
                    <xsl:value-of select = "format-number(//n:total/n:ICMSTot/n:vBCST,'##.##.##0,00')"/>
                  </xsl:if>
                </span>
              </td> 
            </tr>
            <tr>
              <td>
                <label>Valor ICMS Substituição</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vST"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>
                  Valor Total dos Produtos
                </label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vProd"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Valor do Frete</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vFrete"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Valor do Seguro</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vSeg"/>
                  </xsl:call-template>
                </span>
              </td>             
            </tr>
            <tr>
              <td>
                <label>Outras Despesas Acessórias</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vOutro"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Valor Total do IPI</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vIPI"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Valor Total da NFe</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vNF"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Valor Total dos Descontos</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vDesc"/>
                  </xsl:call-template>
                </span>
              </td>              
            </tr>
            <tr>
              <td>
                <label>Valor Total do II</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vII"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Valor do PIS</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vPIS"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Valor da COFINS</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vCOFINS"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Valor Aproximado dos Tributos</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vTotTrib"/>
                  </xsl:call-template>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
        <xsl:variable name = "issqn" select = "//n:total/n:ISSQNtot"/>
        <xsl:if test = "$issqn != ''">
          <br />
          <fieldset>
            <legend>
              ISSQN
            </legend>
            <table class="box">
              <tr class="col-3">
                <td>
                  <label>Valor&#160;Total&#160;Serv.&#160;Não&#160;Tributados&#160;p/&#160;ICMS&#160;</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:ISSQNtot/n:vServ"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td>
                  <label>Base de Cálculo do ISS</label>
                  <span>
                    <xsl:value-of select = "format-number(//n:total/n:ISSQNtot/n:vBC,'##.##.##0,00')"/>
                  </span>
                </td>
                <td>
                  <label>Valor Total do ISS</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:ISSQNtot/n:vISS"/>
                    </xsl:call-template>
                  </span>
                </td>
              </tr>
              <tr>
                <td>
                  <label>Valor do PIS sobre Serviços</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:ISSQNtot/n:vPIS"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td>
                  <label>Valor da COFINS sobre Serviços</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:ISSQNtot/n:vCOFINS"/>
                    </xsl:call-template>
                  </span>
                </td> 
                <td>
                  <label>Data Prestação Serviço</label>
                  <span> 
                    <xsl:call-template name="formatDate">
                      <xsl:with-param name="date" select="//n:total/n:ISSQNtot/n:dCompet"/>
                    </xsl:call-template> 
                  </span>
                </td> 
              </tr>


              <tr>
                <td>
                  <label>Valor Dedução para Redução da BC</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:ISSQNtot/n:vDeducao"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td>
                  <label>Valor Outras Retenções</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:ISSQNtot/n:vOutro"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td>
                  <label>Valor Desconto Incondicionado</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:ISSQNtot/n:vDescIncond"/>
                    </xsl:call-template>
                  </span>
                </td>
              </tr> 
              <tr>
                <td>
                  <label>Valor Desconto Condicionado</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:ISSQNtot/n:vDescCond"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td>
                  <label>Valor Total Retenção ISS</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:ISSQNtot/n:vISSRet"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td>
                  <label>Código Regime Tributação</label>
                  <span>
                    <xsl:variable name="codRegEspTrib" select="//n:total/n:ISSQNtot/n:cRegTrib"/>
                    <xsl:choose>
                      <xsl:when test="$codRegEspTrib = 1">
                        01 - Microempresa Municipal
                      </xsl:when>
                      <xsl:when test="$codRegEspTrib = 2">
                        02 - Estimativa
                      </xsl:when>
                      <xsl:when test="$codRegEspTrib = 3">
                        03 - Sociedade de Profissionais
                      </xsl:when>
                      <xsl:when test="$codRegEspTrib = 4">
                        04 - Cooperativa
                      </xsl:when>
                      <xsl:when test="$codRegEspTrib = 5">
                        05 - Microempresário Individual (MEI)
                      </xsl:when>
                      <xsl:when test="$codRegEspTrib = 6">
                        06 - Microempresário e Empresa de Pequeno Porte (ME/EPP)
                      </xsl:when>
                      <xsl:otherwise>
                        <xsl:value-of select="$codRegEspTrib"/>
                      </xsl:otherwise>
                    </xsl:choose> 
                  </span>
                </td>
              </tr> 
            </table>
          </fieldset>
        </xsl:if>
        <xsl:variable name = "rt" select = "//n:total/n:retTrib"/>
        <xsl:if test = "$rt != ''">
          <br />
          <fieldset>
            <legend>
              Retenção de Tributos
            </legend>
            <table class="box">
              <tr class="col-3">
                <td>
                  <label>Valor Retido PIS</label>
                  <span>
                    <xsl:choose>
                      <xsl:when test="//n:total/n:retTrib/n:vRetPIS !=''">
                        <xsl:call-template name="format2Casas">
                          <xsl:with-param name="num" select="//n:total/n:retTrib/n:vRetPIS"/>
                        </xsl:call-template>
                      </xsl:when>
                    </xsl:choose>
                  </span>
                </td>
                <td>
                  <label>Valor Retido COFINS</label>
                  <span> 
                    <xsl:choose>
                      <xsl:when test="//n:total/n:retTrib/n:vRetCOFINS !=''">
                        <xsl:call-template name="format2Casas">
                          <xsl:with-param name="num" select="//n:total/n:retTrib/n:vRetCOFINS"/>
                        </xsl:call-template>
                      </xsl:when>
                    </xsl:choose> 
                  </span>
                </td>
                <td>
                  <label>Valor Retido CSLL</label>
                  <span> 
                      <xsl:call-template name="format2Casas">
                        <xsl:with-param name="num" select="//n:total/n:retTrib/n:vRetCSLL"/>
                      </xsl:call-template>                    
                  </span>
                </td>
              </tr>
              <tr>
                <td>
                  <label>Base de Cálculo IRRF</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:retTrib/n:vBCIRRF"/>
                    </xsl:call-template> 
                  </span>
                </td>
                <td>
                  <label>Valor Retido IRRF</label>
                  <span>                      
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:retTrib/n:vIRRF"/>
                    </xsl:call-template>                    
                  </span>
                </td>
                <td>
                  <label>Base de Cálculo Previdência Social</label>
                  <span>
                      <xsl:call-template name="format2Casas">
                        <xsl:with-param name="num" select="//n:total/n:retTrib/n:vBCRetPrev"/>
                      </xsl:call-template> 
                  </span>
                </td>
              </tr>
              <tr>
                <td>
                  <label>Valor Retido Previdência Social</label>
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="//n:total/n:retTrib/n:vRetPrev"/>
                    </xsl:call-template> 
                  </span>
                </td>
              </tr>
            </table>
          </fieldset>
        </xsl:if>
      </fieldset>
    </div>
  </xsl:template>
</xsl:stylesheet>