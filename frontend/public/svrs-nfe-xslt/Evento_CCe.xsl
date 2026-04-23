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
  <xsl:template match="/" name="Evento_CCe">
    <xsl:variable name="ultCCe" select="(//n:retEvento/n:infEvento[n:tpEvento=110110]/n:nProt)[last()]"/>

    <xsl:for-each select="//n:eveNFe">
      <xsl:if test="n:evento/n:infEvento/n:tpEvento=110110">
        <xsl:if test="n:retEvento/n:infEvento/n:nProt=$ultCCe or not($ultCCe) ">
          <div id="CCe" class="GeralXslt">
          <fieldset>
            <legend class="titulo-aba">Carta de Correção</legend>
            <table class="box">
              <tr class="col-3">
                <td class="col-2" colspan="3">
                  <label>Orgão Recepção do Evento</label>
                    <span>
                      <xsl:value-of select="n:evento/n:infEvento/n:cOrgao"/> -
                      <xsl:call-template name="nomeUF">
                        <xsl:with-param select="n:evento/n:infEvento/n:cOrgao" name="uf">
                        </xsl:with-param>
                      </xsl:call-template>
                  </span>
                </td>
                <td class="col-4" colspan="3">
                  <label>Ambiente</label>
                  <span>
                    <xsl:value-of select="n:evento/n:infEvento/n:tpAmb"/>
                    <xsl:if test="n:evento/n:infEvento/n:tpAmb='1'"> - Produção</xsl:if>
                    <xsl:if test="n:evento/n:infEvento/n:tpAmb='2'"> - Homologação</xsl:if>
                  </span>
                </td>
                <td class="col-4" colspan="3">
                  <label>Versão</label>
                  <span>
                    <xsl:value-of select="n:evento/@versao"/>
                  </span>
                </td>
              </tr>
            </table >
            <br/>
            <table class="box">
              <tr class="col-3">
                <td colspan="3">
                  <xsl:variable name="cnpjEmit" select="n:evento/n:infEvento/n:CNPJ"/>
                  <xsl:variable name="cpfDest" select="n:evento/n:infEvento/n:CPF"/>
                  <label>Autor Evento (CNPJ / CPF)</label>
                  <span>
                    <xsl:call-template name="formatCnpj">
                      <xsl:with-param name="cnpj" select="$cnpjEmit"/>
                    </xsl:call-template>
                    <xsl:call-template name="formatCpf">
                      <xsl:with-param name="cpf" select="$cpfDest"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td colspan="3">
                  <label>Chave de Acesso</label>
                  <span>
                    <xsl:value-of select="n:evento/n:infEvento/n:chNFe"/>
                  </span>
                </td>
                <td colspan="3">
                  <label>Data Evento</label>
                  <span>
                    <xsl:call-template name="formatDateTimeFuso">
                      <xsl:with-param name="dateTime" select="n:evento/n:infEvento/n:dhEvento"/>
                      <xsl:with-param name="include_as" select="1"/>
                    </xsl:call-template>
                  </span>
                </td>
              </tr>
            </table >
            <br/>
            <table class="box">
              <tr class="col-2">
                <td colspan="2">
                  <label>Tipo de Evento</label>
                  <span>
                    <xsl:value-of select="n:evento/n:infEvento/n:tpEvento"/> - Carta de Correção
                  </span>
                </td>
                <td colspan="2">
                  <label>Sequencial do Evento</label>
                  <span>
                    <xsl:value-of select="n:evento/n:infEvento/n:nSeqEvento"/>
                  </span>
                </td>
              </tr>
            </table >
          </fieldset>
          <fieldset>
              <legend class="titulo-aba-interna">Detalhes do Evento</legend>
            <table class="box">
                <tr class="col-2">
                  <td colspan="2">
                    <label>Descrição do Evento</label>
                    <span>
                      <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:descEvento"/>
                    </span>
                  </td>
                  <td colspan="2">
                    <label>Versão</label>
                    <span>
                      <xsl:value-of select="n:evento/n:infEvento/n:verEvento"/>
                    </span>
                  </td>
                </tr>
              </table>
              <br/>
              <table class="box">
                <tr class="col-1">
                  <td colspan="1">
                    <label>Texto da Carta de Correção</label>
                    <span class="linhaCartaCorrecao">
                      <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:xCorrecao"/>
                    </span>
                  </td>
                </tr>
              </table>
            </fieldset>
          <fieldset>
              <legend class="titulo-aba">Autorização pela SEFAZ</legend>
            <table class="box">
                <tr>
                  <td class="col-2" colspan="3">
                    <label>Mensagem de Autorização</label>
                    <span>
                      <xsl:value-of select="n:retEvento/n:infEvento/n:cStat"/> -
                      <xsl:value-of select="n:retEvento/n:infEvento/n:xMotivo"/>
                    </span>
                  </td>
                  <td class="col-5" colspan="3">
                    <label>Protocolo</label>
                    <span>
                      <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                    </span>
                  </td>
                  <td class="col-3" colspan="3">
                    <label>Data/Hora Autorização</label>
                    <span>
                      <xsl:call-template name="formatDateTimeFuso">
                        <xsl:with-param name="dateTime" select="n:retEvento/n:infEvento/n:dhRegEvento"/>
                        <xsl:with-param name="include_as" select="1"/>
                      </xsl:call-template>
                    </span>
                  </td>
                </tr>
              </table>
            </fieldset>
          <fieldset>
          <legend class="titulo-aba">Condições de uso da Carta de Correção</legend>
          <table class="box">
              <tr class="col-1">
                <td colspan="1">
                  <span>
                    <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:xCondUso"/>
                  </span>
                </td>
              </tr>
            </table>
          </fieldset>        
      </div>
      </xsl:if>
    </xsl:if>
    </xsl:for-each>
  </xsl:template>
</xsl:stylesheet>
