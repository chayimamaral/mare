import setupAPIClient from '../../components/api/api';

const api = () => setupAPIClient(undefined);

export default function RotinaPFService() {
  const getRotinasPFLite = async () => {
    const response = await api().get('/api/listrotinapflite');
    return {
      data: {
        rotinas_pf: response.data?.rotinas_pf ?? [],
        totalRecords: response.data?.totalRecords ?? 0,
      },
    };
  };

  const getRotinasPFAdmin = async (params: { lazyEvent: string }) => {
    const p = JSON.parse(params.lazyEvent) as Record<string, unknown>;
    const q: Record<string, unknown> = { ...p };
    if (q.filters != null && typeof q.filters === 'object') {
      q.filters = JSON.stringify(q.filters);
    }
    const response = await api().get('/api/listrotinaspf', { params: q });
    return {
      data: {
        rotinas_pf: response.data?.rotinas_pf ?? [],
        totalRecords: response.data?.totalRecords ?? 0,
      },
    };
  };

  const getItens = async (rotinaPfId: string) => {
    const response = await api().get('/api/rotinapfitens', {
      params: { rotina_pf_id: rotinaPfId },
    });
    return {
      data: {
        itens: response.data?.itens ?? [],
        totalRecords: response.data?.totalRecords ?? 0,
      },
    };
  };

  const createRotinaPF = async (payload: {
    nome: string;
    categoria: string;
    descricao?: string;
    ativo?: boolean;
  }) => {
    const response = await api().post('/api/rotinapf', {
      params: {
        nome: payload.nome,
        categoria: payload.categoria,
        descricao: payload.descricao ?? '',
        ativo: payload.ativo !== false,
      },
    });
    return { data: response.data };
  };

  const updateRotinaPF = async (payload: {
    id: string;
    nome: string;
    categoria: string;
    descricao?: string;
    ativo?: boolean;
  }) => {
    const response = await api().put('/api/rotinapf', {
      params: {
        id: payload.id,
        nome: payload.nome,
        categoria: payload.categoria,
        descricao: payload.descricao ?? '',
        ativo: payload.ativo !== false,
      },
    });
    return { data: response.data };
  };

  const deleteRotinaPF = async (id: string) => {
    const response = await api().put('/api/deleterotinapf', {
      params: { id },
    });
    return { data: response.data };
  };

  const createItem = async (payload: {
    rotina_pf_id: string;
    ordem: number;
    passo_id?: string;
    descricao?: string;
    tempo_estimado?: number;
  }) => {
    const response = await api().post('/api/rotinapfitem', {
      params: {
        rotina_pf_id: payload.rotina_pf_id,
        ordem: payload.ordem,
        passo_id: payload.passo_id ?? '',
        descricao: payload.descricao ?? '',
        tempo_estimado: payload.tempo_estimado ?? 0,
      },
    });
    return { data: response.data };
  };

  const updateItem = async (payload: {
    item_id: string;
    rotina_pf_id: string;
    ordem: number;
    passo_id?: string;
    descricao?: string;
    tempo_estimado?: number;
  }) => {
    const response = await api().put('/api/rotinapfitem', {
      params: {
        item_id: payload.item_id,
        rotina_pf_id: payload.rotina_pf_id,
        ordem: payload.ordem,
        passo_id: payload.passo_id ?? '',
        descricao: payload.descricao ?? '',
        tempo_estimado: payload.tempo_estimado ?? 0,
      },
    });
    return { data: response.data };
  };

  const deleteItem = async (itemId: string, rotinaPfId: string) => {
    const response = await api().put('/api/deleterotinapfitem', {
      params: { item_id: itemId, rotina_pf_id: rotinaPfId },
    });
    return { data: response.data };
  };

  return {
    getRotinasPFLite,
    getRotinasPFAdmin,
    getItens,
    createRotinaPF,
    updateRotinaPF,
    deleteRotinaPF,
    createItem,
    updateItem,
    deleteItem,
  };
}
