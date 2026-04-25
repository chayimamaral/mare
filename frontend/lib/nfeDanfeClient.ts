import api from '../components/api/apiClient';

export function parseDanfeErrorMessage(e: unknown): string {
    const err = e as { response?: { data?: unknown }; message?: string };
    let msg =
        (err?.response?.data as { error?: string; message?: string })?.error
        ?? (err?.response?.data as { message?: string })?.message
        ?? err?.message
        ?? 'Falha ao gerar DANFE';
    const raw = err?.response?.data;
    if (typeof raw === 'string' && raw.trim().startsWith('{')) {
        try {
            const o = JSON.parse(raw) as { error?: string };
            if (o?.error) {
                msg = o.error;
            }
        } catch {
            /* ignore */
        }
    }
    return String(msg);
}

export async function fetchDanfeHtmlFromRetorno(retorno: string): Promise<string> {
    const { data } = await api.post<string>(
        '/api/serpro/nfe/documento/danfe-html',
        { retorno },
        { responseType: 'text' as any },
    );
    const html = String(data ?? '').trim();
    if (!html) {
        throw new Error('DANFE_VAZIO');
    }
    return html;
}
