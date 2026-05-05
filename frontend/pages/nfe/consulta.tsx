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
import { useAuthScopeKey } from '../../components/hooks/useAuthScopeKey';
import { useRouteClientGuard } from '../../components/hooks/useClientGuards';
import { DanfeTesteView } from '../../components/nfe/DanfeTesteView';
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
    tax_ids?: string[];
};

export default function NFEConsultaPage() {
    useRouteClientGuard();
    const authScope = useAuthScopeKey();
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
    const [danfePane, setDanfePane] = useState<'consulta' | 'matriz'>('consulta');
    const [localPin, setLocalPin] = useState('');
    const [selectedLocalCert, setSelectedLocalCert] = useState('');
    const [localSigning, setLocalSigning] = useState(false);
    /** EF-937: CNPJ/CPF do cliente (somente dígitos) para resolver cert_clientes/{id}.pfx ou titular A3 */
    const [taxIdLocalSign, setTaxIdLocalSign] = useState('');
    const [procuracaoLocal, setProcuracaoLocal] = useState(false);
    /** CNPJ/CPF do contador quando procuracao + A3 com vários certificados no token */
    const [signerTaxIdLocal, setSignerTaxIdLocal] = useState('');
    /** Aviso quando o usuário tenta assinar sem PIN/senha (A3 quase sempre exige PIN no token). */
    const [pinAgenteDialogOpen, setPinAgenteDialogOpen] = useState(false);
    const closeDanfePreview = useCallback(() => {
        setDanfeVisible(false);
        setDanfeData(null);
        setDanfeLoading(false);
        setDanfePane('consulta');
    }, []);


    /** Evita reexecutar o fluxo automático ao trocar só a query (ex.: remover `visualizar`). */
    const autoDanfeRanForKey = useRef<string | null>(null);

    const { data: validacoes } = useQuery({
        queryKey: ['nfe-validacoes-consulta', authScope],
        queryFn: async () => {
            const { data } = await api.get<{ items: NFEValidacaoRegra[] }>('/api/serpro/nfe/validacoes');
            return data.items ?? [];
        },
    });
    const { data: localCertsData, refetch: refetchLocalCerts, isFetching: localCertsLoading } = useQuery({
        queryKey: ['nfe-local-agent-certs', authScope],
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
            setDanfePane('consulta');
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
        if (danfePane !== 'matriz' || !danfeData) {
            return;
        }
        const t = window.setTimeout(() => window.print(), 600);
        return () => window.clearTimeout(t);
    }, [danfePane, danfeData]);

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

    const runAssinarHashNoAgente = async (confirmouSemPin: boolean) => {
        const chave = onlyDigitsChave(chaveNFe);
        if (chave.length !== 44) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Informe uma chave de NF-e com 44 dígitos.', life: 3500 });
            return;
        }
        const taxDigits = onlyDigitsChave(taxIdLocalSign);
        if (procuracaoLocal && (taxDigits.length !== 11 && taxDigits.length !== 14)) {
            toast.current?.show({
                severity: 'warn',
                summary: 'EF-937',
                detail: 'Com procuração, informe o CNPJ ou CPF do cliente (11 ou 14 dígitos).',
                life: 5000,
            });
            return;
        }
        if (!confirmouSemPin && !localPin.trim()) {
            setPinAgenteDialogOpen(true);
            return;
        }
        setLocalSigning(true);
        try {
            const signerDigits = onlyDigitsChave(signerTaxIdLocal);
            const { data } = await api.post('/api/local-agent/sign-hash', {
                raw_text: chave,
                document_id: chave,
                tax_id: taxDigits || undefined,
                procuracao: procuracaoLocal,
                signer_tax_id: signerDigits || undefined,
                certificate_id: selectedLocalCert || undefined,
                pin: localPin.trim() || undefined,
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

    const assinarHashNoAgente = () => void runAssinarHashNoAgente(false);

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
                    header="PIN do token ou senha do .pfx"
                    visible={pinAgenteDialogOpen}
                    modal
                    dismissableMask
                    style={{ width: 'min(32rem, 94vw)' }}
                    onHide={() => setPinAgenteDialogOpen(false)}
                    footer={(
                        <div className="flex justify-content-end gap-2 flex-wrap">
                            <Button type="button" label="Voltar e preencher" severity="secondary" onClick={() => setPinAgenteDialogOpen(false)} />
                            <Button
                                type="button"
                                label="Continuar sem PIN/senha"
                                onClick={() => {
                                    setPinAgenteDialogOpen(false);
                                    void runAssinarHashNoAgente(true);
                                }}
                            />
                        </div>
                    )}
                >
                    <p className="text-600 mt-0 mb-2">
                        O campo ao lado de <strong>Assinar hash no agente</strong> envia o valor ao vecx-agent como <code className="text-sm">pin</code>:
                    </p>
                    <ul className="m-0 pl-3 text-sm text-600">
                        <li><strong>A3 (token/cartão):</strong> na prática é obrigatório informar o PIN do dispositivo para a biblioteca PKCS#11 autenticar a chave privada (salvo sessão já desbloqueada no driver).</li>
                        <li><strong>A1 (.pfx em disco, EF-937):</strong> use a <strong>senha do arquivo PKCS#12</strong> nesse mesmo campo.</li>
                    </ul>
                    <p className="text-600 text-sm mb-0">
                        O agente desktop <strong>não</strong> abre um diálogo próprio para PIN: a coleta fica nesta tela (ou em qualquer cliente que chame a API). Só continue sem preencher se souber que o token já está desbloqueado.
                    </p>
                </Dialog>
                <Dialog
                    header={
                        danfePane === 'matriz'
                            ? 'DANFE — impressão (A4)'
                            : 'DANFE (visualização nativa React)'
                    }
                    visible={danfeVisible}
                    modal={false}
                    style={{ width: danfePane === 'matriz' ? 'min(86rem, 99.5vw)' : '980px' }}
                    contentStyle={{ overflow: 'auto', height: '78vh', minHeight: '78vh', maxHeight: '78vh' }}
                    maximizable
                    dismissableMask
                    closeOnEscape
                    footer={(
                        <div className="flex justify-content-end flex-wrap gap-2">
                            {danfePane === 'matriz' ? (
                                <Button
                                    type="button"
                                    label="Voltar à consulta (abas)"
                                    icon="pi pi-arrow-left"
                                    text
                                    onClick={() => setDanfePane('consulta')}
                                />
                            ) : null}
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
                    {danfePane === 'consulta' ? (
                        <p className="text-600 text-sm mt-0 mb-3">
                            A DANFE é renderizada em componente React com dados normalizados do XML armazenado no tenant. Em{' '}
                            <strong>Imprimir</strong>, abre-se a matriz A4 para pré-visualização de impressão.
                        </p>
                    ) : (
                        <p className="text-600 text-sm mt-0 mb-3">
                            Matriz A4 retrato (MOC 7.0 — Anexo II, 3.8.1) com os dados desta NF-e.
                        </p>
                    )}
                    {danfeLoading ? (
                        <div className="flex flex-column align-items-center gap-3 py-6">
                            <ProgressSpinner style={{ width: '3rem', height: '3rem' }} />
                            <span className="text-600">Montando visualização DANFE…</span>
                        </div>
                    ) : null}
                    {!danfeLoading && danfeData && danfePane === 'consulta' ? (
                        <DanfeView data={danfeData} onImprimirMatriz={() => setDanfePane('matriz')} />
                    ) : null}
                    {!danfeLoading && danfeData && danfePane === 'matriz' ? (
                        <DanfeTesteView data={danfeData} variant="default" />
                    ) : null}
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
                                <div className="font-medium mb-2">Agente local (EF-930 / EF-937)</div>
                                <p className="text-600 text-sm mt-0 mb-3">
                                    Lista certificados A3 no token ou resolve A1 em disco conforme pasta raiz configurada no agente
                                    (<code className="text-sm">cert_clientes</code> / <code className="text-sm">cert_contador</code>).
                                    Com <strong>CNPJ/CPF do cliente</strong> e opcional <strong>procuração</strong>, o backend envia os metadados para o agente escolher o .pfx ou o certificado A3 pelo titular.
                                </p>
                                <div className="grid mb-3">
                                    <div className="col-12 md:col-4">
                                        <label htmlFor="tax-local-sign" className="block mb-2 font-medium">CNPJ/CPF cliente (EF-937)</label>
                                        <InputText
                                            id="tax-local-sign"
                                            className="w-full"
                                            value={taxIdLocalSign}
                                            onChange={(e) => setTaxIdLocalSign(onlyDigitsChave(e.target.value))}
                                            placeholder="Opcional — obrigatório se procuração"
                                            maxLength={14}
                                        />
                                    </div>
                                    <div className="col-12 md:col-4">
                                        <label htmlFor="signer-tax-local" className="block mb-2 font-medium">CNPJ/CPF contador (A3 + procuração)</label>
                                        <InputText
                                            id="signer-tax-local"
                                            className="w-full"
                                            value={signerTaxIdLocal}
                                            onChange={(e) => setSignerTaxIdLocal(onlyDigitsChave(e.target.value))}
                                            placeholder="Opcional"
                                            maxLength={14}
                                        />
                                    </div>
                                    <div className="col-12 md:col-4 flex align-items-end">
                                        <div className="field-checkbox m-0">
                                            <input
                                                type="checkbox"
                                                id="procuracao-local"
                                                checked={procuracaoLocal}
                                                onChange={(e) => setProcuracaoLocal(e.target.checked)}
                                            />
                                            <label htmlFor="procuracao-local" className="ml-2">Procuração (certificado do contador)</label>
                                        </div>
                                    </div>
                                </div>
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
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="local-agent-pin" className="block mb-1 font-medium">PIN (A3) ou senha do .pfx (A1)</label>
                                    <InputText
                                        id="local-agent-pin"
                                        type="password"
                                        value={localPin}
                                        onChange={(e) => setLocalPin(e.target.value)}
                                        placeholder="Obrigatório na maioria dos tokens A3"
                                        className="min-w-14rem"
                                        autoComplete="off"
                                    />
                                    <small className="text-600 block mt-1">
                                        Mesmo valor é enviado ao agente como <code className="text-sm">pin</code> (login PKCS#11 ou senha do PKCS#12).
                                    </small>
                                </div>
                                <div className="flex gap-2 flex-wrap mb-3">
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
