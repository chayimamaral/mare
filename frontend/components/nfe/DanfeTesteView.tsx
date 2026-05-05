import React, { useEffect, useMemo, useRef } from 'react';
import JsBarcode from 'jsbarcode';

import type { NFEDanfeView } from '../../lib/nfeDanfeClient';
import styles from '../../styles/nfe-danfe-teste.module.css';

type Props = {
  /** Opcional: no momento a casca não preenche dados, mas já recebe a estrutura para evolução. */
  data?: NFEDanfeView | null;
  /** `print`: apenas a folha A4 (pré-impressão / página dedicada). */
  variant?: 'default' | 'print';
};

type Box = {
  key: string;
  label: string;
  topCm: number;
  leftCm: number;
  widthCm: number;
  heightCm: number;
  sample?: string;
  kind?: 'normal' | 'produtos' | 'barcode';
};

function cm(n: number): string {
  return `${n.toFixed(2)}cm`;
}

const fmt = (v?: string | number | null) => (v == null || String(v).trim() === '' ? '' : String(v).trim());

function onlyDigits(v?: string | null): string {
  return String(v ?? '').replace(/\D/g, '');
}

function fmtCNPJCPF(v?: string | null): string {
  const d = onlyDigits(v);
  if (d.length === 14) return d.replace(/^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$/, '$1.$2.$3/$4-$5');
  if (d.length === 11) return d.replace(/^(\d{3})(\d{3})(\d{3})(\d{2})$/, '$1.$2.$3-$4');
  return fmt(v);
}

function fmtBRL(v?: string | number | null): string {
  if (v == null || String(v).trim() === '') return '';
  const n = Number(String(v).replace(',', '.'));
  if (!Number.isFinite(n)) return String(v);
  return n.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function fmtDateBR(isoLike?: string | null): string {
  const s = String(isoLike ?? '').trim();
  if (!s) return '';
  const d = new Date(s);
  if (Number.isNaN(d.getTime())) return '';
  return d.toLocaleDateString('pt-BR', { timeZone: 'UTC' });
}

export function DanfeTesteView({ data, variant = 'default' }: Props) {
  const barcodeRef = useRef<SVGSVGElement | null>(null);
  const produtoCols = useMemo(
    () => [
      // Larguras em cm somando ~20.57 (largura do quadro de produtos/serviços).
      { key: 'cod', label: 'CÓDIGO', w: 1.6 },
      { key: 'desc', label: 'DESCRIÇÃO DOS PRODUTOS/SERVIÇOS', w: 5.4 },
      { key: 'ncm', label: 'NCM/SH', w: 1.1 },
      { key: 'cst', label: 'CST', w: 0.9 },
      { key: 'cfop', label: 'CFOP', w: 0.9 },
      { key: 'un', label: 'UN', w: 0.8 },
      { key: 'qtd', label: 'QUANT.', w: 1.2 },
      { key: 'vunit', label: 'V. UNIT.', w: 1.4 },
      { key: 'vdesc', label: 'DESC.', w: 1.1 },
      { key: 'vtotal', label: 'V. TOTAL', w: 1.4 },
      { key: 'bcicms', label: 'B.CALC. ICMS', w: 1.1 },
      { key: 'vicms', label: 'V. ICMS', w: 1.0 },
      { key: 'bcst', label: 'B.CALC. ST', w: 1.1 },
      { key: 'vst', label: 'V. ST', w: 1.0 },
      { key: 'vipi', label: 'V. IPI', w: 1.0 },
      { key: 'aicms', label: 'ALIQ. ICMS', w: 0.7 },
      { key: 'aipi', label: 'ALIQ. IPI', w: 0.7 },
    ],
    [],
  );

  const values = useMemo(() => {
    const id = data?.identificacao;
    const em = data?.emitente;
    const de = data?.destinatario;
    const tot = data?.totais;
    const tr = data?.transporte;
    return {
      chave: fmt(id?.chave),
      natureza: fmt(id?.natureza_operacao),
      emissao: fmtDateBR(id?.emissao_em ?? null),
      entrada: fmtDateBR(id?.saida_entrada_em ?? null),
      // Emitente
      em_nome: fmt(em?.nome),
      em_fantasia: fmt((em as any)?.nome_fantasia),
      em_cnpj: fmtCNPJCPF(em?.cnpj_cpf),
      em_ie: fmt(em?.ie),
      em_ie_st: fmt((em as any)?.ie_substituto),
      // Destinatário
      de_nome: fmt(de?.nome),
      de_cnpj: fmtCNPJCPF(de?.cnpj_cpf),
      de_ie: fmt(de?.ie),
      de_end: fmt(de?.endereco_completo || [de?.logradouro, de?.numero].filter(Boolean).join(', ')),
      de_bairro: fmt(de?.bairro),
      de_cep: fmt(de?.cep),
      de_mun: fmt(de?.municipio_cod_nome || de?.municipio),
      de_fone: fmt(de?.telefone),
      de_uf: fmt(de?.uf),
      // Totais
      v_bc_icms: fmtBRL((tot as any)?.base_icms),
      v_icms: fmtBRL((tot as any)?.valor_icms),
      v_bc_st: fmtBRL((tot as any)?.base_icms_st),
      v_st: fmtBRL((tot as any)?.valor_st),
      v_prod: fmtBRL((tot as any)?.valor_produtos),
      v_frete: fmtBRL((tot as any)?.valor_frete),
      v_seg: fmtBRL((tot as any)?.valor_seguro),
      v_desc: fmtBRL((tot as any)?.valor_desconto),
      v_outros: fmtBRL((tot as any)?.valor_outros),
      v_ipi: fmtBRL((tot as any)?.valor_ipi),
      v_nf: fmtBRL((tot as any)?.valor_nota),
      // Transporte
      tr_nome: fmt((tr as any)?.transportador),
      tr_cnpj: fmtCNPJCPF((tr as any)?.cnpj_cpf),
      tr_end: fmt((tr as any)?.endereco),
      tr_mun: fmt((tr as any)?.municipio),
      tr_uf: fmt((tr as any)?.uf),
      tr_ie: fmt((tr as any)?.ie),
      tr_placa: fmt((tr as any)?.placa),
      tr_antt: fmt((tr as any)?.antt),
      tr_rntc: fmt((tr as any)?.rntc),
      tr_qtd: fmt((tr as any)?.quantidade_volumes),
      tr_especie: fmt((tr as any)?.especie),
      tr_marca: fmt((tr as any)?.marca),
      tr_numeracao: fmt((tr as any)?.numeracao),
      tr_peso_b: fmt((tr as any)?.peso_bruto),
      tr_peso_l: fmt((tr as any)?.peso_liquido),
      // ISSQN (se existir no JSON; alguns layouts não trazem)
      iss_im: fmt((tot as any)?.inscricao_municipal || (em as any)?.im),
      iss_v_serv: fmtBRL((tot as any)?.valor_total_servicos),
      iss_bc: fmtBRL((tot as any)?.base_issqn),
      iss_v: fmtBRL((tot as any)?.valor_issqn),
      // Adicionais
      info_comp: fmt((data as any)?.adicionais?.informacoes_complementares),
      info_fisco: fmt((data as any)?.adicionais?.informacoes_fisco),
    };
  }, [data]);

  useEffect(() => {
    const svg = barcodeRef.current;
    const chave = onlyDigits(values.chave);
    if (!svg) return;
    if (chave.length !== 44) return;
    try {
      JsBarcode(svg, chave, {
        format: 'CODE128C',
        displayValue: false,
        margin: 0,
        background: 'transparent',
        lineColor: '#111',
        // Altura em px; o SVG está em 100% do box, então isso só define proporção interna.
        height: 42,
      });
    } catch {
      // Se falhar, mantém vazio (casca).
    }
  }, [values.chave]);

  const boxes: Box[] = useMemo(
    () => [
      // 3.8.1 (A4 retrato) — primeiras linhas do layout (centímetros).
      // CANHOTO
      { key: 'canhoto-recebemos', label: 'RECEBEMOS DE...', topCm: 0.42, leftCm: 0.25, widthCm: 16.1, heightCm: 0.85 },
      { key: 'canhoto-nfe', label: 'NF-e / Nº / SÉRIE', topCm: 0.42, leftCm: 16.35, widthCm: 4.5, heightCm: 1.7, sample: 'NF-e / Nº 000.000.000 / SÉRIE 000' },
      { key: 'canhoto-data', label: 'DATA DE RECEBIMENTO', topCm: 1.27, leftCm: 0.25, widthCm: 4.1, heightCm: 0.85 },
      { key: 'canhoto-ass', label: 'IDENTIFICAÇÃO E ASSINATURA...', topCm: 1.27, leftCm: 4.35, widthCm: 12.1, heightCm: 0.85 },

      // DADOS DA NF-e — quadros principais (casca)
      { key: 'emitente-quadro', label: 'EMITENTE', topCm: 2.54, leftCm: 0.25, widthCm: 5.33, heightCm: 3.92, sample: [values.em_nome, values.em_fantasia, values.em_cnpj].filter(Boolean).join(' / ') },
      { key: 'danfe-titulo', label: 'DANFE', topCm: 2.54, leftCm: 5.58, widthCm: 2.54, heightCm: 3.92, sample: 'DANFE' },
      { key: 'barcode-quadro', label: 'QUADRO CÓDIGO DE BARRAS (CHAVE)', topCm: 2.54, leftCm: 8.12, widthCm: 12.7, heightCm: 1.48 },
      { key: 'barcode', label: 'CÓDIGO DE BARRAS DA CHAVE', topCm: 2.78, leftCm: 8.62, widthCm: 11.5, heightCm: 1.0, kind: 'barcode' },
      { key: 'chave', label: 'CHAVE DE ACESSO (44)', topCm: 4.02, leftCm: 8.12, widthCm: 12.7, heightCm: 0.85, sample: values.chave },
      // Campos variáveis (3.9) — colocados logo abaixo da chave, na área antes da "Natureza da operação".
      { key: 'campo-var-1', label: 'CAMPO 1 (CONTEÚDO VARIÁVEL)', topCm: 4.87, leftCm: 8.12, widthCm: 12.7, heightCm: 0.80, sample: 'Consulta de autenticidade no portal nacional da NF-e / Sefaz' },
      { key: 'campo-var-2', label: 'CAMPO 2 (CONTEÚDO VARIÁVEL)', topCm: 5.67, leftCm: 8.12, widthCm: 12.7, heightCm: 0.79, sample: fmt((data as any)?.identificacao?.protocolo) || 'PROTOCOLO / DATA-HORA AUTORIZAÇÃO' },
      { key: 'natureza', label: 'NATUREZA DA OPERAÇÃO', topCm: 6.46, leftCm: 0.25, widthCm: 7.87, heightCm: 0.85, sample: values.natureza },
      { key: 'insc-est', label: 'INSCRIÇÃO ESTADUAL DO EMITENTE', topCm: 7.31, leftCm: 0.25, widthCm: 6.86, heightCm: 0.85, sample: values.em_ie },
      { key: 'insc-st', label: 'INSCRIÇÃO ESTADUAL ST DO EMITENTE', topCm: 7.31, leftCm: 7.11, widthCm: 6.86, heightCm: 0.85, sample: values.em_ie_st },
      { key: 'cnpj-em', label: 'CNPJ DO EMITENTE', topCm: 7.31, leftCm: 13.97, widthCm: 6.86, heightCm: 0.85, sample: values.em_cnpj },
      { key: 'dest-titulo', label: 'DESTINATÁRIO/REMETENTE', topCm: 8.16, leftCm: 0.25, widthCm: 3.3, heightCm: 0.42, sample: 'DESTINATÁRIO/REMETENTE' },
      { key: 'dest-razao', label: 'RAZÃO SOCIAL', topCm: 8.58, leftCm: 0.25, widthCm: 12.32, heightCm: 0.85, sample: values.de_nome },
      { key: 'dest-cnpj', label: 'CNPJ', topCm: 8.58, leftCm: 12.57, widthCm: 5.33, heightCm: 0.85, sample: values.de_cnpj },
      { key: 'dest-emissao', label: 'DATA DA EMISSÃO', topCm: 8.58, leftCm: 17.9, widthCm: 2.92, heightCm: 0.85, sample: values.emissao },
      { key: 'dest-end', label: 'ENDEREÇO', topCm: 9.43, leftCm: 0.25, widthCm: 10.16, heightCm: 0.85, sample: values.de_end },
      { key: 'dest-bairro', label: 'BAIRRO/DISTRITO', topCm: 9.43, leftCm: 10.41, widthCm: 4.83, heightCm: 0.85, sample: values.de_bairro },
      { key: 'dest-cep', label: 'CEP', topCm: 9.43, leftCm: 15.24, widthCm: 2.67, heightCm: 0.85, sample: values.de_cep },
      { key: 'dest-entrada', label: 'DATA ENTRADA/SAÍDA', topCm: 9.43, leftCm: 17.91, widthCm: 2.92, heightCm: 0.85, sample: values.entrada },
      { key: 'dest-mun', label: 'MUNICÍPIO', topCm: 10.28, leftCm: 0.25, widthCm: 7.11, heightCm: 0.85, sample: values.de_mun },
      { key: 'dest-fone', label: 'FONE/FAX', topCm: 10.28, leftCm: 7.36, widthCm: 4.06, heightCm: 0.85, sample: values.de_fone },
      { key: 'dest-uf', label: 'UF', topCm: 10.28, leftCm: 11.42, widthCm: 1.14, heightCm: 0.85, sample: values.de_uf },
      { key: 'dest-ie', label: 'INSCRIÇÃO ESTADUAL', topCm: 10.28, leftCm: 12.56, widthCm: 5.33, heightCm: 0.85, sample: values.de_ie },
      { key: 'dest-hora', label: 'HORA DA ENTRADA/SAÍDA', topCm: 10.28, leftCm: 17.89, widthCm: 2.92, heightCm: 0.85 },

      // FATURA / DUPLICATAS
      { key: 'fat-titulo', label: 'FATURA/DUPLICATAS', topCm: 11.09, leftCm: 0.25, widthCm: 1.0, heightCm: 0.42, sample: 'FATURA/DUPLICATAS' },
      { key: 'fat-quadro', label: 'FATURA', topCm: 11.51, leftCm: 0.25, widthCm: 20.57, heightCm: 0.85 },

      // CÁLCULO DO IMPOSTO
      { key: 'imp-titulo', label: 'CÁLCULO DO IMPOSTO', topCm: 12.36, leftCm: 0.25, widthCm: 5.6, heightCm: 0.42, sample: 'CÁLCULO DO IMPOSTO' },
      { key: 'imp-bc-icms', label: 'BASE DE CÁLCULO DO ICMS', topCm: 12.78, leftCm: 0.25, widthCm: 4.06, heightCm: 0.85, sample: values.v_bc_icms },
      { key: 'imp-v-icms', label: 'VALOR DO ICMS', topCm: 12.78, leftCm: 4.31, widthCm: 4.06, heightCm: 0.85, sample: values.v_icms },
      { key: 'imp-bc-st', label: 'BASE DE CÁLCULO DO ICMS ST', topCm: 12.78, leftCm: 8.37, widthCm: 4.06, heightCm: 0.85, sample: values.v_bc_st },
      { key: 'imp-v-st', label: 'VALOR DO ICMS ST', topCm: 12.78, leftCm: 12.43, widthCm: 4.06, heightCm: 0.85, sample: values.v_st },
      { key: 'imp-v-prod', label: 'VALOR TOTAL DOS PRODUTOS', topCm: 12.78, leftCm: 16.49, widthCm: 4.32, heightCm: 0.85, sample: values.v_prod },
      { key: 'imp-v-frete', label: 'VALOR DO FRETE', topCm: 13.63, leftCm: 0.25, widthCm: 3.3, heightCm: 0.85, sample: values.v_frete },
      { key: 'imp-v-seg', label: 'VALOR DO SEGURO', topCm: 13.63, leftCm: 3.55, widthCm: 3.3, heightCm: 0.85, sample: values.v_seg },
      { key: 'imp-desc', label: 'DESCONTO', topCm: 13.63, leftCm: 6.85, widthCm: 3.3, heightCm: 0.85, sample: values.v_desc },
      { key: 'imp-outros', label: 'OUTRAS DESPESAS ACESSÓRIAS', topCm: 13.63, leftCm: 10.15, widthCm: 3.3, heightCm: 0.85, sample: values.v_outros },
      { key: 'imp-v-ipi', label: 'VALOR DO IPI', topCm: 13.63, leftCm: 13.45, widthCm: 3.3, heightCm: 0.85, sample: values.v_ipi },
      { key: 'imp-v-total', label: 'VALOR TOTAL DA NOTA', topCm: 13.63, leftCm: 16.75, widthCm: 4.06, heightCm: 0.85, sample: values.v_nf },

      // TRANSPORTADOR / VOLUMES
      { key: 'tr-titulo', label: 'TRANSPORTADOR/VOLUMES TRANSPORTADOS', topCm: 14.48, leftCm: 0.25, widthCm: 5.2, heightCm: 0.42, sample: 'TRANSPORTADOR/VOLUMES' },
      { key: 'tr-razao', label: 'RAZÃO SOCIAL', topCm: 14.9, leftCm: 0.25, widthCm: 9.02, heightCm: 0.85, sample: values.tr_nome },
      { key: 'tr-frete-conta', label: 'FRETE POR CONTA DE', topCm: 14.9, leftCm: 9.27, widthCm: 2.79, heightCm: 0.85 },
      { key: 'tr-antt', label: 'CÓDIGO ANTT', topCm: 14.9, leftCm: 12.06, widthCm: 1.78, heightCm: 0.85, sample: values.tr_antt || values.tr_rntc },
      { key: 'tr-placa', label: 'PLACA DO VEÍCULO', topCm: 14.9, leftCm: 13.84, widthCm: 2.29, heightCm: 0.85, sample: values.tr_placa },
      { key: 'tr-uf-placa', label: 'UF', topCm: 14.9, leftCm: 16.13, widthCm: 0.76, heightCm: 0.85, sample: values.tr_uf },
      { key: 'tr-cnpj', label: 'CNPJ/CPF', topCm: 14.9, leftCm: 16.89, widthCm: 3.94, heightCm: 0.85, sample: values.tr_cnpj },
      { key: 'tr-end', label: 'ENDEREÇO', topCm: 15.75, leftCm: 0.25, widthCm: 9.02, heightCm: 0.85, sample: values.tr_end },
      { key: 'tr-mun', label: 'MUNICÍPIO', topCm: 15.75, leftCm: 9.27, widthCm: 6.86, heightCm: 0.85, sample: values.tr_mun },
      { key: 'tr-uf', label: 'UF', topCm: 15.75, leftCm: 16.13, widthCm: 0.76, heightCm: 0.85, sample: values.tr_uf },
      { key: 'tr-ie', label: 'INSCRIÇÃO ESTADUAL', topCm: 15.75, leftCm: 16.89, widthCm: 3.94, heightCm: 0.85, sample: values.tr_ie },
      { key: 'tr-qtd-vol', label: 'QUANTIDADE DE VOLUMES', topCm: 16.6, leftCm: 0.25, widthCm: 2.92, heightCm: 0.85, sample: values.tr_qtd },
      { key: 'tr-especie', label: 'ESPÉCIE', topCm: 16.6, leftCm: 3.17, widthCm: 3.05, heightCm: 0.85 },
      { key: 'tr-marca', label: 'MARCA', topCm: 16.6, leftCm: 6.22, widthCm: 3.05, heightCm: 0.85 },
      { key: 'tr-numeracao', label: 'NUMERAÇÃO', topCm: 16.6, leftCm: 9.27, widthCm: 4.83, heightCm: 0.85 },
      { key: 'tr-peso-bruto', label: 'PESO BRUTO', topCm: 16.6, leftCm: 14.1, widthCm: 3.43, heightCm: 0.85 },
      { key: 'tr-peso-liq', label: 'PESO LÍQUIDO', topCm: 16.6, leftCm: 17.53, widthCm: 3.3, heightCm: 0.85 },

      // PRODUTOS / SERVIÇOS (quadro grande)
      { key: 'prod-titulo', label: 'DADOS DOS PRODUTOS/SERVIÇOS', topCm: 17.45, leftCm: 0.25, widthCm: 4.0, heightCm: 0.42, sample: 'DADOS DOS PRODUTOS/SERVIÇOS' },
      { key: 'prod-quadro', label: 'QUADRO DADOS DOS PRODUTOS/SERVIÇOS', topCm: 17.87, leftCm: 0.25, widthCm: 20.57, heightCm: 6.77, kind: 'produtos' },

      // ISSQN
      { key: 'issqn-titulo', label: 'CÁLCULO DO ISSQN', topCm: 24.64, leftCm: 0.25, widthCm: 2.29, heightCm: 0.42, sample: 'CÁLCULO DO ISSQN' },
      { key: 'issqn-im', label: 'INSCRIÇÃO MUNICIPAL', topCm: 25.06, leftCm: 0.25, widthCm: 5.08, heightCm: 0.85 },
      { key: 'issqn-v-serv', label: 'VALOR TOTAL DOS SERVIÇOS', topCm: 25.06, leftCm: 5.33, widthCm: 5.08, heightCm: 0.85 },
      { key: 'issqn-bc', label: 'BASE DE CÁLCULO DO ISSQN', topCm: 25.06, leftCm: 10.41, widthCm: 5.08, heightCm: 0.85 },
      { key: 'issqn-v', label: 'VALOR DO ISSQN', topCm: 25.06, leftCm: 15.49, widthCm: 5.33, heightCm: 0.85 },

      // DADOS ADICIONAIS
      { key: 'add-titulo', label: 'DADOS ADICIONAIS', topCm: 25.91, leftCm: 0.25, widthCm: 2.29, heightCm: 0.42, sample: 'DADOS ADICIONAIS' },
      { key: 'add-info', label: 'INFORMAÇÕES COMPLEMENTARES', topCm: 26.33, leftCm: 0.25, widthCm: 12.95, heightCm: 3.07, sample: values.info_comp },
      { key: 'add-fisco', label: 'RESERVADO AO FISCO', topCm: 26.33, leftCm: 13.17, widthCm: 7.62, heightCm: 3.07, sample: values.info_fisco },
    ],
    [data, values],
  );

  return (
    <div className={`${styles.wrapper} ${variant === 'print' ? styles.wrapperPrint : ''}`}>
      {variant === 'default' ? (
        <div className="flex flex-wrap justify-content-between align-items-start gap-2 mb-2">
          <div>
            <div className="text-900 font-semibold">DANFE (teste) — A4 retrato</div>
            <div className={styles.hint}>Casca vazia com medidas em cm (eixo 0 no canto superior esquerdo).</div>
          </div>
          <button type="button" className="p-button p-component p-button-sm" onClick={() => window.print()}>
            <span className="p-button-icon p-c pi pi-print" />
            <span className="p-button-label">Imprimir</span>
          </button>
        </div>
      ) : null}

      <div className="overflow-x-auto overflow-y-hidden">
        <div className="flex justify-content-center">
          <div className={styles.page}>
            {boxes.map((b) => (
              <div key={b.key}>
                <div
                  className={`${styles.box} ${b.kind === 'barcode' ? styles.barcodeBox : ''}`}
                  style={{
                    top: cm(b.topCm),
                    left: cm(b.leftCm),
                    width: cm(b.widthCm),
                    height: cm(b.heightCm),
                    padding: b.kind === 'produtos' ? 0 : undefined,
                  }}
                >
                  {b.kind !== 'produtos' && b.kind !== 'barcode' ? (
                    <>
                      <div className={styles.boxLabel}>{b.label}</div>
                      <div className={styles.boxValue}>{b.sample ?? ''}</div>
                    </>
                  ) : null}

                  {b.kind === 'barcode' ? (
                    <svg ref={barcodeRef} className={styles.barcodeSvg} aria-label="Código de barras (chave de acesso)" />
                  ) : null}

                  {b.kind === 'produtos' ? (
                    <div className={styles.prodGrid}>
                      <div className={styles.prodVertLinesAll}>
                        {produtoCols.map((c) => (
                          <div key={c.key} className={styles.prodVert} style={{ width: cm(c.w) }} />
                        ))}
                      </div>
                      <div className={styles.prodHeaderRow}>
                        {produtoCols.map((c) => (
                          <div key={c.key} className={styles.prodCol} style={{ width: cm(c.w) }}>
                            {c.label}
                          </div>
                        ))}
                      </div>
                      <div className={styles.prodRows} />

                      {/* Dados (amostra) dentro do quadro de produtos/serviços */}
                      <div className={styles.prodDataLayer}>
                        {(data?.itens ?? []).slice(0, 9).map((it: any, i: number) => {
                          const top = (0.64 + 0.04) * i; // cm (altura linha + traço)
                          const cols = [
                            fmt(it.codigo),
                            fmt(it.descricao),
                            fmt(it.ncm),
                            fmt(it.icms?.tributacao || it.cst),
                            fmt(it.cfop),
                            fmt(it.unidade),
                            fmt(it.quantidade),
                            fmtBRL(it.valor_unitario),
                            fmtBRL(it.valor_desconto),
                            fmtBRL(it.valor_total),
                            fmtBRL(it.icms?.base_calculo || it.base_icms),
                            fmtBRL(it.icms?.valor || it.valor_icms),
                            fmtBRL(it.icms?.base_calculo_st || it.base_icms_st),
                            fmtBRL(it.icms?.valor_st || it.valor_st),
                            fmtBRL(it.valor_ipi || it.ipi?.v_ipi),
                            fmt(it.icms?.aliquota || it.aliquota_icms),
                            fmt(it.ipi?.aliquota || it.aliquota_ipi),
                          ];
                          return (
                            <div key={`${it.codigo ?? 'it'}-${i}`} className={styles.prodRow} style={{ top: cm(top) }}>
                              {produtoCols.map((c, idx) => (
                                <div key={`${c.key}-${i}`} className={styles.prodCellText} style={{ width: cm(c.w) }}>
                                  {cols[idx] ?? ''}
                                </div>
                              ))}
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  ) : null}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

