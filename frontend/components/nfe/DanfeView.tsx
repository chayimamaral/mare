import React, { useCallback, useState } from 'react';
import { TabPanel, TabView } from 'primereact/tabview';

import type { NFEDanfeView } from '../../lib/nfeDanfeClient';
import styles from '../../styles/nfe-danfe.module.css';

type Props = {
  data: NFEDanfeView;
  /** Abre a matriz A4 (MESMO dialog) para impressão; evita nova aba/guard. */
  onImprimirMatriz?: () => void;
};

const fmt = (v?: string) => (v && String(v).trim() ? String(v).trim() : '—');

function fmtBRL(v?: string): string {
  if (v == null || !String(v).trim()) return '—';
  const s = String(v).trim();
  const n = Number(s);
  if (Number.isFinite(n)) {
    return n.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  }
  return s;
}

function fmtBRLFlexible(v?: string): string {
  if (v == null || !String(v).trim()) return '—';
  const s = String(v).trim();
  const n = Number(s);
  if (!Number.isFinite(n)) return s;
  return n.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 10 });
}

function fmtQtd(v?: string): string {
  if (v == null || !String(v).trim()) return '—';
  const s = String(v).trim();
  const n = Number(s);
  if (!Number.isFinite(n)) return s;
  return n.toLocaleString('pt-BR', { minimumFractionDigits: 0, maximumFractionDigits: 4 });
}

export function DanfeView({ data, onImprimirMatriz }: Props) {
  const itens = data.itens.slice(0, 200);
  const itensLimitados = data.itens.length > itens.length;
  const [expIdx, setExpIdx] = useState<Record<number, boolean>>({});

  const toggleExp = useCallback((i: number) => {
    setExpIdx((p) => ({ ...p, [i]: !p[i] }));
  }, []);

  const id = data.identificacao;
  const em = data.emitente;
  const de = data.destinatario;
  const tot = data.totais;
  const tr = data.transporte;

  const imprimirDanfe = useCallback(() => {
    if (onImprimirMatriz) {
      onImprimirMatriz();
      return;
    }
    window.print();
  }, [onImprimirMatriz]);

  return (
    <div className={`${styles.wrapper} surface-0 border-round border-1 surface-border p-3`}>
      <div className="flex justify-content-between align-items-start flex-wrap gap-3 mb-3">
        <div>
          <h3 className="m-0">Consulta da NF-e</h3>
          <small className="text-600">Visualização conforme DANFE (abas)</small>
        </div>
        <button type="button" className="p-button p-component p-button-sm" onClick={imprimirDanfe}>
          <span className="p-button-icon p-c pi pi-print" />
          <span className="p-button-label">Imprimir</span>
        </button>
      </div>

      <fieldset className={styles.fieldset}>
        <legend>Identificação</legend>
        <div className={styles.row}>
          <div className={`${styles.cell} ${styles.w80}`}>
            <div className={styles.label}>Chave de Acesso</div>
            <strong>{fmt(id.chave)}</strong>
          </div>
          <div className={`${styles.cell} ${styles.w20}`}>
            <div className={styles.label}>Número</div>
            <strong>{fmt(id.numero)}</strong>
          </div>
        </div>
      </fieldset>

      <TabView
        className={styles.tabview}
        pt={{
          inkbar: { className: 'danfe-tabview-inkbar-hidden', style: { display: 'none' } },
          nav: { className: 'danfe-tabview-nav-no-bullets' },
        }}
      >
        <TabPanel header="NFe">
          <fieldset className={styles.fieldset}>
            <legend>Dados da NF-e</legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w10}`}>
                <div className={styles.label}>Modelo</div>
                <strong>{fmt(id.modelo)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w10}`}>
                <div className={styles.label}>Série</div>
                <strong>{fmt(id.serie)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w10}`}>
                <div className={styles.label}>Número</div>
                <strong>{fmt(id.numero)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w35}`}>
                <div className={styles.label}>Data de Emissão</div>
                <strong>{fmt(id.emissao_em)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w35}`}>
                <div className={styles.label}>Data/Hora Saída/Entrada</div>
                <strong>{fmt(id.saida_entrada_em)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w100}`}>
                <div className={styles.label}>Valor Total da Nota Fiscal</div>
                <strong>{fmtBRLFlexible(tot.valor_nota)}</strong>
              </div>
            </div>
          </fieldset>

          <fieldset className={styles.fieldset}>
            <legend>Emitente</legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>CNPJ</div>
                <strong>{fmt(em.cnpj_cpf)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w45}`}>
                <div className={styles.label}>Nome / Razão Social</div>
                <strong>{fmt(em.nome)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w20}`}>
                <div className={styles.label}>Inscrição Estadual</div>
                <strong>{fmt(em.ie)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w10}`}>
                <div className={styles.label}>UF</div>
                <strong>{fmt(em.uf)}</strong>
              </div>
            </div>
          </fieldset>

          <fieldset className={styles.fieldset}>
            <legend>Destinatário</legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>CNPJ</div>
                <strong>{fmt(de.cnpj_cpf)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w45}`}>
                <div className={styles.label}>Nome / Razão Social</div>
                <strong>{fmt(de.nome)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w20}`}>
                <div className={styles.label}>Inscrição Estadual</div>
                <strong>{fmt(de.ie)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w10}`}>
                <div className={styles.label}>UF</div>
                <strong>{fmt(de.uf)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w33}`}>
                <div className={styles.label}>Destino da Operação</div>
                <strong>{fmt(id.destino_operacao)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w33}`}>
                <div className={styles.label}>Consumidor Final</div>
                <strong>{fmt(id.consumidor_final)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w33}`}>
                <div className={styles.label}>Presença do Comprador</div>
                <strong>{fmt(id.presenca_comprador)}</strong>
              </div>
            </div>
          </fieldset>

          <fieldset className={styles.fieldset}>
            <legend>Emissão</legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Processo</div>
                <strong>{fmt(id.processo_emissao)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Versão do Processo</div>
                <strong>{fmt(id.versao_processo)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Tipo de Emissão</div>
                <strong>{fmt(id.tipo_emissao)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Finalidade</div>
                <strong>{fmt(id.finalidade)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Natureza da Operação</div>
                <strong>{fmt(id.natureza_operacao)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Tipo da Operação</div>
                <strong>{fmt(id.tipo_operacao)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Forma de Pagamento</div>
                <strong>{fmt(id.forma_pagamento)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Digest Value NF-e</div>
                <strong className={styles.digest}>{fmt(id.digest_value)}</strong>
              </div>
            </div>
          </fieldset>

          <fieldset className={styles.fieldset}>
            <legend>
              Situação atual: {fmt(id.situacao)}
              {id.ambiente ? ` (Ambiente: ${fmt(id.ambiente)})` : ''}
            </legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Eventos da NF-e</div>
                <strong>{fmt(id.evento_descricao)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Protocolo</div>
                <strong>{fmt(id.protocolo)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Data Autorização</div>
                <strong>{fmt(id.data_autorizacao)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Data Inclusão BD</div>
                <strong>{fmt(id.data_inclusao_bd)}</strong>
              </div>
            </div>
          </fieldset>

          {data.cobranca.pagamentos.length > 0 ? (
            <fieldset className={styles.fieldset}>
              <legend>Pagamentos (detalhe)</legend>
              <div className="overflow-auto">
                <table className={`${styles.table} ${styles.tableHover} w-full text-sm`}>
                  <thead>
                    <tr>
                      <th className="text-left">Forma</th>
                      <th className="text-right">Valor</th>
                      <th className="text-left">CNPJ Credenciadora</th>
                      <th className="text-left">Bandeira</th>
                      <th className="text-left">Autorização</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.cobranca.pagamentos.map((p, i) => (
                      <tr key={`${p.forma}-${i}`}>
                        <td>{fmt(p.forma)}</td>
                        <td className="text-right">{fmtBRL(p.valor)}</td>
                        <td>{fmt(p.cnpj_credenciadora)}</td>
                        <td>{fmt(p.bandeira)}</td>
                        <td>{fmt(p.autorizacao)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </fieldset>
          ) : null}
        </TabPanel>

        <TabPanel header="Emitente">
          <fieldset className={styles.fieldset}>
            <legend>Dados do emitente</legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Nome / Razão Social</div>
                <strong>{fmt(em.nome)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Nome Fantasia</div>
                <strong>{fmt(em.nome_fantasia)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>CNPJ</div>
                <strong>{fmt(em.cnpj_cpf)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Endereço</div>
                <strong>{fmt(em.endereco_completo || [em.logradouro, em.numero].filter(Boolean).join(', '))}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Bairro / Distrito</div>
                <strong>{fmt(em.bairro)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>CEP</div>
                <strong>{fmt(em.cep)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Município</div>
                <strong>{fmt(em.municipio_cod_nome || em.municipio)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Telefone</div>
                <strong>{fmt(em.telefone)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>UF</div>
                <strong>{fmt(em.uf)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>País</div>
                <strong>{fmt(em.pais_cod_nome || em.pais_nome)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Inscrição Estadual</div>
                <strong>{fmt(em.ie)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Inscrição Estadual do Substituto</div>
                <strong>{fmt(em.ie_substituto)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Inscrição Municipal</div>
                <strong>{fmt(em.im)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Município Fato Gerador do ICMS</div>
                <strong>{fmt(em.cod_mun_fato_gerador_icms)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>CNAE Fiscal</div>
                <strong>{fmt(em.cnae)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Código de Regime Tributário</div>
                <strong>{fmt(em.crt_descricao || em.crt)}</strong>
              </div>
            </div>
          </fieldset>
        </TabPanel>

        <TabPanel header="Destinatário">
          <fieldset className={styles.fieldset}>
            <legend>Dados do destinatário</legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w100}`}>
                <div className={styles.label}>Nome / Razão Social</div>
                <strong>{fmt(de.nome)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>CNPJ</div>
                <strong>{fmt(de.cnpj_cpf)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Endereço</div>
                <strong>{fmt(de.endereco_completo || [de.logradouro, de.numero].filter(Boolean).join(', '))}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Bairro / Distrito</div>
                <strong>{fmt(de.bairro)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>CEP</div>
                <strong>{fmt(de.cep)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Município</div>
                <strong>{fmt(de.municipio_cod_nome || de.municipio)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Telefone</div>
                <strong>{fmt(de.telefone)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>UF</div>
                <strong>{fmt(de.uf)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>País</div>
                <strong>{fmt(de.pais_cod_nome || de.pais_nome)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Indicador IE</div>
                <strong className="line-height-3">{fmt(de.indicador_ie_descricao || de.indicador_ie_dest)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Inscrição Estadual</div>
                <strong>{fmt(de.ie)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w100}`}>
                <div className={styles.label}>Inscrição SUFRAMA</div>
                <strong>{fmt(de.isuf)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>IM</div>
                <strong>{fmt(de.im)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>E-mail</div>
                <strong>{fmt(de.email)}</strong>
              </div>
            </div>
          </fieldset>
        </TabPanel>

        <TabPanel header="Produtos e Serviços">
          <fieldset className={styles.fieldset}>
            <legend>Dados dos produtos e serviços</legend>
            <div className="overflow-auto">
              <table className={`${styles.table} w-full text-sm`}>
                <thead>
                  <tr>
                    <th className="text-left w-3rem">Num.</th>
                    <th className="text-left">Descrição</th>
                    <th className="text-right">Qtd.</th>
                    <th className="text-left">Unidade Comercial</th>
                    <th className="text-right">Valor (R$)</th>
                    <th className="w-6rem" />
                  </tr>
                </thead>
                <tbody>
                  {itens.map((it, idx) => (
                    <React.Fragment key={`${it.codigo}-${idx}`}>
                      <tr>
                        <td>{idx + 1}</td>
                        <td>{fmt(it.descricao)}</td>
                        <td className="text-right">{fmtQtd(it.quantidade)}</td>
                        <td>{fmt(it.unidade)}</td>
                        <td className="text-right">{fmtBRLFlexible(it.valor_total)}</td>
                        <td className="text-center">
                          <button
                            type="button"
                            className="p-button p-component p-button-text p-button-sm"
                            onClick={() => toggleExp(idx)}
                          >
                            {expIdx[idx] ? 'Ocultar' : 'Detalhes'}
                          </button>
                        </td>
                      </tr>
                      {expIdx[idx] ? (
                        <tr>
                          <td colSpan={6} className="p-0 border-none">
                            <div className="p-3 surface-50 border-top-1 surface-border">
                              <div className={styles.subBlock}>
                                <div className={styles.subBlockTitle}>Detalhamento do item</div>
                                <div className={styles.row}>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Código do Produto</div>
                                    <strong>{fmt(it.codigo)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Código NCM</div>
                                    <strong>{fmt(it.ncm)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Código EX da TIPI</div>
                                    <strong>{fmt(it.extipi)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>CFOP</div>
                                    <strong>{fmt(it.cfop)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Outras Despesas Acessórias</div>
                                    <strong>{fmtBRLFlexible(it.valor_outros)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Valor do Desconto</div>
                                    <strong>{fmtBRLFlexible(it.valor_desconto)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Valor Total do Frete</div>
                                    <strong>{fmtBRLFlexible(it.valor_frete)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Valor do Seguro</div>
                                    <strong>{fmtBRLFlexible(it.valor_seguro)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w100}`}>
                                    <div className={styles.label}>Indicador de Composição do Valor Total da NF-e</div>
                                    <strong>{fmt(it.indicador_total_desc || it.indicador_total_nf)}</strong>
                                  </div>
                                </div>
                              </div>

                              <div className={styles.subBlock}>
                                <div className={styles.subBlockTitle}>Informações comerciais e tributáveis</div>
                                <div className={styles.row}>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Código EAN Comercial</div>
                                    <strong>{fmt(it.cean)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Unidade Comercial</div>
                                    <strong>{fmt(it.unidade)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Quantidade Comercial</div>
                                    <strong>{fmtQtd(it.quantidade)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Código EAN Tributável</div>
                                    <strong>{fmt(it.cean_trib)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Unidade Tributável</div>
                                    <strong>{fmt(it.u_trib)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Quantidade Tributável</div>
                                    <strong>{fmtQtd(it.q_trib)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Valor Unit. Comercializ.</div>
                                    <strong>{fmtBRLFlexible(it.valor_unitario)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Valor Unit. Tributação</div>
                                    <strong>{fmtBRLFlexible(it.v_un_trib)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Valor Aproximado dos Tributos</div>
                                    <strong>{fmtBRLFlexible(it.valor_total_tributos)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Número do Pedido Compra</div>
                                    <strong>{fmt(it.x_ped)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Item do Pedido de Compra</div>
                                    <strong>{fmt(it.n_item_ped)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Número da FCI</div>
                                    <strong>{fmt(it.n_fci)}</strong>
                                  </div>
                                </div>
                              </div>

                              <div className={styles.subBlock}>
                                <div className={styles.subBlockTitle}>ICMS normal e ST</div>
                                <div className={styles.row}>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Origem da Mercadoria</div>
                                    <strong>{fmt(it.icms?.origem)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Tributação do ICMS</div>
                                    <strong>{fmt(it.icms?.tributacao)}</strong>
                                  </div>
                                </div>
                              </div>

                              <div className={styles.subBlock}>
                                <div className={styles.subBlockTitle}>Imposto sobre Produtos Industrializados (IPI)</div>
                                <div className={styles.row}>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Classe de Enquadramento</div>
                                    <strong>{fmt(it.ipi?.cl_enq)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Código de Enquadramento</div>
                                    <strong>{fmt(it.ipi?.c_enq)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Código do Selo</div>
                                    <strong>{fmt(it.ipi?.c_selo)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>CNPJ do Produtor</div>
                                    <strong>{fmt(it.ipi?.cnpj_prod)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Qtd. Selo</div>
                                    <strong>{fmt(it.ipi?.q_selo)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>CST</div>
                                    <strong>{fmt(it.ipi?.cst)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Qtd Total Unidade Padrão</div>
                                    <strong>{fmtQtd(it.ipi?.q_unid)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Valor por Unidade</div>
                                    <strong>{fmtBRLFlexible(it.ipi?.v_unid)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w33}`}>
                                    <div className={styles.label}>Valor IPI</div>
                                    <strong>{fmtBRLFlexible(it.ipi?.v_ipi || it.valor_ipi)}</strong>
                                  </div>
                                </div>
                                <div className={`${styles.row} ${styles.rowTop}`}>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Base de Cálculo</div>
                                    <strong>{fmtBRLFlexible(it.ipi?.v_bc)}</strong>
                                  </div>
                                  <div className={`${styles.cell} ${styles.w50}`}>
                                    <div className={styles.label}>Alíquota</div>
                                    <strong>{fmt(it.ipi?.p_ipi || it.aliquota_ipi)}</strong>
                                  </div>
                                </div>
                              </div>

                              <div className={styles.subBlock}>
                                <div className={styles.subBlockTitle}>PIS / COFINS</div>
                                <div className={`${styles.cell} ${styles.w100} border-none`}>
                                  <div className={styles.label}>PIS — CST</div>
                                  <strong>{fmt(it.pis_cst)}</strong>
                                </div>
                                <div className={`${styles.cell} ${styles.w100} ${styles.rowTop} border-none`}>
                                  <div className={styles.label}>COFINS — CST</div>
                                  <strong>{fmt(it.cofins_cst)}</strong>
                                </div>
                              </div>
                            </div>
                          </td>
                        </tr>
                      ) : null}
                    </React.Fragment>
                  ))}
                </tbody>
              </table>
            </div>
            {itensLimitados ? <small className="text-600 block mt-2">Pré-visualização limitada a 200 itens.</small> : null}
          </fieldset>
        </TabPanel>

        <TabPanel header="Totais">
          <fieldset className={styles.fieldset}>
            <legend>Totais — ICMS e demais</legend>
            <div className={styles.grid4}>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Base de Cálculo ICMS</div>
                <div className={styles.totValue}>{fmtBRL(tot.base_icms)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor do ICMS</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_icms)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor ICMS Desonerado</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_icms_desonerado)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Base de Cálculo ICMS ST</div>
                <div className={styles.totValue}>{fmtBRL(tot.base_icms_st)}</div>
              </div>

              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor ICMS Substit.</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_st)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor Total Produtos</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_produtos)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor do Frete</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_frete)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor do Seguro</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_seguro)}</div>
              </div>

              <div className={styles.totCell}>
                <div className={styles.totLabel}>Outras Desp. Aces.</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_outros)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor Total do IPI</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_ipi)}</div>
              </div>
              <div className={`${styles.totCell} ${styles.totDestaque}`}>
                <div className={styles.totLabel}>Valor Total da NFe</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_nota)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor Total dos Descontos</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_desconto)}</div>
              </div>

              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor Total do II</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_ii)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor do PIS</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_pis)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor da COFINS</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_cofins)}</div>
              </div>
              <div className={styles.totCell}>
                <div className={styles.totLabel}>Valor Aproximado Tributos</div>
                <div className={styles.totValue}>{fmtBRL(tot.valor_total_tributos)}</div>
              </div>
            </div>
          </fieldset>
        </TabPanel>

        <TabPanel header="Transporte">
          <fieldset className={styles.fieldset}>
            <legend>Dados do transporte</legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w100}`}>
                <div className={styles.label}>Modalidade do Frete</div>
                <strong>{fmt(tr.modalidade)}</strong>
              </div>
            </div>
          </fieldset>

          <fieldset className={styles.fieldset}>
            <legend>Transportador</legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>CNPJ</div>
                <strong>{fmt(tr.cnpj_cpf)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w50}`}>
                <div className={styles.label}>Razão Social / Nome</div>
                <strong>{fmt(tr.transportador)}</strong>
              </div>
            </div>
            <div className={`${styles.transGrid12} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.transCol3}`}>
                <div className={styles.label}>Inscrição Estadual</div>
                <strong>{fmt(tr.ie)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.transCol6}`}>
                <div className={styles.label}>Endereço Completo</div>
                <strong>{fmt(tr.endereco)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.transCol3}`}>
                <div className={styles.label}>Município</div>
                <strong>{fmt(tr.municipio)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>UF</div>
                <strong>{fmt(tr.uf)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Placa</div>
                <strong>{fmt(tr.placa)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>RNTC</div>
                <strong>{fmt(tr.rntc)}</strong>
              </div>
              <div className={`${styles.cell} ${styles.w25}`}>
                <div className={styles.label}>Qtd. volumes (soma)</div>
                <strong>{fmt(tr.quantidade_volumes)}</strong>
              </div>
            </div>
          </fieldset>

          <fieldset className={styles.fieldset}>
            <legend>Volumes</legend>
            {tr.volumes.length === 0 ? (
              <p className={`m-0 ${styles.cellMuted}`}>Nenhum volume informado.</p>
            ) : (
              tr.volumes.map((v, i) => (
                <div key={`${v.numero}-${i}`} className={styles.volCard}>
                  <div className={styles.volTitle}>Volume {i + 1}</div>
                  <div className={styles.row}>
                    <div className={`${styles.cell} ${styles.w33}`}>
                      <div className={styles.label}>Quantidade</div>
                      <strong>{fmt(v.quantidade)}</strong>
                    </div>
                    <div className={`${styles.cell} ${styles.w33}`}>
                      <div className={styles.label}>Espécie</div>
                      <strong>{fmt(v.especie)}</strong>
                    </div>
                    <div className={`${styles.cell} ${styles.w33}`}>
                      <div className={styles.label}>Marca dos Volumes</div>
                      <strong>{fmt(v.marca)}</strong>
                    </div>
                  </div>
                  <div className={`${styles.row} ${styles.rowTop}`}>
                    <div className={`${styles.cell} ${styles.w33}`}>
                      <div className={styles.label}>Numeração</div>
                      <strong>{fmt(v.numero)}</strong>
                    </div>
                    <div className={`${styles.cell} ${styles.w33}`}>
                      <div className={styles.label}>Peso Líquido</div>
                      <strong>{fmt(v.peso_liquido)}</strong>
                    </div>
                    <div className={`${styles.cell} ${styles.w33}`}>
                      <div className={styles.label}>Peso Bruto</div>
                      <strong>{fmt(v.peso_bruto)}</strong>
                    </div>
                  </div>
                </div>
              ))
            )}
          </fieldset>
        </TabPanel>

        <TabPanel header="Cobrança">
          <fieldset className={styles.fieldset}>
            <legend>Dados de cobrança — Duplicatas</legend>
            {data.cobranca.duplicatas.length === 0 ? (
              <p className="m-0 text-600">Nenhuma duplicata informada para esta nota.</p>
            ) : (
              <div className="overflow-auto">
                <table className={`${styles.table} ${styles.tableHover} w-full text-sm`}>
                  <thead>
                    <tr>
                      <th className="text-left">Número</th>
                      <th className="text-left">Vencimento</th>
                      <th className="text-right">Valor</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.cobranca.duplicatas.map((d, idx) => (
                      <tr key={`${d.numero}-${idx}`}>
                        <td>{fmt(d.numero)}</td>
                        <td>{fmt(d.vencimento)}</td>
                        <td className="text-right">{fmtBRL(d.valor)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </fieldset>
        </TabPanel>

        <TabPanel header="Informações Adicionais">
          <fieldset className={styles.fieldset}>
            <legend>Informações adicionais</legend>
            <div className={styles.row}>
              <div className={`${styles.cell} ${styles.w100}`}>
                <div className={styles.label}>Formato de Impressão DANFE</div>
                <strong>{fmt(data.adicionais.tp_imp)}</strong>
              </div>
            </div>
            <div className={`${styles.row} ${styles.rowTop}`}>
              <div className={`${styles.cell} ${styles.w100}`}>
                <div className={styles.label}>Informações complementares de interesse do contribuinte</div>
                {!(data.adicionais.informacoes_complementares ?? '').trim() ? (
                  <span className="text-500">—</span>
                ) : (
                  <div className={`${styles.preWrap} text-900`}>{data.adicionais.informacoes_complementares}</div>
                )}
              </div>
            </div>
            {data.adicionais.informacoes_fisco?.trim() ? (
              <div className={`${styles.row} ${styles.rowTop}`}>
                <div className={`${styles.cell} ${styles.w100}`}>
                  <div className={styles.label}>Informações adicionais de interesse do Fisco</div>
                  <div className={styles.preWrap}>
                    <strong>{data.adicionais.informacoes_fisco}</strong>
                  </div>
                </div>
              </div>
            ) : null}
          </fieldset>
        </TabPanel>
      </TabView>
    </div>
  );
}
