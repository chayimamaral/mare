import setupAPIClient from '../../components/api/api';
import { Vec } from '../../types/types';

export interface MonitorOperacaoListResponse {
  itens: Vec.MonitorOperacaoItem[];
  total: number;
}

export default function MonitorOperacaoService() {
  const list = async (limit = 50, offset = 0) => {
    const api = setupAPIClient(undefined);
    const response = await api.get<MonitorOperacaoListResponse>('/api/monitor/operacoes', {
      params: { limit, offset },
    });
    return response.data ?? { itens: [], total: 0 };
  };

  return { list };
}
