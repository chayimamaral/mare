<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet version="1.0"
      xmlns:xsl="http://www.w3.org/1999/XSL/Transform"  
      xmlns:n="http://www.portalfiscal.inf.br/nfe"
      xmlns:date="http://exslt.org/formatoData"
      xmlns:chave="http://exslt.org/chaveacesso" 
      xmlns:r="http://www.serpro.gov.br/nfe/remessanfe.xsd"
      exclude-result-prefixes="date" >
	<xsl:decimal-format decimal-separator="," grouping-separator="." />
	<xsl:template match="/">
		<script language="javascript">
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

		<table width="98%"  align="center">
			<tr>
				<td height="25" class="TituloAreaRestritacentro">
					<B class="textoVerdana8bold">
						<STRONG class="textoVerdana8bold">Dados dos Produtos e Serviços</STRONG>
					</B>
				</td>
			</tr>
		</table>

		<table width="98%"  align="center">
			<tr>
				<td class="TituloAreaRestrita" valign ="top" style="width: 04%;">
					<span class="TextoFundoBrancoNegrito">
						Num.<br />
					</span>
				</td>
				<td class="TituloAreaRestrita" valign ="top" style="width: 56%;">
					<span class="TextoFundoBrancoNegrito">
						Descrição<br />
					</span>
				</td>
				<td class="TituloAreaRestrita" valign ="top" style="width: 10%;">
					<span class="TextoFundoBrancoNegrito">
						Qtd.<br />
					</span>
				</td>
				<td class="TituloAreaRestrita" valign ="top" style="width:10%;">
					<span class="TextoFundoBrancoNegrito">
						Unidade Comercial<br />
					</span>
				</td>
				<td class="TituloAreaRestrita" valign ="top" style="width: 20%;">
					<span class="TextoFundoBrancoNegrito">
						Valor(R$)<br />
					</span>
				</td>
			</tr>
		</table>
		<xsl:for-each select="//n:infNFe/n:det">
			<xsl:variable name="chaves" select="@nItem"/>
			<xsl:variable name="chaves2" select="concat('A',' + ',$chaves)"/>
			<table align="center"  class="textoVerdana7" width="98%" id="{$chaves2}">
				<tr>
					
					<td valign ="top" style="width: 03%;">
						<span class="linha">
							<xsl:value-of select="position()" />
						</span>
					</td>
					<xsl:for-each select ="n:prod/n:xProd">
						<td valign ="top" style="width: 56%;">
							<span class="linha">
								<xsl:value-of select="text()" />
							</span>
						</td>
					</xsl:for-each>
					<xsl:for-each select ="n:prod/n:qCom">
						<td valign ="top" style="width: 10%;">
							<span class="linha">
								<xsl:value-of select="format-number(text(),'##.##.##0,000')" />
							</span>
						</td>
					</xsl:for-each>
				
					<td valign ="top" style="width:10%;">
						<span class="linha">
							<xsl:value-of select="n:prod/n:uCom" />
						</span>
					</td>
					<td valign ="top" style="width: 20%;">
						<span class="linha">
							<xsl:value-of select="format-number(n:prod/n:vProd,'##.##.##0,00')" />
						</span>
					</td>
				</tr>
			</table>
			<xsl:variable name="tab" select="@nItem"/>
			<table align="center"  class="textoVerdana7" width="98%" >
				<tr>
					<td>
						<table  width="98%"  align="center"  >
							<tr class="TextoFundoBranco">
								<xsl:for-each select ="n:prod/n:cProd">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Código do Produto <br />
										</span>
										<span class="linha">
											<xsl:value-of select="text()" />
										</span>
									</td>
								</xsl:for-each>
								<!--NCM-->
								<xsl:for-each select ="n:prod/n:NCM">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Código NCM<br />
										</span>
										<span class="linha">
											<xsl:value-of select="text()" />
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:EXTIPI">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Código EX da TIPI <br />
										</span>
										<span class="linha">
											<xsl:value-of select="text()" />
										</span>
									</td>
								</xsl:for-each>
								<!--CFOP-->
								<td valign ="top" style="height: 37px; width: 33%;">
									<span class="TextoFundoBrancoNegrito">
										CFOP<br />
									</span>
									<span class="linha">
										<xsl:value-of select="n:prod/n:CFOP" />
									</span>
								</td>
							</tr>



							<tr class="TextoFundoBranco">
								<xsl:for-each select ="n:prod/n:vDesc">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Valor do Desconto<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:vFrete">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Valor do Frete <br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:vSeg">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Valor do Seguro<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
										</span>
									</td>
								</xsl:for-each>
							</tr>

							<tr class="TextoFundoBranco">
								<xsl:for-each select ="n:prod/n:vOutro">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Outras despesas acessórias<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:indTot">
									<td colspan="2" valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
                      Indicador de Composição do Valor Total da NF-e <br />
                    </span>
										<span class="linha">
											<xsl:variable name="indTot" select ="text()" />
											<xsl:if test="$indTot='0'">
                        0 – O valor do item (vProd) não compõe o valor total da NF-e (vProd)
                      </xsl:if>
											<xsl:if test="$indTot='1'">
                        1  – O valor do item (vProd) compõe o valor total da NF-e (vProd)
                      </xsl:if>
										</span>
									</td>
								</xsl:for-each>
							</tr>

							<tr class="TextoFundoBranco">
								<xsl:for-each select ="n:prod/n:cEAN">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Código EAN Comercial<br />
										</span>
										<span class="linha">
											<xsl:value-of select="text()"/>
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:uCom">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Unidade Comercial<br />
										</span>
										<span class="linha">
											<xsl:value-of select="text()"/>
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:qCom">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Quantidade Comercial<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(text(),'##.##.##0,0000')" />
										</span>
									</td>
								</xsl:for-each>
							</tr>
							<tr class="TextoFundoBranco">

								<xsl:for-each select ="n:prod/n:cEANTrib">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Código EAN Tributável<br />
										</span>
										<span class="linha">
											<xsl:value-of select="text()"/>
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:uTrib">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Unidade Tributável<br />
										</span>
										<span class="linha">
											<xsl:value-of select="text()"/>
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:qTrib">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Quantidade Tributável<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(text(),'##.##.##0,0000')" />
										</span>
									</td>
								</xsl:for-each>
							</tr>
							<tr class="TextoFundoBranco">
								<xsl:for-each select ="n:prod/n:vUnCom">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Valor unitário de comercialização<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(text(),'##.##.##0,0000')" />
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:vUnTrib">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Valor unitário de tributação<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(text(),'##.##.##0,0000')" />
										</span>
									</td>
								</xsl:for-each>
								<td></td>
							</tr>
							<tr class="TextoFundoBranco">
								<xsl:for-each select ="n:prod/n:xPed">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Número do Pedido de Compra<br />
										</span>
										<span class="linha">
											<xsl:value-of select="text()" />
										</span>
									</td>
								</xsl:for-each>
								<xsl:for-each select ="n:prod/n:nItemPed">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Item do Pedido de Compra<br />
										</span>
										<span class="linha">
											<xsl:value-of select="text()" />
										</span>
									</td>
								</xsl:for-each>
								<td></td>
							</tr>
						</table>

						<!--ICMS DA OPERAÇÃO PRÓPRIA ICMS 00 e ICMS 10-->
						<xsl:variable name="icms" select="n:imposto/n:ICMS"/>
						<xsl:if test="$icms!=''">
							<table align="center"  width="98%">
								<tr class="TextoFundoBranco">
									<td  class="TituloAreaRestrita2">
										ICMS NORMAL e ST
									</td>
								</tr>
							</table>

							<!--ICMS 00-->
							<xsl:variable name="icms00" select="n:imposto/n:ICMS/n:ICMS00"/>
							<xsl:if test="$icms00!=''">
								<table align="center"  class="textoVerdana7" width="98%">
									<tr class="TextoFundoBranco" >
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS00/n:orig" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Origem da Mercadoria<br />
												</span>
												<span class="linha">
													<xsl:variable name="origmerc" select="text()"/>
													<xsl:if test="$origmerc='0'">
														0–Nacional
													</xsl:if>
													<xsl:if test="$origmerc='1'">
														1–Estrangeira – Importação direta
													</xsl:if>
													<xsl:if test="$origmerc='2'">
														2–Estrangeira – Adquirida no mercado interno.
													</xsl:if>
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS00/n:CST" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Tributação do ICMS<br />
												</span>
												<span class="linha">
													<xsl:variable name="origmerc" select="text()"/>
													<xsl:if test="$origmerc='00'">
														00 – Tributada integralmente
													</xsl:if>
													<xsl:if test="$origmerc='10'">
														10 - Tributada e com cobrança do ICMS por substituição tributária
													</xsl:if>
													<xsl:if test="$origmerc='20'">
														20 - Com redução de base de cálculo
													</xsl:if>
													<xsl:if test="$origmerc='30'">
														30 - Isenta ou não tributada e com cobrança do ICMS por substituição tributária
													</xsl:if>
													<xsl:if test="$origmerc='40'">
														40 - Isenta
													</xsl:if>
													<xsl:if test="$origmerc='41'">
														41 - Não tributada
													</xsl:if>
													<xsl:if test="$origmerc='50'">
														50 - Suspensão
													</xsl:if>
													<xsl:if test="$origmerc='51'">
														51 - Diferimento
													</xsl:if>
													<xsl:if test="$origmerc='60'">
														60 - ICMS cobrado anteriormente por substituição tributária
													</xsl:if>
													<xsl:if test="$origmerc='70'">
														70 - Com redução de base de cálculo e cobrança do ICMS por substituição tributária
													</xsl:if>
													<xsl:if test="$origmerc='90'">
														90 - Outros
													</xsl:if>
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS00/n:modBC" >
											<td valign ="top" style="width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Modalidade Definição da BC ICMS NORMAL<br />
												</span>
												<span class="linha">
													<xsl:variable name="modBcop" select="text()"/>
													<xsl:if test="$modBcop='0'">
														0 - Margem Valor Agregado (%)
													</xsl:if>
													<xsl:if test="$modBcop='1'">
														1 - Pauta (Valor)
													</xsl:if>
													<xsl:if test="$modBcop='2'">
														2 - Preço Tabelado Máx. (valor)
													</xsl:if>
													<xsl:if test="$modBcop='3'">
														3 - Valor da operação
													</xsl:if>
												</span>
											</td>
										</xsl:for-each>
									</tr>
									<tr class="TextoFundoBranco" >
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS00/n:vBC" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Base de Cálculo do ICMS Normal<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS00/n:pICMS" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Alíquota do ICMS Normal <br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS00/n:vICMS" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Valor do ICMS Normal<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
									</tr>
								</table>
							</xsl:if>

							<!--ICMS DA SUBSTITUIÇÃO TRIBUTÁRIA  ICMS 10-->
							<xsl:variable name="icms10" select="n:imposto/n:ICMS/n:ICMS10"/>
							<xsl:if test="$icms10!=''">
								<table align="center"  class="textoVerdana7" width="98%">
									<tr class="TextoFundoBranco" >
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:orig" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Origem da Mercadoria<br />
												</span>
												<span class="linha">
													<xsl:variable name="origmercst" select="text()"/>
													<xsl:if test="$origmercst='0'">
														0–Nacional
													</xsl:if>
													<xsl:if test="$origmercst='1'">
														1–Estrangeira – Importação direta
													</xsl:if>
													<xsl:if test="$origmercst='2'">
														2–Estrangeira – Adquirida no mercado interno.
													</xsl:if>
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:CST" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Tributação do ICMS <br />
												</span>
												<span class="linha">
													<xsl:variable name="origmerc" select="text()"/>
													<xsl:if test="$origmerc='10'">
														10 - Tributada e com cobrança do ICMS por substituição tributária
													</xsl:if>
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:modBC" >
											<td valign ="top" style="width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Modalidade Definição da BC ICMS NORMAL <br />
												</span>
												<span class="linha">
													<xsl:variable name="modBc01" select="text()"/>
													<xsl:if test="$modBc01='0'">
														0 - Margem Valor Agregado (%)
													</xsl:if>
													<xsl:if test="$modBc01='1'">
														1 - Pauta (Valor)
													</xsl:if>
													<xsl:if test="$modBc01='2'">
														2 - Preço Tabelado Máx. (valor)
													</xsl:if>
													<xsl:if test="$modBc01='3'">
														3 - valor da operação
													</xsl:if>
												</span>
											</td>
										</xsl:for-each>
									</tr>
									<tr class="TextoFundoBranco">
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:vBC" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Base de Cálculo do ICMS Normal<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:pICMS" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Alíquota ICMS Normal <br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:vICMS" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Valor do ICMS Normal<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
									</tr>
									<tr class="TextoFundoBranco">
									</tr>
									<tr class="TextoFundoBranco">
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:vBCST" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Base de Cálculo do ICMS ST<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:pICMSST" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Alíquota do ICMS ST<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:vICMSST" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Valor do ICMS ST <br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
									</tr>
									<tr class="TextoFundoBranco">
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:pRedBCST" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Percentual Redução de BC do ICMS ST <br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:pMVAST" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Percentual do MVA do ICMS ST<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS10/n:modBCST" >
											<td valign ="top" style="width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Modalidade Definição da BC ICMS ST <br />
												</span>
												<span class="linha">
													<xsl:variable name="modBcst" select="text()"/>
													<xsl:if test="$modBcst='0'">
														0 – Preço tabelado ou máx. sugerido
													</xsl:if>
													<xsl:if test="$modBcst='1'">
														1 - Lista Negativa (valor)
													</xsl:if>
													<xsl:if test="$modBcst='2'">
														2 - Lista Positiva (valor)
													</xsl:if>
													<xsl:if test="$modBcst='3'">
														3 - Lista Neutra (valor)
													</xsl:if>
													<xsl:if test="$modBcst='4'">
														4 - Margem Valor Agregado (%)
													</xsl:if>
													<xsl:if test="$modBcst='5'">
														5 - Pauta (valor)
													</xsl:if>
												</span>
											</td>
										</xsl:for-each>
									</tr>
								</table>
							</xsl:if>
						</xsl:if>

						<!--ICMS20 DA OPERAÇÃO PRÓPRIA Com redução de base de cálculo-->
						<xsl:variable name="icms20" select="n:imposto/n:ICMS/n:ICMS20"/>
						<xsl:if test="$icms20!=''">

							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS20/orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmerc='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmerc='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS20/n:CST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Tributação do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='20'">
													20 - Com redução de base de cálculo
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS20/n:modBC" >
										<td valign ="top" style="width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade Definição da BC do ICMS <br />
											</span>
											<span class="linha">
												<xsl:variable name="modBcop" select="text()"/>
												<xsl:if test="$modBcop='0'">
													0 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBcop='1'">
													1 - Pauta (Valor)
												</xsl:if>
												<xsl:if test="$modBcop='2'">
													2 - Preço Tabelado Máx. (valor)
												</xsl:if>
												<xsl:if test="$modBcop='3'">
													3 - Valor da operação
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Base de Cálculo<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(n:imposto/n:ICMS/n:ICMS20/n:vBC,'##.##.##0,00')" />
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Alíquota<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(n:imposto/n:ICMS/n:ICMS20/n:pICMS,'##.##.##0,00')" />
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Valor<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(n:imposto/n:ICMS/n:ICMS20/n:vICMS,'##.##.##0,00')" />
										</span>
									</td>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS20/n:pRedBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual Redução de BC do ICMS Normal <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!--ICMS30 DA SUBSTITUIÇÃO TRIBUTÁRIA-->
						<xsl:variable name="icmst30" select="n:imposto/n:ICMS/n:ICMS30"/>
						<xsl:if test="$icmst30!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS30/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS30/n:CST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Tributação do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='30'">
													30 - Isenta ou não tributada e com cobrança do ICMS por substituição tributária
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS30/n:modBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade de Determinação da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:variable name="modBcst" select="text()"/>
												<xsl:if test="$modBcst='0'">
													0 – Preço tabelado ou máx. sugerido
												</xsl:if>
												<xsl:if test="$modBcst='1'">
													1 - Lista Negativa (valor)
												</xsl:if>
												<xsl:if test="$modBcst='2'">
													2 - Lista Positiva (valor)
												</xsl:if>
												<xsl:if test="$modBcst='3'">
													3 - Lista Neutra (valor)
												</xsl:if>
												<xsl:if test="$modBcst='4'">
													4 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBcst='5'">
													5 - Pauta (valor)
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS30/n:pRedBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Redução de BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS30/n:pMVAST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Margem de Valor Adicionado do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS30/n:vBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS30/n:pICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do Imposto do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS30/n:vICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!-- ICMS40 DA SUBSTITUIÇÃO TRIBUTÁRIA CST – 40 - Isenta 41 - Não tributada 50 - Suspensão 51 - Diferimento -->
						<xsl:variable name="icmst40" select="n:imposto/n:ICMS/n:ICMS40"/>
						<xsl:if test="$icmst40!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS40/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria <br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS40/n:CST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Tributação do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='40'">
													40 - Isenta
												</xsl:if>
												<xsl:if test="$origmerc='41'">
													41 - Não tributada
												</xsl:if>
												<xsl:if test="$origmerc='50'">
													50 - Suspensão
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS40/n:vICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS40/n:motDesICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Motivo da desoneração do ICMS <br />
											</span>
											<span class="linha">
												<xsl:variable name="motDesICMS" select="text()"/>
												<xsl:if test="$motDesICMS='1'">
													1 – Táxi
												</xsl:if>
												<xsl:if test="$motDesICMS='2'">
													2 – Deficiente Físico
												</xsl:if>
												<xsl:if test="$motDesICMS='3'">
													3 – Produtor Agropecuário
												</xsl:if>
												<xsl:if test="$motDesICMS='4'">
													4 – Frotista/Locadora
												</xsl:if>
												<xsl:if test="$motDesICMS='5'">
													5 – Diplomático/Consular
												</xsl:if>
												<xsl:if test="$motDesICMS='6'">
													6 – Utilitários e Motocicletas da  Amazônia Ocidental e Áreas de Livre Comércio
												</xsl:if>
												<xsl:if test="$motDesICMS='7'">
													7 – SUFRAMA
												</xsl:if>
												<xsl:if test="$motDesICMS='9'">
													9 – outros
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!-- ICMS51 DA SUBSTITUIÇÃO TRIBUTÁRIA CST – 51-->
						<xsl:variable name="icmst51" select="n:imposto/n:ICMS/n:ICMS51"/>
						<xsl:if test="$icmst51!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS51/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria <br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS51/n:CST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Tributação do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='51'">
													51 - Diferimento
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS51/n:modBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade de determinação da BC do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="modBC" select ="n:imposto/n:ICMS/n:ICMS51/n:modBC" />
												<xsl:if test="$modBC='0'">
													0 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBC='1'">
													1 - Pauta (Valor)
												</xsl:if>
												<xsl:if test="$modBC='2'">
													2 - Preço Tabelado Máx. (valor)
												</xsl:if>
												<xsl:if test="$modBC='3'">
													3 - valor da operação
												</xsl:if>

											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS51/n:pRedBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Redução de BC<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS51/n:vICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>

								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS51/n:pICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do imposto<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS51/n:vICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>

							</table>
						</xsl:if>

						<!--ICMS60 DA SUBSTITUIÇÃO TRIBUTÁRIA-->
						<xsl:variable name="icmst60" select="n:imposto/n:ICMS/n:ICMS60"/>
						<xsl:if test="$icmst60!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS60/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS60/n:CST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Tributação do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="CST" select="text()"/>
												<xsl:if test="$CST='60'">
													60 - ICMS cobrado anteriormente por substituição tributária
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS60/n:vBCSTRet" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST retido<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS60/n:vICMSSTRet" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST retido<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!--ICMS70 DA SUBSTITUIÇÃO TRIBUTÁRIA Com redução de base de cálculo e cobrança do ICMS por substituição tributária-->
						<xsl:variable name="icmst70" select="n:imposto/n:ICMS/n:ICMS70"/>
						<xsl:if test="$icmst70!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:CST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Tributação do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='70'">
													70 - Com redução de base de cálculo e cobrança do ICMS por substituição tributária ICMS por substituição tributária
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:modBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade<br />
											</span>
											<span class="linha">
												<xsl:variable name="modBc01" select="text()"/>
												<xsl:if test="$modBc01='0'">
													0 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBc01='1'">
													1 - Pauta (Valor)
												</xsl:if>
												<xsl:if test="$modBc01='2'">
													2 - Preço Tabelado Máx. (valor)
												</xsl:if>
												<xsl:if test="$modBc01='3'">
													3 - valor da operação
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:pRedBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual Redução de BC do ICMS Normal <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:vBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Base de Cálculo<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:pICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:vICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:modBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade de Determinação da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:variable name="modBcst" select="text()"/>
												<xsl:if test="$modBcst='0'">
													0 – Preço tabelado ou máx. sugerido
												</xsl:if>
												<xsl:if test="$modBcst='1'">
													1 - Lista Negativa (valor)
												</xsl:if>
												<xsl:if test="$modBcst='2'">
													2 - Lista Positiva (valor)
												</xsl:if>
												<xsl:if test="$modBcst='3'">
													3 - Lista Neutra (valor)
												</xsl:if>
												<xsl:if test="$modBcst='4'">
													4 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBcst='5'">
													5 - Pauta (valor)
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:pRedBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Redução de BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:pMVAST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Margen de Valor Adicionado do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:vBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:pICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do Imposto do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco">
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS70/n:vICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!--ICMS90 DA SUBSTITUIÇÃO TRIBUTÁRIA – Outros-->
						<xsl:variable name="icmst90" select="n:imposto/n:ICMS/n:ICMS90"/>
						<xsl:if test="$icmst90!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:CST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Tributação do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='90'">
													90 - Outros
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:modBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade <br />
											</span>
											<span class="linha">
												<xsl:variable name="modBc09" select="text()"/>
												<xsl:if test="$modBc09='0'">
													0 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBc09='1'">
													1 - Pauta (Valor)
												</xsl:if>
												<xsl:if test="$modBc09='2'">
													2 - Preço Tabelado Máx. (valor)
												</xsl:if>
												<xsl:if test="$modBc09='3'">
													3 - valor da operação
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:pRedBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual Redução de BC do ICMS Normal <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:vBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Base de Cálculo<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:pICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:vICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:modBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade de Determinação da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:variable name="modBcst9" select="text()"/>
												<xsl:if test="$modBcst9='0'">
													0 – Preço tabelado ou máx. sugerido
												</xsl:if>
												<xsl:if test="$modBcst9='1'">
													1 - Lista Negativa (valor)
												</xsl:if>
												<xsl:if test="$modBcst9='2'">
													2 - Lista Positiva (valor)
												</xsl:if>
												<xsl:if test="$modBcst9='3'">
													3 - Lista Neutra (valor)
												</xsl:if>
												<xsl:if test="$modBcst9='4'">
													4 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBcst9='5'">
													5 - Pauta (valor)
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:pRedBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Redução de BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:pMVAST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Margen de Valor Adicionado do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:vBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:pICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do Imposto do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMS90/n:vICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!--ICMSPart -->
						<xsl:variable name="icmsPart" select="n:imposto/n:ICMS/n:ICMSPart"/>
						<xsl:if test="$icmsPart!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:CST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Tributação do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='10'">
													10 - Tributada e com cobrança do ICMS por substituição tributária
												</xsl:if>
												<xsl:if test="$origmerc='90'">
													90 - Outros
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:modBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade <br />
											</span>
											<span class="linha">
												<xsl:variable name="modBcPart" select="text()"/>
												<xsl:if test="$modBcPart='0'">
													0 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBcPart='1'">
													1 - Pauta (Valor)
												</xsl:if>
												<xsl:if test="$modBcPart='2'">
													2 - Preço Tabelado Máx. (valor)
												</xsl:if>
												<xsl:if test="$modBcPart='3'">
													3 - valor da operação
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:pRedBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual Redução de BC <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:vBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:pICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do Imposto<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:vICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:modBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade de Determinação da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:variable name="modBcstPart" select="text()"/>
												<xsl:if test="$modBcstPart='0'">
													0 – Preço tabelado ou máx. sugerido
												</xsl:if>
												<xsl:if test="$modBcstPart='1'">
													1 - Lista Negativa (valor)
												</xsl:if>
												<xsl:if test="$modBcstPart='2'">
													2 - Lista Positiva (valor)
												</xsl:if>
												<xsl:if test="$modBcstPart='3'">
													3 - Lista Neutra (valor)
												</xsl:if>
												<xsl:if test="$modBcstPart='4'">
													4 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBcstPart='5'">
													5 - Pauta (valor)
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:pRedBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Redução de BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:pMVAST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Margen de Valor Adicionado do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:vBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:pICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do Imposto do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:vICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:pBCOp" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da BC operação própria<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSPart/n:UFST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												UF<br />
											</span>
											<span class="linha">
												<xsl:value-of select="text()" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!--ICMSST -->
						<xsl:variable name="icmsst" select="n:imposto/n:ICMS/n:ICMSST"/>
						<xsl:if test="$icmsst!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSST/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSST/n:CST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Tributação do ICMS<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='41'">
													41 – Não Tributado
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSST/n:vBCSTRet" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSST/n:vICMSSTRet" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST retido na UF remetente <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSST/n:vBCSTDest" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST da UF destino<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSST/n:vICMSSTDest" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST da UF destino<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!--ICMSSN101 -->
						<xsl:variable name="icmssn101" select="n:imposto/n:ICMS/n:ICMSSN101"/>
						<xsl:if test="$icmssn101!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN101/n:orig" >
										<td valign ="top" style="height: 37px; width: 50%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN101/n:CSOSN" >
										<td valign ="top" style="height: 37px; width: 50%;">
											<span class="TextoFundoBrancoNegrito">
												Código de Situação da Operação – Simples Nacional<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='101'">
													101- Tributada pelo Simples Nacional com permissão de crédito
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN101/n:pCredSN" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota aplicável de cálculo do crédito <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN101/n:vCredICMSSN" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor de crédito do ICMS<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!--ICMSSN102 -->
						<xsl:variable name="icmssn102" select="n:imposto/n:ICMS/n:ICMSSN102"/>
						<xsl:if test="$icmssn102!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN102/n:orig" >
										<td valign ="top" style="height: 37px; width: 50%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN102/n:CSOSN" >
										<td valign ="top" style="height: 37px; width: 50%;">
											<span class="TextoFundoBrancoNegrito">
												Código de Situação da Operação – Simples Nacional<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='102'">
													102- Tributada pelo Simples Nacional sem permissão de crédito
												</xsl:if>
												<xsl:if test="$origmerc='103'">
													103 – Isenção do ICMS no Simples Nacional para faixa de receita bruta
												</xsl:if>
												<xsl:if test="$origmerc='300'">
													300 – Imune
												</xsl:if>
												<xsl:if test="$origmerc='400'">
													400 – Não tributada pelo Simples Nacional
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!--ICMSSN201 -->
						<xsl:variable name="icmssn201" select="n:imposto/n:ICMS/n:ICMSSN201"/>
						<xsl:if test="$icmssn201!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:CSOSN" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Código de Situação da Operação<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='201'">
													201- Tributada pelo Simples Nacional com permissão de crédito e com cobrança do ICMS por Substituição Tributária
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:modBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade de determinação da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:variable name="modBCST" select="text()"/>
												<xsl:if test="$modBCST='0'">
													0 – Preço tabelado ou máximo sugerido
												</xsl:if>
												<xsl:if test="$modBCST='1'">
													1 - Lista Negativa (valor)
												</xsl:if>
												<xsl:if test="$modBCST='2'">
													2 - Lista Positiva (valor)
												</xsl:if>
												<xsl:if test="$modBCST='3'">
													3 - Lista Neutra (valor)
												</xsl:if>
												<xsl:if test="$modBCST='4'">
													4 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBCST='5'">
													5 - Pauta (valor)
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:pMVAST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da margem de valor Adicionado do ICMS ST <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:pRedBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Redução de BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:vBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:pICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do imposto do ICMS ST <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:vICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:pCredSN" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota aplicável de cálculo do crédito<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN201/n:vCredICMSSN" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor crédito do ICMS <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<td>
									</td>
									<td>
									</td>
								</tr>
							</table>
						</xsl:if>

						<!--ICMSSN202 -->
						<xsl:variable name="icmssn202" select="n:imposto/n:ICMS/n:ICMSSN202"/>
						<xsl:if test="$icmssn202!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN202/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN202/n:CSOSN" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Código de Situação da Operação<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='202'">
													202- Tributada pelo Simples Nacional sem permissão de crédito e com cobrança do ICMS por Substituição Tributária
												</xsl:if>
												<xsl:if test="$origmerc='203'">
													203- Isenção do ICMS nos Simples Nacional para faixa de receita bruta e com cobrança do ICMS por Substituição Tributária
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN202/n:modBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade de determinação da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:variable name="modBCST" select="text()"/>
												<xsl:if test="$modBCST='0'">
													0 – Preço tabelado ou máximo sugerido
												</xsl:if>
												<xsl:if test="$modBCST='1'">
													1 - Lista Negativa (valor)
												</xsl:if>
												<xsl:if test="$modBCST='2'">
													2 - Lista Positiva (valor)
												</xsl:if>
												<xsl:if test="$modBCST='3'">
													3 - Lista Neutra (valor)
												</xsl:if>
												<xsl:if test="$modBCST='4'">
													4 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBCST='5'">
													5 - Pauta (valor)
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN202/n:pMVAST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da margem de valor Adicionado do ICMS ST <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN202/n:pRedBCST" >
										<td valign ="top" style="height: 37px; width:33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Redução de BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN202/n:vBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN202/n:pICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do imposto do ICMS ST <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN202/n:vICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<td></td>
								</tr>
							</table>
						</xsl:if>

						<!--ICMSSN500 -->
						<xsl:variable name="icmssn500" select="n:imposto/n:ICMS/n:ICMSSN500"/>
						<xsl:if test="$icmssn500!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN500/n:orig" >
										<td valign ="top" style="height: 37px; width: 50%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN500/n:CSOSN" >
										<td valign ="top" style="height: 37px; width: 50%;">
											<span class="TextoFundoBrancoNegrito">
												Código de Situação da Operação<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='500'">
													500 – ICMS cobrado anteriormente por substituição tributária (substituído) ou por antecipação
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN500/n:vBCSTRet" >
										<td valign ="top" style="height: 37px; width: 50%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS ST retido <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN500/n:vICMSSTRet" >
										<td valign ="top" style="height: 37px; width:50%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST retido <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!--ICMSSN900 -->
						<xsl:variable name="icmssn900" select="n:imposto/n:ICMS/n:ICMSSN900"/>
						<xsl:if test="$icmssn900!=''">
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:orig" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Origem da Mercadoria<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="text()"/>
												<xsl:if test="$origmercst='0'">
													0–Nacional
												</xsl:if>
												<xsl:if test="$origmercst='1'">
													1–Estrangeira – Importação direta
												</xsl:if>
												<xsl:if test="$origmercst='2'">
													2–Estrangeira – Adquirida no mercado interno.
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:CSOSN" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Código de Situação da Operação<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmerc" select="text()"/>
												<xsl:if test="$origmerc='90'">
													90 - Outros
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:modBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade de determinação da BC do ICMS <br />
											</span>
											<span class="linha">
												<xsl:variable name="modBCST" select="text()"/>
												<xsl:if test="$modBCST='0'">
													0 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBCST='1'">
													1 - Pauta (valor)
												</xsl:if>
												<xsl:if test="$modBCST='2'">
													2 - Preço Tabelado Máx.
												</xsl:if>
												<xsl:if test="$modBCST='3'">
													3 - valor da operação
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:vBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:pRedBC" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Redução de BC <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:pICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do imposto<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:vICMS" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:modBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Modalidade de determinação da BC do ICMS ST <br />
											</span>
											<span class="linha">
												<xsl:variable name="modBCST" select="text()"/>
												<xsl:if test="$modBCST='0'">
													0 – Preço tabelado ou máximo sugerido
												</xsl:if>
												<xsl:if test="$modBCST='1'">
													1 - Lista Negativa (valor)
												</xsl:if>
												<xsl:if test="$modBCST='2'">
													2 - Lista Positiva (valor)
												</xsl:if>
												<xsl:if test="$modBCST='3'">
													3 - Lista Neutra (valor);
												</xsl:if>
												<xsl:if test="$modBCST='4'">
													4 - Margem Valor Agregado (%)
												</xsl:if>
												<xsl:if test="$modBCST='5'">
													5 - Pauta (valor)
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:pMVAST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da margem de valor Adicionado do ICMS ST <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:pRedBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Percentual da Redução de BC  <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:vBCST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:pICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota do imposto<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:vICMSST" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS  <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:vBCSTRet" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor da BC do ICMS retido<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:vICMSSTRet" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ICMS ST retido<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:pCredSN" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Alíquota aplicável de cálculo do crédito  <br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ICMS/n:ICMSSN900/n:vCredICMSSN" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor de crédito do ICMS<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<td></td>
								</tr>
							</table>
						</xsl:if>

						<!--IPI - IMPOSTO SOBRE PRODUTOS INDUSTRIALIZADOS-->
						<xsl:variable name="ipi" select="n:imposto/n:IPI"/>
						<xsl:if test="$ipi!=''">
							<table align="center"  width="98%">
								<tr class="TextoFundoBranco">
									<td  class="TituloAreaRestrita2">
										IMPOSTO SOBRE PRODUTOS INDUSTRIALIZADOS
									</td>
								</tr>
							</table>
							<xsl:if test="$ipi!=''">
								<table align="center"  class="textoVerdana7" width="98%">
									<tr class="TextoFundoBranco">
										<xsl:for-each select ="n:imposto/n:IPI/n:clEnq" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Classe de Enquadramento<br />
												</span>
												<span class="linha">
													<xsl:value-of select="text()" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:IPI/n:cEnq" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Código de Enquadramento<br />
												</span>
												<span class="linha">
													<xsl:value-of select="text()" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:IPI/n:cSelo" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Código do Selo<br />
												</span>
												<span class="linha">
													<xsl:value-of select="text()" />
												</span>
											</td>
										</xsl:for-each>
									</tr>
									<tr class="TextoFundoBranco" >
										<xsl:for-each select ="n:imposto/n:IPI/n:CNPJProd">
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													CNPJ do Produtor
													<br />
												</span>
												<span class="linha">
													<xsl:value-of select="chave:formatarCnpj(text())" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:IPI/n:qSelo" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Qtd. Selo<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:IPI/n:IPITrib/n:qUnid" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Qtd Total Unidade Padrão<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,000')" />
												</span>
											</td>
										</xsl:for-each>
									</tr>
									<tr class="TextoFundoBranco">
										<xsl:for-each select ="n:imposto/n:IPI/n:IPITrib/n:vBC" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Base de Cálculo <br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:IPI/n:IPITrib/n:pIPI" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Alíquota<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:IPI/n:IPITrib/n:vIPI">
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Valor IPI<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
												</span>
											</td>
										</xsl:for-each>
									</tr>
									<tr class="TextoFundoBranco">
										<xsl:for-each select ="n:imposto/n:IPI/n:IPITrib/n:vUnid" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													Valor por Unidade<br />
												</span>
												<span class="linha">
													<xsl:value-of select="format-number(text(),'##.##.##0,0000')" />
												</span>
											</td>
										</xsl:for-each>
										<xsl:for-each select ="n:imposto/n:IPI/n:IPITrib/n:CST" >
											<td valign ="top" style="height: 37px; width: 33%;">
												<span class="TextoFundoBrancoNegrito">
													CST<br />
												</span>
												<span class="linha">
													<xsl:variable name="varcst" select="text()"/>
													<xsl:if test="$varcst='00'">
														00 - Entrada com recuperação de crédito
													</xsl:if>
													<xsl:if test="$varcst='49'">
														49 - Outras entradas
													</xsl:if>
													<xsl:if test="$varcst='50'">
														50 - Saída tributada
													</xsl:if>
													<xsl:if test="$varcst='99'">
														99 - Outras saídas
													</xsl:if>
												</span>
											</td>
										</xsl:for-each>
									</tr>
								</table>
							</xsl:if>
							<xsl:variable name="ipint" select="n:imposto/n:IPI/n:IPINT"/>
							<xsl:if test="$ipint!=''">
								<table align="center"  class="textoVerdana7" width="98%">
									<tr class="TextoFundoBranco" >
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												CST<br />
											</span>
											<span class="linha">
												<xsl:variable name="origmercst" select="n:imposto/n:IPI/n:IPINT/n:CST"/>
												<xsl:if test="$origmercst='01'">
													01-Entrada tributada com alíquota zero
												</xsl:if>
												<xsl:if test="$origmercst='02'">
													02-Entrada isenta
												</xsl:if>
												<xsl:if test="$origmercst='03'">
													03-Entrada não-tributada
												</xsl:if>
												<xsl:if test="$origmercst='04'">
													04-Entrada imune
												</xsl:if>
												<xsl:if test="$origmercst='05'">
													05-Entrada com suspensão
												</xsl:if>
												<xsl:if test="$origmercst='51'">
													51-Saída tributada com alíquota zero
												</xsl:if>
												<xsl:if test="$origmercst='52'">
													52-Saída isenta
												</xsl:if>
												<xsl:if test="$origmercst='53'">
													53-Saída não-tributada
												</xsl:if>
												<xsl:if test="$origmercst='54'">
													54-Saída imune
												</xsl:if>
												<xsl:if test="$origmercst='55'">
													55-Saída com suspensão
												</xsl:if>
											</span>
										</td>
									</tr>
								</table>
							</xsl:if>
						</xsl:if>

						<!-- DECLARAÇÃO DE IMPORTAÇÃO -->
            <xsl:variable name="DI" select="n:prod/n:DI"/>
            <xsl:if test="$DI!=''">
              <table align="center"  width="98%">
                <tr class="TextoFundoBranco">
                  <td  class="TituloAreaRestrita2">
                    DECLARAÇÃO DE IMPORTAÇÃO
                  </td>
                </tr>
              </table>
              <xsl:for-each select ="n:prod/n:DI">
                <table align="center"  class="textoVerdana7" width="98%">
                  <tr>
                    <td>
                      <span class="TextoFundoBrancoNegrito">DI/DSI/DA </span>
                    </td>
                    <td>
                      <span class="TextoFundoBrancoNegrito">Data Registro</span>
                    </td>
                    <td>
                      <span class="TextoFundoBrancoNegrito">Local Desembaraço Aduaneiro</span>
                    </td>
                    <td>
                      <span class="TextoFundoBrancoNegrito">Data </span>
                    </td>
                    <td>
                      <span class="TextoFundoBrancoNegrito">UF </span>
                    </td>
                    <td>
                      <span class="TextoFundoBrancoNegrito">Código do Exportador</span>
                    </td>
                  </tr>

                  <tr>
                    <td>
                      <span class="linha">
                        <xsl:value-of select="n:nDI" />
                      </span>
                    </td>
                    <td>
                      <span class="linha">
                        <xsl:value-of select="date:formatdate(n:dDI,'dd/MM/yyyy')" />
                      </span>
                    </td>
                    <td>
                      <span class="linha">
                        <xsl:value-of select="n:xLocDesemb" />
                      </span>
                    </td>
                    <td>
                      <span class="linha">
                        <xsl:value-of select="date:formatdate(n:dDesemb,'dd/MM/yyyy')" />
                      </span>
                    </td>
                    <td>
                      <span class="linha">
                        <xsl:value-of select="n:UFDesemb" />
                      </span>
                    </td>
                    <td>
                      <span class="linha">
                        <xsl:value-of select="n:cExportador" />
                      </span>
                    </td>
                  </tr>

                  <tr>
                    <!-- ADIÇÕES DA DI -->
                    <xsl:variable name="adi" select="n:adi" />
                    <xsl:if test="$adi!=''">
                      <td colspan="6" align="right">
                        <table width="75%">
                          <tr class="TextoFundoBranco">
                            <td class="TituloAreaRestrita2">
                              ADIÇÕES DA DI
                            </td>
                          </tr>

                          <tr class="TextoFundoBranco">
                            <td>
                              <table width="100%">
                                <tr>
                                  <td valign ="top"  align="left">
                                    <span class="TextoFundoBrancoNegrito">
                                      No. Adição
                                      <br />
                                    </span>
                                  </td>
                                  <td valign ="top"  align="left">
                                    <span class="TextoFundoBrancoNegrito">
                                      Item
                                      <br />
                                    </span>
                                  </td>
                                  <td valign ="top"  align="left">
                                    <span class="TextoFundoBrancoNegrito">
                                      Código Fabricante Estrangeiro
                                      <br />
                                    </span>

                                  </td>
                                  <td valign ="top"  align="left">
                                    <span class="TextoFundoBrancoNegrito">
                                      Valor do Desconto
                                      <br />
                                    </span>

                                  </td>
                                </tr>

                                <xsl:for-each select ="n:adi">
                                  <tr>
                                    <xsl:for-each select ="n:nAdicao" >
                                      <td valign ="top" align="left">
                                        <span class="linha">
                                          <xsl:value-of select="text()" />
                                        </span>
                                      </td>
                                    </xsl:for-each>

                                    <xsl:for-each select ="n:nSeqAdic" >
                                      <td valign ="top" align="left">
                                        <span class="linha">
                                          <xsl:value-of select="text()" />
                                        </span>
                                      </td>
                                    </xsl:for-each>

                                    <xsl:for-each select ="n:cFabricante" >
                                      <td valign ="top" align="left">
                                        <span class="linha">
                                          <xsl:value-of select="text()" />
                                        </span>
                                      </td>
                                    </xsl:for-each>

                                    <xsl:for-each select ="n:vDescDI" >
                                      <td valign ="top" align="left">
                                        <span class="linha">
                                          <xsl:value-of select="format-number(text(),'##.##.##0,00')" />
                                        </span>
                                      </td>
                                    </xsl:for-each>
                                  </tr>
                                </xsl:for-each>

                              </table>
                            </td>
                          </tr>


                        </table>
                      </td>
                    </xsl:if>
                  </tr>

                </table>
              </xsl:for-each>
            </xsl:if>
            <!-- FIM DA DECLARAÇÃO DE IMPORTAÇÃO -->

            <!--IMPOSTO DE IMPORTAÇÃO-->
            <xsl:variable name="iimp" select="n:imposto/n:II"/>
            <xsl:if test="$iimp!=''">
              <table align="center"  width="98%">
                <tr class="TextoFundoBranco">
                  <td  class="TituloAreaRestrita2">
                    IMPOSTO DE IMPORTAÇÃO
                  </td>
                </tr>
              </table>
              <table align="center"  class="textoVerdana7" width="98%">
                <tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:II/n:vBC">
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Base de Cálculo<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:II/n:vDespAdu">
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Despesas Aduaneiras<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:II/n:vII">
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Imposto de Importação<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:II/n:vIOF">
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												IOF<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>



						<!--dados da importação-->
						<xsl:variable name="imp" select="n:prod/n:DI"/>
						<xsl:if test="$imp!=''">
							<table align="center"  width="98%">
								<tr>
									<td  class="TituloAreaRestrita2">
										DOCUMENTO DE IMPORTAÇÃO
									</td>
								</tr>
							</table>
							<table align="center"  class="textoVerdana7" width="98%">
								<tbody>
									<tr class="TextoFundoBranco" >
										<span class="TextoFundoBrancoNegrito">
										</span>
									</tr>
									<tr>
										<xsl:for-each select ="n:prod/n:DI/n:nDI">
											<td valign ="top">
												<span class="TextoFundoBrancoNegrito">
													NÚMERO DI <xsl:value-of select="position()" />
												</span>
												<span class="linha">
													<xsl:value-of select="text()" />
													<br/>
												</span>
												<br/>
												<xsl:if test="position() mod 3 = 0 ">
													<tr></tr>
												</xsl:if>
											</td>
										</xsl:for-each>
									</tr>
								</tbody>
							</table>
						</xsl:if>

						<!--ISSQN-->
						<xsl:variable name="ISSQN" select="n:imposto/n:ISSQN"/>
						<xsl:if test="$ISSQN!=''">
							<table align="center"  width="98%">
								<tr>
									<td  class="TituloAreaRestrita2">
										ISSQN
									</td>
								</tr>
							</table>
							<table align="center"  class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco" >
									<xsl:for-each select ="n:imposto/n:ISSQN/n:vISSQN">
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Valor do ISSQN<br />
											</span>
											<span class="linha">
												<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
											</span>
										</td>
									</xsl:for-each>
									<xsl:for-each select ="n:imposto/n:ISSQN/n:cSitTrib">
										<td valign ="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Código de Tributação do ISSQN<br />
											</span>
											<span class="linha">
												<xsl:variable name="cSitTrib" select="text()"/>
												<xsl:if test="$cSitTrib='N'">
													N – NORMAL
												</xsl:if>
												<xsl:if test="$cSitTrib='R'">
													R – RETIDA
												</xsl:if>
												<xsl:if test="$cSitTrib='S'">
													S –SUBSTITUTA
												</xsl:if>
												<xsl:if test="$cSitTrib='I'">
													I – ISENTA
												</xsl:if>
											</span>
										</td>
									</xsl:for-each>
								</tr>
							</table>
						</xsl:if>

						<!-- IDENTIFICAÇÃO DOS MEDICAMENTOS -->
						<xsl:variable name="med" select="n:prod/n:med" />
						<xsl:if test="$med!=''">
							<table align="center" width="98%">
								<tr class="TextoFundoBranco">
									<td class="TituloAreaRestrita2">DETALHAMENTO ESPECÍFICO DOS MEDICAMENTOS</td>
								</tr>
							</table>
							<xsl:for-each select="n:prod/n:med">
								<xsl:variable name="chavesmed" select="child::node()[1]" />
								<xsl:variable name="chaves2med" select="concat('A',' + ',$chavesmed)"/>
								<table align="center"  class="textoVerdana7" width="98%" id="{$chaves2med}">
									<tr>										
										<td valign="top" style="width: 100%;">
											<span class="TextoFundoBrancoNegrito">
												Medicamento
												<xsl:value-of select="position()" />
												<br />
											</span>
										</td>
									</tr>
								</table>
								<xsl:variable name="tabmed" select="child::node()[1]"/>
								<xsl:variable name="tabmed2" select="concat(position(),' + ',$tabmed, ' + ', $chaves)"/>
								<table align="center"  class="textoVerdana7" width="98%">
									<tr>
										<td >
											<span class="TextoFundoBrancoNegrito">
												Nro. do Lote
												<br />
											</span>
											<span class="linha">
												<xsl:value-of select="child::node()[1]" />
											</span>
											<br/>
										</td>
										<td>
											<span class="TextoFundoBrancoNegrito">
												Quant. produto no lote
												<br />
											</span>
											<span class="linha">
												<xsl:value-of select=" format-number(child::node()[2],'##.##.##0,000')" />
											</span>
										</td>
										<td>
											<span class="TextoFundoBrancoNegrito">
												Data de fabricação
												<br />
											</span>
											<span class="linha">
												<xsl:value-of select="date:formatdate(child::node()[3],'dd/MM/yyyy')" />
											</span>
										</td>
									</tr>
									<tr>
										<td>
											<span class="TextoFundoBrancoNegrito">
												Data de validade
												<br />
											</span>
											<span class="linha">
												<xsl:value-of select="date:formatdate(child::node()[4],'dd/MM/yyyy')" />
											</span>
										</td>
										<td>
											<span class="TextoFundoBrancoNegrito">
												Preço Máximo Consumidor
												<br />
											</span>
											<span class="linha">
												<xsl:value-of select=" format-number(child::node()[5],'##.##.##0,00')" />
											</span>
										</td>
									</tr>
								</table>
							</xsl:for-each>
						</xsl:if>

						<!-- IDENTIFICAÇÃO DO VEÍCULO -->
						<xsl:variable name="veicProd" select="n:prod/n:veicProd" />
						<xsl:if test="$veicProd!=''">
							<table align="center" width="98%">
								<tr>
									<td class="TituloAreaRestrita2">DETALHAMENTO ESPECÍFICO DOS VEÍCULOS NOVOS</td>
								</tr>
							</table>
							<table align="center" class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco">
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Tipo da Operação
											<br />
										</span>
										<span class="linha">
											<xsl:variable name="operacao" select="n:prod/n:veicProd/n:tpOp" />
											<xsl:if test="$operacao='1'">1-Venda concessionária</xsl:if>
											<xsl:if test="$operacao='2'">2-Faturamento direto para consumidor final</xsl:if>
											<xsl:if test="$operacao='3'">3-Venda direta para grandes consumidores</xsl:if>
											<xsl:if test="$operacao='0'">0-Outros</xsl:if>
										</span>
									</td>
									<xsl:for-each select="n:prod/n:veicProd/n:chassi">
										<td valign="top" style="height: 37px; width: 33%;">
											<span class="TextoFundoBrancoNegrito">
												Chassi do veículo
												<br />
											</span>
											<span class="linha">
												<xsl:value-of select="text()" />
											</span>
										</td>
									</xsl:for-each>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Cilindradas <br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:cilin" />
										</span>
									</td>
								</tr>
								<tr class="TextoFundoBranco">
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Cor
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:cCor" />
										</span>
									</td>
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Descrição da cor
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:xCor" />
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Código da Cor - DENATRAN<br />
										</span>
										<span class="linha">
											<xsl:variable name="corveiculo" select="n:prod/n:veicProd/n:cCorDEN"/>
											<xsl:if test="$corveiculo='01'">
												01-AMARELO
											</xsl:if>
											<xsl:if test="$corveiculo='02'">
												02-AZUL
											</xsl:if>
											<xsl:if test="$corveiculo='03'">
												03-BEGE
											</xsl:if>
											<xsl:if test="$corveiculo='04'">
												04-BRANCA
											</xsl:if>
											<xsl:if test="$corveiculo='05'">
												05-CINZA
											</xsl:if>
											<xsl:if test="$corveiculo='06'">
												06-DOURADA
											</xsl:if>
											<xsl:if test="$corveiculo='07'">
												07-GRENA
											</xsl:if>
											<xsl:if test="$corveiculo='08'">
												08-LARANJA
											</xsl:if>
											<xsl:if test="$corveiculo='09'">
												09-MARROM
											</xsl:if>
											<xsl:if test="$corveiculo='10'">
												10-PRATA
											</xsl:if>
											<xsl:if test="$corveiculo='11'">
												11-PRETA
											</xsl:if>
											<xsl:if test="$corveiculo='12'">
												12-ROSA
											</xsl:if>
											<xsl:if test="$corveiculo='13'">
												13-ROXA
											</xsl:if>
											<xsl:if test="$corveiculo='14'">
												14-VERDE
											</xsl:if>
											<xsl:if test="$corveiculo='15'">
												15-VERMELHA
											</xsl:if>
											<xsl:if test="$corveiculo='16'">
												16-FANTASIA
											</xsl:if>
										</span>
									</td>

								</tr>
								<tr class="TextoFundoBranco">
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Peso Líquido
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(n:prod/n:veicProd/n:pesoL,'##.##.##0,0000')" />
										</span>
									</td>
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Peso Bruto
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(n:prod/n:veicProd/n:pesoB,'##.##.##0,0000')" />
										</span>
									</td>
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Serial (série)
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:nSerie" />
										</span>
									</td>
								</tr>
								<tr class="TextoFundoBranco">
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Tipo de Combustível
											<br />
										</span>
										<span class="linha">
											<xsl:variable name="tipoComb" select="n:prod/n:veicProd/n:tpComb" />
											<!-- <xsl:if test="$tipoComb='01'">
							   01-Álcool
						   </xsl:if>
						   <xsl:if test="$tipoComb='02'">
							   02-Gasolina
						   </xsl:if>
						   <xsl:if test="$tipoComb='03'">
							   03-Diesel
						   </xsl:if>
						   <xsl:if test="$tipoComb='16'">
							   16-Álcool/Gasolina
						   </xsl:if>
						   <xsl:if test="$tipoComb='17'">
							   17-Gasolina/Álcool/GNV
						   </xsl:if>
						   <xsl:if test="$tipoComb='18'">
							   18-Gasolina/Elétrico
						   </xsl:if> -->
											<xsl:value-of select="$tipoComb" />
										</span>
									</td>
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Número de Motor
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:nMotor" />
										</span>
									</td>
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Capacidade Máxima de Tração
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(n:prod/n:veicProd/n:CMT,'##.##.##0,0000')" />
										</span>
									</td>
								</tr>
								<tr class="TextoFundoBranco">
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Distância entre eixos
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="format-number(n:prod/n:veicProd/n:dist,'##.##.##0,0000')" />
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Ano Modelo de Fabricação <br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:anoMod" />
										</span>
									</td>
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Ano de Fabricação
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:anoFab" />
										</span>
									</td>
								</tr>
								<tr class="TextoFundoBranco">
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Tipo de Pintura
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:tpPint" />
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Tipo de Veículo <br />
										</span>
										<span class="linha">
											<xsl:variable name="tipoVeic" select="n:prod/n:veicProd/n:tpVeic" />
											<!--<xsl:if test="$tipoVeic='06'">
								 06-AUTOMÓVEL
							 </xsl:if>
							 <xsl:if test="$tipoVeic='14'">
								 14-CAMINHÃO
							 </xsl:if>
							 <xsl:if test="$tipoVeic='13'">
								 13-CAMINHONETA
							 </xsl:if>
							 <xsl:if test="$tipoVeic='24'">
								 24-CARGA / CAM
							 </xsl:if>
							 <xsl:if test="$tipoVeic='02'">
								 02-CICLOMOTO
							 </xsl:if>
							 <xsl:if test="$tipoVeic='22'">
								 22-ESP / ÔNIBUS
							 </xsl:if>
							 <xsl:if test="$tipoVeic='07'">
								 07-MICROÔNIBUS
							 </xsl:if>
							 <xsl:if test="$tipoVeic='23'">
								 23-MISTO / CAM
							 </xsl:if>
							 <xsl:if test="$tipoVeic='04'">
								 04-MOTOCICLO
							 </xsl:if>
							 <xsl:if test="$tipoVeic='03'">
								 03-MOTONETA
							 </xsl:if>
							 <xsl:if test="$tipoVeic='08'">
								 08-ÔNIBUS
							 </xsl:if>
							 <xsl:if test="$tipoVeic='10'">
								 10-REBOQUE
							 </xsl:if>
							 <xsl:if test="$tipoVeic='05'">
								 05-TRICICLO
							 </xsl:if>
							 <xsl:if test="$tipoVeic='17'">
								 17-C. TRATOR
							 </xsl:if>-->
											<xsl:value-of select="$tipoVeic" />
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Espécie de Veículo<br />
										</span>
										<span class="linha">
											<xsl:variable name="espVeic" select="n:prod/n:veicProd/n:espVeic" />
											<xsl:if test="$espVeic='1'">
												1-PASSAGEIRO
											</xsl:if>
											<xsl:if test="$espVeic='2'">
												2-CARGA
											</xsl:if>
											<xsl:if test="$espVeic='3'">
												3-MISTO
											</xsl:if>
											<xsl:if test="$espVeic='4'">
												4-CORRIDA
											</xsl:if>
											<xsl:if test="$espVeic='5'">
												5-TRAÇÃO
											</xsl:if>
											<xsl:if test="$espVeic='6'">
												6-ESPECIAL
											</xsl:if>
										</span>
									</td>
								</tr>
								<tr class="TextoFundoBranco">
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Condição do VIN (Vehicle Identification Number)<br />
										</span>
										<span class="linha">
											<xsl:variable name="tpVin" select="n:prod/n:veicProd/n:VIN" />
											<xsl:if test="$tpVin='R'">
												R-Remarcado
											</xsl:if>
											<xsl:if test="$tpVin='N'">
												N-Normal
											</xsl:if>
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Condição do Veículo<br />
										</span>
										<span class="linha">
											<xsl:variable name="condveiculo" select="n:prod/n:veicProd/n:condVeic"/>
											<xsl:if test="$condveiculo='1'">
												1-Acabado
											</xsl:if>
											<xsl:if test="$condveiculo='2'">
												2-Inacabado
											</xsl:if>
											<xsl:if test="$condveiculo='3'">
												3-Semi-acabado
											</xsl:if>
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Código Marca Modelo <br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:cMod" />
										</span>
									</td>
								</tr>
								<tr class="TextoFundoBranco">
									<td valign="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Potência Motor (CV)
											<br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:pot" />
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Capacidade máxima de lotação <br />
										</span>
										<span class="linha">
											<xsl:value-of select="n:prod/n:veicProd/n:lota" />
										</span>
									</td>
									<td valign ="top" style="height: 37px; width: 33%;">
										<span class="TextoFundoBrancoNegrito">
											Restrição <br />
										</span>
										<span class="linha">
											<xsl:variable name="restricao" select="n:prod/n:veicProd/n:tpRest" />
											<xsl:if test="$restricao='0'">
												0-Não há
											</xsl:if>
											<xsl:if test="$restricao='1'">
												1-Alienação Fiduciária
											</xsl:if>
											<xsl:if test="$restricao='2'">
												2-Arrendamento Mercantil
											</xsl:if>
											<xsl:if test="$restricao='3'">
												3-Reserva de Domínio
											</xsl:if>
											<xsl:if test="$restricao='4'">
												4-Penhor de Veículos
											</xsl:if>
											<xsl:if test="$restricao='9'">
												9-Outras
											</xsl:if>
										</span>
									</td>
								</tr>
							</table>
						</xsl:if>

						<!-- IDENTIFICAÇÃO DO ARMAMENTO -->
						<xsl:variable name="arma" select="n:prod/n:arma" />
						<xsl:if test="$arma!=''">
							<table align="center" width="98%">
								<tr class="TextoFundoBranco">
									<td class="TituloAreaRestrita2">DETALHAMENTO ESPECÍFICO DE ARMAMENTO</td>
								</tr>
							</table>
							<xsl:for-each select="n:prod/n:arma">
								<xsl:variable name="chavesarma" select="child::node()[2]" />
								<xsl:variable name="chaves2arma" select="concat('A',' + ',$chavesarma)"/>
								<table align="center"  class="textoVerdana7" width="98%">
									<tr>										
										<td valign="top" style="width: 100%;">
											<span class="TextoFundoBrancoNegrito">
												Armamento
												<xsl:value-of select="position()" />
												<br />
											</span>
										</td>
									</tr>
								</table>
								<xsl:variable name="tabarma" select="child::node()[2]"/>
								<xsl:variable name="tabarma2" select="concat(position(),' + ',$tabarma)"/>
								<table align="center"  class="textoVerdana7" width="98%">
									<tr>
										<td  style="height: 37px">
											<span class="TextoFundoBrancoNegrito">
												Tipo de Arma de Fogo
												<br />
											</span>
											<span class="linha">
												<xsl:variable name="tparma" select="child::node()[1]"/>
												<xsl:if test="$tparma='0'">
													0 – Uso permitido
												</xsl:if>
												<xsl:if test="$tparma='1'">
													1 – Uso restrito
												</xsl:if>
											</span>
											<br/>
										</td>
										<td>
											<span class="TextoFundoBrancoNegrito">
												Número de Série da Arma
												<br />
											</span>
											<span class="linha">
												<xsl:value-of select="child::node()[2]" />
											</span>
										</td>
										<td>
											<span class="TextoFundoBrancoNegrito">
												Número de Série do Cano
												<br />
											</span>
											<span class="linha">
												<xsl:value-of select="child::node()[3]" />
											</span>
										</td>
									</tr>
									<tr class="TextoFundoBranco">
										<td colspan="3" valign ="top"  style="width: 100%;">
											<span class="TextoFundoBrancoNegrito">
												Descrição
											</span>
											<span class="linha">
												<xsl:value-of select="child::node()[4]" />
											</span>
											<br/>
										</td>
									</tr>
								</table>
							</xsl:for-each>
						</xsl:if>

						<!-- IDENTIFICAÇÃO DO COMBUSTÍVEL -->
						<xsl:variable name="comb" select="n:prod/n:comb" />
						<xsl:if test="$comb!=''">
							<table align="center" width="98%">
								<tr class="TextoFundoBranco">
									<td class="TituloAreaRestrita2">
										DETALHAMENTO ESPECÍFICO DE COMBUSTÍVEL
									</td>
								</tr>
							</table>
							<table align="center" class="textoVerdana7" width="98%">
								<tr class="TextoFundoBranco">
									<td>
										<table width="100%">
											<tr>
												<xsl:for-each select ="n:prod/n:comb/n:cProdANP">
													<td valign ="top" style="height: 37px; width: 25%;">
														<span class="TextoFundoBrancoNegrito">
															Código do Produto da ANP
															<br />
														</span>
														<span class="linha">
															<xsl:value-of select="text()" />
														</span>
													</td>
												</xsl:for-each>
												<xsl:for-each select ="n:prod/n:comb/n:CODIF">
													<td valign ="top" style="height: 37px; width: 25%;">
														<span class="TextoFundoBrancoNegrito">
															CODIF
															<br />
														</span>
														<span class="linha">
															<xsl:value-of select="text()" />
														</span>
													</td>
												</xsl:for-each>
												<xsl:for-each select ="n:prod/n:comb/n:qTemp">
													<td valign ="top" style="height: 37px; width: 25%;">
														<span class="TextoFundoBrancoNegrito">
															Quant. Combustível Faturada
															<br />
														</span>
														<span class="linha">
															<xsl:value-of select="format-number(text(),'##.##.##0,0000')" />
														</span>
													</td>
												</xsl:for-each>
												<xsl:for-each select ="n:prod/n:comb/n:UFCons">
													<td valign ="top" style="height: 37px; width: 25%;">
														<span class="TextoFundoBrancoNegrito">
															UF de consumo
															<br />
														</span>
														<span class="linha">
															<xsl:value-of select="text()" />
														</span>
													</td>
												</xsl:for-each>
											</tr>
										</table>
									</td>
								</tr>

								<tr>
									<!-- CIDE -->
									<xsl:variable name="cide" select="n:prod/n:comb/n:CIDE" />
									<xsl:if test="$cide!=''">
										<td colspan="6" align="right">
											<table width="75%">
												<tr class="TextoFundoBranco">
													<td class="TituloAreaRestrita2">
														CIDE
													</td>
												</tr>
												<tr class="TextoFundoBranco">
													<td>
														<table width="100%">
															<tr>
																<xsl:for-each select ="n:prod/n:comb/n:CIDE/n:qBCprod">
																	<td valign ="top" align="left" style="height: 37px; width: 33%;">
																		<span class="TextoFundoBrancoNegrito">
																			Quant. Base de Cálculo
																			<br />
																		</span>
																		<span class="linha">
																			<xsl:value-of select="format-number(text(),'##.##.##0,0000')" />
																		</span>
																	</td>
																</xsl:for-each>
																<xsl:for-each select ="n:prod/n:comb/n:CIDE/n:vAliqProd">
																	<td valign ="top" style="height: 37px; width: 33%;" align="left">
																		<span class="TextoFundoBrancoNegrito">
																			Valor da Alíquota (R$)
																			<br />
																		</span>
																		<span class="linha">
																			<xsl:value-of select="format-number(text(),'##.##.##0,0000')" />
																		</span>
																	</td>
																</xsl:for-each>
																<xsl:for-each select ="n:prod/n:comb/n:CIDE/n:vCIDE">
																	<td valign ="top" style="height: 37px; width: 33%;" align="left">
																		<span class="TextoFundoBrancoNegrito">
																			Valor
																			<br />
																		</span>
																		<span class="linha">
																			<xsl:value-of select="format-number(text(),'##.##.##0,00')" />
																		</span>
																	</td>
																</xsl:for-each>
															</tr>
														</table>
													</td>
												</tr>
											</table>
										</td>
									</xsl:if>
								</tr>

							</table>
						</xsl:if>


					</td>
				</tr>
			</table>
		</xsl:for-each>

	</xsl:template>
</xsl:stylesheet>