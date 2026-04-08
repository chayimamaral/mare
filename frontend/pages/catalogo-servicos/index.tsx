import React, { useMemo, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { TreeTable } from 'primereact/treetable';
import { TreeNode } from 'primereact/treenode';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';
import { Dialog } from 'primereact/dialog';
import { InputText } from 'primereact/inputtext';
import { Dropdown } from 'primereact/dropdown';
import { Toast } from 'primereact/toast';
import { ConfirmDialog, confirmDialog } from 'primereact/confirmdialog';
import CatalogoServicoService, { CatalogoServico } from '../../services/cruds/CatalogoServicoService';
import { canSSRAuth } from '../../components/utils/canSSRAuth';
import setupAPIClient from '../../components/api/api';

type FormState = Omit<CatalogoServico, 'id'>;

const SECOES_FIXAS = [
    'Integra-SN',
    'Integra-MEI',
    'Integra-DCTFWeb',
    'Integra-Procurações',
    'Integra-Sicalc',
    'Integra-CaixaPostal',
    'Integra-Pagamento',
    'Integra-Contador-Gerenciador',
    'Integra-SITFIS',
    'Integra-Parcelamentos',
    'Integra-Redesim',
    'Integra-e-Processo',
];

const emptyForm: FormState = {
    secao: '',
    sequencial: 1,
    codigo: '',
    id_sistema: '',
    id_servico: '',
    situacao_implantacao: '',
    data_implantacao: '',
    tipo: '',
    descricao: '',
};

function toNodes(items: CatalogoServico[]): TreeNode[] {
    const agrupado = new Map<string, CatalogoServico[]>();
    for (const item of items) {
        const key = item.secao?.trim() || 'Sem seção';
        const lista = agrupado.get(key) ?? [];
        lista.push(item);
        agrupado.set(key, lista);
    }

    const secoesOrdenadas = Array.from(agrupado.keys()).sort((a, b) => a.localeCompare(b, 'pt-BR', { sensitivity: 'base' }));
    return secoesOrdenadas.map((secao) => {
        const children = (agrupado.get(secao) ?? [])
            .sort((a, b) => a.sequencial - b.sequencial)
            .map((s) => ({
                key: s.id,
                leaf: true,
                data: {
                    ...s,
                    isSecao: false,
                },
            }));
        return {
            key: `secao:${secao}`,
            leaf: false,
            data: {
                secao,
                descricao: `${children.length} serviço(s)`,
                isSecao: true,
            },
            children,
        } as TreeNode;
    });
}

export default function CatalogoServicosPage() {
    const toast = useRef<Toast>(null);
    const svc = useMemo(() => CatalogoServicoService(), []);
    const [secaoFiltro, setSecaoFiltro] = useState<string>('TODAS');
    const [dialogVisible, setDialogVisible] = useState(false);
    const [submitted, setSubmitted] = useState(false);
    const [editingId, setEditingId] = useState<string | null>(null);
    const [form, setForm] = useState<FormState>(emptyForm);
    const { data: roleData, isFetching: isFetchingRole } = useQuery({
        queryKey: ['catalogo-servicos-user-role'],
        queryFn: async () => {
            const api = setupAPIClient(undefined);
            const { data } = await api.get('/api/usuariorole');
            return data?.logado?.role ?? '';
        },
        staleTime: 0,
        gcTime: 0,
        refetchOnMount: 'always',
        refetchOnWindowFocus: true,
        refetchOnReconnect: 'always',
    });
    const podeManter = !isFetchingRole && roleData === 'SUPER';

    const { data, isFetching, refetch } = useQuery<CatalogoServico[]>({
        queryKey: ['catalogo-servicos', secaoFiltro],
        queryFn: () => svc.list(secaoFiltro === 'TODAS' ? '' : secaoFiltro),
    });

    const nodes = useMemo(() => toNodes(data ?? []), [data]);

    const opcoesSecao = useMemo(() => {
        return [{ label: 'Todas as seções', value: 'TODAS' }, ...SECOES_FIXAS.map((s) => ({ label: s, value: s }))];
    }, []);

    const abrirNovo = () => {
        if (!podeManter) return;
        setSubmitted(false);
        setEditingId(null);
        setForm({ ...emptyForm, secao: secaoFiltro !== 'TODAS' ? secaoFiltro : '' });
        setDialogVisible(true);
    };

    const abrirEditar = (row: CatalogoServico) => {
        if (!podeManter) return;
        setSubmitted(false);
        setEditingId(row.id);
        setForm({
            secao: row.secao,
            sequencial: row.sequencial,
            codigo: row.codigo,
            id_sistema: row.id_sistema,
            id_servico: row.id_servico,
            situacao_implantacao: row.situacao_implantacao,
            data_implantacao: row.data_implantacao || '',
            tipo: row.tipo,
            descricao: row.descricao,
        });
        setDialogVisible(true);
    };

    const excluir = (id: string) => {
        if (!podeManter) return;
        confirmDialog({
            header: 'Confirmar exclusão',
            message: 'Deseja excluir este serviço do catálogo?',
            icon: 'pi pi-exclamation-triangle',
            acceptLabel: 'Excluir',
            rejectLabel: 'Cancelar',
            acceptClassName: 'p-button-danger',
            accept: async () => {
                try {
                    await svc.remove(id);
                    toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Serviço excluído.', life: 3000 });
                    await refetch();
                } catch (e: any) {
                    const msg = e?.response?.data?.error || e?.response?.data?.message || 'Falha ao excluir';
                    toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 5000 });
                }
            },
        });
    };

    const salvar = async () => {
        if (!podeManter) return;
        setSubmitted(true);
        if (!form.secao.trim() || form.sequencial <= 0 || !form.codigo.trim() || !form.id_sistema.trim() || !form.id_servico.trim() || !form.situacao_implantacao.trim() || !form.tipo.trim() || !form.descricao.trim()) {
            return;
        }
        try {
            if (editingId) {
                await svc.update({ id: editingId, ...form });
            } else {
                await svc.create(form);
            }
            toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Serviço salvo.', life: 3000 });
            setDialogVisible(false);
            setForm(emptyForm);
            setEditingId(null);
            await refetch();
        } catch (e: any) {
            const msg = e?.response?.data?.error || e?.response?.data?.message || 'Falha ao salvar';
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 5000 });
        }
    };

    const isInvalid = (v: string) => submitted && !v.trim();

    const descricaoTemplate = (node: TreeNode) => {
        const d = node.data as any;
        return d.isSecao ? <strong>{d.descricao}</strong> : d.descricao;
    };

    const acoesTemplate = (node: TreeNode) => {
        const d = node.data as any;
        if (d.isSecao || !podeManter) return null;
        return (
            <div className="flex gap-2">
                <Button type="button" icon="pi pi-pencil" rounded severity="success" onClick={() => abrirEditar(d as CatalogoServico)} />
                <Button type="button" icon="pi pi-trash" rounded severity="warning" onClick={() => excluir(d.id)} />
            </div>
        );
    };

    return (
        <div className="grid">
            <div className="col-12">
                <div className="card vecontab-catalogo-servico-card">
                    <Toast ref={toast} />
                    <ConfirmDialog />
                    <h1 className="text-2xl font-bold text-900 m-0 mb-3">Catálogo de Serviços - Integra Contador</h1>
                    <div className="flex flex-wrap gap-3 align-items-end mb-3">
                        <div className="flex flex-column gap-1" style={{ minWidth: '18rem' }}>
                            <label htmlFor="filtro-secao" className="text-sm font-semibold">Seção</label>
                            <Dropdown
                                inputId="filtro-secao"
                                value={secaoFiltro}
                                options={opcoesSecao}
                                optionLabel="label"
                                optionValue="value"
                                onChange={(e) => setSecaoFiltro(e.value)}
                                className="w-full"
                            />
                        </div>
                        {podeManter && (
                            <Button type="button" label="Incluir" icon="pi pi-plus" severity="success" onClick={abrirNovo} />
                        )}
                    </div>
                    <TreeTable value={nodes} stripedRows loading={isFetching} tableStyle={{ minWidth: '72rem' }}>
                        <Column field="secao" header="Seção" expander sortable style={{ minWidth: '14rem' }} />
                        <Column field="sequencial" header="Sequencial" sortable style={{ width: '8rem' }} />
                        <Column field="codigo" header="Código" sortable style={{ width: '8rem' }} />
                        <Column field="id_sistema" header="idSistema" sortable style={{ minWidth: '10rem' }} />
                        <Column field="id_servico" header="idServico" sortable style={{ minWidth: '12rem' }} />
                        <Column field="situacao_implantacao" header="Situação e Data" sortable style={{ minWidth: '12rem' }} />
                        <Column field="tipo" header="Tipo" sortable style={{ width: '9rem' }} />
                        <Column header="Descrição" body={descricaoTemplate} sortable field="descricao" style={{ minWidth: '16rem' }} />
                        {podeManter && <Column header="Ações" body={acoesTemplate} style={{ width: '8rem' }} />}
                    </TreeTable>
                    <Dialog
                        visible={dialogVisible}
                        onHide={() => setDialogVisible(false)}
                        header={editingId ? 'Alterar serviço' : 'Novo serviço'}
                        style={{ width: 'min(96vw, 52rem)' }}
                        modal
                        footer={
                            <div className="flex gap-2 justify-content-end">
                                <Button type="button" label="Cancelar" text onClick={() => setDialogVisible(false)} />
                                <Button type="button" label="Salvar" icon="pi pi-check" onClick={() => void salvar()} />
                            </div>
                        }
                    >
                        <div className="grid">
                            <div className="col-12 md:col-4">
                                <label htmlFor="secao">Seção</label>
                                <Dropdown
                                    id="secao"
                                    value={form.secao}
                                    options={SECOES_FIXAS.map((s) => ({ label: s, value: s }))}
                                    optionLabel="label"
                                    optionValue="value"
                                    onChange={(e) => setForm((p) => ({ ...p, secao: e.value ?? '' }))}
                                    className={`w-full ${isInvalid(form.secao) ? 'p-invalid' : ''}`}
                                />
                            </div>
                            <div className="col-12 md:col-2">
                                <label htmlFor="sequencial">Sequencial</label>
                                <InputText id="sequencial" type="number" value={String(form.sequencial)} onChange={(e) => setForm((p) => ({ ...p, sequencial: Number(e.target.value || 0) }))} />
                            </div>
                            <div className="col-12 md:col-2">
                                <label htmlFor="codigo">Código</label>
                                <InputText id="codigo" value={form.codigo} onChange={(e) => setForm((p) => ({ ...p, codigo: e.target.value }))} className={isInvalid(form.codigo) ? 'p-invalid' : ''} />
                            </div>
                            <div className="col-12 md:col-4">
                                <label htmlFor="tipo">Tipo</label>
                                <InputText id="tipo" value={form.tipo} onChange={(e) => setForm((p) => ({ ...p, tipo: e.target.value }))} className={isInvalid(form.tipo) ? 'p-invalid' : ''} />
                            </div>
                            <div className="col-12 md:col-4">
                                <label htmlFor="idsistema">idSistema</label>
                                <InputText id="idsistema" value={form.id_sistema} onChange={(e) => setForm((p) => ({ ...p, id_sistema: e.target.value }))} className={isInvalid(form.id_sistema) ? 'p-invalid' : ''} />
                            </div>
                            <div className="col-12 md:col-4">
                                <label htmlFor="idservico">idServico</label>
                                <InputText id="idservico" value={form.id_servico} onChange={(e) => setForm((p) => ({ ...p, id_servico: e.target.value }))} className={isInvalid(form.id_servico) ? 'p-invalid' : ''} />
                            </div>
                            <div className="col-12 md:col-4">
                                <label htmlFor="data_implantacao">Data de Implantação</label>
                                <InputText id="data_implantacao" type="date" value={form.data_implantacao || ''} onChange={(e) => setForm((p) => ({ ...p, data_implantacao: e.target.value }))} />
                            </div>
                            <div className="col-12">
                                <label htmlFor="situacao">Situação e Data da Implantação</label>
                                <InputText id="situacao" value={form.situacao_implantacao} onChange={(e) => setForm((p) => ({ ...p, situacao_implantacao: e.target.value }))} className={isInvalid(form.situacao_implantacao) ? 'p-invalid' : ''} />
                            </div>
                            <div className="col-12">
                                <label htmlFor="descricao">Descrição</label>
                                <InputText id="descricao" value={form.descricao} onChange={(e) => setForm((p) => ({ ...p, descricao: e.target.value }))} className={isInvalid(form.descricao) ? 'p-invalid' : ''} />
                            </div>
                        </div>
                    </Dialog>
                    <div className="vecontab-catalogo-servico-fab-wrap">
                        <Button
                            type="button"
                            icon="pi pi-refresh"
                            tooltip="Atualizar"
                            tooltipOptions={{ position: 'left' }}
                            className="p-button-rounded p-button-text"
                            loading={isFetching}
                            onClick={() => void refetch()}
                        />
                    </div>
                </div>
            </div>
            <style jsx global>{`
                .vecontab-catalogo-servico-card {
                    position: relative;
                    padding-bottom: 3rem;
                }
                .vecontab-catalogo-servico-fab-wrap {
                    position: absolute;
                    right: 1rem;
                    bottom: 0.75rem;
                    z-index: 2;
                }
            `}</style>
        </div>
    );
}

export const getServerSideProps = canSSRAuth(async (ctx) => {
    try {
        const apiClient = setupAPIClient(ctx);
        await apiClient.get('/api/registro');
        return { props: {} };
    } catch (err) {
        console.log(err);
        return {
            redirect: {
                destination: '/',
                permanent: false,
            },
        };
    }
});
