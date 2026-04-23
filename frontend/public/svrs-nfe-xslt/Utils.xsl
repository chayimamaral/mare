<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <xsl:decimal-format decimal-separator="," grouping-separator="."/>
  <xsl:output method="html"/>

  <!-- identificação do ambiente, com logo -->
  <xsl:template match="ID_AMBIENTE" name="ID_AMBIENTE">
    <xsl:variable name="tpAmbiente" select="//n:ide/n:tpAmb"/>
    <span id="tipo-ambiente">
      <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAgFJREFUeNpsU79LI0EU/jabpPHSKfGUg1Mwyilcc4diqV44jNqJjZ2IbowQ8A+4nK3lgikEsRArC9GcV52dmPTxB6ImIpiIeNEqP3fWN7PZZTfJYx8zO/O+b755857075cbDaaSK+Qymi1OHuaTsVhVLDSiVf+Xn5GB0DokydWEvjr+reTTRzBJGglU/+BkpD+4gtrLLlj5zgF2eT+hP7hMM51IEhaJeYzaOTQVCQTDqL39hVa+gQ7mcK1yj2ohgUBQAcUqJzHPpqlAgPsmlqC9JujkrHXq/voxdAbMxibFPytnUCscIvCDpwicRChQ+sYXoRWOwIokmzHLmcYH3blWJJL/BwjQgRzLFcisdAutlOH3c9xbgPnHZdhMK2UhFa/tSaQAEeQkMMCo79lNsmINAorSWQsCk4S1INBtBFyirmstCXg5GHtOAvNalgKXt4uy/OAIm1e3xVh+2nHWhMdvKeCvEL87PYDc9hWSt7uuxvC96IJw+5rk6YDs+45sShRTnCsI5y+TosJ6R2dQY1VS8miwy0a6UM+By/sR7g8ETibwmD7lfbFmljKRpKj+JaVnZIpI6L0recxtbBlXyG0TuBNu3zdkCJw7PxNgaqiivRdWcxekhEg+D4dEK7Lqs9jwtE+LMZP6Y4KjBK40NhNPdZQCeJCC1ibA5BVz4V2AAQBo8gRz1Ov2wQAAAABJRU5ErkJggg==" />
      Ambiente de
      <xsl:if test="$tpAmbiente='1'">Produção</xsl:if>
      <xsl:if test="$tpAmbiente='2'">Homologação</xsl:if>
    </span>
  </xsl:template>

  <!-- cabeçalho reutilizável -->
  <xsl:template match="CABECALHO_NFE" name="CABECALHO_NFE">
    <xsl:variable name="tpAmbiente" select="/n:NFe/n:infNFe/n:ide/n:tpAmb"/>
    <div class="GeralXslt">
        <h1>
          Consulta da NF-e
          <xsl:if test="$tpAmbiente='2'">
            (Ambiente de Homologação)
          </xsl:if>
        </h1>
        <fieldset>
          <legend>Dados Gerais</legend>
          <table class="box">
            <tr>
              <td>
                <label>Chave de Acesso</label>
                <span>
                  <xsl:call-template name="formatNfe">
                    <xsl:with-param name="nfe" select="substring-after(//n:infNFe/@Id,'NFe')"/>
                  </xsl:call-template>
                </span>
              </td>
              <td class="fixo-nro-serie">
                <label>Número</label>
                <span>
                  <xsl:value-of select = "//n:infNFe/n:ide/n:nNF"/>
                </span>
              </td>
              <td class="fixo-versao-xml">
                <label>Versão XML</label>
                <span>
                  <xsl:value-of select="//n:infNFe/@versao"/>
                </span>
              </td>
            </tr>
          </table>
        </fieldset>
    </div>
  </xsl:template>

  <!-- formatação de data simples "YYYY-MM-DD" para "DD/MM/YYYY" -->
  <xsl:template match="formatDate" name="formatDate">
    <xsl:param name="date"/>
    <xsl:if test="string-length($date) != 0">
      <xsl:variable name="year" select="substring-before($date, '-')"/>
      <xsl:variable name="month" select="substring-before(substring-after($date, '-'), '-')"/>
      <xsl:variable name="day" select="substring-after(substring-after($date, '-'), '-')"/>
      <xsl:value-of select="concat('', $day, '/', $month, '/', $year, '')"/>
    </xsl:if>
  </xsl:template>

  <!-- formatação de tempo em uma SQL Date "YYYY-MM-DDTHH:MM:SS" para "HH:MM:SS" -->
  <xsl:template name="formatTime">
    <xsl:param name="dateTime"/>
    <xsl:if test="string-length($dateTime) != 0">
      <xsl:value-of select="concat('', substring-after($dateTime, 'T'), '')"/>
    </xsl:if>
  </xsl:template>

  <!-- formatação completa de SQL Date "YYYY-MM-DDTHH:MM:SS" para "DD/MM/YYYY HH:MM:SS" -->
  <xsl:template match="formatDateTime" name="formatDateTime">
    <xsl:param name="dateTime"/>
    <xsl:param name="include_as"/>
    <xsl:if test="string-length($dateTime) != 0">
      <xsl:variable name="date" select="substring-before($dateTime, 'T')"/>
      <xsl:variable name="year" select="substring-before($date, '-')"/>
      <xsl:variable name="month" select="substring-before(substring-after($date, '-'), '-')"/>
      <xsl:variable name="day" select="substring-after(substring-after($date, '-'), '-')"/>
      <xsl:variable name="time" select="substring-after($dateTime, 'T')"/>
      <xsl:value-of select="concat('', $day, '/', $month, '/', $year, '', ' ')"/>
      <xsl:if test="string($include_as) != ''">
        às
      </xsl:if>
      <xsl:value-of select="concat('', substring-before(concat($time, '.000'), '.'), '')"/>
    </xsl:if>
  </xsl:template>

  <!-- formatação completa de SQL Date "YYYY-MM-DDTHH:MM:SS" para "DD/MM/YYYY HH:MM:SS" -->
  <xsl:template match="formatDateTimeFuso" name="formatDateTimeFuso">
    <xsl:param name="dateTime"/>
    <xsl:param name="include_as"/>
    <xsl:if test="string-length($dateTime) != 0">
      <xsl:variable name="date" select="substring-before($dateTime, 'T')"/>
      <xsl:variable name="year" select="substring-before($date, '-')"/>
      <xsl:variable name="month" select="substring-before(substring-after($date, '-'), '-')"/>
      <xsl:variable name="day" select="substring-after(substring-after($date, '-'), '-')"/>
      <xsl:value-of select="concat('', $day, '/', $month, '/', $year, '', ' ')"/>
      <xsl:if test="string($include_as) != ''">
        às
      </xsl:if>
      <xsl:variable name="horacompleta" select="concat('', substring-after($dateTime, 'T'), '')"/>
      <xsl:variable name="fuso" select="substring-after(substring-before($horacompleta, '-'), '')"/>
      <xsl:value-of select="$horacompleta"/>

    </xsl:if>
  </xsl:template>

  <!-- formatação de moeda com duas casas decimais (##,##) -->
  <xsl:template match="format2Casas" name="format2Casas">
    <xsl:param name="num"/>
    <xsl:if test="string-length($num) != 0">
      <xsl:value-of select="concat('', format-number($num,'##.##.##0,00'), '')"/>
    </xsl:if>
  </xsl:template>

  <!-- formatação de numero com 3 casas decimais (##,###) -->
  <xsl:template match="format3Casas" name="format3Casas">
    <xsl:param name="num"/>
    <xsl:if test="string-length($num) != 0">
      <xsl:value-of select="concat('', format-number($num,'###.###.###.##0,000'), '')"/>
    </xsl:if>
  </xsl:template>

  <!-- formatação de numero com quatro casas decimais (##,####) -->
  <xsl:template match="format4Casas" name="format4Casas">
    <xsl:param name="num"/>
    <xsl:if test="string-length($num) != 0">
      <xsl:value-of select="concat('', format-number($num,'###.###.###.##0,0000'), '')"/>
    </xsl:if>
  </xsl:template>

  <!-- formatação de numero com dez casas decimais  -->
  <xsl:template match="format10Casas" name="format10Casas">
    <xsl:param name="num"/>
    <xsl:if test="string-length($num) != 0">
      <xsl:value-of select="concat('', format-number($num,'###.###.###.##0,0000000000'), '')"/>
    </xsl:if>
  </xsl:template>

  <!-- formatação de CNPJ (##.###.###-####/## -->
  <xsl:template match="formatCnpj" name="formatCnpj">
    <xsl:param name="cnpj"/>
    <xsl:if test="string-length($cnpj) != 0">
      <xsl:value-of select="concat(substring($cnpj,1,2), '.', substring($cnpj,3,3), '.', substring($cnpj,6,3), '/', substring($cnpj,9,4), '-', substring($cnpj,13,2))"/>
    </xsl:if>
  </xsl:template>

  <!-- formatação de CPF (###.###.###/##)-->
  <xsl:template match="formatCpf" name="formatCpf">
    <xsl:param name="cpf"/>
    <xsl:if test="string-length($cpf) != 0">
      <xsl:value-of select="concat(substring($cpf, 1, 3), '.', substring($cpf, 4, 3), '.', substring($cpf, 7, 3), '-', substring($cpf, 10, 2))"/>
    </xsl:if>
  </xsl:template>

  <!-- formatação de CEP (#####-###)-->
  <xsl:template match="formatCep" name="formatCep">
    <xsl:param name="cep"/>
    <xsl:if test="string-length($cep) != 0">
      <xsl:value-of select="concat(substring($cep,1,5), '-', substring($cep,6,3))"/>
    </xsl:if>
  </xsl:template>


  <!-- retorna um campo específico da chave de acesso -->
  <xsl:template match="retornaCampoChaveAcesso" name="retornaCampoChaveAcesso">
    <xsl:param name="chave_acesso"/>
    <xsl:param name="campo"/>
    
    <xsl:if test="string-length($chave_acesso) = 44">
      <xsl:choose>
        <xsl:when test="$campo = 'cUF'">
          <xsl:value-of select="substring($chave_acesso,1,2)"/>
        </xsl:when>
        <xsl:when test="$campo = 'AA'">
          <xsl:value-of select="substring($chave_acesso,3,2)"/>
        </xsl:when>
        <xsl:when test="$campo = 'MM'">
          <xsl:value-of select="substring($chave_acesso,5,2)"/>
        </xsl:when>
        <xsl:when test="$campo = 'AAMM'">
          <xsl:value-of select="substring($chave_acesso,3,4)"/>
        </xsl:when>
        <xsl:when test="$campo = 'CNPJ'">
          <xsl:value-of select="substring($chave_acesso,7, 14)"/>
        </xsl:when>
        <xsl:when test="$campo = 'modelo'">
          <xsl:value-of select="substring($chave_acesso,21,2)"/>
        </xsl:when>
        <xsl:when test="$campo = 'serie'">
          <xsl:value-of select="substring($chave_acesso,23, 3)"/>
        </xsl:when>
        <xsl:when test="$campo = 'nNF'">
          <xsl:value-of select="substring($chave_acesso,26, 9)"/>
        </xsl:when>
        <xsl:when test="$campo = 'tpEmis'">
          <xsl:value-of select="substring($chave_acesso,35, 1)"/>
        </xsl:when>
        <xsl:when test="$campo = 'cDV'">
          <xsl:value-of select="substring($chave_acesso,44, 1)"/>
        </xsl:when>
      </xsl:choose>
    </xsl:if>
  </xsl:template>

 
  <!-- retorna nome do mes baseado em codigo de dois digitos -->
  <xsl:template match="nomeMes" name="nomeMes">
    <xsl:param name="mes"/>
    <xsl:if test="$mes='01'">Jan</xsl:if>
    <xsl:if test="$mes='02'">Fev</xsl:if>
    <xsl:if test="$mes='03'">Mar</xsl:if>
    <xsl:if test="$mes='04'">Abr</xsl:if>
    <xsl:if test="$mes='05'">Mai</xsl:if>
    <xsl:if test="$mes='06'">Jun</xsl:if>
    <xsl:if test="$mes='07'">Jul</xsl:if>
    <xsl:if test="$mes='08'">Ago</xsl:if>
    <xsl:if test="$mes='09'">Set</xsl:if>
    <xsl:if test="$mes='10'">Out</xsl:if>
    <xsl:if test="$mes='11'">Nov</xsl:if>
    <xsl:if test="$mes='12'">Dez</xsl:if>
  </xsl:template>

  <!-- retorna nome do estado baseado no codigo da UF(#####-###)-->
  <xsl:template match="nomeUF" name="nomeUF">
    <xsl:param name="uf"/>
    <xsl:if test="$uf=11">RONDÔNIA</xsl:if>
    <xsl:if test="$uf=12">ACRE</xsl:if>
    <xsl:if test="$uf=13">AMAZONAS</xsl:if>
    <xsl:if test="$uf=14">RORAIMA</xsl:if>
    <xsl:if test="$uf=15">PARÁ</xsl:if>
    <xsl:if test="$uf=16">AMAPÁ</xsl:if>
    <xsl:if test="$uf=17">TOCANTINS</xsl:if>
    <xsl:if test="$uf=21">MARANHÃO</xsl:if>
    <xsl:if test="$uf=22">PIAUÍ</xsl:if>
    <xsl:if test="$uf=23">CEARÁ</xsl:if>
    <xsl:if test="$uf=24">RIO GRANDE DO NORTE</xsl:if>
    <xsl:if test="$uf=25">PARAÍBA</xsl:if>
    <xsl:if test="$uf=26">PERNAMBUCO</xsl:if>
    <xsl:if test="$uf=27">ALAGOAS</xsl:if>
    <xsl:if test="$uf=28">SERGIPE</xsl:if>
    <xsl:if test="$uf=29">BAHIA</xsl:if>
    <xsl:if test="$uf=31">MINAS GERAIS</xsl:if>
    <xsl:if test="$uf=32">ESPIRITO SANTO</xsl:if>
    <xsl:if test="$uf=33">RIO DE JANEIRO</xsl:if>
    <xsl:if test="$uf=35">SÃO PAULO</xsl:if>
    <xsl:if test="$uf=41">PARANÁ</xsl:if>
    <xsl:if test="$uf=42">SANTA CATARINA</xsl:if>
    <xsl:if test="$uf=43">RIO GRANDE DO SUL</xsl:if>
    <xsl:if test="$uf=50">MATO GROSSO DO SUL</xsl:if>
    <xsl:if test="$uf=51">MATO GROSSO</xsl:if>
    <xsl:if test="$uf=52">GOIÁS</xsl:if>
    <xsl:if test="$uf=53">DISTRITO FEDERAL</xsl:if>
    <xsl:if test="$uf=91">AMBIENTE NACIONAL</xsl:if>
    <xsl:if test="$uf=0 or $uf=99">EXTERIOR</xsl:if>
  </xsl:template>

  <!-- retorna SIGLA do estado baseado no codigo da UF(#####-###)-->
  <xsl:template match="siglaUF" name="siglaUF">
    <xsl:param name="uf"/>
    <xsl:if test="$uf=11">RO</xsl:if>
    <xsl:if test="$uf=12">AC</xsl:if>
    <xsl:if test="$uf=13">AM</xsl:if>
    <xsl:if test="$uf=14">RR</xsl:if>
    <xsl:if test="$uf=15">PA</xsl:if>
    <xsl:if test="$uf=16">AP</xsl:if>
    <xsl:if test="$uf=17">TO</xsl:if>
    <xsl:if test="$uf=21">MA</xsl:if>
    <xsl:if test="$uf=22">PI</xsl:if>
    <xsl:if test="$uf=23">CE</xsl:if>
    <xsl:if test="$uf=24">RN</xsl:if>
    <xsl:if test="$uf=25">PB</xsl:if>
    <xsl:if test="$uf=26">PE</xsl:if>
    <xsl:if test="$uf=27">AL</xsl:if>
    <xsl:if test="$uf=28">SE</xsl:if>
    <xsl:if test="$uf=29">BA</xsl:if>
    <xsl:if test="$uf=31">MG</xsl:if>
    <xsl:if test="$uf=32">ES</xsl:if>
    <xsl:if test="$uf=33">RJ</xsl:if>
    <xsl:if test="$uf=35">SP</xsl:if>
    <xsl:if test="$uf=41">PR</xsl:if>
    <xsl:if test="$uf=42">SC</xsl:if>
    <xsl:if test="$uf=43">RS</xsl:if>
    <xsl:if test="$uf=50">MS</xsl:if>
    <xsl:if test="$uf=51">MT</xsl:if>
    <xsl:if test="$uf=52">GO</xsl:if>
    <xsl:if test="$uf=53">DF</xsl:if>
    <xsl:if test="$uf=0 or $uf=99">EX</xsl:if>
  </xsl:template>

  <!-- retorna CODIGO IBGE do estado baseado na sigla da UF-->
  <xsl:template match="codigoUF" name="codigoUF">
    <xsl:param name="uf"/>
    <xsl:if test="$uf='RO'">11</xsl:if>
    <xsl:if test="$uf='AC'">12</xsl:if>
    <xsl:if test="$uf='AM'">13</xsl:if>
    <xsl:if test="$uf='RR'">14</xsl:if>
    <xsl:if test="$uf='PA'">15</xsl:if>
    <xsl:if test="$uf='AP'">16</xsl:if>
    <xsl:if test="$uf='TO'">17</xsl:if>
    <xsl:if test="$uf='MA'">21</xsl:if>
    <xsl:if test="$uf='PI'">22</xsl:if>
    <xsl:if test="$uf='CE'">23</xsl:if>
    <xsl:if test="$uf='RN'">24</xsl:if>
    <xsl:if test="$uf='PB'">25</xsl:if>
    <xsl:if test="$uf='PE'">26</xsl:if>
    <xsl:if test="$uf='AL'">27</xsl:if>
    <xsl:if test="$uf='SE'">28</xsl:if>
    <xsl:if test="$uf='BA'">29</xsl:if>
    <xsl:if test="$uf='MG'">31</xsl:if>
    <xsl:if test="$uf='ES'">32</xsl:if>
    <xsl:if test="$uf='RJ'">33</xsl:if>
    <xsl:if test="$uf='SP'">35</xsl:if>
    <xsl:if test="$uf='PR'">41</xsl:if>
    <xsl:if test="$uf='SC'">42</xsl:if>
    <xsl:if test="$uf='RS'">43</xsl:if>
    <xsl:if test="$uf='MS'">50</xsl:if>
    <xsl:if test="$uf='MT'">51</xsl:if>
    <xsl:if test="$uf='GO'">52</xsl:if>
    <xsl:if test="$uf='DF'">53</xsl:if>
  </xsl:template>

  <!-- formatação de números de telefone "####-####", "(##)####-####" e "(###)####-####" -->
  <xsl:template match="formatFone" name="formatFone">
    <xsl:param name="fone"/>
    <xsl:if test="string-length($fone) != 0">
      <xsl:choose>
        <xsl:when test="string-length($fone) = 8">
          <xsl:value-of select="concat(substring($fone,1,4), '-', substring($fone,5,4))"/>
        </xsl:when>
        <xsl:when test="string-length($fone) = 10">
          <xsl:value-of select="concat('(', substring($fone,1,2),')', substring($fone,3,4), '-', substring($fone,7,4))"/>
        </xsl:when>
        <xsl:when test="string-length($fone) = 11">
          <xsl:value-of select="concat('(', substring($fone,1,3),')', substring($fone,4,4), '-', substring($fone,8,4))"/>
        </xsl:when>
        <xsl:otherwise>
          <xsl:value-of select="$fone"/>
        </xsl:otherwise>
      </xsl:choose>
    </xsl:if>
  </xsl:template>

  <!-- formatação da chave de acesso da NFe -->
  <xsl:template match="formatNfe" name="formatNfe">
    <xsl:param name="nfe"/>
    <xsl:if test="string-length($nfe) != 0">
      <xsl:value-of select="concat(
									substring($nfe,1,4), ' ',
									substring($nfe,5,4), ' ',
									substring($nfe,9,4), ' ',
									substring($nfe,13,4), ' ',
									substring($nfe,17,4), ' ',
									substring($nfe,21,4), ' ',
									substring($nfe,25,4), ' ',
									substring($nfe,29,4), ' ',
									substring($nfe,33,4), ' ',
									substring($nfe,37,4), ' ',
									substring($nfe,41,4))"/>
    </xsl:if>
  </xsl:template>
  <xsl:template match="quebraLinha" name="quebraLinha">
    <xsl:param name="infCpl"/>
    <div style="word-wrap: break-word">
      <xsl:value-of select="$infCpl"/>
    </div>
  </xsl:template>

  <xsl:template match="HEADER_IMPRESSAO" name="HEADER_IMPRESSAO">
    <xsl:param name="titulo"/>
    <!-- titulo que aparece no header (opcional) -->
    <xsl:param name="tipo"/>
    <!-- "simples", "completo" -->
    <div id="print-header">
      <div id="print-header-logo"></div>
      <div id="print-header-titulo">
        Governo do Estado do Rio Grande do Sul<br />
        Secretaria da Fazenda
      </div>
    </div>
    <div id="print-nfe-info">
      <xsl:if test="$titulo != ''" >
        <span id="print-header-titulo-param">
          <xsl:value-of select="$titulo"/>
        </span>
      </xsl:if>
      <xsl:choose>
        <xsl:when test="$tipo = 'simples'">
          <table>
            <tr>
              <td>
                <b>Chave&#160;de&#160;Acesso:</b>
              </td>
              <td>
                <xsl:call-template name="formatNfe">
                  <xsl:with-param name="nfe" select="substring-after(//n:infNFe/@Id,'NFe')"/>
                </xsl:call-template>
              </td>
              <td>
                <b>Número&#160;NF-e:</b>
              </td>
              <td>
                <xsl:value-of select = "//n:infNFe/n:ide/n:nNF"/>
              </td>
            </tr>
            <tr>
              <td>
                <b>Data&#160;de&#160;Emissão:</b>
              </td>
              <td>
                <xsl:variable name="dEmi" select="//n:infNFe/n:ide/n:dEmi"/>
                <xsl:call-template name="formatDate">
                  <xsl:with-param name="date" select="$dEmi"/>
                </xsl:call-template>
              </td>
            </tr>
          </table>
        </xsl:when>
        <xsl:when test="$tipo = 'completo'">
          <table>
            <tr>
              <td>
                <b>Chave de Acesso:</b>
              </td>
              <td>
                <xsl:call-template name="formatNfe">
                  <xsl:with-param name="nfe" select="substring-after(//n:infNFe/@Id,'NFe')"/>
                </xsl:call-template>
              </td>
            </tr>
            <tr>
              <td>
                <b>Número NF-e:</b>
              </td>
              <td>
                <xsl:value-of select = "//n:infNFe/n:ide/n:nNF"/>
              </td>
              <td>
                <b>Série:</b>
              </td>
              <td>
                <xsl:value-of select="//n:infNFe/n:ide/n:serie"/>
              </td>
            </tr>
            <tr>
              <td>
                <b>Data de Emissão:</b>
              </td>
              <td>
                <xsl:variable name="dEmi" select="//n:infNFe/n:ide/n:dEmi"/>
                <xsl:call-template name="formatDate">
                  <xsl:with-param name="date" select="$dEmi"/>
                </xsl:call-template>
              </td>
              <td>
                <b>Número do Protocolo:</b>
              </td>
              <td>
                <xsl:value-of select="//n:infProt/n:nProt"/>
              </td>
            </tr>
          </table>
        </xsl:when>
        <xsl:otherwise >
          <hr />
        </xsl:otherwise>
      </xsl:choose>
    </div>
  </xsl:template>

</xsl:stylesheet>