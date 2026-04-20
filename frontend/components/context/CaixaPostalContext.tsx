import { createContext, ReactNode, useCallback, useContext, useEffect, useRef, useState } from 'react';
import CaixaPostalService from '../../services/cruds/CaixaPostalService';
import AuthContext from './AuthContext';

interface CaixaPostalContextData {
    naoLidas: number;
    refreshCount: () => void;
}

const CaixaPostalContext = createContext<CaixaPostalContextData>({
    naoLidas: 0,
    refreshCount: () => { },
});

export function CaixaPostalProvider({ children }: { children: ReactNode }) {
    const [naoLidas, setNaoLidas] = useState(0);
    const { isAuthenticated } = useContext(AuthContext);
    const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

    const refreshCount = useCallback(async () => {
        if (!isAuthenticated) return;
        try {
            const svc = CaixaPostalService();
            const count = await svc.contarNaoLidas();
            setNaoLidas(count);
        } catch {
            // silencioso
        }
    }, [isAuthenticated]);

    useEffect(() => {
        refreshCount();
        intervalRef.current = setInterval(refreshCount, 60_000);
        return () => {
            if (intervalRef.current) clearInterval(intervalRef.current);
        };
    }, [refreshCount]);

    return (
        <CaixaPostalContext.Provider value={{ naoLidas, refreshCount }}>
            {children}
        </CaixaPostalContext.Provider>
    );
}

export function useCaixaPostal() {
    return useContext(CaixaPostalContext);
}

export default CaixaPostalContext;
