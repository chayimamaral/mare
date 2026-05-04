import setupAPIClient from '../../components/api/api';
import type { Vec } from '../../types/types';

export default function EnquadramentoJuridicoPorteService() {
  const list = async (anoVigencia?: number) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.get<{ items: Vec.EnquadramentoJuridicoPorte[] }>('/api/enquadramentos-juridicos-porte', {
      params: anoVigencia != null && anoVigencia > 0 ? { ano_vigencia: anoVigencia } : {},
    });
    return { data: { items: response.data?.items ?? [] } };
  };

  const create = async (params: {
    sigla: string;
    descricao: string;
    limite_inicial: number;
    limite_final: number | null;
    ano_vigencia: number;
  }) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.post('/api/enquadramento-juridico-porte', { params });
    return { data: response.data };
  };

  const update = async (params: {
    id: string;
    sigla: string;
    descricao: string;
    limite_inicial: number;
    limite_final: number | null;
    ano_vigencia: number;
    ativo?: boolean;
  }) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.put('/api/enquadramento-juridico-porte', { params });
    return { data: response.data };
  };

  const remove = async (id: string) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.delete('/api/enquadramento-juridico-porte', { params: { id } });
    return { data: response.data };
  };

  return { list, create, update, remove };
}
