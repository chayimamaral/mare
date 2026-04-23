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
  <xsl:template match="/" name="Evento_Registro_Passagem_Sistema_Externo">

   <xsl:variable name="prot" select="n:retEvento/n:infEvento/n:nProt" />
    <div id="EventoRegistroPassagemNFe{$prot}" class="GeralXslt">

        <input type="button" value="Imprimir" style="border: 1px solid black; background-color: #cccccc; margin-right: 10px;" onclick="window.print()" />
        <input type="button" value="Fechar" style="border: 1px solid black; background-color: #cccccc;" onclick="window.close()"/>

          <style type="text/css">
        .NoBorderTop
        {
        border-top:  0px !important;
        }
        .NoBorderBottom
        {
        border-bottom:  0px !important;
        }

        .NoBorderTopBottom
        {
        border-top:  0px !important;
        border-bottom:  0px !important;
        }

      </style>      
      <fieldset>
        <legend class="titulo-aba">VERIFICAÇÃO DE TRÂNSITO</legend>
        <table class="box NoBorderBottom">
          <tbody>
            <tr class="col-1">
              <td width="400">
                <label>Posto Fiscal</label>
              </td>
              <td width="150">
                <label>Sentido na Via</label>
              </td>
              <td>
                <label>Protocolo do Evento</label>
              </td>
            </tr>
            <tr class="col-1">
              <td>
                <span class="multiline">
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:cPostoUF"/> -
                  <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:xPostoUF" disable-output-escaping="yes"/>
                </span>
              </td>
              <td>
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
              <td>
                <span>
                  <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                </span>
              </td>
            </tr>
          </tbody>
        </table> 
        
        
        <xsl:variable name="cnpjEmit" select="n:evento/n:infEvento/n:CNPJ"/>
        <xsl:variable name="cpfEmit" select="n:evento/n:infEvento/n:CPF"/>
        <table class="box NoBorderTopBottom">
          <tbody>
            <tr class="col-1">
              <td>
                <label>UF</label>
              </td>
              <td>
                <xsl:choose>
                  <xsl:when test="$cnpjEmit!=''">
                    <label>CNPJ Emitente</label>
                  </xsl:when>
                  <xsl:when test="$cpfEmit!=''">
                    <label>CPF Emitente</label>
                  </xsl:when>
                  <xsl:otherwise>
                    <label>CPF / CNPJ Emitente</label>
                  </xsl:otherwise>
                </xsl:choose> 
              </td>
              <td>
                <label>Série/Número</label>
              </td>
              <td>
                <label>Tipo Emissão</label>
              </td>
            </tr>
            <tr class="col-1">
              <td>
                <span>
                  <xsl:call-template name="codigoUF">
                    <xsl:with-param name="uf" select="//n:infNFe/n:emit/n:enderEmit/n:UF"/>
                  </xsl:call-template> -
                  <xsl:value-of select="//n:infNFe/n:emit/n:enderEmit/n:UF"/>
                </span>
              </td>
              <td>
                <span> 
                    <xsl:call-template name="formatCnpj">
                      <xsl:with-param name="cnpj" select="$cnpjEmit"/>
                    </xsl:call-template>
                    <xsl:call-template name="formatCpf">
                      <xsl:with-param name="cpf" select="$cpfEmit"/>
                    </xsl:call-template>                  
                </span>
              </td>
              <td>
                <span>
                  <xsl:value-of select="//n:serie"/>/<xsl:value-of select="//n:nNF"/>
                </span>
              </td>
              <td>
                <span>
                  <xsl:for-each select="//n:tpEmis">
                    <xsl:variable name="foremis" select="text()"/>
                    <xsl:choose>
                      <xsl:when test="$foremis='1'">1 - Normal</xsl:when>
                      <xsl:when test="$foremis='2'">2 - Contingência FS</xsl:when>
                      <xsl:when test="$foremis='3'">3 - Contingência SCAN</xsl:when>
                      <xsl:when test="$foremis='4'">4 - Contingência DPEC</xsl:when>
                      <xsl:when test="$foremis='5'">5 - Contingência FS-DA</xsl:when>
                      <xsl:when test="$foremis='6'">6 - Contingência SVC-AN</xsl:when>
                      <xsl:when test="$foremis='7'">7 - Contingência SVC-RS</xsl:when>
                      <xsl:when test="$foremis='9'">9 - Contingência NFC-e off-line</xsl:when>
                      <xsl:otherwise>
                        <xsl:value-of select="$foremis"/>
                      </xsl:otherwise>
                    </xsl:choose>
                  </xsl:for-each>
                </span>
              </td>
            </tr>
          </tbody>
        </table> 
        
        
        <table class="box NoBorderTopBottom">
          <tbody>
            <tr class="col-1">
              <td>
                <label>
                  UF
                </label>
              </td>
                  <xsl:variable name="destinatario" select="//n:infNFe/n:dest"/>
                  <xsl:if test="$destinatario != ''">
              <td>
                  <xsl:variable name="cnpj" select="//n:infNFe/n:dest/n:CNPJ"/>
                  <xsl:variable name="cpf" select="//n:infNFe/n:dest/n:CPF"/>
                  <xsl:variable name="idestrang" select="//n:infNFe/n:dest/n:idEstrangeiro"/>

                  <xsl:choose>
                    <xsl:when test="$cnpj!=''">
                      <label>CNPJ Destinatário</label>
                    </xsl:when>
                    <xsl:when test="$cpf!=''" >
                      <label>CPF Destinatário</label>
                    </xsl:when>
                    <xsl:when test="$idestrang!=''" >
                      <label>Id. Estrangeiro Destinatário</label>
                    </xsl:when>
                    <xsl:otherwise>
                      <label>CNPJ/CPF/Id. Estrangeiro Destinatário</label>
                    </xsl:otherwise>
                  </xsl:choose>
                </td>
              </xsl:if>

              <td>
                <label>Tipo Operação</label>
              </td>
              <td>
                <label>Dt. Emissão</label>
              </td>
            </tr>
            <tr class="col-1">
              <td>
                <span>
                  <xsl:call-template name="codigoUF">
                    <xsl:with-param name="uf" select="//n:infNFe/n:dest/n:enderDest/n:UF"/> 
                  </xsl:call-template> - 
                  <xsl:value-of select="//n:infNFe/n:dest/n:enderDest/n:UF"/>
                </span>
              </td>
              <td>
                <span>
                  <xsl:variable name="cnpj" select="//n:infNFe/n:dest/n:CNPJ"/>
                  <xsl:variable name="cpf" select="//n:infNFe/n:dest/n:CPF"/>
                  <xsl:variable name="idestrang" select="//n:infNFe/n:dest/n:idEstrangeiro"/>
                  <xsl:call-template name="formatCnpj">
                    <xsl:with-param name="cnpj" select="$cnpj"/>
                  </xsl:call-template>
                  <xsl:call-template name="formatCpf">
                    <xsl:with-param name="cpf" select="$cpf"/>
                  </xsl:call-template>
                  <xsl:value-of select="$idestrang"/>
                </span>
              </td>
              <td>
                <span>
                  <xsl:for-each select="//n:ide/n:tpNF">
                    <xsl:variable name="tipdocfis" select="text()"/>
                    <xsl:choose>
                      <xsl:when test="$tipdocfis = 0">
                        0 - Entrada
                      </xsl:when>
                      <xsl:when test="$tipdocfis = 1">
                        1 - Saída
                      </xsl:when>
                      <xsl:otherwise>
                        <xsl:value-of select="$tipdocfis"/>
                      </xsl:otherwise>
                    </xsl:choose>
                  </xsl:for-each>
                </span>
              </td>
              <td>
                <span>
                  <xsl:variable name="dEmi" select="//n:infNFe/n:ide/n:dEmi"/>
                  <xsl:variable name="dhEmi" select="//n:infNFe/n:ide/n:dhEmi"/>
                  <xsl:if test="$dEmi != ''">
                    <xsl:call-template name="formatDate">
                      <xsl:with-param name="date" select="$dEmi"/>
                    </xsl:call-template>
                  </xsl:if>
                  <xsl:if test="$dhEmi != ''">
                    <xsl:call-template name="formatDateTimeFuso">
                      <xsl:with-param name="dateTime" select="$dhEmi"/>
                    </xsl:call-template>
                  </xsl:if>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        
        <table class="box NoBorderTopBottom">
          <tbody>
            <tr class="col-1">
              <td>
                <label>Valor Total</label>
              </td>
              <td>
                <label>Protocolo NF-e</label>
              </td>
              <td>
                <label>Data-Hora Passagem</label>
              </td>
            </tr>
            <tr class="col-1">
              <td>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="//n:total/n:ICMSTot/n:vNF"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <span>
                  <xsl:value-of select="//n:protNFe/n:infProt/n:nProt"/>
                </span>
              </td>
              <td>
                <span>
                  <xsl:call-template name="formatDateTime">
                    <xsl:with-param name="dateTime" select="n:evento/n:infEvento/n:detEvento/n:dhPas"/>
                    <xsl:with-param name="include_as" select="1"/>
                  </xsl:call-template>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        
        <table class="box NoBorderTop">
          <tbody>
            <tr class="col-1">
              <td>
                <label>Modal</label>
              </td>


              <xsl:if test="n:evento/n:infEvento/n:detEvento/n:modalRodov">
                <td>
                  <label>Veículo</label>
                </td>
                <td>
                  <label>Carreta</label>
                </td>
                <td>
                  <label>Carreta 2</label>
                </td> 
              </xsl:if>
              <xsl:if test="n:evento/n:infEvento/n:detEvento/n:modalOutro">
                <td>
                  <label>Identificação do Transporte</label>
                </td> 
              </xsl:if> 
            </tr>
            <tr class="col-1">
              <td>
                <span>
                  <xsl:if test="n:evento/n:infEvento/n:detEvento/n:modalRodov">
                    Rodoviário
                  </xsl:if>
                  <xsl:if test="n:evento/n:infEvento/n:detEvento/n:modalOutro">
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
                  </xsl:if>
                </span>
              </td>
              <xsl:if test="n:evento/n:infEvento/n:detEvento/n:modalRodov">
                <td>
                  <span>
                    <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:placaVeic"/> - <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:UFVeic"/>
                  </span>
                </td>
                <td>
                  <span>
                    <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:placaCarreta"/> - <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:UFCarreta"/>
                  </span>
                </td>
                <td>
                  <span>
                    <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:placaCarreta2"/> - <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalRodov/n:UFCarreta2"/>
                  </span>
                </td>
              </xsl:if>
              <xsl:if test="n:evento/n:infEvento/n:detEvento/n:modalOutro">
                <td>
                  <span>
                    <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:modalOutro/n:xIdent"/>
                  </span>
                </td> 
              </xsl:if>
            </tr>
          </tbody>
        </table>  
        
     
      </fieldset>
    </div>
  </xsl:template>
</xsl:stylesheet>
