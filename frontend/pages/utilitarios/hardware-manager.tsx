import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { Card } from 'primereact/card';
import { Column } from 'primereact/column';
import { DataTable } from 'primereact/datatable';
import { ProgressSpinner } from 'primereact/progressspinner';
import { Toast } from 'primereact/toast';
import api from '../../components/api/apiClient';
import { useRouteClientGuard } from '../../components/hooks/useClientGuards';

type HardwareDevice = {
  device_id: string;
  path: string;
  status: string;
};

type LocalAgentCert = {
  id: string;
  label?: string;
  subject?: string;
  serial_hex?: string;
  slot_id?: number;
  token_label?: string;
};

type LocalAgentCertResponse = {
  items: LocalAgentCert[];
};

type WailsRuntimeLike = {
  EventsOn?: (eventName: string, cb: (...args: any[]) => void) => void | (() => void) | Promise<(() => void) | void>;
};

function getLocalAgentBaseURL(): string {
  const raw = String(process.env.NEXT_PUBLIC_LOCAL_AGENT_BASE_URL ?? '').trim();
  if (!raw) {
    return 'http://127.0.0.1:9999';
  }
  return raw.replace(/\/+$/, '');
}

/** Deve coincidir com AGENT_SHARED_SECRET do vecx-agent (se estiver definido lá). */
function getLocalAgentSharedSecret(): string {
  return String(process.env.NEXT_PUBLIC_LOCAL_AGENT_SECRET ?? '').trim();
}

function describeFetchFailure(err: unknown): string {
  const msg = err instanceof Error ? err.message : String(err);
  const lower = msg.toLowerCase();
  if (lower.includes('failed to fetch') || lower === 'networkerror when attempting to fetch resource.') {
    return (
      'O navegador não conseguiu completar o pedido (Failed to fetch). ' +
      'Causas frequentes: (1) vecx-agent parado ou porta errada; (2) CORS — confira AGENT_ALLOWED_ORIGINS no agente; ' +
      '(3) se o agente usa AGENT_SHARED_SECRET, o pedido precisa do cabeçalho X-Local-Agent-Secret e o CORS do agente deve permitir esse cabeçalho; ' +
      '(4) páginas HTTP remotas: o Chrome pode bloquear chamadas a 127.0.0.1 (rede privada); use HTTPS no site ou teste em localhost.'
    );
  }
  return msg;
}

export default function HardwareManagerPage() {
  useRouteClientGuard();
  const toast = useRef<Toast>(null);
  const [items, setItems] = useState<HardwareDevice[]>([]);
  const [scanning, setScanning] = useState(false);

  const {
    data: userRole = null,
    isLoading: roleLoading,
  } = useQuery<string | null>({
    queryKey: ['hardware-manager-user-role'],
    queryFn: async () => {
      const r = await api.get('/api/usuariorole');
      const raw = r.data?.logado?.role;
      if (typeof raw !== 'string') return null;
      const norm = raw.trim().toUpperCase();
      return norm || null;
    },
    staleTime: 0,
    retry: 2,
  });

  const isSuper = userRole === 'SUPER';

  const scanDevices = useCallback(async () => {
    setScanning(true);
    try {
      const baseURL = getLocalAgentBaseURL();
      const secret = getLocalAgentSharedSecret();
      const headers: Record<string, string> = {
        Accept: 'application/json',
      };
      if (secret !== '') {
        headers['X-Local-Agent-Secret'] = secret;
      }
      const resp = await fetch(`${baseURL}/certificates`, {
        method: 'GET',
        credentials: 'omit',
        headers,
      });
      if (!resp.ok) {
        if (resp.status === 401) {
          throw new Error(
            'vecx-agent recusou o segredo local (HTTP 401). ' +
              'Defina NEXT_PUBLIC_LOCAL_AGENT_SECRET no build do frontend com o mesmo valor de AGENT_SHARED_SECRET do agente, ' +
              'ou deixe AGENT_SHARED_SECRET vazio no agente para desativar esta exigência.',
          );
        }
        let body = await resp.text().catch(() => '');
        try {
          const j = JSON.parse(body) as { error?: string };
          if (j?.error === 'unauthorized') {
            body =
              'não autorizado no vecx-agent — verifique NEXT_PUBLIC_LOCAL_AGENT_SECRET / AGENT_SHARED_SECRET.';
          }
        } catch {
          /* manter body */
        }
        throw new Error(body || `vecx-agent respondeu HTTP ${resp.status}`);
      }
      const data = (await resp.json()) as LocalAgentCertResponse;

      const mapped = Array.isArray(data?.items)
        ? data.items.map((c) => ({
            device_id: String(c.id ?? '—'),
            path: [c.label, c.token_label, c.subject].filter(Boolean).join(' | ') || 'Certificado local',
            status: c.serial_hex ? `serial ${c.serial_hex}` : 'disponível',
          }))
        : [];
      const detectedCount = mapped.length;
      setItems(mapped);

      toast.current?.show({
        severity: 'success',
        summary: 'Escaneamento concluído',
        detail: detectedCount > 0
          ? `Dispositivos detectados: ${detectedCount}.`
          : 'Nenhum dispositivo detectado nesta varredura.',
        life: 2500,
      });
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e);
      const detail =
        msg.toLowerCase().includes('failed to fetch') || msg.toLowerCase().includes('networkerror')
          ? describeFetchFailure(e)
          : msg || describeFetchFailure(e);
      toast.current?.show({
        severity: 'error',
        summary: 'Erro no escaneamento',
        detail,
        life: 8000,
      });
    } finally {
      setScanning(false);
    }
  }, []);

  useEffect(() => {
    if (!isSuper) return;
    void scanDevices();
  }, [isSuper, scanDevices]);

  useEffect(() => {
    if (!isSuper) return;
    const rt = (window as unknown as { runtime?: WailsRuntimeLike }).runtime;
    if (!rt?.EventsOn) return;

    let unsubscribe: (() => void) | undefined;
    const onInserted = () => {
      void scanDevices();
    };

    const maybePromise = rt.EventsOn('hardware:usb-inserted', onInserted);
    if (typeof maybePromise === 'function') {
      unsubscribe = maybePromise;
    } else if (maybePromise && typeof (maybePromise as Promise<any>).then === 'function') {
      void (maybePromise as Promise<(() => void) | void>).then((fn) => {
        if (typeof fn === 'function') unsubscribe = fn;
      });
    }

    return () => {
      if (unsubscribe) unsubscribe();
    };
  }, [isSuper, scanDevices]);

  const header = useMemo(() => (
    <div className="flex flex-column md:flex-row md:justify-content-between md:align-items-center gap-2">
      <div>
        <h5 className="m-0">Hardware Manager</h5>
        <small className="text-600">Consulta direta no vecx-agent local (127.0.0.1:9999).</small>
      </div>
      <Button
        type="button"
        label="Escanear vecx-agent"
        icon="pi pi-refresh"
        onClick={() => void scanDevices()}
        loading={scanning}
      />
    </div>
  ), [scanDevices, scanning]);

  return (
    <div className="grid">
      <div className="col-12">
        <Toast ref={toast} />
        <Card>
          {roleLoading ? (
            <div className="flex flex-column align-items-center gap-3 py-6">
              <ProgressSpinner style={{ width: '3rem', height: '3rem' }} />
              <span className="text-600">Carregando permissões…</span>
            </div>
          ) : null}

          {!roleLoading && !isSuper ? (
            <div className="p-4 border-round border-1 surface-border bg-red-50 text-red-700">
              Acesso negado: esta página é restrita ao perfil SUPER.
            </div>
          ) : null}

          {!roleLoading && isSuper ? (
            <div>
              {header}
              <div className="mt-3">
                <DataTable
                  value={items}
                  dataKey="device_id"
                  loading={scanning}
                  stripedRows
                  size="small"
                  emptyMessage="Nenhum dispositivo detectado."
                >
                  <Column field="device_id" header="ID do Dispositivo" style={{ minWidth: '14rem' }} />
                  <Column field="path" header="Nome/Caminho" style={{ minWidth: '24rem' }} />
                  <Column field="status" header="Status" style={{ minWidth: '10rem' }} />
                </DataTable>
              </div>
            </div>
          ) : null}
        </Card>
      </div>
    </div>
  );
}

