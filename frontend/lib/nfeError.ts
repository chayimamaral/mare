export type NFEApiErrorPayload = {
    error?: string;
    message?: string;
    codigo?: string;
    mensagem?: string;
    etapa_validacao?: string;
    acao_sugerida?: string;
    origem?: string;
};

export function parseNFEApiError(error: unknown): { title: string; detail: string } {
    const payload = (
        typeof error === 'object' &&
        error &&
        'response' in error &&
        typeof (error as { response?: unknown }).response === 'object'
            ? ((error as { response?: { data?: NFEApiErrorPayload } }).response?.data ?? {})
            : {}
    ) as NFEApiErrorPayload;

    const base = payload.mensagem || payload.error || payload.message || 'Falha ao processar NF-e';
    const parts = [base];
    if (payload.codigo) {
        parts.push(`Código: ${payload.codigo}`);
    }
    if (payload.etapa_validacao) {
        parts.push(`Etapa: ${payload.etapa_validacao}`);
    }
    if (payload.acao_sugerida) {
        parts.push(`Ação: ${payload.acao_sugerida}`);
    }
    return {
        title: payload.codigo ? `Erro ${payload.codigo}` : 'Erro',
        detail: parts.join(' | '),
    };
}
