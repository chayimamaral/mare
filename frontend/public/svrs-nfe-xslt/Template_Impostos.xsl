<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:output method="html"/>
  <xsl:template match="ORIGEM_MERCADORIA" name="ORIGEM_MERCADORIA">
    <xsl:param name="orig"/>
    <xsl:choose>
      <xsl:when test="$orig = 0">0 - Nacional</xsl:when>
      <xsl:when test="$orig = 1">1 - Estrangeira - Importação direta</xsl:when>
      <xsl:when test="$orig = 2">2 - Estrangeira - Adquirida no Mercado Interno</xsl:when>
      <xsl:when test="$orig = 3">3 - Nacional, mercadoria ou bem com Conteúdo de Importação superior a 40% e inferior ou igual a 70%</xsl:when>
      <xsl:when test="$orig = 4">4 - Nacional, com produção em conformidade com processo produtivo básico previsto na legislação</xsl:when>
      <xsl:when test="$orig = 5">5 - Nacional, com Conteúdo de Importação inferior ou igual a 40%</xsl:when>
      <xsl:when test="$orig = 6">6 - Estrangeira - Importação direta, sem similar nacional, constante em lista da CAMEX e gás natural</xsl:when>
      <xsl:when test="$orig = 7">7 - Estrangeira - Adquirida no mercado interno, sem similar nacional, constante em lista da CAMEX e gás natural</xsl:when>
      <xsl:when test="$orig = 8">8 - Nacional, com Conteúdo de Importação superior a 70%</xsl:when>
      <xsl:otherwise><xsl:value-of select="$orig"/></xsl:otherwise>
    </xsl:choose> 
  </xsl:template>
  <xsl:template match="MOD_DEF_BC_ICMS_NORMAL" name="MOD_DEF_BC_ICMS_NORMAL">
    <xsl:param name="modBC"/>
    <xsl:choose>
      <xsl:when test="$modBC = 0">0 - Margem Valor Agregado(%)</xsl:when>
      <xsl:when test="$modBC = 1">1 - Pauta (valor)</xsl:when>
      <xsl:when test="$modBC = 2">2 - Preço Tabelado Máx. (valor)</xsl:when>
      <xsl:when test="$modBC = 3">3 - Valor da Operação</xsl:when>
      <xsl:otherwise><xsl:value-of select="$modBC"/></xsl:otherwise>
    </xsl:choose> 
  </xsl:template> 
  <xsl:template match="MOT_DES_ICMS" name="MOT_DES_ICMS">
    <xsl:param name="motDestICMS"/> 
    <xsl:choose>
      <xsl:when test="$motDestICMS = 3">3 - Uso na agropecuária</xsl:when>
      <xsl:when test="$motDestICMS = 6">6 - Utilitários e Motocicletas da Amazônia Ocidental e Áreas de Livre Comércio</xsl:when>
      <xsl:when test="$motDestICMS = 7">7 - SUFRAMA</xsl:when>
      <xsl:when test="$motDestICMS = 9">9 - Outros</xsl:when>
      <xsl:when test="$motDestICMS = 12">12 - Órgão de fomento e desenvolvimento agropecuário</xsl:when>
      <xsl:otherwise><xsl:value-of select="$motDestICMS"/></xsl:otherwise>
    </xsl:choose>
  </xsl:template>
  <xsl:template match="MOD_DEF_BC_ICMS_ST" name="MOD_DEF_BC_ICMS_ST">
    <xsl:param name="modBcst"/>
    <xsl:choose>
      <xsl:when test="$modBcst = 0">0 - Preço tabelado ou máx. sugerido</xsl:when>
      <xsl:when test="$modBcst = 1">1 - Lista Negativa (valor)</xsl:when>
      <xsl:when test="$modBcst = 2">2 - Lista Positiva (valor)</xsl:when>
      <xsl:when test="$modBcst = 3">3 - Lista Neutra (valor)</xsl:when>
      <xsl:when test="$modBcst = 4">4 - Margem Valor Agregado (%)</xsl:when>
      <xsl:when test="$modBcst = 5">5 - Pauta (valor)</xsl:when> 
      <xsl:otherwise><xsl:value-of select="$modBcst"/></xsl:otherwise>
    </xsl:choose> 
  </xsl:template>
  <xsl:template match="ICMS_CST" name="ICMS_CST">
    <xsl:param name="cst"/>
    <xsl:choose>
      <xsl:when test="$cst = 0">00 - Tributada integralmente</xsl:when>
      <xsl:when test="$cst = 10">10 - Tributada e com cobrança do ICMS por substituição tributária</xsl:when>
      <xsl:when test="$cst = 20">20 - Com redução de base de cálculo</xsl:when>
      <xsl:when test="$cst = 30">30 - Isenta ou não tributada e com cobrança do ICMS por substituição tributária</xsl:when>
      <xsl:when test="$cst = 40">40 - Isenta</xsl:when>
      <xsl:when test="$cst = 41">41 - Não tributada</xsl:when>
      <xsl:when test="$cst = 50">50 - Suspensão</xsl:when>
      <xsl:when test="$cst = 51">51 - Diferimento</xsl:when>
      <xsl:when test="$cst = 60">60 - ICMS cobrado anteriormente por substituição tributária</xsl:when>
      <xsl:when test="$cst = 70">70 - Com redução de base de cálculo e cobrança do ICMS por substituição tributária</xsl:when>
      <xsl:when test="$cst = 90">90 - Outros</xsl:when>
      <xsl:otherwise><xsl:value-of select="$cst"/></xsl:otherwise>
    </xsl:choose> 
  </xsl:template>
  <xsl:template match="ICMSPART_CST" name="ICMSPART_CST">
    <xsl:param name="cst"/>
      <xsl:if test="$cst='10'">10 - Tributada e com cobrança do ICMS por substituição tributária - Partilha</xsl:if>
      <xsl:if test="$cst='90'">90 - Outros - Partilha</xsl:if>
  </xsl:template>
  <xsl:template match="ICMSST_CST" name="ICMSST_CST">
    <xsl:param name="cst"/>
      <xsl:if test="$cst='41'">41 - Não Tributada - ICMS-ST</xsl:if>
  </xsl:template>
  <xsl:template match="ICMS00" name="ICMS00">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
          <xsl:variable name="orig" select="n:imposto/n:ICMS/n:ICMS00/n:orig"/>
          <xsl:call-template name="ORIGEM_MERCADORIA">
            <xsl:with-param name="orig" select="$orig"/>
          </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Tributação do ICMS</label>
          <span>
              <xsl:variable name="cst" select="n:imposto/n:ICMS/n:ICMS00/n:CST"/>
              <xsl:call-template name="ICMS_CST">
                <xsl:with-param name="cst" select="$cst"/>
              </xsl:call-template>
           </span>
        </td>
        <td>
          <label>Modalidade Definição da BC ICMS NORMAL</label>
          <span>
          <xsl:variable name="modBC" select="n:imposto/n:ICMS/n:ICMS00/n:modBC"/>
          <xsl:call-template name="MOD_DEF_BC_ICMS_NORMAL">
            <xsl:with-param name="modBC" select="$modBC"/>
          </xsl:call-template>
           </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Base de Cálculo do ICMS Normal</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS00/n:vBC">
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota do ICMS Normal</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS00/n:pICMS">
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor do ICMS Normal</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS00/n:vICMS"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template match="ICMS10" name="ICMS10">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
            <xsl:variable name="orig" select="n:imposto/n:ICMS/n:ICMS10/n:orig"/>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="$orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Tributação do ICMS</label>
          <span>
              <xsl:variable name="cst" select="n:imposto/n:ICMS/n:ICMS10/n:CST"/>
              <xsl:call-template name="ICMS_CST">
                <xsl:with-param name="cst" select="$cst"/>
              </xsl:call-template>
           </span>
        </td>
        <td>
          <label>Modalidade Definição da BC ICMS NORMAL</label>
          <span>
            <xsl:variable name="modBC" select="n:imposto/n:ICMS/n:ICMS10/n:modBC"/>
              <xsl:call-template name="MOD_DEF_BC_ICMS_NORMAL">
                <xsl:with-param name="modBC" select="$modBC"/>
              </xsl:call-template>
           </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Base de Cálculo do ICMS Normal</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS10/n:vBC">
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota do ICMS Normal</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS10/n:pICMS">
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor do ICMS Normal</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS10/n:vICMS"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Base de Cálculo do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS10/n:vBCST">
              <xsl:if test="text() != ''">
                 <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS10/n:pICMSST">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor do ICMS ST</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS10/n:vICMSST"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual Redução de BC do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS10/n:pRedBCST">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Percentual do MVA do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS10/n:pMVAST">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Modalidade Definição da BC ICMS ST</label>
          <span>
              <xsl:for-each select="n:imposto/n:ICMS/n:ICMS10/n:modBCST">
              <xsl:call-template name="MOD_DEF_BC_ICMS_ST">
                <xsl:with-param name="modBcst" select="text()"/>
              </xsl:call-template>
            </xsl:for-each>
            </span>
        </td>
      </tr>
      <tr class="col-3">
      </tr>
    </table>
  </xsl:template>
  <xsl:template match="ICMS20" name="ICMS20">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
          <xsl:variable name="orig" select="n:imposto/n:ICMS/n:ICMS20/n:orig"/>
          <xsl:call-template name="ORIGEM_MERCADORIA">
            <xsl:with-param name="orig" select="$orig"/>
          </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Tributação do ICMS</label>
          <span>
              <xsl:variable name="cst" select="n:imposto/n:ICMS/n:ICMS20/n:CST"/>
              <xsl:call-template name="ICMS_CST">
                <xsl:with-param name="cst" select="$cst"/>
              </xsl:call-template>
           </span>
        </td>
        <td>
          <label>Modalidade Definição da BC do ICMS</label>
          <span>
              <xsl:variable name="modBC" select="n:imposto/n:ICMS/n:ICMS20/n:modBC"/>
              <xsl:call-template name="MOD_DEF_BC_ICMS_NORMAL">
                <xsl:with-param name="modBC" select="$modBC"/>
              </xsl:call-template>
           </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Base de Cálculo</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS20/n:vBC">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS20/n:pICMS">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS20/n:vICMS"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual Redução de BC do ICMS Normal</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS20/n:pRedBC">
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor ICMS Desonerado</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS20/n:vICMSDeson"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Motivo Desoneração ICMS</label>
          <span>
            <xsl:call-template name="MOT_DES_ICMS">
              <xsl:with-param name="motDestICMS" select="n:imposto/n:ICMS/n:ICMS20/n:motDesICMS"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template match="ICMS30" name="ICMS30">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
              <xsl:variable name="orig" select="n:imposto/n:ICMS/n:ICMS30/n:orig"/>
              <xsl:call-template name="ORIGEM_MERCADORIA">
                <xsl:with-param name="orig" select="$orig"/>
              </xsl:call-template>
           </span>
        </td>
        <td>
          <label>Tributação do ICMS</label>
          <span>
              <xsl:variable name="cst" select="n:imposto/n:ICMS/n:ICMS30/n:CST"/>
              <xsl:call-template name="ICMS_CST">
                <xsl:with-param name="cst" select="$cst"/>
              </xsl:call-template>
           </span>
        </td>
        <td>
          <label>Modalidade Definição da BC do ICMS ST</label>
          <span>
              <xsl:variable name="modBcst" select="n:imposto/n:ICMS/n:ICMS30/n:modBCST"/>
              <xsl:call-template name="MOD_DEF_BC_ICMS_ST">
                <xsl:with-param name="modBcst" select="$modBcst"/>
              </xsl:call-template>
           </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual de Redução de BC do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS30/n:pRedBCST">
              <xsl:if test="text() != ''">
                 <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Percentual do MVA do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS30/n:pMVAST">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor da BC do ICMS ST</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS30/n:vBCST"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Alíquota do Imposto do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS30/n:pICMSST">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor do ICMS ST</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS30/n:vICMSST"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Valor ICMS Desonerado</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS30/n:vICMSDeson"/>
            </xsl:call-template>
          </span>
        </td> 
      </tr>
      <tr class="col-3">
        <td>
          <label>Motivo Desoneração ICMS</label>
          <span>
            <xsl:call-template name="MOT_DES_ICMS">
              <xsl:with-param name="motDestICMS" select="n:imposto/n:ICMS/n:ICMS30/n:motDesICMS"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template match="ICMS40" name="ICMS40">
    <table class="box">
      <tr class="col-2">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
              <xsl:variable name="orig" select="n:imposto/n:ICMS/n:ICMS40/n:orig"/>
              <xsl:call-template name="ORIGEM_MERCADORIA">
                <xsl:with-param name="orig" select="$orig"/>
              </xsl:call-template>
           </span>
        </td>
        <td colspan="2">
          <label>Tributação do ICMS</label>
          <span>
              <xsl:variable name="cst" select="n:imposto/n:ICMS/n:ICMS40/n:CST"/>
              <xsl:call-template name="ICMS_CST">
                <xsl:with-param name="cst" select="$cst"/>
              </xsl:call-template>
           </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor ICMS desoneração</label>
          <span>

            <xsl:if test="n:imposto/n:ICMS/n:ICMS40/n:vICMS">
              <xsl:call-template name="format2Casas">
                <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS40/n:vICMS"/>
              </xsl:call-template>
            </xsl:if>
            <xsl:if test="n:imposto/n:ICMS/n:ICMS40/n:vICMSDeson">
              <xsl:call-template name="format2Casas">
                <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS40/n:vICMSDeson"/>
              </xsl:call-template>
            </xsl:if>
          </span>
        </td>
        <xsl:for-each select="n:imposto/n:ICMS/n:ICMS40/n:motDesICMS">
          <td>
            <label>Motivo da desoneração do ICMS</label>
            <span>
              <xsl:variable name="mDes" select="text()"/>
              <xsl:choose>
                <xsl:when test="$mDes = 1">1 - Táxi</xsl:when>
                <xsl:when test="$mDes = 2">2 - Deficiente Físico</xsl:when>
                <xsl:when test="$mDes = 3">3 - Produtor Agropecuário</xsl:when>
                <xsl:when test="$mDes = 4">4 - Frotista/Locadora</xsl:when>
                <xsl:when test="$mDes = 5">5 - Diplomático/Consular</xsl:when>
                <xsl:when test="$mDes = 6">6 - Utilitários e Motocicletas da  Amazônia Ocidental e Áreas de Livre Comércio</xsl:when>
                <xsl:when test="$mDes = 7">7 - SUFRAMA</xsl:when>
                <xsl:when test="$mDes = 8">8 - Venda a Órgão Público</xsl:when>
                <xsl:when test="$mDes = 9">9 - Outros</xsl:when>
                <xsl:when test="$mDes = 10">10 - Deficiente Condutor</xsl:when>
                <xsl:when test="$mDes = 11">11 - Deficiente Não Condutor (Convênio ICMS 38/12)</xsl:when>
                <xsl:otherwise><xsl:value-of select="$mDes"/></xsl:otherwise>               
              </xsl:choose> 
            </span>
          </td>
        </xsl:for-each> 
      </tr>
    </table>
  </xsl:template>
  <xsl:template match="ICMS51" name="ICMS51">
    <table class="box">
      <tr class="col-2">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
            <xsl:variable name="orig" select="n:imposto/n:ICMS/n:ICMS51/n:orig"/>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="$orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Tributação do ICMS</label>
          <span>
            <xsl:variable name="cst" select="n:imposto/n:ICMS/n:ICMS51/n:CST"/>
            <xsl:call-template name="ICMS_CST">
              <xsl:with-param name="cst" select="$cst"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
    </table>
    <br/>
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Modalidade de determinação da BC do ICMS</label>
          <span>
            <xsl:variable name="modBC" select="n:imposto/n:ICMS/n:ICMS51/n:modBC"/>
            <xsl:call-template name="MOD_DEF_BC_ICMS_NORMAL">
              <xsl:with-param name="modBC" select="$modBC"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Percentual da Redução de BC</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS51/n:pRedBC">
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor da BC do ICMS</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS51/n:vBC"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Alíquota do imposto</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS51/n:pICMS">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor do ICMS</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS51/n:vICMS"/>
            </xsl:call-template>
          </span>
        </td> 
        <td>
          <label>Valor do ICMS da Operação</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS51/n:vICMSOp"/>
            </xsl:call-template>
          </span>
        </td>
      </tr> 
      <tr class="col-3">
        <td>
          <label>Percentual do Diferimento</label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS51/n:pDif"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Valor do ICMS Diferido</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS51/n:vICMSDif"/>
            </xsl:call-template>
          </span>
        </td>
      </tr> 
    </table>
  </xsl:template>
  <xsl:template match="ICMS60" name="ICMS60">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
              <xsl:variable name="orig" select="n:imposto/n:ICMS/n:ICMS60/n:orig"/>
              <xsl:call-template name="ORIGEM_MERCADORIA">
                <xsl:with-param name="orig" select="$orig"/>
              </xsl:call-template>
          </span>
        </td>
        <td colspan="2">
          <label>Tributação do ICMS</label>
          <span>
              <xsl:variable name="cst" select="n:imposto/n:ICMS/n:ICMS60/n:CST"/>
              <xsl:call-template name="ICMS_CST">
                <xsl:with-param name="cst" select="$cst"/>
              </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Valor da BC do ICMS ST retido</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS60/n:vBCSTRet"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor do ICMS ST retido</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS60/n:vICMSSTRet"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template match="ICMS70" name="ICMS70">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
              <xsl:variable name="orig" select="n:imposto/n:ICMS/n:ICMS70/n:orig"/>
              <xsl:call-template name="ORIGEM_MERCADORIA">
                <xsl:with-param name="orig" select="$orig"/>
              </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Tributação do ICMS</label>
          <span>
              <xsl:variable name="cst" select="n:imposto/n:ICMS/n:ICMS70/n:CST"/>
              <xsl:call-template name="ICMS_CST">
                <xsl:with-param name="cst" select="$cst"/>
              </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Modalidade</label>
          <span>
              <xsl:variable name="modBC" select="n:imposto/n:ICMS/n:ICMS70/n:modBC"/>
              <xsl:call-template name="MOD_DEF_BC_ICMS_NORMAL">
                <xsl:with-param name="modBC" select="$modBC"/>
              </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual Redução de BC do ICMS Normal</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS70/n:pRedBC">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Base de Cálculo</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS70/n:vBC">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS70/n:pICMS">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor do ICMS</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS70/n:vICMS"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Modalidade&#160;de&#160;Determinação&#160;da&#160;BC&#160;do&#160;ICMS&#160;ST</label>
          <span>
                <xsl:call-template name="MOD_DEF_BC_ICMS_ST">
                <xsl:with-param name="modBcst" select="n:imposto/n:ICMS/n:ICMS70/n:modBCST"/>
              </xsl:call-template>
           </span>
        </td>
        <td>
          <label>Percentual da Redução de BC do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS70/n:pRedBCST">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual da MVA do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS70/n:pMVAST">
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor da BC do ICMS ST</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS70/n:vBCST"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Alíquota do Imposto do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS70/n:pICMSST">
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor do ICMS ST</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS70/n:vICMSST"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Valor ICMS Desonerado</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS70/n:vICMSDeson"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Motivo Desoneração ICMS</label>
          <span>
            <xsl:call-template name="MOT_DES_ICMS">
              <xsl:with-param name="motDestICMS" select="n:imposto/n:ICMS/n:ICMS70/n:motDesICMS"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template match="ICMS90" name="ICMS90">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
            <xsl:variable name="orig" select="n:imposto/n:ICMS/n:ICMS90/n:orig"/>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="$orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Tributação do ICMS</label>
          <span>
              <xsl:variable name="cst" select="n:imposto/n:ICMS/n:ICMS90/n:CST"/>
              <xsl:call-template name="ICMS_CST">
                <xsl:with-param name="cst" select="$cst"/>
              </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Modalidade</label>
          <span>
              <xsl:variable name="modBC" select="n:imposto/n:ICMS/n:ICMS90/n:modBC"/>
              <xsl:call-template name="MOD_DEF_BC_ICMS_NORMAL">
                <xsl:with-param name="modBC" select="$modBC"/>
              </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual Redução de BC do ICMS Normal</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS90/n:pRedBC">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Base de Cálculo</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS90/n:vBC">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS90/n:pICMS">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor do ICMS</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS90/n:vICMS"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Modalidade&#160;de&#160;Determinação&#160;da&#160;BC&#160;do&#160;ICMS&#160;ST</label>
          <span>
              <xsl:call-template name="MOD_DEF_BC_ICMS_ST">
                <xsl:with-param name="modBcst" select="n:imposto/n:ICMS/n:ICMS90/n:modBCST"/>
              </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Percentual da Redução de BC do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS90/n:pRedBCST">
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual da MVA do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS90/n:pMVAST">
              <xsl:if test="text() != ''">
                 <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor da BC do ICMS ST</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS90/n:vBCST"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Alíquota do Imposto do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMS90/n:pICMSST">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor do ICMS ST</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS90/n:vICMSST"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Valor ICMS Desonerado</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMS90/n:vICMSDeson"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Motivo Desoneração ICMS</label>
          <span>
            <xsl:call-template name="MOT_DES_ICMS">
              <xsl:with-param name="motDestICMS" select="n:imposto/n:ICMS/n:ICMS90/n:motDesICMS"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template name="ICMSPart" match="ICMSPart">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span >
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="n:imposto/n:ICMS/n:ICMSPart/n:orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>
            Tributação do ICMS
          </label>
          <span>
            <xsl:call-template name="ICMSPART_CST" >
              <xsl:with-param name="cst" select="n:imposto/n:ICMS/n:ICMSPart/n:CST"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>
            Modalidade 
          </label>
          <span>
            <xsl:call-template name="MOD_DEF_BC_ICMS_NORMAL">
              <xsl:with-param name="modBC" select="n:imposto/n:ICMS/n:ICMSPart/n:modBC"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>
            Percentual Redução de BC 
          </label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSPart/n:pRedBC"/>
            </xsl:call-template>  
          </span>
        </td>
        <td>
          <label>
            Valor da BC do ICMS
          </label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSPart/n:vBC"/>
            </xsl:call-template>  
          </span>
        </td>
        <td>
          <label>
            Alíquota do Imposto
          </label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSPart/n:pICMS"/>
            </xsl:call-template> 
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>
            Valor
          </label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSPart/n:vICMS"/>
            </xsl:call-template> 
          </span>
        </td>
        <td>
          <label>
            Modalidade de Determinação BC do ICMS ST
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSPart/n:modBCST" >
              <xsl:call-template name="MOD_DEF_BC_ICMS_ST">
                <xsl:with-param name="modBcst" select="text()"/>
              </xsl:call-template>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>
            Percentual Redução de BC do ICMS ST
          </label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSPart/n:pRedBCST"/>
            </xsl:call-template> 
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>
            Percentual Margem de Valor Adicionado ICMS ST
          </label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSPart/n:pMVAST"/>
            </xsl:call-template> 
          </span>
        </td>
        <td>
          <label>
            Valor da BC do ICMS ST
          </label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSPart/n:vBCST"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>
            Alíquota do Imposto do ICMS ST
          </label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSPart/n:pICMSST"/>
            </xsl:call-template> 
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>
            Valor do ICMS ST<br />
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSPart/n:vICMSST" >
              <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>
            Percentual da BC operação própria<br />
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSPart/n:pBCOp" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>
            UF<br />
          </label>
          <span>
            <xsl:value-of select="n:imposto/n:ICMS/n:ICMSPart/n:UFST"/>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template name="ICMSST" match="ICMSST">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>
            Origem da Mercadoria
          </label>
          <span>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="n:imposto/n:ICMS/n:ICMSST/n:orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>
            Tributação do ICMS
          </label>
          <span>
            <xsl:call-template name="ICMSST_CST">
              <xsl:with-param name="cst" select="n:imposto/n:ICMS/n:ICMSST/n:CST"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>
            Valor da BC do ICMS ST
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSST/n:vBCSTRet" >
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>
            Valor do ICMS ST retido na UF remetente
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSST/n:vICMSSTRet" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor da BC do ICMS ST da UF destino</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSST/n:vBCSTDest" >
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor do ICMS ST da UF destino</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSST/n:vICMSSTDest" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template name="ICMS_SN_101" match="ICMS_SN_101">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="n:imposto/n:ICMS/n:ICMSSN101/n:orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Código de Situação da Operação ? Simples Nacional</label>
          <span>
            <xsl:variable name="origmerc" select="n:imposto/n:ICMS/n:ICMSSN101/n:CSOSN"/>
            <xsl:if test="$origmerc='101'">
              101- Tributada pelo Simples Nacional com permissão de crédito
            </xsl:if>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Alíquota aplicável de cálculo do crédito </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN101/n:pCredSN" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor de crédito do ICMS</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN101/n:vCredICMSSN" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template name="ICMS_SN_102" match="ICMS_SN_102">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="n:imposto/n:ICMS/n:ICMSSN102/n:orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Código de Situação da Operação - Simples Nacional</label>
          <span>
            <xsl:variable name="origmerc" select="n:imposto/n:ICMS/n:ICMSSN102/n:CSOSN"/>
            <xsl:choose>
              <xsl:when test="$origmerc='102'">
                102 - Tributada pelo Simples Nacional sem permissão de crédito
              </xsl:when>
              <xsl:when test="$origmerc='103'">
                103 - Isenção do ICMS no Simples Nacional para faixa de receita bruta
              </xsl:when>
              <xsl:when test="$origmerc='300'">
                300 - Imune
              </xsl:when>
              <xsl:when test="$origmerc='400'">
                400 - Não tributada pelo Simples Nacional
              </xsl:when>              
              <xsl:otherwise>
                <xsl:value-of select="$origmerc"/>
              </xsl:otherwise>
            </xsl:choose> 
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template name="ICMS_SN_201" match="ICMS_SN_201">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="n:imposto/n:ICMS/n:ICMSSN201/n:orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Código de Situação da Operação</label>
          <span>
            <xsl:variable name="origmerc" select="n:imposto/n:ICMS/n:ICMSSN201/n:CSOSN"/>
            <xsl:if test="$origmerc='201'">
              201 - Tributada pelo Simples Nacional com permissão de crédito e com cobrança do ICMS por Substituição Tributária
            </xsl:if>
          </span>
        </td>
        <td>
          <label>Modalidade de determinação da BC do ICMS ST</label>
          <span>
            <xsl:variable name="modBCST" select="n:imposto/n:ICMS/n:ICMSSN201/n:modBCST"/>
            <xsl:call-template name="MOD_DEF_BC_ICMS_ST">
              <xsl:with-param name="modBcst" select="$modBCST"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual Margem Valor Adicionado do ICMS ST </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN201/n:pMVAST" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Percentual da Redução de BC do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN201/n:pRedBCST" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor da BC do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN201/n:vBCST" >
              <xsl:if test="text() != ''">
                    <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Alíquota do imposto do ICMS ST </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN201/n:pICMSST" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>
            Valor do ICMS ST
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN201/n:vICMSST" >
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota aplicável de cálculo do crédito</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN201/n:pCredSN" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>
            Valor crédito do ICMS
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN201/n:vCredICMSSN" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template name="ICMS_SN_202" match="ICMS_SN_202">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig"  select="n:imposto/n:ICMS/n:ICMSSN202/n:orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Código de Situação da Operação</label>
          <span>
            <xsl:variable name="origmerc" select="n:imposto/n:ICMS/n:ICMSSN202/n:CSOSN"/>
            <xsl:if test="$origmerc='202'">
              202 - Tributada pelo Simples Nacional sem permissão de crédito e com cobrança do ICMS por Substituição Tributária
            </xsl:if>
            <xsl:if test="$origmerc='203'">
              203 - Isenção do ICMS nos Simples Nacional para faixa de receita bruta e com cobrança do ICMS por Substituição Tributária
            </xsl:if>
          </span>
        </td>
        <td>
          <label>Modalidade Determinação BC do ICMS ST</label>
          <span>
              <xsl:call-template name="MOD_DEF_BC_ICMS_ST">
                <xsl:with-param name="modBcst" select="n:imposto/n:ICMS/n:ICMSSN202/n:modBCST"/>
              </xsl:call-template>            
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual Margem Valor Adicionado ICMS ST </label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSSN202/n:pMVAST"/>
            </xsl:call-template> 
          </span>
        </td>
        <td>
          <label>Percentual Redução de BC do ICMS ST</label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSSN202/n:pRedBCST"/>
            </xsl:call-template> 
          </span>
        </td>
        <td>
          <label>Valor BC do ICMS ST</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSSN202/n:vBCST"/>
            </xsl:call-template> 
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Alíquota Imposto ICMS ST </label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSSN202/n:pICMSST"/>
            </xsl:call-template>  
          </span>
        </td>
        <td >
          <label>Valor do ICMS ST </label>
          <span>
            <xsl:call-template name="format4Casas">
              <xsl:with-param name="num" select="n:imposto/n:ICMS/n:ICMSSN202/n:vICMSST"/>
            </xsl:call-template> 
          </span>
        </td>
        <td></td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template name="ICMS_SN_500" match="ICMS_SN_500">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="n:imposto/n:ICMS/n:ICMSSN500/n:orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Código de Situação da Operação</label>
          <span>
            <xsl:variable name="origmerc" select="n:imposto/n:ICMS/n:ICMSSN500/n:CSOSN"/>
            <xsl:if test="$origmerc='500'">
              500 - ICMS cobrado anteriormente por substituição tributária (substituído) ou por antecipação
            </xsl:if>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor da BC do ICMS ST retido </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN500/n:vBCSTRet" >
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor do ICMS ST retido </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN500/n:vICMSSTRet" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template name="ICMS_SN_900" match="ICMS_SN_900">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Origem da Mercadoria</label>
          <span>
            <xsl:call-template name="ORIGEM_MERCADORIA">
              <xsl:with-param name="orig" select="n:imposto/n:ICMS/n:ICMSSN900/n:orig"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Código de Situação da Operação</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:CSOSN" >
              <xsl:variable name="origmerc" select="text()"/>
              <xsl:if test="$origmerc='900'">
                900 - Outros
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Modalidade de determinação da BC do ICMS </label>
          <span>
            <xsl:call-template name="MOD_DEF_BC_ICMS_NORMAL">
              <xsl:with-param name="modBC" select="n:imposto/n:ICMS/n:ICMSSN900/n:modBC"/>
            </xsl:call-template>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor da BC do ICMS </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:vBC" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Percentual da Redução de BC </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:pRedBC" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota do imposto</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:pICMS" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor do ICMS </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:vICMS" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Modalidade de determinação da BC do ICMS ST </label>
          <span>
            <xsl:call-template name="MOD_DEF_BC_ICMS_ST">
              <xsl:with-param name="modBcst" select="n:imposto/n:ICMS/n:ICMSSN900/n:modBCST"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Percentual Margem Valor Adicionado do ICMS ST </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:pMVAST" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Percentual da Redução de BC</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:pRedBC" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Valor da BC do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:vBCST" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota do imposto ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:pICMSST" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor do ICMS ST</label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:vICMSST" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>
            Alíquota aplicável de cálculo do crédito
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:pCredSN" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>
            Valor de crédito do ICMS
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:vCredICMSSN" >
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.##.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>        
      </tr>
      <tr class="col-3">
        <td>
          <label>
            Percentual redução BC do ICMS ST
          </label>
          <span>
            <xsl:for-each select="n:imposto/n:ICMS/n:ICMSSN900/n:pRedBCST" >
              <xsl:if test="text() != ''">
                <xsl:value-of select="format-number(text(),'##.##.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
  <xsl:template match="PIS" name="PIS">
    <table class="box">
      <tr class="col-3">
        <td colspan="3">
          <label>CST</label>
          <span>
            <xsl:choose>
              <xsl:when test="n:PISAliq/n:CST = 01">01 - Operação Tributável (base de cálculo = valor da operação alíquota normal (cumulativo/não cumulativo))</xsl:when>
              <xsl:when test="n:PISAliq/n:CST = 02">02 - Operação Tributável (base de cálculo = valor da operação (alíquota diferenciada))</xsl:when>
              <xsl:when test="n:PISQtde/n:CST = 03">03 - Operação Tributável (base de cálculo = quantidade vendida x alíquota por unidade de produto)</xsl:when>
              <xsl:when test="n:PISNT/n:CST = 04">04 - Operação Tributável (tributação monofásica (alíquota zero))</xsl:when>
              <xsl:when test="n:PISNT/n:CST = 05">05 - Operação Tributável (Substituição Tributária)</xsl:when>
              <xsl:when test="n:PISNT/n:CST = 06">06 - Operação Tributável (alíquota zero)</xsl:when>
              <xsl:when test="n:PISNT/n:CST = 07">07 - Operação Isenta da Contribuição</xsl:when>
              <xsl:when test="n:PISNT/n:CST = 08">08 - Operação Sem Incidência da Contribuição</xsl:when>
              <xsl:when test="n:PISNT/n:CST = 09">09 - Operação com Suspensão da Contribuição</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 49">49 - Outras Operações de Saída</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 50">50 - Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Tributada no Mercado Interno</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 51">51 - Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Não Tributada no Mercado Interno</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 52">52 - Operação com Direito a Crédito – Vinculada Exclusivamente a Receita de Exportação</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 53">53 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 54">54 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas no Mercado Interno e de Exportação</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 55">55 - Operação com Direito a Crédito - Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 56">56 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno, e de Exportação</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 60">60 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Tributada no Mercado Interno</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 61">61 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Não-Tributada no Mercado Interno</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 62">62 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita de Exportação</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 63">63 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 64">64 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas no Mercado Interno e de Exportação</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 65">65 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 66">66 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno, e de Exportação</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 67">67 - Crédito Presumido - Outras Operações</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 70">70 - Operação de Aquisição sem Direito a Crédito</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 71">71 - Operação de Aquisição com Isenção</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 72">72 - Operação de Aquisição com Suspensão</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 73">73 - Operação de Aquisição a Alíquota Zero</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 74">74 - Operação de Aquisição; sem Incidência da Contribuição</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 75">75 - Operação de Aquisição por Substituição Tributária</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 98">98 - Outras Operações de Entrada</xsl:when>
              <xsl:when test="n:PISOutr/n:CST = 99">99 - Outras Operações</xsl:when>
            </xsl:choose>
          </span>
        </td>
      </tr>
      <xsl:variable name="cst" select="n:PISAliq/n:CST"/>
      <xsl:variable name="cstQtde" select="n:PISQtde/n:CST"/>
      <xsl:variable name="cstOutr" select="n:PISOutr/n:CST"/>
      <xsl:choose>
        <xsl:when test="$cst=01 or $cst=02">
          <tr class="col-3">
            <td>
              <label>Base de Cálculo</label>
              <span>
                <xsl:variable name="vBC" select="n:PISAliq/n:vBC"/>
                <xsl:if test="$vBC != ''">
                  <xsl:value-of select="format-number($vBC,'#.###.###.###.##0,00')"/>
                </xsl:if>
              </span>
            </td>
            <td>
              <label>Alíquota</label>
              <span>
                <xsl:variable name="pPIS" select="n:PISAliq/n:pPIS"/>
                <xsl:if test="$pPIS != ''">
                    <xsl:value-of select="format-number($pPIS,'##0,0000')"/>
                </xsl:if>
              </span>
            </td>
            <td>
              <label>Valor</label>
              <span>
                <xsl:variable name="vPIS" select="n:PISAliq/n:vPIS"/>
                <xsl:if test="$vPIS != ''">
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="$vPIS"/>
                  </xsl:call-template>
                </xsl:if>
              </span>
            </td>
          </tr>
        </xsl:when>
        <xsl:when test="$cstQtde=03">
          <tr class="col-3">
            <td>
              <label>Quantidade Vendida</label>
              <span>
                <xsl:call-template name="format4Casas">
                  <xsl:with-param name="num" select="n:PISQtde/n:qBCProd"/>
                </xsl:call-template>
              </span>
            </td>
            <td>
              <label>Alíquota (R$)</label>
              <span>
                <xsl:call-template name="format4Casas">
                  <xsl:with-param name="num" select="n:PISQtde/n:vAliqProd"/>
                </xsl:call-template> 
              </span>
            </td>
            <td>
              <label>Valor do PIS</label>
              <span>
                <xsl:call-template name="format2Casas">
                  <xsl:with-param name="num" select="n:PISQtde/n:vPIS"/>
                </xsl:call-template> 
              </span>
            </td>
          </tr>
        </xsl:when>
        <xsl:when test="$cstOutr!='' ">
          <xsl:variable name="vBCop" select="n:PISOutr/n:vBC"/>
          <xsl:if test="$vBCop != ''">
            <tr class="col-3">
              <td>
                <label>Base de Cálculo</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="n:PISOutr/n:vBC"/>
                  </xsl:call-template> 
                </span>
              </td>
              <td>
                <label>Alíquota (%)</label>
                <span>
                  <xsl:call-template name="format4Casas">
                    <xsl:with-param name="num" select="n:PISOutr/n:pPIS"/>
                  </xsl:call-template> 
                </span>
              </td>
              <td>
                <label>Valor do PIS</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="n:PISOutr/n:vPIS"/>
                  </xsl:call-template> 
                </span>
              </td>
            </tr>
          </xsl:if>
          
          <xsl:variable name="qBCProdop" select="n:PISOutr/n:qBCProd"/>
          <xsl:if test="$qBCProdop != ''">
            <tr class="col-3">
              <td>
                <label>Quantidade Vendida</label>
                <span>
                  <xsl:call-template name="format4Casas">
                    <xsl:with-param name="num" select="n:PISOutr/n:qBCProd"/>
                  </xsl:call-template> 
                </span>
              </td>
              <td>
                <label>Alíquota (R$)</label>
                <span>
                  <xsl:call-template name="format4Casas">
                    <xsl:with-param name="num" select="n:PISOutr/n:vAliqProd"/>
                  </xsl:call-template> 
                </span>
              </td>
              <td>
                <label>Valor do PIS</label>
                <span>
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="n:PISOutr/n:vPIS"/>
                  </xsl:call-template> 
                </span>
              </td>
            </tr>
          </xsl:if>
        </xsl:when>
      </xsl:choose>
    </table>
  </xsl:template>
  <xsl:template match="PISST" name="PISST">
    <xsl:variable name="vBCop" select="n:vBC"/>
    <xsl:if test="$vBCop != ''">
      <table class="box">
        <tr class="col-3">
          <td>
            <label>Base de Cálculo</label>
            <span>
              <xsl:variable name="vBC" select="n:vBC"/>
              <xsl:if test="$vBC != ''">
                <xsl:value-of select="format-number($vBC,'#.###.###.###.##0,00')"/>
              </xsl:if>
            </span>
          </td>
          <td>
            <label>Alíquota (%)</label>
            <span>
              <xsl:variable name="pPIS" select="n:pPIS"/>
              <xsl:if test="$pPIS != ''">
                                <xsl:value-of select="format-number($pPIS,'##0,0000')"/>
              </xsl:if>
            </span>
          </td>
          <td>
            <label>Valor</label>
            <span>
              <xsl:variable name="vPIS" select="n:vPIS"/>
              <xsl:if test="$vPIS != ''">
                <xsl:call-template name="format2Casas">
                  <xsl:with-param name="num" select="$vPIS"/>
                </xsl:call-template>
              </xsl:if>
            </span>
          </td>
        </tr>
      </table>
    </xsl:if>
    <xsl:variable name="qBCProdop" select="n:qBCProd"/>
    <xsl:if test="$qBCProdop != ''">
      <table class="box">
        <tr class="col-3">
          <td>
            <label>Quantidade Vendida</label>
            <span>
              <xsl:variable name="qBCProd" select="n:qBCProd"/>
              <xsl:if test="$qBCProd != ''">
                <xsl:value-of select="format-number($qBCProd,'###.###.###.##0,00##')"/>
              </xsl:if>
            </span>
          </td>
          <td>
            <label>Alíquota (R$)</label>
            <span>
              <xsl:variable name="vAliqProd" select="n:vAliqProd"/>
              <xsl:if test="$vAliqProd != ''">
                <xsl:value-of select="format-number($vAliqProd,'##.###.###.##0,00##')"/>
              </xsl:if>
            </span>
          </td>
          <td>
          </td>
        </tr>
      </table>
    </xsl:if>
  </xsl:template>
  <xsl:template match="COFINS" name="COFINS">
    <table class="box">
      <tr class="col-3">
        <td colspan="3">
          <label>CST</label>
          <span>
            <xsl:choose>
              <xsl:when test="n:COFINSAliq/n:CST=01">01 - Operação Tributável (base de cálculo = valor da operação alíquota normal (cumulativo/não cumulativo))</xsl:when>
              <xsl:when test="n:COFINSAliq/n:CST=02">02 - Operação Tributável (base de cálculo = valor da operação (alíquota diferenciada))</xsl:when>
              <xsl:when test="n:COFINSQtde/n:CST=03">03 - Operação Tributável (base de cálculo = quantidade vendida x alíquota por unidade de produto)</xsl:when>
              <xsl:when test="n:COFINSNT/n:CST = 04">04 - Operação Tributável (tributação monofásica (alíquota zero))</xsl:when>
              <xsl:when test="n:COFINSNT/n:CST = 05">05 - Operação Tributável (Substituição Tributária)</xsl:when>
              <xsl:when test="n:COFINSNT/n:CST = 06">06 - Operação Tributável (alíquota zero)</xsl:when>
              <xsl:when test="n:COFINSNT/n:CST = 07">07 - Operação Isenta da Contribuição</xsl:when>
              <xsl:when test="n:COFINSNT/n:CST = 08">08 - Operação Sem Incidência da Contribuição</xsl:when>
              <xsl:when test="n:COFINSNT/n:CST = 09">09 - Operação com Suspensão da Contribuição</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 49">49 - Outras Operações de Saída</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 50">50 - Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Tributada no Mercado Interno</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 51">51 - Operação com Direito a Crédito - Vinculada Exclusivamente a Receita Não Tributada no Mercado Interno</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 52">52 - Operação com Direito a Crédito – Vinculada Exclusivamente a Receita de Exportação</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 53">53 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 54">54 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas no Mercado Interno e de Exportação</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 55">55 - Operação com Direito a Crédito - Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 56">56 - Operação com Direito a Crédito - Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno, e de Exportação</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 60">60 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Tributada no Mercado Interno</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 61">61 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita Não-Tributada no Mercado Interno</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 62">62 - Crédito Presumido - Operação de Aquisição Vinculada Exclusivamente a Receita de Exportação</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 63">63 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 64">64 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas no Mercado Interno e de Exportação</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 65">65 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Não-Tributadas no Mercado Interno e de Exportação</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 66">66 - Crédito Presumido - Operação de Aquisição Vinculada a Receitas Tributadas e Não-Tributadas no Mercado Interno, e de Exportação</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 67">67 - Crédito Presumido - Outras Operações</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 70">70 - Operação de Aquisição sem Direito a Crédito</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 71">71 - Operação de Aquisição com Isenção</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 72">72 - Operação de Aquisição com Suspensão</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 73">73 - Operação de Aquisição a Alíquota Zero</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 74">74 - Operação de Aquisição; sem Incidência da Contribuição</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 75">75 - Operação de Aquisição por Substituição Tributária</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 98">98 - Outras Operações de Entrada</xsl:when>
              <xsl:when test="n:COFINSOutr/n:CST = 99">99 - Outras Operações</xsl:when>
            </xsl:choose>
          </span>
        </td>
      </tr>
      <xsl:variable name="cstAliq" select="n:COFINSAliq/n:CST"/>
      <xsl:variable name="cstQtde" select="n:COFINSQtde/n:CST"/>
      <xsl:variable name="cstOutr" select="n:COFINSOutr/n:CST"/>
      <xsl:choose>
        <xsl:when test="$cstAliq=01 or $cstAliq=02">
          <tr class="col-3">
            <td>
              <label>Base de Cálculo</label>
              <span>
                <xsl:variable name="vBC" select="n:COFINSAliq/n:vBC"/>
                <xsl:if test="$vBC != ''">
                  <xsl:value-of select="format-number($vBC,'#.###.###.###.##0,00')"/>
                </xsl:if>
              </span>
            </td>
            <td>
              <label>Alíquota</label>
              <span>
                <xsl:variable name="pCOFINS" select="n:COFINSAliq/n:pCOFINS"/>
                <xsl:if test="$pCOFINS != ''">
                  <xsl:value-of select="format-number($pCOFINS,'##0,0000')"/>
                </xsl:if>
              </span>
            </td>
            <td>
              <label>Valor</label>
              <span>
                <xsl:variable name="vCOFINS" select="n:COFINSAliq/n:vCOFINS"/>
                <xsl:if test="$vCOFINS != ''">
                  <xsl:call-template name="format2Casas">
                    <xsl:with-param name="num" select="$vCOFINS"/>
                  </xsl:call-template>
                </xsl:if>
              </span>
            </td>
          </tr>
        </xsl:when>
        <xsl:when test="$cstQtde=03">
          <tr class="col-3">
            <td>
              <label>Quantidade Vendida</label>
              <span>
                <xsl:call-template name="format4Casas">
                  <xsl:with-param name="num" select="n:COFINSQtde/n:qBCProd"/>
                </xsl:call-template> 
              </span>
            </td>
            <td>
              <label>Alíquota (R$)</label>
              <span>
                <xsl:call-template name="format4Casas">
                  <xsl:with-param name="num" select="n:COFINSQtde/n:vAliqProd"/>
                </xsl:call-template> 
              </span>
            </td>
            <td>
              <label>Valor</label>
              <span>
                <xsl:call-template name="format2Casas">
                  <xsl:with-param name="num" select="n:COFINSQtde/n:vCOFINS"/>
                </xsl:call-template>
              </span>
            </td>
          </tr>
        </xsl:when>
        <xsl:when test="$cstOutr!=''">
          <xsl:variable name="vBCop" select="n:COFINSOutr/n:vBC"/>
          <xsl:if test="$vBCop != ''">
            <tr class="col-3">
              <td>
                <label>Base de Cálculo</label>
                <span>
                  <xsl:variable name="vBC" select="n:COFINSOutr/n:vBC"/>
                  <xsl:if test="$vBC != ''">
                    <xsl:value-of select="format-number($vBC,'#.###.###.###.##0,00')"/>
                  </xsl:if>
                </span>
              </td>
              <td>
                <label>Alíquota (%)</label>
                <span>
                  <xsl:variable name="pCOFINS" select="n:COFINSOutr/n:pCOFINS"/>
                  <xsl:if test="$pCOFINS != ''">
                                        <xsl:value-of select="format-number($pCOFINS,'##0,0000')"/>
                  </xsl:if>
                </span>
              </td>
              <td>
                <label>Valor</label>
                <span>
                  <xsl:variable name="vCOFINS" select="n:COFINSOutr/n:vCOFINS"/>
                  <xsl:if test="$vCOFINS != ''">
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="$vCOFINS"/>
                    </xsl:call-template>
                  </xsl:if>
                </span>
              </td>
            </tr>
          </xsl:if>
          <xsl:variable name="qBCProdop" select="n:COFINSOutr/n:qBCProd"/>
          <xsl:if test="$qBCProdop != ''">
            <tr class="col-3">
              <td>
                <label>Quantidade Vendida</label>
                <span>
                  <xsl:call-template name="format4Casas">
                    <xsl:with-param name="num" select="n:COFINSOutr/n:qBCProd"/>
                  </xsl:call-template> 
                </span>
              </td>
              <td>
                <label>Alíquota (R$)</label>
                <span>
                  <xsl:call-template name="format4Casas">
                    <xsl:with-param name="num" select="n:COFINSOutr/n:vAliqProd"/>
                  </xsl:call-template> 
                </span>
              </td>
              <td>
                <label>Valor</label>
                <span>
                  <xsl:variable name="vCOFINS" select="n:COFINSOutr/n:vCOFINS"/>
                  <xsl:if test="$vCOFINS != ''">
                    <xsl:call-template name="format2Casas">
                      <xsl:with-param name="num" select="$vCOFINS"/>
                    </xsl:call-template>
                  </xsl:if>
                </span>
              </td>
            </tr>
          </xsl:if>
        </xsl:when>
      </xsl:choose>
    </table>
  </xsl:template>
  <xsl:template match="COFINSST" name="COFINSST">
    <xsl:variable name="vBCop" select="n:vBC"/>
    <xsl:if test="$vBCop != ''">
      <table class="box">
        <tr class="col-3">
          <td>
            <label>Base de Cálculo</label>
            <span>
              <xsl:for-each select="n:vBC">
                <xsl:if test="text() != ''">
                    <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
                </xsl:if>
              </xsl:for-each>
            </span>
          </td>
          <td>
            <label>Alíquota (%)</label>
            <span>
              <xsl:for-each select="n:pCOFINS">
                <xsl:if test="text() != ''">
                    <xsl:value-of select="format-number(text(),'##0,0000')"/>
                </xsl:if>
              </xsl:for-each>
            </span>
          </td>
          <td>
            <label>Valor</label>
            <span>
              <xsl:call-template name="format2Casas">
                <xsl:with-param name="num" select="n:vCOFINS"/>
              </xsl:call-template>
            </span>
          </td>
        </tr>
      </table>
    </xsl:if>
    <xsl:variable name="qBCProdop" select="n:qBCProd"/>
    <xsl:if test="$qBCProdop != ''">
      <table class="box">
        <tr class="col-3">
          <td>
            <label>Quantidade Vendida</label>
            <span>
              <xsl:for-each select="n:qBCProd">
                <xsl:if test="text() != ''">
                    <xsl:value-of select="format-number(text(),'###.###.###.##0,00##')"/>
                </xsl:if>
              </xsl:for-each>
            </span>
          </td>
          <td>
            <label>Alíquota (R$)</label>
            <span>
              <xsl:for-each select="n:vAliqProd">
                <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'##.###.###.##0,0000')"/>
                </xsl:if>
              </xsl:for-each>
            </span>
          </td>
          <td>
          </td>
        </tr>
      </table>
    </xsl:if>
  </xsl:template>
  <xsl:template match="ISSQN" name="ISSQN">
    <table class="box">
      <tr class="col-3">
        <td>
          <label>Código de Tributação do ISSQN</label>
          <span>
            <xsl:for-each select="n:cSitTrib">
              <xsl:variable name="cSitTrib" select="text()"/>
              <xsl:choose>
                <xsl:when test="$cSitTrib='N'">
                  N - NORMAL
                </xsl:when>
                <xsl:when test="$cSitTrib='R'">
                  R - RETIDA
                </xsl:when>
                <xsl:when test="$cSitTrib='S'">
                  S - SUBSTITUTA
                </xsl:when>
                <xsl:when test="$cSitTrib='I'">
                  I - ISENTA
                </xsl:when>
                <xsl:otherwise>
                  <xsl:value-of select="$cSitTrib"/>
                </xsl:otherwise>
              </xsl:choose> 
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Base de Cálculo</label>
          <span>
            <xsl:for-each select="n:vBC">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'#.###.###.###.##0,00')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
        <td>
          <label>Alíquota</label>
          <span>
            <xsl:for-each select="n:vAliq">
              <xsl:if test="text() != ''">
                  <xsl:value-of select="format-number(text(),'#.###.###.###.##0,0000')"/>
              </xsl:if>
            </xsl:for-each>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Valor</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:vISSQN"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Município</label>
          <span>
            <xsl:value-of select="n:cMunFG"/>
          </span>
        </td>
        <td>
          <label>Serviço</label>
          <span>
            <xsl:value-of select="n:cListServ"/>
          </span>
        </td>
        <td></td>
      </tr>
  
      <tr class="col-3">
        <td>
          <label>Valor dedução para redução da BC</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:vDeducao"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Valor outras retenções</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:vOutro"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Valor desconto incondicionado</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:vDescIncond"/>
            </xsl:call-template>
          </span>
        </td>       
      </tr> 

      <tr class="col-3">
        <td>
          <label>Valor desconto condicionado</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:vDescCond"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Valor retenção ISS</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:vISSRet"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Indicador da exigibilidade do ISS</label>
          <span>

            <xsl:variable name="indISS" select="n:indISS"/>
            <xsl:choose>
              <xsl:when test="$indISS = 1">
                01 = Exigível
              </xsl:when>
              <xsl:when test="$indISS = 2">
                02 = Não incidência
              </xsl:when>
              <xsl:when test="$indISS = 3">
                03 = Isenção
              </xsl:when>
              <xsl:when test="$indISS = 4">
                04 = Exportação
              </xsl:when>
              <xsl:when test="$indISS = 5">
                05 = Imunidade
              </xsl:when>
              <xsl:when test="$indISS = 6">
                06 = Exigibilidade Suspensa por Decisão Judicial
              </xsl:when>
              <xsl:when test="$indISS = 7">
                07 = Exigibilidade Suspensa por Processo Administrativo
              </xsl:when>
              <xsl:otherwise>
                <xsl:value-of select="$indISS"/>
              </xsl:otherwise>
            </xsl:choose>  
          </span>
        </td>
      </tr> 
      <tr class="col-3">
        <td>
          <label>Código Serviço Prestado</label>
          <span>
            <xsl:value-of select="n:cServico"/>
          </span>
        </td>
        <td>
          <label>Código Município Imposto</label>
          <span>
            <xsl:value-of select="n:cMun"/>
          </span>
        </td>
        <td>
          <label>Código País Serviço</label>
          <span>
            <xsl:value-of select="n:cPais"/>
          </span>
        </td>
      </tr>
      <tr class="col-3">
        <td>
          <label>Número Processo Administrativo Suspensão</label>
          <span>
            <xsl:value-of select="n:nProcesso"/>
          </span>
        </td>
        <td>
          <label>Indicador de Incentivo Fiscal</label>
          <span>
            <xsl:variable name="indIncentivo" select="n:indIncentivo"/>
            <xsl:choose>
              <xsl:when test="$indIncentivo = 1">
                1 = Sim
              </xsl:when>
              <xsl:when test="$indIncentivo = 2">
                2 = Não
              </xsl:when>
              <xsl:otherwise>
                <xsl:value-of select="$indIncentivo"/>
              </xsl:otherwise>
            </xsl:choose> 
          </span>
        </td>
        <td></td>
      </tr> 
    </table>
  </xsl:template>
  <xsl:template match="TribDevol" name="TribDevol">
    <table class="box">
      <tr class="col-2">
        <td>
          <label>Percentual da Mercadoria Devolvida</label>
          <span>
            <xsl:call-template name="format2Casas">
              <xsl:with-param name="num" select="n:pDevol"/>
            </xsl:call-template>
          </span>
        </td>
        <td>
          <label>Valor do IPI Devolvido</label>
          <span>
            <xsl:for-each select="n:IPI/n:vIPIDevol">
              <xsl:call-template name="format2Casas">
                  <xsl:with-param name="num" select="text()"/>
              </xsl:call-template>
            </xsl:for-each>
          </span>
        </td>
      </tr>
    </table>
  </xsl:template>
</xsl:stylesheet>