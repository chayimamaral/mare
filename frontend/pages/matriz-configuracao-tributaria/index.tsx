import React, { useState, useEffect, useRef } from 'react';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';
import { Dialog } from 'primereact/dialog';
import { InputText } from 'primereact/inputtext';
import { InputSwitch } from 'primereact/inputswitch';
import { InputNumber } from 'primereact/inputnumber';
import { Dropdown } from 'primereact/dropdown';
import { Toast } from 'primereact/toast';
import { ConfirmDialog, confirmDialog } from 'primereact/confirmdialog';
import MatrizConfiguracaoTributariaService from '../../services/cruds/MatrizConfiguracaoTributariaService';
import TipoEmpresaService from '../../services/cruds/TipoEmpresaService';
import EnquadramentoJuridicoPorteService from '../../services/cruds/EnquadramentoJuridicoPorteService';
import RegimeTributarioService from '../../services/cruds/RegimeTributarioService';

type MatrizItem = {
    id: string;
    nome: string;
    natureza_juridica_id: string;
    natureza_juridica: string;
    enquadramento_porte_id: string;
    enquadramento_porte: string;
    regime_tributario_id: string;
    regime_tributario: string;
    aliquota_base: number;
    possui_fator_r: boolean;
    aliquota_fator_r: number;
    ativo: boolean;
};

type DropdownOption = {
    id: string;
    descricao?: string;
    nome?: string;
    sigla?: string;
};

const formatAliquota = (n: number): string => {
    return `${(n * 100).toFixed(2)}%`;
};

const MatrizConfiguracaoTributaria = () => {
    const toast = useRef<Toast>(null);
    const [items, setItems] = useState<MatrizItem[]>([]);
    const [loading, setLoading] = useState(false);
    const [totalRecords, setTotalRecords] = useState(0);
    const [first, setFirst] = useState(0);
    const [rows, setRows] = useState(25);
    const [sortField, setSortField] = useState('nome');
    const [sortOrder, setSortOrder] = useState(1);
    const [nomeFilter, setNomeFilter] = useState('');

    const [dialogVisible, setDialogVisible] = useState(false);
    const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create');
    const [submitting, setSubmitting] = useState(false);

    const [naturezaOptions, setNaturezaOptions] = useState<DropdownOption[]>([]);
    const [porteOptions, setPorteOptions] = useState<DropdownOption[]>([]);
    const [regimeOptions, setRegimeOptions] = useState<DropdownOption[]>([]);

    const [formData, setFormData] = useState({
        id: '',
        nome: '',
        natureza_juridica_id: null as DropdownOption | null,
        enquadramento_porte_id: null as DropdownOption | null,
        regime_tributario_id: null as DropdownOption | null,
        aliquota_base: 0,
        possui_fator_r: false,
        aliquota_fator_r: 0,
        ativo: true,
    });

    const serviceRef = useRef(MatrizConfiguracaoTributariaService());
    const tipoEmpresaServiceRef = useRef(TipoEmpresaService());
    const porteServiceRef = useRef(EnquadramentoJuridicoPorteService());
    const regimeServiceRef = useRef(RegimeTributarioService());

    const loadData = async () => {
        setLoading(true);
        try {
            const response = await serviceRef.current.list({
                lazyEvent: JSON.stringify({
                    first,
                    rows,
                    sortField,
                    sortOrder,
                    filters: nomeFilter ? { nome: { value: nomeFilter, matchMode: 'contains' } } : {},
                }),
            });
            setItems(response.data.items || []);
            setTotalRecords(response.data.totalRecords || 0);
        } catch (err) {
            toast.current?.show({
                severity: 'error',
                summary: 'Erro',
                detail: 'Erro ao carregar matriz tributaria',
                life: 3000,
            });
        } finally {
            setLoading(false);
        }
    };

    const loadDropdowns = async () => {
        try {
            const [naturezaRes, porteRes, regimeRes] = await Promise.all([
                tipoEmpresaServiceRef.current.getTiposEmpresaLite(),
                porteServiceRef.current.list(),
                regimeServiceRef.current.getRegimes({ lazyEvent: JSON.stringify({ first: 0, rows: 200, sortField: 'nome', sortOrder: 1 }) }),
            ]);

            setNaturezaOptions(naturezaRes.data.tiposEmpresa || []);
            setPorteOptions(porteRes.data.items || []);
            setRegimeOptions(regimeRes.data.regimes || []);
        } catch (err) {
            console.error('Erro ao carregar dropdowns:', err);
        }
    };

    useEffect(() => {
        loadData();
        loadDropdowns();
    }, []);

    useEffect(() => {
        loadData();
    }, [first, rows, sortField, sortOrder]);

    const onPage = (event: { first: number; rows: number; page: number }) => {
        setFirst(event.first);
        setRows(event.rows);
    };

    const onSort = (event: { sortField: string; sortOrder: number }) => {
        setSortField(event.sortField);
        setSortOrder(event.sortOrder);
    };

    const openNew = () => {
        setDialogMode('create');
        setFormData({
            id: '',
            nome: '',
            natureza_juridica_id: null,
            enquadramento_porte_id: null,
            regime_tributario_id: null,
            aliquota_base: 0,
            possui_fator_r: false,
            aliquota_fator_r: 0,
            ativo: true,
        });
        setDialogVisible(true);
    };

    const openEdit = async (item: MatrizItem) => {
        const natureza = naturezaOptions.find((o) => o.id === item.natureza_juridica_id) || null;
        const porte = porteOptions.find((o) => o.id === item.enquadramento_porte_id) || null;
        const regime = regimeOptions.find((o) => o.id === item.regime_tributario_id) || null;

        setDialogMode('edit');
        setFormData({
            id: item.id,
            nome: item.nome,
            natureza_juridica_id: natureza,
            enquadramento_porte_id: porte,
            regime_tributario_id: regime,
            aliquota_base: item.aliquota_base,
            possui_fator_r: item.possui_fator_r,
            aliquota_fator_r: item.aliquota_fator_r,
            ativo: item.ativo,
        });
        setDialogVisible(true);
    };

    const confirmDelete = (item: MatrizItem) => {
        confirmDialog({
            message: `Deseja realmente excluir a configuracao "${item.nome}"?`,
            header: 'Confirmar Exclusao',
            icon: 'pi pi-exclamation-triangle',
            acceptLabel: 'Sim',
            rejectLabel: 'Nao',
            accept: async () => {
                try {
                    await serviceRef.current.remove(item.id);
                    toast.current?.show({
                        severity: 'success',
                        summary: 'Sucesso',
                        detail: 'Configuracao excluida com sucesso',
                        life: 3000,
                    });
                    loadData();
                } catch (err) {
                    toast.current?.show({
                        severity: 'error',
                        summary: 'Erro',
                        detail: 'Erro ao excluir configuracao',
                        life: 3000,
                    });
                }
            },
        });
    };

    const save = async () => {
        if (!formData.nome || !formData.natureza_juridica_id || !formData.enquadramento_porte_id || !formData.regime_tributario_id) {
            toast.current?.show({
                severity: 'warn',
                summary: 'Atencao',
                detail: 'Preencha todos os campos obrigatorios',
                life: 3000,
            });
            return;
        }

        setSubmitting(true);
        try {
            const params = {
                nome: formData.nome,
                natureza_juridica_id: formData.natureza_juridica_id.id,
                enquadramento_porte_id: formData.enquadramento_porte_id.id,
                regime_tributario_id: formData.regime_tributario_id.id,
                aliquota_base: formData.aliquota_base,
                possui_fator_r: formData.possui_fator_r,
                aliquota_fator_r: formData.aliquota_fator_r,
                ativo: formData.ativo,
            };

            if (dialogMode === 'create') {
                await serviceRef.current.create(params);
                toast.current?.show({
                    severity: 'success',
                    summary: 'Sucesso',
                    detail: 'Configuracao criada com sucesso',
                    life: 3000,
                });
            } else {
                await serviceRef.current.update({ id: formData.id, ...params });
                toast.current?.show({
                    severity: 'success',
                    summary: 'Sucesso',
                    detail: 'Configuracao atualizada com sucesso',
                    life: 3000,
                });
            }

            setDialogVisible(false);
            loadData();
        } catch (err: any) {
            const msg = err?.response?.data?.error || err?.message || 'Erro ao salvar configuracao';
            toast.current?.show({
                severity: 'error',
                summary: 'Erro',
                detail: msg,
                life: 5000,
            });
        } finally {
            setSubmitting(false);
        }
    };

    const filterByNome = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setNomeFilter(value);
        setFirst(0);
    };

    const applyFilter = () => {
        loadData();
    };

    const actionBodyTemplate = (rowData: MatrizItem) => {
        return (
            <div className="flex gap-2">
                <Button
                    icon="pi pi-pencil"
                    rounded
                    severity="success"
                    className="mr-2"
                    tooltip="Alterar"
                    tooltipOptions={{ position: 'top' }}
                    onClick={() => openEdit(rowData)}
                />
                <Button
                    icon="pi pi-trash"
                    rounded
                    severity="warning"
                    tooltip="Excluir"
                    tooltipOptions={{ position: 'top' }}
                    onClick={() => confirmDelete(rowData)}
                />
            </div>
        );
    };

    const aliquotaBaseTemplate = (rowData: MatrizItem) => (
        <span>{formatAliquota(rowData.aliquota_base)}</span>
    );

    const aliquotaFatorRTemplate = (rowData: MatrizItem) => {
        if (!rowData.possui_fator_r) return <span className="text-gray-400">-</span>;
        return <span>{formatAliquota(rowData.aliquota_fator_r)}</span>;
    };

    const fatorRTemplate = (rowData: MatrizItem) => (
        <i
            className={`pi ${rowData.possui_fator_r ? 'pi-check text-green-500' : 'pi-times text-red-400'}`}
        />
    );

    const ativoTemplate = (rowData: MatrizItem) => (
        <i
            className={`pi ${rowData.ativo ? 'pi-check text-green-500' : 'pi-times text-red-400'}`}
        />
    );

    const paginatorLeft = (
        <Button
            type="button"
            icon="pi pi-refresh"
            className="p-button-text"
            tooltip="Atualizar"
            tooltipOptions={{ position: 'top' }}
            onClick={loadData}
        />
    );

    return (
        <div className="p-4">
            <Toast ref={toast} />
            <ConfirmDialog />

            <div className="flex justify-content-between align-items-center mb-4">
                <h1 className="text-2xl font-bold m-0">Matriz de Configuracao Tributaria</h1>
                <Button
                    label="Nova Configuracao"
                    icon="pi pi-plus"
                    severity="success"
                    onClick={openNew}
                />
            </div>

            <div className="flex align-items-center gap-2 mb-3">
                <span className="p-input-icon-left">
                    <i className="pi pi-search" />
                    <InputText
                        value={nomeFilter}
                        onChange={filterByNome}
                        placeholder="Buscar por nome..."
                        onKeyDown={(e) => e.key === 'Enter' && applyFilter()}
                    />
                </span>
                <Button icon="pi pi-search" severity="info" onClick={applyFilter} />
            </div>

            <DataTable
                value={items}
                lazy
                paginator
                rows={rows}
                first={first}
                totalRecords={totalRecords}
                onPage={onPage}
                onSort={onSort}
                sortField={sortField}
                sortOrder={sortOrder}
                loading={loading}
                dataKey="id"
                emptyMessage="Nenhuma configuracao encontrada."
                className="p-datatable-sm"
                paginatorLeft={paginatorLeft}
                paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport"
                currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords}"
            >
                <Column field="nome" header="Nome" sortable style={{ minWidth: '200px' }} />
                <Column field="natureza_juridica" header="Natureza Juridica" sortable style={{ minWidth: '180px' }} />
                <Column field="enquadramento_porte" header="Porte" sortable style={{ minWidth: '150px' }} />
                <Column field="regime_tributario" header="Regime Tributario" sortable style={{ minWidth: '180px' }} />
                <Column field="aliquota_base" header="Aliquota Base" body={aliquotaBaseTemplate} sortable style={{ minWidth: '120px' }} />
                <Column field="possui_fator_r" header="Fator R" body={fatorRTemplate} style={{ minWidth: '90px', textAlign: 'center' }} />
                <Column field="aliquota_fator_r" header="Aliquota Fator R" body={aliquotaFatorRTemplate} style={{ minWidth: '140px' }} />
                <Column field="ativo" header="Ativo" body={ativoTemplate} style={{ minWidth: '80px', textAlign: 'center' }} />
                <Column header="Acoes" body={actionBodyTemplate} style={{ minWidth: '120px' }} />
            </DataTable>

            <Dialog
                header={dialogMode === 'create' ? 'Nova Configuracao Tributaria' : 'Editar Configuracao Tributaria'}
                visible={dialogVisible}
                style={{ width: '550px' }}
                modal
                className="p-fluid"
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
                            onClick={save}
                            loading={submitting}
                        />
                    </div>
                }
            >
                <div className="field">
                    <label htmlFor="nome">Nome da Configuracao</label>
                    <InputText
                        id="nome"
                        value={formData.nome}
                        onChange={(e) => setFormData((prev) => ({ ...prev, nome: e.target.value }))}
                        placeholder="Ex: Simples Nacional - Anexo III - Servicos"
                        className="w-full"
                    />
                </div>

                <div className="field">
                    <label htmlFor="natureza_juridica">Natureza Juridica</label>
                    <Dropdown
                        id="natureza_juridica"
                        value={formData.natureza_juridica_id}
                        onChange={(e) => setFormData((prev) => ({ ...prev, natureza_juridica_id: e.value }))}
                        options={naturezaOptions}
                        optionLabel="descricao"
                        placeholder="Selecione a natureza juridica"
                        filter
                        showClear
                        className="w-full"
                    />
                </div>

                <div className="field">
                    <label htmlFor="porte">Enquadramento / Porte</label>
                    <Dropdown
                        id="porte"
                        value={formData.enquadramento_porte_id}
                        onChange={(e) => setFormData((prev) => ({ ...prev, enquadramento_porte_id: e.value }))}
                        options={porteOptions}
                        optionLabel="sigla"
                        placeholder="Selecione o porte"
                        filter
                        showClear
                        className="w-full"
                        itemTemplate={(option: DropdownOption) => (
                            <span>{option.sigla}{option.descricao ? ` - ${option.descricao}` : ''}</span>
                        )}
                        valueTemplate={(option: DropdownOption) => {
                            if (!option) return <span>Selecione o porte</span>;
                            return <span>{option.sigla}{option.descricao ? ` - ${option.descricao}` : ''}</span>;
                        }}
                    />
                </div>

                <div className="field">
                    <label htmlFor="regime">Regime Tributario</label>
                    <Dropdown
                        id="regime"
                        value={formData.regime_tributario_id}
                        onChange={(e) => setFormData((prev) => ({ ...prev, regime_tributario_id: e.value }))}
                        options={regimeOptions}
                        optionLabel="nome"
                        placeholder="Selecione o regime tributario"
                        filter
                        showClear
                        className="w-full"
                    />
                </div>

                <div className="field">
                    <label htmlFor="aliquota_base">Aliquota Base (%)</label>
                    <InputNumber
                        id="aliquota_base"
                        value={formData.aliquota_base * 100}
                        onValueChange={(e) => setFormData((prev) => ({ ...prev, aliquota_base: (e.value || 0) / 100 }))}
                        suffix="%"
                        minFractionDigits={2}
                        maxFractionDigits={4}
                        min={0}
                        max={100}
                        className="w-full"
                    />
                </div>

                <div className="field-checkbox flex align-items-center gap-2">
                    <input
                        type="checkbox"
                        id="possui_fator_r"
                        checked={formData.possui_fator_r}
                        onChange={(e) => {
                            const checked = e.target.checked;
                            setFormData((prev) => ({
                                ...prev,
                                possui_fator_r: checked,
                                aliquota_fator_r: checked ? prev.aliquota_fator_r : 0,
                            }));
                        }}
                    />
                    <label htmlFor="possui_fator_r">Possui Fator R?</label>
                </div>

                {formData.possui_fator_r && (
                    <div className="field">
                        <label htmlFor="aliquota_fator_r">Aliquota Alternativa (Fator R) (%)</label>
                        <InputNumber
                            id="aliquota_fator_r"
                            value={formData.aliquota_fator_r * 100}
                            onValueChange={(e) => setFormData((prev) => ({ ...prev, aliquota_fator_r: (e.value || 0) / 100 }))}
                            suffix="%"
                            minFractionDigits={2}
                            maxFractionDigits={4}
                            min={0}
                            max={100}
                            className="w-full"
                        />
                        <small className="text-gray-500">
                            Usada quando a folha de pagamento for &ge; 28% do faturamento
                        </small>
                    </div>
                )}

                <div className="field-checkbox flex align-items-center gap-2 mt-3">
                    <input
                        type="checkbox"
                        id="ativo"
                        checked={formData.ativo}
                        onChange={(e) => setFormData((prev) => ({ ...prev, ativo: e.target.checked }))}
                    />
                    <label htmlFor="ativo">Ativo</label>
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
                onClick={loadData}
            />
        </div>
    );
};

export default MatrizConfiguracaoTributaria;
