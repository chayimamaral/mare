import setupAPIClient from '../../components/api/api';

export type CatalogoServico = {
    id: string;
    secao: string;
    sequencial: number;
    codigo: string;
    id_sistema: string;
    id_servico: string;
    data_implantacao?: string;
    tipo: string;
    descricao: string;
    ativo?: boolean;
};

export default function CatalogoServicoService() {
    const list = async (opts?: { secao?: string; incluirInativos?: boolean }) => {
        const api = setupAPIClient(undefined);
        const { data } = await api.get('/api/catalogo-servicos', {
            params: {
                secao: opts?.secao ?? '',
                incluir_inativos: opts?.incluirInativos ? '1' : '',
            },
        });
        return data?.servicos ?? [];
    };

    const create = async (params: Omit<CatalogoServico, 'id'>) => {
        const api = setupAPIClient(undefined);
        return api.post('/api/catalogo-servico', { params });
    };

    const update = async (params: CatalogoServico) => {
        const api = setupAPIClient(undefined);
        return api.put('/api/catalogo-servico', { params });
    };

    const remove = async (id: string) => {
        const api = setupAPIClient(undefined);
        return api.put('/api/deletecatalogo-servico', { params: { id } });
    };

    return { list, create, update, remove };
}
