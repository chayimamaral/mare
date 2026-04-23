<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  
  <xsl:include href="Utils.xsl" />
  <xsl:include href="css.xsl" />

  <xsl:decimal-format decimal-separator="," grouping-separator="."/>

  <xsl:output method="html" indent="no"/>

  <xsl:param name="fonte"/>
  
  <xsl:template match="/">

    <div class="divXSLT">

      <xsl:call-template name="css" />
      
      <table class="titulo">
        <thead>
          <tr>
            <td colspan="4">
              <xsl:choose>
                <xsl:when test="$fonte = 'local'">
                  Pesquisa de Passagens na UF
                </xsl:when>
                <xsl:otherwise>
                  Pesquisa de Passagens no AN
                </xsl:otherwise>
              </xsl:choose>
              <xsl:choose>
                <xsl:when test="count(//n:Passagens/n:Passagem) > 1">
                  - (<xsl:value-of select="count(//n:Passagens/n:Passagem)"/> Registros de Passagem)
                </xsl:when>
                <xsl:when test="count(//n:Passagens/n:Passagem) = 1">
                  - (<xsl:value-of select="count(//n:Passagens/n:Passagem)"/> Registro de Passagem)
                </xsl:when>
              </xsl:choose>
            </td>
          </tr>
        </thead>
      </table>

      <table class="dados">
        <thead>
          <td>
            <label>Chave de Acesso</label>
          </td>
        </thead>
        <tbody>
          <td>
            <xsl:call-template name="formatNfe">
              <xsl:with-param name="nfe" select="//n:chNFe"/>
            </xsl:call-template>
          </td>
        </tbody>
      </table>
      <p/>
      
      <xsl:choose>
        <xsl:when test="//n:cStat = '109' or //n:cStat = '9999' ">

          <table class="dados">
            <thead>
              <td>
                <label>Situação</label>
              </td>
            </thead>
            <tbody>
              <td>
                <xsl:value-of select="//n:cStat"/> - <xsl:value-of select="//n:xMotivo"/>
              </td>
            </tbody>
          </table>
          <p/>
          
        </xsl:when>
        <xsl:otherwise>
          <xsl:choose>
            <xsl:when test="count(//n:Passagens/n:Passagem) = 0">
              <table class="dados">
                <thead>
                  <td>
                    <label>Situação</label>
                  </td>
                </thead>
                <tbody>
                  <td>
                    <xsl:value-of select="//n:cStat"/> - <xsl:value-of select="//n:xMotivo"/>
                  </td>
                </tbody>
              </table>
              <p/>
            </xsl:when>
          </xsl:choose>
        </xsl:otherwise>
      </xsl:choose>

      <xsl:for-each select="//n:Passagens/n:Passagem">
        <xsl:sort select="n:dPass" order="descending" />

        <table class="titulo">
          <thead>
            <tr>
              <td colspan="4">
                <xsl:choose>
                  <xsl:when test="@Id != ''">
                    Registro de Passagem número: <xsl:value-of select="@Id" />
                  </xsl:when>
                  <xsl:otherwise>
                    Registro de Passagem número: <xsl:value-of select="count(//n:Passagens/n:Passagem)+1 - position()"/>
                  </xsl:otherwise>
                </xsl:choose>
              </td>
            </tr>
          </thead>
        </table>
        
        <table class="dados">
          <thead>
            <td width="100">
              <label>UF Passagem</label>
            </td>
            <td width="400">
              <label>Posto Fiscal</label>
            </td>
            <td>
              <label>Data Passagem</label>
            </td>
          </thead>
          <tbody>
            <td>
              <xsl:value-of select="n:cUF"/> - <xsl:call-template name="siglaUF">
                <xsl:with-param name="uf" select="n:cUF"/>
              </xsl:call-template>
            </td>
            <td>
              <xsl:value-of select="n:cUnidFiscal"/> - <xsl:value-of select="n:xUnidFiscal"/>
            </td>
            <td>
              <xsl:call-template name="formatDateTime">
                <xsl:with-param name="dateTime" select="n:dPass"/>
                <xsl:with-param name="include_as" select="1"/>
              </xsl:call-template>
            </td>
          </tbody>
        </table>

        <xsl:variable name="rodo" select="count(n:mdTransp/n:rodoviario)"/>
        <xsl:variable name="outro" select="n:mdTransp/n:outro/n:cMod"/>
        <table class="dados">
          <thead>
            <td width="120">
              <label>Modal</label>
            </td>

            <xsl:choose>
              <xsl:when test="$rodo = 1">
                <td>
                  <label>Veículo</label>
                </td>
                <td>
                  <label>Carreta</label>
                </td>
                <td>
                  <label>Carreta 2</label>
                </td>
              </xsl:when>
              <xsl:when test="$outro != ''">
                <td>
                  <label>Identificação do Transporte</label>
                </td>
              </xsl:when>
            </xsl:choose>
          </thead>
          <tbody>
            <td>
              <xsl:choose>
                <xsl:when test="$rodo = 1">Rodoviário</xsl:when>
                <xsl:when test="$outro = 'F'">Ferroviário</xsl:when>
                <xsl:when test="$outro = 'D'">Dutoviário</xsl:when>
                <xsl:when test="$outro = 'AE'">Aéreo</xsl:when>
                <xsl:when test="$outro = 'AQ'">Aquaviário</xsl:when>
                <xsl:when test="$outro = 'O'">Outro</xsl:when>
              </xsl:choose>
            </td>
            <xsl:choose>
              <xsl:when test="$rodo = 1">
                <td>
                  <xsl:value-of select="n:mdTransp/n:rodoviario/n:pVeic"/> / <xsl:call-template name="siglaUF">
                    <xsl:with-param name="uf" select="n:mdTransp/n:rodoviario/n:cUFVeic"/>
                  </xsl:call-template>
                </td>
                <td>
                  <xsl:choose>
                    <xsl:when test="n:mdTransp/n:rodoviario/n:pCarreta != ''">
                      <xsl:value-of select="n:mdTransp/n:rodoviario/n:pCarreta"/> / <xsl:call-template name="siglaUF">
                        <xsl:with-param name="uf" select="n:mdTransp/n:rodoviario/n:cUFCarreta"/>
                      </xsl:call-template>
                    </xsl:when>
                    <xsl:otherwise>-</xsl:otherwise>
                  </xsl:choose>
                </td>
                <td>
                  <xsl:choose>
                    <xsl:when test="n:mdTransp/n:rodoviario/n:pCarreta2 != ''">
                      <xsl:value-of select="n:mdTransp/n:rodoviario/n:pCarreta2"/> / <xsl:call-template name="siglaUF">
                        <xsl:with-param name="uf" select="n:mdTransp/n:rodoviario/n:cUFCarreta2"/>
                      </xsl:call-template>
                    </xsl:when>
                    <xsl:otherwise>-</xsl:otherwise>
                  </xsl:choose>
                </td>
              </xsl:when>
              <xsl:when test="$outro != ''">
                <td>
                  <xsl:if test="n:mdTransp/n:outro/n:xIdent=''">-</xsl:if>
                  <xsl:value-of select="n:mdTransp/n:outro/n:xIdent"/>
                </td>
              </xsl:when>
            </xsl:choose>
          </tbody>
        </table>

        <table class="dados">
          <thead>
            <td width="150">
              <label>Tipo Transmissão</label>
            </td>
            <td width="120">
              <label>Sentido Passagem</label>
            </td>
            <td>
              <label>Operador</label>
            </td>
          </thead>
          <tbody>
            <td>
              <xsl:variable name="tTrasm" select="n:tTrasm"/>
              <xsl:choose>
                <xsl:when test="$tTrasm = 'N'">Normal</xsl:when>
                <xsl:when test="$tTrasm = 'A'">Atrasada</xsl:when>
                <xsl:otherwise>-</xsl:otherwise>
              </xsl:choose>
            </td>
            <td>
              <xsl:variable name="tSentido" select="n:tSentido"/>
              <xsl:choose>
                <xsl:when test="$tSentido = 'E'">Entrada na UF</xsl:when>
                <xsl:when test="$tSentido = 'S'">Saída da UF</xsl:when>
                <xsl:when test="$tSentido = 'I'">Indeterminado</xsl:when>
                <xsl:otherwise>-</xsl:otherwise>
              </xsl:choose>
            </td>
            <td>
              <xsl:call-template name="formatCpf">
                <xsl:with-param name="cpf" select="n:cpfFunc"/>
              </xsl:call-template> - <xsl:value-of select="n:xFunc"/>
            </td>
          </tbody>
        </table>

        <table class="dados">
          <thead>
            <td width="400">
              <label>Observação Passagem</label>
            </td>
            <td width="120">
              <label>Retorno</label>
            </td>
            <td>
              <label>UF Destinatário</label>
            </td>
          </thead>
          <tbody>
            <td>
              <xsl:if test="n:xObs=''">-</xsl:if>
              <xsl:value-of select="n:xObs"/>
            </td>
            <td>
              <xsl:variable name="tIndRet" select="n:tIndRet"/>
              <xsl:choose>
                <xsl:when test="$tIndRet = 'D'">Devolução</xsl:when>
                <xsl:when test="$tIndRet = 'R'">Retorno</xsl:when>
                <xsl:otherwise>-</xsl:otherwise>
              </xsl:choose>
            </td>
            <td>
              <xsl:variable name="ufDest" select="n:cUFDest"/>
              <xsl:choose>
                <xsl:when test="$ufDest = '0'">EX</xsl:when>
                <xsl:otherwise>
                  <xsl:value-of select="$ufDest"/> - <xsl:call-template name="siglaUF">
                    <xsl:with-param name="uf" select="$ufDest"/>
                  </xsl:call-template>
                </xsl:otherwise>
              </xsl:choose>
            </td>
          </tbody>
        </table>

        <p />
        
      </xsl:for-each>

    </div>
      
  </xsl:template>
  
</xsl:stylesheet>
