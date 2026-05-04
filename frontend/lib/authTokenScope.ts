/**
 * Chave estável e distinta por JWT para queryKey / remount.
 * Não usar prefixo curto do token: o header Base64 de HS256 costuma ser idêntico entre usuários.
 */
export function authTokenScopeKey(token: string): string {
    const t = String(token ?? '').trim();
    if (!t) {
        return 'anon';
    }
    let h = 2166136261;
    for (let i = 0; i < t.length; i++) {
        h ^= t.charCodeAt(i);
        h = Math.imul(h, 16777619);
    }
    return `tk:${(h >>> 0).toString(16)}:${t.length}`;
}
