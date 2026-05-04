import { useEffect, useState } from 'react';
import { authTokenScopeKey } from '../../lib/authTokenScope';
import { AUTH_SESSION_CHANGED_EVENT } from '../context/AuthContext';

function readToken(): string {
    if (typeof window === 'undefined') {
        return '';
    }
    const fromSession = String(window.sessionStorage.getItem('vecontab_token') ?? '').trim();
    if (fromSession) {
        return fromSession;
    }
    return String(window.localStorage.getItem('vecontab_token') ?? '').trim();
}

export function useAuthScopeKey(): string {
    const [scope, setScope] = useState(() => authTokenScopeKey(readToken()));

    useEffect(() => {
        const refresh = () => {
            setScope(authTokenScopeKey(readToken()));
        };
        refresh();
        if (typeof window !== 'undefined') {
            window.addEventListener(AUTH_SESSION_CHANGED_EVENT, refresh);
        }
        return () => {
            if (typeof window !== 'undefined') {
                window.removeEventListener(AUTH_SESSION_CHANGED_EVENT, refresh);
            }
        };
    }, []);

    return scope;
}

