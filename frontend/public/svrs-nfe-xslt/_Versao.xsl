<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
	xmlns:fo="http://www.w3.org/1999/XSL/Format"
	xmlns:n="http://www.portalfiscal.inf.br/nfe"
	xmlns:s="http://www.w3.org/2000/09/xmldsig#"
	version="2.0"
	exclude-result-prefixes="fo n s">
  <!--
  Versão 1.0.0 -> Criação do arquivo
  Versão 1.1.1 -> Alteração no layout e criação de alerta nas abas para notas canceladas ou denegadas
  Versão 1.1.2 -> Correção na formatação do grupo de Cana de Açucar e inclusão dos eventos de confirmação de recebimento
  Versão 1.1.3 -> Alteração para informar se a NFe foi denegada devido ao destinatario ou ao emitente
  Versão 1.1.4 -> Inclusão da descrição de CST != 99 para os grupos COFINSOutr e PISOutr
  Versão 1.1.5 -> Correção para visualização das informações ICMSSN101,ICMSSN102,ICMSSN202,ICMSSN500  
  Versão 1.1.6 -> Inclusão da visualização de cancelamento como evento em "Eventos da NF-e" na aba NF-e 
                  Alteração do Label "Valor do ICMS ST" para "Valor ICMS desoneração" no elemento "imposto/ICMS/ICMS40/vICMS" do XML
  Versão 1.1.7 -> Label CPF/CNPJ editado para “CNPJ” ou “CPF” nas abas "NF-e", "Emitente", "Destinatário" e "Transporte"
                  Alterado nome de Evento: de “Recusa de Recebimento pelo Destinatário” para “Operação não Realizada”;
  Versão 1.1.8(31-07-2012) -> Correção na visualiação da data do evento da ultima carta de correção enviada
  Versão 1.1.9 ->  Correção no arquivo Template_Impostos.xsl nos valores de Percentual Redução de BC do ICMS Normal (ICMS20)
  
  Versão 3.0.0 - > Incluídos campos relativos a NFC-e
  
  Versão 3.0.2 - > Correção no arquivo ProdutosServicos que não estava mostrando as informações de CIDE no detalhamento de combustiveis
                   Inclusão de novos código de Origem de Mercadoria (3,4,5,6,7), conforme NT 2012/005
  Versão 3.0.3 - > Visualização dos dados de Cancelamento em uma Pop-up                   
                   Alteração no layout da Pop-up de carta de correção
  Versão 3.0.4 - > Correção na aba "Cobrança" para correta visualização de mais de uma forma de pagamento
                   Inclusão do campo "Valor Total dos Tributos" nas abas "Totais" e "Produtos e Serviços"
               - > Inclusão dos eventos de CT-e(610600 e 610601)na aba "NF-e"
  Versão 3.0.5 - > Correção no arquivo ProdutosServicos.xsl para apresentar os dados de Base de Cálculo e Alíquota referente ao IPI.     
  Versão 3.0.6 - > Permite a apresentação da existencia de todos tipos de eventos para a NFe.                    
  Versão 3.0.7 - > Alteradas as descrições para origem de mercadoria 3,4 e 5. 
                   Incluída origem de mercadoria 8. 
                   Adicionado campo "Número FCI".
  Versão:3.0.8 - > Correção visualização campo Data de Emissão para casos com eventos CT-e Autorizado relacionados.
                   Adição de descrição para o campo cStat com valor 150 - "Autorização de Uso Fora de Prazo".
                   Criação dos XSLT's para visualização dos detalhes dos eventos a seguir:
                        1. Confirmação da Operação pelo Destinatário
                        2. Ciência da Operação pelo Destinatário
                        3. Desconhecimento da Operação pelo Destinatário
                        4. Operação não Realizada
                        5. Registro Passagem NF-e
                        6. Cancelamento Registro Passagem NF-e
                        7. CT-e Autorizado
                        8. CT-e Cancelado
                        9. Vistoria SUFRAMA
                        10. Internalização SUFRAMA
                   Identificação do ambiente onde os XSLT’s estão sendo executados (consulta pública ou intranet). Isso permite mostrar apenas detalhes de eventos pertinentes para cada uma das áreas. 
                        1. Visualização Parcial
                        2. Visualização Completa
                        3. Não Visível
  Versão:3.0.9 - > Cancelamento Registro Passagem agora não tem mais seus detalhes apresentados quando o tipo de consulta for pública.
                   Correção na apresentação do último evento de Carta de Correção Eletrônica.
                   Correção na listagem de volumes da aba Transporte.
  Versão:3.1.0 - > 
                   Correção do teste para apresentação/ocultamento dos dados relacionados a impostos e informações adicionais.  
                   Compilação de campos que não eram apresentados pelo XSLT da versão 3.10 da Nota Fiscal e anteriores.
                   Localização de campos com valor sem descrição e inclusão da descrição para os mesmos (campos com tabela de domínio).
                   Tratamento e  inclusão de cláusulas xsl:otherwise, corrigindo a apresentação de valores onde a inexistência de algum código, por exemplo (inclusão futura de um código tPag=14) impossibilitaria a visualização do valor sem a geração de uma nova versão do arquivo. 
                   Novo layout mais limpo e melhoria no código CSS (incluindo melhor separação do código CSS e HTML).
                   Separação e melhor organização dos scripts JS.
                   Inclusão de XSLT para visualização dos eventos MDF-e Autorizado, MDF-e Cancelado e Registro de Passagem NF-e BRId.
                   Correção na formatação das casas decimais de diversos campos numéricos.
                   Correção para apresentação de campos de diferentes versões com nomes diferentes.
  Versão:3.1.1 - > 
                   Inclusão do evento EPEC.
                   Inclusão campos evento Registro Passagem.
  Versão:3.1.2 - > 
                   Alteração do nome da colune Data/Hora na listagem dos eventos para Data Autorização.
                   Inclusão da coluna Data Inclusão BD, na listagem dos eventos da Nota Fiscal.                     
				           Correção da apresentação do campo Autor do Evento, para os eventos de CC-e e Cancelamento.
  -->
  <xsl:output method="html"/>
  <xsl:template match="Versao" name="Versao">
    v3.1.2 
  </xsl:template>
</xsl:stylesheet>