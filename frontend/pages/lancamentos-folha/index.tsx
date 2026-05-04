import React, { useState, useEffect, useRef } from 'react';
import { TreeTable } from 'primereact/treetable';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';
import { Dialog } from 'primereact/dialog';
import { InputText } from 'primereact/inputtext';
import { InputNumber } from 'primereact/inputnumber';
import { InputTextarea } from 'primereact/inputtextarea';
import { Dropdown } from 'primereact/dropdown';
import { Toast } from 'primereact/toast';
import { ConfirmDialog, confirmDialog } from 'primereact/confirmdialog';
import LancamentoFolhaService from '../../services/cruds/LancamentoFolhaService';
import ClienteService from '../../services/cruds/ClienteService';

const formatMoney = (n: number): string => {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(n);
};

type TreeNode = {
    key: string;
    data: {
        tipo: 'cliente' | 'lancamento';
        cliente_id?: string;
        nome?: string;
        documento?: string;
        total_folha?: number;
        total_faturamento?: number;
        id?: string;
        competencia?: string;
        valor_folha?: number;
        valor_faturamento?: number;
        observacoes?: string;
    };
    children?: TreeNode[];
    leaf?: boolean;
};

type ClienteOption = {
    id: string;
    nome: string;
    documento: string;
};

const LancamentosFolha = () => {
    const toast = useRef<Toast>(null);
    const [treeData, setTreeData] = useState<TreeNode[]>([]);
    const [loading, setLoading] = useState(false);
    const [dialogVisible, setDialogVisible] = useState(false);
    const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create');
    const [selectedCliente, setSelectedCliente] = useState<ClienteOption | null>(null);
    const [clientesOptions, setClientesOptions] = useState<ClienteOption[]>([]);
    const [selectedNode, setSelectedNode] = useState<TreeNode | null>(null);
    const [showClienteSelector, setShowClienteSelector] = useState(false);

    const [formData, setFormData] = useState({
        id: '',
        cliente_id: '',
        cliente_nome: '',
        competencia: '',
        valor_folha: 0,
        valor_faturamento: 0,
        observacoes: '',
    });

    const lancamentoFolhaServiceRef = useRef(LancamentoFolhaService());
    const clienteServiceRef = useRef(ClienteService());

    const loadTree = async () => {
        setLoading(true);
        try {
            const response = await lancamentoFolhaServiceRef.current.getTree();
            setTreeData(response.data.tree || []);
        } catch (err) {
            toast.current?.show({
                severity: 'error',
                summary: 'Erro',
                detail: 'Erro ao carregar lançamentos de folha',
                life: 3000,
            });
        } finally {
            setLoading(false);
        }
    };

    const loadClientes = async () => {
        try {
            const clientes = await clienteServiceRef.current.list(500, 0);
            const options: ClienteOption[] = (clientes as Array<{ id: string; nome: string; documento: string }>).map((c) => ({
                id: c.id,
                nome: c.nome,
                documento: c.documento || '',
            }));
            setClientesOptions(options);
        } catch (err) {
            console.error('Erro ao carregar clientes:', err);
        }
    };

    useEffect(() => {
        loadTree();
        loadClientes();
    }, []);

    const openNewDialog = () => {
        setDialogMode('create');
        setSelectedCliente(null);
        setFormData({
            id: '',
            cliente_id: '',
            cliente_nome: '',
            competencia: '',
            valor_folha: 0,
            valor_faturamento: 0,
            observacoes: '',
        });
        setShowClienteSelector(true);
        setDialogVisible(true);
    };

    const openNewForCliente = (clienteId: string, clienteNome: string) => {
        setDialogMode('create');
        setSelectedCliente({ id: clienteId, nome: clienteNome, documento: '' });
        setFormData({
            id: '',
            cliente_id: clienteId,
            cliente_nome: clienteNome,
            competencia: '',
            valor_folha: 0,
            valor_faturamento: 0,
            observacoes: '',
        });
        setShowClienteSelector(false);
        setDialogVisible(true);
    };

    const confirmCliente = () => {
        if (!selectedCliente) return;
        setFormData((prev) => ({
            ...prev,
            cliente_id: selectedCliente.id,
            cliente_nome: selectedCliente.nome,
        }));
        setShowClienteSelector(false);
    };

    const openEditDialog = async (node: TreeNode) => {
        if (node.data.tipo !== 'lancamento' || !node.data.id) return;

        // Tenta encontrar o nome do cliente a partir do nó pai
        const parentNome = treeData.find((t) =>
            t.children?.some((c) => c.key === node.key)
        )?.data?.nome || '';

        setDialogMode('edit');
        setShowClienteSelector(false);
        setFormData({
            id: node.data.id,
            cliente_id: node.data.cliente_id || '',
            cliente_nome: parentNome || node.data.cliente_id || '',
            competencia: node.data.competencia || '',
            valor_folha: node.data.valor_folha || 0,
            valor_faturamento: node.data.valor_faturamento || 0,
            observacoes: node.data.observacoes || '',
        });

        setDialogVisible(true);
    };

    const confirmDelete = (node: TreeNode) => {
        if (node.data.tipo !== 'lancamento' || !node.data.id) return;

        confirmDialog({
            message: `Deseja realmente excluir o lançamento de ${node.data.competencia}?`,
            header: 'Confirmar Exclusao',
            icon: 'pi pi-exclamation-triangle',
            acceptLabel: 'Sim',
            rejectLabel: 'Nao',
            accept: async () => {
                try {
                    await lancamentoFolhaServiceRef.current.deleteLancamento(node.data.id!);
                    toast.current?.show({
                        severity: 'success',
                        summary: 'Sucesso',
                        detail: 'Lançamento removido com sucesso',
                        life: 3000,
                    });
                    loadTree();
                } catch (err) {
                    toast.current?.show({
                        severity: 'error',
                        summary: 'Erro',
                        detail: 'Erro ao excluir lançamento',
                        life: 3000,
                    });
                }
            },
        });
    };

    const saveLancamento = async () => {
        try {
            if (dialogMode === 'create') {
                if (!formData.cliente_id) {
                    toast.current?.show({
                        severity: 'warn',
                        summary: 'Atencao',
                        detail: 'Selecione um cliente',
                        life: 3000,
                    });
                    return;
                }

                if (!formData.competencia) {
                    toast.current?.show({
                        severity: 'warn',
                        summary: 'Atencao',
                        detail: 'Informe a competencia (MM/AAAA)',
                        life: 3000,
                    });
                    return;
                }

                await lancamentoFolhaServiceRef.current.createLancamento({
                    cliente_id: formData.cliente_id,
                    competencia: formData.competencia,
                    valor_folha: formData.valor_folha,
                    valor_faturamento: formData.valor_faturamento,
                    observacoes: formData.observacoes,
                });
                toast.current?.show({
                    severity: 'success',
                    summary: 'Sucesso',
                    detail: 'Lançamento criado com sucesso',
                    life: 3000,
                });
            } else {
                await lancamentoFolhaServiceRef.current.updateLancamento({
                    id: formData.id,
                    competencia: formData.competencia,
                    valor_folha: formData.valor_folha,
                    valor_faturamento: formData.valor_faturamento,
                    observacoes: formData.observacoes,
                });
                toast.current?.show({
                    severity: 'success',
                    summary: 'Sucesso',
                    detail: 'Lançamento atualizado com sucesso',
                    life: 3000,
                });
            }

            setDialogVisible(false);
            loadTree();
        } catch (err: any) {
            const msg = err?.response?.data?.error || err?.message || 'Erro ao salvar lançamento';
            toast.current?.show({
                severity: 'error',
                summary: 'Erro',
                detail: msg,
                life: 5000,
            });
        }
    };

    const actionBodyTemplate = (node: TreeNode) => {
        if (node.data.tipo !== 'lancamento') return null;

        return (
            <div className="flex gap-2">
                <Button
                    icon="pi pi-pencil"
                    rounded
                    severity="success"
                    className="mr-2"
                    tooltip="Alterar"
                    tooltipOptions={{ position: 'top' }}
                    onClick={() => openEditDialog(node)}
                />
                <Button
                    icon="pi pi-trash"
                    rounded
                    severity="warning"
                    tooltip="Excluir"
                    tooltipOptions={{ position: 'top' }}
                    onClick={() => confirmDelete(node)}
                />
            </div>
        );
    };

    const clienteActionTemplate = (node: TreeNode) => {
        if (node.data.tipo !== 'cliente') return null;

        return (
            <Button
                icon="pi pi-plus"
                rounded
                severity="info"
                tooltip="Novo Lançamento"
                tooltipOptions={{ position: 'top' }}
                onClick={() => openNewForCliente(node.data.cliente_id!, node.data.nome || '')}
            />
        );
    };

    const valorFolhaTemplate = (node: TreeNode) => {
        if (node.data.tipo === 'cliente') {
            return <span className="font-bold">{formatMoney(node.data.total_folha || 0)}</span>;
        }
        return <span>{formatMoney(node.data.valor_folha || 0)}</span>;
    };

    const valorFaturamentoTemplate = (node: TreeNode) => {
        if (node.data.tipo === 'cliente') {
            return <span className="font-bold">{formatMoney(node.data.total_faturamento || 0)}</span>;
        }
        return <span>{formatMoney(node.data.valor_faturamento || 0)}</span>;
    };

    const nomeTemplate = (node: TreeNode) => {
        if (node.data.tipo === 'cliente') {
            return (
                <span className="font-bold">
                    {node.data.nome}
                    {node.data.documento ? ` (${node.data.documento})` : ''}
                </span>
            );
        }
        return <span className="pl-3">{node.data.competencia}</span>;
    };

    return (
        <div className="p-4">
            <Toast ref={toast} />
            <ConfirmDialog />

            <div className="flex justify-content-between align-items-center mb-4">
                <h1 className="text-2xl font-bold m-0">Lançamentos da Folha</h1>
                <Button
                    label="Novo Lançamento"
                    icon="pi pi-plus"
                    severity="success"
                    onClick={openNewDialog}
                />
            </div>

            <TreeTable
                value={treeData}
                loading={loading}
                lazy
                scrollable
                scrollHeight="flex"
                sortMode="single"
                selectionMode="single"
                selectionKeys={selectedNode?.key || ''}
                onSelectionChange={(e) => setSelectedNode(e.value as TreeNode | null)}
                className="p-treetable-sm"
                globalFilter={null}
                emptyMessage="Nenhum lançamento encontrado."
            >
                <Column
                    field="nome"
                    header="Cliente / Competencia"
                    expander
                    body={nomeTemplate}
                    style={{ minWidth: '300px' }}
                />
                <Column
                    field="valor_folha"
                    header="Valor da Folha"
                    body={valorFolhaTemplate}
                    style={{ minWidth: '150px' }}
                    sortable
                />
                <Column
                    field="valor_faturamento"
                    header="Faturamento Mensal"
                    body={valorFaturamentoTemplate}
                    style={{ minWidth: '150px' }}
                    sortable
                />
                <Column
                    field="observacoes"
                    header="Observações"
                    style={{ minWidth: '200px' }}
                />
                <Column
                    header="Ações"
                    body={actionBodyTemplate}
                    style={{ minWidth: '120px' }}
                />
                <Column
                    header=""
                    body={clienteActionTemplate}
                    style={{ minWidth: '60px' }}
                />
            </TreeTable>

            <Dialog
                header={dialogMode === 'create' ? 'Novo Lançamento de Folha' : 'Editar Lançamento de Folha'}
                visible={dialogVisible}
                style={{ width: '580px' }}
                modal
                onHide={() => setDialogVisible(false)}
                footer={
                    <div className="flex justify-content-end gap-2">
                        <Button
                            label="Cancelar"
                            icon="pi pi-times"
                            severity="secondary"
                            onClick={() => setDialogVisible(false)}
                        />
                        <Button
                            label="Salvar"
                            icon="pi pi-check"
                            severity="success"
                            onClick={saveLancamento}
                        />
                    </div>
                }
            >
                <div className="grid formgrid p-fluid">
                    <div className="field col-12 flex align-items-center mb-3">
                        <label htmlFor="cliente" className="col-4 text-right mr-2 font-bold">Cliente:</label>
                        <div className="col-8 p-0">
                            {showClienteSelector ? (
                                <div className="flex gap-2">
                                    <Dropdown
                                        id="cliente"
                                        value={selectedCliente}
                                        onChange={(e) => setSelectedCliente(e.value)}
                                        options={clientesOptions}
                                        optionLabel="nome"
                                        placeholder="Selecione um cliente"
                                        filter
                                        className="flex-1"
                                    />
                                    <Button
                                        icon="pi pi-check"
                                        severity="success"
                                        disabled={!selectedCliente}
                                        onClick={confirmCliente}
                                    />
                                </div>
                            ) : (
                                <InputText
                                    id="cliente_display"
                                    value={formData.cliente_nome || formData.cliente_id}
                                    disabled
                                />
                            )}
                        </div>
                    </div>

                    <div className="field col-12 flex align-items-center mb-3">
                        <label htmlFor="competencia" className="col-4 text-right mr-2 font-bold">Competência:</label>
                        <div className="col-8 p-0">
                            <InputText
                                id="competencia"
                                value={formData.competencia}
                                onChange={(e) => setFormData((prev) => ({ ...prev, competencia: e.target.value }))}
                                placeholder="MM/AAAA"
                            />
                        </div>
                    </div>

                    <div className="field col-12 flex align-items-center mb-3">
                        <label htmlFor="valor_folha" className="col-4 text-right mr-2 font-bold">Valor Total da Folha (R$):</label>
                        <div className="col-8 p-0">
                            <InputNumber
                                id="valor_folha"
                                value={formData.valor_folha}
                                onValueChange={(e) => setFormData((prev) => ({ ...prev, valor_folha: e.value || 0 }))}
                                mode="currency"
                                currency="BRL"
                                locale="pt-BR"
                                minFractionDigits={2}
                                inputStyle={{ width: '100%' }}
                            />
                        </div>
                    </div>

                    <div className="field col-12 flex align-items-center mb-3">
                        <label htmlFor="valor_faturamento" className="col-4 text-right mr-2 font-bold">Faturamento Mensal (R$):</label>
                        <div className="col-8 p-0">
                            <InputNumber
                                id="valor_faturamento"
                                value={formData.valor_faturamento}
                                onValueChange={(e) => setFormData((prev) => ({ ...prev, valor_faturamento: e.value || 0 }))}
                                mode="currency"
                                currency="BRL"
                                locale="pt-BR"
                                minFractionDigits={2}
                                inputStyle={{ width: '100%' }}
                            />
                        </div>
                    </div>

                    <div className="field col-12 flex align-items-start mb-3">
                        <label htmlFor="observacoes" className="col-4 text-right mr-2 font-bold" style={{ marginTop: '0.5rem' }}>Observações:</label>
                        <div className="col-8 p-0">
                            <InputTextarea
                                id="observacoes"
                                value={formData.observacoes}
                                onChange={(e) => setFormData((prev) => ({ ...prev, observacoes: e.target.value }))}
                                rows={3}
                                style={{ width: '100%', minWidth: '280px' }}
                            />
                        </div>
                    </div>
                </div>
            </Dialog>

            <Button
                icon="pi pi-refresh"
                className="p-button-rounded p-button-text p-button-lg"
                style={{
                    position: 'fixed',
                    bottom: '2rem',
                    right: '2rem',
                    zIndex: 100,
                }}
                tooltip="Atualizar"
                tooltipOptions={{ position: 'left' }}
                onClick={loadTree}
            />
        </div>
    );
};

export default LancamentosFolha;
