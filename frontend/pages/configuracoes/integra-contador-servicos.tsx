import { useMemo, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Card } from 'primereact/card';
import { Toast } from 'primereact/toast';
import { Dropdown } from 'primereact/dropdown';
import { InputText } from 'primereact/inputtext';
import { InputTextarea } from 'primereact/inputtextarea';
import { Button } from 'primereact/button';
import { Tag } from 'primereact/tag';
import setupAPIClient from '../../components/api/api';
import { canSSRAuth } from '../../components/utils/canSSRAuth';
import CatalogoServicoService, { CatalogoServico } from '../../services/cruds/CatalogoServicoService';

type Ambiente = 'trial' | 'producao';
type NI = { numero: string; tipo: number };

type IntegraCallResponse = {
    status_code: number;
    headers: Record<string, string>;
    raw_body: string;
};

type ServicoProcuracao = {
    id_sistema: string;
    id_servico: string;
    cod_procuracao: string;
    nome_servico: string;
};

function exigeProcuracao(codProcuracao: string): boolean {
    const raw = (codProcuracao || '').trim().toLowerCase();
    if (!raw) return false;
    if (raw === 'n/a' || raw === 'na' || raw === 'não se aplica' || raw === 'aguardando definição') return false;
    if (raw.startsWith('sim')) return true;
    return /^\d+$/.test(raw);
}

function interpretarStatusHTTP(status: number): { severidade: 'success' | 'info' | 'warn' | 'error'; titulo: string; orientacao: string } {
    switch (status) {
        case 200:
            return { severidade: 'success', titulo: '200 OK', orientacao: 'Requisição concluída com sucesso.' };
        case 202:
            return { severidade: 'info', titulo: '202 Accepted', orientacao: 'Processando. Aguarde tempoEspera e consulte novamente.' };
        case 204:
            return { severidade: 'info', titulo: '204 No Content', orientacao: 'Sem conteúdo por enquanto. Verifique eTag/header de espera.' };
        case 304:
            return { severidade: 'info', titulo: '304 Not Modified', orientacao: 'Sem alteração no conteúdo. Pode reutilizar cache/token informado.' };
        case 400:
            return { severidade: 'warn', titulo: '400 Bad Request', orientacao: 'Valide sintaxe e payload de entrada.' };
        case 401:
            return { severidade: 'warn', titulo: '401 Unauthorized', orientacao: 'Reautentique e obtenha novos access_token/jwt_token.' };
        case 403:
            return { severidade: 'warn', titulo: '403 Forbidden', orientacao: 'Verifique permissão e procuração eCAC para autor/contribuinte.' };
        case 404:
            return { severidade: 'error', titulo: '404 Not Found', orientacao: 'Recurso/serviço não encontrado para a rota usada.' };
        case 429:
            return { severidade: 'warn', titulo: '429 Too Many Requests', orientacao: 'Limite de chamadas excedido. Aplique backoff e tente depois.' };
        case 500:
            return { severidade: 'error', titulo: '500 Internal Server Error', orientacao: 'Falha interna do provedor. Repetir após intervalo.' };
        case 503:
            return { severidade: 'error', titulo: '503 Service Unavailable', orientacao: 'Serviço indisponível no momento. Tente novamente mais tarde.' };
        default:
            if (status >= 200 && status < 300) {
                return { severidade: 'success', titulo: `${status}`, orientacao: 'Requisição aceita com sucesso.' };
            }
            return { severidade: 'error', titulo: `${status}`, orientacao: 'Resposta não mapeada. Validar detalhes no retorno bruto.' };
    }
}

function presetPorServico(idServico: string): string {
    const key = (idServico || '').trim().toUpperCase();
    switch (key) {
        case 'GERARDASPDF21':
        case 'GERARDASCODBARRA22':
            return JSON.stringify({ periodoApuracao: '201901' }, null, 2);
        case 'ATUBENEFICIO23':
            return JSON.stringify(
                {
                    anoCalendario: 2021,
                    infoBeneficio: [
                        { periodoApuracao: '202101', indicadorBeneficio: true },
                        { periodoApuracao: '202102', indicadorBeneficio: true },
                    ],
                },
                null,
                2,
            );
        case 'DIVIDAATIVA24':
            return JSON.stringify({ anoCalendario: '2020' }, null, 2);
        case 'CONSEXTRATO16':
            return JSON.stringify({ numeroDas: '99999999999999999' }, null, 2);
        case 'TRANSDECLARACAO11':
            return JSON.stringify(
                {
                    periodoApuracao: '202401',
                    declaracao: {
                        /* preencher conforme layout oficial do serviço */
                    },
                },
                null,
                2,
            );
        default:
            return JSON.stringify(
                {
                    exemplo: 'preencha conforme documentação do serviço selecionado',
                },
                null,
                2,
            );
    }
}

export default function IntegraContadorServicosPage() {
    const toast = useRef<Toast>(null);
    const api = setupAPIClient(undefined);
    const catalogoSvc = useMemo(() => CatalogoServicoService(), []);
    const [ambiente, setAmbiente] = useState<Ambiente>('trial');
    const [secao, setSecao] = useState<string>('');
    const [servicoId, setServicoId] = useState<string>('');
    const [versaoSistema, setVersaoSistema] = useState('1.0');
    const [dadosJSON, setDadosJSON] = useState('{\n  "periodoApuracao": "201901"\n}');
    const [contratante, setContratante] = useState<NI>({ numero: '00000000000100', tipo: 2 });
    const [autorPedidoDados, setAutorPedidoDados] = useState<NI>({ numero: '00000000000100', tipo: 2 });
    const [contribuinte, setContribuinte] = useState<NI>({ numero: '00000000000100', tipo: 2 });
    const [accessTokenManual, setAccessTokenManual] = useState('');
    const [jwtTokenManual, setJwtTokenManual] = useState('');
    const [autenticarProcuradorToken, setAutenticarProcuradorToken] = useState('');
    const [retorno, setRetorno] = useState('');
    const [headersUltimaResposta, setHeadersUltimaResposta] = useState<Record<string, string>>({});
    const [loading, setLoading] = useState(false);

    const { data: catalogo = [] } = useQuery<CatalogoServico[]>({
        queryKey: ['integra-contador-catalogo-servicos'],
        queryFn: () => catalogoSvc.list({ incluirInativos: false }),
    });
    const { data: servicosProcuracao = [] } = useQuery<ServicoProcuracao[]>({
        queryKey: ['integra-contador-servicos-procuracao'],
        queryFn: async () => {
            const { data } = await api.get('/api/integra-contador/servicos-procuracao');
            return data?.servicos_procuracao ?? [];
        },
    });

    const secoes = useMemo(
        () =>
            Array.from(new Set(catalogo.map((s) => s.secao).filter(Boolean)))
                .sort((a, b) => a.localeCompare(b, 'pt-BR', { sensitivity: 'base' }))
                .map((s) => ({ label: s, value: s })),
        [catalogo],
    );

    const servicosFiltrados = useMemo(
        () =>
            catalogo
                .filter((s) => !secao || s.secao === secao)
                .sort((a, b) => {
                    if (a.sequencial !== b.sequencial) return a.sequencial - b.sequencial;
                    return a.codigo.localeCompare(b.codigo, 'pt-BR', { sensitivity: 'base' });
                })
                .map((s) => ({
                    label: (() => {
                        const proc = servicosProcuracao.find(
                            (p) =>
                                (p.id_sistema || '').trim().toUpperCase() === (s.id_sistema || '').trim().toUpperCase() &&
                                (p.id_servico || '').trim().toUpperCase() === (s.id_servico || '').trim().toUpperCase(),
                        );
                        const cod = (proc?.cod_procuracao || '').trim();
                        const flag = exigeProcuracao(cod) ? `[PROC ${cod}] ` : '';
                        return `${flag}${s.sequencial} - ${s.codigo} - ${s.descricao}`;
                    })(),
                    value: s.id,
                })),
        [catalogo, secao, servicosProcuracao],
    );

    const servicoSelecionado = useMemo(() => catalogo.find((s) => s.id === servicoId) ?? null, [catalogo, servicoId]);
    const servicoProcSelecionado = useMemo(() => {
        if (!servicoSelecionado) return null;
        return servicosProcuracao.find(
            (s) =>
                (s.id_sistema || '').trim().toUpperCase() === (servicoSelecionado.id_sistema || '').trim().toUpperCase() &&
                (s.id_servico || '').trim().toUpperCase() === (servicoSelecionado.id_servico || '').trim().toUpperCase(),
        ) ?? null;
    }, [servicoSelecionado, servicosProcuracao]);
    const servicoExigeProcuracao = useMemo(
        () => exigeProcuracao(servicoProcSelecionado?.cod_procuracao ?? ''),
        [servicoProcSelecionado?.cod_procuracao],
    );
    const presetAtual = useMemo(() => presetPorServico(servicoSelecionado?.id_servico ?? ''), [servicoSelecionado?.id_servico]);

    const autenticar = async () => {
        setLoading(true);
        try {
            const { data } = await api.post('/api/integra-contador/autenticar');
            setAccessTokenManual(data?.access_token ?? '');
            setJwtTokenManual(data?.jwt_token ?? '');
            setRetorno(JSON.stringify(data, null, 2));
            toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Autenticação SAPI executada.', life: 3000 });
        } catch (e: any) {
            const msg = e?.response?.data?.error || e?.response?.data?.message || 'Falha na autenticação SAPI';
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 7000 });
        } finally {
            setLoading(false);
        }
    };

    const executarServico = async () => {
        if (!servicoSelecionado) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Selecione um serviço do catálogo.', life: 4000 });
            return;
        }
        const operacao = (servicoSelecionado.tipo || '').trim().toLowerCase();
        if (!['apoiar', 'consultar', 'declarar', 'emitir', 'monitorar'].includes(operacao)) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Tipo do serviço inválido para chamada.', life: 4000 });
            return;
        }
        const autorDiferenteContribuinte = autorPedidoDados.numero.trim() !== contribuinte.numero.trim();
        if (servicoExigeProcuracao && autorDiferenteContribuinte && !autenticarProcuradorToken.trim()) {
            toast.current?.show({
                severity: 'warn',
                summary: 'Procuração obrigatória',
                detail: 'Este serviço exige procuração quando autor e contribuinte são diferentes. Informe o autenticar_procurador_token.',
                life: 7000,
            });
            return;
        }

        let dadosObj: unknown;
        try {
            dadosObj = JSON.parse(dadosJSON);
        } catch {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'O campo Dados deve conter JSON válido.', life: 4000 });
            return;
        }

        setLoading(true);
        try {
            const payload = {
                ambiente,
                operacao,
                access_token: accessTokenManual.trim(),
                jwt_token: jwtTokenManual.trim(),
                autenticar_procurador_token: autenticarProcuradorToken.trim(),
                payload: {
                    contratante,
                    autorPedidoDados,
                    contribuinte,
                    pedidoDados: {
                        idSistema: servicoSelecionado.id_sistema,
                        idServico: servicoSelecionado.id_servico,
                        versaoSistema: versaoSistema.trim() || '1.0',
                        dados: JSON.stringify(dadosObj),
                    },
                },
            };
            const { data } = await api.post<IntegraCallResponse>('/api/integra-contador/chamar', payload);
            setRetorno(JSON.stringify(data, null, 2));
            setHeadersUltimaResposta(data?.headers ?? {});
            const leitura = interpretarStatusHTTP(Number(data?.status_code || 0));
            toast.current?.show({
                severity: leitura.severidade,
                summary: leitura.titulo,
                detail: leitura.orientacao,
                life: 6000,
            });
        } catch (e: any) {
            const msg = e?.response?.data?.error || e?.response?.data?.message || 'Falha ao executar serviço';
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 7000 });
            setHeadersUltimaResposta({});
        } finally {
            setLoading(false);
        }
    };

    const aplicarPreset = () => {
        if (!servicoSelecionado) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Selecione um serviço para aplicar o modelo.', life: 3500 });
            return;
        }
        setDadosJSON(presetAtual);
        if (servicoSelecionado.id_sistema?.trim()) {
            setVersaoSistema('1.0');
        }
        toast.current?.show({ severity: 'info', summary: 'Modelo aplicado', detail: `Preset carregado para ${servicoSelecionado.id_servico}.`, life: 2500 });
    };

    const eTagResposta = headersUltimaResposta?.etag || headersUltimaResposta?.ETag || '';
    const tempoEsperaResposta =
        headersUltimaResposta?.tempoespera ||
        headersUltimaResposta?.tempoEspera ||
        headersUltimaResposta?.['x-tempo-espera'] ||
        '';

    const copiarRetorno = async () => {
        if (!retorno.trim()) {
            toast.current?.show({ severity: 'info', summary: 'Sem conteúdo', detail: 'Não há retorno para copiar.', life: 2500 });
            return;
        }
        try {
            await navigator.clipboard.writeText(retorno);
            toast.current?.show({ severity: 'success', summary: 'Copiado', detail: 'Retorno copiado para a área de transferência.', life: 2500 });
        } catch {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Não foi possível copiar automaticamente.', life: 3000 });
        }
    };

    return (
        <div className="grid">
            <div className="col-12">
                <Toast ref={toast} />
                <Card title="Integra Contador - Execução por Catálogo">
                    <p className="text-600 mt-0 mb-4">
                        Esta tela executa serviços do Integra Contador com base no cadastro de `Catálogo de Serviços`.
                    </p>
                    <div className="grid">
                        <div className="col-12 md:col-2">
                            <label htmlFor="ambiente" className="block mb-2 font-medium">Ambiente</label>
                            <Dropdown
                                id="ambiente"
                                value={ambiente}
                                options={[
                                    { label: 'Trial', value: 'trial' },
                                    { label: 'Produção', value: 'producao' },
                                ]}
                                optionLabel="label"
                                optionValue="value"
                                className="w-full"
                                onChange={(e) => setAmbiente(e.value as Ambiente)}
                            />
                        </div>
                        <div className="col-12 md:col-4">
                            <label htmlFor="secao" className="block mb-2 font-medium">Seção</label>
                            <Dropdown
                                id="secao"
                                value={secao}
                                options={secoes}
                                showClear
                                placeholder="Todas as seções"
                                className="w-full"
                                onChange={(e) => {
                                    setSecao(e.value ?? '');
                                    setServicoId('');
                                }}
                            />
                        </div>
                        <div className="col-12 md:col-6">
                            <label htmlFor="servico" className="block mb-2 font-medium">Serviço do catálogo</label>
                            <Dropdown
                                id="servico"
                                value={servicoId}
                                options={servicosFiltrados}
                                filter
                                showClear
                                placeholder="Selecione"
                                className="w-full"
                                onChange={(e) => setServicoId(e.value ?? '')}
                            />
                        </div>
                        <div className="col-12 md:col-3">
                            <label className="block mb-2 font-medium">Modelo de payload</label>
                            <Button type="button" label="Aplicar preset do serviço" icon="pi pi-bolt" className="w-full" onClick={aplicarPreset} disabled={!servicoSelecionado} />
                        </div>

                        <div className="col-12 md:col-4">
                            <label htmlFor="id-sistema" className="block mb-2 font-medium">idSistema</label>
                            <InputText id="id-sistema" className="w-full" readOnly value={servicoSelecionado?.id_sistema ?? ''} />
                        </div>
                        <div className="col-12 md:col-4">
                            <label htmlFor="id-servico" className="block mb-2 font-medium">idServico</label>
                            <InputText id="id-servico" className="w-full" readOnly value={servicoSelecionado?.id_servico ?? ''} />
                        </div>
                        <div className="col-12 md:col-2">
                            <label htmlFor="tipo-operacao" className="block mb-2 font-medium">Operação</label>
                            <InputText id="tipo-operacao" className="w-full" readOnly value={servicoSelecionado?.tipo ?? ''} />
                        </div>
                        <div className="col-12 md:col-2">
                            <label htmlFor="cod-procuracao" className="block mb-2 font-medium">Cód. procuração</label>
                            <InputText id="cod-procuracao" className="w-full" readOnly value={servicoProcSelecionado?.cod_procuracao ?? 'n/a'} />
                        </div>
                        <div className="col-12 md:col-2">
                            <label htmlFor="versao-sistema" className="block mb-2 font-medium">Versão</label>
                            <InputText id="versao-sistema" className="w-full" value={versaoSistema} onChange={(e) => setVersaoSistema(e.target.value)} />
                        </div>
                        <div className="col-12">
                            <label htmlFor="nome-servico-procuracao" className="block mb-2 font-medium">Nome do serviço (procuração eCAC)</label>
                            <InputText id="nome-servico-procuracao" className="w-full" readOnly value={servicoProcSelecionado?.nome_servico ?? 'n/a'} />
                            <div className="mt-2">
                                {servicoExigeProcuracao ? (
                                    <Tag value="Exige procuração" severity="warning" />
                                ) : (
                                    <Tag value="Sem exigência de procuração" severity="success" />
                                )}
                            </div>
                        </div>

                        <div className="col-12 md:col-4">
                            <label htmlFor="contratante-numero" className="block mb-2 font-medium">Contratante (número/tipo)</label>
                            <div className="p-inputgroup">
                                <InputText id="contratante-numero" value={contratante.numero} onChange={(e) => setContratante((p) => ({ ...p, numero: e.target.value.replace(/\D/g, '') }))} />
                                <InputText value={String(contratante.tipo)} onChange={(e) => setContratante((p) => ({ ...p, tipo: Number(e.target.value || 2) }))} style={{ maxWidth: '4rem' }} />
                            </div>
                        </div>
                        <div className="col-12 md:col-4">
                            <label htmlFor="autor-numero" className="block mb-2 font-medium">Autor pedido (número/tipo)</label>
                            <div className="p-inputgroup">
                                <InputText id="autor-numero" value={autorPedidoDados.numero} onChange={(e) => setAutorPedidoDados((p) => ({ ...p, numero: e.target.value.replace(/\D/g, '') }))} />
                                <InputText value={String(autorPedidoDados.tipo)} onChange={(e) => setAutorPedidoDados((p) => ({ ...p, tipo: Number(e.target.value || 2) }))} style={{ maxWidth: '4rem' }} />
                            </div>
                        </div>
                        <div className="col-12 md:col-4">
                            <label htmlFor="contribuinte-numero" className="block mb-2 font-medium">Contribuinte (número/tipo)</label>
                            <div className="p-inputgroup">
                                <InputText id="contribuinte-numero" value={contribuinte.numero} onChange={(e) => setContribuinte((p) => ({ ...p, numero: e.target.value.replace(/\D/g, '') }))} />
                                <InputText value={String(contribuinte.tipo)} onChange={(e) => setContribuinte((p) => ({ ...p, tipo: Number(e.target.value || 2) }))} style={{ maxWidth: '4rem' }} />
                            </div>
                        </div>

                        <div className="col-12 md:col-4">
                            <label htmlFor="access-token-manual" className="block mb-2 font-medium">Access token (opcional)</label>
                            <InputText id="access-token-manual" className="w-full" value={accessTokenManual} onChange={(e) => setAccessTokenManual(e.target.value)} />
                        </div>
                        <div className="col-12 md:col-4">
                            <label htmlFor="jwt-token-manual" className="block mb-2 font-medium">JWT token (opcional)</label>
                            <InputText id="jwt-token-manual" className="w-full" value={jwtTokenManual} onChange={(e) => setJwtTokenManual(e.target.value)} />
                        </div>
                        <div className="col-12 md:col-4">
                            <label htmlFor="autenticar-procurador-token" className="block mb-2 font-medium">autenticar_procurador_token (quando aplicável)</label>
                            <InputText
                                id="autenticar-procurador-token"
                                className="w-full"
                                value={autenticarProcuradorToken}
                                onChange={(e) => setAutenticarProcuradorToken(e.target.value)}
                            />
                        </div>

                        <div className="col-12">
                            <label htmlFor="dados-json" className="block mb-2 font-medium">Dados (JSON não escapado)</label>
                            <InputTextarea id="dados-json" className="w-full" rows={10} value={dadosJSON} onChange={(e) => setDadosJSON(e.target.value)} />
                            {servicoSelecionado && (
                                <small className="text-600 block mt-2">
                                    Serviço selecionado: <strong>{servicoSelecionado.id_servico}</strong> ({servicoSelecionado.descricao})
                                </small>
                            )}
                            <small className="text-600 block mt-1">
                                Prefixo `[PROC codigo]` no seletor indica exigência de procuração para o serviço.
                            </small>
                        </div>
                        <div className="col-12 flex gap-2 flex-wrap">
                            <Button type="button" label="Autenticar SAPI" icon="pi pi-key" onClick={autenticar} loading={loading} />
                            <Button type="button" label="Executar Serviço" icon="pi pi-play" severity="success" onClick={executarServico} loading={loading} />
                        </div>
                        <div className="col-12">
                            {(eTagResposta || tempoEsperaResposta) && (
                                <div className="mb-2 p-2 surface-100 border-round">
                                    <small className="block text-700">
                                        {tempoEsperaResposta ? `tempoEspera: ${tempoEsperaResposta}` : ''}
                                        {tempoEsperaResposta && eTagResposta ? ' | ' : ''}
                                        {eTagResposta ? `eTag: ${eTagResposta}` : ''}
                                    </small>
                                </div>
                            )}
                            <div className="flex align-items-center justify-content-between gap-2 mb-2">
                                <label htmlFor="retorno-api" className="font-medium m-0">Retorno da API</label>
                                <div className="flex gap-2">
                                    <Button type="button" label="Copiar" icon="pi pi-copy" text onClick={() => void copiarRetorno()} />
                                    <Button
                                        type="button"
                                        label="Limpar"
                                        icon="pi pi-trash"
                                        text
                                        severity="secondary"
                                        onClick={() => {
                                            setRetorno('');
                                            setHeadersUltimaResposta({});
                                        }}
                                    />
                                </div>
                            </div>
                            <InputTextarea id="retorno-api" className="w-full" rows={16} value={retorno} onChange={(e) => setRetorno(e.target.value)} />
                        </div>
                    </div>
                </Card>
            </div>
        </div>
    );
}

export const getServerSideProps = canSSRAuth(async (ctx) => {
    const apiClient = setupAPIClient(ctx);
    try {
        await apiClient.get('/api/registro');
    } catch (err: unknown) {
        const ax = err as { response?: { status?: number; data?: { error?: string } } };
        const msg = ax?.response?.data?.error ?? '';
        if (!(ax?.response?.status === 400 && msg.includes('no rows in result set'))) {
            return { redirect: { destination: '/', permanent: false } };
        }
    }

    const { data } = await apiClient.get('/api/usuariorole');
    const role = data?.logado?.role;
    if (role !== 'SUPER') {
        return { redirect: { destination: '/', permanent: false } };
    }

    return { props: {} };
});
