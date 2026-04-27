import { useQuery } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { DataTable, DataTablePageEvent, DataTableSortEvent } from 'primereact/datatable';
import { Dialog } from 'primereact/dialog';
import { Dropdown } from 'primereact/dropdown';
import { InputText } from 'primereact/inputtext';
import { InputTextarea } from 'primereact/inputtextarea';
import { ProgressSpinner } from 'primereact/progressspinner';
import { Toast } from 'primereact/toast';
import { Toolbar } from 'primereact/toolbar';
import { useRouter } from 'next/router';
import React, { useMemo, useRef, useState } from 'react';

import api from '../../components/api/apiClient';
import { DanfeView } from '../../components/nfe/DanfeView';
import { useRouteClientGuard } from '../../components/hooks/useClientGuards';
import { fetchDanfeJsonByChave, parseDanfeErrorMessage, type NFEDanfeView } from '../../lib/nfeDanfeClient';
import { parseNFEApiError } from '../../lib/nfeError';

type ManifestacaoRow = {
    id: string;
    chave_nfe: string;
    tp_evento: string;
    cnpj_dest: string;
    cstat_lote: number;
    x_motivo_lote?: string;
    cstat_evento: number;
    x_motivo_evento?: string;
    n_prot?: string;
    criado_em: string;
};

type NFEValidacaoRegra = {
    id: string;
    etapa: string;
    codigo_regra: string;
    titulo: string;
    descricao: string;
};

const MANIFEST_TIPO_OPTIONS = [
    { label: '210200 — Confirmação da operação', value: '210200' },
    { label: '210210 — Ciência da operação', value: '210210' },
    { label: '210220 — Desconhecimento (exige justificativa)', value: '210220' },
    { label: '210240 — Operação não realizada (exige justificativa)', value: '210240' },
];

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

const onlyDigitsChave = (v: string) => String(v ?? '').replace(/\D/g, '');

export default function NFEManutencaoPage() {
    useRouteClientGuard();
    const router = useRouter();
    const toast = useRef<Toast>(null);

    const [tipoArquivo, setTipoArquivo] = useState('');
    const [emissaoIni, setEmissaoIni] = useState('');
    const [emissaoFim, setEmissaoFim] = useState('');
    const [chaveNFeFiltroDraft, setChaveNFeFiltroDraft] = useState('');
    const [chaveNFeFiltroApplied, setChaveNFeFiltroApplied] = useState('');
    const [cnpjEmitenteDraft, setCnpjEmitenteDraft] = useState('');
    const [cnpjEmitenteApplied, setCnpjEmitenteApplied] = useState('');
    const [cnpjDestinatarioDraft, setCnpjDestinatarioDraft] = useState('');
    const [cnpjDestinatarioApplied, setCnpjDestinatarioApplied] = useState('');

    const [first, setFirst] = useState(0);
    const [rows, setRows] = useState(20);
    const [sortField, setSortField] = useState('data_download');
    const [sortOrder, setSortOrder] = useState<-1 | 1>(-1);

    const [detalhe, setDetalhe] = useState<NFEGestaoRow | null>(null);
    const [danfeVisible, setDanfeVisible] = useState(false);
    const [danfeLoading, setDanfeLoading] = useState(false);
    const [danfeData, setDanfeData] = useState<NFEDanfeView | null>(null);
    const [danfeChave, setDanfeChave] = useState('');
    const closeDanfePreview = () => {
        setDanfeVisible(false);
        setDanfeData(null);
        setDanfeChave('');
        setDanfeLoading(false);
    };


    const [manifestTp, setManifestTp] = useState('210210');
    const [manifestAmbiente, setManifestAmbiente] = useState<'producao' | 'homologacao'>('producao');
    const [manifestXJust, setManifestXJust] = useState('');
    const [manifestSimular, setManifestSimular] = useState(false);
    const [manifestNSeq, setManifestNSeq] = useState('1');
    const [manifestSending, setManifestSending] = useState(false);

    const queryKey = useMemo(
        () => ({
            first,
            rows,
            sortField,
            sortOrder,
            tipoArquivo,
            chaveNFeFiltroApplied,
            emissaoIni,
            emissaoFim,
            cnpjEmitenteApplied,
            cnpjDestinatarioApplied,
        }),
        [first, rows, sortField, sortOrder, tipoArquivo, chaveNFeFiltroApplied, emissaoIni, emissaoFim, cnpjEmitenteApplied, cnpjDestinatarioApplied],
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
                    chave_nfe: onlyDigitsChave(queryKey.chaveNFeFiltroApplied) || undefined,
                    emissao_ini: queryKey.emissaoIni || undefined,
                    emissao_fim: queryKey.emissaoFim || undefined,
                    cnpj_emitente: queryKey.cnpjEmitenteApplied || undefined,
                    cnpj_destinatario: queryKey.cnpjDestinatarioApplied || undefined,
                },
            });
            return res;
        },
    });

    const chaveDetalhe = detalhe?.chave_nfe ?? '';

    const { data: manifestData, refetch: refetchManifest } = useQuery({
        queryKey: ['nfe-manifestacao', chaveDetalhe],
        enabled: chaveDetalhe.length === 44,
        queryFn: async () => {
            const { data: res } = await api.get<{ items: ManifestacaoRow[]; totalRecords: number }>('/api/serpro/nfe/manifestacao', {
                params: { chave: chaveDetalhe },
            });
            return res;
        },
    });
    const { data: validacoesData } = useQuery({
        queryKey: ['nfe-validacoes-manutencao'],
        queryFn: async () => {
            const { data } = await api.get<{ items: NFEValidacaoRegra[] }>('/api/serpro/nfe/validacoes');
            return data.items ?? [];
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

    const applyTextFilters = () => {
        setChaveNFeFiltroApplied(chaveNFeFiltroDraft);
        setCnpjEmitenteApplied(cnpjEmitenteDraft);
        setCnpjDestinatarioApplied(cnpjDestinatarioDraft);
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
                    chave_nfe: onlyDigitsChave(queryKey.chaveNFeFiltroApplied) || undefined,
                    emissao_ini: queryKey.emissaoIni || undefined,
                    emissao_fim: queryKey.emissaoFim || undefined,
                    cnpj_emitente: queryKey.cnpjEmitenteApplied || undefined,
                    cnpj_destinatario: queryKey.cnpjDestinatarioApplied || undefined,
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

    const visualizarDanfeDaChave = async (chaveRaw: string) => {
        const chave = onlyDigitsChave(chaveRaw);
        if (chave.length !== 44) {
            toast.current?.show({
                severity: 'warn',
                summary: 'Chave inválida',
                detail: 'É necessário uma chave com 44 dígitos.',
                life: 4000,
            });
            return;
        }
        setDetalhe(null);
        setDanfeChave(chave);
        setDanfeVisible(true);
        setDanfeData(null);
        setDanfeLoading(true);
        try {
            const payload = await fetchDanfeJsonByChave(chave);
            setDanfeData(payload);
        } catch (e: unknown) {
            closeDanfePreview();
            const parsed = parseNFEApiError(e);
            const msg = parsed.detail || parseDanfeErrorMessage(e);
            toast.current?.show({ severity: 'error', summary: 'Visualização DANFE', detail: msg, life: 9000 });
        } finally {
            setDanfeLoading(false);
        }
    };

    React.useEffect(() => {
        if (!danfeVisible || !danfeLoading) {
            return;
        }
        const t = window.setTimeout(() => {
            closeDanfePreview();
            toast.current?.show({
                severity: 'warn',
                summary: 'Pré-visualização encerrada',
                detail: 'A visualização DANFE foi encerrada por tempo excedido. Tente novamente.',
                life: 7000,
            });
        }, 25000);
        return () => window.clearTimeout(t);
    }, [danfeVisible, danfeLoading]);

    const enviarManifestacao = async () => {
        if (!detalhe) {
            return;
        }
        const chave = onlyDigitsChave(detalhe.chave_nfe);
        const cnpj = onlyDigitsChave(detalhe.cnpj_destinatario);
        if (chave.length !== 44) {
            toast.current?.show({ severity: 'warn', summary: 'Chave', detail: 'Chave inválida.', life: 4000 });
            return;
        }
        if (cnpj.length !== 14 && cnpj.length !== 11) {
            toast.current?.show({ severity: 'warn', summary: 'Destinatário', detail: 'CNPJ/CPF destinatário inválido na linha.', life: 5000 });
            return;
        }
        if (manifestTp === '210220' || manifestTp === '210240') {
            if (manifestXJust.trim().length < 15) {
                toast.current?.show({
                    severity: 'warn',
                    summary: 'Justificativa',
                    detail: 'Informe x_just com pelo menos 15 caracteres para este tpEvento.',
                    life: 5000,
                });
                return;
            }
        }
        const nSeq = parseInt(manifestNSeq.replace(/\D/g, ''), 10) || 1;
        setManifestSending(true);
        try {
            await api.post('/api/serpro/nfe/manifestar-destinatario', {
                chave_nfe: chave,
                tp_evento: manifestTp,
                cnpj_destinatario: cnpj,
                ambiente: manifestAmbiente,
                x_just: manifestXJust.trim(),
                n_seq_evento: nSeq,
                simular: manifestSimular,
            });
            toast.current?.show({
                severity: 'success',
                summary: 'Manifestação',
                detail: manifestSimular ? 'Modo simulado: registro gravado no histórico.' : 'Envio concluído; verifique cStat no histórico abaixo.',
                life: 6000,
            });
            void refetchManifest();
        } catch (e: unknown) {
            const parsed = parseNFEApiError(e);
            toast.current?.show({ severity: 'error', summary: 'Manifestação', detail: parsed.detail, life: 10000 });
        } finally {
            setManifestSending(false);
        }
    };

    const detalhesBody = (row: NFEGestaoRow) => (
        <Button
            icon="pi pi-search"
            rounded
            severity="info"
            type="button"
            aria-label="Detalhes"
            onClick={() => {
                closeDanfePreview();
                setDetalhe(row);
            }}
        />
    );

    return (
        <div className="grid crud-demo">
            <div className="col-12">
                <div className="card">
                    <Toast ref={toast} />
                    <Dialog
                        header={`DANFE — chave ${danfeChave || '…'}`}
                        visible={danfeVisible}
                        modal={false}
                        style={{ width: 'min(96vw, 960px)' }}
                        contentStyle={{ overflow: 'auto', maxHeight: '92vh' }}
                        maximizable
                        dismissableMask
                        closeOnEscape
                        footer={(
                            <div className="flex justify-content-end">
                                <Button
                                    type="button"
                                    label="Fechar pré-visualização"
                                    icon="pi pi-times"
                                    text
                                    onClick={closeDanfePreview}
                                />
                            </div>
                        )}
                        onHide={closeDanfePreview}
                    >
                        <p className="text-600 text-sm mt-0 mb-3">
                            Documento carregado do tenant e renderizado em visualização DANFE nativa React.
                        </p>
                        {danfeLoading ? (
                            <div className="flex flex-column align-items-center gap-3 py-6">
                                <ProgressSpinner style={{ width: '3rem', height: '3rem' }} />
                                <span className="text-600">Montando visualização DANFE…</span>
                            </div>
                        ) : null}
                        {!danfeLoading && danfeData ? <DanfeView data={danfeData} /> : null}
                    </Dialog>
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
                                <label htmlFor="chaveNfeFiltro" className="text-sm text-600 mb-2 block">
                                    Chave NF-e
                                </label>
                                <InputText
                                    id="chaveNfeFiltro"
                                    className="w-full"
                                    value={chaveNFeFiltroDraft}
                                    maxLength={44}
                                    onChange={(e) => {
                                        setChaveNFeFiltroDraft(onlyDigitsChave(e.target.value));
                                    }}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') {
                                            applyTextFilters();
                                        }
                                    }}
                                    placeholder="44 dígitos ou parcial"
                                />
                            </div>
                            <div className="field mb-0 min-w-0">
                                <label htmlFor="cnpjEmit" className="text-sm text-600 mb-2 block">
                                    CNPJ emitente
                                </label>
                                <InputText
                                    id="cnpjEmit"
                                    className="w-full"
                                    value={cnpjEmitenteDraft}
                                    onChange={(e) => {
                                        setCnpjEmitenteDraft(e.target.value);
                                    }}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') {
                                            applyTextFilters();
                                        }
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
                                    value={cnpjDestinatarioDraft}
                                    onChange={(e) => {
                                        setCnpjDestinatarioDraft(e.target.value);
                                    }}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') {
                                            applyTextFilters();
                                        }
                                    }}
                                    placeholder="Somente números ou parcial"
                                />
                            </div>
                            <div className="field mb-0 min-w-0 flex align-items-end">
                                <Button type="button" label="Aplicar filtros" icon="pi pi-filter" onClick={applyTextFilters} className="w-full" />
                            </div>
                        </div>
                    </div>

                    <Toolbar className="mb-4" left={leftToolbar} right={rightToolbar} />

                    <div className="flex flex-column md:flex-row md:justify-content-between md:align-items-center mb-3">
                        <div>
                            <h5 className="m-0">Manutenção de NFe</h5>
                            <p className="m-0 mt-1 text-600 text-sm">
                                Detalhes e DANFE abrem em janelas sobre esta tela para você continuar consultando outras notas na lista.
                            </p>
                        </div>
                    </div>
                    <div className="surface-50 border-1 surface-border border-round p-3 mb-3">
                        <div className="font-medium mb-2">Regras de validação ativas (catálogo global)</div>
                        {validacoesData && validacoesData.length > 0 ? (
                            <ul className="m-0 pl-3 text-sm">
                                {validacoesData.slice(0, 8).map((r) => (
                                    <li key={r.id}>{r.titulo} - {r.descricao}</li>
                                ))}
                            </ul>
                        ) : (
                            <span className="text-600 text-sm">Nenhuma regra cadastrada.</span>
                        )}
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
                        footer={<Button label="Fechar" icon="pi pi-times" text type="button" onClick={() => setDetalhe(null)} />}
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
                                <div className="border-top-1 surface-border pt-3 mt-2">
                                    <p className="text-600 font-semibold m-0 mb-2">Manifestação do destinatário (SVRS)</p>
                                    <p className="text-600 text-xs m-0 mb-3">
                                        Envio via Recepção de Evento 4.00 (SVRS). Exige certificado A1 do tenant, exceto em modo simulado.
                                    </p>
                                    <div className="flex flex-column gap-2">
                                        <div>
                                            <label className="text-xs text-600 block mb-1">Tipo de evento (tpEvento)</label>
                                            <Dropdown
                                                value={manifestTp}
                                                options={MANIFEST_TIPO_OPTIONS}
                                                className="w-full"
                                                onChange={(e) => setManifestTp(String(e.value ?? '210210'))}
                                            />
                                        </div>
                                        <div>
                                            <label className="text-xs text-600 block mb-1">Ambiente SEFAZ</label>
                                            <Dropdown
                                                value={manifestAmbiente}
                                                options={[
                                                    { label: 'Produção (tpAmb 1)', value: 'producao' },
                                                    { label: 'Homologação (tpAmb 2)', value: 'homologacao' },
                                                ]}
                                                className="w-full"
                                                onChange={(e) => setManifestAmbiente((e.value as 'producao' | 'homologacao') ?? 'producao')}
                                            />
                                        </div>
                                        <div>
                                            <label className="text-xs text-600 block mb-1">nSeqEvento</label>
                                            <InputText
                                                className="w-full"
                                                value={manifestNSeq}
                                                onChange={(e) => setManifestNSeq(String(e.target.value ?? '1'))}
                                                placeholder="1"
                                            />
                                        </div>
                                        {(manifestTp === '210220' || manifestTp === '210240') && (
                                            <div>
                                                <label className="text-xs text-600 block mb-1">Justificativa (mín. 15 caracteres)</label>
                                                <InputTextarea
                                                    className="w-full"
                                                    rows={3}
                                                    value={manifestXJust}
                                                    onChange={(e) => setManifestXJust(e.target.value)}
                                                    autoResize
                                                />
                                            </div>
                                        )}
                                        <div className="field-checkbox mb-0">
                                            <input
                                                id="manifestSim"
                                                type="checkbox"
                                                checked={manifestSimular}
                                                onChange={(e) => setManifestSimular(e.currentTarget.checked)}
                                            />
                                            <label htmlFor="manifestSim" className="ml-2">
                                                Simular (não chama SVRS; grava histórico de teste)
                                            </label>
                                        </div>
                                        <Button
                                            label="Enviar manifestação"
                                            icon="pi pi-send"
                                            type="button"
                                            loading={manifestSending}
                                            className="w-full"
                                            onClick={() => void enviarManifestacao()}
                                        />
                                    </div>
                                    {manifestData?.items && manifestData.items.length > 0 ? (
                                        <div className="mt-3 text-xs">
                                            <span className="text-600 font-semibold">Histórico recente</span>
                                            <ul className="list-none p-0 m-0 mt-2">
                                                {manifestData.items.slice(0, 8).map((m) => (
                                                    <li key={m.id} className="mb-2 border-bottom-1 surface-border pb-2">
                                                        <span className="block">
                                                            {formatDateTimeBR(m.criado_em)} — tp {m.tp_evento} — evento cStat {m.cstat_evento}
                                                            {m.n_prot ? ` — prot. ${m.n_prot}` : ''}
                                                        </span>
                                                        {m.x_motivo_evento ? (
                                                            <span className="text-600 block white-space-normal">{m.x_motivo_evento}</span>
                                                        ) : null}
                                                    </li>
                                                ))}
                                            </ul>
                                        </div>
                                    ) : null}
                                </div>
                                <Button
                                    label="Visualizar DANFE"
                                    icon="pi pi-eye"
                                    type="button"
                                    outlined
                                    className="w-full"
                                    onClick={() => void visualizarDanfeDaChave(detalhe.chave_nfe)}
                                />
                            </div>
                        )}
                    </Dialog>
                </div>
            </div>
        </div>
    );
}
