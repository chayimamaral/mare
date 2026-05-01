import { useEffect } from 'react';
import { useRouter } from 'next/router';
import { parseCookies } from 'nookies';
import { useQuery } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import api from '../api/apiClient';
import {
    clearAuthTokenCookies,
    getAuthTokenFromParsedCookies,
} from '../../constants/authCookie';
import {
    FEATURE,
    PATH_REQUIRES_CORE,
    PATH_REQUIRES_FEATURE_SLUG,
    sessionAllowsFeature,
} from '../../constants/featureAccess';
import { AUTH_SESSION_CHANGED_EVENT } from '../context/AuthContext';

type UserRole = 'SUPER' | 'ADMIN' | 'USER' | 'REPRESENTANTE';

const GUEST_ONLY_ROUTES = new Set<string>(['/auth/login', '/auth/register']);

const AUTH_REQUIRED_ROUTES = new Set<string>([
    '/agenda',
    '/agenda-arvore',
    '/catalogo-servicos',
    '/cliente-pf',
    '/clientes',
    '/cnae',
    '/caixa-postal',
    '/compromissos',
    '/compromissos-empresas',
    '/compromissos-por-natureza',
    '/compromissos-visao',
    '/configuracoes/api-integra-contador',
    '/configuracoes/certificado-digital',
    '/configuracoes/geracao-guias',
    '/configuracoes/integra-contador-servicos',
    '/configuracoes/integra-contador-tabela-consumo',
    '/empresas',
    '/estados',
    '/feriados',
    '/grupopassos',
    '/matriz-conformidade-fiscal',
    '/monitor',
    '/utilitarios/hardware-manager',
    '/utilitarios/monitoramento-global',
    '/nfe/consulta',
    '/nfe/manutencao',
    '/nfe/sincronizacao',
    '/municipios',
    '/obrigacoes',
    '/passos',
    '/regimes-tributarios',
    '/registro',
    '/representantes',
    '/rotinas',
    '/rotinas-pf',
    '/salario-minimo',
    '/tenants',
    '/tipoempresa',
    '/usuarios',
]);

const ROLE_RESTRICTED_ROUTES: Partial<Record<string, UserRole[]>> = {
    '/admin/broadcast': ['SUPER'],
    '/catalogo-servicos': ['SUPER'],
    '/configuracoes/api-integra-contador': ['SUPER'],
    '/configuracoes/integra-contador-servicos': ['SUPER'],
    '/configuracoes/integra-contador-tabela-consumo': ['SUPER'],
    '/matriz-conformidade-fiscal': ['SUPER'],
    '/representantes': ['SUPER'],
    '/monitor': ['SUPER', 'ADMIN'],
    '/tenants': ['SUPER', 'REPRESENTANTE'],
    '/usuarios': ['SUPER', 'ADMIN'],
    '/utilitarios/hardware-manager': ['SUPER'],
    '/utilitarios/monitoramento-global': ['SUPER'],
};

type GuardSession = {
    role: string;
    featureSlugs: string[] | undefined;
};

type GuestSessionValidation = {
    valid: boolean;
    unauthorized: boolean;
};

export function useRouteClientGuard(): void {
    const router = useRouter();
    const pathname = router.pathname;
    const isGuestOnly = GUEST_ONLY_ROUTES.has(pathname);
    const needsAuth = AUTH_REQUIRED_ROUTES.has(pathname) || Boolean(ROLE_RESTRICTED_ROUTES[pathname]);
    const rolesPermitidas = ROLE_RESTRICTED_ROUTES[pathname];
    const featureSlugForPath = PATH_REQUIRES_FEATURE_SLUG[pathname];
    const needsCore = PATH_REQUIRES_CORE.has(pathname);
    const needsFeatureCheck = Boolean(featureSlugForPath || needsCore);

    const cookieToken = getAuthTokenFromParsedCookies(parseCookies());
    const token =
        cookieToken ||
        (typeof window !== 'undefined' ? String(window.localStorage.getItem('vecontab_token') ?? '').trim() : '');

    const { data: guardData, isFetching: guardLoading } = useQuery({
        // Inclui token na chave para não reaproveitar cache de outro login
        // (ex.: usuário anterior USER causando redirect indevido para '/' após login SUPER).
        queryKey: ['route-role-feature-guard', pathname, token],
        enabled: !!token && (!!rolesPermitidas || needsFeatureCheck),
        queryFn: async (): Promise<GuardSession> => {
            const { data } = await api.get<{
                logado?: { role?: string; feature_slugs?: string[] };
            }>('/api/usuariorole');
            const logado = data?.logado;
            const role = String(logado?.role ?? '').trim().toUpperCase();
            const hasKey = Boolean(logado && Object.prototype.hasOwnProperty.call(logado, 'feature_slugs'));
            const featureSlugs = hasKey ? logado?.feature_slugs ?? [] : undefined;
            return { role, featureSlugs };
        },
    });

    const {
        data: guestSessionValidation,
        isFetching: guestSessionValidationLoading,
        isFetched: guestSessionValidationFetched,
    } = useQuery({
        queryKey: ['guest-session-validation', pathname, token],
        enabled: isGuestOnly && !!token,
        retry: false,
        queryFn: async (): Promise<GuestSessionValidation> => {
            try {
                await api.get('/api/me');
                return { valid: true, unauthorized: false };
            } catch (error) {
                const axiosErr = error as AxiosError;
                const status = axiosErr?.response?.status;
                if (status === 401 || status === 403) {
                    return { valid: false, unauthorized: true };
                }
                throw error;
            }
        },
    });

    useEffect(() => {
        const cookies = parseCookies();
        const cookieToken = getAuthTokenFromParsedCookies(cookies);
        const token =
            cookieToken ||
            (typeof window !== 'undefined' ? String(window.localStorage.getItem('vecontab_token') ?? '').trim() : '');

        if (isGuestOnly) {
            if (!token) {
                return;
            }
            if (guestSessionValidationLoading) {
                return;
            }
            if (guestSessionValidation?.valid) {
                void router.replace('/');
                return;
            }
            if (guestSessionValidationFetched && guestSessionValidation?.unauthorized) {
                // Token invalido no /auth/login nao pode redirecionar para dashboard nem piscar layout autenticado.
                clearAuthTokenCookies(null);
                if (typeof window !== 'undefined') {
                    window.sessionStorage.removeItem('vecontab_token');
                    window.localStorage.removeItem('vecontab_token');
                    window.dispatchEvent(new Event(AUTH_SESSION_CHANGED_EVENT));
                }
            }
            return;
        }

        if (needsAuth && !token) {
            void router.replace('/auth/login');
            return;
        }

        if ((rolesPermitidas || needsFeatureCheck) && guardLoading) {
            return;
        }

        if (rolesPermitidas && guardData && !rolesPermitidas.includes(guardData.role as UserRole)) {
            void router.replace('/');
            return;
        }

        if (needsFeatureCheck && guardData) {
            if (featureSlugForPath && !sessionAllowsFeature(guardData.role, guardData.featureSlugs, featureSlugForPath)) {
                void router.replace('/');
                return;
            }
            if (needsCore && !sessionAllowsFeature(guardData.role, guardData.featureSlugs, FEATURE.core)) {
                void router.replace('/');
            }
        }
    }, [
        isGuestOnly,
        needsAuth,
        guardData,
        guardLoading,
        rolesPermitidas,
        needsFeatureCheck,
        featureSlugForPath,
        needsCore,
        guestSessionValidation,
        guestSessionValidationLoading,
        guestSessionValidationFetched,
        router,
    ]);
}

export function useTenantIdQuery() {
    const cookieToken = getAuthTokenFromParsedCookies(parseCookies());
    const token =
        cookieToken ||
        (typeof window !== 'undefined' ? String(window.localStorage.getItem('vecontab_token') ?? '').trim() : '');

    return useQuery<string>({
        queryKey: ['tenant-id-client'],
        enabled: !!token,
        queryFn: async () => {
            try {
                const { data } = await api.get('/api/usuariotenant');
                const tenantFromTenantEndpoint = String(data?.tenantid ?? '').trim();
                if (tenantFromTenantEndpoint) {
                    return tenantFromTenantEndpoint;
                }
            } catch {
                // tenta fallback abaixo
            }

            const { data } = await api.get('/api/me');
            const tenantFallback =
                data?.usuarios?.[0]?.resultado?.tenant?.id ??
                data?.tenant?.id ??
                data?.tenantid ??
                '';
            return String(tenantFallback).trim();
        },
        retry: 2,
    });
}

export function useUserIdQuery() {
    const cookieToken = getAuthTokenFromParsedCookies(parseCookies());
    const token =
        cookieToken ||
        (typeof window !== 'undefined' ? String(window.localStorage.getItem('vecontab_token') ?? '').trim() : '');

    return useQuery<string>({
        queryKey: ['user-id', token],
        enabled: !!token,
        queryFn: async () => {
            const { data } = await api.get('/api/usuariorole');
            return String(data?.logado?.id ?? '').trim();
        },
        retry: 2,
    });
}
