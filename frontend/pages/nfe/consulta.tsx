import { useRouter } from 'next/router';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Card } from 'primereact/card';
import { InputText } from 'primereact/inputtext';
import { Button } from 'primereact/button';
import { Toast } from 'primereact/toast';
import { InputTextarea } from 'primereact/inputtextarea';
import { Checkbox } from 'primereact/checkbox';
import { Dropdown } from 'primereact/dropdown';
import { Dialog } from 'primereact/dialog';
import { ProgressSpinner } from 'primereact/progressspinner';

import api from '../../components/api/apiClient';
import { useRouteClientGuard } from '../../components/hooks/useClientGuards';
import { DanfeView } from '../../components/nfe/DanfeView';
import { fetchDanfeJsonByChave, parseDanfeErrorMessage, type NFEDanfeView } from '../../lib/nfeDanfeClient';
import { parseNFEApiError } from '../../lib/nfeError';

type AmbienteNFe = 'trial' | 'producao';

const onlyDigitsChave = (v: string) => String(v ?? '').replace(/\D/g, '');

function queryTruthyFlag(raw: string | string[] | undefined): boolean {
    if (raw === undefined) {
        return false;
    }
    const s = Array.isArray(raw) ? raw[0] : raw;
    if (typeof s !== 'string') {
        return false;
    }
    const t = s.trim().toLowerCase();
    return t === '1' || t === 'true' || t === 'sim' || t === 'yes';
}

type NFEDocResponse = {
    id: string;
    chave_nfe: string;
    payload_json: unknown;
    payload_xml?: string;
    evento_codigo?: string;
    evento_descricao?: string;
    recebido_em?: string;
    /** true quando os dados vieram do banco do tenant (já baixados antes), sem nova chamada à SERPRO */
    ja_baixada?: boolean;
};

type NFEValidacaoRegra = {
    id: string;
    etapa: string;
    codigo_regra: string;
    titulo: string;
    descricao: string;
};

type LocalAgentCert = {
    id: string;
    label: string;
    subject: string;
    serial_hex: string;
    slot_id: number;
    token_label: string;
};

export default function NFEConsultaPage() {
    useRouteClientGuard();
    const router = useRouter();

    const toast = useRef<Toast>(null);
    const [ambiente, setAmbiente] = useState<AmbienteNFe>('trial');
    const [chaveNFe, setChaveNFe] = useState('');
    const [requestTag, setRequestTag] = useState('');
    const [assinar, setAssinar] = useState(false);
    const [loading, setLoading] = useState(false);
    const [retorno, setRetorno] = useState('');
    const [danfeVisible, setDanfeVisible] = useState(false);
    const [danfeLoading, setDanfeLoading] = useState(false);
    const [danfeData, setDanfeData] = useState<NFEDanfeView | null>(null);
    const [localPin, setLocalPin] = useState('');
    const [selectedLocalCert, setSelectedLocalCert] = useState('');
    const [localSigning, setLocalSigning] = useState(false);
    const closeDanfePreview = useCallback(() => {
        setDanfeVisible(false);
        setDanfeData(null);
        setDanfeLoading(false);
    }, []);


    /** Evita reexecutar o fluxo automático ao trocar só a query (ex.: remover `visualizar`). */
    const autoDanfeRanForKey = useRef<string | null>(null);

    const { data: validacoes } = useQuery({
        queryKey: ['nfe-validacoes-consulta'],
        queryFn: async () => {
            const { data } = await api.get<{ items: NFEValidacaoRegra[] }>('/api/serpro/nfe/validacoes');
            return data.items ?? [];
        },
    });
    const { data: localCertsData, refetch: refetchLocalCerts, isFetching: localCertsLoading } = useQuery({
        queryKey: ['nfe-local-agent-certs'],
        enabled: false,
        queryFn: async () => {
            const { data } = await api.get<{ items: LocalAgentCert[] }>('/api/local-agent/certificates');
            return data.items ?? [];
        },
    });
    const localCerts = localCertsData ?? [];

    useEffect(() => {
        const raw = router.query.chave;
        const ch = Array.isArray(raw) ? raw[0] : raw;
        if (typeof ch !== 'string') {
            return;
        }
        const digits = onlyDigitsChave(ch);
        if (digits.length === 44) {
            setChaveNFe(digits);
        }
    }, [router.query.chave]);

    const abrirDanfePorChave = useCallback(
        async (chave: string) => {
            setDanfeVisible(true);
            setDanfeData(null);
            setDanfeLoading(true);
            try {
                const payload = await fetchDanfeJsonByChave(chave);
                setDanfeData(payload);
                return true;
            } catch (e: unknown) {
                toast.current?.show({
                    severity: 'error',
                    summary: 'Erro',
                    detail: parseDanfeErrorMessage(e),
                    life: 9000,
                });
                closeDanfePreview();
                return false;
            } finally {
                setDanfeLoading(false);
            }
        },
        [closeDanfePreview],
    );

    useEffect(() => {
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
    }, [danfeVisible, danfeLoading, closeDanfePreview]);

    useEffect(() => {
        if (!router.isReady) {
            return;
        }
        if (!queryTruthyFlag(router.query.visualizar)) {
            autoDanfeRanForKey.current = null;
            return;
        }
        const rawCh = router.query.chave;
        const chStr = Array.isArray(rawCh) ? rawCh[0] : rawCh;
        if (typeof chStr !== 'string') {
            return;
        }
        const chave = onlyDigitsChave(chStr);
        if (chave.length !== 44) {
            return;
        }
        const runKey = `auto-danfe|${chave}`;
        if (autoDanfeRanForKey.current === runKey) {
            return;
        }
        autoDanfeRanForKey.current = runKey;

        let cancelled = false;
        (async () => {
            setChaveNFe(chave);
            setLoading(true);
            try {
                const { data } = await api.get<NFEDocResponse>('/api/serpro/nfe/documento', { params: { chave } });
                if (cancelled) {
                    return;
                }
                const json = JSON.stringify(data, null, 2);
                setRetorno(json);
                toast.current?.show({
                    severity: 'info',
                    summary: 'NF-e no tenant',
                    detail: 'Gerando visualização DANFE…',
                    life: 3000,
                });
                const ok = await abrirDanfePorChave(chave);
                if (cancelled) {
                    return;
                }
                if (!ok) {
                    autoDanfeRanForKey.current = null;
                    return;
                }
                void router.replace({ pathname: '/nfe/consulta', query: { chave } }, undefined, { shallow: true });
            } catch (e: unknown) {
                if (cancelled) {
                    return;
                }
                autoDanfeRanForKey.current = null;
                const parsed = parseNFEApiError(e);
                toast.current?.show({ severity: 'error', summary: parsed.title, detail: parsed.detail, life: 7000 });
            } finally {
                setLoading(false);
            }
        })();

        return () => {
            cancelled = true;
        };
    }, [router.isReady, router.query.chave, router.query.visualizar, router.replace, abrirDanfePorChave]);

    const consultar = async () => {
        const chave = onlyDigitsChave(chaveNFe);
        if (chave.length !== 44) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Informe uma chave de NF-e com 44 dígitos.', life: 3500 });
            return;
        }
        setLoading(true);
        try {
            const { data } = await api.post<NFEDocResponse>('/api/serpro/nfe/consultar', {
                ambiente,
                chave_nfe: chave,
                request_tag: requestTag.trim(),
                assinar,
            });
            setRetorno(JSON.stringify(data, null, 2));
            if (data?.ja_baixada) {
                toast.current?.show({
                    severity: 'warn',
                    summary: 'NF-e já estava baixada',
                    detail: 'Esta nota já estava armazenada no tenant; exibindo os dados salvos, sem nova consulta à SERPRO.',
                    life: 5000,
                });
            } else {
                toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'NF-e consultada e armazenada no schema do tenant.', life: 3500 });
            }
        } catch (e: unknown) {
            const parsed = parseNFEApiError(e);
            toast.current?.show({ severity: 'error', summary: parsed.title, detail: parsed.detail, life: 9000 });
        } finally {
            setLoading(false);
        }
    };

    const buscarPersistida = async () => {
        const chave = onlyDigitsChave(chaveNFe);
        if (chave.length !== 44) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Informe uma chave de NF-e com 44 dígitos.', life: 3500 });
            return;
        }
        setLoading(true);
        try {
            const { data } = await api.get<NFEDocResponse>('/api/serpro/nfe/documento', { params: { chave } });
            setRetorno(JSON.stringify(data, null, 2));
            toast.current?.show({ severity: 'info', summary: 'Consulta local', detail: 'NF-e carregada do banco do tenant.', life: 3000 });
        } catch (e: unknown) {
            const parsed = parseNFEApiError(e);
            toast.current?.show({ severity: 'warn', summary: parsed.title, detail: parsed.detail, life: 7000 });
        } finally {
            setLoading(false);
        }
    };

    const exportarXML = async () => {
        const chave = onlyDigitsChave(chaveNFe);
        if (chave.length !== 44) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Informe uma chave de NF-e com 44 dígitos.', life: 3500 });
            return;
        }
        setLoading(true);
        try {
            const { data } = await api.get<string>('/api/serpro/nfe/documento/xml', {
                params: { chave },
                responseType: 'text' as any,
            });
            const xmlText = String(data ?? '');
            setRetorno(xmlText);
            const blob = new Blob([xmlText], { type: 'application/xml;charset=utf-8' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `nfe_${chave}.xml`;
            a.click();
            URL.revokeObjectURL(url);
            toast.current?.show({ severity: 'success', summary: 'Exportado', detail: 'XML gerado com sucesso.', life: 3000 });
        } catch (e: unknown) {
            const parsed = parseNFEApiError(e);
            toast.current?.show({ severity: 'error', summary: parsed.title, detail: parsed.detail, life: 7000 });
        } finally {
            setLoading(false);
        }
    };

    const assinarHashNoAgente = async () => {
        const chave = onlyDigitsChave(chaveNFe);
        if (chave.length !== 44) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Informe uma chave de NF-e com 44 dígitos.', life: 3500 });
            return;
        }
        setLocalSigning(true);
        try {
            const { data } = await api.post('/api/local-agent/sign-hash', {
                raw_text: chave,
                certificate_id: selectedLocalCert || undefined,
                pin: localPin || undefined,
            });
            setRetorno(JSON.stringify(data, null, 2));
            toast.current?.show({ severity: 'success', summary: 'Agente local', detail: 'Hash assinado com sucesso.', life: 3500 });
        } catch (e: unknown) {
            const parsed = parseNFEApiError(e);
            toast.current?.show({ severity: 'error', summary: 'Agente local', detail: parsed.detail, life: 9000 });
        } finally {
            setLocalSigning(false);
        }
    };

    const visualizarDanfe = async () => {
        const chave = onlyDigitsChave(chaveNFe);
        if (chave.length !== 44) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Informe uma chave de NF-e com 44 dígitos.', life: 3500 });
            return;
        }
        await abrirDanfePorChave(chave);
    };

    return (
        <div className="grid">
            <div className="col-12">
                <Toast ref={toast} />
                <Dialog
                    header="DANFE (visualização nativa React)"
                    visible={danfeVisible}
                    modal={false}
                    style={{ width: '980px' }}
                    contentStyle={{ overflow: 'auto', height: '78vh', minHeight: '78vh', maxHeight: '78vh' }}
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
                        A DANFE é renderizada em componente React com dados normalizados do XML armazenado no tenant.
                    </p>
                    {danfeLoading ? (
                        <div className="flex flex-column align-items-center gap-3 py-6">
                            <ProgressSpinner style={{ width: '3rem', height: '3rem' }} />
                            <span className="text-600">Montando visualização DANFE…</span>
                        </div>
                    ) : null}
                    {!danfeLoading && danfeData ? <DanfeView data={danfeData} /> : null}
                </Dialog>
                <Card title="Consulta NF-e (Tenant)">
                    <p className="text-600 mt-0 mb-4">
                        Consulta a NF-e na SERPRO e persiste o JSON/XML no schema do tenant atual.
                        O ambiente segue o mesmo padrão do Integra Contador (trial ou produção), salvo se o backend
                        tiver <code className="text-sm">SERPRO_NFE_API_BASE_URL</code> definida — nesse caso essa URL fixa prevalece.
                    </p>
                    <div className="grid">
                        <div className="col-12 md:col-2">
                            <label htmlFor="ambiente-nfe" className="block mb-2 font-medium">Ambiente</label>
                            <Dropdown
                                id="ambiente-nfe"
                                value={ambiente}
                                options={[
                                    { label: 'Trial', value: 'trial' },
                                    { label: 'Produção', value: 'producao' },
                                ]}
                                optionLabel="label"
                                optionValue="value"
                                className="w-full"
                                onChange={(e) => setAmbiente(e.value as AmbienteNFe)}
                            />
                        </div>
                        <div className="col-12 md:col-4">
                            <label htmlFor="chave-nfe" className="block mb-2 font-medium">Chave NF-e (44 dígitos)</label>
                            <InputText
                                id="chave-nfe"
                                className="w-full"
                                value={chaveNFe}
                                maxLength={44}
                                onChange={(e) => setChaveNFe(onlyDigitsChave(e.target.value))}
                                placeholder="Somente números"
                            />
                        </div>
                        <div className="col-12 md:col-3">
                            <label htmlFor="request-tag" className="block mb-2 font-medium">x-request-tag (opcional)</label>
                            <InputText
                                id="request-tag"
                                className="w-full"
                                value={requestTag}
                                maxLength={32}
                                onChange={(e) => setRequestTag(e.target.value)}
                            />
                        </div>
                        <div className="col-12 md:col-3 flex align-items-end">
                            <div className="field-checkbox m-0">
                                <Checkbox
                                    inputId="assinar"
                                    checked={assinar}
                                    onChange={(e) => setAssinar(Boolean(e.checked))}
                                />
                                <label htmlFor="assinar" className="ml-2">x-signature=1</label>
                            </div>
                        </div>
                        <div className="col-12 flex gap-2 flex-wrap">
                            <Button type="button" label="Consultar SERPRO" icon="pi pi-search" onClick={consultar} loading={loading} />
                            <Button type="button" label="Buscar no Tenant" icon="pi pi-database" severity="secondary" onClick={buscarPersistida} loading={loading} />
                            <Button type="button" label="Exportar XML" icon="pi pi-download" severity="help" onClick={exportarXML} loading={loading} />
                            <Button
                                type="button"
                                label="Visualizar DANFE"
                                icon="pi pi-eye"
                                severity="secondary"
                                disabled={!retorno.trim() || loading || danfeLoading}
                                loading={danfeLoading}
                                onClick={() => void visualizarDanfe()}
                            />
                        </div>
                        <div className="col-12">
                            <div className="surface-50 border-1 surface-border border-round p-3">
                                <div className="font-medium mb-2">Regras de validação ativas (catálogo global)</div>
                                {validacoes && validacoes.length > 0 ? (
                                    <ul className="m-0 pl-3 text-sm">
                                        {validacoes.slice(0, 6).map((r) => (
                                            <li key={r.id}>{r.titulo} - {r.descricao}</li>
                                        ))}
                                    </ul>
                                ) : (
                                    <span className="text-600 text-sm">Nenhuma regra cadastrada.</span>
                                )}
                            </div>
                        </div>
                        <div className="col-12">
                            <div className="surface-50 border-1 surface-border border-round p-3">
                                <div className="font-medium mb-2">Agente local A3 (EF-930)</div>
                                <p className="text-600 text-sm mt-0 mb-3">
                                    Usa o agente local para listar certificados do token/smartcard e assinar o hash SHA-256 da chave da NF-e.
                                </p>
                                <div className="flex gap-2 flex-wrap mb-3">
                                    <Button
                                        type="button"
                                        label="Listar certificados locais"
                                        icon="pi pi-refresh"
                                        severity="secondary"
                                        loading={localCertsLoading}
                                        onClick={() => void refetchLocalCerts()}
                                    />
                                    <Dropdown
                                        value={selectedLocalCert}
                                        options={localCerts.map((c) => ({ label: `${c.label} (${c.token_label})`, value: c.id }))}
                                        onChange={(e) => setSelectedLocalCert(String(e.value ?? ''))}
                                        placeholder="Selecione um certificado (opcional)"
                                        className="min-w-20rem"
                                    />
                                    <InputText
                                        type="password"
                                        value={localPin}
                                        onChange={(e) => setLocalPin(e.target.value)}
                                        placeholder="PIN do token (opcional)"
                                        className="min-w-14rem"
                                    />
                                    <Button
                                        type="button"
                                        label="Assinar hash no agente"
                                        icon="pi pi-lock"
                                        loading={localSigning}
                                        onClick={() => void assinarHashNoAgente()}
                                    />
                                </div>
                            </div>
                        </div>
                        <div className="col-12">
                            <label htmlFor="retorno" className="block mb-2 font-medium">Retorno</label>
                            <InputTextarea
                                id="retorno"
                                className="w-full"
                                rows={18}
                                value={retorno}
                                onChange={(e) => setRetorno(e.target.value)}
                            />
                        </div>
                    </div>
                </Card>
            </div>
        </div>
    );
}
