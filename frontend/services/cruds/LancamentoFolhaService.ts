import setupAPIClient from '../../components/api/api';

export default function LancamentoFolhaService() {

    const getTree = async () => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.get('/api/lancamentos-folha');
        const { tree } = response.data;
        return { data: { tree } };
    };

    const createLancamento = async (params: {
        cliente_id: string;
        competencia: string;
        valor_folha: number;
        valor_faturamento: number;
        observacoes: string;
    }) => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.post('/api/lancamento-folha', {
            params: { ...params },
        });
        return { data: response.data };
    };

    const updateLancamento = async (params: {
        id: string;
        competencia: string;
        valor_folha: number;
        valor_faturamento: number;
        observacoes: string;
    }) => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.put('/api/lancamento-folha', {
            params: { ...params },
        });
        return { data: response.data };
    };

    const deleteLancamento = async (id: string) => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.delete('/api/lancamento-folha', {
            params: { id },
        });
        return { data: response.data };
    };

    const getLancamento = async (id: string) => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.get('/api/lancamento-folha', {
            params: { id },
        });
        return { data: response.data };
    };

    return {
        getTree,
        createLancamento,
        updateLancamento,
        deleteLancamento,
        getLancamento,
    };
}
