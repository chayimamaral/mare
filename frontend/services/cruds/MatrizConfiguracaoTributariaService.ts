import setupAPIClient from '../../components/api/api';

export default function MatrizConfiguracaoTributariaService() {

    const list = async (params: { lazyEvent: string }) => {
        const apiClient = setupAPIClient(undefined);
        const lazy = typeof params.lazyEvent === 'string' ? JSON.parse(params.lazyEvent) : params.lazyEvent;
        const response = await apiClient.get('/api/matriz-configuracao-tributaria', {
            params: {
                first: lazy.first ?? 0,
                rows: lazy.rows ?? 25,
                sortField: lazy.sortField ?? 'nome',
                sortOrder: lazy.sortOrder ?? 1,
                filters: typeof lazy.filters === 'string' ? lazy.filters : JSON.stringify(lazy.filters ?? {}),
            },
        });
        const { items, totalRecords } = response.data;
        return { data: { items, totalRecords } };
    };

    const create = async (params: Record<string, unknown>) => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.post('/api/matriz-configuracao-tributaria', { params });
        return { data: response.data };
    };

    const update = async (params: Record<string, unknown>) => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.put('/api/matriz-configuracao-tributaria', { params });
        return { data: response.data };
    };

    const remove = async (id: string) => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.delete('/api/matriz-configuracao-tributaria', {
            params: { id },
        });
        return { data: response.data };
    };

    const getByID = async (id: string) => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.get('/api/matriz-configuracao-tributaria/item', {
            params: { id },
        });
        return { data: response.data };
    };

    return { list, create, update, remove, getByID };
}
