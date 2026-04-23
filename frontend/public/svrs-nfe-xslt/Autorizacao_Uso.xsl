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
  <xsl:template match="/" name="AutUso">
    <table align="center" width="740" border="0" cellSpacing="0" cellPadding="0">
      <tr>
        <td valign="top">
          <table width="560" align="center" borderColor="#f2f2f2" border="2" cellSpacing="0" cellPadding="0">
            <tr>
              <td align="center" valign="top">
                <table width="90%" border="0" cellSpacing="0" cellPadding="0" >
                  <tr>
                    <td height="120" align="center" valign="middle">
                      <xsl:if test="substring(//n:NFe/n:infNFe/@Id,4,2)=43">
                        <img width="51" height="65" src="imagens/brasao.jpg" />
                      </xsl:if>
                      <xsl:if test="substring(//n:NFe/n:infNFe/@Id,4,2)!=43">
                        <img src="imagens/logo_nfe_p.gif" />
                      </xsl:if>
                    </td>
                    <td height="120" align="center" valign="middle">
                      <strong>
                        GOVERNO DO ESTADO -
                        <xsl:call-template name="nomeUF">
                          <xsl:with-param name="uf" select="substring(//n:NFe/n:infNFe/@Id,4,2)"/>
                        </xsl:call-template>
                        <br></br>
                        SECRETARIA DA FAZENDA<br></br>
                        <br></br>
                        Nota Fiscal Eletrônica - Autorização de Uso<br></br>
                      </strong>

                    </td>
                  </tr>
                </table>
                <table width="570" border="0" cellSpacing="0" cellPadding="5">
                  <tr>
                    <td width="140" bgColor="#ffffff" colspan="2">
                      <strong>Chave de Acesso</strong>
                    </td>
                    <td bgColor="#ffffff" colspan="6">
                      <xsl:call-template name="formatNfe">
                        <xsl:with-param name="nfe" select="//n:chNFe"/>
                      </xsl:call-template>

                    </td>
                  </tr>
                  <tr>
                    <td width="140" bgColor="#ffffff" colspan="2">
                      <strong>Número NF-e: </strong>
                    </td>
                    <td width="50">
                      <xsl:value-of select="//n:nNF"/>
                    </td>
                    <td bgColor="#ffffff" colspan="2">
                      <strong>Série: </strong>
                    </td>
                    <td width="50" bgColor="#ffffff">
                      <xsl:value-of select="//n:serie"/>
                    </td>
                    <td width="140" bgColor="#ffffff">
                      <strong>Data de Emissão: </strong>
                    </td>
                    <td width="140" bgColor="#ffffff">
                      <xsl:variable name="dEmi" select="//n:dEmi"/>
                      <xsl:call-template name="formatDate">
                        <xsl:with-param name="date" select="$dEmi"/>
                      </xsl:call-template>
                    </td>

                  </tr>

                  <tr>
                    <td bgColor="#ffffff" colspan="2">
                      <strong>Número do Protocolo:</strong>
                    </td>
                    <td colspan="5" bgColor="#ffffff">
                      <xsl:value-of select="//n:nProt"/>
                    </td>
                  </tr>
                  <tr>
                    <td colspan="7"></td>
                  </tr>
                  <tr>
                    <td colspan="8" bgColor="#e8e8e8">
                      <strong>Emitente</strong>
                    </td>
                  </tr>
                  <tr>
                    <td width="13" bgColor="#ffffff">
                    </td>
                    <td width="140" bgColor="#ffffff">
                      <strong>Razão Social:</strong>
                    </td>
                    <td bgColor="#ffffff" colspan="6">
                      <xsl:value-of select="//n:infNFe/n:emit/n:xNome"/>
                    </td>
                  </tr>
                  <tr>
                    <td width="13" bgColor="#ffffff">
                    </td>
                    <td width="140" bgColor="#ffffff">
                      <strong>CNPJ:</strong>
                    </td>
                    <td colspan="5" bgColor="#ffffff">
                      <xsl:call-template name="formatCnpj">
                        <xsl:with-param name="cnpj" select="//n:infNFe/n:emit/n:CNPJ"/>
                      </xsl:call-template>
                    </td>
                  </tr>
                  <tr>
                    <td width="13">
                    </td>
                    <td width="120" bgColor="#ffffff">
                      <strong>UF:</strong>
                    </td>
                    <td colspan="5" bgColor="#ffffff">
                      <xsl:value-of select="//n:infNFe/n:emit/n:enderEmit/n:UF"/>
                    </td>
                  </tr>
                  <tr>
                    <td colspan="8">
                    </td>
                  </tr>
                  <tr>
                    <td colspan="8" bgColor="#e8e8e8">
                      <strong>Destinatário</strong>
                    </td>
                  </tr>
                  <tr>
                    <td width="13" bgColor="#ffffff">
                    </td>
                    <td bgColor="#ffffff" width="140">
                      <strong>Nome/Razão Social:</strong>
                    </td>
                    <td colspan="6" bgColor="#ffffff">
                      <xsl:value-of select="//n:infNFe/n:dest/n:xNome"/>
                    </td>
                  </tr>
                  <tr>
                    <td width="13">
                    </td>
                    <td width="140" bgColor="#ffffff">
                      <strong>

                        <xsl:if test="//n:infNFe/n:dest/n:CPF != ''">
                          CPF:
                        </xsl:if>
                        <xsl:if test="//n:infNFe/n:dest/n:CNPJ != ''">
                          CNPJ:
                        </xsl:if>

                      </strong>
                    </td>
                    <td bgColor="#ffffff" colspan="5">
                      <xsl:if test="//n:infNFe/n:dest/n:CPF != ''">
                        <xsl:call-template name="formatCpf">
                          <xsl:with-param name="cpf" select="//n:infNFe/n:dest/n:CPF">
                          </xsl:with-param>
                        </xsl:call-template>
                      </xsl:if>
                      <xsl:if test="//n:infNFe/n:dest/n:CNPJ != ''">
                        <xsl:call-template name="formatCnpj">
                          <xsl:with-param name="cnpj" select="//n:infNFe/n:dest/n:CNPJ">
                          </xsl:with-param>
                        </xsl:call-template>
                      </xsl:if>
                    </td>
                  </tr>
                  <tr>
                    <td width="13">
                    </td>
                    <td width="140" bgColor="#ffffff">
                      <strong>UF:</strong>
                    </td>
                    <td colspan="5" bgColor="#ffffff">
                      <xsl:value-of select="//n:infNFe/n:dest/n:enderDest/n:UF"/>
                    </td>
                  </tr>
                  <tr>
                    <td colspan="8">
                    </td>
                  </tr>
                  <tr>
                    <td colspan="8" bgColor="#e8e8e8">
                      <strong>Data e hora da Autorização de Uso</strong>
                    </td>
                  </tr>
                  <tr>
                    <td width="13" bgColor="#ffffff">
                    </td>
                    <td width="90" bgColor="#ffffff">
                      <strong>Autorização:</strong>
                    </td>
                    <td colspan="5" bgColor="#ffffff">
                      <xsl:call-template name="formatDateTime">
                        <xsl:with-param name="dateTime" select="//n:dhRecbto"/>
                      </xsl:call-template>
                    </td>
                  </tr>
                  <tr>
                    <td colspan="7">
                    </td>
                  </tr>
                  <tr>
                    <td width="13" bgColor="#ffffff">
                    </td>
                    <td>
                      <strong>
                        <em>Digest Value</em> da NF-e:
                      </strong>
                    </td>
                    <td colspan="5">
                      <xsl:value-of select="//n:digVal"/>
                    </td>
                  </tr>
                  
                  <table width="530" border="0" cellSpacing="0" cellPadding="5">
                    <tr>
                      <td bgColor="#ffffff">
                       <strong>Base de Cálculo do ICMS:</strong>
                      </td>
                      <td bgColor="#ffffff">
                        <strong>Valor do ICMS:</strong>
                      </td>
                      <td bgColor="#ffffff" colspan="2">
                        <strong>Valor total da NF-e:</strong>
                      </td>
                    </tr>
                    <tr>
                      <td bgColor="#ffffff">
                         <xsl:call-template name="format2Casas">
                        <xsl:with-param name="num" select="//n:infNFe/n:total/n:ICMSTot/n:vBC" />
                      </xsl:call-template>
                      </td>
                      <td bgColor="#ffffff">
                        <xsl:call-template name="format2Casas">
                          <xsl:with-param name="num" select="//n:infNFe/n:total/n:ICMSTot/n:vICMS"/>
                        </xsl:call-template>
                      </td>
                      <td bgColor="#ffffff" colspan="2">
                        <xsl:call-template name="format2Casas">
                        <xsl:with-param name="num" select="//n:infNFe/n:total/n:ICMSTot/n:vNF"/>
                      </xsl:call-template>
                      </td>
                    </tr>
                    <tr>
                      <tr colspan="3"></tr>
                    </tr>
                  </table>
                </table>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </xsl:template>
</xsl:stylesheet>