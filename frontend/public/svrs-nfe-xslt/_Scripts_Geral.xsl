<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:template match="SCRIPTS_GERAL" name="SCRIPTS_GERAL">
    <script type="text/javascript" src="/svrs-nfe-xslt/danfe-assets/danfe-runtime.js"></script>




    <!--<script language="javascript" src="http://localhost/apl/nfe/programas/Scripts/nfe-vis.js?rand=6545" />
    <script language="javascript" src="http://localhost/apl/nfe/programas/Scripts/jquery-1.2.6.pack.js?rand=3420248" />
    <script language="javascript" src="http://localhost/apl/nfe/programas/Scripts/jquery-1.2.6.min.js?rand=5849672" />
    <script language="javascript" src="http://localhost/apl/nfe/programas/Scripts/nfe-vis.js?rand=7854325" />-->
    
  </xsl:template>
</xsl:stylesheet>