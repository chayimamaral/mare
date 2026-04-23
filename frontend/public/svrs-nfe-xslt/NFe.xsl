<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
xmlns:fo="http://www.w3.org/1999/XSL/Format"
xmlns:n="http://www.portalfiscal.inf.br/nfe"
xmlns:s="http://www.w3.org/2000/09/xmldsig#"
version="2.0"
exclude-result-prefixes="fo n s">  
  <xsl:import href="Evento_CCe.xsl"/>
  <xsl:import href="Evento_Cancelamento.xsl"/>
  <xsl:import href="Evento_EPEC.xsl"/>
  <xsl:import href="NFe_Cancelamento.xsl"/>
  <xsl:import href="Evento_Confirmacao.xsl"/>
  <xsl:import href="Evento_Ciencia.xsl"/>
  <xsl:import href="Evento_Desconhecimento.xsl"/>
  <xsl:import href="Evento_Nao_Realizado.xsl"/>
  <xsl:import href="Evento_CTe_Autorizado.xsl"/>
  <xsl:import href="Evento_CTe_Cancelado.xsl"/>
  <xsl:import href="Evento_Vistoria_SUFRAMA.xsl"/>
  <xsl:import href="Evento_Internalizacao_SUFRAMA.xsl"/>
  <xsl:import href="Evento_Registro_Passagem.xsl"/>
  <xsl:import href="Evento_Registro_Passagem_Sistema_Externo.xsl"/>
  <xsl:import href="Evento_MDFe_Autorizado.xsl"/>
  <xsl:import href="Evento_MDFe_Cancelado.xsl"/>  
  <xsl:import href="Evento_Registro_Passagem_BRId.xsl"/>
  <xsl:import href="Evento_Cancelamento_Registro_Passagem.xsl"/>   
  <xsl:include href="Utils.xsl"/>
  <xsl:include href="_Estilos_Geral.xsl"/>
  <xsl:include href="_Scripts_Geral.xsl"/>  
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:output method="html"/>
  <xsl:template match="/" name="NFe">

    <xsl:param name="ambiente" select="'publico'"/>
    <div id="NFe" class="GeralXslt">
      <fieldset >
        <legend class="titulo-aba">Dados da NF-e</legend>
        <table class="box">
          <tr>
            <td>
              <label>Modelo</label>
              <span>
                <xsl:value-of select="//n:mod"/>
              </span>
            </td>
            <td>
              <label>Série</label>
              <span>
                <xsl:value-of select="//n:serie"/>
              </span>
            </td>
            <td>
              <label>Número</label>
              <span>
                <xsl:value-of select="//n:nNF"/>
              </span>
            </td>
            <td>
              <label>Data de Emissão</label>
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
            <xsl:variable name="dhSaiEnt" select="//n:dhSaiEnt"/>
            <xsl:choose>
              <xsl:when test="$dhSaiEnt != ''">
                  <td>
                    <label>
                      Data/Hora de Saída ou da Entrada
                    </label>
                    <span>
                      <xsl:call-template name="formatDateTimeFuso">
                        <xsl:with-param name="dateTime" select="$dhSaiEnt"/>
                      </xsl:call-template>
                    </span>
                  </td>
              </xsl:when>
              <xsl:otherwise>
                  <xsl:variable name="dSaiEnt" select="//n:dSaiEnt"/>
                  <xsl:variable name="hSaiEnt" select="//n:hSaiEnt"/>
                  <td>
                    <label>
                      Data<xsl:if test="$hSaiEnt != ''">/Hora </xsl:if> Saída/Entrada
                    </label>
                    <span>
                      <xsl:call-template name="formatDate">
                        <xsl:with-param name="date" select="$dSaiEnt"/>
                      </xsl:call-template>
                      <xsl:if test="$hSaiEnt != ''">
                        às <xsl:value-of select="//n:hSaiEnt"/>
                      </xsl:if>
                    </span>
                  </td>
              </xsl:otherwise>
            </xsl:choose> 
            <td>
              <label>Valor&#160;Total&#160;da&#160;Nota&#160;Fiscal&#160;&#160;</label>
              <span>
                <xsl:call-template name="format2Casas">
                  <xsl:with-param name="num" select="//n:infNFe/n:total/n:ICMSTot/n:vNF"/>
                </xsl:call-template>
              </span>
            </td>
          </tr>
        </table>
      </fieldset>
      <fieldset>
        <legend>Emitente</legend>
        <table class="box">
          <tr>
            <td class="col-5"> 
              <xsl:variable name="cnpjEmit" select="//n:infNFe/n:emit/n:CNPJ"/>
              <xsl:variable name="cpfDest" select="//n:infNFe/n:emit/n:CPF"/>
              <xsl:if test="$cnpjEmit!=''">
                <label>CNPJ</label>
              </xsl:if>
              <xsl:if test="$cpfDest!=''">
                <label>CPF</label>
              </xsl:if>
              <span>
                <xsl:call-template name="formatCnpj">
                  <xsl:with-param name="cnpj" select="$cnpjEmit"/>
                </xsl:call-template>
                <xsl:call-template name="formatCpf">
                  <xsl:with-param name="cpf" select="$cpfDest"/>
                </xsl:call-template>
              </span>
            </td>
            <td class="col-2">
              <label>Nome / Razão Social</label>
              <span>
                <xsl:value-of select="//n:infNFe/n:emit/n:xNome"/>
              </span>
            </td>
            <td class="col-5">
              <label>Inscrição Estadual</label>
              <span>
                <xsl:value-of select="//n:infNFe/n:emit/n:IE"/>
              </span>
            </td>
            <td class="col-10">
              <label>UF</label>
              <span>
                <xsl:value-of select="//n:infNFe/n:emit/n:enderEmit/n:UF"/>
              </span>
            </td>
          </tr>
        </table>
      </fieldset>
 
      <xsl:variable name="destinatario" select="//n:infNFe/n:dest"/>
      <xsl:if test="$destinatario != ''">      
        <fieldset>
          <legend>Destinatário</legend>
          <table class="box">
            <tr>
              <td class="col-5">
                <xsl:variable name="cnpj" select="//n:infNFe/n:dest/n:CNPJ"/>
                <xsl:variable name="cpf" select="//n:infNFe/n:dest/n:CPF"/>
                <xsl:variable name="idestrang" select="//n:infNFe/n:dest/n:idEstrangeiro"/>

                <xsl:choose>
                  <xsl:when test="$cnpj!=''">
                    <label>CNPJ</label>
                  </xsl:when>
                  <xsl:when test="$cpf!=''" >
                    <label>CPF</label>
                  </xsl:when>
                  <xsl:when test="$idestrang!=''" >
                    <label>Id. Estrangeiro</label>
                  </xsl:when>
                  <xsl:otherwise>
                    <label>CNPJ/CPF/Id. Estrangeiro</label>
                  </xsl:otherwise>
                </xsl:choose>

                <span>
                  <xsl:call-template name="formatCnpj">
                    <xsl:with-param name="cnpj" select="$cnpj"/>
                  </xsl:call-template>
                  <xsl:call-template name="formatCpf">
                    <xsl:with-param name="cpf" select="$cpf"/>
                  </xsl:call-template>
                  <xsl:value-of select="$idestrang"/>
                </span>
              </td>
              <td class="col-2">
                <label>Nome / Razão Social</label>
                <span>
                  <xsl:value-of select="//n:infNFe/n:dest/n:xNome"/>
                </span>
              </td>
              <td class="col-5">
                <label>Inscrição Estadual</label>
                <span>
                  <xsl:value-of select="//n:infNFe/n:dest/n:IE"/>
                </span>
              </td>
              <td class="col-10">
                <label>UF</label>
                <span>
                  <xsl:value-of select="//n:infNFe/n:dest/n:enderDest/n:UF"/>
                </span>
              </td>
            </tr>
            <tr>
              <td class="col-5">
                <label>Destino da operação</label>
                <span>
                  <xsl:variable name="idDest" select="//n:idDest"/>
                  <xsl:choose>
                    <xsl:when test="$idDest='1'">
                      1 - Operação Interna
                    </xsl:when>
                    <xsl:when test="$idDest='2'">
                      2 - Operação Interestadual
                    </xsl:when>
                    <xsl:when test="$idDest='3'">
                      3 - Operação com Exterior
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$idDest"/>
                    </xsl:otherwise>
                  </xsl:choose>  
                </span>
              </td>
              <td class="col-2">
                <label>Consumidor final</label>
                <span>
                  <xsl:variable name="indFinal" select="//n:indFinal"/>
                  <xsl:choose>
                    <xsl:when test="$indFinal='0'">
                      0 - Normal
                    </xsl:when>
                    <xsl:when test="$indFinal='1'">
                      1 - Consumidor final
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$indFinal"/>
                    </xsl:otherwise>
                  </xsl:choose> 
                </span>
              </td>
              <td colspan="2" class="col-3">
                <label>Presença do Comprador</label>
                <span>
                  <xsl:variable name="indPres" select="//n:indPres"/>
                  <xsl:choose>
                    <xsl:when test="$indPres='0'">
                      0 - Não se aplica
                    </xsl:when>
                    <xsl:when test="$indPres='1'">
                      1 - Operação presencial
                    </xsl:when>
                    <xsl:when test="$indPres='2'">
                      2 - Operação pela internet
                    </xsl:when>
                    <xsl:when test="$indPres='3'">
                      3 - Operação não presencial (teleatendimento)
                    </xsl:when>
                    <xsl:when test="$indPres='4'">
                      4 - NFC-e com entrega a domicílio
                    </xsl:when>
                    <xsl:when test="$indPres='9'">
                      9 - Operação não presencial (outros)
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$indPres"/>
                    </xsl:otherwise>
                  </xsl:choose> 
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
      </xsl:if>
      <fieldset>
        <legend>Emissão</legend>
        <table class="box">
          <tr>
            <td>
              <label>Processo</label>
              <span>
                <xsl:for-each select="//n:procEmi">
                  <xsl:variable name="proces" select="text()"/>
                  <xsl:choose>
                    <xsl:when test="$proces = 0">
                      0 - com aplicativo do Contribuinte
                    </xsl:when>
                    <xsl:when test="$proces = 1">
                      1 - avulsa pelo Fisco
                    </xsl:when>
                    <xsl:when test="$proces = 2">
                      2 - avulsa, pelo Contribuinte com seu Certificado Digital, através do site do Fisco
                    </xsl:when>
                    <xsl:when test="$proces = 3">
                      3 - pelo Contribuinte com aplicativo fornecido pelo Fisco
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$proces"/>
                    </xsl:otherwise>
                  </xsl:choose> 
                </xsl:for-each>
              </span>
            </td>
            <td>
              <label>Versão do Processo</label>
              <span>
                <xsl:for-each select="//n:verProc">
                  <xsl:value-of select="text()"/>
                </xsl:for-each>
              </span>
            </td>
            <td>
              <label>
                Tipo de Emissão<!--Forma-->
              </label>
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
                    <xsl:otherwise><xsl:value-of select="$foremis"/></xsl:otherwise>
                  </xsl:choose> 
                </xsl:for-each>
              </span>
            </td>
            <td>
              <label>Finalidade</label>
              <span>
                <xsl:for-each select="//n:finNFe">
                  <xsl:variable name="finNFe" select="text()"/>
                  <xsl:choose>
                    <xsl:when test="$finNFe='1'">
                      1 - Normal
                    </xsl:when>
                    <xsl:when test="$finNFe='2'">
                      2 - complementar
                    </xsl:when>
                    <xsl:when test="$finNFe='3'">
                      3 - de Ajuste
                    </xsl:when>
                    <xsl:when test="$finNFe='4'">
                      4 - devolução de mercadoria
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$finNFe"/>
                    </xsl:otherwise>
                  </xsl:choose> 
                </xsl:for-each>
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <label>Natureza da Operação</label>
              <span>
                <xsl:for-each select="//n:natOp">
                  <xsl:value-of select="text()"/>
                </xsl:for-each>
              </span>
            </td>
            <td>
              <label>Tipo da Operação</label>
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
              <label>Forma de Pagamento</label>
              <span>
                <xsl:for-each select="//n:indPag">
                  <xsl:variable name="indPag" select="text()"/>
                  <xsl:choose>
                    <xsl:when test="$indPag = '0'">
                      0 - À vista
                    </xsl:when>
                    <xsl:when test="$indPag = '1'">
                      1 - A prazo
                    </xsl:when>
                    <xsl:when test="$indPag = '2'">
                      2 - Outros
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="$indPag"/>
                    </xsl:otherwise>
                  </xsl:choose> 
                </xsl:for-each>
              </span>
            </td>
            <td>
              <label>
                <i>Digest</i> Value da NF-e
              </label>
              <span>
                <xsl:variable name="digestNfe" select="//s:Signature/s:SignedInfo/s:Reference/s:DigestValue"/>
                <xsl:if test="$digestNfe != ''">
                  <xsl:value-of select="$digestNfe"/>
                </xsl:if>
              </span>
            </td>
          </tr>
        </table>
      </fieldset>
      <fieldset>
        <xsl:variable name="nota_cancelada" select="//n:cancNFe/n:infCanc/n:xServ"/>
        <legend>
          Situação Atual:
          <xsl:if  test="//n:infProt/n:cStat=301">&#160;DENEGADA</xsl:if>
          <xsl:if  test="//n:infProt/n:cStat=302">&#160;DENEGADA</xsl:if>
          <xsl:if  test="//n:infProt/n:cStat!=301 and //n:infProt/n:cStat!=302">
            <!-- PODE TER SIDO CANCELADA POR EVENTO-->
            <xsl:variable name="evento_cancelamento" select="//n:retEvento/n:infEvento[n:tpEvento=110111]/n:nProt"/>
            <xsl:if test="count($evento_cancelamento)=1">&#160;CANCELADA</xsl:if>
            <xsl:if test="count($evento_cancelamento)=0">
              <xsl:if test="count($nota_cancelada)=1">&#160;CANCELADA</xsl:if>
              <xsl:if test="count($nota_cancelada)=0">&#160;AUTORIZADA</xsl:if>
            </xsl:if>
          </xsl:if>


          (Ambiente de autorização:
          <xsl:variable name="tpAmb" select="//n:infNFe/n:ide/n:tpAmb"/>
          <xsl:choose>
            <xsl:when test="$tpAmb = 1"> produção</xsl:when>
            <xsl:when test="$tpAmb = 2"> homologação</xsl:when>
            <xsl:otherwise> <xsl:value-of select="$tpAmb"/></xsl:otherwise>
          </xsl:choose>)
        </legend>
        <table class="box">
          <tr>
            <td>
              <label>Eventos da NF-e</label>
            </td>
            <td>
              <label>Protocolo</label>
            </td>
            <td>
              <label>Data Autorização</label>
            </td>
            <td>
              <label>Data Inclusão BD</label>
            </td>
          </tr>
          <tr>
            <td>
              <span>
                <xsl:variable name="cStat" select="//n:infProt/n:cStat"/>
                <xsl:choose>
                  <xsl:when test="$cStat = '301'">
                    Denegação de Uso - Situação do emitente (Cod: 110101)
                  </xsl:when>
                  <xsl:when test="$cStat = '302'">
                    Denegação de Uso - Situação do destinatário (Cod: 110101)
                  </xsl:when>
                  <xsl:when test="$cStat = '100'">
                    Autorização de Uso (Cód.: 110100)
                  </xsl:when>
                  <xsl:when test="$cStat = '150'">
                    Autorização de Uso Fora de Prazo (Cód.: 110100)
                  </xsl:when>  
                  <xsl:otherwise>
                    <xsl:value-of select="$cStat"/> - <xsl:value-of select="//n:infProt/n:xMotivo"/> (Cód.: )
                  </xsl:otherwise>
                </xsl:choose> 
              </span>
            </td>
            <td>
              <span>
                <xsl:value-of select="//n:infProt/n:nProt"/>
                <xsl:element name="input">
                  <xsl:attribute name="type">
                    <xsl:value-of select="'hidden'"/>
                  </xsl:attribute>
                  <xsl:attribute name="id">
                    <xsl:value-of select="'nProt'"/>
                  </xsl:attribute>
                  <xsl:attribute name="value">
                    <xsl:value-of select="//n:infProt/n:nProt"/>
                  </xsl:attribute>
                </xsl:element>
              </span>
            </td>
            <td>
              <span>
                <xsl:if test="//n:mod=65">
                  <xsl:call-template name="formatDateTimeFuso">
                    <xsl:with-param name="dateTime" select="//n:infProt/n:dhRecbto"/>
                    <xsl:with-param name="include_as" select="1"/>
                  </xsl:call-template>
                </xsl:if>
                <xsl:if test="//n:mod=55">
                  <xsl:call-template name="formatDateTime">
                    <xsl:with-param name="dateTime" select="//n:infProt/n:dhRecbto"/>
                    <xsl:with-param name="include_as" select="1"/>
                  </xsl:call-template>
                </xsl:if> 
              </span>
            </td>
            <td>
              <span>
                <xsl:value-of select="//n:infProt/extDthInclusaoBdAutorizacao"/>                
              </span>
            </td>
          </tr>
          <xsl:if test="$nota_cancelada != ''">
            <tr>
              <td>
                <span>Cancelamento pelo emitente (Cód.: 110111)</span>
              </td>
              <td>
                <span>
                  <a id="lnkCce"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoCanc', 110111);">
                    <xsl:value-of select="//n:retCancNFe/n:infCanc/n:nProt"/>
                  </a>
                </span>
              </td>
              <td>
                <span>
                  <xsl:call-template name="formatDateTime">
                    <xsl:with-param name="dateTime" select="//n:retCancNFe/n:infCanc/n:dhRecbto"/>
                    <xsl:with-param name="include_as" select="1"/>
                  </xsl:call-template>
                </span>
              </td>
            </tr>
          </xsl:if>

          <!--INICIO - Acréscimo Eventos-->
          <xsl:variable name="ultCCC" select="(//n:retEvento/n:infEvento[n:tpEvento=110110]/n:nProt)[last()]"/>
          <xsl:for-each select="//n:NFeLog/n:eveNFe">

            <xsl:variable name="prot" select="n:retEvento/n:infEvento/n:nProt" />
            <xsl:variable name="tipo" select="n:retEvento/n:infEvento/n:tpEvento" />
              <tr>
              <td>
                <span>
                  <xsl:choose>
                    <xsl:when test="$tipo = '110110'">
                      Carta de Correção Eletrônica (Cód.: 110110)
                    </xsl:when>
                    <xsl:when test="$tipo = '210200'">
                      Confirmação da Operação pelo Destinatário(Cód.: 210200)
                    </xsl:when>
                    <xsl:when test="$tipo = '210210'">
                      Ciência da Operação pelo Destinatário (Cód.: 210210)
                    </xsl:when>
                    <xsl:when test="$tipo = '210220'">
                      Desconhecimento da Operação pelo Destinatário (Cód.: 210220)
                    </xsl:when>
                    <xsl:when test="$tipo = '210240'">
                      Operação não Realizada (Cód.: 210240)
                    </xsl:when>
                    <xsl:when test="$tipo = '110111'">
                      Cancelamento pelo emitente (Cód.: 110111)
                    </xsl:when>
                    <xsl:when test="$tipo = '990910'">
                      Internalização SUFRAMA (Cód.: 990910)
                    </xsl:when>
                    <xsl:when test="$tipo = '610500'">
                      Registro Passagem NF-e (Cód.: 610500)
                    </xsl:when>
                    <xsl:when test="$tipo = '610550'">
                      Registro Passagem NF-e BRId(Cód.: 610550)
                    </xsl:when>
                    <xsl:when test="$tipo = '610501'">
                      Cancelamento Registro Passagem NF-e (Cód.: 610501)
                    </xsl:when>
                    <xsl:when test="$tipo = '110140'">
                      EPEC-Emissão em Contingência (Cód.: 110140)
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="n:evento/n:infEvento/n:detEvento/n:descEvento"/> (Cód.: <xsl:value-of select="$tipo"/>)
                    </xsl:otherwise>
                  </xsl:choose> 
                </span>
              </td>
              <td>
                <span>
                  <xsl:choose>
                    <xsl:when test="$tipo = '110110'">
                      <xsl:variable name="countCCe" select="//n:infEvento[n:tpEvento=110110]/n:nProt" />
                      
                      <xsl:choose>
                        <xsl:when test="$ultCCC = $prot">
                          <a id="lnkCce"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('CCe', EventosEnum.CCE);">
                            <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                          </a>
                          <div style="display:none">
                            <xsl:call-template name="Evento_CCe"/>
                          </div> 
                        </xsl:when>
                        <xsl:otherwise>
                          <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                        </xsl:otherwise>
                      </xsl:choose>
                    </xsl:when>
                    <xsl:when test="$tipo = '110111'">
                      <a id="lnkCce"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoCanc', EventosEnum.CANC);">
                        <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                      </a>
                      <div style="display:none">
                        <xsl:call-template name="Evento_Cancelamento"/>                        
                      </div>
                    </xsl:when>
                    <xsl:when test="$tipo = '210200'">
                      <a id="lnkCce"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('ConfirmacaoOperacao{$prot}', EventosEnum.CONF_DEST);">
                        <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                      </a>
                      <div style="display:none">
                        <xsl:call-template name="Evento_Confirmacao"/>
                      </div>
                    </xsl:when>
                    <xsl:when test="$tipo = '210210'">
                      <a id="lnkCce"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('CienciaOperacao{$prot}', EventosEnum.CIENCIA_OP_DEST);">
                        <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                      </a>
                      <div style="display:none">
                        <xsl:call-template name="Evento_Ciencia"/>
                      </div>
                    </xsl:when>
                    <xsl:when test="$tipo = '210220'">
                      <a id="lnkCce"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('DesconhecimentoOperacao{$prot}', EventosEnum.DESC_OP_DEST);">
                        <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                      </a>
                      <div style="display:none">
                        <xsl:call-template name="Evento_Desconhecimento"/>
                      </div>
                    </xsl:when>
                    <xsl:when test="$tipo = '210240'">
                      <a id="l
                         nkOperacaoNaoRealizada"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('OperacaoNaoRealizada{$prot}', EventosEnum.OP_NREALIZADA);" >
                        <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                      </a>
                      <div style="display:none">
                        <xsl:call-template name="Evento_Nao_Realizado"/>
                      </div>
                    </xsl:when>
                    <xsl:when test="$tipo = '610600'">
                      <xsl:choose>
                        <xsl:when test="$ambiente = 'intranet'">
                          <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoCTeAutorizado{$prot}', EventosEnum.CTE_AUT);">
                            <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                          </a>
                          <div style="display:none">
                            <xsl:call-template name="Evento_CTe_Autorizado"/>
                          </div>
                        </xsl:when>
                        <xsl:otherwise>
                          <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                        </xsl:otherwise>
                      </xsl:choose> 
                    </xsl:when>
                    <xsl:when test="$tipo = '610601'"> 
                      <xsl:choose>
                        <xsl:when test="$ambiente = 'intranet'">
                          <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoCTeCancelado{$prot}', EventosEnum.CANC_CTE_AUT);">
                            <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                          </a>
                          <div style="display:none">
                            <xsl:call-template name="Evento_CTe_Cancelado"/>
                          </div>
                        </xsl:when>
                        <xsl:otherwise>
                          <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                        </xsl:otherwise>
                      </xsl:choose> 
                    </xsl:when> 
                    
                    <xsl:when test="$tipo = '990900'">
                      <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoVistoriaSUFRAMA{$prot}', EventosEnum.VIST_SUFRAMA);">
                        <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                      </a>
                      <div style="display:none">
                        <xsl:call-template name="Evento_Vistoria_SUFRAMA"/>
                      </div>
                    </xsl:when>
                    <xsl:when test="$tipo = '990910'">
                      <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoInternalizacaoSUFRAMA{$prot}', EventosEnum.INT_SUFRAMA);">
                        <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                      </a>
                      <div style="display:none">
                        <xsl:call-template name="Evento_Internalizacao_SUFRAMA"/>
                      </div>
                    </xsl:when>
                    <xsl:when test="$tipo = '610500'">
                      <xsl:choose>
                        <xsl:when test="$ambiente = 'intranet'">
                          <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoRegistroPassagemNFe{$prot}', EventosEnum.REG_PAS);">
                            <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                          </a>
                          <div style="display:none">
                            <xsl:call-template name="Evento_Registro_Passagem"/>
                          </div>
                        </xsl:when>
                        <xsl:when test="$ambiente = 'sistema_externo' and n:evento/n:infEvento/n:detEvento/n:cOrgaoAutor = 43 ">
                          <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoRegistroPassagemNFe{$prot}', EventosEnum.REG_PAS);">
                            <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                          </a>
                          <div style="display:none">
                            <xsl:call-template name="Evento_Registro_Passagem_Sistema_Externo"/>
                          </div>
                        </xsl:when>
                        <xsl:otherwise>
                          <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                        </xsl:otherwise>
                      </xsl:choose>
                    </xsl:when>
                    <xsl:when test="$tipo = '610550'">
                      <xsl:choose>
                        <xsl:when test="$ambiente = 'intranet'">
                          <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoRegistroPassagemNFeBRId{$prot}', EventosEnum.REG_PAS_BRID);">
                            <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                          </a>
                          <div style="display:none">
                            <xsl:call-template name="Evento_Registro_Passagem_BRId"/>
                          </div>
                        </xsl:when>
                        <xsl:otherwise>
                          <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                        </xsl:otherwise>
                      </xsl:choose>
                    </xsl:when>
                    <xsl:when test="$tipo = '610501'"> 
                      <xsl:choose>
                        <xsl:when test="$ambiente = 'intranet'">
                          <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoCancelamentoRegistroPassagem{$prot}', EventosEnum.CANC_REG_PAS);">
                            <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                          </a>
                          <div style="display:none">
                            <xsl:call-template name="Evento_Cancelamento_Registro_Passagem"/>
                          </div>
                        </xsl:when>
                        <xsl:otherwise>
                          <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                        </xsl:otherwise>
                      </xsl:choose> 
                    </xsl:when> 
                    <xsl:when test="$tipo = '610610'">
                      <xsl:choose>
                        <xsl:when test="$ambiente = 'intranet'">
                          <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoMDFeAutorizado{$prot}', EventosEnum.MDFE_AUT);">
                            <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                          </a>
                          <div style="display:none">
                            <xsl:call-template name="Evento_MDFe_Autorizado"/>
                          </div>
                        </xsl:when>
                        <xsl:otherwise>
                          <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                        </xsl:otherwise>
                      </xsl:choose>
                    </xsl:when>
                    <xsl:when test="$tipo = '610611'">
                      <xsl:choose>
                        <xsl:when test="$ambiente = 'intranet'">
                          <a id="lnkCTeAutorizado"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EventoMDFeCancelado{$prot}', EventosEnum.MDFE_CANC);">
                            <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                          </a>
                          <div style="display:none">
                            <xsl:call-template name="Evento_MDFe_Cancelado"/>
                          </div>
                        </xsl:when>
                        <xsl:otherwise>
                          <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                        </xsl:otherwise>
                      </xsl:choose>
                    </xsl:when>
                    <xsl:when test="$tipo = '110140'">
                      <a id="lnkCce"  class="linkCce"  href="javascript:;" onClick="javascript: visualizaEvento('EPEC{$prot}', EventosEnum.EPEC);">
                        <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                      </a>
                      <div style="display:none">
                        <xsl:call-template name="Evento_EPEC"/>
                      </div>
                    </xsl:when>
                    <xsl:otherwise>
                      <xsl:value-of select="n:retEvento/n:infEvento/n:nProt"/>
                    </xsl:otherwise>
                  </xsl:choose>                  
                </span>
              </td> 
              <td>
                <span>
                  <xsl:call-template name="formatDateTimeFuso">
                    <xsl:with-param name="dateTime" select="n:retEvento/n:infEvento/n:dhRegEvento"/>
                    <xsl:with-param name="include_as" select="1"/>
                  </xsl:call-template>
                </span>
              </td>
              <td>
                <span>
                  <xsl:value-of select="n:extDthInclusaoBdEvento"/>
                </span>
              </td>
            </tr>
          </xsl:for-each>
          <!--TERMINO - Acréscimo Eventos-->

          <!-- 
            INICIO BLOCO EXCLUSIVO RS - REMOVER AO DISTRIBUIR XSLT
          -->

          <!-- CONSULTA AO CMT - VERIFICACAO DE OTC  -->
          <xsl:for-each select="//n:CmtData/n:rp_regPassagem">

            <!-- Preparando json para popup -->
            <xsl:element name="input">
              <xsl:attribute name="type">
                <xsl:value-of select="'hidden'"/>
              </xsl:attribute>
              <xsl:attribute name="id">
                <xsl:value-of select="concat('CmtOtc_', n:rp_codIntCmtRpas)"/>
              </xsl:attribute>
              <xsl:attribute name="value">
                <xsl:value-of select="concat(n:rp_NFe/n:rp_infNFe/n:rp_ide/n:rp_cUF, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_ide/n:rp_serie, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_ide/n:rp_nNF, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_ide/n:rp_tpEmis, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_ide/n:rp_tpNF, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_emit/n:rp_CNPJ, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_emit/n:rp_CPF, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_dest/n:rp_CNPJ, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_dest/n:rp_CPF, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_dest/n:rp_cUF, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_total/n:rp_ICMSTot/n:rp_vNF, ';',
                                          n:rp_postoFiscal, ';',
                                          n:rp_incidente/n:rp_msg, ';',
                                          n:rp_meioTransp/n:rp_cMod, ';',
                                          n:rp_meioTransp/n:rp_pVeic, ';',
                                          n:rp_meioTransp/n:rp_cUFVeic, ';',
                                          n:rp_meioTransp/n:rp_pCarreta, ';',
                                          n:rp_meioTransp/n:rp_cUFCarreta, ';',
                                          n:rp_meioTransp/n:rp_pCarreta2, ';',
                                          n:rp_meioTransp/n:rp_cUFCarreta2, ';',
                                          n:rp_meioTransp/n:rp_xIdent, ';',
                                          n:rp_dthRegPass, ';',
                                          n:rp_NFe/n:rp_infNFe/n:rp_ide/n:rp_dEmi, ';',
                                          n:rp_codIntCmtRpas, ';',
                                          n:rp_tipoSentidoPista)"
                              />
              </xsl:attribute>
            </xsl:element>
            <!-- Fim preparacao -->

            <tr>
              <td>
                <span>
                  <xsl:element name="a">
                    <xsl:attribute name="class">
                      <xsl:value-of select="'linkCce'" />
                    </xsl:attribute>
                    <xsl:attribute name="href">
                      <xsl:value-of select="'javascript:;'" />
                    </xsl:attribute>
                    <xsl:attribute name="onClick">
                      <xsl:value-of select="concat('visualizaCmt(', n:rp_codIntCmtRpas, ');')"/>
                    </xsl:attribute>
                    <!--<xsl:value-of select="n:rp_incidente/n:rp_msg"/>-->
                    Verificação de Trânsito
                  </xsl:element>
                </span>
              </td>
              <td>
                <span>
                  <xsl:element name="a">
                    <xsl:attribute name="class">
                      <xsl:value-of select="'linkCce'" />
                    </xsl:attribute>
                    <xsl:attribute name="href">
                      <xsl:value-of select="'javascript:;'" />
                    </xsl:attribute>
                    <xsl:attribute name="onClick">
                      <xsl:value-of select="concat('visualizaCmt(', n:rp_codIntCmtRpas, ');')"/>
                    </xsl:attribute>
                    <xsl:value-of select="concat('RP: ', n:rp_codIntCmtRpas)"/>
                  </xsl:element>
                </span>
              </td>

              <td>
                <span>
                  <xsl:call-template name="formatDateTime">
                    <xsl:with-param name="dateTime" select="n:rp_dthRegPass"/>
                    <xsl:with-param name="include_as" select="1"/>
                  </xsl:call-template>
                </span>
              </td>
            </tr>
            <!-- *** -->

          </xsl:for-each>

          <!-- 
            FIM BLOCO EXCLUSIVO RS - REMOVER AO DISTRIBUIR XSLT
          -->

        </table>
      </fieldset>
    </div> 
  </xsl:template> 
</xsl:stylesheet>
