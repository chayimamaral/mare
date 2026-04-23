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
  <xsl:template match="/" name="Evento_Registro_Passagem_BRId">

    <xsl:variable name="prot" select="n:retEvento/n:infEvento/n:nProt" />

    <div id="EventoRegistroPassagemNFeBRId{$prot}" class="GeralXslt">
      <fieldset>
        <legend class="titulo-aba">Registro Passagem NF-e BRId</legend>
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
        </table>
        <br/>
        <table class="box">
          <tr class="col-2">
            <td colspan="2">
              <label>Tipo de Evento</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:tpEvento"/> - Registro Passagem NF-e
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
              <label>Código UF Localização Antena</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:cOrgaoAutor"/>
              </span>
            </td>
            <td colspan="2">
              <label>Identificação da Antena</label>
              <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:cIdAntena"/> - <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:xIdAntena"/>
              </span>
            </td>            
          </tr>
          <tr class="col-2">
            <td colspan="2">
              <label>Latitude </label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:latGPS"/> 
              </span>
            </td>
            <td colspan="2">
              <label>Longitude</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:longGPS"/>
              </span>
            </td>            
          </tr>
          <tr class="col-2">
            <td colspan="2">
              <label>Data Passagem</label>
              <span>
                <xsl:call-template name="formatDateTimeFuso">
                  <xsl:with-param name="dateTime" select="n:evento/n:infEvento/n:detEvento/n:dhPas"/>
                  <xsl:with-param name="include_as" select="1"/>
                </xsl:call-template>
              </span>
            </td>
            <td colspan="2">
              <label>Placa do Veículo</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:placaVeic"/>
              </span>
            </td>
          </tr>
          <tr class="col-2">
            <td colspan="2">
              <label>Indicador de Transmissão Off-line</label>
              <span>
                <xsl:variable name="indTransOffline" select="n:evento/n:infEvento/n:detEvento/n:indOffline"/>
                <xsl:choose>
                  <xsl:when test="$indTransOffline = 1">
                    1 - Transmissão do Evento off-line
                  </xsl:when>
                  <xsl:when test="$indTransOffline = 0">
                    0 - Transmissão do Evento on-line
                  </xsl:when>
                  <xsl:otherwise>
                    <xsl:value-of select="$indTransOffline"/>
                  </xsl:otherwise>
                </xsl:choose>
              </span>
            </td>
            <td colspan="2">
              <label>Sentido na Via</label>
              <span>
                <xsl:variable name="sentidoVia" select="n:evento/n:infEvento/n:detEvento/n:sentidoVia"/>
                <xsl:choose>
                  <xsl:when test="$sentidoVia = 'E'">
                    E - Entrada na UF
                  </xsl:when>
                  <xsl:when test="$sentidoVia = 'S'">
                    S - Saída da UF
                  </xsl:when>
                  <xsl:when test="$sentidoVia = 'I'">
                    I - Indeterminado
                  </xsl:when>
                  <xsl:otherwise>
                    <xsl:value-of select="$sentidoVia"/>
                  </xsl:otherwise>
                </xsl:choose>
              </span>
            </td>
          </tr> 
          <tr class="col-2">
            <td colspan="2">
              <label>Chave de Acesso do MDF-e</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:chMDFe"/>
              </span>
            </td>
            <td colspan="2">
              <label>Chave de Acesso do CT-e</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:chCTe"/>
              </span>
            </td>
          </tr>
          <tr class="col-2">
            <td colspan="2">
              <label>NSU Registro Base BackOffice</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:NSU"/>
              </span>
            </td>
            <td colspan="2">
              <label>NSU Registro Antena  Base BackOffice</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:NSUAntena"/>
              </span>
            </td>
          </tr>  
          <tr class="col-4">
            <td colspan="4">
              <label>Observação</label>
              <span>
                <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:xObs" disable-output-escaping="yes"/>
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
    </div>
  </xsl:template>
</xsl:stylesheet>
