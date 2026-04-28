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

type ScanResponse = {
  items: HardwareDevice[];
};

type WailsRuntimeLike = {
  EventsOn?: (eventName: string, cb: (...args: any[]) => void) => void | (() => void) | Promise<(() => void) | void>;
};

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
      const endpointCandidates = [
        '/api/hardware/dispositivos-locais',
        '/hardware/dispositivos-locais',
      ];
      let data: ScanResponse | null = null;
      let lastError: unknown = null;
      for (const ep of endpointCandidates) {
        try {
          const resp = await api.get<ScanResponse>(ep);
          data = resp.data;
          break;
        } catch (e: unknown) {
          const ax = e as { response?: { status?: number } };
          lastError = e;
          // Se nao for 404, nao adianta tentar o proximo endpoint.
          if (ax?.response?.status !== 404) {
            throw e;
          }
        }
      }
      if (!data) {
        throw lastError ?? new Error('rota de hardware indisponivel');
      }

      setItems(Array.isArray(data?.items) ? data.items : []);
      toast.current?.show({
        severity: 'success',
        summary: 'Escaneamento concluído',
        detail: Array.isArray(data?.items) && data.items.length > 0
          ? `Dispositivos detectados: ${data.items.length}.`
          : 'Nenhum dispositivo USB local detectado nesta varredura.',
        life: 2500,
      });
    } catch (e: unknown) {
      const err = e as { response?: { status?: number; data?: { error?: string } } };
      const detail = err?.response?.status === 404
        ? 'Endpoint de hardware nao encontrado no backend em execucao. Recompile e reinicie o vecx-backend.'
        : (err?.response?.data?.error ?? 'Falha ao listar dispositivos locais.');
      toast.current?.show({
        severity: 'error',
        summary: 'Erro no escaneamento',
        detail,
        life: 5000,
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
        <small className="text-600">Detecção local de dispositivos (modo desktop/binário).</small>
      </div>
      <Button
        type="button"
        label="Escanear Dispositivos Locais"
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

