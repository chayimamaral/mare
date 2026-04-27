import { useQuery } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { DataTable, DataTablePageEvent } from 'primereact/datatable';
import { Dropdown } from 'primereact/dropdown';
import { InputText } from 'primereact/inputtext';
import { Toast } from 'primereact/toast';
import React, { useMemo, useRef, useState } from 'react';

import api from '../../components/api/apiClient';
import { useRouteClientGuard } from '../../components/hooks/useClientGuards';

type SyncEstadoRow = {
    id: string;
    provider: string;
    uf: string;
    cnpj: string;
    ultimo_nsu: string;
    ultimo_cstat?: number;
    ultimo_motivo?: string;
    ultima_verificacao?: string | null;
    proxima_consulta_apos?: string | null;
    /** qtDfeRet da última resposta SC (regra de intervalo 118). */
    ultima_qt_dfe_ret?: number;
};

type SyncResponse = {
    provider: string;
    uf: string;
    cnpj: string;
    anterior_nsu: string;
    novo_nsu: string;
    total_recebidos: number;
    total_persistidos: number;
    cstat: number;
    x_motivo?: string;
    ultima_qt_dfe_ret?: number;
};

const PROVIDER_OPTIONS = [
    { label: 'SC (SEF-SC)', value: 'SC' },
    { label: 'Nacional (stub)', value: 'NACIONAL' },
];

function formatDateTime(iso?: string | null): string {
    if (!iso) {
        return '—';
    }
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) {
        return String(iso);
    }
    return d.toLocaleString('pt-BR');
}

function onlyDigits(v: string): string {
    return String(v ?? '').replace(/\D/g, '');
}

function formatCNPJCPF(v: string): string {
    const d = onlyDigits(v);
    if (d.length === 14) {
        return d.replace(/^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$/, '$1.$2.$3/$4-$5');
    }
    if (d.length === 11) {
        return d.replace(/^(\d{3})(\d{3})(\d{3})(\d{2})$/, '$1.$2.$3-$4');
    }
    return v || '—';
}

export default function NFESincronizacaoPage() {
    useRouteClientGuard();
    const toast = useRef<Toast>(null);

    const [provider, setProvider] = useState('SC');
    const [uf, setUF] = useState('SC');
    const [cnpj, setCNPJ] = useState('');
    const [ambiente, setAmbiente] = useState('producao');
    const [simular, setSimular] = useState(true);
    const [loadingSync, setLoadingSync] = useState(false);

    const [fProvider, setFProvider] = useState('');
    const [fUF, setFUF] = useState('');
    const [fCNPJ, setFCNPJ] = useState('');
    const [appliedFProvider, setAppliedFProvider] = useState('');
    const [appliedFUF, setAppliedFUF] = useState('');
    const [appliedFCNPJ, setAppliedFCNPJ] = useState('');
    const [first, setFirst] = useState(0);
    const [rows, setRows] = useState(20);

    const queryKey = useMemo(
        () => ({ first, rows, provider: appliedFProvider, uf: appliedFUF, cnpj: appliedFCNPJ }),
        [first, rows, appliedFProvider, appliedFUF, appliedFCNPJ],
    );

    const { data, isFetching, refetch } = useQuery({
        queryKey: ['nfe-sync-estado', queryKey],
        queryFn: async () => {
            const { data: res } = await api.get<{ items: SyncEstadoRow[]; totalRecords: number }>('/api/serpro/nfe/sync-estado', {
                params: {
                    first: queryKey.first,
                    rows: queryKey.rows,
                    provider: queryKey.provider || undefined,
                    uf: queryKey.uf || undefined,
                    cnpj: queryKey.cnpj || undefined,
                },
            });
            return res;
        },
    });

    const onPage = (e: DataTablePageEvent) => {
        setFirst(e.first);
        setRows(e.rows);
    };

    const applyFiltros = () => {
        setAppliedFProvider(fProvider);
        setAppliedFUF(fUF);
        setAppliedFCNPJ(fCNPJ);
        setFirst(0);
    };

    const executarSincronizacao = async () => {
        const cnpjDigits = onlyDigits(cnpj);
        if (cnpjDigits.length !== 14 && cnpjDigits.length !== 11) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Informe um CNPJ/CPF válido.', life: 4000 });
            return;
        }
        setLoadingSync(true);
        try {
            const { data: res } = await api.post<SyncResponse>('/api/serpro/nfe/sincronizar-provider', {
                provider,
                uf: (uf || 'SC').toUpperCase(),
                cnpj: cnpjDigits,
                ambiente,
                simular,
            });
            const qtPart = ` | qtDfeRet: ${res.ultima_qt_dfe_ret ?? '—'}`;
            toast.current?.show({
                severity: 'success',
                summary: 'Sincronização concluída',
                detail: `cStat ${res.cstat}${res.x_motivo ? ` — ${res.x_motivo}` : ''} | Recebidos: ${res.total_recebidos} | Persistidos: ${res.total_persistidos} | NSU: ${res.anterior_nsu} -> ${res.novo_nsu}${qtPart}`,
                life: 8000,
            });
            setFProvider(res.provider);
            setFUF(res.uf);
            setFCNPJ(res.cnpj);
            setAppliedFProvider(res.provider);
            setAppliedFUF(res.uf);
            setAppliedFCNPJ(res.cnpj);
            setFirst(0);
            void refetch();
        } catch (e: any) {
            const msg = e?.response?.data?.error || e?.response?.data?.message || 'Falha ao sincronizar';
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 7000 });
        } finally {
            setLoadingSync(false);
        }
    };

    return (
        <div className="grid crud-demo">
            <div className="col-12">
                <div className="card">
                    <Toast ref={toast} />
                    <div className="mb-4">
                        <h5 className="m-0">Sincronização NFe por Provider</h5>
                        <p className="m-0 mt-2 text-600 text-sm">
                            Execute a sincronização automática por provider e acompanhe o estado de NSU por CNPJ/UF.
                        </p>
                    </div>

                    <div
                        className="mb-4"
                        style={{
                            display: 'grid',
                            gridTemplateColumns: 'repeat(auto-fit, minmax(14rem, 1fr))',
                            gap: '1rem',
                        }}
                    >
                        <div className="field mb-0">
                            <label htmlFor="providerSync" className="text-sm text-600 mb-2 block">Provider</label>
                            <Dropdown
                                id="providerSync"
                                value={provider}
                                options={PROVIDER_OPTIONS}
                                onChange={(e) => setProvider(String(e.value ?? 'SC'))}
                                className="w-full"
                            />
                        </div>
                        <div className="field mb-0">
                            <label htmlFor="ufSync" className="text-sm text-600 mb-2 block">UF</label>
                            <InputText
                                id="ufSync"
                                value={uf}
                                maxLength={2}
                                className="w-full"
                                onChange={(e) => setUF(String(e.target.value ?? '').toUpperCase())}
                                placeholder="SC"
                            />
                        </div>
                        <div className="field mb-0">
                            <label htmlFor="cnpjSync" className="text-sm text-600 mb-2 block">CNPJ/CPF</label>
                            <InputText
                                id="cnpjSync"
                                value={cnpj}
                                className="w-full"
                                onChange={(e) => setCNPJ(onlyDigits(e.target.value))}
                                placeholder="Somente números"
                            />
                        </div>
                        <div className="field mb-0">
                            <label htmlFor="ambienteSync" className="text-sm text-600 mb-2 block">Ambiente</label>
                            <Dropdown
                                id="ambienteSync"
                                value={ambiente}
                                options={[
                                    { label: 'Produção', value: 'producao' },
                                    { label: 'Homologação', value: 'homologacao' },
                                ]}
                                onChange={(e) => setAmbiente(String(e.value ?? 'producao'))}
                                className="w-full"
                            />
                        </div>
                    </div>

                    <div className="field-checkbox mb-4">
                        <input
                            id="simularSync"
                            type="checkbox"
                            checked={simular}
                            onChange={(e) => setSimular(e.currentTarget.checked)}
                        />
                        <label htmlFor="simularSync" className="ml-2">
                            Modo simulado (não exige certificado/chaves reais)
                        </label>
                    </div>

                    <div className="mb-4">
                        <Button
                            type="button"
                            label="Executar sincronização"
                            icon="pi pi-play"
                            loading={loadingSync}
                            onClick={() => void executarSincronizacao()}
                        />
                    </div>

                    <div className="mb-3 text-sm text-600">
                        <strong>Estado da sincronização</strong> (checkpoint por provider/UF/CNPJ).
                    </div>

                    <div
                        className="mb-3"
                        style={{
                            display: 'grid',
                            gridTemplateColumns: 'repeat(auto-fit, minmax(14rem, 1fr))',
                            gap: '1rem',
                        }}
                    >
                        <div className="field mb-0">
                            <label htmlFor="fProvider" className="text-sm text-600 mb-2 block">Provider</label>
                            <InputText
                                id="fProvider"
                                className="w-full"
                                value={fProvider}
                                onChange={(e) => setFProvider(String(e.target.value ?? '').toUpperCase())}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') {
                                        applyFiltros();
                                    }
                                }}
                            />
                        </div>
                        <div className="field mb-0">
                            <label htmlFor="fUF" className="text-sm text-600 mb-2 block">UF</label>
                            <InputText
                                id="fUF"
                                className="w-full"
                                value={fUF}
                                maxLength={2}
                                onChange={(e) => setFUF(String(e.target.value ?? '').toUpperCase())}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') {
                                        applyFiltros();
                                    }
                                }}
                            />
                        </div>
                        <div className="field mb-0">
                            <label htmlFor="fCnpj" className="text-sm text-600 mb-2 block">CNPJ/CPF</label>
                            <InputText
                                id="fCnpj"
                                className="w-full"
                                value={fCNPJ}
                                onChange={(e) => setFCNPJ(onlyDigits(e.target.value))}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter') {
                                        applyFiltros();
                                    }
                                }}
                            />
                        </div>
                        <div className="field mb-0 flex align-items-end">
                            <Button type="button" label="Aplicar filtros" icon="pi pi-filter" onClick={applyFiltros} />
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
                        emptyMessage="Nenhum checkpoint de sincronização encontrado."
                        stripedRows
                        size="small"
                    >
                        <Column field="provider" header="Provider" style={{ minWidth: '8rem' }} />
                        <Column field="uf" header="UF" style={{ minWidth: '5rem' }} />
                        <Column field="cnpj" header="CNPJ/CPF" body={(r: SyncEstadoRow) => formatCNPJCPF(r.cnpj)} style={{ minWidth: '11rem' }} />
                        <Column field="ultimo_nsu" header="Último NSU" style={{ minWidth: '8rem' }} />
                        <Column field="ultimo_cstat" header="cStat" style={{ minWidth: '6rem' }} />
                        <Column
                            field="ultima_qt_dfe_ret"
                            header="qtDfeRet"
                            body={(r: SyncEstadoRow) => (r.ultima_qt_dfe_ret != null ? String(r.ultima_qt_dfe_ret) : '—')}
                            style={{ minWidth: '6rem' }}
                        />
                        <Column field="ultimo_motivo" header="Motivo" style={{ minWidth: '20rem' }} />
                        <Column field="ultima_verificacao" header="Última verificação" body={(r: SyncEstadoRow) => formatDateTime(r.ultima_verificacao)} style={{ minWidth: '12rem' }} />
                        <Column field="proxima_consulta_apos" header="Próxima consulta" body={(r: SyncEstadoRow) => formatDateTime(r.proxima_consulta_apos)} style={{ minWidth: '12rem' }} />
                    </DataTable>

                    <div className="mt-2">
                        <Button type="button" icon="pi pi-refresh" tooltip="Atualizar" className="p-button-text" onClick={() => void refetch()} />
                    </div>
                </div>
            </div>
        </div>
    );
}
