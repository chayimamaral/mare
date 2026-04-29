import React, { useRef } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { Card } from 'primereact/card';
import { Column } from 'primereact/column';
import { DataTable } from 'primereact/datatable';
import { Toast } from 'primereact/toast';
import api from '../../components/api/apiClient';
import { useRouteClientGuard } from '../../components/hooks/useClientGuards';

type ActiveSessionRow = {
  id: string;
  user_id: string;
  user_email: string;
  tenant_id: string;
  tenant_name: string;
  tenant_cnpj: string;
  tenant_contato: string;
  logged_at: string;
  ip_address: string;
  active: boolean;
};

function formatDt(iso: string) {
  if (!iso) return '—';
  try {
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) return iso;
    return d.toLocaleString('pt-BR');
  } catch {
    return iso;
  }
}

export default function MonitoramentoGlobalPage() {
  useRouteClientGuard();
  const toast = useRef<Toast>(null);

  const { data: userRole = null } = useQuery<string | null>({
    queryKey: ['monitoramento-global-role'],
    queryFn: async () => {
      const r = await api.get('/api/usuariorole');
      const raw = r.data?.logado?.role;
      if (typeof raw !== 'string') return null;
      return raw.trim().toUpperCase() || null;
    },
    staleTime: 0,
    retry: 2,
  });

  const isSuper = userRole === 'SUPER';

  const {
    data: rows = [],
    isFetching,
    refetch,
  } = useQuery<ActiveSessionRow[]>({
    queryKey: ['monitoramento-global-sessoes'],
    enabled: isSuper,
    queryFn: async () => {
      const { data } = await api.get<ActiveSessionRow[]>('/api/monitoramento-global/sessoes');
      return Array.isArray(data) ? data : [];
    },
    refetchInterval: 5000,
    refetchIntervalInBackground: true,
  });

  return (
    <div className="grid">
      <div className="col-12">
        <Card title="Monitoramento global">
          <Toast ref={toast} />
          <p className="text-600 mt-0">
            Sessões ativas registradas no banco de auditoria (VECX_AUDIT). Atualização automática a cada 5 segundos.
          </p>
          <div className="flex justify-content-end mb-3">
            <Button
              type="button"
              label="Atualizar agora"
              icon="pi pi-refresh"
              onClick={() => {
                void refetch();
                toast.current?.show({ severity: 'info', summary: 'Lista', detail: 'Atualizando…', life: 1500 });
              }}
              loading={isFetching}
              disabled={!isSuper}
            />
          </div>
          <DataTable
            value={rows}
            dataKey="id"
            emptyMessage={isSuper ? 'Nenhuma sessão ativa no momento.' : 'Acesso restrito a SUPER.'}
            loading={isFetching}
            stripedRows
          >
            <Column field="user_email" header="Usuário" sortable style={{ minWidth: '14rem' }} />
            <Column field="tenant_name" header="Escritório" sortable style={{ minWidth: '12rem' }} />
            <Column field="tenant_cnpj" header="CNPJ" style={{ minWidth: '10rem' }} />
            <Column field="tenant_contato" header="Contato" style={{ minWidth: '10rem' }} />
            <Column
              field="logged_at"
              header="Login"
              sortable
              body={(r: ActiveSessionRow) => formatDt(r.logged_at)}
              style={{ minWidth: '11rem' }}
            />
            <Column field="ip_address" header="IP" style={{ minWidth: '8rem' }} />
          </DataTable>
        </Card>
      </div>
    </div>
  );
}
