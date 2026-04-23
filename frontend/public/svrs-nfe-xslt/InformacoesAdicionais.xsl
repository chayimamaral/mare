<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:output method="html"/>
  <xsl:include href="_Versao.xsl"/>
  <xsl:include href="Utils.xsl"/>
  <xsl:include href="_Estilos_Geral.xsl"/>
  <xsl:include href="_Scripts_Geral.xsl"/>
  <xsl:template match="/" name="Informacoes_Adicionais">
    <div id="Inf" class="GeralXslt">
      <fieldset>
        <legend class="titulo-aba">
          Informações Adicionais
        </legend>
        <div id="Versao">
          XSLT:<xsl:call-template name="Versao"/> 
        </div>
        <table class="box">
          <tr>
            <td colspan="3">
              <label>Formato de Impressão DANFE</label>
              <span>
                <xsl:variable name="tpImp" select="//n:ide/n:tpImp"/>
                <xsl:choose>
                  <xsl:when test="$tpImp = 0">0 - Sem geração de DANFE</xsl:when>
                  <xsl:when test="$tpImp = 1">1 - DANFE normal, retrato</xsl:when>
                  <xsl:when test="$tpImp = 2">2 - DANFE normal, paisagem</xsl:when>
                  <xsl:when test="$tpImp = 3">3 - DANFE Simplificado</xsl:when>
                  <xsl:when test="$tpImp = 4">4 - DANFE NFC-e</xsl:when>
                  <xsl:when test="$tpImp = 5">5 - DANFE NFC-e resumido</xsl:when>
                  <xsl:when test="$tpImp = 6">6 - DANFE NFC-e, mensagem eletrônica</xsl:when>
                  <xsl:otherwise><xsl:value-of select="$tpImp"/></xsl:otherwise>
                </xsl:choose> 
              </span>
            </td>
          </tr>
          <xsl:for-each select="//n:ide/n:dhCont">
            <tr class="col-2">
              <td>
                <label>Entrada em Contingência</label>
                <span>
                  <xsl:call-template name="formatDateTimeFuso">
                    <xsl:with-param name="dateTime" select="text()"/>
                  </xsl:call-template>
                </span>
              </td>
              <td colspan="2">
                <label>Justificativa</label>
                <span>
                  <xsl:value-of select="//n:ide/n:xJust"/>&#160;
                </span>
              </td>
            </tr>
          </xsl:for-each>
        </table> 
        
        <xsl:variable name="autObterXML" select="//n:autXML"/>
        <xsl:if test="$autObterXML != ''">
          <fieldset>
            <legend class="titulo-aba-interna">Autorizados a acessar o XML da NF-e</legend>
            <table class="box">  
              <xsl:for-each select="//n:autXML">
                <xsl:variable name="posicao" select="position()"/>
                <xsl:variable name="cnpjAut" select="n:CNPJ"/>
                <xsl:variable name="cpfAut" select="n:CPF"/>
                <tr>
                  <td> 
                    <label>
                        Autorizado <xsl:value-of select="$posicao"/> -
                        <xsl:if test="$cnpjAut!=''">
                          CNPJ 
                        </xsl:if>
                        <xsl:if test="$cpfAut!=''">
                          CPF
                        </xsl:if>
                    </label>
                    <span>
                      <xsl:call-template name="formatCnpj">
                        <xsl:with-param name="cnpj" select="$cnpjAut"/>
                      </xsl:call-template>
                      <xsl:call-template name="formatCpf">
                        <xsl:with-param name="cpf" select="$cpfAut"/>
                      </xsl:call-template>
                    </span>
                  </td>
                </tr>
              </xsl:for-each>
              </table>
          </fieldset>
        </xsl:if>

        <xsl:for-each select="//n:exporta">
          <fieldset>
            <legend class="titulo-aba-interna">Exportação</legend>
            <table class="box">
              <tr class="col-3"> 
                  <td>
                    <label>Local de Embarque</label>
                    <span>
                      <xsl:choose>
                        <xsl:when test="n:xLocEmbarq != ''">
                          <xsl:value-of select="n:xLocEmbarq"/>
                        </xsl:when>
                        <xsl:when test="n:xLocExporta != ''">
                          <xsl:value-of select="n:xLocExporta"/>
                        </xsl:when>
                      </xsl:choose> 
                    </span>
                  </td> 
                  <td>
                    <label>UF de Embarque</label>
                    <span>

                      <xsl:choose>
                        <xsl:when test="n:UFEmbarq != ''">
                          <xsl:value-of select="n:UFEmbarq"/>
                        </xsl:when>
                        <xsl:when test="n:UFSaidaPais != ''">
                          <xsl:value-of select="n:UFSaidaPais"/>
                        </xsl:when>
                      </xsl:choose> 
                    </span>
                  </td>
                <td>
                  <label>Descrição Local Despacho</label>
                  <span>
                    <xsl:value-of select="n:xLocDespacho"/>
                  </span>
                </td>
              </tr>
            </table>
          </fieldset>
          </xsl:for-each>
        <xsl:variable name="compra" select="//n:compra"/>
        <xsl:if test="$compra != ''">
          <fieldset>
            <legend class="titulo-aba-interna">Informações de Compra</legend>
            <table class="box">
              <tr class="col-3">
               
                  <td>
                    <label>Nota de Empenho</label>
                    <span>
                      <xsl:value-of select="//n:compra/n:xNEmp"/>
                    </span>
                  </td> 
                  <td>
                    <label>Pedido</label>
                    <span>
                      <xsl:value-of select="//n:compra/n:xPed"/>
                    </span>
                  </td> 
                  <td>
                    <label>Contrato</label>
                    <span>
                      <xsl:value-of select="//n:compra/n:xCont"/>
                    </span>
                  </td> 
              </tr>
            </table>
          </fieldset>
        </xsl:if>
        <xsl:variable name="infocana" select="//n:cana"/>
        <xsl:if test="$infocana !=''">
          <fieldset>
            <legend class="titulo-aba-interna">Informações do Registro de Aquisição de Cana</legend>
            <table class="box" style="border-bottom:0px">
              <tr>
                <xsl:for-each select="//n:cana/n:safra">
                  <td>
                    <label>Identificação da Safra</label>
                    <span>
                      <xsl:value-of select="text()"/>
                    </span>
                  </td>
                </xsl:for-each>
                <xsl:for-each select="//n:cana/n:ref">
                  <td>
                    <label>Mês e Ano de Referência</label>
                    <span>
                      <xsl:value-of select="text()"/>
                    </span>
                  </td>
                </xsl:for-each>
              </tr>
            </table>
            <table class="box" style="border-top:0px">
              <tr class="col-3">
                  <td>
                    <label>Quantidade Total do Mês</label>
                    <span> 
                      <xsl:call-template name="format10Casas">
                        <xsl:with-param name="num" select="//n:cana/n:qTotMes"/>
                      </xsl:call-template>  
                    </span>
                  </td>
                  <td>
                    <label>Quantidade Total Anterior</label>
                    <span>
                      <xsl:call-template name="format10Casas">
                        <xsl:with-param name="num" select="//n:cana/n:qTotAnt"/>
                      </xsl:call-template> 
                    </span>
                  </td>
                  <td>
                    <label>Quantidade Total Geral</label>
                    <span>
                      <xsl:call-template name="format10Casas">
                        <xsl:with-param name="num" select="//n:cana/n:qTotGer"/>
                      </xsl:call-template>                      
                    </span>
                  </td>                
              </tr>
              <tr class="col-3">
                <xsl:for-each select="//n:cana/n:vFor">
                  <td>
                    <label>Valor dos Fornecimentos</label>
                    <span>
                      <xsl:call-template name="format2Casas">
                        <xsl:with-param name="num" select="text()"/>
                      </xsl:call-template>
                    </span>
                  </td>
                </xsl:for-each>
                <xsl:for-each select="//n:cana/n:vTotDed">
                  <td>
                    <label>Valor Total da Dedução</label>
                    <span>
                      <xsl:call-template name="format2Casas">
                        <xsl:with-param name="num" select="text()"/>
                      </xsl:call-template>
                    </span>
                  </td>
                </xsl:for-each>
                <xsl:for-each select="//n:cana/n:vLiqFor">
                  <td>
                    <label>Valor Líquido dos Fornecimentos</label>
                    <span>
                      <xsl:call-template name="format2Casas">
                        <xsl:with-param name="num" select="text()"/>
                      </xsl:call-template>
                    </span>
                  </td>
                </xsl:for-each>
              </tr>
            </table>
            <xsl:variable name="deduccana" select="//n:cana/n:deduc"/>
            <xsl:if test="$deduccana != ''">
              <fieldset>
                <legend class="titulo-aba-interna">Deduções - Taxas e Contribuições</legend>
                <table class="box">
                  <tr class="col-2">
                    <td>
                      <label>Descrição da Dedução</label>
                    </td>
                    <td>
                      <label>Valor da Dedução</label>
                    </td>
                  </tr>
                  <xsl:for-each select="//n:cana/n:deduc">
                    <tr class="col-2">
                      <td>
                        <span>
                          <xsl:value-of select="n:xDed"/>
                        </span>
                      </td>
                      <td>
                        <span>
                          <xsl:call-template name="format2Casas">
                            <xsl:with-param name="num" select="n:vDed"/>
                          </xsl:call-template>
                        </span>
                      </td>
                    </tr>
                  </xsl:for-each>
                </table>
              </fieldset>
            </xsl:if>
            <xsl:variable name="diariocana" select="//n:cana/n:forDia"/>
            <xsl:if test="$diariocana != ''">
              <fieldset>
                <legend class="titulo-aba-interna">Fornecimento Diário de Cana</legend>
                <table class="box">
                  <tr class="col-2">
                    <td>
                      <label>Dia</label>
                    </td>
                    <td>
                      <label>Quantidade</label>
                    </td>
                  </tr>
                  <xsl:for-each select="//n:cana/n:forDia">
                    <tr class="col-2">
                      <td>
                        <span>
                          <xsl:value-of select="@dia"/>
                        </span>
                      </td>
                      <td>
                        <span>
                          <xsl:value-of select="n:qtde"/>
                        </span>
                      </td>
                    </tr>
                  </xsl:for-each>
                </table>
              </fieldset>
            </xsl:if>
          </fieldset>
        </xsl:if>
        <xsl:variable name="procreferenciado" select="//n:infAdic/n:procRef"/>
        <xsl:if test="$procreferenciado != ''">
          <fieldset>
            <legend class="titulo-aba-interna">Processo Referenciado</legend>
            <xsl:for-each select="//n:infAdic/n:procRef">
              <xsl:variable name="procRef" select="position()"/>
              <div>
                <h5 class="tggle">
                  Processo Referenciado <xsl:value-of select=" position()"/>
                </h5>
                <table class="toggable box">
                  <tr class="col-3">
                    <td colspan="2">
                      <label>Identificador / Ato concessório</label>
                      <span>
                        <xsl:value-of select="n:nProc"/>
                      </span>
                    </td>
                    <td>
                      <label>Indicador da Origem</label>
                      <span>
                        <xsl:variable name="vindProc" select="n:indProc"/>
                        <xsl:choose>
                          <xsl:when test="$vindProc = 0">
                            0 - SEFAZ
                          </xsl:when>
                          <xsl:when test="$vindProc = 1">
                            1 - Justiça Federal
                          </xsl:when>
                          <xsl:when test="$vindProc = 2">
                            2 - Justiça Estadual
                          </xsl:when>
                          <xsl:when test="$vindProc = 3">
                            3 - Secex/RFB
                          </xsl:when>
                          <xsl:when test="$vindProc = 9">
                            9 - Outros
                          </xsl:when>
                          <xsl:otherwise>
                            <xsl:value-of select="$vindProc"/>
                          </xsl:otherwise>
                        </xsl:choose> 
                      </span>
                    </td>
                  </tr>
                </table>
              </div>
            </xsl:for-each>
          </fieldset>
        </xsl:if>
        <xsl:variable name="infoadicionalfisco" select="//n:infAdic/n:infAdFisco"/>
        <xsl:if test="$infoadicionalfisco != ''">
          <fieldset>
            <legend class="titulo-aba-interna">Informações Adicionais de Interesse do Fisco</legend>
            <table class="box">
              <tr>
                <td>
                  <label>Descrição</label>
                  <span>
                    <xsl:value-of select="//n:infAdic/n:infAdFisco"/>
                  </span>
                </td>
              </tr>
            </table>
          </fieldset>
        </xsl:if>
        <xsl:variable name="observacoesfisco" select="//n:infAdic/n:obsFisco"/>
        <xsl:if test="$observacoesfisco != ''">
          <fieldset>
            <legend class="titulo-aba-interna">Observações do Fisco</legend>
            <table class="box">
              <tr>
                <td>
                  <label>Campo</label>
                </td>
                <td>
                  <label>Texto</label>
                </td>
              </tr>
              <xsl:for-each select="//n:infAdic/n:obsFisco">
                <tr class="col-2">
                  <td class="fixo-info-adic-proc-ref">
                    <span>
                      <xsl:value-of select="@xCampo"/>
                    </span>
                  </td>
                  <td>
                    <span>
                      <xsl:value-of select="n:xTexto"/>
                    </span>
                  </td>
                </tr>
              </xsl:for-each>
            </table>
          </fieldset>
        </xsl:if>
        <xsl:variable name="infocpl" select="//n:infAdic/n:infCpl"/>
        <xsl:if test="$infocpl != ''">
          <fieldset>
            <legend>Informações Complementares de Interesse do Contribuinte</legend> 
            <table class="box">
              <tr>
                <td>
                    <label>Descrição</label>
                    <span>
                      <xsl:variable name="infCpl" select="//n:infAdic/n:infCpl"/>
                      <xsl:call-template name="quebraLinha">
                        <xsl:with-param name="infCpl" select="$infCpl"/>
                      </xsl:call-template>
                    </span>
                </td>
              </tr>
            </table> 
          </fieldset>
        </xsl:if>
        <xsl:variable name="obscontribuinte" select="//n:infAdic/n:obsCont"/>
        <xsl:if test="$obscontribuinte != ''">
          <fieldset>
            <legend class="titulo-aba-interna">Observações do Contribuinte</legend>
            <table class="box">
              <tr>
                <td class="fixo-info-adic-proc-ref">
                  <label>Campo</label>
                </td>
                <td>
                  <label>Texto</label>
                </td>
              </tr>
              <xsl:for-each select="//n:infAdic/n:obsCont">
                <xsl:variable name="obsContribuinte" select="position()"/>
                <tr class="col-2">
                  <td>
                    
                      <span>
                        <xsl:value-of select="@xCampo"/>
                      </span>
                    
                  </td>
                  <td>
                    <span>
                      <xsl:value-of select="n:xTexto"/>
                    </span>
                  </td>
                </tr>
              </xsl:for-each>
            </table>
          </fieldset>
        </xsl:if>

        <!-- Documentos Referenciados -->

        <xsl:variable name="NFreferenciada" select="//n:ide/n:NFref"/>
        <xsl:if test="$NFreferenciada != ''">
          <fieldset>
            <legend>Documentos Fiscais Referenciados</legend>
            <xsl:for-each select="//n:ide/n:NFref">
              <!-- NF - Nota Fiscal -->
              <xsl:for-each select="n:refNF">
                <fieldset>
                  <legend class="titulo-aba-interna">Nota Fiscal</legend>
                  <table class="box">
                    <tr class="col-3">
                      <xsl:for-each select="n:cUF">
                        <td>
                          <label>Código da UF</label>
                          <span>
                            <xsl:value-of select="text()"/>
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select="n:AAMM">
                        <td>
                          <label>Ano / Mês</label>
                          <span>
                            <xsl:variable name="dAAMM" select="text()"/>
                            <xsl:value-of select="$dAAMM"/>
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select="n:CNPJ">
                        <td>
                          <label>CNPJ</label>
                          <span>
                            <xsl:call-template name="formatCnpj">
                              <xsl:with-param name="cnpj" select="text()"/>
                            </xsl:call-template>
                          </span>
                        </td>
                      </xsl:for-each>
                    </tr>
                    <tr>
                      <xsl:for-each select="n:mod">
                        <td>
                          <label>Modelo</label>
                          <span>
                            <xsl:variable name="vmod" select="text()"/>
                            <xsl:choose>
                              <xsl:when test="$vmod='01'">
                                01 - Modelo 01.
                              </xsl:when>
                              <xsl:otherwise>
                                <xsl:value-of select="$vmod"/>
                              </xsl:otherwise>
                            </xsl:choose> 
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select="n:serie">
                        <td>
                          <label>Série</label>
                          <span>
                            <xsl:value-of select="text()"/>
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select="n:nNF">
                        <td>
                          <label>Número</label>
                          <span>
                            <xsl:value-of select="text()"/>
                          </span>
                        </td>
                      </xsl:for-each>
                    </tr>
                  </table>
                </fieldset>
              </xsl:for-each>

              <!-- NFE - Nota Fiscal Eletrônica -->
              <xsl:for-each  select ="n:refNFe">
                <fieldset>
                  <legend class="titulo-aba-interna">Nota Fiscal Eletrônica</legend>
                  <table class="box">
                    <tr class="col-3">
                      <td colspan="3">
                        <label>Chave de Acesso</label>
                        <span>
                          <xsl:call-template name="formatNfe">
                            <xsl:with-param name="nfe" select="text()"/>
                          </xsl:call-template>
                        </span>
                      </td>
                    </tr>
                  </table>

                </fieldset>
              </xsl:for-each>

              <!-- NFP - Nota Fiscal de Produtor Rural -->
              <xsl:for-each select ="n:refNFP">
                <fieldset>
                  <legend class="titulo-aba-interna">Nota Fiscal de Produtor Rural</legend>
                  <table class="box">
                    <tr>
                      <td>
                        <table class="box">
                          <tr class="col-4">
                            <td>
                              <label>Código da UF</label>
                              <span>
                                <xsl:for-each select="n:cUF">
                                  <xsl:value-of select="text()"/>
                                </xsl:for-each>
                              </span>
                            </td>
                            <td>
                              <label>Ano / Mês</label>
                              <span>
                                <xsl:for-each select="n:AAMM">
                                  <xsl:variable name="dAAMM" select="text()"/>
                                  <xsl:value-of select="$dAAMM"/>
                                </xsl:for-each>
                              </span>
                            </td>
                            <td> 
                              <xsl:variable name="cnpjProd" select="n:CNPJ"/>
                              <xsl:variable name="cpfProd" select="n:CPF"/>
                              <xsl:if test="$cnpjProd!=''">
                                <label>CNPJ</label>
                              </xsl:if>
                              <xsl:if test="$cpfProd!=''">
                                <label>CPF</label>
                              </xsl:if>
                              <span>
                                <xsl:call-template name="formatCnpj">
                                  <xsl:with-param name="cnpj" select="$cnpjProd"/>
                                </xsl:call-template>
                                <xsl:call-template name="formatCpf">
                                  <xsl:with-param name="cpf" select="$cpfProd"/>
                                </xsl:call-template>
                              </span> 
                            </td>
                            <td>
                              <label>IE</label>
                              <span>
                                <xsl:for-each select="n:IE">
                                  <xsl:value-of select="text()"/>
                                </xsl:for-each>
                              </span>
                            </td>
                          </tr>
                        </table>
                        <br/>
                        <table class="box">
                          <tr class="col-3">
                            <td>
                              <label>Modelo do Documento Fiscal</label>
                              <span>
                                <xsl:variable name="vmod" select="n:mod"/>

                                <xsl:choose>
                                  <xsl:when test="$vmod = 1">
                                    01 - NF
                                  </xsl:when>
                                  <xsl:when test="$vmod = 4">
                                    04 - NF de Produtor
                                  </xsl:when>
                                  <xsl:otherwise>
                                    <xsl:value-of select="$vmod"/>
                                  </xsl:otherwise>
                                </xsl:choose> 
                              </span>
                            </td>
                            <td>
                              <label>Série do Documento Fiscal</label>
                              <span>
                                <xsl:for-each select="n:serie">
                                  <xsl:value-of select="text()"/>
                                </xsl:for-each>
                              </span>
                            </td>
                            <td>
                              <label>Número do Documento Fiscal</label>
                              <span>
                                <xsl:for-each select="n:nNF">
                                  <xsl:value-of select="text()"/>
                                </xsl:for-each>
                              </span>
                            </td>
                          </tr>
                        </table>
                      </td>
                    </tr>
                  </table>
                </fieldset>
              </xsl:for-each>

              <!-- CTE - Conhecimento de Transporte Eletrônico -->
              <xsl:for-each select ="n:refCTe">
                <fieldset>
                  <legend class="titulo-aba-interna">Conhecimento de Transporte Eletrônico</legend>
                  <table class="box">
                    <tr class="col-3">
                      <td colspan="3">
                        <label>Chave de Acesso</label>
                        <span>
                          <xsl:call-template name="formatNfe">
                            <xsl:with-param name="nfe" select="text()"/>
                          </xsl:call-template>
                        </span>
                      </td>
                    </tr>
                  </table>
                </fieldset>
              </xsl:for-each>

              <!-- ECF - Emissor de Cupom Fiscal -->
              <xsl:for-each select="n:refECF">
                <fieldset>
                  <legend>Informações do Cupom Fiscal</legend>
                  <table class="box">
                    <tr class="col-3">
                        <xsl:variable name="refECFMod" select="n:mod"/>
                        <xsl:if test="$refECFMod != ''">
                          <td>
                            <label>Modelo&#160;de&#160;Documento&#160;Fiscal</label>
                            <span>
                              <xsl:choose>
                                <xsl:when test="$refECFMod = '2B'">
                                  2B - Cupom Fiscal emitido por máquina registradora (não ECF)  
                                </xsl:when>
                                <xsl:when test="$refECFMod = '2C'">
                                  2C - Cupom Fiscal PDV
                                </xsl:when>
                                <xsl:when test="$refECFMod = '2D'">
                                  2D - Cupom Fiscal (emitido por ECF)
                                </xsl:when>
                                <xsl:otherwise>
                                  <xsl:value-of select="$refECFMod"/>
                                </xsl:otherwise>
                              </xsl:choose> 
                            </span>
                          </td>
                        </xsl:if>  
                        <xsl:for-each select="n:nECF">
                          <td>
                              <label>Número&#160;de&#160;Ordem&#160;Sequencial&#160;do&#160;ECF</label>
                              <span>
                                  <xsl:value-of select="text()"/>
                              </span>
                          </td>
                      </xsl:for-each> 
                      <xsl:for-each select="n:nCOO">
                          <td>
                            <label>
                              Número&#160;do&#160;Contador&#160;de&#160;Ordem&#160;de&#160;Operação
                            </label>
                            <span>
                              <xsl:value-of select="text()"/>
                            </span>
                          </td>
                      </xsl:for-each>
                    </tr>
                  </table>
                </fieldset>
              </xsl:for-each> 
            </xsl:for-each>
          </fieldset>
        </xsl:if>
      </fieldset>
    </div>
  </xsl:template>
</xsl:stylesheet>