import dynamic from 'next/dynamic';
import { useRouter } from 'next/router';
import { useEffect, useRef, useState } from 'react';
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

const DanfeHtmlIframe = dynamic(
    () => import('../../components/nfe/DanfeHtmlIframe').then((m) => m.DanfeHtmlIframe),
    { ssr: false },
);

type AmbienteNFe = 'trial' | 'producao';

const onlyDigitsChave = (v: string) => String(v ?? '').replace(/\D/g, '');

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
    const [danfeHtml, setDanfeHtml] = useState<string | null>(null);
    /** Força novo iframe a cada HTML recebido (evita cache/estado estranho do documento embutido). */
    const [danfeFrameKey, setDanfeFrameKey] = useState(0);

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
        } catch (e: any) {
            const msg = e?.response?.data?.error || e?.response?.data?.message || 'Falha ao consultar NF-e';
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 7000 });
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
        } catch (e: any) {
            const msg = e?.response?.data?.error || e?.response?.data?.message || 'NF-e não encontrada no tenant';
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: msg, life: 5000 });
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
        } catch (e: any) {
            const msg = e?.response?.data?.error || e?.response?.data?.message || 'Falha ao exportar XML';
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 5000 });
        } finally {
            setLoading(false);
        }
    };

    const visualizarDanfe = async () => {
        const texto = retorno.trim();
        if (!texto) {
            toast.current?.show({
                severity: 'warn',
                summary: 'Retorno vazio',
                detail: 'Consulte a SERPRO, busque no tenant ou cole JSON/XML no quadro Retorno.',
                life: 5000,
            });
            return;
        }
        setDanfeVisible(true);
        setDanfeHtml(null);
        setDanfeLoading(true);
        try {
            const { data } = await api.post<string>(
                '/api/serpro/nfe/documento/danfe-html',
                { retorno: texto },
                { responseType: 'text' as any },
            );
            const html = String(data ?? '').trim();
            if (!html) {
                toast.current?.show({ severity: 'warn', summary: 'DANFE vazio', detail: 'Resposta sem HTML.', life: 5000 });
                setDanfeVisible(false);
                return;
            }
            setDanfeHtml(html);
            setDanfeFrameKey((k) => k + 1);
        } catch (e: any) {
            const raw = e?.response?.data;
            let msg = raw?.error || raw?.message || e?.message || 'Falha ao gerar DANFE';
            if (typeof raw === 'string' && raw.trim().startsWith('{')) {
                try {
                    const o = JSON.parse(raw);
                    if (o?.error) msg = o.error;
                } catch {
                    /* ignore */
                }
            }
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 9000 });
            setDanfeVisible(false);
        } finally {
            setDanfeLoading(false);
        }
    };

    return (
        <div className="grid">
            <div className="col-12">
                <Toast ref={toast} />
                <Dialog
                    header="DANFE (Saxon + XSLT SVRS)"
                    visible={danfeVisible}
                    style={{ width: 'min(96vw, 960px)' }}
                    contentStyle={{ overflow: 'auto', maxHeight: '92vh' }}
                    maximizable
                    onHide={() => {
                        setDanfeVisible(false);
                        setDanfeHtml(null);
                        setDanfeFrameKey(0);
                    }}
                >
                    <p className="text-600 text-sm mt-0 mb-3">
                        O servidor interpreta o quadro <strong>Retorno</strong>: JSON da consulta (com <code className="text-xs">payload_json</code>{' '}
                        trial SERPRO), JSON bruto com <code className="text-xs">nfeProc</code>, ou XML com namespace SEFAZ. Geração com Saxon-HE +{' '}
                        <a href="https://dfe-portal.svrs.rs.gov.br/Schemas/PRNFE/XSLT/" target="_blank" rel="noreferrer">
                            XSLT SVRS
                        </a>
                        . Uso demonstrativo com notas de teste; documento fiscal oficial segue outro fluxo.
                    </p>
                    {danfeLoading ? (
                        <div className="flex flex-column align-items-center gap-3 py-6">
                            <ProgressSpinner style={{ width: '3rem', height: '3rem' }} />
                            <span className="text-600">Gerando DANFE no servidor…</span>
                        </div>
                    ) : null}
                    {!danfeLoading && danfeHtml ? <DanfeHtmlIframe key={danfeFrameKey} html={danfeHtml} /> : null}
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
