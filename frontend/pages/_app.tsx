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
import { AuthProvider } from '../components/context/AuthContext';
// import userPersistedState from '../components/utils/usePersistedState';

type Props = AppProps & {
    Component: Page;
};

export default function App({ Component, pageProps }: Props) {
    const [queryClient] = React.useState(
        () =>
            new QueryClient({
                defaultOptions: {
                    queries: {
                        refetchOnWindowFocus: false,
                        retry: 1,
                        // Adicionando os 5 minutos de cache "fresco"
                        staleTime: 1000 * 60 * 5, 
                        // Opcional: tempo que o dado fica em memória após sumir da tela (v5)
                        gcTime: 1000 * 60 * 10,
                    }
                }
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

    if (Component.getLayout) {
        return (
            <QueryClientProvider client={queryClient}>
                <AuthProvider>
                    <LayoutProvider>{Component.getLayout(<Component {...pageProps} />)}</LayoutProvider>;
                </AuthProvider>
            </QueryClientProvider>
        )
    } else {
        return (
            <QueryClientProvider client={queryClient}>
                <AuthProvider>
                    <LayoutProvider>
                        <Layout>
                            <Component {...pageProps} />
                        </Layout>
                    </LayoutProvider>
                </AuthProvider>
            </QueryClientProvider>
        );
    }
}
