<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
xmlns:fo="http://www.w3.org/1999/XSL/Format"
xmlns:n="http://www.portalfiscal.inf.br/nfe"
xmlns:s="http://www.w3.org/2000/09/xmldsig#"
version="2.0"
exclude-result-prefixes="fo n s">
 
  <xsl:import href="Evento_CCe.xsl"/>
  <xsl:import href="Evento_Cancelamento.xsl"/>
  <xsl:import href="NFe_Cancelamento.xsl"/>
  <xsl:import href="Evento_EPEC.xsl"/>
  <xsl:import href="Evento_Confirmacao.xsl"/>
  <xsl:import href="Evento_Ciencia.xsl"/>
  <xsl:import href="Evento_Desconhecimento.xsl"/>
  <xsl:import href="Evento_Nao_Realizado.xsl"/>
  <xsl:import href="Evento_CTe_Autorizado.xsl"/>
  <xsl:import href="Evento_CTe_Cancelado.xsl"/>
  <xsl:import href="Evento_Vistoria_SUFRAMA.xsl"/>
  <xsl:import href="Evento_Internalizacao_SUFRAMA.xsl"/>
  <xsl:import href="Evento_Registro_Passagem.xsl"/>
  <xsl:import href="Evento_Registro_Passagem_BRId.xsl"/>
  <xsl:import href="Evento_MDFe_Autorizado.xsl"/>
  <xsl:import href="Evento_MDFe_Cancelado.xsl"/>
  <xsl:import href="Evento_Cancelamento_Registro_Passagem.xsl"/>

  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:output method="html"/>

  <xsl:include href="Utils.xsl"/>
  <xsl:include href="_Estilos_Geral.xsl"/>
  <xsl:include href="_Scripts_Geral.xsl"/>
  <xsl:template match="/" name="VisualizacaoEvento">
   
        <link rel="stylesheet" type="text/css" href="/svrs-nfe-xslt/danfe-assets/xslt.css" />
        <xsl:for-each select="n:eveNFe">
          <xsl:variable name="tipo" select="n:retEvento/n:infEvento/n:tpEvento" />
          
            <xsl:choose>
                  <xsl:when test="$tipo = '110110'"> 
                    <xsl:call-template name="Evento_CCe"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '110111'">
                    <xsl:call-template name="Evento_Cancelamento"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '110140'">
                    <xsl:call-template name="Evento_EPEC"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '210200'">
                    <xsl:call-template name="Evento_Confirmacao"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '210210'">
                    <xsl:call-template name="Evento_Ciencia"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '210220'">
                    <xsl:call-template name="Evento_Desconhecimento"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '210240'">
                    <xsl:call-template name="Evento_Nao_Realizado"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '610600'">
                    <xsl:call-template name="Evento_CTe_Autorizado"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '610601'">
                    <xsl:call-template name="Evento_CTe_Cancelado"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '990900'">
                    <xsl:call-template name="Evento_Vistoria_SUFRAMA"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '990910'">
                    <xsl:call-template name="Evento_Internalizacao_SUFRAMA"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '610500'">
                    <xsl:call-template name="Evento_Registro_Passagem"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '610550'">
                    <xsl:call-template name="Evento_Registro_Passagem_BRId"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '610501'">
                    <xsl:call-template name="Evento_Cancelamento_Registro_Passagem"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '610610'">
                    <xsl:call-template name="Evento_MDFe_Autorizado"/>
                  </xsl:when>
                  <xsl:when test="$tipo = '610611'">
                    <xsl:call-template name="Evento_MDFe_Cancelado"/>
                  </xsl:when>
            </xsl:choose>
          </xsl:for-each>  
  </xsl:template> 
</xsl:stylesheet>
