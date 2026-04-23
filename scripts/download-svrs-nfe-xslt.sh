#!/usr/bin/env bash
# Baixa o pacote XSLT NF-e do SVRS para frontend/public/svrs-nfe-xslt
# Fonte: https://dfe-portal.svrs.rs.gov.br/Schemas/PRNFE/XSLT/NFe/
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DEST="$ROOT/frontend/public/svrs-nfe-xslt"
BASE="https://dfe-portal.svrs.rs.gov.br/Schemas/PRNFE/XSLT/NFe"
mkdir -p "$DEST"
FILES=(
  _Estilos_Geral.xsl _ImpressaoAutuso.xsl _ImpressaoNFe.xsl _ImpressaoResumo.xsl _Scripts_Geral.xsl _Versao.xsl
  _Visualizacao_Internet.xsl _Visualizacao_Intranet.xsl _Visualizacao_RegPass.xsl _Visualizacao_Resumo.xsl _Visualizacao_Sistema_Externo.xsl
  Autorizacao_Uso.xsl Avulsa.xsl Cobranca.xsl consNFePassAnt_Passagens.xsl css.xsl DestinatarioRemetente.xsl Dpec.xsl Emitente.xsl
  Evento_Cancelamento.xsl Evento_Cancelamento_Registro_Passagem.xsl Evento_CCe.xsl Evento_Ciencia.xsl Evento_Confirmacao.xsl
  Evento_CTe_Autorizado.xsl Evento_CTe_Cancelado.xsl Evento_nao_Realizado.xsl Evento_EPEC.xsl Evento_Internalizacao_SUFRAMA.xsl
  Evento_MDFe_Autorizado.xsl Evento_MDFe_Cancelado.xsl Evento_Registro_Passagem.xsl Evento_Registro_Passagem_BRId.xsl
  Evento_Registro_Passagem_Sistema_Externo.xsl Evento_Vistoria_SUFRAMA.xsl Evento_Desconhecimento.xsl
  InformacoesAdicionais.xsl InformacoesAdicionaisImp.xsl NFe.xsl NFe_Cancelamento.xsl ProdutosServicos.xsl ProdutosServicosImp.xsl
  RegPass_Dpec.xsl RegPass_Passagens.xsl Template_Impostos.xsl Totais.xsl Transporte.xsl Utils.xsl Vistoria_SUFRAMA.xsl VisualizacaoEvento.xsl
)
for f in "${FILES[@]}"; do
  curl -fsS "$BASE/$f" -o "$DEST/$f"
done
# NFe.xsl importa Evento_Nao_Realizado.xsl (case); duplicar a partir do arquivo publicado
cp "$DEST/Evento_nao_Realizado.xsl" "$DEST/Evento_Nao_Realizado.xsl"
echo "OK -> $DEST ($(ls -1 "$DEST" | wc -l) arquivos)"
