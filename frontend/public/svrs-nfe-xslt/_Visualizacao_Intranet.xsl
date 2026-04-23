<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  
  <!-- Abas superiores--> 
  <xsl:import href="NFe.xsl"/>
  <xsl:import href="Emitente.xsl"/>
  <xsl:import href="DestinatarioRemetente.xsl"/>
  <xsl:import href="ProdutosServicos.xsl"/>
  <xsl:import href="Totais.xsl"/>
  <xsl:import href="Transporte.xsl"/>
  <xsl:import href="InformacoesAdicionais.xsl"/>
  <xsl:import href="Cobranca.xsl"/>
  <xsl:import href="Avulsa.xsl"/>
  
  <!--Eventos-->
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
  <xsl:import href="Evento_Registro_Passagem_Sistema_Externo.xsl"/>
  <xsl:import href="Evento_MDFe_Autorizado.xsl"/>
  <xsl:import href="Evento_MDFe_Cancelado.xsl"/>  
  <xsl:import href="Evento_Registro_Passagem_BRId.xsl"/>
  <xsl:import href="Evento_Cancelamento_Registro_Passagem.xsl"/>

  <!--Configuração -->
  <xsl:import href="_Estilos_Geral.xsl"/>
  <xsl:import href="_Scripts_Geral.xsl"/>
  
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:output method="html" indent="no"/>
  <xsl:template match="/">
    <xsl:call-template name="ESTILOS_GERAL"/>
    <xsl:call-template name="SCRIPTS_GERAL"/>
    <xsl:call-template name="NFe">
      <xsl:with-param name="ambiente" select="'intranet'"/>
    </xsl:call-template>
    <xsl:call-template name="Emitente"/>
    <xsl:call-template name="Destinatario"/>
    <xsl:call-template name="Produtos_e_Servicos"/>
    <xsl:call-template name="Totais"/>
    <xsl:call-template name="Transporte"/>
    <xsl:call-template name="Cobranca"/>
    <xsl:call-template name="Informacoes_Adicionais"/>
    <xsl:call-template name="NFe_Cancelamento"/>
    <xsl:for-each select ="//n:avulsa">
        <xsl:call-template name="Avulsa"/>
    </xsl:for-each> 
    
  </xsl:template>
</xsl:stylesheet>