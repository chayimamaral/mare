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
  <xsl:template match="/" name="Avulsa">
    <div id="Avulsa">
      <div class="GeralXslt"> 
        <fieldset>
          <legend class="titulo-aba">Dados de Nota Fiscal Avulsa</legend>
        <table class="box">
          <tr class="col-2">
            <xsl:for-each select="//n:avulsa/n:xOrgao">
              <td>
                <label>Órgão Emitente</label>
                <span>
                  <xsl:value-of select="text()"/>
                </span>
              </td>
            </xsl:for-each>
              <td class="fixo-cpf-cnpj">
                <label>CNPJ</label>
                <span>
                  <xsl:call-template name="formatCnpj">
                    <xsl:with-param name="cnpj" select="//n:avulsa/n:CNPJ"/>
                  </xsl:call-template>
                </span>
              </td>            
          </tr>
          <tr>
              <td>
                <label>Repartição Fiscal do Emitente</label>
                <span>
                  <xsl:value-of select="//n:avulsa/n:repEmi"/>
                </span>
              </td>
              <td>
                <label>Matrícula do Funcionário</label>
                <span>
                  <xsl:value-of select="//n:avulsa/n:matr"/>
                </span>
              </td>
          </tr>
          <tr> 
            <td>
              <label>Nome do Funcionário</label>
              <span>
                <xsl:value-of select="//n:avulsa/n:xAgente"/>
              </span>
            </td> 
            <td>
              <label>Fone / Fax</label>
              <span>
                <xsl:call-template name="formatFone">
                  <xsl:with-param name="fone" select="//n:avulsa/n:fone"/>
                </xsl:call-template>
              </span>
            </td>
            
          </tr>
          <tr>
            <td>
              <label>UF</label>
              <span>
                <xsl:value-of select="//n:avulsa/n:UF"/>
              </span>
            </td> 
            <td>
              <label>
                Número do Documento Arrecadação  
              </label>
              <span>
                <xsl:value-of select="//n:avulsa/n:nDAR"/>
              </span>
            </td>
          </tr>
          <tr>
              <td>
                <label>Valor Total do Documento Arrecadação</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:avulsa/n:vDAR"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <label>Data de Emissão do Documento Arrecadação</label>
                <span>
                  <xsl:call-template name="formatDate">
                    <xsl:with-param name="date" select="//n:avulsa/n:dEmi"/>
                  </xsl:call-template>
                </span>
              </td> 
          </tr>
          <tr>
              <td>
                <label>Data do Pagamento do Documento Arrecadação</label>
                <span>
                  <xsl:call-template name="formatDate">
                    <xsl:with-param name="date" select="//n:avulsa/n:dPag"/>
                  </xsl:call-template>
                </span>
              </td>            
          </tr>
        </table>
      </fieldset>
      </div>
    </div>
  </xsl:template>
</xsl:stylesheet>
