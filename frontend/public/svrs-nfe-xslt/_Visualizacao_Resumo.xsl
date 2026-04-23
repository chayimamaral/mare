<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
	<xsl:import href="NFe.xsl"/>
	<xsl:import href="_Estilos_Geral.xsl"/>
	<xsl:import href="_Scripts_Geral.xsl"/>
	<xsl:decimal-format decimal-separator="," grouping-separator="."/>
	<xsl:output method="html"/>
	<xsl:template match="/">
		<xsl:call-template name ="ESTILOS_GERAL"/>
		<xsl:call-template name ="SCRIPTS_GERAL"/>
		<xsl:call-template name="NFe"/>
		</xsl:template>
</xsl:stylesheet>