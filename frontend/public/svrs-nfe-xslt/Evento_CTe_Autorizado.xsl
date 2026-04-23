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
  <xsl:template match="/" name="Evento_CTe_Autorizado">

    <xsl:variable name="prot" select="n:retEvento/n:infEvento/n:nProt" />
    <div id="EventoCTeAutorizado{$prot}" class="GeralXslt">
      <fieldset>
        <legend class="titulo-aba">CT-e Autorizado</legend>
        <table class="box">
          <tr>
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
        </table>
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
        </table>
        <br/>
        <table class="box">
          <tr class="col-2">
            <td colspan="2">
              <label>Tipo de Evento</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:tpEvento"/> - <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:descEvento"/>
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
        <legend>Detalhes do Evento</legend>
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
          <tr class="col-2">
            <td colspan="2">
              <label>Chave de Acesso CT-e</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:CTe/n:chCTe"/>
              </span>
            </td>
            <td colspan="2">
              <label>Protocolo do CT-e</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:CTe/n:nProt"/>
              </span>
            </td>
          </tr>
          <tr class="col-2">
            <td colspan="2">
              <label>Data de emissão do CT-e</label>
              <span>

                <xsl:call-template name="formatDateTimeFuso">
                  <xsl:with-param name="dateTime" select="n:evento/n:infEvento/n:detEvento/n:CTe/n:dhEmi"/>
                  <xsl:with-param name="include_as" select="1"/>
                </xsl:call-template>
              </span>
            </td>
            <td colspan="2">
              <label>Data da autorização do CT-e</label>
              <span>
                <xsl:call-template name="formatDateTimeFuso">
                  <xsl:with-param name="dateTime" select="n:evento/n:infEvento/n:detEvento/n:CTe/n:dhRecbto"/>
                  <xsl:with-param name="include_as" select="1"/>
                </xsl:call-template>
              </span>
            </td>
          </tr>
          <tr class="col-2">
            <td colspan="2">
              <label>Modal</label>
              <span>
                <xsl:variable name="modalidade" select="n:evento/n:infEvento/n:detEvento/n:CTe/n:modal"/>
                <xsl:choose>
                  <xsl:when test="$modalidade = '01'">
                    01 - Rodoviário
                  </xsl:when>
                  <xsl:when test="$modalidade = '02'">
                    02 - Aéreo
                  </xsl:when>
                  <xsl:when test="$modalidade = '03'">
                    03 - Aquaviário
                  </xsl:when>
                  <xsl:when test="$modalidade = '04'">
                    04 - Ferroviário
                  </xsl:when>
                  <xsl:when test="$modalidade = '05'">
                    05 - Dutoviário
                  </xsl:when>
                  <xsl:when test="$modalidade = '06'">
                    06 - Multimodal
                  </xsl:when>
                  <xsl:otherwise>
                    <xsl:value-of select="$modalidade"/>
                  </xsl:otherwise>
                </xsl:choose>
              </span>
            </td>
          </tr>
        </table>
      </fieldset>
      <fieldset>
        <legend>Emitente CT-e</legend>
        <table class="box">
          <tr class="col-2">
            <td colspan="2">
              <label>CNPJ</label>
              <span>
                <xsl:call-template name="formatCnpj">
                  <xsl:with-param name="cnpj" select="n:evento/n:infEvento/n:detEvento/n:emit/n:CNPJ"/>
                </xsl:call-template>
              </span>
            </td>
            <td colspan="2">
              <label>IE</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:emit/n:IE"/>
              </span>
            </td>
          </tr>
          <tr class="col-2">
            <td colspan="2">
              <label>Nome do Emitente</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:emit/n:xNome"/>
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
                <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/> <br/>
              </span>
            </td>
            <td class="col-3" colspan="3">
              <label>Data/Hora Autorização</label>
              <span>
                <xsl:call-template name="formatDateTimeFuso">
                  <xsl:with-param name="dateTime" select="n:retEvento/n:infEvento[n:tpEvento=610600]/n:dhRegEvento"/>
                  <xsl:with-param name="include_as" select="1"/>
                </xsl:call-template>
              </span>
            </td>
          </tr>
        </table>
      </fieldset>
    </div> 
  </xsl:template>
</xsl:stylesheet>
