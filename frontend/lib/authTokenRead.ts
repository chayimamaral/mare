import { parseCookies } from 'nookies';
import { getAuthTokenFromParsedCookies } from '../constants/authCookie';

/**
 * Lê o token no browser na mesma ordem do AuthContext (sessão → cookie → localStorage).
 * No SSR retorna string vazia.
 */
export function readAuthTokenForGuard(): string {
    if (typeof window === 'undefined') {
        return '';
    }
    const cookieToken = getAuthTokenFromParsedCookies(parseCookies()) ?? '';
    const sessionToken = String(window.sessionStorage.getItem('vecontab_token') ?? '').trim();
    const localToken = String(window.localStorage.getItem('vecontab_token') ?? '').trim();
    return String(sessionToken || cookieToken || localToken || '').trim();
}
