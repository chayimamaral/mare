<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <xsl:include href="Utils.xsl"/>
  <xsl:include href="_Estilos_Geral.xsl"/>
  <xsl:include href="_Scripts_Geral.xsl"/>
  <xsl:include href="Template_Impostos.xsl"/>
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:template match="/" name="Produtos_e_Servicos">
    <div id="Prod" class="GeralXslt">
      <fieldset>
        <legend class="titulo-aba">Dados dos Produtos e Serviços</legend>
        <div>
          <table class="prod-serv-header box">
            <tr>
              <td class="fixo-prod-serv-numero">
                <label>Num.</label>
              </td>
              <td class="fixo-prod-serv-descricao">
                <label>Descrição</label>
              </td>
              <td class="fixo-prod-serv-qtd">
                <label>Qtd.</label>
              </td>
              <td class="fixo-prod-serv-uc">
                <label>Unidade Comercial</label>
              </td>
              <td class="fixo-prod-serv-vb">
                <label>Valor(R$)</label>
              </td>
            </tr>
          </table>

          <xsl:for-each select="//n:infNFe/n:det">
            <table class="toggle box">
              <tr>
                <xsl:if test="position() mod 2 != 1">
                  <xsl:attribute name="class">highlighted</xsl:attribute>
                </xsl:if>
                <td class="fixo-prod-serv-numero">
                  <span>
                    <xsl:value-of select="position()"/>
                  </span>
                </td>
                <td class="fixo-prod-serv-descricao">
                  <span>
                    <xsl:value-of select="n:prod/n:xProd"/>
                  </span>
                </td>
                <td class="fixo-prod-serv-qtd">
                  <span>
                    <xsl:call-template name="format4Casas">
                      <xsl:with-param name="num" select ="n:prod/n:qCom"/>
                    </xsl:call-template>
                  </span>
                </td>
                <td class="fixo-prod-serv-uc">
                  <span>
                    <xsl:value-of select="n:prod/n:uCom"/>
                  </span>
                </td>
                <td class="fixo-prod-serv-vb">
                  <span>
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="n:prod/n:vProd"/>
                    </xsl:call-template>
                  </span>
                </td>
              </tr>
            </table>
            <table class="toggable box" style="background-color:#ECECEC">
              <tr>
                <td>
                  <table class="box">
                    <tr class="col-4">
                      <td colspan="4">
                        <label>Código do Produto</label>
                        <span>
                          <xsl:value-of select="n:prod/n:cProd"/>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Código NCM</label>
                        <span>
                          <xsl:value-of select="n:prod/n:NCM"/>
                        </span>
                      </td>
                      <td colspan="4">
                        <xsl:if test="n:prod/n:genero">
                          <label>Gênero</label>
                          <span>
                            <xsl:value-of select="n:prod/n:genero"/>
                          </span>
                        </xsl:if>                        
                        <xsl:if test="n:prod/n:NVE">
                          <label>NVE</label>
                          <span>
                            <xsl:for-each select="n:prod/n:NVE">
                              <xsl:value-of select="text()"/><xsl:if test="position() != last()">, </xsl:if>
                            </xsl:for-each> 
                          </span>
                        </xsl:if>
                      </td>
                    </tr>
                    <tr>
                      <td colspan="4">
                        <label>Código EX da TIPI</label>
                        <span>
                          <xsl:value-of select="n:prod/n:EXTIPI"/>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>CFOP</label>
                        <span>
                          <xsl:value-of select="n:prod/n:CFOP"/>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Outras Despesas Acessórias</label>
                        <span>
                          <xsl:call-template name="format2Casas">
                            <xsl:with-param name="num" select="n:prod/n:vOutro"/>
                          </xsl:call-template>
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td colspan="4">
                        <label>Valor do Desconto</label>
                        <span>
                          <xsl:call-template name="format2Casas">
                            <xsl:with-param name="num" select="n:prod/n:vDesc"/>
                          </xsl:call-template>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Valor Total do Frete</label>
                        <span>
                          <xsl:call-template name="format2Casas">
                            <xsl:with-param name="num" select="n:prod/n:vFrete"/>
                          </xsl:call-template>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Valor do Seguro</label>
                        <span>
                          <xsl:call-template name="format2Casas">
                            <xsl:with-param name="num" select="n:prod/n:vSeg"/>
                          </xsl:call-template>
                        </span>
                      </td>
                    </tr>
                  </table>
                  <br/>
                  <table class="box">
                    <tr class="col-12">
                      <td colspan="12">
                        <label>
                          Indicador de Composição do Valor Total da NF-e
                        </label>
                        <span>
                          <xsl:variable name="indTot" select="n:prod/n:indTot"/>
                          <xsl:choose>
                            <xsl:when test="$indTot ='0'">
                              0 - O valor do item (vProd) não compõe o valor total da NF-e (vProd)
                            </xsl:when>
                            <xsl:when test="$indTot ='1'">
                              1 - O valor do item (vProd) compõe o valor total da NF-e (vProd)
                            </xsl:when>
                            <xsl:otherwise>
                              <xsl:value-of select="$indTot"/>
                            </xsl:otherwise>
                          </xsl:choose> 
                        </span>
                      </td>
                    </tr>
                    <tr class="col-3">
                      <td colspan="4">
                        <label>Código EAN Comercial</label>
                        <span>
                          <xsl:value-of select="n:prod/n:cEAN"/>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Unidade Comercial</label>
                        <span>
                          <xsl:value-of select="n:prod/n:uCom"/>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Quantidade Comercial</label>
                        <span>
                          <xsl:call-template name="format4Casas">
                            <xsl:with-param name="num" select="n:prod/n:qCom"/>
                          </xsl:call-template>
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td colspan="4">
                        <label>Código EAN Tributável</label>
                        <span>
                          <xsl:value-of select="n:prod/n:cEANTrib"/>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Unidade Tributável</label>
                        <span>
                          <xsl:value-of select="n:prod/n:uTrib"/>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Quantidade Tributável</label>
                        <span>
                          <xsl:call-template name="format4Casas">
                            <xsl:with-param name="num" select="n:prod/n:qTrib"/>
                          </xsl:call-template>
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td colspan="4">
                        <label>Valor unitário de comercialização</label>
                        <span>
                          <xsl:call-template name="format10Casas">
                            <xsl:with-param name="num" select="n:prod/n:vUnCom"/>
                          </xsl:call-template>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Valor unitário de tributação</label>
                        <span>
                          <xsl:call-template name="format10Casas">
                            <xsl:with-param name="num" select="n:prod/n:vUnTrib"/>
                          </xsl:call-template>
                        </span>
                      </td>
                      <td colspan="4"></td>
                    </tr>
                    <tr>
                      <td colspan="4">
                        <label>Número do pedido de compra</label>
                        <span>
                          <xsl:value-of select="n:prod/n:xPed"/>
                        </span>
                      </td>
                      <td colspan="4">
                        <label>Item do pedido de compra</label>
                        <span>
                          <xsl:value-of select="n:prod/n:nItemPed"/>
                        </span>
                      </td>
                      <td>
                        <label>Valor Aproximado dos Tributos</label>
                        <span>
                          <xsl:call-template name="format2Casas">
                            <xsl:with-param name="num" select="n:imposto/n:vTotTrib"/>
                          </xsl:call-template>                          
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td colspan="4">
                        <label>Número da FCI</label>
                        <span>
                          <xsl:value-of select="n:prod/n:nFCI"/>
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td colspan="12">
                        <fieldset>
                          <legend>ICMS Normal e ST</legend>
                          <xsl:choose>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMS00">
                              <xsl:call-template name="ICMS00"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMS10">
                              <xsl:call-template name="ICMS10"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMS20">
                              <xsl:call-template name="ICMS20"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMS30">
                              <xsl:call-template name="ICMS30"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMS40">
                              <xsl:call-template name="ICMS40"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMS51">
                              <xsl:call-template name="ICMS51"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMS60">
                              <xsl:call-template name="ICMS60"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMS70">
                              <xsl:call-template name="ICMS70"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMS90">
                              <xsl:call-template name="ICMS90"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMSPart" >
                              <xsl:call-template name="ICMSPart"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMSST" >
                              <xsl:call-template name="ICMSST"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMSSN101" >
                              <xsl:call-template name="ICMS_SN_101"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMSSN102" >
                              <xsl:call-template name="ICMS_SN_102"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMSSN201" >
                              <xsl:call-template name="ICMS_SN_201"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMSSN202" >
                              <xsl:call-template name="ICMS_SN_202"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMSSN500" >
                              <xsl:call-template name="ICMS_SN_500"/>
                            </xsl:when>
                            <xsl:when test="n:imposto/n:ICMS/n:ICMSSN900" >
                              <xsl:call-template name="ICMS_SN_900"/>
                            </xsl:when>
                          </xsl:choose>
                        </fieldset>

                        <!-- IPI -->
                        <xsl:variable name="ipi" select="n:imposto/n:IPI"/>
                        <xsl:if test="$ipi!=''">
                          <fieldset>
                            <legend>Imposto Sobre Produtos Industrializados</legend>
                            <table class="box">
                              <tr class="col-3">
                                <td>
                                  <label>Classe de Enquadramento</label>
                                  <span>
                                    <xsl:value-of select="n:imposto/n:IPI/n:clEnq"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Código de Enquadramento</label>
                                  <span>
                                    <xsl:value-of select="n:imposto/n:IPI/n:cEnq"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Código do Selo</label>
                                  <span>
                                    <xsl:value-of select="n:imposto/n:IPI/n:cSelo"/>
                                  </span>
                                </td>
                              </tr>
                              <tr>
                                <td>
                                  <label>CNPJ do Produtor</label>
                                  <span>
                                    <xsl:call-template name="formatCnpj">
                                      <xsl:with-param name="cnpj" select="n:imposto/n:IPI/n:CNPJProd"/>
                                    </xsl:call-template>
                                  </span>
                                </td>
                                <td>
                                  <label>Qtd. Selo</label>
                                  <span>
                                    <xsl:for-each select="n:imposto/n:IPI/n:qSelo">
                                      <xsl:value-of select="format-number(text(),'###.###.###.##0')"/>
                                    </xsl:for-each>
                                  </span>
                                </td>
                                <td>
                                  <xsl:choose>
                                    <xsl:when test="n:imposto/n:IPI/n:IPITrib/n:CST = (00 or 49 or 50 or 99)">
                                      <label>CST</label>
                                      <span>
                                        <xsl:variable name="origmercst" select="n:imposto/n:IPI/n:IPITrib/n:CST"/>
                                        <xsl:choose>
                                          <xsl:when test="$origmercst = 00">
                                            00 - Entrada com recuperação de crédito
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 49">
                                            49 - Outras entradas
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 50">
                                            50 - Saída tributada
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 99">
                                            99 - Outras saídas
                                          </xsl:when>
                                          <xsl:otherwise>
                                            <xsl:value-of select="$origmercst"/>
                                          </xsl:otherwise>
                                        </xsl:choose>  
                                      </span>
                                    </xsl:when>
                                    <xsl:otherwise>
                                      <label>CST</label>
                                      <span>
                                        <xsl:variable name="origmercst" select="n:imposto/n:IPI/n:IPINT/n:CST"/>
                                        <xsl:choose>
                                          <xsl:when test="$origmercst = 01">
                                            01-Entrada tributada com alíquota zero
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 02">
                                            02-Entrada isenta
                                          </xsl:when>
                                          <xsl:when test="$origmercst =  03">
                                            03-Entrada não tributada
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 04">
                                            04-Entrada imune
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 05">
                                            05-Entrada com suspensão
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 51">
                                            51-Saída tributada com alíquota zero
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 52">
                                            52-Saída isenta
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 53">
                                            53-Saída não-tributada
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 54">
                                            54-Saída imune
                                          </xsl:when>
                                          <xsl:when test="$origmercst = 55">
                                            55-Saída com suspensão
                                          </xsl:when>
                                          <xsl:otherwise>
                                            <xsl:value-of select="$origmercst"/>
                                          </xsl:otherwise>
                                        </xsl:choose>  
                                      </span>
                                    </xsl:otherwise>
                                  </xsl:choose>
                                </td>
                              </tr>
                              <!--xxxxxxxxxxxx  qtd Unidade padrão x vlr Unidade padrão. xxxxxxxxxxxxxxxx-->
                              <tr>
                                <td>
                                  <label>Qtd Total Unidade Padrão</label>
                                  <span>
                                    <xsl:call-template name="format4Casas">
                                      <xsl:with-param name="num" select="n:imposto/n:IPI/n:IPITrib/n:qUnid"/>
                                    </xsl:call-template>
                                  </span>
                                </td>
                                <td>
                                  <label>Valor por Unidade</label>
                                  <span>
                                    <xsl:call-template name="format4Casas">
                                      <xsl:with-param name="num" select="n:imposto/n:IPI/n:IPITrib/n:vUnid"/>
                                    </xsl:call-template>
                                  </span>
                                </td>
                                <td>
                                  <label>Valor IPI</label>
                                  <span>
                                    <xsl:call-template name="format2Casas">
                                      <xsl:with-param name="num" select="n:imposto/n:IPI/n:IPITrib/n:vIPI"/>
                                    </xsl:call-template>
                                  </span>
                                </td>
                              </tr>
                              <!--xxxxxxxxxxxx BASE DE CALCULO x ALIQUOTA xxxxxxxxxxxxxxxx-->
                              <tr>
                                <td>
                                 <label>Base de Cálculo</label>
                                    <span>
                                      <xsl:for-each select="n:imposto/n:IPI/n:IPITrib/n:vBC">
                                        <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
                                      </xsl:for-each>
                                    </span>
                                </td>
                                <td>
                                  <label>Alíquota</label>
                                  <span>
                                    <xsl:for-each select="n:imposto/n:IPI/n:IPITrib/n:pIPI">
                                        <xsl:value-of select="format-number(text(),'##0,0000')"/>
                                    </xsl:for-each>
                                  </span>
                                </td>
                              </tr>
                            </table>
                          </fieldset>
                        </xsl:if>
                        <xsl:variable name="ii" select="n:imposto/n:II"/>
                        <xsl:if test="$ii!=''">
                          <fieldset>
                            <legend>Imposto de Importação</legend>
                            <table class="box">
                              <tr class="col-3">
                                <td>
                                  <label>Base de Cálculo</label>
                                  <span>
                                    <xsl:call-template name="format2Casas">
                                      <xsl:with-param name="num" select="n:imposto/n:II/n:vBC"/>
                                    </xsl:call-template> 
                                  </span>
                                </td>
                                <td>
                                  <label>Despesas Aduaneiras</label>
                                  <span>
                                    <xsl:call-template name="format2Casas">
                                      <xsl:with-param name="num" select="n:imposto/n:II/n:vDespAdu"/>
                                    </xsl:call-template>
                                  </span>
                                </td>
                                <td>
                                  <label>Imposto de Importação</label>
                                  <span>
                                    <xsl:call-template name="format2Casas">
                                      <xsl:with-param name="num" select="n:imposto/n:II/n:vII"/>
                                    </xsl:call-template>
                                  </span>
                                </td>
                              </tr>
                              <tr>
                                <td>
                                  <label>IOF</label>
                                  <span>
                                    <xsl:for-each select="n:imposto/n:II/n:vIOF">
                                      <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
                                    </xsl:for-each>
                                  </span>
                                </td>
                              </tr>
                            </table>
                          </fieldset>
                        </xsl:if>
                        <xsl:variable name="di" select="n:prod/n:DI"/>
                        <xsl:if test="$di!=''">
                          <fieldset>
                            <legend>Documentos de Importação</legend>
                            <xsl:for-each select="n:prod/n:DI">
                              <div>
                                <h5 class="toggle">
                                  Documento de Importação nº <xsl:value-of select="position()"/>
                                </h5>
                                <table class="toggable box">
                                  <tr class="col-3">
                                    <td>
                                      <label>Número</label>
                                      <span>
                                        <xsl:value-of select="n:nDI"/>
                                      </span>
                                    </td>
                                    <td>
                                      <label>Data de registro</label>
                                      <span>
                                        <xsl:call-template name="formatDate">
                                          <xsl:with-param name="date" select="n:dDI"/>
                                        </xsl:call-template>
                                      </span>
                                    </td>
                                    <td colspan="3">
                                      <label>Código do Exportador</label>
                                      <span>
                                        <xsl:value-of select="n:cExportador"/>
                                      </span>
                                    </td>
                                  </tr>
                                  <tr>
                                    <td>
                                      <label>UF do desembaraço</label>
                                      <span>
                                        <xsl:value-of select="n:UFDesemb"/>
                                      </span>
                                    </td>
                                    <td>
                                      <label>Data do desembaraço</label>
                                      <span>
                                        <xsl:call-template name="formatDate">
                                          <xsl:with-param name="date" select="n:dDesemb"/>
                                        </xsl:call-template>
                                      </span>
                                    </td>
                                    <td colspan="3">
                                      <label>Local do desembaraço aduaneiro</label>
                                      <span>
                                        <xsl:value-of select="n:xLocDesemb"/>
                                      </span>
                                    </td>
                                  </tr> 
                                  <tr>
                                    <td>
                                      <label>Via Transporte Internacional</label>
                                      <span>
                                        
                                        <xsl:variable name="viaTransp" select="n:tpViaTransp"/>
                                        <xsl:choose>
                                          <xsl:when test="$viaTransp = 1">
                                            01 = Marítima
                                          </xsl:when>
                                          <xsl:when test="$viaTransp = 2">
                                            02 = Fluvial
                                          </xsl:when>
                                          <xsl:when test="$viaTransp = 3">
                                            03 = Lacustre
                                          </xsl:when>
                                          <xsl:when test="$viaTransp = 4">
                                            04 = Aérea
                                          </xsl:when>
                                          <xsl:when test="$viaTransp = 5">
                                            05 = Postal
                                          </xsl:when>
                                          <xsl:when test="$viaTransp = 6">
                                            06 = Ferroviária
                                          </xsl:when>
                                          <xsl:when test="$viaTransp = 7">
                                            07 = Rodoviária
                                          </xsl:when>
                                          <xsl:when test="$viaTransp = 8">
                                            08 = Conduto / Rede Transmissão
                                          </xsl:when>
                                          <xsl:when test="$viaTransp = 9">
                                            09 = Meios Próprios
                                          </xsl:when>
                                          <xsl:when test="$viaTransp = 10">
                                            10 = Entrada / Saída ficta
                                          </xsl:when>
                                          <xsl:otherwise>
                                            <xsl:value-of select="$viaTransp"/>
                                          </xsl:otherwise>
                                        </xsl:choose>
                                      </span>
                                    </td>
                                    <td>
                                      <label>Valor da AFRMM</label>
                                      <span>
                                        <xsl:call-template name="format2Casas">
                                          <xsl:with-param name="num" select="n:vAFRMM"/>
                                        </xsl:call-template>
                                      </span>
                                    </td>
                                    <td colspan="3">
                                      <label>Forma de Importação Intermediação</label>
                                      <span>
                                        <xsl:variable name="tpIntermedio" select="n:tpIntermedio"/>
                                        <xsl:choose>
                                          <xsl:when test="$tpIntermedio = 1">
                                            01 = Importação por conta própria
                                          </xsl:when>
                                          <xsl:when test="$tpIntermedio = 2">
                                            02 = Importação por conta e ordem
                                          </xsl:when>
                                          <xsl:when test="$tpIntermedio = 3">
                                            03 = Importação por encomenda
                                          </xsl:when>
                                          <xsl:otherwise>
                                            <xsl:value-of select="$tpIntermedio"/>
                                          </xsl:otherwise>
                                        </xsl:choose> 
                                      </span>
                                    </td>
                                  </tr>
                                  <tr>
                                    <td>
                                      <label>CNPJ Adquirente/Encomendante</label>
                                      <span>
                                        <xsl:call-template name="formatCnpj">
                                          <xsl:with-param name="cnpj" select="n:CNPJ"/>
                                        </xsl:call-template>
                                      </span>
                                    </td>
                                    <td>
                                      <label>Sigla UF Adquirente/Encomendante</label>
                                      <span>
                                        <xsl:value-of select="n:UFTerceiro"/>
                                      </span>
                                    </td>
                                  </tr> 
                                  <tr>
                                    <td colspan="5">
                                      <xsl:for-each select="n:adi">
                                        <div>
                                          <h5 class="toggle">
                                            Adição <xsl:value-of select=" position()"/>
                                          </h5>
                                          <table class="toggable box">
                                            <tr class="col-5">
                                              <td style="width:16%">
                                                <label>Adição</label>
                                                <span>
                                                  <xsl:value-of select="n:nAdicao"/>
                                                </span>
                                              </td>
                                              <td style="width:16%">
                                                <label>Item</label>
                                                <span>
                                                  <xsl:value-of select="n:nSeqAdic"/>
                                                </span>
                                              </td>
                                              <td style="width:28%">
                                                <label>Código Fabricante Estrangeiro</label>
                                                <span>
                                                  <xsl:value-of select="n:cFabricante"/>
                                                </span>
                                              </td>
                                              <td>
                                                <label>Valor do Desconto</label>
                                                <span>
                                                  <xsl:variable name="vDesc" select="n:vDescDI"/>
                                                  <xsl:if test="$vDesc != ''">
                                                    <xsl:value-of select="format-number( $vDesc,'#.###.###.###.##0,00')"/>
                                                  </xsl:if>
                                                </span>
                                              </td>
                                              <td>
                                                <label>Número Drawback</label>
                                                <span>
                                                  <xsl:value-of select="n:nDraw"/>                                                  
                                                </span>
                                              </td>
                                            </tr>
                                          </table>
                                        </div>
                                      </xsl:for-each>
                                    </td>
                                  </tr>
                                </table>
                              </div>
                            </xsl:for-each>
                          </fieldset>
                        </xsl:if>

                        <xsl:variable name="de" select="n:prod/n:detExport"/>
                        <xsl:if test="$de!=''">
                          <fieldset>
                            <legend>Documentos de Exportação</legend>
                            <xsl:for-each select="n:prod/n:detExport">
                              <br/>
                              <div>
                                <h5 class="toggle">
                                  Documento de Exportação nº <xsl:value-of select="position()"/>
                                </h5>
                                <table class="toggable box">
                                  <tr class="col-3">
                                    <td>
                                      <label>Número Drawback</label>
                                      <span>
                                        <xsl:value-of select="n:nDraw"/>
                                      </span>
                                    </td>
                                  </tr> 
                                  <xsl:for-each select="n:exportInd"> 
                                    <tr class="col-3">
                                      <td>
                                        <label>Número Registro Exportação</label>
                                        <span>
                                          <xsl:value-of select="n:nRE"/>
                                        </span>
                                      </td>
                                      <td>
                                        <label>Chave de Acesso Recebida</label>
                                        <span>
                                          <xsl:value-of select="n:chNFe"/>
                                        </span>
                                      </td>
                                      <td>
                                        <label>Quantidade Item</label>
                                        <span>
                                          <xsl:call-template name="format4Casas">
                                            <xsl:with-param name="num" select="n:qExport"/>
                                          </xsl:call-template> 
                                        </span>
                                      </td>
                                    </tr> 
                                  </xsl:for-each> 
                                </table>
                              </div> 
                            </xsl:for-each>
                          </fieldset>
                        </xsl:if> 
                      </td>
                    </tr> 
                    <tr>
                      <xsl:for-each select="n:imposto/n:PIS">
                        <td colspan="12">
                          <fieldset>
                            <legend class="toggle">PIS</legend>
                            <div class="toggable">
                              <xsl:call-template name="PIS"/>
                            </div>
                          </fieldset>
                        </td>
                      </xsl:for-each>
                    </tr>
                    <tr>
                      <xsl:for-each select="n:imposto/n:PISST">
                        <td colspan="12">
                          <fieldset>
                            <legend class="toggle">PISST</legend>
                            <div class="toggable">
                              <xsl:call-template name="PISST"/>
                            </div>
                          </fieldset>
                        </td>
                      </xsl:for-each>
                    </tr>
                    <tr>
                      <xsl:for-each select="n:imposto/n:COFINS">
                        <td colspan="12">
                          <fieldset>
                            <legend class="toggle">COFINS</legend>
                            <div class="toggable">
                              <xsl:call-template name="COFINS"/>
                            </div>
                          </fieldset>
                        </td>
                      </xsl:for-each>
                    </tr>
                    <tr>
                      <xsl:for-each select="n:imposto/n:COFINSST">
                        <td colspan="12">
                          <fieldset>
                            <legend class="toggle">COFINSST</legend>
                            <div class="toggable">
                              <xsl:call-template name="COFINSST"/>
                            </div>
                          </fieldset>
                        </td>
                      </xsl:for-each>
                    </tr>
                    <tr>
                      <xsl:for-each select="n:imposto/n:ISSQN">
                        <td colspan="12">
                          <fieldset>
                            <legend class="toggle">ISSQN</legend>
                            <div class="toggable">
                              <xsl:call-template name="ISSQN"/>
                            </div>
                          </fieldset>
                        </td>
                      </xsl:for-each>
                    </tr>
                    <tr>
                      <xsl:for-each select="n:impostoDevol">
                        <td colspan="12">
                          <fieldset>
                            <legend class="toggle">Imposto Devolvido</legend>
                            <div class="toggable">
                              <xsl:call-template name="TribDevol"/>
                            </div>
                          </fieldset>
                        </td>
                      </xsl:for-each>
                    </tr>
                    <tr>
                      <td colspan="12">
                        <xsl:for-each select="n:prod/n:veicProd">
                          <fieldset>
                            <legend>Detalhamento Específico dos Veículos Novos</legend>
                            <table class="box">
                              <tr class="col-3">
                                <td>
                                  <label>Tipo da Operação</label>
                                  <span>
                                    <xsl:variable name="operacao" select="n:tpOp"/>
                                    <xsl:choose>
                                      <xsl:when test="$operacao = 1">
                                        1 - Venda concessionária
                                      </xsl:when>
                                      <xsl:when test="$operacao = 2">
                                        2 - Faturamento direto para consumidor final
                                      </xsl:when>
                                      <xsl:when test="$operacao = 3">
                                        3 - Venda direta para grandes consumidores
                                      </xsl:when>
                                      <xsl:when test="$operacao = 0">
                                        0 - Outros
                                      </xsl:when>
                                      <xsl:otherwise>
                                         <xsl:value-of select="$operacao"/>
                                      </xsl:otherwise>
                                    </xsl:choose> 
                                  </span>
                                </td>
                                <xsl:for-each select="n:chassi">
                                  <td>
                                    <label>Chassi do veículo</label>
                                    <span>
                                      <xsl:value-of select="text()"/>
                                    </span>
                                  </td>
                                </xsl:for-each>
                                <td>
                                  <label>Cilindradas</label>
                                  <span>
                                    <xsl:value-of select="n:cilin"/>                                   
                                  </span>
                                </td>
                              </tr>
                              <tr>
                                <xsl:for-each select="n:cCor">
                                  <td>
                                    <label>Cor</label>
                                    <span>
                                      <xsl:value-of select="text()"/>
                                    </span>
                                  </td>
                                </xsl:for-each>
                                <xsl:for-each select="n:xCor">
                                  <td>
                                    <label>Descrição da cor</label>
                                    <span>
                                      <xsl:value-of select="text()"/>
                                    </span>
                                  </td>
                                </xsl:for-each>
                                <td>
                                  <label>Código da Cor</label>
                                  <span>
                                    <xsl:variable name="corveiculo" select="n:cCorDENATRAN"/>
                                    <xsl:choose>
                                      <xsl:when test="$corveiculo = 1">
                                        01 - AMARELO
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 2">
                                        02 - AZUL
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 3">
                                        03 - BEGE
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 4">
                                        04 -BRANCA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 5">
                                        05 - CINZA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 6">
                                        06 - DOURADA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 7">
                                        07 - GRENA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 8">
                                        08 - LARANJA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 9">
                                        09 - MARROM
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 10">
                                        10 - PRATA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 11">
                                        11 - PRETA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 12">
                                        12 - ROSA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 13">
                                        13 - ROXA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 14">
                                        14 - VERDE
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 15">
                                        15 - VERMELHA
                                      </xsl:when>
                                      <xsl:when test="$corveiculo = 16">
                                        16 - FANTASIA
                                      </xsl:when>
                                      <xsl:otherwise>
                                        <xsl:value-of select="$corveiculo"/>
                                      </xsl:otherwise>
                                    </xsl:choose> 
                                  </span>
                                </td>
                              </tr>
                              <tr>
                                <td>
                                  <label>Peso Líquido</label>
                                  <span>
                                    <xsl:call-template name="format4Casas">
                                      <xsl:with-param name="num" select="n:pesoL"/>
                                    </xsl:call-template>
                                  </span>
                                </td>
                                <td>
                                  <label>Peso Bruto</label>
                                  <span>
                                    <xsl:call-template name="format4Casas">
                                      <xsl:with-param name="num" select="n:pesoB"/>
                                    </xsl:call-template>
                                  </span>
                                </td>

                                <td>
                                  <label>Serial (Série)</label>
                                  <span>
                                    <xsl:value-of select="n:nSerie"/>
                                  </span>
                                </td>
                              </tr>
                              <tr>
                                <td>
                                  <label>Tipo de Combustível</label>
                                  <span>
                                    <xsl:value-of select="n:tpComb"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Número de Motor</label>
                                  <span>
                                    <xsl:value-of select="n:nMotor"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Capacidade Máxima de Tração</label>
                                  <span>
                                    <xsl:value-of select="n:CMT"/> 
                                  </span>
                                </td>
                              </tr>
                              <tr>
                                <td>
                                  <label>Distância entre eixos</label>
                                  <span>
                                    <xsl:value-of select="n:dist"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Ano Modelo de Fabricação</label>
                                  <span>
                                    <xsl:value-of select="n:anoMod"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Ano de Fabricação</label>
                                  <span>
                                    <xsl:value-of select="n:anoFab"/>
                                  </span>
                                </td>
                              </tr>
                              <tr>
                                <td>
                                  <label>Tipo de Pintura</label>
                                  <span>
                                    <xsl:value-of select="n:tpPint"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Tipo de Veículo</label>
                                  <span>
                                    <xsl:value-of select="n:tpVeic"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Espécie de Veículo</label>
                                  <span>
                                    <xsl:variable name="espVeic" select="n:espVeic"/>
                                    <xsl:choose>
                                      <xsl:when test="$espVeic = 1">
                                        1-PASSAGEIRO
                                      </xsl:when>
                                      <xsl:when test="$espVeic = 2">
                                        2-CARGA
                                      </xsl:when>
                                      <xsl:when test="$espVeic = 3">
                                        3-MISTO
                                      </xsl:when>
                                      <xsl:when test="$espVeic = 4">
                                        4-CORRIDA
                                      </xsl:when>
                                      <xsl:when test="$espVeic = 5">
                                        5-TRAÇÃO
                                      </xsl:when>
                                      <xsl:when test="$espVeic = 6">
                                        6-ESPECIAL
                                      </xsl:when>
                                      <xsl:otherwise>
                                        <xsl:value-of select="$espVeic"/>
                                      </xsl:otherwise>
                                    </xsl:choose> 
                                  </span>
                                </td>
                              </tr>
                              <tr>
                                <td>
                                  <label>
                                    Condição&#160;do&#160;VIN&#160;(<i class="eng">Vehicle&#160;Identification&#160;Number</i>)
                                  </label>
                                  <span>
                                    <xsl:variable name="tpVin" select="n:VIN"/>
                                    <xsl:if test="$tpVin='R'">R-Remarcado</xsl:if>
                                    <xsl:if test="$tpVin='N'">N-Normal</xsl:if>
                                  </span>
                                </td>
                                <td>
                                  <label>Condição do Veículo</label>
                                  <span>
                                    <xsl:variable name="condveiculo" select="n:condVeic"/>
                                    <xsl:choose>
                                      <xsl:when test="$condveiculo='1'">
                                        1-Acabado
                                      </xsl:when>
                                      <xsl:when test="$condveiculo='2'">
                                        2-Inacabado
                                      </xsl:when>
                                      <xsl:when test="$condveiculo='3'">
                                        3-Semi-acabado
                                      </xsl:when>
                                      <xsl:otherwise>
                                        <xsl:value-of select="$condveiculo"/>
                                      </xsl:otherwise>
                                    </xsl:choose> 
                                  </span>
                                </td>
                                <td>
                                  <label>Código Marca Modelo</label>
                                  <span>
                                    <xsl:value-of select="n:cMod"/>
                                  </span>
                                </td>
                              </tr>
                              <tr>
                                <td>
                                  <label>Potência Motor</label>
                                  <span>
                                    <xsl:value-of select="n:pot"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Capacidade Máxima de Lotação</label>
                                  <span>
                                    <xsl:value-of select="n:lota"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Restrição</label>
                                  <span>
                                    <xsl:variable name="restricao" select="n:tpRest"/>
                                    <xsl:choose>
                                      <xsl:when test="$restricao = 0">
                                        0-Não há
                                      </xsl:when>
                                      <xsl:when test="$restricao = 1">
                                        1-Alienação Fiduciária
                                      </xsl:when>
                                      <xsl:when test="$restricao = 2">
                                        2-Arrendamento Mercantil
                                      </xsl:when>
                                      <xsl:when test="$restricao = 3">
                                        3-Reserva de Domínio
                                      </xsl:when>
                                      <xsl:when test="$restricao = 4">
                                        4-Penhor de Veículos
                                      </xsl:when>
                                      <xsl:when test="$restricao = 9">
                                        9-Outras
                                      </xsl:when>
                                      <xsl:otherwise>
                                        <xsl:value-of select="$restricao"/>
                                      </xsl:otherwise>
                                    </xsl:choose> 
                                  </span>
                                </td>
                              </tr>
                            </table>
                          </fieldset>
                        </xsl:for-each>
                        <xsl:variable name="meds" select="n:prod/n:med"/>
                        <xsl:if test="$meds != ''">
                          <fieldset>
                            <legend>Detalhamento específico dos medicamentos</legend>
                            <xsl:for-each select="n:prod/n:med">
                              <div>
                                <h5 class="toggle">Medicamento <xsl:value-of select="position()"/></h5>
                                <table class="toggable box">
                                  <tr class="col-3">
                                    <td>
                                      <label>Nro. do Lote</label>
                                      <span>
                                        <xsl:value-of select="n:nLote"/>
                                      </span>
                                    </td>
                                    <td>
                                      <label>Quantidade de produtos no lote</label>
                                      <span>
                                        <xsl:call-template name="format3Casas">
                                          <xsl:with-param name="num" select ="n:qLote" />
                                        </xsl:call-template>
                                      </span>
                                    </td>
                                    <td>
                                      <label>Data de fabricação</label>
                                      <span>
                                        <xsl:call-template name="formatDate">
                                          <xsl:with-param name="date" select="n:dFab"/>
                                        </xsl:call-template>
                                      </span>
                                    </td>
                                  </tr>
                                  <tr class="col-3">
                                    <td>
                                      <label>Data de validade</label>
                                      <span>
                                        <xsl:call-template name="formatDate">
                                          <xsl:with-param name="date" select="n:dVal"/>
                                        </xsl:call-template>
                                      </span>
                                    </td>
                                    <td>
                                      <label>Preço Máximo Consumidor</label>
                                      <span>
                                        <xsl:call-template name="format2Casas">
                                          <xsl:with-param name="num" select="n:vPMC"/>
                                        </xsl:call-template>
                                      </span>
                                    </td>
                                  </tr>
                                </table>
                              </div>
                            </xsl:for-each>
                          </fieldset>
                        </xsl:if>

                        <xsl:variable name="recopi" select="n:prod/n:nRECOPI"/>
                        <xsl:if test="$recopi != ''">
                          <fieldset>
                            <legend>Número do RECOPI</legend>
                            <table class="toggable box">
                              <tr>
                                <td>
                                  <span>
                                      <xsl:value-of select="$recopi"/>
                                  </span>
                                </td>
                              </tr>
                            </table>
                          </fieldset>
                        </xsl:if>
                        
                        
                        <xsl:variable name="armas" select="n:prod/n:arma"/>
                        <xsl:if test="$armas != ''">
                          <fieldset>
                            <legend>Detalhamento específico dos armamentos</legend>
                            <xsl:for-each select="n:prod/n:arma">
                              <div>
                                <h5 class="toggle">
                                  Armamento <xsl:value-of select="position()"/>
                                </h5>
                                <table class="toggable box">
                                  <tr>
                                    <td>
                                      <label>Tipo de Arma de Fogo</label>
                                      <span>
                                        <xsl:variable name="tparma" select="n:tpArma"/>
                                        <xsl:choose>
                                          <xsl:when test="$tparma = 0">
                                            0 - Uso permitido
                                          </xsl:when>
                                          <xsl:when test="$tparma = 1">
                                            1 - Uso restrito
                                          </xsl:when>
                                          <xsl:otherwise>
                                            <xsl:value-of select="$tparma"/>
                                          </xsl:otherwise>
                                        </xsl:choose> 
                                      </span>
                                    </td>
                                    <td>
                                      <label>Número de Série da Arma</label>
                                      <span>
                                        <xsl:value-of select="n:nSerie"/>
                                      </span>
                                    </td>
                                    <td>
                                      <label>Número de Série do Cano</label>
                                      <span>
                                        <xsl:value-of select="n:nCano"/>
                                      </span>
                                    </td>
                                  </tr>
                                  <tr>
                                    <td colspan="3">
                                      <label>Descrição</label>
                                      <span>
                                        <xsl:value-of select="n:descr"/>
                                      </span>
                                    </td>
                                  </tr>
                                </table>
                              </div>
                            </xsl:for-each>
                          </fieldset>
                        </xsl:if>
                        <xsl:for-each select="n:prod/n:comb">
                          <fieldset>
                            <legend>Detalhamento específico de combustível</legend>
                            <table class="box">
                              <tr class="col-3">
                                <td>
                                  <label>Código do Produto da ANP</label>
                                  <span>
                                    <xsl:value-of select="n:cProdANP"/>
                                  </span>
                                </td>
                                <td>
                                  <label>Percentual Gás Natural</label>
                                  <span>
                                    <xsl:call-template name="format4Casas">
                                      <xsl:with-param name="num" select="n:pMixGN"/>
                                    </xsl:call-template> 
                                  </span>
                                </td>
                                <td>
                                  <label>CODIF</label>
                                  <span>
                                    <xsl:value-of select="n:CODIF"/>
                                  </span>
                                </td> 
                              </tr>
                              <tr class="col-3">
                                <td>
                                  <label>Quantidade&#160;Combustível&#160;Faturada</label>
                                  <span>
                                    <xsl:call-template name="format4Casas">
                                      <xsl:with-param name="num" select="n:qTemp"/>
                                    </xsl:call-template>
                                  </span>
                                </td>
                                <td>
                                  <label>UF de Consumo</label>
                                  <span>
                                    <xsl:value-of select="n:UFCons"/>
                                  </span>
                                </td>
                              </tr>

                              <xsl:if test="count(n:CIDE)=1">
                                <tr class="col-4">
                                  <td colspan="4">
                                    <fieldset>
                                      <legend>CIDE</legend>
                                      <table  class="box">
                                        <tr class="col-3">
                                          <td>
                                            <label>Quant. Base de Cálculo</label>
                                            <span>
                                              <xsl:call-template name="format4Casas">
                                                <xsl:with-param name="num" select="n:CIDE/n:qBCProd"/>
                                              </xsl:call-template>
                                            </span>
                                          </td>
                                          <td>
                                            <label>Valor da Alíquota (R$)</label>
                                            <span>
                                              <xsl:call-template name="format4Casas">
                                                <xsl:with-param name="num" select="n:CIDE/n:vAliqProd"/>
                                              </xsl:call-template>
                                            </span>
                                          </td>
                                          <td>
                                            <label>Valor</label>
                                            <span>
                                              <xsl:call-template name="format2Casas">
                                                <xsl:with-param name="num" select="n:CIDE/n:vCIDE"/>
                                              </xsl:call-template>
                                            </span>
                                          </td>
                                        </tr>
                                      </table>
                                    </fieldset>
                                  </td>
                                </tr>
                              </xsl:if>
                            </table>
                          </fieldset>
                        </xsl:for-each>
                        <xsl:for-each select="n:infAdProd">
                          <fieldset>
                            <legend class="titulo-aba-interna">Informações adicionais do produto</legend>
                            <table class="box">
                              <tr>
                                <td>
                                  <label>Descrição</label>
                                  <span>
                                    <xsl:value-of select="text()"/>
                                  </span>
                                </td>
                              </tr>
                            </table>
                          </fieldset>
                        </xsl:for-each>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
          </xsl:for-each>
        </div>
      </fieldset>
    </div>
  </xsl:template>
</xsl:stylesheet>
