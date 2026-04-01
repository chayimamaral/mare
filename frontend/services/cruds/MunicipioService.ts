import setupAPIClient from '../../components/api/api';
//import { setupAPIClient } from '../api';

function mutationList(data: { municipios?: unknown[]; cidades?: unknown[]; totalRecords?: number }) {
    const municipios = data.municipios ?? data.cidades ?? [];
    const totalRecords = data.totalRecords ?? 0;
    return { municipios, totalRecords };
}

export default function MunicipioService() {

    const getMunicipios = async (params) => {

        // do jeito que receber o params, ele vem como string, entao tem que converter para objeto

        const apiClient = setupAPIClient(undefined);
        // ao passar o params, ele vem como string, entao tem que converter para objeto com JSON.parse
        const response = await apiClient.get('/api/cidades', {
            params: JSON.parse(params.lazyEvent),
        });

        // não precisa converter para objeto, pois o axios já faz isso
        const { municipios, totalRecords } = response.data

        return {
            data: {
                municipios,
                totalRecords
            }
        }
    }

    const createMunicipio = async (params) => {
        const apiClient = setupAPIClient(undefined);

        const response = await apiClient.post('/api/cidade', {
            params: {
                ...params
            }
        });

        return {
            data: mutationList(response.data)
        };
    }

    const updateMunicipio = async (params) => {
        const apiClient = setupAPIClient(undefined);

        const response = await apiClient.put('/api/cidade', {
            params: {
                ...params
            }

        });

        return {
            data: mutationList(response.data)
        };
    }

    const deleteMunicipio = async (params) => {
        const apiClient = setupAPIClient(undefined);
        const response = await apiClient.delete('/api/cidade', {
            params: {
                ...params
            }

        });

        return {
            data: mutationList(response.data)
        };
    }

    const getMunicipiosLite = async () => {

        // do jeito que receber o params, ele vem como string, entao tem que converter para objeto

        const apiClient = setupAPIClient(undefined);
        // ao passar o params, ele vem como string, entao tem que converter para objeto com JSON.parse
        const response = await apiClient.get('/api/cidadeslite');

        // não precisa converter para objeto, pois o axios já faz isso
        const { municipios, totalRecords } = response.data

        return {
            data: {
                municipios
            }
        }
    }

    return {
        getMunicipios,
        createMunicipio,
        updateMunicipio,
        deleteMunicipio,
        getMunicipiosLite
    }
}
//export default MunicipioService;

