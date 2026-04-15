import setupAPIClient from '../../components/api/api';

export default function SalarioMinimoService() {
  const list = async () => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.get('/api/salarios-minimos');
    const { salarios } = response.data;
    return { data: { salarios: salarios ?? [] } };
  };

  const create = async (params: { ano: number; valor: number }) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.post('/api/salario-minimo', { params });
    return { data: response.data };
  };

  const update = async (params: { id: string; ano: number; valor: number }) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.put('/api/salario-minimo', { params });
    return { data: response.data };
  };

  const remove = async (id: string) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.delete('/api/salario-minimo', { params: { id } });
    return { data: response.data };
  };

  return { list, create, update, remove };
}
