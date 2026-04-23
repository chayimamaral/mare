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
  <xsl:template match="/" name="Destinatario">
  <xsl:variable name="op" select="//n:ide/n:tpNF"/>
    <div id="DestRem" op="{$op}" class="GeralXslt">
        <fieldset>
          <legend class="titulo-aba">Dados do <xsl:if test="$op = 0">Remetente</xsl:if><xsl:if test="$op = 1">Destinatário</xsl:if></legend>
          <table class="box">
            <tr class="col-2">
              <td colspan="3">
                <label>Nome / Razão Social</label>
                <span>
                  <xsl:value-of select = "//n:dest/n:xNome"/>
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td  colspan="2">
                <xsl:variable name="cnpjDest" select="//n:dest/n:CNPJ"/>
                <xsl:variable name="cpfDest" select="//n:dest/n:CPF"/>
                <xsl:variable name="idestrang" select="//n:infNFe/n:dest/n:idEstrangeiro"/>
                <xsl:choose>
                  <xsl:when test="$cnpjDest!=''">
                    <label>CNPJ</label>
                  </xsl:when>
                  <xsl:when test="$cpfDest!=''">
                    <label>CPF</label>
                  </xsl:when>
                  <xsl:when test="$idestrang!=''" >
                    <label>Id. Estrangeiro</label>
                  </xsl:when>
                  <xsl:otherwise>
                    <label>CNPJ/CPF/Id. Estrangeiro</label>
                  </xsl:otherwise>
                </xsl:choose>
                <span>
                  <xsl:call-template name="formatCnpj">
                    <xsl:with-param name="cnpj" select="$cnpjDest"/>
                  </xsl:call-template>
                  <xsl:call-template name="formatCpf">
                    <xsl:with-param name="cpf" select="$cpfDest"/>
                  </xsl:call-template>
                  <xsl:value-of select="$idestrang"/>
                </span>
              </td>
              <td>
                <label>Endereço</label>
                <span>
                  <xsl:variable name="endD" select="//n:dest/n:enderDest/n:xLgr"/>
                  <xsl:if test="$endD != ''">
                    <xsl:value-of select="$endD"/>,&#160;
                  </xsl:if>
                  <xsl:value-of select = "//n:dest/n:enderDest/n:nro"/>&#160;
                  <xsl:value-of select = "//n:dest/n:enderDest/n:xCpl"/>
                </span>
              </td>
            </tr>
            <tr>
              <td colspan="2">
                <label>Bairro / Distrito</label>
                <span>
                  <xsl:value-of select = "//n:dest/n:enderDest/n:xBairro"/>
                </span>
              </td>
              <td>
                <label>CEP</label>
                <span>
                  <xsl:for-each select="//n:dest/n:enderDest/n:CEP">
                    <xsl:call-template name="formatCep">
                      <xsl:with-param name="cep" select="text()"/>
                    </xsl:call-template>
                  </xsl:for-each>
                </span>
              </td>
            </tr>
            <tr>
              <td colspan="2">
                <label>Município</label>
                <span>
                  <xsl:value-of select="//n:dest/n:enderDest/n:cMun"/>
                  <xsl:if test="//n:dest/n:enderDest/n:xMun != ''">
                    -
                  </xsl:if>
                  <xsl:value-of select="//n:dest/n:enderDest/n:xMun"/>
                </span>
              </td>
              <td>
                <label>Telefone</label>
                <span>
                  <xsl:for-each select="//n:dest/n:enderDest/n:fone">
                    <xsl:call-template name="formatFone">
                      <xsl:with-param name="fone" select="text()"/>
                    </xsl:call-template>
                  </xsl:for-each>
                </span>
              </td>
            </tr>
            <tr>
              <td colspan="2">
                <label>UF</label>
                <span>
                  <xsl:value-of select = "//n:dest/n:enderDest/n:UF"/>
                </span>
              </td>
              <td>
                <label>País</label>
                <span>
                  <xsl:for-each select="//n:dest/n:enderDest/n:xPais">
                    <xsl:variable name="cPaisE" select="//n:dest/n:enderDest/n:cPais"/>
                    <xsl:if test="$cPaisE != ''">
                      <xsl:value-of select="$cPaisE"/>
                      -
                    </xsl:if>
                    <xsl:value-of select = "//n:dest/n:enderDest/n:xPais"/>
                  </xsl:for-each>
                </span>
              </td>
            </tr>
            <tr>
              <td style="width:25%">
                <label>Indicador IE</label>
                <span>
                  <xsl:variable name="indIeDest" select="//n:dest/n:indIEDest"/>
                  <xsl:choose>
                    <xsl:when test="$indIeDest = 1">
                      01 - Contribuinte ICMS (informar a IE do destinatário)
                    </xsl:when>
                    <xsl:when test="$indIeDest = 2">
                      02 - Contribuinte isento de Inscrição no cadastro de Contribuintes do ICMS
                    </xsl:when>
                    <xsl:when test="$indIeDest = 9">
                      09 - Não Contribuinte, que pode ou não possuir Inscrição Estadual no Cadastro de Contribuintes do ICMS
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select = "$indIeDest"/>
                    </xsl:otherwise>
                  </xsl:choose> 
                </span>
              </td>
              <td style="width:25%">
                <label>Inscrição Estadual</label>
                <span>
                  <xsl:value-of select = "//n:dest/n:IE"/>
                </span>
              </td>            
              <td  style="width:50%">
                <label>Inscrição SUFRAMA</label>
                <span>
                  <xsl:value-of select = "//n:dest/n:ISUF"/>
                </span>
              </td>
            </tr>
            <tr>
              <td>
                <label>IM</label>
                <span>
                  <xsl:value-of select = "//n:dest/n:IM"/>
                </span>
              </td>
              <td colspan="2">
                <label>E-mail</label>
                <span>
                  <xsl:value-of select = "//n:dest/n:email"/>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>      
      <xsl:variable name="retirada" select="//n:retirada"/>
      <xsl:if test="$retirada != ''">
        <fieldset>
          <legend class="titulo-aba">Local de Retirada</legend>
          <table class="box">
            <tr>
              <td style="width: 35%">
                <xsl:variable name="cnpjRet" select="//n:retirada/n:CNPJ"/>
                <xsl:variable name="cpfRet" select="//n:retirada/n:CPF"/>
                <label>
                  <xsl:if test="$cnpjRet != ''">CNPJ</xsl:if>
                  <xsl:if test="$cpfRet != ''">CPF</xsl:if> 
                </label>
                <span>
                  <xsl:call-template name="formatCnpj">
                    <xsl:with-param name="cnpj" select="$cnpjRet"/>
                  </xsl:call-template>
                  <xsl:call-template name="formatCpf">
                    <xsl:with-param name="cpf" select="$cpfRet"/>
                  </xsl:call-template>
                </span>
              </td>
              <td colspan="2">
                <label>Logradouro</label>
                <span>
                  <xsl:variable name="endRet" select="//n:retirada/n:xLgr"/>
                  <xsl:if test="$endRet != ''">
                    <xsl:value-of select="$endRet"/>
                    ,&#160;
                  </xsl:if>
                  <xsl:value-of select="//n:retirada/n:nro"/>
                  &#160;
                  <xsl:value-of select="//n:retirada/n:xCpl"/>
                </span>
              </td>
            </tr>
            <tr>
              <td>
                <label>Bairro</label>
                <span>
                  <xsl:value-of select="//n:retirada/n:xBairro"/>
                </span>
              </td>
              <td>
                <label>Município</label>
                <span>
                  <xsl:variable name="mun1" select="//n:retirada/n:cMun"/>
                  <xsl:if test="$mun1 != ''">
                    <xsl:value-of select="$mun1"/>
                    -
                  </xsl:if>
                  <xsl:value-of select="//n:retirada/n:xMun"/>
                </span>
              </td>
              <td style="width: 7%; min-width: 45px;">
                <label>UF</label>
                <span>
                  <xsl:value-of select="//n:retirada/n:UF"/>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
      </xsl:if>
      <xsl:variable name="entrega" select="//n:entrega"/>
      <xsl:if test="$entrega != ''">
        <fieldset>
          <legend class="titulo-aba">Local de Entrega</legend>
          <table class="box">
            <tr>
              <td style="width: 35%;">
                <xsl:variable name="cnpjEnt" select="//n:entrega/n:CNPJ"/>
                <xsl:variable name="cpfEnt" select="//n:entrega/n:CPF"/>
                <label>
                  <xsl:if test="$cnpjEnt != ''">CNPJ</xsl:if>
                  <xsl:if test="$cpfEnt != ''">CPF</xsl:if>
                </label>
                <span>
                  <xsl:call-template name="formatCnpj">
                    <xsl:with-param name="cnpj" select="$cnpjEnt"/>
                  </xsl:call-template>
                  <xsl:call-template name="formatCpf">
                    <xsl:with-param name="cpf" select="$cpfEnt"/>
                  </xsl:call-template>
                </span>
              </td>
              <td colspan="2">
                <label>Logradouro</label>
                <span>
                  <xsl:variable name="endEnt" select="//n:entrega/n:xLgr"/>
                  <xsl:if test="$endEnt != ''">
                    <xsl:value-of select="$endEnt"/>
                    ,&#160;
                  </xsl:if>
                  <xsl:value-of select="//n:entrega/n:nro"/>
                  &#160;
                  <xsl:value-of select="//n:entrega/n:xCpl"/>
                </span>
              </td>
            </tr>
            <tr>
              <td>
                <label>Bairro</label>
                <span>
                  <xsl:value-of select="//n:entrega/n:xBairro"/>
                </span>
              </td>
              <td>
                <label>Município</label>
                <span>
                  <xsl:variable name="mun2" select="//n:entrega/n:cMun"/>
                  <xsl:if test="$mun2 != ''">
                    <xsl:value-of select="$mun2"/>
                    -
                  </xsl:if>
                  <xsl:value-of select="//n:entrega/n:xMun"/>
                </span>
              </td>
              <td style="width: 7%; min-width: 45px;">
                <label>UF</label>
                <span>
                  <xsl:value-of select="//n:entrega/n:UF"/>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
      </xsl:if>
    </div>
  </xsl:template>
</xsl:stylesheet>
