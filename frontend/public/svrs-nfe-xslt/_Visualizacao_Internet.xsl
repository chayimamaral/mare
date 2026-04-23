<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <xsl:import href="NFe.xsl"/>
  <xsl:import href="Emitente.xsl"/>
  <xsl:import href="DestinatarioRemetente.xsl"/>
  <xsl:import href="ProdutosServicos.xsl"/>
  <xsl:import href="Totais.xsl"/>
  <xsl:import href="Transporte.xsl"/>
  <xsl:import href="InformacoesAdicionais.xsl"/>
  <xsl:import href="Cobranca.xsl"/>
  <!-- Evento_*.xsl já são importados por NFe.xsl; repetir aqui só duplica Utils/estilos e dispara SXWN9019. -->

  <xsl:import href="Avulsa.xsl"/>
  <xsl:import href="_Estilos_Geral.xsl"/>
  <xsl:import href="_Scripts_Geral.xsl"/>
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:output method="html" indent="no"/>
  <xsl:template match="/">
    <html xmlns="">
      <xsl:call-template name ="ESTILOS_GERAL"/>
      <xsl:call-template name ="SCRIPTS_GERAL"/>
     
      <xsl:variable name="nota_cancelada" select="//n:cancNFe/n:infCanc/n:xServ"/>
      <xsl:variable name="nota_denegada" select="//n:infProt/n:cStat"/>
      
      <xsl:if  test="$nota_denegada = 301 or $nota_denegada = 302">
        <script language="javascript">
          jQuery(document).ready(function () {
            alert('Nota Fiscal denegada pela SEFAZ')
              jQuery('li[id^="tab_"]').click(function () {
                alert('Nota Fiscal denegada pela SEFAZ');
              });
          });
        </script>
      </xsl:if>
      <xsl:if  test="$nota_cancelada!=''">
        <script language="javascript">
          jQuery(document).ready(function () {
            alert('Nota Fiscal cancelada pelo emitente')
              jQuery('li[id^="tab_"]').click(function () {
                alert('Nota Fiscal cancelada pelo emitente');
              });
          });
        </script>
      </xsl:if>

      <body>
        <xsl:call-template name ="CABECALHO_NFE"/> 
        <br /> 
        <ul id="botoes_nft">
          <li class="nftselected" id="tab_0" onclick="mostraAba(0)">
            <b>
              NFe
            </b>
          </li>
          <li id="tab_1" onclick="mostraAba(1)">
            <b>Emitente</b>
          </li>
          <li id="tab_2" onclick="mostraAba(2)">
             <b>Destinatário</b>
          </li>
          <li id="tab_3" onclick="mostraAba(3)">
            <b>Produtos&#160;e&#160;Serviços</b>
          </li>
          <li id="tab_4" onclick="mostraAba(4)">
            <b>Totais</b>
          </li>
          <li id="tab_5" onclick="mostraAba(5)">
            <b>Transporte</b>
          </li>
          <li id="tab_6" onclick="mostraAba(6)">
            <b>Cobrança</b>
          </li>
          <li id="tab_7" onclick="mostraAba(7)">
            <b>Informações&#160;Adicionais</b>
          </li>

          <xsl:for-each select ="//n:avulsa">
            <li id="tab_8" onclick="mostraAba(8)">
              <b>Avulsa</b>
            </li>
          </xsl:for-each> 
        </ul> 
        <div class="aba_container" style="clear:both;">
          <div id="aba_nft_0" class="nft" style="display:block; ">
            <xsl:call-template name="NFe">
              <xsl:with-param name="ambiente" select="'publico'"/>
            </xsl:call-template>
          </div>
          <div id="aba_nft_1" class="nft">
            <xsl:call-template name="Emitente"/>
          </div>
          <div id="aba_nft_2" class="nft">
            <xsl:call-template name="Destinatario"/>
          </div>
          <div id="aba_nft_4" class="nft">
            <xsl:call-template name="Totais"/>
          </div>
          <div id="aba_nft_5" class="nft">
            <xsl:call-template name="Transporte"/>
          </div>
          <div id="aba_nft_6" class="nft">
            <xsl:call-template name="Cobranca"/>
          </div>
          <div id="aba_nft_7" class="nft">
            <xsl:call-template name="Informacoes_Adicionais"/>
          </div>
          <xsl:for-each select ="//n:avulsa">
            <div id="aba_nft_8" class="nft">
              <xsl:call-template name="Avulsa"/>
            </div>
          </xsl:for-each>
          
          <xsl:for-each select ="//n:evento">
            <div id="aba_nft_9" class="nft">
              <xsl:call-template name="Evento_CCe"/>
            </div>
          </xsl:for-each>
          
          <div id="aba_nft_3" class="nft">
            <xsl:call-template name="Produtos_e_Servicos"/>
          </div>
          <div id="aba_nft_10" class="nft">
            <xsl:call-template name="Evento_Cancelamento"/>
            <xsl:call-template name="NFe_Cancelamento"/>
          </div>
        </div> 
      </body>
    </html> 
  </xsl:template>
</xsl:stylesheet>