import React, { useCallback, useRef, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { Dialog } from 'primereact/dialog';
import { Dropdown } from 'primereact/dropdown';
import { InputText } from 'primereact/inputtext';
import { InputTextarea } from 'primereact/inputtextarea';
import { TabPanel, TabView } from 'primereact/tabview';
import { Toast } from 'primereact/toast';
import { classNames } from 'primereact/utils';
import { Tag } from 'primereact/tag';
import setupAPIClient from '../../components/api/api';
import { useCaixaPostal } from '../../components/context/CaixaPostalContext';
import CaixaPostalService, { CaixaPostalMensagem } from '../../services/cruds/CaixaPostalService';

type TenantOpcao = { id: string; nome: string };

const QUERY_KEY = 'caixa-postal-mensagens';

export default function CaixaPostal() {
    const toast = useRef<Toast>(null);
    const queryClient = useQueryClient();
    const { refreshCount } = useCaixaPostal();

    const { data: userRole } = useQuery<string>({
        queryKey: ['usuariorole-caixa-postal'],
        queryFn: async () => {
            const api = setupAPIClient(undefined);
            const { data } = await api.get('/api/usuariorole');
            return data?.logado?.role ?? '';
        },
    });

    const isSuper = userRole === 'SUPER';

    const [dialogEnviar, setDialogEnviar] = useState(false);
    const [titulo, setTitulo] = useState('');
    const [conteudo, setConteudo] = useState('');
    const [tenantSelecionado, setTenantSelecionado] = useState<TenantOpcao | null>(null);
    const [enviando, setEnviando] = useState(false);

    const svc = CaixaPostalService();

    const { data: mensagens = [], isLoading, refetch } = useQuery<CaixaPostalMensagem[]>({
        queryKey: [QUERY_KEY],
        queryFn: () => svc.listar(),
    });

    const { data: tenants = [] } = useQuery<TenantOpcao[]>({
        queryKey: ['tenants-opcoes-caixa-postal'],
        enabled: isSuper,
        queryFn: async () => {
            const api = setupAPIClient(undefined);
            const { data } = await api.get<{ id: string; nome: string }[]>('/api/tenants');
            return data ?? [];
        },
    });

    const inbox = mensagens.filter((m) => m.tipo === 'INBOX');
    const outbox = mensagens.filter((m) => m.tipo === 'OUTBOX');
    const naoLidasCount = inbox.filter((m) => !m.lida).length;

    const marcarLida = useCallback(
        async (msg: CaixaPostalMensagem) => {
            if (msg.lida) return;
            try {
                await svc.marcarComoLida(msg.id);
                queryClient.invalidateQueries({ queryKey: [QUERY_KEY] });
                refreshCount();
            } catch {
                toast.current?.show({ severity: 'error', summary: 'Erro', detail: 'Não foi possível marcar como lida.' });
            }
        },
        [queryClient, refreshCount, svc]
    );

    const abrirEnviar = () => {
        setTitulo('');
        setConteudo('');
        setTenantSelecionado(null);
        setDialogEnviar(true);
    };

    const enviar = async () => {
        if (!titulo.trim() || !conteudo.trim()) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Título e conteúdo são obrigatórios.' });
            return;
        }
        setEnviando(true);
        try {
            await svc.enviar({
                tenant_id: isSuper ? (tenantSelecionado?.id ?? '') : '',
                titulo: titulo.trim(),
                conteudo: conteudo.trim(),
            });
            toast.current?.show({ severity: 'success', summary: 'Enviado', detail: 'Mensagem enviada com sucesso.' });
            setDialogEnviar(false);
            queryClient.invalidateQueries({ queryKey: [QUERY_KEY] });
            refreshCount();
        } catch {
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: 'Não foi possível enviar a mensagem.' });
        } finally {
            setEnviando(false);
        }
    };

    const tenantOpcoes: TenantOpcao[] = [{ id: '', nome: 'Todos os tenants (global)' }, ...tenants];

    const formatarData = (dateStr: string) => {
        try {
            return new Intl.DateTimeFormat('pt-BR', {
                day: '2-digit',
                month: '2-digit',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
            }).format(new Date(dateStr));
        } catch {
            return dateStr;
        }
    };

    const renderMensagem = (msg: CaixaPostalMensagem) => (
        <div
            key={msg.id}
            className={classNames('caixa-postal-mensagem-card', {
                lida: msg.lida,
                outbox: msg.tipo === 'OUTBOX',
            })}
        >
            <div className="flex justify-content-between align-items-start mb-1">
                <div className="flex align-items-center gap-2">
                    {msg.tipo === 'INBOX' && !msg.lida && (
                        <Button
                            icon="pi pi-envelope"
                            rounded
                            text
                            severity="info"
                            tooltip="Marcar como lida"
                            tooltipOptions={{ position: 'top' }}
                            size="small"
                            onClick={() => marcarLida(msg)}
                        />
                    )}
                    {msg.tipo === 'INBOX' && msg.lida && (
                        <Button
                            icon="pi pi-envelope-open"
                            rounded
                            text
                            severity="secondary"
                            disabled
                            size="small"
                        />
                    )}
                    {msg.tipo === 'OUTBOX' && (
                        <Button
                            icon="pi pi-send"
                            rounded
                            text
                            severity="success"
                            disabled
                            size="small"
                        />
                    )}
                    <strong>{msg.titulo}</strong>
                    {msg.is_global && (
                        <Tag value="Global" severity="info" className="text-xs" />
                    )}
                    {!msg.lida && msg.tipo === 'INBOX' && (
                        <Tag value="Nova" severity="danger" className="text-xs" />
                    )}
                </div>
                <small className="text-color-secondary">{formatarData(msg.criado_em)}</small>
            </div>
            <div className="text-sm mb-1">
                <span className="text-color-secondary">
                    {msg.tipo === 'INBOX' ? `De: ${msg.remetente_nome}` : `Para: ${msg.is_global ? 'Todos' : 'VEC Sistemas'}`}
                </span>
            </div>
            <div className="text-sm" style={{ whiteSpace: 'pre-wrap' }}>
                {msg.conteudo}
            </div>
        </div>
    );

    const footerDialog = (
        <div className="flex justify-content-end gap-2">
            <Button label="Cancelar" icon="pi pi-times" severity="secondary" outlined onClick={() => setDialogEnviar(false)} />
            <Button label="Enviar" icon="pi pi-send" loading={enviando} onClick={enviar} />
        </div>
    );

    const paginatorLeft = (
        <Button
            type="button"
            icon="pi pi-refresh"
            tooltip="Atualizar"
            className="p-button-text"
            onClick={() => refetch()}
        />
    );

    return (
        <div className="grid">
            <Toast ref={toast} />
            <div className="col-12">
                <div className="card">
                    <div className="flex justify-content-between align-items-center mb-3">
                        <div>
                            <h5 className="m-0">Caixa Postal</h5>
                            {naoLidasCount > 0 && (
                                <small className="text-color-secondary">{naoLidasCount} mensagem(ns) não lida(s)</small>
                            )}
                        </div>
                        <Button
                            label={isSuper ? 'Nova Mensagem' : 'Enviar para VEC Sistemas'}
                            icon="pi pi-send"
                            onClick={abrirEnviar}
                        />
                    </div>

                    <TabView className="w-full caixa-postal-tabview"
                        pt={{
                            root: { style: { maxWidth: '100%' } },
                            navContainer: {
                                className: 'w-full',
                                style: { boxSizing: 'border-box', paddingLeft: 0, borderBottom: 'none' },
                            },
                            navContent: {
                                style: { flex: '1 1 auto', minWidth: 0, overflowX: 'auto', overflowY: 'hidden' },
                            },
                            inkbar: { style: { display: 'none' } },
                            nav: {
                                style: {
                                    display: 'flex',
                                    flexWrap: 'nowrap',
                                    width: '100%',
                                    justifyContent: 'flex-start',
                                    alignItems: 'flex-end',
                                    columnGap: '1.5rem',
                                    listStyle: 'none',
                                    margin: 0,
                                    paddingLeft: 0,
                                    paddingRight: 0,
                                },
                            },
                        }}
                    >
                        <TabPanel
                            headerStyle={{ whiteSpace: 'nowrap' }}
                            header={
                                <span>
                                    Recebidas
                                    {naoLidasCount > 0 && (
                                        <span className="ml-2 caixa-postal-badge" style={{ position: 'static', display: 'inline-flex' }}>
                                            {naoLidasCount}
                                        </span>
                                    )}
                                </span>
                            }
                        >
                            {isLoading ? (
                                <p>Carregando...</p>
                            ) : inbox.length === 0 ? (
                                <p className="text-color-secondary text-center mt-4">Nenhuma mensagem recebida.</p>
                            ) : (
                                <div>{inbox.map(renderMensagem)}</div>
                            )}
                        </TabPanel>
                        <TabPanel header="Enviadas" headerStyle={{ whiteSpace: 'nowrap' }}>
                            {isLoading ? (
                                <p>Carregando...</p>
                            ) : outbox.length === 0 ? (
                                <p className="text-color-secondary text-center mt-4">Nenhuma mensagem enviada.</p>
                            ) : (
                                <div>{outbox.map(renderMensagem)}</div>
                            )}
                        </TabPanel>
                    </TabView>
                </div>
            </div>

            <div className="fixed" style={{ bottom: '1.5rem', right: '1.5rem' }}>
                {paginatorLeft}
            </div>

            <Dialog
                visible={dialogEnviar}
                onHide={() => setDialogEnviar(false)}
                header={isSuper ? 'Enviar Aviso/Mensagem' : 'Enviar mensagem para VEC Sistemas'}
                style={{ width: '42rem' }}
                footer={footerDialog}
                draggable={false}
                resizable={false}
            >
                <div className="flex flex-column gap-3 pt-2">
                    {isSuper && (
                        <div className="field">
                            <label htmlFor="tenant-destino" className="font-medium">
                                Destinatário
                            </label>
                            <Dropdown
                                id="tenant-destino"
                                value={tenantSelecionado}
                                options={tenantOpcoes}
                                onChange={(e) => setTenantSelecionado(e.value)}
                                optionLabel="nome"
                                placeholder="Todos os tenants (global)"
                                className="w-full"
                            />
                        </div>
                    )}
                    <div className="field">
                        <label htmlFor="titulo-msg" className="font-medium">
                            Título <span className="text-red-500">*</span>
                        </label>
                        <InputText
                            id="titulo-msg"
                            value={titulo}
                            onChange={(e) => setTitulo(e.target.value)}
                            className="w-full"
                            maxLength={200}
                            placeholder="Título da mensagem"
                        />
                    </div>
                    <div className="field">
                        <label htmlFor="conteudo-msg" className="font-medium">
                            Conteúdo <span className="text-red-500">*</span>
                        </label>
                        <InputTextarea
                            id="conteudo-msg"
                            value={conteudo}
                            onChange={(e) => setConteudo(e.target.value)}
                            rows={6}
                            className="w-full"
                            placeholder="Digite sua mensagem..."
                            autoResize
                        />
                    </div>
                </div>
            </Dialog>
        </div>
    );
}
