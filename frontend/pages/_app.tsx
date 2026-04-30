import type { AppProps } from 'next/app';
import { LayoutConfig, type Page } from '../types/types';
import React from 'react';
import { LayoutProvider } from '../layout/context/layoutcontext';
import Layout from '../layout/layout';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import 'primereact/resources/primereact.css';
import 'primeflex/primeflex.css';
import 'primeicons/primeicons.css';
import '../styles/layout/layout.scss';
import { AUTH_SESSION_CHANGED_EVENT, AuthProvider } from '../components/context/AuthContext';
import { CaixaPostalProvider } from '../components/context/CaixaPostalContext';
import { useRouteClientGuard } from '../components/hooks/useClientGuards';
import { useIdleLogout } from '../components/hooks/useIdleLogout';
// import userPersistedState from '../components/utils/usePersistedState';

type Props = AppProps & {
    Component: Page;
};

function AppContent({ Component, pageProps }: Props) {
    useRouteClientGuard();
    useIdleLogout(); // Monitor de inatividade padrão (15 mins)


    if (Component.getLayout) {
        return <LayoutProvider>{Component.getLayout(<Component {...pageProps} />)}</LayoutProvider>;
    }

    return (
        <LayoutProvider>
            <Layout>
                <Component {...pageProps} />
            </Layout>
        </LayoutProvider>
    );
}

export default function App({ Component, pageProps }: Props) {
    const [authSessionStamp, setAuthSessionStamp] = React.useState('boot');
    const [queryClient] = React.useState(
        () =>
            new QueryClient({
                defaultOptions: {
                    queries: {
                        refetchOnWindowFocus: false,
                        retry: 1,
                        /**
                         * staleTime 0: dados nascem obsoletos — ao montar a página, refetch típico.
                         * Evita perfil/permissões “congelados” (ex.: usuariorole null em cache por minutos).
                         * gcTime mantém resultado em memória após desmontar (dedupe, voltar à página).
                         * “Zero cache” absoluto não existe no TanStack; para ainda mais agressivo use gcTime: 0.
                         */
                        staleTime: 0,
                        gcTime: 1000 * 60 * 10,
                    },
                },
            })
    );

    const defaultLayoutConfig: LayoutConfig = {
        theme: 'dark',
        ripple: true,
        inputStyle: 'outlined',
        menuMode: 'layout-menu-light',
        colorScheme: 'light',
        scale: 10,
    };

    // const [layoutConfig, setLayoutConfig] = userPersistedState<LayoutConfig>('theme', {
    //     theme: 'dark',
    //     ripple: true,
    //     inputStyle: 'outlined',
    //     menuMode: 'layout-menu-light',
    //     colorScheme: 'light',
    //     scale: 10
    // });

    // const [layoutConfig, setLayoutConfig] = userPersistedState<LayoutConfig>('theme', defaultLayoutConfig);

    // useEffect(() => {
    //   if (!layoutConfig) {
    //     setLayoutConfig(defaultLayoutConfig);
    //   }
    //   localStorage.setItem('layoutConfig', JSON.stringify(layoutConfig));
    // }, [layoutConfig]);

    React.useEffect(() => {
        if (typeof window === 'undefined') {
            return;
        }
        const readTokenFingerprint = () => {
            const token = String(window.sessionStorage.getItem('vecontab_token') ?? window.localStorage.getItem('vecontab_token') ?? '').trim();
            return token ? `tk:${token.slice(0, 16)}` : 'anon';
        };
        const handleAuthSessionChanged = () => {
            queryClient.clear();
            setAuthSessionStamp(`s:${readTokenFingerprint()}:${Date.now()}`);
        };
        setAuthSessionStamp(`s:${readTokenFingerprint()}:${Date.now()}`);
        window.addEventListener(AUTH_SESSION_CHANGED_EVENT, handleAuthSessionChanged);
        return () => {
            window.removeEventListener(AUTH_SESSION_CHANGED_EVENT, handleAuthSessionChanged);
        };
    }, [queryClient]);

    return (
        <QueryClientProvider client={queryClient}>
            <AuthProvider>
                <CaixaPostalProvider>
                    <AppContent key={authSessionStamp} Component={Component} pageProps={pageProps} />
                </CaixaPostalProvider>
            </AuthProvider>
        </QueryClientProvider>
    );
}

