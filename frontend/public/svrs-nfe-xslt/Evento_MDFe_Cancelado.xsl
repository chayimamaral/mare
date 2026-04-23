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
  <xsl:template match="/" name="Evento_MDFe_Cancelado">

    <xsl:variable name="prot" select="n:retEvento/n:infEvento/n:nProt" />
    <div id="EventoMDFeCancelado{$prot}" class="GeralXslt">
      <fieldset>
        <legend class="titulo-aba">MDF-e Cancelado</legend>
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
          <tr>
            <td class="col-2" colspan="2">
              <label>Código Autor do Evento</label>
              <span>
                <xsl:variable name="cCodAutor" select="n:evento/n:infEvento/n:detEvento/n:cOrgaoAutor"/>
                <xsl:choose>
                  <xsl:when test="$cCodAutor = 91">
                    91 = AN (Serpro)
                  </xsl:when>
                  <xsl:when test="$cCodAutor = 92">
                    92 = AN do MDF-e (Procergs)
                  </xsl:when>
                  <xsl:otherwise>
                    <xsl:value-of select="$cCodAutor"/>
                  </xsl:otherwise>
                </xsl:choose>
              </span>
            </td>
            <td class="col-2" colspan="2">
              <label>Tipo Autor</label>
              <span>
                <xsl:variable name="tpAutor" select="n:evento/n:infEvento/n:detEvento/n:tpAutor"/>
                <xsl:choose>
                  <xsl:when test="$tpAutor = 1">
                    1 = Empresa Emitente
                  </xsl:when>
                  <xsl:when test="$tpAutor = 2">
                    2 = Empresa Destinatária
                  </xsl:when>
                  <xsl:when test="$tpAutor = 3">
                    3 = Empresa
                  </xsl:when>
                  <xsl:when test="$tpAutor = 5">
                    5 = Fisco
                  </xsl:when>
                  <xsl:when test="$tpAutor = 6">
                    6 = RFB
                  </xsl:when>
                  <xsl:when test="$tpAutor = 9">
                    9 = Outros
                  </xsl:when>
                  <xsl:otherwise>
                    <xsl:value-of select="$tpAutor"/>
                  </xsl:otherwise>
                </xsl:choose>
              </span>
            </td>
          </tr>
          <tr>
            <td class="col-2" colspan="2">
              <label>Versão Aplicativo Autor Evento</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:verAplic"/>
              </span>
            </td>
          </tr>
        </table>
        <br/>
        <table class="box">
          <tr class="col-2">
            <td colspan="2">
              <label>Chave de Acesso MDF-e</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:MDFe/n:chMDFe"/>
              </span>
            </td>
            <td colspan="2">
              <label>Protocolo de Cancelamento MDF-e</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:MDFe/n:nProtCanc"/>
              </span>
            </td>
          </tr>
        </table>
      </fieldset>
      <fieldset>
        <legend class="titulo-aba">Autorização pela SEFAZ</legend>
        <table class="box">
          <tr>
            <td width="50%" colspan="3">
              <label>Mensagem de Autorização</label>
              <span>
                <xsl:value-of select="n:retEvento/n:infEvento/n:cStat"/> -
                <xsl:value-of select="n:retEvento/n:infEvento/n:xMotivo"/>
              </span>
            </td>
            <td width="20%" colspan="3">
              <label>Protocolo</label>
              <span>
                <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
              </span>
            </td>
            <td width="30%" colspan="3">
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
    </div>
  </xsl:template>
</xsl:stylesheet>
