import setupAPIClient from '../../components/api/api';

export default function SerproServicoEnquadramentoService() {
  const list = async (enquadramento_id: string, regime_tributario_id: string) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.get('/api/serpro-servico-enquadramento', {
      params: { enquadramento_id, regime_tributario_id },
    });
    return { data: { servicos_ids: response.data?.servicos_ids ?? [] } };
  };

  const save = async (params: { enquadramento_id: string; regime_tributario_id: string; servicos_ids: string[] }) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.put('/api/serpro-servico-enquadramento', { params });
    return { data: response.data };
  };

  return { list, save };
}
