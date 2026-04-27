import React from 'react';
import { TabPanel, TabView } from 'primereact/tabview';

import type { NFEDanfeView } from '../../lib/nfeDanfeClient';
import styles from '../../styles/nfe-danfe.module.css';

type Props = {
    data: NFEDanfeView;
};

const fmt = (v?: string) => (v && v.trim() ? v : '—');

export function DanfeView({ data }: Props) {
    const itens = data.itens.slice(0, 200);
    const itensLimitados = data.itens.length > itens.length;

    return (
        <div className={`${styles.wrapper} surface-0 border-round border-1 surface-border p-3`}>
            <div className="flex justify-content-between align-items-start flex-wrap gap-3 mb-3">
                <div>
                    <h3 className="m-0">Consulta da NF-e</h3>
                    <small className="text-600">Visualização da DANFE em abas</small>
                </div>
                <button type="button" className="p-button p-component p-button-sm" onClick={() => window.print()}>
                    <span className="p-button-icon p-c pi pi-print" />
                    <span className="p-button-label">Imprimir</span>
                </button>
            </div>

            <fieldset className={styles.fieldset}>
                <legend>Dados Gerais</legend>
                <div className={styles.row}>
                    <div className={`${styles.cell} ${styles.w80}`}>
                        <div className={styles.label}>Chave de Acesso</div>
                        <strong>{fmt(data.identificacao.chave)}</strong>
                    </div>
                    <div className={`${styles.cell} ${styles.w20}`}>
                        <div className={styles.label}>Número</div>
                        <strong>{fmt(data.identificacao.numero)}</strong>
                    </div>
                </div>
            </fieldset>

            <TabView className={styles.tabview} pt={{ inkbar: { style: { display: 'none' } } }}>
                <TabPanel header="NFe">
                    <fieldset className={styles.fieldset}>
                        <legend>Dados da NF-e</legend>
                        <div className={styles.row}>
                            <div className={`${styles.cell} ${styles.w10}`}><div className={styles.label}>Modelo</div><strong>{fmt(data.identificacao.modelo)}</strong></div>
                            <div className={`${styles.cell} ${styles.w10}`}><div className={styles.label}>Série</div><strong>{fmt(data.identificacao.serie)}</strong></div>
                            <div className={`${styles.cell} ${styles.w10}`}><div className={styles.label}>Número</div><strong>{fmt(data.identificacao.numero)}</strong></div>
                            <div className={`${styles.cell} ${styles.w35}`}><div className={styles.label}>Data de Emissão</div><strong>{fmt(data.identificacao.emissao_em)}</strong></div>
                            <div className={`${styles.cell} ${styles.w35}`}><div className={styles.label}>Situação</div><strong>{fmt(data.identificacao.situacao)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w50}`}><div className={styles.label}>Data/Hora de Saída ou Entrada</div><strong>{fmt(data.identificacao.saida_entrada_em)}</strong></div>
                            <div className={`${styles.cell} ${styles.w50}`}><div className={styles.label}>Valor Total da Nota Fiscal</div><strong>{fmt(data.totais.valor_nota)}</strong></div>
                        </div>
                    </fieldset>
                    <fieldset className={styles.fieldset}>
                        <legend>Emissão</legend>
                        <div className={styles.row}>
                            <div className={`${styles.cell} ${styles.w40}`}><div className={styles.label}>Natureza da Operação</div><strong>{fmt(data.identificacao.natureza_operacao)}</strong></div>
                            <div className={`${styles.cell} ${styles.w20}`}><div className={styles.label}>Ambiente</div><strong>{fmt(data.identificacao.ambiente)}</strong></div>
                            <div className={`${styles.cell} ${styles.w40}`}><div className={styles.label}>Protocolo</div><strong>{fmt(data.identificacao.protocolo)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Tipo da Operação</div><strong>{fmt(data.identificacao.tipo_operacao)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Destino da Operação</div><strong>{fmt(data.identificacao.destino_operacao)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Consumidor Final</div><strong>{fmt(data.identificacao.consumidor_final)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Presença do Comprador</div><strong>{fmt(data.identificacao.presenca_comprador)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Processo de Emissão</div><strong>{fmt(data.identificacao.processo_emissao)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Versão do Processo</div><strong>{fmt(data.identificacao.versao_processo)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Finalidade</div><strong>{fmt(data.identificacao.finalidade)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Forma de Pagamento</div><strong>{fmt(data.identificacao.forma_pagamento)}</strong></div>
                        </div>
                    </fieldset>
                    <fieldset className={styles.fieldset}>
                        <legend>Situação Atual: {fmt(data.identificacao.situacao)} (Ambiente de autorização: {fmt(data.identificacao.ambiente)})</legend>
                        <div className={styles.row}>
                            <div className={`${styles.cell} ${styles.w50}`}>
                                <div className={styles.label}>Eventos da NF-e</div>
                                <strong>{fmt(data.identificacao.evento_descricao)}</strong>
                            </div>
                            <div className={`${styles.cell} ${styles.w20}`}>
                                <div className={styles.label}>Protocolo</div>
                                <strong>{fmt(data.identificacao.protocolo)}</strong>
                            </div>
                            <div className={`${styles.cell} ${styles.w20}`}>
                                <div className={styles.label}>Data Autorização</div>
                                <strong>{fmt(data.identificacao.data_autorizacao)}</strong>
                            </div>
                            <div className={`${styles.cell} ${styles.w10}`}>
                                <div className={styles.label}>cStat</div>
                                <strong>{fmt(data.identificacao.codigo_status)}</strong>
                            </div>
                        </div>
                    </fieldset>
                </TabPanel>

                <TabPanel header="Emitente">
                    <fieldset className={styles.fieldset}>
                        <legend>Emitente</legend>
                        <div className={styles.row}>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>CNPJ/CPF</div><strong>{fmt(data.emitente.cnpj_cpf)}</strong></div>
                            <div className={`${styles.cell} ${styles.w45}`}><div className={styles.label}>Nome / Razão Social</div><strong>{fmt(data.emitente.nome)}</strong></div>
                            <div className={`${styles.cell} ${styles.w20}`}><div className={styles.label}>Inscrição Estadual</div><strong>{fmt(data.emitente.ie)}</strong></div>
                            <div className={`${styles.cell} ${styles.w10}`}><div className={styles.label}>UF</div><strong>{fmt(data.emitente.uf)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w70}`}><div className={styles.label}>Endereço</div><strong>{fmt(data.emitente.logradouro)}, {fmt(data.emitente.numero)} - {fmt(data.emitente.bairro)}</strong></div>
                            <div className={`${styles.cell} ${styles.w20}`}><div className={styles.label}>Município</div><strong>{fmt(data.emitente.municipio)}</strong></div>
                            <div className={`${styles.cell} ${styles.w10}`}><div className={styles.label}>CEP</div><strong>{fmt(data.emitente.cep)}</strong></div>
                        </div>
                    </fieldset>
                </TabPanel>

                <TabPanel header="Destinatário">
                    <fieldset className={styles.fieldset}>
                        <legend>Destinatário</legend>
                        <div className={styles.row}>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>CNPJ/CPF</div><strong>{fmt(data.destinatario.cnpj_cpf)}</strong></div>
                            <div className={`${styles.cell} ${styles.w45}`}><div className={styles.label}>Nome / Razão Social</div><strong>{fmt(data.destinatario.nome)}</strong></div>
                            <div className={`${styles.cell} ${styles.w20}`}><div className={styles.label}>Inscrição Estadual</div><strong>{fmt(data.destinatario.ie)}</strong></div>
                            <div className={`${styles.cell} ${styles.w10}`}><div className={styles.label}>UF</div><strong>{fmt(data.destinatario.uf)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Indicador IE Dest.</div><strong>{fmt(data.destinatario.indicador_ie_dest)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w70}`}><div className={styles.label}>Endereço</div><strong>{fmt(data.destinatario.logradouro)}, {fmt(data.destinatario.numero)} - {fmt(data.destinatario.bairro)}</strong></div>
                            <div className={`${styles.cell} ${styles.w20}`}><div className={styles.label}>Município</div><strong>{fmt(data.destinatario.municipio)}</strong></div>
                            <div className={`${styles.cell} ${styles.w10}`}><div className={styles.label}>CEP</div><strong>{fmt(data.destinatario.cep)}</strong></div>
                        </div>
                    </fieldset>
                </TabPanel>

                <TabPanel header="Produtos e Serviços">
                    <div className={`${styles.box} overflow-auto`}>
                        <table className={`${styles.table} w-full text-sm`}>
                            <thead>
                                <tr>
                                    <th className="text-left">Código</th>
                                    <th className="text-left">Descrição</th>
                                    <th className="text-left">NCM</th>
                                    <th className="text-left">EAN</th>
                                    <th className="text-left">CFOP</th>
                                    <th className="text-right">Qtd</th>
                                    <th className="text-left">Un</th>
                                    <th className="text-right">Vlr Unit</th>
                                    <th className="text-right">Vlr Total</th>
                                    <th className="text-right">Desc.</th>
                                    <th className="text-right">Frete</th>
                                    <th className="text-right">Outros</th>
                                </tr>
                            </thead>
                            <tbody>
                                {itens.map((it, idx) => (
                                    <tr key={`${it.codigo}-${idx}`}>
                                        <td>{fmt(it.codigo)}</td>
                                        <td>{fmt(it.descricao)}</td>
                                        <td>{fmt(it.ncm)}</td>
                                        <td>{fmt(it.cean)}</td>
                                        <td>{fmt(it.cfop)}</td>
                                        <td className="text-right">{fmt(it.quantidade)}</td>
                                        <td>{fmt(it.unidade)}</td>
                                        <td className="text-right">{fmt(it.valor_unitario)}</td>
                                        <td className="text-right">{fmt(it.valor_total)}</td>
                                        <td className="text-right">{fmt(it.valor_desconto)}</td>
                                        <td className="text-right">{fmt(it.valor_frete)}</td>
                                        <td className="text-right">{fmt(it.valor_outros)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        {itensLimitados ? <small className="text-600 block mt-2">Pré-visualização limitada a 200 itens.</small> : null}
                    </div>
                </TabPanel>

                <TabPanel header="Totais">
                    <fieldset className={styles.fieldset}>
                        <legend>Totais</legend>
                        <div className={styles.row}>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Base ICMS</div><strong>{fmt(data.totais.base_icms)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>ICMS</div><strong>{fmt(data.totais.valor_icms)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>IPI</div><strong>{fmt(data.totais.valor_ipi)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>PIS</div><strong>{fmt(data.totais.valor_pis)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>ICMS Desonerado</div><strong>{fmt(data.totais.valor_icms_desonerado)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Base ICMS ST</div><strong>{fmt(data.totais.base_icms_st)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>COFINS</div><strong>{fmt(data.totais.valor_cofins)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Frete</div><strong>{fmt(data.totais.valor_frete)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w33}`}><div className={styles.label}>ICMS ST</div><strong>{fmt(data.totais.valor_st)}</strong></div>
                            <div className={`${styles.cell} ${styles.w33}`}><div className={styles.label}>Imposto de Importação (II)</div><strong>{fmt(data.totais.valor_ii)}</strong></div>
                            <div className={`${styles.cell} ${styles.w33}`}><div className={styles.label}>Outras Despesas</div><strong>{fmt(data.totais.valor_outros)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w33}`}><div className={styles.label}>Desconto</div><strong>{fmt(data.totais.valor_desconto)}</strong></div>
                            <div className={`${styles.cell} ${styles.w33}`}><div className={styles.label}>Valor Total dos Tributos</div><strong>{fmt(data.totais.valor_total_tributos)}</strong></div>
                            <div className={`${styles.cell} ${styles.w33}`}><div className={styles.label}>Valor Total NF-e</div><strong>{fmt(data.totais.valor_nota)}</strong></div>
                        </div>
                    </fieldset>
                </TabPanel>

                <TabPanel header="Transporte">
                    <fieldset className={styles.fieldset}>
                        <legend>Transporte</legend>
                        <div className={styles.row}>
                            <div className={`${styles.cell} ${styles.w20}`}><div className={styles.label}>Modalidade</div><strong>{fmt(data.transporte.modalidade)}</strong></div>
                            <div className={`${styles.cell} ${styles.w40}`}><div className={styles.label}>Transportador</div><strong>{fmt(data.transporte.transportador)}</strong></div>
                            <div className={`${styles.cell} ${styles.w20}`}><div className={styles.label}>CNPJ/CPF</div><strong>{fmt(data.transporte.cnpj_cpf)}</strong></div>
                            <div className={`${styles.cell} ${styles.w10}`}><div className={styles.label}>Placa</div><strong>{fmt(data.transporte.placa)}</strong></div>
                            <div className={`${styles.cell} ${styles.w10}`}><div className={styles.label}>UF</div><strong>{fmt(data.transporte.uf)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Inscrição Estadual</div><strong>{fmt(data.transporte.ie)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>RNTC</div><strong>{fmt(data.transporte.rntc)}</strong></div>
                            <div className={`${styles.cell} ${styles.w50}`}><div className={styles.label}>Endereço / Município</div><strong>{fmt(data.transporte.endereco)} - {fmt(data.transporte.municipio)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w100}`}><div className={styles.label}>Quantidade de Volumes</div><strong>{fmt(data.transporte.quantidade_volumes)}</strong></div>
                        </div>
                        {data.transporte.volumes.length > 0 ? (
                            <div className={`${styles.row} ${styles.rowTop}`}>
                                <div className={`${styles.cell} ${styles.w100} overflow-auto`}>
                                    <div className={styles.label}>Volumes</div>
                                    <table className={`${styles.table} w-full text-sm`}>
                                        <thead>
                                            <tr>
                                                <th>Qtd</th><th>Espécie</th><th>Marca</th><th>Numeração</th><th>Peso Líquido</th><th>Peso Bruto</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {data.transporte.volumes.map((v, i) => (
                                                <tr key={`${v.numero}-${i}`}>
                                                    <td>{fmt(v.quantidade)}</td>
                                                    <td>{fmt(v.especie)}</td>
                                                    <td>{fmt(v.marca)}</td>
                                                    <td>{fmt(v.numero)}</td>
                                                    <td>{fmt(v.peso_liquido)}</td>
                                                    <td>{fmt(v.peso_bruto)}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        ) : null}
                    </fieldset>
                </TabPanel>

                <TabPanel header="Cobrança">
                    <fieldset className={styles.fieldset}>
                        <legend>Cobrança</legend>
                        <div className={styles.row}>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Nº Fatura</div><strong>{fmt(data.cobranca.numero_fatura)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Valor Original</div><strong>{fmt(data.cobranca.valor_original)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Valor Desconto</div><strong>{fmt(data.cobranca.valor_desconto)}</strong></div>
                            <div className={`${styles.cell} ${styles.w25}`}><div className={styles.label}>Valor Líquido</div><strong>{fmt(data.cobranca.valor_liquido)}</strong></div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w100}`}>
                                <div className={styles.label}>Duplicatas</div>
                                {data.cobranca.duplicatas.length > 0 ? (
                                    <ul className="m-0 pl-3">
                                        {data.cobranca.duplicatas.map((d, idx) => (
                                            <li key={`${d.numero}-${idx}`}>Nº {fmt(d.numero)} | Venc.: {fmt(d.vencimento)} | Valor: {fmt(d.valor)}</li>
                                        ))}
                                    </ul>
                                ) : <strong>—</strong>}
                            </div>
                        </div>
                        <div className={`${styles.row} ${styles.rowTop}`}>
                            <div className={`${styles.cell} ${styles.w100} overflow-auto`}>
                                <div className={styles.label}>Formas de Pagamento</div>
                                {data.cobranca.pagamentos.length > 0 ? (
                                    <table className={`${styles.table} w-full text-sm`}>
                                        <thead>
                                            <tr>
                                                <th>Forma</th><th>Valor</th><th>CNPJ Credenciadora</th><th>Bandeira</th><th>Autorização</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {data.cobranca.pagamentos.map((p, i) => (
                                                <tr key={`${p.forma}-${i}`}>
                                                    <td>{fmt(p.forma)}</td>
                                                    <td>{fmt(p.valor)}</td>
                                                    <td>{fmt(p.cnpj_credenciadora)}</td>
                                                    <td>{fmt(p.bandeira)}</td>
                                                    <td>{fmt(p.autorizacao)}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                ) : <strong>—</strong>}
                            </div>
                        </div>
                    </fieldset>
                </TabPanel>

                <TabPanel header="Informações Adicionais">
                    <fieldset className={styles.fieldset}>
                        <legend>Informações Adicionais</legend>
                        <div className={`${styles.cell} ${styles.w100}`}>
                            <div className={styles.label}>Informações complementares</div>
                            <strong>{fmt(data.adicionais.informacoes_complementares)}</strong>
                        </div>
                        <div className={`${styles.cell} ${styles.w100} ${styles.rowTop}`}>
                            <div className={styles.label}>Informações adicionais do Fisco</div>
                            <strong>{fmt(data.adicionais.informacoes_fisco)}</strong>
                        </div>
                    </fieldset>
                </TabPanel>
            </TabView>
        </div>
    );
}
