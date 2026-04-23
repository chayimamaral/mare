<xsl:stylesheet
   xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
   xmlns:xs = "http://www.w3.org/2001/XMLSchema"
   xmlns:n="http://www.portalfiscal.inf.br/nfe"
   version="2.0"
   exclude-result-prefixes="#all">

  <xsl:output indent="yes" />
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:output method="html"/>

  <xsl:param name="tpAmb"/>
  
  <xsl:include href="Utils.xsl"/>

  <xsl:template match="/" name="nfeDpec">

    <div id="DPEC" style="width: 800px">

      <fieldset>
        <legend class="titulo-aba">DPEC</legend>
        <table>
          <tr>
            <td width="150">
              <label>Ambiente</label>
              <span>
                <xsl:if test="$tpAmb='1'">1 - Produção</xsl:if>
                <xsl:if test="$tpAmb='2'">2 - Homologação</xsl:if>
              </span>
            </td>
            <td>
              <label>Resposta da Consulta</label>
              <span>
                <xsl:value-of select="//n:cStat"/> - <xsl:value-of select="//n:xMotivo"/>
              </span>
            </td>
            <td width="100">
              <label>Versão XML</label>
              <span>
                <xsl:value-of select="//@versao"/>
              </span>
            </td>
          </tr>

        </table>
      </fieldset>

      <fieldset>
        <legend class="titulo-aba-interna">DPEC Localizado</legend>
        <table>
          <tr>
            <td>
              <label>Data Hora Registro</label>
              <span>
                <xsl:call-template name="formatDateTime">
                  <xsl:with-param name="dateTime" select="n:dhRegDPEC"/>
                  <xsl:with-param name="include_as" select="1"/>
                </xsl:call-template>
              </span>
            </td>
            <td>
              <label>Número Registro</label>
              <span>
                <xsl:value-of select="n:nRegDPEC"/>
              </span>
            </td>
          </tr>
          
        </table>
      </fieldset>

      <xsl:variable name="chNFe" select="n:chNFe"/>
      <fieldset>
        <legend class="titulo-aba-interna">Dados da NF-e</legend>
        <table>
          <tr>
            <td colspan="5">
              <label>Chave de Acesso</label>
              <span>
                <xsl:call-template name="formatNfe">
                  <xsl:with-param name="nfe" select="$chNFe"/>
                </xsl:call-template>
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <label>Número</label>
              <span>
                <xsl:value-of select="substring($chNFe,26,9)"/>
              </span>
            </td>
            <td>
              <label>Série</label>
              <span>
                <xsl:value-of select="substring($chNFe,23,3)"/>
              </span>
            </td>
            <td>
              <label>Data de Emissão</label>
              <span>
                <xsl:call-template name="formatDateTime">
                  <xsl:with-param name="dateTime" select="n:dhRegDPEC"/>
                  <xsl:with-param name="include_as" select="1"/>
                </xsl:call-template>
              </span>
            </td>
            <td>
              <label>Tipo Emissão</label>
              <span>
                4 - DPEC
              </span>
            </td>
          </tr>
        </table>
      </fieldset>

      <fieldset>
        <legend class="titulo-aba-interna">Emitente</legend>
        <table>
          <tr>
            <td width="150">
              <label>CNPJ</label>
              <span>
                <xsl:call-template name="formatCnpj">
                  <xsl:with-param name="cnpj" select="n:CNPJ"/>
                </xsl:call-template>
              </span>
            </td>
            <td width="150">
              <label>Inscrição Estadual</label>
              <span>
                <xsl:value-of select="n:IE"/>
              </span>
            </td>
            <td width="80">
              <label>UF</label>
              <span>
                <xsl:value-of select="n:cUF"/> - <xsl:call-template name="siglaUF">
                  <xsl:with-param name="uf" select="n:cUF"/>
                </xsl:call-template>
              </span>
            </td>
          </tr>
        </table>
      </fieldset>

      <fieldset>
        <legend class="titulo-aba-interna">Destinatário</legend>
        <xsl:variable name="cnpj" select="n:CNPJDest"/>
        <xsl:variable name="cpf" select="n:CPF"/>
        <xsl:variable name="ufDest" select="n:UF"/>
        <xsl:if test="$ufDest!='EX'">
          <table>
            <tr>
              <td width="150">
                <xsl:if test="$cnpj!=''">
                  <label>CNPJ</label>
                </xsl:if>
                <xsl:if test="$cpf!=''">
                  <label>CPF</label>
                </xsl:if>
                <span>
                  <xsl:call-template name="formatCnpj">
                    <xsl:with-param name="cnpj" select="$cnpj"/>
                  </xsl:call-template>
                  <xsl:call-template name="formatCpf">
                    <xsl:with-param name="cpf" select="$cpf"/>
                  </xsl:call-template>
                </span>
              </td>
              <td width="80">
                <label>UF</label>
                <span>
                  <xsl:call-template name="codigoUF">
                    <xsl:with-param name="uf" select="$ufDest"/>
                  </xsl:call-template> - <xsl:value-of select="$ufDest"/>
                </span>
              </td>
            </tr>
          </table>
        </xsl:if>
        <xsl:if test="$ufDest='EX'">
          <table>
            <tr>
              <td colspan="3">
                <label>UF</label>
                <span>
                  EX
                </span>
              </td>
            </tr>
          </table>
        </xsl:if>
      </fieldset>

      <fieldset>
        <legend class="titulo-aba-interna">Totais</legend>
        <table>
          <tr>
            <td>
              <label>Valor Total da NF-e</label>
              <span>
                <xsl:call-template name="format2Casas">
                  <xsl:with-param name="num" select="n:vNF"/>
                </xsl:call-template>
              </span>
            </td>
            <td>
              <label>Valor Total do ICMS</label>
              <span>
                <xsl:call-template name="format2Casas">
                  <xsl:with-param name="num" select="n:vICMS"/>
                </xsl:call-template>
              </span>
            </td>
            <td>
              <label>Valor Total do ICMS-ST</label>
              <span>
                <xsl:call-template name="format2Casas">
                  <xsl:with-param name="num" select="n:vST"/>
                </xsl:call-template>
              </span>
            </td>
          </tr>
        </table>
      </fieldset>

    </div>
  </xsl:template>

</xsl:stylesheet>
