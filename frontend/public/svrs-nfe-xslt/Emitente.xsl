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
  <xsl:template match="/" name="Emitente">
    <div id="Emitente" class="GeralXslt">
      <fieldset>
        <legend class="titulo-aba">Dados do Emitente</legend>
        <table class="box">
          <tr class="col-2">
            <td>
              <label>Nome / Razão Social</label>
              <span>
                <xsl:value-of select = "//n:emit/n:xNome"/>
              </span>
            </td>
            <td>
              <label>Nome Fantasia</label>
              <span>
                <xsl:value-of select = "//n:emit/n:xFant"/>
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <xsl:variable name="cnpj" select="//n:emit/n:CNPJ"/>
              <xsl:variable name="cpf" select="//n:emit/n:CPF"/> 
              <xsl:choose>
                <xsl:when test="$cnpj!=''">
                  <label>CNPJ</label>
                  <span>
                      <xsl:call-template name="formatCnpj">
                         <xsl:with-param name="cnpj" select="$cnpj"/>
                      </xsl:call-template>
                  </span>
                </xsl:when>
                <xsl:when test="$cpf!=''">
                  <label>CPF</label>
                  <span> 
                    <xsl:call-template name="formatCpf">
                        <xsl:with-param name="cpf" select="$cpf"/>
                    </xsl:call-template>
                  </span>
                </xsl:when> 
              </xsl:choose> 
            </td>
            <td>
              <label>Endereço</label>
              <span>
                <xsl:variable name="endE" select="//n:emit/n:enderEmit/n:xLgr"/>
                <xsl:if test="$endE != ''">
                  <xsl:value-of select="$endE"/>,&#160;
                </xsl:if>
                <xsl:value-of select = "//n:emit/n:enderEmit/n:nro"/>&#160;
                <xsl:value-of select = "//n:emit/n:enderEmit/n:xCpl"/>
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <label>Bairro / Distrito</label>
              <span>
                <xsl:value-of select = "//n:emit/n:enderEmit/n:xBairro"/>
              </span>
            </td>
            <td>
              <label>CEP</label>
              <span>
                <xsl:call-template name="formatCep">
                  <xsl:with-param name="cep" select="//n:emit/n:enderEmit/n:CEP"/>
                </xsl:call-template>
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <label>Município</label>
              <span>
                <xsl:value-of select="//n:emit/n:enderEmit/n:cMun"/>
                -
                <xsl:value-of select = "//n:emit/n:enderEmit/n:xMun"/>
              </span>
            </td>
            <td>
              <label>Telefone</label>
              <span>
                <xsl:call-template name="formatFone">
                  <xsl:with-param name="fone" select="//n:emit/n:enderEmit/n:fone"/>
                </xsl:call-template>
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <label>UF</label>
              <span>
                <xsl:value-of select = "//n:emit/n:enderEmit/n:UF"/>
              </span>
            </td>
            <td>
              <label>País</label>
              <span>
                <xsl:for-each select="//n:emit/n:enderEmit/n:cPais">
                  <xsl:variable name="cPaisE" select="//n:emit/n:enderEmit/n:cPais"/>
                  <xsl:if test="$cPaisE != ''">
                    <xsl:value-of select="$cPaisE"/>
                    -
                  </xsl:if>
                </xsl:for-each>
                
                <xsl:value-of select = "//n:emit/n:enderEmit/n:xPais"/>                
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <label>Inscrição Estadual</label>
              <span>
                <xsl:value-of select = "//n:emit/n:IE"/>
              </span>
            </td>
            <td>
              <label>Inscrição Estadual do Substituto Tributário</label>
              <span>
                <xsl:value-of select = "//n:emit/n:IEST"/>
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <label>Inscrição Municipal</label>
              <span>
                <xsl:value-of select = "//n:emit/n:IM"/>
              </span>
            </td>
            <td>
              <label>Município da Ocorrência do Fato Gerador do ICMS</label>
              <span>
                <xsl:value-of select="//n:ide/n:cMunFG"/>
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <label>CNAE Fiscal</label>
              <span>
                <xsl:value-of select = "//n:emit/n:CNAE"/>
              </span>
            </td>
            <td>
              <label>Código de Regime Tributário</label>
              <span>
                <xsl:for-each select="//n:emit/n:CRT">
                  <xsl:variable name="crt" select="text()"/>
                  <xsl:choose>
                    <xsl:when test="$crt = 1">
                      1 - Simples Nacional
                    </xsl:when>
                    <xsl:when test="$crt = 2">
                      2 - Simples Nacional - excesso de sublimite de receita bruta
                    </xsl:when>
                    <xsl:when test="$crt = 3">
                      3 - Regime Normal
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select = "$crt"/>
                    </xsl:otherwise>
                  </xsl:choose> 
                </xsl:for-each>
              </span>
            </td>
          </tr>
        </table>
      </fieldset>
    </div>
  </xsl:template>
</xsl:stylesheet>