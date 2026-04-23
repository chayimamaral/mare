<?xml version="1.0" encoding="utf-8"?>

<xsl:stylesheet version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:n="http://www.portalfiscal.inf.br/nfe"
     xmlns:date="http://exslt.org/formatoData"
      xmlns:chave="http://exslt.org/chaveacesso"
    xmlns:r="http://www.serpro.gov.br/nfe/remessanfe.xsd" >
  <xsl:decimal-format decimal-separator="," grouping-separator="." />
  <xsl:template match="/">

    <script language = "JavaScript">
      function ocultarExibir(idDaTabela)
      {
      if (document.getElementById(idDaTabela).style.display == 'inline')
      {
      document.getElementById(idDaTabela).style.display = 'none';
      document.getElementById("label_exibir_" + idDaTabela).style.display = 'inline';
      document.getElementById("label_ocultar_" + idDaTabela).style.display = 'none';
      }
      else
      {
      document.getElementById(idDaTabela).style.display = 'inline';
      document.getElementById("label_exibir_" + idDaTabela).style.display = 'none';
      document.getElementById("label_ocultar_" + idDaTabela).style.display = 'inline';
      }
      }
    </script>

    <table align="center"  width="98%">
      <tr>
        <td height="25" class="TituloAreaRestritacentro">
          <B class="textoVerdana8bold">
            <STRONG class="textoVerdana8bold">Informações Adicionais</STRONG>
          </B>
        </td>
      </tr>
    </table>

    <table align="center"  width="98%">
      <xsl:for-each select ="//n:infNFe/n:ide/n:tpImp">
        <td valign ="top" width="33%"  >
          <span class="TextoFundoBrancoNegrito">
            Formato de Impressão<br />
          </span>
          <span class="linha">
            <xsl:variable name="forimp" select="text()"/>
            <xsl:if test="$forimp='1'">
              1 - Retrato
            </xsl:if>
            <xsl:if test="$forimp='2'">
              2 - Paisagem
            </xsl:if>
          </span>
        </td>
      </xsl:for-each>
      <xsl:for-each select ="//n:infNFe/n:ide/n:tpEmis">
        <td valign ="top" width="33%"  >
          <span class="TextoFundoBrancoNegrito">
            Formato de Emissão<br />
          </span>
          <span class="linha">
            <xsl:value-of select="text()"/>
          </span>
        </td>
      </xsl:for-each>
      <xsl:for-each select ="//n:infNFe/n:ide/n:cDV">
        <td valign ="top" width="33%"  >
          <span class="TextoFundoBrancoNegrito">
            Dígito Verificador da Chave de Acesso<br />
          </span>
          <span class="linha">
            <xsl:value-of select="text()" />
          </span>
        </td>
      </xsl:for-each>

    </table>

    <table align="center"  width="98%">
      <xsl:for-each select ="//n:infNFe/n:ide/n:tpAmb">
        <td valign ="top" width="33%"  >
          <span class="TextoFundoBrancoNegrito">
            Identificação do Ambiente<br />
          </span>
          <span class="linha">
            <xsl:variable name="ambiente" select="text()"/>
            <xsl:if test="$ambiente='1'">
              1 – Produção
            </xsl:if>
            <xsl:if test="$ambiente='2'">
              2 - Homologação
            </xsl:if>
          </span>
        </td>
      </xsl:for-each>
      <xsl:for-each select ="//n:infNFe/n:ide/n:finNFe">
        <td valign ="top" width="33%"  >
          <span class="TextoFundoBrancoNegrito">
            Finalidade<br />
          </span>
          <span class="linha">
            <xsl:value-of select="text()" />
          </span>
        </td>
      </xsl:for-each>
      <xsl:for-each select ="//n:infNFe/n:ide/n:procEmi">
        <td valign ="top" width="33%"  >
          <span class="TextoFundoBrancoNegrito">
            Processo<br />
          </span>
          <span class="linha">
            <xsl:value-of select="text()" />
          </span>
        </td>
      </xsl:for-each>

    </table>

    <table align="center"  width="98%">
      <xsl:for-each select ="//n:infNFe/n:ide/n:verProc">
        <td valign ="top" width="33%"  >
          <span class="TextoFundoBrancoNegrito">
            Versão<br />
          </span>
          <span class="linha">
            <xsl:value-of select="text()" />
          </span>
        </td>
      </xsl:for-each>

      <xsl:for-each select ="//n:infNFe/n:ide/n:dhCont">
        <td valign ="top" width="33%"  >
          <span class="TextoFundoBrancoNegrito">
            Entrada em Contingência<br />
          </span>
          <span class="linha">
            <xsl:value-of select="text()" />
            <xsl:value-of select="date:formatdate(text(),'dd/MM/yyyy HH:mm:ss')" />
          </span>
        </td>
      </xsl:for-each>

      <td valign ="top" width="33%" >
      </td>
    </table >

    <table align="center"  width="98%">
      <xsl:for-each select ="//n:infNFe/n:ide/n:xJust">
        <td valign ="top" width="100%"  >
          <span class="TextoFundoBrancoNegrito">
            Justificativa<br />
          </span>
          <span class="linha">
            <xsl:value-of select="text()" />
          </span>
        </td>
      </xsl:for-each>
      <td valign ="top" width="33%"  >

      </td>
      <td valign ="top" width="33%" >
      </td>
    </table>


    <xsl:variable name="exporta" select="//n:NFe/n:infNFe/n:exporta" />
    <xsl:if test="$exporta!=''">

      <table align="center"  width="98%">
        <tr class="TextoFundoBranco">
          <td  class="TituloAreaRestrita2">
            EXPORTAÇÃO
          </td>
        </tr>
      </table>

      <table align="center"  width="98%">
        <xsl:for-each select ="//n:NFe/n:infNFe/n:exporta/n:xLocEmbarq">
          <td valign ="top" width="50%"  >
            <span class="TextoFundoBrancoNegrito">
              Local de Embarque<br />
            </span>
            <span class="linha">
              <xsl:value-of select="text()" />
            </span>
          </td>
        </xsl:for-each>
        <xsl:for-each select ="//n:NFe/n:infNFe/n:exporta/n:UFEmbarq">
          <td valign ="top" width="50%"  >
            <span class="TextoFundoBrancoNegrito">
              UF de Embarque<br />
            </span>
            <span class="linha">
              <xsl:value-of select="text()" />
            </span>
          </td>
        </xsl:for-each>
        <tr></tr>
      </table>
    </xsl:if>

    <xsl:variable name="compra" select="//n:NFe/n:infNFe/n:compra" />
    <xsl:if test="$compra!=''">

      <table align="center"  width="98%">
        <tr class="TextoFundoBranco">
          <td  class="TituloAreaRestrita2">
            INFORMAÇÕES DE COMPRA
          </td>
        </tr>
      </table>

      <table align="center"  width="98%">
        <xsl:for-each select ="//n:NFe/n:infNFe/n:compra/n:xNEmp">
          <td valign ="top" width="33%"  >
            <span class="TextoFundoBrancoNegrito">
              Nota de Empenho<br />
            </span>
            <span class="linha">
              <xsl:value-of select="text()" />
            </span>
          </td>
        </xsl:for-each>
        <xsl:for-each select ="//n:NFe/n:infNFe/n:compra/n:xPed">
          <td valign ="top" width="33%"  >
            <span class="TextoFundoBrancoNegrito">
              Pedido<br />
            </span>
            <span class="linha">
              <xsl:value-of select="text()" />
            </span>
          </td>
        </xsl:for-each>
        <xsl:for-each select ="//n:NFe/n:infNFe/n:compra/n:xCont">
          <td valign ="top" width="33%"  >
            <span class="TextoFundoBrancoNegrito">
              Contrato<br />
            </span>
            <span class="linha">
              <xsl:value-of select="text()" />
            </span>
          </td>
        </xsl:for-each>
        <tr></tr>
      </table>
    </xsl:if>

    <!-- CANA-->
    <xsl:variable name="cana" select="//n:NFe/n:infNFe/n:cana" />
    <xsl:if test="$cana!=''">

      <table align="center"  width="98%">
        <tr class="TextoFundoBranco">
          <td  class="TituloAreaRestrita2">
            INFORMAÇÕES DO REGISTRO DE AQUISIÇÃO DE CANA
          </td>
        </tr>
      </table>

      <table align="center"  class="textoVerdana7" width="98%">
        <tr>
          <td valign ="top"  width="50%" >
            <span class="TextoFundoBrancoNegrito">
              Identificação da safra
            </span>
          </td>
          <td valign ="top"  width="50%" >
            <span class="TextoFundoBrancoNegrito">
              Mês e ano de referência
            </span>
          </td>
        </tr>
        <xsl:for-each select ="//n:NFe/n:infNFe/n:cana">
          <tr>
            <xsl:for-each select ="n:safra" >
              <td valign ="top"  width="50%"  >
                <span class="linha">
                  <xsl:value-of select="text()" />
                </span>
              </td>
            </xsl:for-each>
            <xsl:for-each select ="n:ref" >
              <td valign ="top"  width="50%"  >
                <span class="linha">
                  <xsl:value-of select="text()" />
                </span>
              </td>
            </xsl:for-each>
          </tr>
        </xsl:for-each>

        <tr>
          <table align="center"  width="98%">
            <tr>
              <xsl:for-each select ="//n:NFe/n:infNFe/n:cana/n:qTotMes">
                <td valign ="top" width="33%"  >
                  <span class="TextoFundoBrancoNegrito">
                    Quantidade Total do Mês<br />
                  </span>
                  <span class="linha">
                    <xsl:value-of select="format-number(text(),'##.##.##0,0000000000')" />
                  </span>
                </td>
              </xsl:for-each>
              <xsl:for-each select ="//n:NFe/n:infNFe/n:cana/n:qTotAnt">
                <td valign ="top" width="33%"  >
                  <span class="TextoFundoBrancoNegrito">
                    Quantidade Total Anterior<br />
                  </span>
                  <span class="linha">
                    <xsl:value-of select="format-number(text(),'##.##.##0,0000000000')" />
                  </span>
                </td>
              </xsl:for-each>
              <xsl:for-each select ="//n:NFe/n:infNFe/n:cana/n:qTotGer">
                <td valign ="top" width="33%"  >
                  <span class="TextoFundoBrancoNegrito">
                    Quantidade Total Geral<br />
                  </span>
                  <span class="linha">
                    <xsl:value-of select="format-number(text(),'##.##.##0,0000000000')" />
                  </span>
                </td>
              </xsl:for-each>
            </tr>
            <tr></tr>
            <tr>
              <xsl:for-each select ="//n:NFe/n:infNFe/n:cana/n:vFor">
                <td valign ="top" width="33%"  >
                  <span class="TextoFundoBrancoNegrito">
                    Valor dos Fornecimentos<br />
                  </span>
                  <span class="linha">
                    <xsl:value-of select="format-number(text(),'##.##.##0,00')" />
                  </span>
                </td>
              </xsl:for-each>
              <xsl:for-each select ="//n:NFe/n:infNFe/n:cana/n:vTotDed">
                <td valign ="top" width="33%"  >
                  <span class="TextoFundoBrancoNegrito">
                    Valor Total da Dedução<br />
                  </span>
                  <span class="linha">
                    <xsl:value-of select="format-number(text(),'##.##.##0,00')" />
                  </span>
                </td>
              </xsl:for-each>
              <xsl:for-each select ="//n:NFe/n:infNFe/n:cana/n:vLiqFor">
                <td valign ="top" width="33%"  >
                  <span class="TextoFundoBrancoNegrito">
                    Valor Líquido dos Fornecimentos<br />
                  </span>
                  <span class="linha">
                    <xsl:value-of select="format-number(text(),'##.##.##0,00')" />
                  </span>
                </td>
              </xsl:for-each>
            </tr>
          </table>
        </tr>





        <tr>
          <xsl:variable name="deduc" select="//n:NFe/n:infNFe/n:cana/n:deduc" />
          <xsl:if test="$deduc!=''">
            <td colspan="6" align="right">
              <table width="95%">
                <tr>
                  <td>
                    <table align="center"  width="98%">
                      <tr class="TextoFundoBranco">
                        <td  class="TituloAreaRestrita2">
                          DEDUÇÕES – TAXAS E CONTRIBUIÇÕES
                        </td>
                      </tr>
                    </table>
                    <table align="center"  class="textoVerdana7" width="98%">
                      <tr>
                        <td valign ="top"  width="50%" align="left" >
                          <span class="TextoFundoBrancoNegrito">
                            Descrição da Dedução
                          </span>
                        </td>
                        <td valign ="top"  width="50%" align="left" >
                          <span class="TextoFundoBrancoNegrito">
                            Valor da Dedução
                          </span>
                        </td>
                      </tr>
                      <xsl:for-each select="//n:NFe/n:infNFe/n:cana/n:deduc" >
                        <tr>
                          <xsl:for-each select ="n:xDed" >
                            <td valign ="top"  width="50%" align="left" >
                              <span class="linha">
                                <xsl:value-of select="text()" />
                              </span>
                            </td>
                          </xsl:for-each>
                          <xsl:for-each select ="n:vDed" >
                            <td valign ="top"  width="50%" align="left" >
                              <span class="linha">
                                <xsl:value-of select="format-number(text(),'##.##.##0,00')" />
                              </span>
                            </td>
                          </xsl:for-each>
                        </tr>
                      </xsl:for-each>
                      <tr></tr>
                    </table>

                  </td >
                </tr >
              </table >
            </td >
          </xsl:if>
        </tr>

        <tr>
          <xsl:variable name="fordia" select="//n:NFe/n:infNFe/n:cana/n:forDia" />
          <xsl:if test="$fordia!=''">
            <td colspan="6" align="right">
              <table width="95%">
                <tr>
                  <td>
                    <table align="center"  class="textoVerdana7" width="98%">
                      <tr class="TextoFundoBranco">
                        <td  class="TituloAreaRestrita2">
                          FORNECIMENTO DIÁRIO DE CANA
                        </td>
                      </tr>
                    </table>
                    <table align="center"  class="textoVerdana7" width="98%">
                      <tr>
                        <td valign ="top" align="left"  width="50%" >
                          <span class="TextoFundoBrancoNegrito">
                            Dia
                          </span>
                        </td>
                        <td valign ="top"  align="left" width="50%" >
                          <span class="TextoFundoBrancoNegrito">
                            Quantidade
                          </span>
                        </td>
                      </tr>
                      <xsl:for-each select="//n:NFe/n:infNFe/n:cana/n:forDia" >
                        <tr>
                          <td valign ="top" align="left"  width="50%"  >
                            <span class="linha">
                              <xsl:value-of select="@dia" />
                            </span>
                          </td>
                          <xsl:for-each select ="n:qtde" >
                            <td valign ="top" align="left" width="50%"  >
                              <span class="linha">
                                <xsl:value-of select="format-number(text(),'##.##.##0,0000000000')" />
                              </span>
                            </td>
                          </xsl:for-each>
                        </tr>
                      </xsl:for-each>
                      <tr></tr>
                    </table>
                  </td>
                </tr>
              </table >
            </td >
          </xsl:if>
        </tr>



      </table >
    </xsl:if>

    <!--CABEÇALHO DA ABA (O MESMO PARA TODAS AS ABAS EXCETO A ABA NFe, A QUAL INFORMA O NÚMERO DA NFe NOS DADOS DA NFe)-->
    <xsl:variable name = "procreferenciado" select = "/n:NFe/n:infNFe/n:infAdic/n:procRef"/>
    <xsl:if test = "$procreferenciado != ''">

      <table align="center"  width="98%">
        <tr class="TextoFundoBranco">
          <td  class="TituloAreaRestrita2">
            PROCESSO REFERENCIADO
          </td>
        </tr>
      </table>
      <xsl:for-each select = "//n:infAdic/n:procRef">
        <xsl:variable name = "chaves" select = "position()"/>
        <xsl:variable name = "chaves2" select = "concat('A',' + ',$chaves)"/>
        <table align = "center" class = "textoVerdana7" width = "98%" id = "{$chaves2}">
          <tr>
            <td valign = "top" style = "height: 20px; width: 100%;">
              <span class="TextoFundoBrancoNegrito">
                Processo Referenciado <xsl:value-of select = " position()"/>
              </span>
              <br/>
            </td>
          </tr>
        </table>

        <xsl:variable name = "tab" select = "position()"/>
        <table class = "tabelaInterna" align = "center" style = "width: 98%;">
          <tr>
            <td width="2%"></td>
            <td valign ="top" width="50%"  >
              <span class="TextoFundoBrancoNegrito">
                Identificador/Ato concessório<br />
              </span>
              <span class="linha">
                <xsl:value-of select = "n:nProc"/>
              </span>
            </td>
            <td valign ="top" width="50%"  >
              <span class="TextoFundoBrancoNegrito">
                Indicador da Origem
              </span>
              <span class = "linha">
                <xsl:variable name = "vindProc" select = "n:indProc"/>
                <xsl:if test="$vindProc='0'">
                  0 - SEFAZ
                </xsl:if>
                <xsl:if test="$vindProc='1'">
                  1 - Justiça Federal
                </xsl:if>
                <xsl:if test="$vindProc='2'">
                  2 - Justiça Estadual
                </xsl:if>
                <xsl:if test="$vindProc='3'">
                  3 - Secex/RFB
                </xsl:if>
                <xsl:if test="$vindProc='9'">
                  9 - Outros
                </xsl:if>
              </span>
            </td> 
          </tr>
        </table>
      </xsl:for-each>
    </xsl:if > 
    <xsl:variable name="fisco" select="//n:infNFe/n:infAdic/n:infAdFisco" />
    <xsl:if test="$fisco!=''">
      <table align="center"  width="98%">
        <tr class="TextoFundoBranco">
          <td  class="TituloAreaRestrita2">
            INFORMAÇÕES ADICIONAIS DE INTERESSE DO FISCO
          </td>
        </tr>
      </table>
      <table align="center"  class="textoVerdana7" width="98%">
        <tr>
          <td valign ="top"  >
            <span class="TextoFundoBrancoNegrito">
              Descrição
            </span>
            <span class="linha">
              <xsl:value-of select="$fisco" />
            </span>
          </td>
        </tr>
        <tr></tr>
      </table>
    </xsl:if>

    <xsl:variable name="infofisco" select="//n:infNFe/n:infAdic/n:obsFisco" />
    <xsl:if test="$infofisco!=''">
      <table align="center"  width="98%">
        <tr class="TextoFundoBranco">
          <td  class="TituloAreaRestrita2">
            OBSERVAÇÕES DO FISCO
          </td>
        </tr>
      </table>
      <table align="center"  class="textoVerdana7" width="98%">
        <tr>
          <td valign ="top"  width="50%"  >
            <span class="TextoFundoBrancoNegrito">
              Campo
            </span>
          </td>
          <td valign ="top"  width="50%"  >
            <span class="TextoFundoBrancoNegrito">
              Texto
            </span>
          </td>
        </tr>
        <xsl:for-each select="//n:infNFe/n:infAdic/n:obsFisco" >
          <tr>
            <td valign ="top"  width="50%"  >
              <span  class="linha">
                <xsl:value-of select ="@xCampo" />
              </span>
            </td>
            <xsl:for-each select ="n:xTexto" >
              <td valign ="top"  width="50%"  >
                <span class="linha">
                  <xsl:value-of select="text()" />
                </span>
              </td>
            </xsl:for-each>
          </tr>
        </xsl:for-each>
        <tr></tr>
      </table>
    </xsl:if>

    <xsl:variable name="contribuinte" select="//n:infNFe/n:infAdic/n:infCpl" />
    <xsl:if test="$contribuinte!=''">
      <table align="center"  width="98%">
        <tr class="TextoFundoBranco">
          <td  class="TituloAreaRestrita2">
            INFORMAÇÕES COMPLEMENTARES DE INTERESSE DO CONTRIBUINTE
          </td>
        </tr>
      </table>
      <table align="center"  class="textoVerdana7" width="98%">
        <td valign ="top"  >
          <span class="TextoFundoBrancoNegrito">
            Descrição
          </span>
          <span class="linha">
            <xsl:value-of select="$contribuinte" />
          </span>
        </td>
        <tr></tr>
      </table>
    </xsl:if>

    <xsl:variable name="infocont" select="//n:infNFe/n:infAdic/n:obsCont" />
    <xsl:if test="$infocont!=''">
      <table align="center"  width="98%">
        <tr class="TextoFundoBranco">
          <td  class="TituloAreaRestrita2">
            OBSERVAÇÕES DO CONTRIBUINTE
          </td>
        </tr>
      </table>
      <table align="center"  class="textoVerdana7" width="98%">
        <tr>
          <td valign ="top"  width="50%" >
            <span class="TextoFundoBrancoNegrito">
              Campo
            </span>
          </td>
          <td valign ="top"  width="50%" >
            <span class="TextoFundoBrancoNegrito">
              Texto
            </span>
          </td>
        </tr>
        <xsl:for-each select="//n:infNFe/n:infAdic/n:obsCont" >
          <tr>
            <td valign ="top"  width="50%"  >
              <span  class="linha">
                <xsl:value-of select ="@xCampo" />
              </span>
            </td>
            <xsl:for-each select ="n:xTexto" >
              <td valign ="top"  width="50%"  >
                <span class="linha">
                  <xsl:value-of select="text()" />
                </span>
              </td>
            </xsl:for-each>
          </tr>
        </xsl:for-each>
        <tr></tr>
      </table>
    </xsl:if>

    <!-- Nota Fiscal Referenciada-->
    <xsl:variable name="nferef" select="//n:infNFe/n:ide/n:NFref" />
    <xsl:if test="$nferef!=''">

      <!-- <xsl:variable name="ideuf" select="//n:infNFe/n:ide/n:NFref/n:refNFe/n:cUF" />
      <xsl:variable name="ideaamm" select="//n:infNFe/n:ide/n:NFref/n:refNFe/n:AAMM" />
      <xsl:variable name="idecnpj" select="//n:infNFe/n:ide/n:NFref/n:refNFe/n:CNPJ" />
      <xsl:variable name="idemod" select="//n:infNFe/n:ide/n:NFref/n:refNFe/n:mod" />
      <xsl:variable name="ideserie" select="//n:infNFe/n:ide/n:NFref/n:refNFe/n:serie" />
      <xsl:variable name="idennf" select="//n:infNFe/n:ide/n:NFref/n:refNFe/n:nNF" />
      <xsl:variable name="idetpimp" select="//n:infNFe/n:ide/n:tpImp" />
      <xsl:variable name="idetpemis" select="//n:infNFe/n:ide/n:tpEmis" />
      <xsl:variable name="idecdv" select="//n:infNFe/n:ide/n:cDV" />
      <xsl:variable name="idetpamb" select="//n:infNFe/n:ide/n:tpAmb" />
      <xsl:variable name="idefinnfe" select="//n:infNFe/n:ide/n:finNFe" />
      <xsl:variable name="ideprocemi" select="//n:infNFe/n:ide/n:procEmi" />
      <xsl:variable name="ideverproc" select="//n:infNFe/n:ide/n:verProc" />
	  or($idetpimp)or($idetpemis)or($idecdv)or($idetpamb)or($idefinnfe)or($ideprocemi)or($ideverproc) 
      <xsl:if test="($ideuf)or($ideaamm)or($idecnpj)or($idemod)or($ideserie)or($idennf)!='' ">-->

      <table align="center"  width="98%">
        <tr class="TextoFundoBranco">
          <td  class="TituloAreaRestrita2">
            DOCUMENTOS FISCAIS REFERENCIADOS
          </td>
        </tr>

        <TR>
          <td colspan="6" align="right">
            <table width="95%">
              <tr>
                <td>

                  <table align="center"  width="98%">
                    <tr class="TextoFundoBranco">
                      <td  class="TituloAreaRestrita2">
                        NOTA FISCAL
                      </td>
                    </tr>
                  </table>

                  <table align="center"  width="98%">
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNF/n:cUF">
                      <td valign ="top" width="33%"  align="left">
                        <span class="TextoFundoBrancoNegrito">
                          Código da UF<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select = "text()"/>
                        </span>
                      </td>
                    </xsl:for-each>
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNF/n:AAMM">
                      <td valign ="top" width="33%" align="left" >
                        <span class="TextoFundoBrancoNegrito">
                          Ano / Mês<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select="text()" />
                        </span>
                      </td>
                    </xsl:for-each>
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNF/n:CNPJ">
                      <td valign ="top" width="33%"  align="left">
                        <span class="TextoFundoBrancoNegrito">
                          CNPJ<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select="chave:formatarCnpj(text())" />
                        </span>
                      </td>
                    </xsl:for-each>

                  </table>

                  <table align="center"  width="98%">
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNF/n:mod">
                      <td valign ="top" width="33%" align="left" >
                        <span class="TextoFundoBrancoNegrito">
                          Modelo<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select="text()" />
                        </span>
                      </td>
                    </xsl:for-each>
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNF/n:serie">
                      <td valign ="top" width="33%"  align="left">
                        <span class="TextoFundoBrancoNegrito">
                          Série<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select="text()" />
                        </span>
                      </td>
                    </xsl:for-each>
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNF/n:nNF">
                      <td valign ="top" width="33%"  align="left">
                        <span class="TextoFundoBrancoNegrito">
                          Número<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select="text()" />
                        </span>
                      </td>
                    </xsl:for-each>
                  </table>

                  <table align="center"  width="98%">
                    <tr class="TextoFundoBranco">
                      <td  class="TituloAreaRestrita2">
                        NOTA FISCAL ELETRÔNICA
                      </td>
                    </tr>
                  </table>

                  <table align="center"  width="98%">
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNFe">
                      <td valign ="top" width="100%" align="left" >
                        <span class="TextoFundoBrancoNegrito">
                          Chave de Acesso<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select="chave:formatarChaveAcesso(text())" />
                        </span>
                      </td>
                    </xsl:for-each>
                  </table>

                  <xsl:variable name="refnfep" select="//n:infNFe/n:ide/n:NFref/n:refNFP" />
                  <xsl:if test="$refnfep!=''">
                    <table align="center"  width="98%">
                      <tr class="TextoFundoBranco">
                        <td  class="TituloAreaRestrita2">
                          NOTA FISCAL DE PRODUTOR RURAL
                        </td>
                      </tr>
                    </table>
                    <table align="center"  width="98%">
                      <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNFP/n:cUF">
                        <td valign ="top" width="25%" align="left" >
                          <span class="TextoFundoBrancoNegrito">
                            Código da UF<br />
                          </span>
                          <span class="linha">
                            <xsl:value-of select = "text()"/>
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNFP/n:AAMM">
                        <td valign ="top" width="25%" align="left" >
                          <span class="TextoFundoBrancoNegrito">
                            Ano / Mês<br />
                          </span>
                          <span class="linha">
                            <xsl:value-of select="text()" />
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNFP/n:CNPJ">
                        <td valign ="top" width="25%"  align="left">
                          <span class="TextoFundoBrancoNegrito">
                            CNPJ<br />
                          </span>
                          <span class="linha">
                            <xsl:value-of select="chave:formatarCnpj(text())" />
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNFP/n:CPF">
                        <td valign ="top" width="25%" align="left" >
                          <span class="TextoFundoBrancoNegrito">
                            CPF<br />
                          </span>
                          <span class="linha">
                            <xsl:value-of select="chave:formatarCPF(text())" />
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNFP/n:IE">
                        <td valign ="top" width="25%" align="left" >
                          <span class="TextoFundoBrancoNegrito">
                            IE<br />
                          </span>
                          <span class="linha">
                            <xsl:value-of select = "text()"/>
                          </span>
                        </td>
                      </xsl:for-each>
                    </table>

                    <table align="center"  width="98%">
                      <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNFP/n:mod">
                        <td valign ="top" width="33%"  align="left">
                          <span class="TextoFundoBrancoNegrito">
                            Modelo do Documento Fiscal<br />
                          </span>
                          <span class="linha">
                            <xsl:variable name="mod" select="text()"/>
                            <xsl:if test="$mod='04'">
                              04 – NF de Produtor
                            </xsl:if>
                            <xsl:if test="$mod='01'">
                              01- para NF avulsa
                            </xsl:if>
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNFP/n:serie">
                        <td valign ="top" width="33%"  align="left">
                          <span class="TextoFundoBrancoNegrito">
                            Série do Documento Fiscal<br />
                          </span>
                          <span class="linha">
                            <xsl:value-of select="text()" />
                          </span>
                        </td>
                      </xsl:for-each>
                      <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refNFP/n:nNF">
                        <td valign ="top" width="33%" align="left">
                          <span class="TextoFundoBrancoNegrito">
                            Número do Documento Fiscal<br />
                          </span>
                          <span class="linha">
                            <xsl:value-of select = "text()"/>
                          </span>
                        </td>
                      </xsl:for-each>
                    </table>

                  </xsl:if>


                  <table align="center"  width="98%">
                    <tr class="TextoFundoBranco">
                      <td  class="TituloAreaRestrita2">
                        INFORMAÇÕES DO CUPOM FISCAL
                      </td>
                    </tr>
                  </table>
                  <table align="center"  width="98%">
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refECF/n:mod">
                      <td valign ="top" width="33%"  align="left">
                        <span class="TextoFundoBrancoNegrito">
                          Modelo do Documento Fiscal<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select = "text()"/>
                        </span>
                      </td>
                    </xsl:for-each>
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refECF/n:nECF">
                      <td valign ="top" width="33%"  align="left">
                        <span class="TextoFundoBrancoNegrito">
                          Número de Ordem Seqüencial do ECF<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select = "text()"/>
                        </span>
                      </td>
                    </xsl:for-each>
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refECF/n:nCOO">
                      <td valign ="top" width="33%"  align="left">
                        <span class="TextoFundoBrancoNegrito">
                          Número do Contador de Ordem de Operação<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select="text()" />
                        </span>
                      </td>
                    </xsl:for-each>
                  </table>

                  <table align="center"  width="98%">
                    <tr class="TextoFundoBranco">
                      <td  class="TituloAreaRestrita2">
                        CONHECIMENTO DE TRANSPORTE ELETRÔNICO
                      </td>
                    </tr>
                  </table>

                  <table align="center"  width="98%">
                    <xsl:for-each select ="//n:infNFe/n:ide/n:NFref/n:refCTe">
                      <td valign ="top" width="50%"  align="left">
                        <span class="TextoFundoBrancoNegrito">
                          Chave de Acesso<br />
                        </span>
                        <span class="linha">
                          <xsl:value-of select="chave:formatarChaveAcesso(text())" />
                        </span>
                      </td>
                    </xsl:for-each>
                  </table>





                </td>
              </tr>
            </table>
          </td>
        </TR>
      </table >

    </xsl:if>


  </xsl:template>
</xsl:stylesheet>

