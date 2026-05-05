import React, { useEffect, useRef, useState } from 'react';
import { InputText } from 'primereact/inputtext';
import { Button } from 'primereact/button';

import { Toast } from 'primereact/toast';
import { TabPanel, TabView } from 'primereact/tabview';

import { InputTextarea } from 'primereact/inputtextarea';
import RegistroService from '../../services/cruds/RegistroService';
import { isValidCNPJ, onlyDigits } from '../../constants/documento';

interface Registro {
    tenantid: string,
    razaosocial: string,
    fantasia: string,
    endereco: string,
    bairro: string,
    cidade: string,
    estado: string,
    cep: string,
    telefone: string,
    email: string,
    cnpj: string,
    ie: string,
    im: string,
    observacoes: string,
    enviar_resumo_mensal: boolean
}

function Registro() {

    const [dropdownItem, setDropdownItem] = useState(null);
    const toast = useRef<Toast>(null);
    const [isInvalid, setIsInvalid] = useState(false);

    const [registro, setRegistro] = useState<Registro>({
        tenantid: '',
        razaosocial: '',
        fantasia: '',
        endereco: '',
        bairro: '',
        cidade: '',
        estado: '',
        cep: '',
        telefone: '',
        email: '',
        cnpj: '',
        ie: '',
        im: '',
        observacoes: '',
        enviar_resumo_mensal: false
    });
    const [activeTabIndex, setActiveTabIndex] = useState(0);

    const registroService = RegistroService();

    async function handleUpdateEmpresa() {
        const cnpjDigits = onlyDigits(registro.cnpj);
        if (cnpjDigits.length === 0) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Informe o CNPJ.', life: 3500 });
            return;
        }
        if (!isValidCNPJ(cnpjDigits)) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'CNPJ inválido.', life: 3500 });
            return;
        }
        try {
            await registroService.gravaRegistro(registro);
            toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Dados atualizados com Sucesso', life: 3000 });
        } catch (error: any) {
            toast.current?.show({
                severity: 'error',
                summary: 'Erro',
                detail: error?.response?.data?.message || 'Falha ao gravar os dados',
                life: 3500,
            });
        }
    }

    useEffect(() => {
        loadLazyRegistro()
    }, [])

    const toSafeString = (value: unknown): string => {
        if (value == null) return '';
        if (typeof value === 'string') return value;
        if (typeof value === 'number' || typeof value === 'boolean') return String(value);
        if (typeof value === 'object') {
            const maybeNullString = value as { String?: unknown; Valid?: unknown; string?: unknown; valid?: unknown };
            if (typeof maybeNullString.String === 'string') {
                const isValid = typeof maybeNullString.Valid === 'boolean' ? maybeNullString.Valid : true;
                return isValid ? maybeNullString.String : '';
            }
            if (typeof maybeNullString.string === 'string') {
                const isValid = typeof maybeNullString.valid === 'boolean' ? maybeNullString.valid : true;
                return isValid ? maybeNullString.string : '';
            }
        }
        return '';
    };

    const normalizeRegistro = (raw: any): Registro => ({
        tenantid: toSafeString(raw?.tenantid),
        razaosocial: toSafeString(raw?.razaosocial),
        fantasia: toSafeString(raw?.fantasia),
        endereco: toSafeString(raw?.endereco),
        bairro: toSafeString(raw?.bairro),
        cidade: toSafeString(raw?.cidade),
        estado: toSafeString(raw?.estado),
        cep: toSafeString(raw?.cep),
        telefone: toSafeString(raw?.telefone),
        email: toSafeString(raw?.email),
        cnpj: toSafeString(raw?.cnpj),
        ie: toSafeString(raw?.ie),
        im: toSafeString(raw?.im),
        observacoes: toSafeString(raw?.observacoes),
        enviar_resumo_mensal: Boolean(raw?.enviar_resumo_mensal),
    });

    const loadLazyRegistro = async () => {
        try {
            const { dados } = await registroService.getRegistro(registro)
            setRegistro(normalizeRegistro(dados));
        } catch (error) {
            toast.current?.show({ severity: 'error', summary: 'Erro', detail: 'Erro ao carregar os dados', life: 3000 });
        }
    }

    const onInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>, nome: keyof Registro) => {

        const val = (e.target && e.target.value) || '';
        let _registro = { ...registro };
        _registro[nome] = val;

        setRegistro(_registro);

    }

    const onToggleResumoMensal = (e: React.ChangeEvent<HTMLInputElement>) => {
        setRegistro((prev) => ({ ...prev, enviar_resumo_mensal: e.target.checked }));
    };

    return (
        <div className="grid">
            <div className="col-12">
                <div className="card">
                    <TabView
                        activeIndex={activeTabIndex}
                        onTabChange={(e) => setActiveTabIndex(e.index)}
                        className="w-full"
                        pt={{
                            navContainer: {
                                className: 'w-full',
                                style: {
                                    boxSizing: 'border-box',
                                    borderBottom: 'none',
                                    paddingLeft: 0,
                                    paddingRight: 0,
                                    marginBottom: '0.5rem',
                                },
                            },
                            navContent: {
                                style: {
                                    flex: '1 1 auto',
                                    minWidth: 0,
                                    overflowX: 'auto',
                                    overflowY: 'hidden',
                                },
                            },
                            inkbar: { style: { display: 'none' } },
                            nav: {
                                style: {
                                    display: 'flex',
                                    flexWrap: 'nowrap',
                                    width: '100%',
                                    justifyContent: 'flex-start',
                                    alignItems: 'flex-end',
                                    columnGap: '1rem',
                                    listStyle: 'none',
                                    margin: 0,
                                    paddingLeft: 0,
                                    paddingRight: 0,
                                },
                            },
                            panelContainer: {
                                style: {
                                    marginTop: '0.25rem',
                                },
                            },
                        }}
                    >
                        <TabPanel header="Dados da Empresa" headerStyle={{ whiteSpace: 'nowrap' }}>
                            <div className="p-fluid formgrid grid">
                                <div className="field col-12 md:col-6">
                                    <label htmlFor="razaosocial_">Razão Social</label>
                                    <InputText value={registro.razaosocial ?? ''} onChange={(e) => onInputChange(e, 'razaosocial')} id="razaosocial_" type="text" />
                                </div>
                                <div className="field col-12 md:col-6">
                                    <label htmlFor="fantasia_">Nome Fantasia</label>
                                    <InputText value={registro.fantasia ?? ''} onChange={(e) => onInputChange(e, 'fantasia')} id="fantasia_" type="text" />
                                </div>
                                <div className="field col-12 md:col-6">
                                    <label htmlFor="endereco_">Endereço</label>
                                    <InputText value={registro.endereco ?? ''} onChange={(e) => onInputChange(e, 'endereco')} id="endereco_" type="text" />
                                </div>
                                <div className="field col-12 md:col-6">
                                    <label htmlFor="_bairro_">Bairro</label>
                                    <InputText value={registro.bairro ?? ''} onChange={(e) => onInputChange(e, 'bairro')} id="bairro_" type="text" />
                                </div>
                                <div className="field col-12 md:col-6">
                                    <label htmlFor="cidade_">Cidade</label>
                                    <InputText value={registro.cidade ?? ''} onChange={(e) => onInputChange(e, 'cidade')} id="cidade_" type="text" />
                                </div>
                                <div className="field col-12 md:col-3">
                                    <label htmlFor="estado_">Estado</label>
                                    <InputText value={registro.estado ?? ''} onChange={(e) => onInputChange(e, 'estado')} id="estado_" type='text' />
                                </div>
                                <div className="field col-12 md:col-3">
                                    <label htmlFor="cep_">CEP</label>
                                    <InputText value={registro.cep ?? ''} onChange={(e) => onInputChange(e, 'cep')} id="cep_" type="text" />
                                </div>
                                <div className="field col-12 md:col-6">
                                    <label htmlFor="telefone_">Telefone</label>
                                    <InputText value={registro.telefone ?? ''} onChange={(e) => onInputChange(e, 'telefone')} id="telefone" type="text" />
                                </div>
                                <div className="field col-12 md:col-6">
                                    <label htmlFor="email_">Email</label>
                                    <InputText value={registro.email ?? ''} onChange={(e) => onInputChange(e, 'email')} id="email_" type="text" />
                                </div>

                                <div className="field col-12 md:col-4">
                                    <label htmlFor="cnpj_">CNPJ</label>
                                    <InputText
                                        value={registro.cnpj ?? ''}
                                        onChange={(e) => onInputChange(e, 'cnpj')}
                                        id="cnpj_"
                                        type="text"
                                        inputMode="numeric"
                                        maxLength={18}
                                        className={isInvalid || (onlyDigits(registro.cnpj ?? '').length > 0 && !isValidCNPJ(registro.cnpj ?? '')) ? 'p-invalid' : ''}
                                    />
                                </div>
                                <div className="field col-12 md:col-4">
                                    <label htmlFor="ie_">Inscrição Estadual</label>
                                    <InputText value={registro.ie ?? ''} onChange={(e) => onInputChange(e, 'ie')} id="ie_" type="text" />
                                </div>
                                <div className="field col-12 md:col-4">
                                    <label htmlFor="im_">Inscrição Municipal</label>
                                    <InputText value={registro.im ?? ''} onChange={(e) => onInputChange(e, 'im')} id="im_" type="text" />
                                </div>
                                <div className="field col-12">
                                    <label htmlFor="observacoes_">Observações</label>
                                    <InputTextarea name='observacoes' value={registro.observacoes ?? ''} onChange={(e) => onInputChange(e, 'observacoes')} id="observacoes_" rows={4} />
                                </div>
                            </div>
                        </TabPanel>
                        <TabPanel header="Configurações" headerStyle={{ whiteSpace: 'nowrap' }}>
                            <div className="p-fluid formgrid grid">
                                <div className="field col-12">
                                    <div className="field-checkbox mb-3">
                                        <input
                                            id="registro_enviar_resumo_mensal"
                                            type="checkbox"
                                            checked={Boolean(registro.enviar_resumo_mensal)}
                                            onChange={onToggleResumoMensal}
                                        />
                                        <label htmlFor="registro_enviar_resumo_mensal" className="ml-2">
                                            Enviar resumo mensal de notas fiscais recebidas para os clientes
                                        </label>
                                    </div>
                                </div>
                            </div>
                        </TabPanel>
                    </TabView>
                    <div className="mt-3">
                        <Toast ref={toast} />
                        <Button label="Gravar" onClick={handleUpdateEmpresa}></Button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Registro;

    // Aqui não é necessário nenhum processamento adicional
