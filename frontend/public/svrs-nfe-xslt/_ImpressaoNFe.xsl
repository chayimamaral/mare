<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
	<xsl:import href="NFe.xsl"/>
	<xsl:import href="Emitente.xsl"/>
	<xsl:import href="Destinatario.xsl"/>
	<xsl:import href="Produtos_e_Servicos.xsl"/>
	<xsl:import href="Totais.xsl"/>
	<xsl:import href="Transporte.xsl"/>
	<xsl:import href="Informacoes_Adicionais.xsl"/>
	<xsl:import href="Cobranca.xsl"/>
  <xsl:import href="Evento_CCe.xsl"/>
	<xsl:import href="_Estilos_Geral.xsl"/>
	<xsl:import href="_Scripts_Geral.xsl"/>
	<xsl:import href="Avulsa.xsl"/>
	<xsl:decimal-format decimal-separator="," grouping-separator="."/>
	<xsl:output method="html"/>
	<xsl:template match="/">
		<xsl:call-template name="HEADER_IMPRESSAO" >
			<xsl:with-param name="tipo">simples</xsl:with-param>
		</xsl:call-template>
		<xsl:call-template name ="ESTILOS_GERAL"/>
		<xsl:call-template name ="SCRIPTS_GERAL"/>
		<xsl:call-template name="NFe"/>
		<xsl:call-template name="Emitente"/>
		<xsl:call-template name="Destinatario"/>
		<xsl:call-template name="Produtos_e_Servicos"/>
		<xsl:call-template name="Totais"/>
		<xsl:call-template name="Transporte"/>
		<xsl:call-template name="Cobranca"/>
		<xsl:call-template name="Informacoes_Adicionais"/>
		<xsl:call-template name="Evento_CCe"/>
    <xsl:call-template name="Avulsa"/>
	</xsl:template>
</xsl:stylesheet>