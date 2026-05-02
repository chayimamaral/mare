import { createContext, ReactNode, useEffect, useState } from 'react';
import { setCookie, parseCookies } from 'nookies';
import Router from 'next/router';
import { AxiosError } from 'axios';
import { useQueryClient } from '@tanstack/react-query';

import api from '../api/apiClient';
import { AuthError } from '../../lib/authErrors';
import {
  AUTH_TOKEN_COOKIE,
  clearAuthTokenCookies,
  clearLegacyAuthTokenCookieBrowser,
  getAuthTokenFromParsedCookies,
} from '../../constants/authCookie';

interface AuthContextData {
  user?: UserProps | undefined;
  isAuthenticated: boolean;
  /** true após decidir sessão inicial (/api/me ou ausência de token). Usado para não pintar o dashboard antes da validação. */
  authBootstrapped: boolean;
  signIn: (credentials: SignInProps) => Promise<void>;
  signUp: (credentials: SignUpProps) => Promise<SignUpResult>;
  logoutUser: () => Promise<void>;
}

interface UserProps {
  id: string;
  nome: string;
  email: string;
  tenant?: Tenant | null;
}

interface Tenant {
  id: string;
  nome?: string;
  schema_name?: string;
  schemaName?: string;
}

interface SubscriptionProps {
  id: string;
  status: string;
}

type AuthProviderProps = {
  children: ReactNode;
}

interface SignInProps {
  email: string;
  password: string;
}

interface SignUpProps {
  nome: string;
  email: string;
  password: string;
  empresa_nome: string;

}

interface SignUpResult {
  id: string;
  nome: string;
  email: string;
  role: string;
  tenantid: string;
  tenant_schema?: string;
  active: boolean;
}

export const AUTH_SESSION_CHANGED_EVENT = 'vecx:auth-session-changed';

const AuthContext = createContext({} as AuthContextData);

export function signOut() {

  try {
    clearAuthTokenCookies(null);
    if (typeof window !== 'undefined') {
      window.sessionStorage.removeItem('vecontab_token');
      window.localStorage.removeItem('vecontab_token');
      window.dispatchEvent(new Event(AUTH_SESSION_CHANGED_EVENT));
    }
    Router.replace('/auth/login');

  } catch (err) {

  }
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<UserProps>()
  const [authBootstrapped, setAuthBootstrapped] = useState(false);
  const isAuthenticated = !!user;
  const queryClient = useQueryClient();

  useEffect(() => {
    const cookieToken = getAuthTokenFromParsedCookies(parseCookies());
    const sessionToken =
      typeof window !== 'undefined' ? String(window.sessionStorage.getItem('vecontab_token') ?? '').trim() : '';
    const token =
      sessionToken ||
      cookieToken ||
      (typeof window !== 'undefined' ? String(window.localStorage.getItem('vecontab_token') ?? '').trim() : '');

    if (!token) {
      setAuthBootstrapped(true);
      return;
    }

    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    api.get('/api/me').then(response => {
      const { id, nome, email, tenant } = response.data ?? {}
      setUser({ id, nome, email, tenant })

    })
      .catch((err) => {
        const axiosErr = err as AxiosError;
        if (axiosErr?.response?.status === 401) {
          signOut();
        }
      })
      .finally(() => {
        setAuthBootstrapped(true);
      });
  }, [])

  const signIn = async ({ email, password }: SignInProps) => {

    //async function signIn({ email, password }: SignInProps) {
    try {
      // Evita qualquer reaproveitamento de cache entre sessões/tenants diferentes.
      queryClient.clear();

      const response = await api.post("/api/session", {
        email,
        password,
      })

      const { id, nome, token, empresa } = response.data;


      setCookie(undefined, AUTH_TOKEN_COOKIE, token, {
        maxAge: 60 * 60 * 24 * 30, // Expirar em 1 mês
        path: '/',
        sameSite: 'lax',
      })
      clearLegacyAuthTokenCookieBrowser();
      try {
        window.sessionStorage.setItem('vecontab_token', token);
        window.localStorage.setItem('vecontab_token', token);
        window.dispatchEvent(new Event(AUTH_SESSION_CHANGED_EVENT));
      } catch {
        // ignore
      }

      api.defaults.headers.common['Authorization'] = `Bearer ${token}`

      // Define o usuário imediatamente com os dados do login para o topbar não ficar em branco
      setUser({ id, nome, email });

      try {
        // Enriquece com dados completos (tenant) passando o token explicitamente
        const me = await api.get('/api/me', { headers: { Authorization: `Bearer ${token}` } });
        const meData = me.data ?? {};
        const tenant = meData?.tenant ?? undefined;
        setUser({
          id: meData?.id ?? id,
          nome: meData?.nome ?? nome,
          email: meData?.email ?? email,
          tenant,
        });
      } catch {
        // já definido acima com dados do login
      }


      Router.push('/')


    } catch (err) {
      const axiosErr = err as AxiosError<{ error?: string; message?: string; code?: string }>;
      const message =
        axiosErr.response?.data?.error ||
        axiosErr.response?.data?.message ||
        axiosErr.message ||
        'Erro ao autenticar';
      const code =
        typeof axiosErr.response?.data?.code === 'string'
          ? axiosErr.response.data.code.trim()
          : undefined;
      throw new AuthError(message, code);

    }
  }

  async function signUp({ nome, email, password, empresa_nome }: SignUpProps): Promise<SignUpResult> {
    try {
      const response = await api.post("/api/registro", {
        nome,
        email,
        password,
        empresa_nome,
      })

      Router.push('/auth/login')
      return response.data as SignUpResult;

    } catch (err) {
      const axiosErr = err as AxiosError<{ error?: string; message?: string }>;
      const message =
        axiosErr.response?.data?.error ||
        axiosErr.response?.data?.message ||
        (err instanceof Error ? err.message : 'Erro ao registrar usuário');
      throw new Error(message)
    }
  }

  async function logoutUser() {
    try {
      queryClient.clear();
      setUser(undefined);
      api.defaults.headers.common['Authorization'] = '';
      clearAuthTokenCookies(null);
      if (typeof window !== 'undefined') {
        window.sessionStorage.removeItem('vecontab_token');
        window.localStorage.removeItem('vecontab_token');
        window.dispatchEvent(new Event(AUTH_SESSION_CHANGED_EVENT));
      }
      // Navega imediatamente para evitar qualquer frame da área autenticada.
      await Router.replace('/auth/login');
      // Auditoria de encerramento de sessão não pode bloquear o logout visual.
      void api.post('/api/session/end').catch(() => {
        // segue mesmo com auditoria indisponivel
      });
    } catch (err) {
      //console.log("Erro ao Sair", err)
      const message = err instanceof Error ? err.message : 'Erro ao sair';
      throw new Error(message)

    }
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated,
        authBootstrapped,
        signIn,
        signUp,
        logoutUser,
      }}>
      {children}
    </AuthContext.Provider>
  )
}

export default AuthContext;