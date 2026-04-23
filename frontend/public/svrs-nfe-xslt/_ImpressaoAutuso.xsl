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
	<xsl:template match="/">
		<xsl:call-template name="HEADER_IMPRESSAO">
			<xsl:with-param name="titulo">AUTORIZAÇÃO DE USO</xsl:with-param>
			<xsl:with-param name="tipo">completo</xsl:with-param>
		</xsl:call-template>
		<div class="tbl-aut-uso">
			<!-- ** Emitente ** -->
			<span class="subtitle">Emitente</span>
			<table>
				<tr>
					<td class="itm">Nome/Razão Social:</td>
					<td class="val" colspan="2">
						<xsl:value-of select = "//n:infNFe/n:emit/n:xNome"/>
					</td>
				</tr>
				<tr>
					<td class="itm">CNPJ / CPF:</td>
					<td class="val" colspan="2">
						<xsl:variable name="cnpj" select="//n:infNFe/n:emit/n:CNPJ"/>
						<xsl:variable name="cpf" select="//n:infNFe/n:emit/n:CPF"/>
						<xsl:call-template name="formatCnpj">
							<xsl:with-param name="cnpj" select="$cnpj"/>
						</xsl:call-template>
						<xsl:call-template name="formatCpf">
							<xsl:with-param name="cpf" select="$cpf"/>
						</xsl:call-template>
					</td>
				</tr>
				<tr>
					<td class="itm">UF:</td>
					<td class="val" colspan="2">
						<xsl:value-of select = "//n:infNFe/n:emit/n:enderEmit/n:UF"/>
					</td>
				</tr>
			</table>
			<!-- ** Destinatário ** -->
			<span class="subtitle">Destinatário</span>
			<table>
				<tr>
					<td class="itm">Nome/Razão Social:</td>
					<td class="val" colspan="2">
						<xsl:value-of select = "//n:infNFe/n:dest/n:xNome"/>
					</td>
				</tr>
				<tr>
					<td class="itm">CNPJ / CPF:</td>
					<td class="val" colspan="2">
						<xsl:variable name="cnpj" select="//n:infNFe/n:dest/n:CNPJ"/>
						<xsl:variable name="cpf" select="//n:infNFe/n:dest/n:CPF"/>
						<xsl:call-template name="formatCnpj">
							<xsl:with-param name="cnpj" select="$cnpj"/>
						</xsl:call-template>
						<xsl:call-template name="formatCpf">
							<xsl:with-param name="cpf" select="$cpf"/>
						</xsl:call-template>
					</td>
				</tr>
				<tr>
					<td class="itm">UF:</td>
					<td class="val" colspan="2">
						<xsl:value-of select = "//n:infNFe/n:dest/n:enderDest/n:UF"/>
					</td>
				</tr>
			</table>
			<!-- ** Autorização ** -->
			<span class="subtitle">Data e Hora da Autorização de Uso</span>
			<table>
				<tr>
					<td class="itm">Autorização:</td>
					<td class="val">
						<xsl:call-template name="formatDateTime">
							<xsl:with-param name="dateTime" select="//n:infProt/n:dhRecbto"/>
							<xsl:with-param name="include_as" select="1"/>
						</xsl:call-template>
					</td>
				</tr>
			</table>
			<!-- ** Hash Code da NF-e ** -->
			<span class="subtitle">Hash Code da NF-e:</span>
			<table>
				<tr>
					<td class="val">
						<xsl:value-of select="//s:Signature/s:SignedInfo/s:Reference/s:DigestValue"/>
					</td>
				</tr>
			</table>
			<!-- valores do ICMS e total da NF-e -->
			<table>
				<tr>
					<td class="itm">Base de cálculo do ICMS:</td>
					<td class="itm">Valor do ICMS:</td>
					<td class="itm">Valor total da NF-e:</td>
				</tr>
				<tr>
					<td class="val">
						<xsl:value-of select = "format-number(//n:infNFe/n:total/n:ICMSTot/n:vBC,'##.##.##0,00')"/>
					</td>
					<td class="val">
						<xsl:call-template name="format2Casas">
							<xsl:with-param name="num" select="//n:infNFe/n:total/n:ICMSTot/n:vICMS"/>
						</xsl:call-template>
					</td>
					<td class="val">
						<xsl:call-template name="format2Casas">
							<xsl:with-param name="num" select="//n:infNFe/n:total/n:ICMSTot/n:vNF"/>
						</xsl:call-template>
					</td>
				</tr>
			</table>
		</div>
		<br />
	</xsl:template>
</xsl:stylesheet>