import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { DataTable } from 'primereact/datatable';
import { Toast } from 'primereact/toast';
import React, { useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import MonitorOperacaoService from '../../services/cruds/MonitorOperacaoService';
import { Vec } from '../../types/types';

const fmtDetalhe = (d?: Record<string, unknown>) => {
  if (!d || typeof d !== 'object') {
    return '';
  }
  try {
    return JSON.stringify(d);
  } catch {
    return '';
  }
};

const MonitorPage = () => {
  const [itensFallback, setItensFallback] = useState<Vec.MonitorOperacaoItem[]>([]);
  const [totalRecordsFallback, setTotalRecordsFallback] = useState(0);
  const [first, setFirst] = useState(0);
  const [rows, setRows] = useState(25);
  const toast = useRef<Toast>(null);

  const load = async () => {
    const { itens: lista, total } = await MonitorOperacaoService().list(rows, first);
    return {
      itens: lista ?? [],
      totalRecords: typeof total === 'number' ? total : 0,
    };
  };

  const { data, isFetching, refetch } = useQuery({
    queryKey: ['monitor-operacoes', rows, first],
    queryFn: async () => {
      try {
        const next = await load();
        setItensFallback(next.itens);
        setTotalRecordsFallback(next.totalRecords);
        return next;
      } catch (e: unknown) {
        const ax = e as { response?: { data?: { error?: string } } };
        toast.current?.show({
          severity: 'error',
          summary: 'Erro',
          detail: ax?.response?.data?.error ?? 'Falha ao carregar o monitor',
          life: 5000,
        });
        return { itens: [], totalRecords: 0 };
      }
    },
  });

  const paginatorLeft = (
    <Button
      type="button"
      icon="pi pi-refresh"
      tooltip="Atualizar"
      className="p-button-text"
      onClick={() => refetch()}
      loading={isFetching}
    />
  );

  const statusBody = (row: Vec.MonitorOperacaoItem) => (
    <span className={row.status === 'ERRO' ? 'text-red-500' : undefined}>{row.status}</span>
  );

  return (
    <div className="grid">
      <div className="col-12">
        <div className="card">
          <h5>Monitor de operações</h5>
          <p className="text-color-secondary text-sm mb-3">
            Registros de geração de compromissos, agenda e execuções automáticas. Cada linha tem um tenant_id;
            SUPER vê todos os clientes; ADMIN vê apenas seus clientes.
          </p>
          <Toast ref={toast} />
          <DataTable
            value={data?.itens ?? itensFallback}
            loading={isFetching}
            dataKey="id"
            paginator
            rows={rows}
            first={first}
            totalRecords={data?.totalRecords ?? totalRecordsFallback}
            lazy
            paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
            currentPageReportTemplate="{first} a {last} de {totalRecords}"
            rowsPerPageOptions={[10, 25, 50, 100]}
            paginatorLeft={paginatorLeft}
            onPage={(e) => {
              setFirst(e.first);
              setRows(e.rows);
            }}
            emptyMessage="Nenhum registro"
          >
            <Column
              field="criado_em"
              header="Data"
              body={(r: Vec.MonitorOperacaoItem) =>
                r.criado_em ? new Date(r.criado_em).toLocaleString('pt-BR') : ''
              }
              style={{ minWidth: '10rem' }}
            />
            <Column field="tenant_nome" header="Tenant" body={(r) => r.tenant_nome ?? r.tenant_id ?? '—'} />
            <Column field="user_id" header="Usuário" body={(r) => r.user_id ?? '—'} />
            <Column field="origem" header="Origem" />
            <Column field="tipo" header="Tipo" style={{ minWidth: '12rem' }} />
            <Column field="status" header="Status" body={statusBody} />
            <Column field="mensagem" header="Mensagem" style={{ minWidth: '14rem' }} />
            <Column
              header="Detalhe"
              body={(r: Vec.MonitorOperacaoItem) => (
                <span className="text-sm" style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>
                  {fmtDetalhe(r.detalhe)}
                </span>
              )}
            />
          </DataTable>
        </div>
      </div>
    </div>
  );
};

export default MonitorPage;
