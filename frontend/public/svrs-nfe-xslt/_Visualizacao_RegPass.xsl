<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
	<xsl:import href="NFe.xsl"/>
  <xsl:import href="ProdutosServicos.xsl"/>
  <xsl:import href="RegPass_Dpec.xsl"/>
	<xsl:import href="_Estilos_Geral.xsl"/>
	<xsl:import href="_Scripts_Geral.xsl"/>
	<xsl:decimal-format decimal-separator="," grouping-separator="."/>
	<xsl:output method="html"/>

  <xsl:param name="tpAmb"/>
 
	<xsl:template match="/">
		<xsl:call-template name="ESTILOS_GERAL"/>
		<xsl:call-template name="SCRIPTS_GERAL"/>

    <xsl:variable name="temNFe" select="count(//n:NFe)"/>
    <xsl:for-each select ="//n:NFe">
      <xsl:call-template name="NFe"/>
      <xsl:call-template name="Produtos_e_Servicos"/>
    </xsl:for-each>

    <xsl:variable name="temDpec" select="count(//n:nfeDpec)"/>
    <xsl:for-each select ="//n:nfeDpec">
      <xsl:call-template name="nfeDpec">
        <xsl:with-param name="tpAmb" select="$tpAmb" />
      </xsl:call-template>
    </xsl:for-each>

    <xsl:choose>
      <xsl:when test="$temNFe=0 and $temDpec=0">
        <div id="NFe">
          <fieldset>
            <legend class="titulo-aba">NF-e/DPEC não encontrada no Ambiente Nacional</legend>
          </fieldset>
        </div>
      </xsl:when>
    </xsl:choose>    
    
	</xsl:template>
</xsl:stylesheet>