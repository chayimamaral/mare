import React, { useCallback, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { InputTextarea } from 'primereact/inputtextarea';

import setupAPIClient from '../api/api';

type ChatMsg = { role: 'user' | 'assistant'; text: string };

function getToken(): string {
  if (typeof window === 'undefined') {
    return '';
  }
  try {
    return String(
      window.sessionStorage.getItem('vecontab_token') ?? window.localStorage.getItem('vecontab_token') ?? '',
    ).trim();
  } catch {
    return '';
  }
}

function parseSSEBuffer(buf: string, onData: (o: Record<string, unknown>) => void): string {
  const blocks = buf.split('\n\n');
  let rest = blocks.pop() ?? '';
  for (const block of blocks) {
    for (const line of block.split('\n')) {
      const t = line.trim();
      if (!t.startsWith('data:')) {
        continue;
      }
      const payload = t.slice(5).trim();
      if (!payload) {
        continue;
      }
      try {
        onData(JSON.parse(payload) as Record<string, unknown>);
      } catch {
        // ignora linha inválida
      }
    }
  }
  return rest;
}

const VirtualAssistantChat: React.FC = () => {
  const api = setupAPIClient(undefined);
  const { data: iaPublic } = useQuery({
    queryKey: ['public-ia'],
    queryFn: async () => {
      const { data } = await api.get<{ iaEnabled: boolean }>('/api/public/ia');
      return data;
    },
    staleTime: 60_000,
  });

  const [open, setOpen] = useState(false);
  const [input, setInput] = useState('');
  const [messages, setMessages] = useState<ChatMsg[]>([]);
  const [sending, setSending] = useState(false);
  const abortRef = useRef<AbortController | null>(null);

  const resetPanel = useCallback(() => {
    abortRef.current?.abort();
    abortRef.current = null;
    setMessages([]);
    setInput('');
    setSending(false);
  }, []);

  const handleClose = useCallback(() => {
    setOpen(false);
    resetPanel();
  }, [resetPanel]);

  const send = useCallback(async () => {
    const text = input.trim();
    if (!text || sending) {
      return;
    }
    const token = getToken();
    if (!token) {
      return;
    }

    setInput('');
    setMessages((m) => [...m, { role: 'user', text }]);
    setSending(true);

    const base = String(api.defaults.baseURL ?? '').replace(/\/+$/, '');
    const url = `${base}/api/ai/chat`;

    const ac = new AbortController();
    abortRef.current = ac;

    let assistantBuf = '';

    setMessages((m) => [...m, { role: 'assistant', text: '' }]);

    try {
      const res = await fetch(url, {
        method: 'POST',
        signal: ac.signal,
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ message: text }),
      });

      if (res.status === 503) {
        setMessages((m) => {
          const copy = [...m];
          const last = copy[copy.length - 1];
          if (last?.role === 'assistant') {
            copy[copy.length - 1] = { role: 'assistant', text: 'Serviço Indisponível.' };
          }
          return copy;
        });
        return;
      }

      if (!res.ok || !res.body) {
        let detail = `Erro HTTP ${res.status}.`;
        try {
          const text = await res.text();
          const trimmed = text.trim();
          if (trimmed.startsWith('{')) {
            const j = JSON.parse(trimmed) as { error?: unknown };
            if (typeof j?.error === 'string' && j.error.trim()) {
              detail = j.error.trim();
            }
          } else if (trimmed) {
            detail = trimmed.slice(0, 500);
          }
        } catch {
          // mantém detail genérico
        }
        setMessages((m) => {
          const copy = [...m];
          const last = copy[copy.length - 1];
          if (last?.role === 'assistant') {
            copy[copy.length - 1] = { role: 'assistant', text: detail };
          }
          return copy;
        });
        return;
      }

      const reader = res.body.getReader();
      const dec = new TextDecoder();
      let raw = '';

      const applyDelta = (delta: string) => {
        assistantBuf += delta;
        setMessages((m) => {
          const copy = [...m];
          const last = copy[copy.length - 1];
          if (last?.role === 'assistant') {
            copy[copy.length - 1] = { role: 'assistant', text: assistantBuf };
          }
          return copy;
        });
      };

      const onSSEObj = (obj: Record<string, unknown>) => {
        if (typeof obj.error === 'string' && obj.error) {
          assistantBuf = String(obj.error);
          setMessages((m) => {
            const copy = [...m];
            const last = copy[copy.length - 1];
            if (last?.role === 'assistant') {
              copy[copy.length - 1] = { role: 'assistant', text: assistantBuf };
            }
            return copy;
          });
          return;
        }
        if (typeof obj.delta === 'string' && obj.delta) {
          applyDelta(obj.delta);
        }
      };

      for (;;) {
        const { done, value } = await reader.read();
        if (done) {
          break;
        }
        raw += dec.decode(value, { stream: true });
        raw = parseSSEBuffer(raw, onSSEObj);
      }
      if (raw.trim()) {
        parseSSEBuffer(raw + '\n\n', onSSEObj);
      }
    } catch (e) {
      if ((e as Error)?.name === 'AbortError') {
        return;
      }
      setMessages((m) => {
        const copy = [...m];
        const last = copy[copy.length - 1];
        if (last?.role === 'assistant') {
          copy[copy.length - 1] = { role: 'assistant', text: 'Falha na conexão com o assistente.' };
        }
        return copy;
      });
    } finally {
      setSending(false);
      abortRef.current = null;
    }
  }, [api, input, sending]);

  if (!iaPublic?.iaEnabled) {
    return null;
  }

  return (
    <>
      <button
        type="button"
        className="vecx-ai-fab border-none cursor-pointer shadow-2 flex align-items-center justify-content-center"
        style={{
          position: 'fixed',
          zIndex: 11000,
          right: '1.25rem',
          bottom: '1.25rem',
          width: '3.25rem',
          height: '3.25rem',
          borderRadius: '50%',
          background: 'var(--primary-color, #0ea5e9)',
          color: '#fff',
        }}
        aria-label="Abrir assistente virtual"
        onClick={() => setOpen(true)}
      >
        <span className="pi pi-headphones" style={{ fontSize: '1.25rem' }} aria-hidden />
      </button>

      {open ? (
        <div
          className="vecx-ai-panel surface-card border-round-lg shadow-8 flex flex-column overflow-hidden"
          style={{
            position: 'fixed',
            zIndex: 11001,
            right: '1.25rem',
            bottom: '5rem',
            width: 'min(100vw - 2rem, 22rem)',
            height: 'min(70vh, 28rem)',
          }}
        >
          <div
            className="flex align-items-center justify-content-between px-3 py-2"
            style={{ borderBottom: '1px solid var(--surface-border, #dee2e6)' }}
          >
            <div className="flex align-items-center gap-2">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img src="/x.png" alt="" width={22} height={22} />
              <span className="font-semibold text-sm">Assistente VECX</span>
            </div>
            <button
              type="button"
              className="border-none bg-transparent cursor-pointer p-0 flex align-items-center justify-content-center"
              style={{ color: 'var(--text-color-secondary)' }}
              aria-label="Fechar assistente"
              onClick={handleClose}
            >
              <span className="pi pi-times" style={{ fontSize: '1.1rem' }} aria-hidden />
            </button>
          </div>

          <div className="flex-1 overflow-auto px-3 py-2 text-sm" style={{ flex: 1, minHeight: 0 }}>
            {messages.length === 0 ? (
              <p className="text-color-secondary m-0">Envie uma dúvida sobre o sistema ou processos contábeis.</p>
            ) : (
              messages.map((m, i) => (
                <div key={i} className={`mb-2 ${m.role === 'user' ? 'text-right' : 'text-left'}`}>
                  <span
                    className="inline-block px-2 py-1 border-round"
                    style={{
                      background: m.role === 'user' ? 'var(--primary-100, #e0f2fe)' : 'var(--surface-100, #f4f4f5)',
                      maxWidth: '100%',
                      whiteSpace: 'pre-wrap',
                      wordBreak: 'break-word',
                    }}
                  >
                    {m.text}
                  </span>
                </div>
              ))
            )}
          </div>

          <div className="p-2 flex flex-column gap-2" style={{ borderTop: '1px solid var(--surface-border, #dee2e6)' }}>
            <InputTextarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              rows={2}
              className="w-full text-sm"
              placeholder="Digite sua pergunta…"
              disabled={sending}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && !e.shiftKey) {
                  e.preventDefault();
                  void send();
                }
              }}
            />
            <Button label={sending ? 'Enviando…' : 'Enviar'} className="w-full" disabled={sending} onClick={() => void send()} />
          </div>
        </div>
      ) : null}
    </>
  );
};

export default VirtualAssistantChat;
