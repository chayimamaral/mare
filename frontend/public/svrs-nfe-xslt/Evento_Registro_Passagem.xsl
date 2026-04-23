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
  <xsl:template match="/" name="Evento_Registro_Passagem">
    <xsl:variable name="prot" select="n:retEvento/n:infEvento/n:nProt" />
    
    <div id="EventoRegistroPassagemNFe{$prot}" class="GeralXslt">
        <fieldset>
          <legend class="titulo-aba">Registro Passagem NF-e</legend>
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
                <label>Órgão Autor Registro de Passagem</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:cOrgaoAutor"/> - 
                  <xsl:call-template name="nomeUF">
                    <xsl:with-param name="uf" select="n:evento/n:infEvento/n:detEvento/n:cOrgaoAutor"/>
                  </xsl:call-template>  
                </span>
              </td>
              <td colspan="2">
                <label>Posto Fiscal</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:cPostoUF"/> -
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:xPostoUF" disable-output-escaping="yes"/>
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td colspan="2">
                <label>Latitude do Local</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:latGPS"/> 
                </span>
              </td>
              <td colspan="2">
                <label>Longitude do Local</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:longGPS"/>                  
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td colspan="2">
                <label>Dados do Operador</label>
                <span>
                  <xsl:call-template name="formatCpf">
                    <xsl:with-param name="cpf" select="n:evento/n:infEvento/n:detEvento/n:CPFOper"/>
                  </xsl:call-template> - <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:xNomeOper"/>
                </span>
              </td>
              <td colspan="2">
                <label>Data Passagem</label>
                <span>
                  <xsl:call-template name="formatDateTimeFuso">
                    <xsl:with-param name="dateTime" select="n:evento/n:infEvento/n:detEvento/n:dhPas"/>
                    <xsl:with-param name="include_as" select="1"/>
                  </xsl:call-template>
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td colspan="2">
                <label>Indicador de Transmissão Off-line</label>
                <span>
                  <xsl:variable name="indTransOffline" select="n:evento/n:infEvento/n:detEvento/n:indOffline"/>
                  <xsl:choose>
                    <xsl:when test="$indTransOffline = '1'">
                      1 - Transmissão do Evento off-line
                    </xsl:when>
                    <xsl:when test="$indTransOffline = '0'">
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
                <label>Indicador de Retorno</label>
                <span>
                  <xsl:variable name="indicadorRetorno" select="n:evento/n:infEvento/n:detEvento/n:indRet"/>
                  <xsl:choose>
                    <xsl:when test="$indicadorRetorno = 'R'">
                      R - Retorno
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$indicadorRetorno"/>
                    </xsl:otherwise>
                  </xsl:choose>
                </span>
              </td>
              <td colspan="2">
                <label>UF Destino</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:UFDest"/>
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td colspan="2">
                <label>Chave MDF-e</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:chMDFe"/>
                </span>
              </td>
              <td colspan="2">
                <label>Chave CT-e</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:chCTe"/>
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


      <xsl:if test="n:evento/n:infEvento/n:detEvento/n:modalRodov">
      
        <fieldset >
          <legend>Modal Rodoviário</legend>
          <table class="box">
            <tr class="col-2">
              <td colspan="2">
                <label>Placa do Veículo</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:placaVeic"/>
                </span>
              </td>
              <td colspan="2">
                <label>UF Veículo</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:UFVeic"/>
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td colspan="2">
                <label>Placa da Carreta</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:placaCarreta"/>
                </span>
              </td>
              <td colspan="2">
                <label>UF Carreta</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:UFCarreta"/>
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td colspan="2">
                <label>Placa da Carreta 2</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:placaCarreta2"/>
                </span>
              </td>
              <td colspan="2">
                <label>UF Carreta 2</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:UFCarreta2"/>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
      </xsl:if>

      <xsl:if test="n:evento/n:infEvento/n:detEvento/n:modalOutro">
        <fieldset>
          <legend>Modal Outros (não Rodoviário)</legend>
          <table class="box">
            <tr class="col-2">
              <td colspan="2">
                <label>Outra Modalidade de Transporte</label>
                <span>
                  <xsl:variable name="modOutTrans" select="n:evento/n:infEvento/n:detEvento/n:modalOutro/n:tpModal"/>
                  <xsl:choose>
                    <xsl:when test="$modOutTrans = 'F'">Ferroviário</xsl:when>
                    <xsl:when test="$modOutTrans = 'D'">Dutoviário</xsl:when>
                    <xsl:when test="$modOutTrans = 'AE'">Aeroviário</xsl:when>
                    <xsl:when test="$modOutTrans = 'AQ'">Aquaviário</xsl:when>
                    <xsl:when test="$modOutTrans = 'O'">Outros</xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalOutro/n:tpModal"/>
                    </xsl:otherwise>
                  </xsl:choose>
                </span>
              </td>
              <td colspan="2">
                <label>Identificação Meio de Transporte</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalOutro/n:xIdent"/>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
      </xsl:if>

      <xsl:if test="n:evento/n:infEvento/n:detEvento/n:ctg">
        <fieldset>
          <legend>Emissão em Contingência</legend>
          <table class="box">
            <tr class="col-2">
              <td colspan="2">
                <label>Número do Formulário de Segurança</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:ctg/n:nFormSeg"/>
                </span>
              </td>
              <td colspan="2">
                <label>UF de Destino</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:ctg/n:UFDest"/>
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td colspan="2">
                <xsl:variable name="cnpjDest" select="n:evento/n:infEvento/n:detEvento/n:ctg/n:CNPJDest"/>
                <xsl:variable name="cpfDest" select="n:evento/n:infEvento/n:detEvento/n:ctg/n:CPFDest"/>

                <xsl:choose>
                  <xsl:when test="$cnpjDest!=''">
                    <label>CNPJ Destinatário</label>
                  </xsl:when>
                  <xsl:when test="$cpfDest!=''">
                    <label>CPF Destinatário</label>
                  </xsl:when>
                  <xsl:otherwise>
                    <label>CPF / CNPJ Destinatário</label>
                  </xsl:otherwise>
                </xsl:choose>  
                <span>
                  <xsl:call-template name="formatCnpj">
                    <xsl:with-param name="cnpj" select="$cnpjDest"/>
                  </xsl:call-template>
                  <xsl:call-template name="formatCpf">
                    <xsl:with-param name="cpf" select="$cpfDest"/>
                  </xsl:call-template>
                </span>
              </td>
              <td colspan="2">
                <label>Valor Total da NF-e</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="n:evento/n:infEvento/n:detEvento/n:ctg/n:vTotalNFe"/>
                  </xsl:call-template>
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td colspan="2">
                <label>Indicador de Destaque de ICMS próprio</label>
                <span>
                  <xsl:variable name="indICMSPrp" select="n:evento/n:infEvento/n:detEvento/n:ctg/n:indICMS"/>
                  <xsl:choose>
                    <xsl:when test="$indICMSPrp = '1'">
                      1 - Sim
                    </xsl:when>
                    <xsl:when test="$indICMSPrp = '2'">
                      2 - Não
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$indICMSPrp"/>
                    </xsl:otherwise>
                  </xsl:choose>
                </span>
              </td>
              <td colspan="2">
                <label>Indicador de Destaque de ICMS-ST</label>
                <span>
                  <xsl:variable name="indICMSST" select="n:evento/n:infEvento/n:detEvento/n:ctg/n:indICMSST"/>
                  <xsl:choose>
                    <xsl:when test="$indICMSST = '1'">
                      1 - Sim
                    </xsl:when>
                    <xsl:when test="$indICMSST = '2'">
                      2 - Não
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$indICMSST"/>
                    </xsl:otherwise>
                  </xsl:choose>
                </span>
              </td>
            </tr>
            <tr class="col-2">
              <td colspan="2">
                <label>Dia Emissão NF-e</label>
                <span>
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:ctg/n:diaEmi"/>
                </span>
              </td>
              <td colspan="2">
                <label>Tipo de Emissão da NF-e</label>
                <span>
                  <xsl:variable name="tpEmisNFe" select="n:evento/n:infEvento/n:detEvento/n:ctg/n:tpEmis"/>
                  <xsl:choose>
                    <xsl:when test="$tpEmisNFe = 2">
                      2 - FS
                    </xsl:when>
                    <xsl:when test="$tpEmisNFe = 5">
                      5 - FSDA
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$tpEmisNFe"/>
                    </xsl:otherwise>
                  </xsl:choose> 
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
      </xsl:if> 
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
