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

  <xsl:template match="/" name="Transporte">
    <div id="Transporte" class="GeralXslt">
      <fieldset>
        <legend class="titulo-aba">Dados do Transporte</legend>
        <table class="box">
          <tr>
            <td>
              <label>Modalidade do Frete</label>
              <span>
                <xsl:variable name = "modalidade" select = "//n:transp/n:modFrete"/>
                <xsl:choose>
                  <xsl:when test="$modalidade='0'">
                    0 - Por Conta do Emitente
                  </xsl:when>
                  <xsl:when test="$modalidade='1'">
                    1 - Por Conta do Destintário
                  </xsl:when>
                  <xsl:when test="$modalidade='2'">
                    2 - Por conta de Terceiros
                  </xsl:when>
                  <xsl:when test="$modalidade='9'">
                    9 - Sem Frete
                  </xsl:when>
                  <xsl:otherwise>
                    <xsl:value-of select = "$modalidade"/>
                  </xsl:otherwise>
                </xsl:choose> 
              </span>
            </td>
          </tr>
        </table>
      </fieldset>
      <xsl:variable name = "transporta" select = "//n:transp/n:transporta"/>
      <xsl:if test = "$transporta != ''">
        <fieldset>
          <legend class="titulo-aba-interna">Transportador</legend>
          <table class="box">
            <tr class="col-3">
              <td>
                <xsl:variable name="cnpj" select="//n:transp/n:transporta/n:CNPJ"/>
                <xsl:variable name="cpf" select="//n:transp/n:transporta/n:CPF"/>
                <xsl:if test="$cnpj!=''">
                  <label>CNPJ</label>
                </xsl:if>
                <xsl:if test="$cpf!=''">
                  <label>CPF</label>
                </xsl:if>
                <span>
                  <xsl:call-template name="formatCnpj">
                    <xsl:with-param name="cnpj" select="$cnpj"/>
                  </xsl:call-template>
                  <xsl:call-template name="formatCpf">
                    <xsl:with-param name="cpf" select="$cpf"/>
                  </xsl:call-template>
                </span>
              </td>
              <td colspan="2">
                <label>Razão Social / Nome</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:transporta/n:xNome"/>
                </span>
              </td>
              <td></td>
            </tr>
            <tr class="col-3">
              <td>
                <label>Inscrição Estadual</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:transporta/n:IE"/>
                </span>
              </td>
              <td>
                <label>Endereço Completo</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:transporta/n:xEnder"/>
                </span>
              </td>
              <td>
                <label>Município</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:transporta/n:xMun"/>
                </span>
              </td>
            </tr>
            <tr>
              <td>
                <label>UF</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:transporta/n:UF"/>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
      </xsl:if>
      <xsl:variable name = "retencaoICMS" select = "//n:transp/n:retTransp"/>
      <xsl:if test = "$retencaoICMS != ''">
        <fieldset>
          <legend class="titulo-aba-interna">Retenção do ICMS</legend>
          <table class="box">
            <tr class="col-3">
              <td>
                <label>Valor do Serviço</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:transp/n:retTransp/n:vServ"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Base de Cálculo</label>
                <span>
                  <xsl:value-of select = "format-number(//n:transp/n:retTransp/n:vBCRet,'##.##.##0,00')"/>
                </span>
              </td>
              <td>
                <label>Alíquota</label>
                <span>
                  <xsl:call-template name="format4Casas">
                    <xsl:with-param name="num" select="//n:transp/n:retTransp/n:pICMSRet"/>
                  </xsl:call-template> 
                </span>
              </td>
            </tr>
            <tr>
              <td>
                <label>Valor ICMS</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:transp/n:retTransp/n:vICMSRet"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>CFOP</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:retTransp/n:CFOP"/>
                </span>
              </td>
              <td>
                <label>Município Ocor. Fato Gerador</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:retTransp/n:cMunFG"/>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
      </xsl:if>
      <xsl:variable name = "veiculo" select = "//n:transp/n:veicTransp"/>
      <xsl:if test = "$veiculo != ''">
        <fieldset>
          <legend class="titulo-aba-interna">Veículo</legend>
          <table class="box">
            <tr class="col-3">
              <td>
                <label>Placa</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:veicTransp/n:placa"/>
                </span>
              </td>
              <td>
                <label>UF</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:veicTransp/n:UF"/>
                </span>
              </td>
              <td>
                <label>RNTC</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:veicTransp/n:RNTC"/>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
      </xsl:if>
      <xsl:variable name = "reboque" select = "//n:transp/n:reboque"/>
      <xsl:if test = "$reboque != ''">
        <fieldset>
          <legend class="titulo-aba-interna">Reboque</legend>
          <table class="box">
            <tr class="col-3">
              <td class="fixo-transp-placa">
                <label>Placa</label>
                <span>
                  <xsl:value-of select = "child::node()[2]"/>
                  <xsl:value-of select="//n:transp/n:reboque/n:placa "/>
                </span>
              </td>
              <td class="fixo-transp-uf">
                <label>UF</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:reboque/n:UF"/>
                </span>
              </td>
              <td>
                <label>RNTC</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:reboque/n:RNTC"/>
                </span>
              </td>
            </tr> 
            <tr>
              <xsl:if test="//n:transp/n:reboque/n:vagao">
              <td>
                <label>Identificação do Vagão</label>
                <span>
                  <xsl:value-of select = "//n:transp/n:reboque/n:vagao"/>
                </span>
              </td>
              </xsl:if>
              <xsl:if test="//n:transp/n:reboque/n:balsa">
                <td>
                  <label>Identificação da Balsa</label>
                  <span>
                    <xsl:value-of select = "//n:transp/n:reboque/n:balsa"/>
                  </span>
                </td>
              </xsl:if>
            </tr> 
          </table>
        </fieldset>
      </xsl:if>
      <xsl:variable name = "volume" select = "//n:transp/n:vol"/>
      <xsl:if test = "$volume != ''">

        <fieldset>
          <legend class="titulo-aba-interna">Volumes</legend> 
          <xsl:for-each select="//n:transp/n:vol"> 
            <table class="box">
              <tr class="col-1">
                <td colspan="3" style="border-bottom:solid 1px #CCC">
                  <label>
                    Volume <xsl:value-of select = "position()"/>
                  </label>
                </td>
              </tr>
              <tr class="col-3">
                <td>
                  <label>Quantidade</label>
                  <span>
                    <xsl:value-of select = "n:qVol"/>
                  </span>
                </td>
                <td>
                  <label>Espécie</label>
                  <span>
                    <xsl:value-of select = "n:esp"/>
                  </span>
                </td>
                <td>
                  <label>Marca dos Volumes</label>
                  <span>
                    <xsl:value-of select = "n:marca"/>
                  </span>
                </td>
              </tr>
              <tr>
                <td>
                  <label>Numeração</label>
                  <span>
                    <xsl:value-of select = "n:nVol"/>
                  </span>
                </td>
                <td>
                  <label>Peso Líquido</label>
                  <span>
                    <xsl:call-template name="format3Casas">
                      <xsl:with-param name="num" select="n:pesoL"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td>
                  <label>Peso Bruto</label>
                  <span>
                    <xsl:call-template name="format3Casas">
                      <xsl:with-param name="num" select="n:pesoB"/>
                    </xsl:call-template>
                  </span>
                </td>
              </tr>
              <tr>
                <xsl:for-each select = "n:lacres/n:nLacre">
                  <td>
                    <label>
                      Número do Lacre <xsl:value-of select = " position()"/>
                    </label>
                    <span>
                      <xsl:value-of select = "text()"/>
                    </span>
                  </td>
                  <xsl:if test = "position() mod 3 = 0 ">
                    <tr>
                      
                    </tr>
                  </xsl:if>
                </xsl:for-each>
              </tr>
            </table> 
            <br/>
          </xsl:for-each>
        </fieldset>

      </xsl:if>
    </div>
  </xsl:template>
</xsl:stylesheet>