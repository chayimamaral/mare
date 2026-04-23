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
  <xsl:template match="/" name="NFe_Cancelamento">
    <xsl:for-each select="//n:cancNFe">
          <div id="EventoCanc" class="GeralXslt" style="display:none">
            <fieldset>
              <legend class="titulo-aba">Cancelamento</legend>
              <table class="box">
                <tr>
                  <td class="col-2" colspan="3">
                    <label>Orgão Recepção do Evento</label>
                    <span>
                        <xsl:value-of select="//n:infNFe/n:ide/n:cUF"/> -
                        <xsl:call-template name="nomeUF">
                          <xsl:with-param select="//n:infNFe/n:ide/n:cUF" name="uf">
                          </xsl:with-param>
                        </xsl:call-template>                      
                    </span>
                  </td>
                  <td class="col-4" colspan="3">
                    <label>Ambiente</label>
                    <span>
                      <xsl:value-of select="//n:cancNFe/n:infCanc/n:tpAmb"/>
                      <xsl:if test="//n:cancNFe/n:infCanc/n:tpAmb='1'"> - Produção</xsl:if>
                      <xsl:if test="//n:cancNFe/n:infCanc/n:tpAmb='2'"> - Homologação</xsl:if>
                    </span>
                  </td>
                  <td class="col-4" colspan="3">
                    <label>Versão</label>
                    <span>
                      <xsl:value-of select="//n:cancNFe/@versao"/>
                    </span>
                  </td>
                </tr>
              </table >
              <br/>
              <table class="box">
                <tr class="col-3">
                  <td colspan="3">
                    <xsl:variable name="cnpjEmit" select="//n:infNFe/n:emit/n:CNPJ"/>
                    <xsl:variable name="cpfDest" select="//n:infNFe/n:emit/n:CPF"/>
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
                      <xsl:value-of select="//n:cancNFe/n:infCanc/n:chNFe"/>
                    </span>
                  </td>
                  <td colspan="3">
                    <label>Data Evento</label>
                    <span>
                      <xsl:call-template name="formatDateTimeFuso">
                        <xsl:with-param name="dateTime" select="//n:retCancNFe/n:infCanc/n:dhRecbto"/>
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
                    <span>110111 - Cancelamento pelo Emitente</span>
                  </td>
                  <td colspan="2">
                    <label>Sequencial do Evento</label>
                    <span>0</span>
                  </td>
                </tr>
              </table>
            </fieldset>
            <fieldset>
                <legend class="titulo-aba-interna">Detalhes do Evento</legend>
                <table class="box">
                  <tr class="col-2">
                    <td colspan="2">
                      <label>Descrição do Evento</label>
                      <span>
                        Cancelamento
                      </span>
                    </td>
                    <td colspan="2">
                      <label>Versão</label>
                      <span>
                        <xsl:value-of select="//n:cancNFe/@versao"/>
                      </span>
                    </td>
                  </tr>
                </table>
                <br/>
                <table class="box">
                  <tr class="col-2">
                    <td colspan="2">
                      <label>Justificativa do Cancelamento</label>
                      <span>
                        <xsl:value-of select="//n:cancNFe/n:infCanc/n:xJust"/>
                      </span>
                    </td>
                    <td colspan="2">
                      <label>Protocolo da NF-e</label>
                      <span>
                        <xsl:value-of select="//n:cancNFe/n:infCanc/n:nProt"/>
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
                        <xsl:value-of select="//n:retCancNFe/n:infCanc/n:cStat"/> -
                        <xsl:value-of select="//n:retCancNFe/n:infCanc/n:xMotivo"/>
                      </span>
                    </td>
                    <td class="col-5" colspan="3">
                      <label>Protocolo</label>
                      <span>
                        <xsl:value-of select="//n:retCancNFe/n:infCanc/n:nProt"/>
                      </span>
                    </td>
                    <td class="col-3" colspan="3">
                      <label>Data/Hora Autorização</label>
                      <span>
                        <xsl:call-template name="formatDateTimeFuso">
                          <xsl:with-param name="dateTime" select="//n:retCancNFe/n:infCanc/n:dhRecbto"/>
                          <xsl:with-param name="include_as" select="1"/>
                        </xsl:call-template>
                      </span>
                    </td>
                  </tr>
                </table>
              </fieldset>
          </div>
    </xsl:for-each>
  </xsl:template>
</xsl:stylesheet>
