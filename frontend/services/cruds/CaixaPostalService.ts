import setupAPIClient from '../../components/api/api';

export interface CaixaPostalMensagem {
    id: string;
    remetente_id?: string;
    remetente_tenantid?: string;
    remetente_nome: string;
    tipo: 'INBOX' | 'OUTBOX';
    is_global: boolean;
    titulo: string;
    conteudo: string;
    lida: boolean;
    lida_por?: string;
    lida_em?: string;
    criado_em: string;
}

export interface EnviarPayload {
    tenant_id?: string;
    is_global?: boolean;
    titulo: string;
    conteudo: string;
}

export default function CaixaPostalService() {
    const listar = async (): Promise<CaixaPostalMensagem[]> => {
        const api = setupAPIClient(undefined);
        const { data } = await api.get('/api/caixa-postal');
        return data ?? [];
    };

    const contarNaoLidas = async (): Promise<number> => {
        const api = setupAPIClient(undefined);
        const { data } = await api.get('/api/caixa-postal/count-nao-lidas');
        return data?.count ?? 0;
    };

    const marcarComoLida = async (id: string): Promise<void> => {
        const api = setupAPIClient(undefined);
        await api.put(`/api/caixa-postal/${id}/ler`);
    };

    const enviar = async (payload: EnviarPayload): Promise<void> => {
        const api = setupAPIClient(undefined);
        await api.post('/api/caixa-postal/enviar', payload);
    };

    return { listar, contarNaoLidas, marcarComoLida, enviar };
}
