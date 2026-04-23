import { useEffect, useRef, type CSSProperties } from 'react';

type Props = {
    html: string;
};

/**
 * HTML gerado no backend (Saxon + XSLT SVRS).
 * Injeta via document.write em about:blank — em vários browsers srcDoc com HTML completo + scripts SEFAZ fica em branco.
 */
export function DanfeHtmlIframe({ html }: Props) {
    const ref = useRef<HTMLIFrameElement>(null);
    const t = html.trim();

    useEffect(() => {
        if (!t) {
            return;
        }
        const frame = ref.current;
        if (!frame) {
            return;
        }

        const inject = () => {
            try {
                const doc = frame.contentDocument;
                if (!doc) {
                    return;
                }
                doc.open();
                doc.write(t);
                doc.close();
            } catch {
                /* ignore */
            }
        };

        inject();
        frame.addEventListener('load', inject, { once: true });
    }, [t]);

    if (!t) {
        return null;
    }

    const style: CSSProperties = {
        width: '100%',
        height: '70vh',
        minHeight: 420,
        border: '1px solid var(--surface-border, #dee2e6)',
        display: 'block',
    };

    return <iframe ref={ref} title="DANFE" src="about:blank" style={style} />;
}
