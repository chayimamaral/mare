<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:template match="ESTILOS_GERAL" name="ESTILOS_GERAL">
    <!-- Recursos locais: o host nfe-extranet.sefazvirtual.rs.gov.br costuma dar timeout fora da rede da SEFAZ. -->
    <link rel="stylesheet" type="text/css" href="/svrs-nfe-xslt/danfe-assets/xslt.css" />
  </xsl:template> 
</xsl:stylesheet>