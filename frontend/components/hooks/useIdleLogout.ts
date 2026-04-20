import { useEffect, useRef, useContext } from 'react';
import AuthContext from '../context/AuthContext';

const DEFAULT_IDLE_TIMEOUT = 15 * 60 * 1000; // 15 minutos

export function useIdleLogout(idleTimeMs: number = DEFAULT_IDLE_TIMEOUT) {
    const { isAuthenticated, logoutUser } = useContext(AuthContext);
    const timeoutRef = useRef<NodeJS.Timeout | null>(null);

    useEffect(() => {
        if (!isAuthenticated) return;

        const resetTimer = () => {
            if (timeoutRef.current) clearTimeout(timeoutRef.current);
            timeoutRef.current = setTimeout(() => {
                logoutUser().catch(() => {});
            }, idleTimeMs);
        };

        const events = ['mousemove', 'keydown', 'wheel', 'DOMMouseScroll', 'mouseWheel', 'mousedown', 'touchstart', 'touchmove', 'MSPointerDown', 'MSPointerMove'];
        
        events.forEach(event => document.addEventListener(event, resetTimer));

        // Inicializa o primeiro contador
        resetTimer();

        return () => {
            if (timeoutRef.current) clearTimeout(timeoutRef.current);
            events.forEach(event => document.removeEventListener(event, resetTimer));
        };
    }, [isAuthenticated, idleTimeMs, logoutUser]);
}
