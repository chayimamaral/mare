import { useQuery } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { DataTable } from 'primereact/datatable';
import { Dialog } from 'primereact/dialog';
import { InputNumber } from 'primereact/inputnumber';
import { Toast } from 'primereact/toast';
import { Toolbar } from 'primereact/toolbar';
import React, { useRef, useState } from 'react';
import { withAuthServerSideProps } from '../../components/utils/crudUtils';
import SalarioMinimoService from '../../services/cruds/SalarioMinimoService';
import { Vec } from '../../types/types';

const SalarioMinimoPage = () => {
  const svc = SalarioMinimoService();
  const toast = useRef<Toast>(null);

  const [dialogVisible, setDialogVisible] = useState(false);
  const [deleteVisible, setDeleteVisible] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [item, setItem] = useState<Vec.SalarioMinimoNacional>({ id: '', ano: new Date().getFullYear(), valor: undefined });

  const { data, isFetching, refetch } = useQuery({
    queryKey: ['salario-minimo-nacional'],
    queryFn: () => svc.list(),
  });

  const openNew = () => {
    setItem({ id: '', ano: new Date().getFullYear(), valor: undefined });
    setSubmitted(false);
    setDialogVisible(true);
  };

  const editItem = (row: Vec.SalarioMinimoNacional) => {
    setItem({ ...row });
    setSubmitted(false);
    setDialogVisible(true);
  };

  const askDelete = (row: Vec.SalarioMinimoNacional) => {
    setItem(row);
    setDeleteVisible(true);
  };

  const save = async () => {
    setSubmitted(true);
    const ano = Number(item.ano ?? 0);
    const valor = Number(item.valor ?? 0);
    if (ano < 1994 || valor <= 0) {
      toast.current?.show({ severity: 'warn', summary: 'Validação', detail: 'Informe ano válido e valor maior que zero.', life: 4000 });
      return;
    }

    try {
      if (item.id) {
        await svc.update({ id: item.id, ano, valor });
      } else {
        await svc.create({ ano, valor });
      }
      toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Registro salvo.', life: 3000 });
      setDialogVisible(false);
      refetch();
    } catch (err: any) {
      const msg = err?.response?.data?.error || err?.message || 'Falha ao salvar salário mínimo.';
      toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 5000 });
    }
  };

  const remove = async () => {
    if (!item.id) return;
    try {
      await svc.remove(item.id);
      toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Registro removido.', life: 3000 });
      setDeleteVisible(false);
      refetch();
    } catch (err: any) {
      const msg = err?.response?.data?.error || err?.message || 'Falha ao remover salário mínimo.';
      toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 5000 });
    }
  };

  const valorBody = (row: Vec.SalarioMinimoNacional) => (
    <span>{Number(row.valor ?? 0).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })}</span>
  );

  const actionsBody = (row: Vec.SalarioMinimoNacional) => (
    <>
      <Button icon="pi pi-pencil" rounded severity="success" className="mr-2" onClick={() => editItem(row)} />
      <Button icon="pi pi-trash" rounded severity="warning" onClick={() => askDelete(row)} />
    </>
  );

  return (
    <div className="grid crud-demo">
      <div className="col-12">
        <div className="card">
          <Toast ref={toast} />

          <Toolbar
            className="mb-4"
            left={() => (
              <Button label="Novo" icon="pi pi-plus" severity="success" onClick={openNew} />
            )}
          />

          <DataTable
            value={data?.data?.salarios ?? []}
            dataKey="id"
            loading={isFetching}
            stripedRows
            size="small"
            emptyMessage="Nenhum salário mínimo cadastrado."
          >
            <Column field="ano" header="Ano vigente" style={{ width: '12rem' }} />
            <Column field="valor" header="Valor" body={valorBody} style={{ width: '14rem' }} />
            <Column body={actionsBody} headerStyle={{ width: '10rem' }} />
          </DataTable>

          <div className="mt-2">
            <Button type="button" icon="pi pi-refresh" tooltip="Atualizar" className="p-button-text" onClick={() => refetch()} />
          </div>

          <Dialog
            visible={dialogVisible}
            header={item.id ? 'Editar salário mínimo' : 'Novo salário mínimo'}
            style={{ width: 'min(95vw, 30rem)' }}
            modal
            onHide={() => setDialogVisible(false)}
            footer={
              <div>
                <Button label="Cancelar" text onClick={() => setDialogVisible(false)} />
                <Button label="Salvar" icon="pi pi-check" text onClick={save} />
              </div>
            }
          >
            <div className="field">
              <label htmlFor="ano">Ano vigente</label>
              <InputNumber
                inputId="ano"
                value={item.ano ?? null}
                useGrouping={false}
                min={1994}
                max={2100}
                onValueChange={(e) => setItem((prev) => ({ ...prev, ano: e.value ?? undefined }))}
                className={submitted && (!item.ano || Number(item.ano) < 1994) ? 'p-invalid' : ''}
              />
            </div>
            <div className="field">
              <label htmlFor="valor">Valor</label>
              <InputNumber
                inputId="valor"
                value={item.valor ?? null}
                mode="currency"
                currency="BRL"
                locale="pt-BR"
                min={0.01}
                minFractionDigits={2}
                maxFractionDigits={2}
                onValueChange={(e) => setItem((prev) => ({ ...prev, valor: e.value ?? undefined }))}
                className={submitted && (!item.valor || Number(item.valor) <= 0) ? 'p-invalid' : ''}
              />
            </div>
          </Dialog>

          <Dialog
            visible={deleteVisible}
            header="Confirmar exclusão"
            style={{ width: 'min(95vw, 28rem)' }}
            modal
            onHide={() => setDeleteVisible(false)}
            footer={
              <div>
                <Button label="Não" text onClick={() => setDeleteVisible(false)} />
                <Button label="Sim" text onClick={remove} />
              </div>
            }
          >
            <span>Excluir salário mínimo do ano <strong>{item.ano}</strong>?</span>
          </Dialog>
        </div>
      </div>
    </div>
  );
};

export default SalarioMinimoPage;

export const getServerSideProps = withAuthServerSideProps(async () => undefined);
