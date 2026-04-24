import { useQuery } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { DataTable, DataTablePageEvent, DataTableSortEvent } from 'primereact/datatable';
import { Dialog } from 'primereact/dialog';
import { Dropdown } from 'primereact/dropdown';
import { InputText } from 'primereact/inputtext';
import { Toast } from 'primereact/toast';
import { Toolbar } from 'primereact/toolbar';
import { useRouter } from 'next/router';
import React, { useMemo, useRef, useState } from 'react';

import api from '../../components/api/apiClient';
import { useRouteClientGuard } from '../../components/hooks/useClientGuards';

export type NFEGestaoRow = {
    id: string;
    chave_nfe: string;
    tipo_arquivo: string;
    numero_nfe: string;
    razao_social_emitente: string;
    cnpj_emitente: string;
    data_emissao?: string | null;
    cnpj_destinatario: string;
    valor_total?: number | null;
    data_download: string;
};

const TIPO_ARQUIVO_OPTIONS = [
    { label: 'Todos', value: '' },
    { label: 'NF-e', value: 'NF-e' },
    { label: 'NFC-e', value: 'NFC-e' },
    { label: 'CT-e', value: 'CT-e' },
    { label: 'CT-e OS', value: 'CT-e OS' },
    { label: 'NFS-e Nacional', value: 'NFS-e Nacional' },
    { label: 'Outro', value: 'Outro' },
];

function formatDateBR(iso?: string | null): string {
    if (!iso) {
        return '—';
    }
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) {
        return iso;
    }
    return d.toLocaleDateString('pt-BR', { timeZone: 'UTC' });
}

function formatDateTimeBR(iso?: string | null): string {
    if (!iso) {
        return '—';
    }
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) {
        return iso;
    }
    return d.toLocaleString('pt-BR');
}

function formatValor(v?: number | null): string {
    if (v == null || Number.isNaN(v)) {
        return '—';
    }
    return v.toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
}

function formatCNPJ(v: string): string {
    const d = String(v ?? '').replace(/\D/g, '');
    if (d.length === 14) {
        return d.replace(/^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$/, '$1.$2.$3/$4-$5');
    }
    if (d.length === 11) {
        return d.replace(/^(\d{3})(\d{3})(\d{3})(\d{2})$/, '$1.$2.$3-$4');
    }
    return v || '—';
}

export default function NFEManutencaoPage() {
    useRouteClientGuard();
    const router = useRouter();
    const toast = useRef<Toast>(null);

    const [tipoArquivo, setTipoArquivo] = useState('');
    const [emissaoIni, setEmissaoIni] = useState('');
    const [emissaoFim, setEmissaoFim] = useState('');
    const [cnpjEmitente, setCnpjEmitente] = useState('');
    const [cnpjDestinatario, setCnpjDestinatario] = useState('');

    const [first, setFirst] = useState(0);
    const [rows, setRows] = useState(20);
    const [sortField, setSortField] = useState('data_download');
    const [sortOrder, setSortOrder] = useState<-1 | 1>(-1);

    const [detalhe, setDetalhe] = useState<NFEGestaoRow | null>(null);

    const queryKey = useMemo(
        () => ({
            first,
            rows,
            sortField,
            sortOrder,
            tipoArquivo,
            emissaoIni,
            emissaoFim,
            cnpjEmitente,
            cnpjDestinatario,
        }),
        [first, rows, sortField, sortOrder, tipoArquivo, emissaoIni, emissaoFim, cnpjEmitente, cnpjDestinatario],
    );

    const { data, isFetching, refetch } = useQuery({
        queryKey: ['nfe-gestao', queryKey],
        queryFn: async () => {
            const { data: res } = await api.get<{ items: NFEGestaoRow[]; totalRecords: number }>('/api/serpro/nfe/gestao', {
                params: {
                    first: queryKey.first,
                    rows: queryKey.rows,
                    sortField: queryKey.sortField,
                    sortOrder: queryKey.sortOrder,
                    tipo_arquivo: queryKey.tipoArquivo || undefined,
                    emissao_ini: queryKey.emissaoIni || undefined,
                    emissao_fim: queryKey.emissaoFim || undefined,
                    cnpj_emitente: queryKey.cnpjEmitente || undefined,
                    cnpj_destinatario: queryKey.cnpjDestinatario || undefined,
                },
            });
            return res;
        },
    });

    const onPage = (e: DataTablePageEvent) => {
        setFirst(e.first);
        setRows(e.rows);
    };

    const onSort = (e: DataTableSortEvent) => {
        setSortField(String(e.sortField ?? 'data_download'));
        setSortOrder(e.sortOrder === 1 ? 1 : -1);
        setFirst(0);
    };

    const exportarExcel = async () => {
        try {
            const { data: res } = await api.get<{ items: NFEGestaoRow[] }>('/api/serpro/nfe/gestao', {
                params: {
                    first: 0,
                    rows: 3000,
                    sortField: queryKey.sortField,
                    sortOrder: queryKey.sortOrder,
                    tipo_arquivo: queryKey.tipoArquivo || undefined,
                    emissao_ini: queryKey.emissaoIni || undefined,
                    emissao_fim: queryKey.emissaoFim || undefined,
                    cnpj_emitente: queryKey.cnpjEmitente || undefined,
                    cnpj_destinatario: queryKey.cnpjDestinatario || undefined,
                },
            });
            const list = res?.items ?? [];
            const sep = ';';
            const header = [
                'Tipo',
                'Numero',
                'Razao social emitente',
                'CNPJ emitente',
                'Data emissao',
                'CNPJ destinatario',
                'Valor',
                'Data download',
                'Chave',
            ];
            const lines = list.map((r) =>
                [
                    r.tipo_arquivo,
                    r.numero_nfe,
                    `"${(r.razao_social_emitente ?? '').replace(/"/g, '""')}"`,
                    r.cnpj_emitente,
                    formatDateBR(r.data_emissao),
                    r.cnpj_destinatario,
                    r.valor_total != null ? String(r.valor_total).replace('.', ',') : '',
                    formatDateTimeBR(r.data_download),
                    r.chave_nfe,
                ].join(sep),
            );
            const bom = '\ufeff';
            const csv = bom + header.join(sep) + '\n' + lines.join('\n');
            const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'manutencao_nfe.csv';
            a.click();
            URL.revokeObjectURL(url);
            toast.current?.show({ severity: 'success', summary: 'Exportação', detail: `${list.length} registro(s) exportado(s).`, life: 3000 });
        } catch (e: unknown) {
            const msg =
                typeof e === 'object' && e && 'response' in e
                    ? String((e as { response?: { data?: { error?: string } } }).response?.data?.error ?? 'Falha ao exportar')
                    : 'Falha ao exportar';
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 5000 });
        }
    };

    const leftToolbar = (
        <Button
            label="Sincronização de NFe"
            icon="pi pi-cloud-download"
            type="button"
            outlined
            onClick={() => void router.push('/nfe/sincronizacao')}
        />
    );

    const rightToolbar = (
        <Button label="Exportar para Excel" icon="pi pi-file-excel" severity="help" type="button" onClick={() => void exportarExcel()} />
    );

    const detalhesBody = (row: NFEGestaoRow) => (
        <Button icon="pi pi-search" rounded severity="info" type="button" aria-label="Detalhes" onClick={() => setDetalhe(row)} />
    );

    return (
        <div className="grid crud-demo">
            <div className="col-12">
                <div className="card">
                    <Toast ref={toast} />
                    <div className="mb-4">
                        <p className="text-600 m-0 text-sm mb-3">
                            Notas fiscais baixadas por consulta SERPRO neste tenant. Os campos são preenchidos automaticamente na baixa.
                        </p>
                        <div
                            className="nfe-gestao-filtros"
                            style={{
                                display: 'grid',
                                gridTemplateColumns: 'repeat(auto-fit, minmax(14rem, 1fr))',
                                gap: '1rem',
                            }}
                        >
                            <div className="field mb-0 min-w-0">
                                <label htmlFor="filtroTipoNfe" className="text-sm text-600 mb-2 block">
                                    Tipo de arquivo
                                </label>
                                <Dropdown
                                    id="filtroTipoNfe"
                                    value={tipoArquivo}
                                    options={TIPO_ARQUIVO_OPTIONS}
                                    onChange={(e) => {
                                        setTipoArquivo(e.value ?? '');
                                        setFirst(0);
                                    }}
                                    className="w-full"
                                />
                            </div>
                            <div className="field mb-0 min-w-0">
                                <label htmlFor="emissaoIni" className="text-sm text-600 mb-2 block">
                                    Data de emissão (inicial)
                                </label>
                                <input
                                    id="emissaoIni"
                                    type="date"
                                    className="p-inputtext p-component w-full"
                                    value={emissaoIni}
                                    onChange={(e) => {
                                        setEmissaoIni(e.target.value);
                                        setFirst(0);
                                    }}
                                />
                            </div>
                            <div className="field mb-0 min-w-0">
                                <label htmlFor="emissaoFim" className="text-sm text-600 mb-2 block">
                                    Data de emissão (final)
                                </label>
                                <input
                                    id="emissaoFim"
                                    type="date"
                                    className="p-inputtext p-component w-full"
                                    value={emissaoFim}
                                    onChange={(e) => {
                                        setEmissaoFim(e.target.value);
                                        setFirst(0);
                                    }}
                                />
                            </div>
                            <div className="field mb-0 min-w-0">
                                <label htmlFor="cnpjEmit" className="text-sm text-600 mb-2 block">
                                    CNPJ emitente
                                </label>
                                <InputText
                                    id="cnpjEmit"
                                    className="w-full"
                                    value={cnpjEmitente}
                                    onChange={(e) => {
                                        setCnpjEmitente(e.target.value);
                                        setFirst(0);
                                    }}
                                    placeholder="Somente números ou parcial"
                                />
                            </div>
                            <div className="field mb-0 min-w-0">
                                <label htmlFor="cnpjDest" className="text-sm text-600 mb-2 block">
                                    CNPJ destinatário
                                </label>
                                <InputText
                                    id="cnpjDest"
                                    className="w-full"
                                    value={cnpjDestinatario}
                                    onChange={(e) => {
                                        setCnpjDestinatario(e.target.value);
                                        setFirst(0);
                                    }}
                                    placeholder="Somente números ou parcial"
                                />
                            </div>
                        </div>
                    </div>

                    <Toolbar className="mb-4" left={leftToolbar} right={rightToolbar} />

                    <div className="flex flex-column md:flex-row md:justify-content-between md:align-items-center mb-3">
                        <div>
                            <h5 className="m-0">Manutenção de NFe</h5>
                            <p className="m-0 mt-1 text-600 text-sm">Listagem alinhada ao módulo de empresas; detalhes no padrão de ações de Municípios.</p>
                        </div>
                    </div>

                    <DataTable
                        value={data?.items ?? []}
                        lazy
                        loading={isFetching}
                        dataKey="id"
                        paginator
                        rows={rows}
                        first={first}
                        totalRecords={data?.totalRecords ?? 0}
                        rowsPerPageOptions={[10, 20, 50]}
                        onPage={onPage}
                        onSort={onSort}
                        sortField={sortField}
                        sortOrder={sortOrder}
                        emptyMessage="Nenhuma NF-e encontrada com os filtros atuais."
                        stripedRows
                        size="small"
                    >
                        <Column field="tipo_arquivo" header="Tipo" sortable style={{ minWidth: '8rem' }} />
                        <Column field="numero_nfe" header="Número" sortable style={{ minWidth: '7rem' }} />
                        <Column field="razao_social_emitente" header="Razão social emitente" sortable style={{ minWidth: '14rem' }} />
                        <Column
                            field="cnpj_emitente"
                            header="CNPJ emitente"
                            sortable
                            style={{ minWidth: '11rem' }}
                            body={(r: NFEGestaoRow) => formatCNPJ(r.cnpj_emitente)}
                        />
                        <Column
                            field="data_emissao"
                            header="Data emissão"
                            sortable
                            style={{ minWidth: '9rem' }}
                            body={(r: NFEGestaoRow) => formatDateBR(r.data_emissao)}
                        />
                        <Column
                            field="cnpj_destinatario"
                            header="CNPJ destinatário"
                            sortable
                            style={{ minWidth: '11rem' }}
                            body={(r: NFEGestaoRow) => formatCNPJ(r.cnpj_destinatario)}
                        />
                        <Column
                            field="valor_total"
                            header="Valor"
                            sortable
                            style={{ minWidth: '8rem' }}
                            body={(r: NFEGestaoRow) => formatValor(r.valor_total)}
                        />
                        <Column
                            field="data_download"
                            header="Data download"
                            sortable
                            style={{ minWidth: '10rem' }}
                            body={(r: NFEGestaoRow) => formatDateTimeBR(r.data_download)}
                        />
                        <Column field="chave_nfe" header="Chave" sortable style={{ minWidth: '12rem' }} />
                        <Column header="Detalhes" body={detalhesBody} style={{ minWidth: '6rem' }} />
                    </DataTable>

                    <div className="mt-2">
                        <Button type="button" icon="pi pi-refresh" tooltip="Atualizar" className="p-button-text" onClick={() => void refetch()} />
                    </div>

                    <Dialog
                        header="Detalhes da NF-e"
                        visible={detalhe != null}
                        style={{ width: 'min(40rem, 95vw)' }}
                        onHide={() => setDetalhe(null)}
                        footer={
                            <Button label="Fechar" icon="pi pi-times" text type="button" onClick={() => setDetalhe(null)} />
                        }
                    >
                        {detalhe && (
                            <div className="flex flex-column gap-2 text-sm">
                                <p>
                                    <span className="text-600">Tipo:</span> {detalhe.tipo_arquivo}
                                </p>
                                <p>
                                    <span className="text-600">Número:</span> {detalhe.numero_nfe || '—'}
                                </p>
                                <p>
                                    <span className="text-600">Razão social emitente:</span> {detalhe.razao_social_emitente || '—'}
                                </p>
                                <p>
                                    <span className="text-600">CNPJ emitente:</span> {formatCNPJ(detalhe.cnpj_emitente)}
                                </p>
                                <p>
                                    <span className="text-600">Data emissão:</span> {formatDateBR(detalhe.data_emissao)}
                                </p>
                                <p>
                                    <span className="text-600">CNPJ destinatário:</span> {formatCNPJ(detalhe.cnpj_destinatario)}
                                </p>
                                <p>
                                    <span className="text-600">Valor:</span> {formatValor(detalhe.valor_total)}
                                </p>
                                <p>
                                    <span className="text-600">Data download:</span> {formatDateTimeBR(detalhe.data_download)}
                                </p>
                                <p className="mb-3">
                                    <span className="text-600">Chave:</span> {detalhe.chave_nfe}
                                </p>
                                <Button
                                    label="Abrir na consulta NF-e"
                                    icon="pi pi-external-link"
                                    type="button"
                                    outlined
                                    className="w-full"
                                    onClick={() => {
                                        void router.push(`/nfe/consulta?chave=${encodeURIComponent(detalhe.chave_nfe)}`);
                                        setDetalhe(null);
                                    }}
                                />
                            </div>
                        )}
                    </Dialog>
                </div>
            </div>
        </div>
    );
}
