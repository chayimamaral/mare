import { useQuery } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { DataTable } from 'primereact/datatable';
import { Dialog } from 'primereact/dialog';
import { InputNumber } from 'primereact/inputnumber';
import { InputText } from 'primereact/inputtext';
import { Toast } from 'primereact/toast';
import { Toolbar } from 'primereact/toolbar';
import React, { useRef, useState } from 'react';
import EnquadramentoJuridicoPorteService from '../../services/cruds/EnquadramentoJuridicoPorteService';
import { useRouteClientGuard } from '../../components/hooks/useClientGuards';
import { Vec } from '../../types/types';

function fmtBRL(n: number | null | undefined): string {
  if (n == null || Number.isNaN(Number(n))) {
    return '—';
  }
  return Number(n).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' });
}

const EnquadramentoJuridicoPage = () => {
  useRouteClientGuard();
  const svc = EnquadramentoJuridicoPorteService();
  const toast = useRef<Toast>(null);
  const anoAtual = new Date().getFullYear();

  const [filtroAno, setFiltroAno] = useState<number | null>(anoAtual);
  const [dialogVisible, setDialogVisible] = useState(false);
  const [deleteVisible, setDeleteVisible] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [semLimiteSuperior, setSemLimiteSuperior] = useState(false);
  const [item, setItem] = useState<Vec.EnquadramentoJuridicoPorte>({
    id: '',
    sigla: '',
    descricao: '',
    limite_inicial: 0,
    limite_final: null,
    ano_vigencia: anoAtual,
    ativo: true,
  });

  const { data, isFetching, refetch } = useQuery({
    queryKey: ['enquadramentos-juridicos-porte', filtroAno],
    queryFn: () => svc.list(filtroAno ?? undefined),
  });

  const openNew = () => {
    setSemLimiteSuperior(false);
    setItem({
      id: '',
      sigla: '',
      descricao: '',
      limite_inicial: 0,
      limite_final: null,
      ano_vigencia: filtroAno ?? anoAtual,
      ativo: true,
    });
    setSubmitted(false);
    setDialogVisible(true);
  };

  const editItem = (row: Vec.EnquadramentoJuridicoPorte) => {
    const limF = row.limite_final;
    setSemLimiteSuperior(limF == null);
    setItem({ ...row, limite_final: limF ?? null });
    setSubmitted(false);
    setDialogVisible(true);
  };

  const askDelete = (row: Vec.EnquadramentoJuridicoPorte) => {
    setItem(row);
    setDeleteVisible(true);
  };

  const save = async () => {
    setSubmitted(true);
    const sigla = String(item.sigla ?? '').trim();
    const desc = String(item.descricao ?? '').trim();
    const ano = Number(item.ano_vigencia ?? 0);
    const li = Number(item.limite_inicial ?? 0);
    const lf: number | null = semLimiteSuperior ? null : Number(item.limite_final ?? 0);
    if (!sigla || !desc || ano < 1995 || ano > 2100) {
      toast.current?.show({
        severity: 'warn',
        summary: 'Validação',
        detail: 'Preencha sigla, descrição e ano de vigência válidos.',
        life: 4000,
      });
      return;
    }
    if (!semLimiteSuperior && (lf == null || Number.isNaN(lf) || lf < li)) {
      toast.current?.show({
        severity: 'warn',
        summary: 'Validação',
        detail: 'Limite final deve ser maior ou igual ao inicial, ou marque “sem limite superior”.',
        life: 5000,
      });
      return;
    }

    try {
      if (item.id) {
        await svc.update({
          id: item.id,
          sigla,
          descricao: desc,
          limite_inicial: li,
          limite_final: semLimiteSuperior ? null : lf,
          ano_vigencia: ano,
          ativo: item.ativo !== false,
        });
      } else {
        await svc.create({
          sigla,
          descricao: desc,
          limite_inicial: li,
          limite_final: semLimiteSuperior ? null : lf,
          ano_vigencia: ano,
        });
      }
      toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Registro salvo.', life: 3000 });
      setDialogVisible(false);
      refetch();
    } catch (err: unknown) {
      const ax = err as { response?: { data?: { error?: string } } };
      const msg = ax?.response?.data?.error || (err instanceof Error ? err.message : 'Falha ao salvar.');
      toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 6000 });
    }
  };

  const remove = async () => {
    if (!item.id) return;
    try {
      await svc.remove(item.id);
      toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Registro removido.', life: 3000 });
      setDeleteVisible(false);
      refetch();
    } catch (err: unknown) {
      const ax = err as { response?: { data?: { error?: string } } };
      const msg = ax?.response?.data?.error || (err instanceof Error ? err.message : 'Falha ao remover.');
      toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 6000 });
    }
  };

  const faixaBody = (row: Vec.EnquadramentoJuridicoPorte) => (
    <span>
      De {fmtBRL(row.limite_inicial)} até {row.limite_final == null ? 'sem limite superior' : fmtBRL(row.limite_final)}
    </span>
  );

  const actionsBody = (row: Vec.EnquadramentoJuridicoPorte) => (
    <>
      <Button icon="pi pi-pencil" rounded severity="success" className="mr-2" type="button" onClick={() => editItem(row)} />
      <Button icon="pi pi-trash" rounded severity="warning" type="button" onClick={() => askDelete(row)} />
    </>
  );

  const rows = data?.data?.items ?? [];

  return (
    <div className="grid crud-demo">
      <div className="col-12">
        <div className="card">
          <Toast ref={toast} />

          <div className="mb-3">
            <h5 className="m-0">Enquadramento Jurídico (porte por faturamento)</h5>
            <p className="m-0 mt-2 text-600 text-sm">
              Faixas de faturamento anual por ano de vigência. Dados globais (schema public). O MEI é caso especial com faixa própria,
              conforme legislação.
            </p>
          </div>

          <Toolbar
            className="mb-4"
            left={() => (
              <div className="flex flex-wrap align-items-center gap-3">
                <Button label="Novo" icon="pi pi-plus" severity="success" type="button" onClick={openNew} />
                <span className="text-sm text-600">Filtrar por ano</span>
                <InputNumber
                  value={filtroAno}
                  onValueChange={(e) => setFiltroAno(e.value == null ? null : Number(e.value))}
                  useGrouping={false}
                  min={1995}
                  max={2100}
                  placeholder="Todos"
                  className="w-8rem"
                />
                <Button
                  type="button"
                  label="Limpar filtro"
                  className="p-button-text"
                  onClick={() => setFiltroAno(null)}
                />
              </div>
            )}
          />

          <DataTable
            value={rows}
            dataKey="id"
            loading={isFetching}
            stripedRows
            size="small"
            emptyMessage="Nenhum registro encontrado para o filtro."
          >
            <Column field="sigla" header="Sigla" style={{ width: '10rem' }} />
            <Column field="descricao" header="Descrição" />
            <Column header="Faixa (R$)" body={faixaBody} style={{ minWidth: '16rem' }} />
            <Column field="ano_vigencia" header="Ano vigência" style={{ width: '10rem' }} />
            <Column body={actionsBody} header="Ações" headerStyle={{ width: '10rem' }} />
          </DataTable>

          <div className="mt-2">
            <Button type="button" icon="pi pi-refresh" tooltip="Atualizar" className="p-button-text" onClick={() => void refetch()} />
          </div>

          <Dialog
            visible={dialogVisible}
            header={item.id ? 'Editar enquadramento' : 'Novo enquadramento'}
            style={{ width: 'min(95vw, 32rem)' }}
            modal
            onHide={() => setDialogVisible(false)}
            footer={
              <div>
                <Button label="Cancelar" text type="button" onClick={() => setDialogVisible(false)} />
                <Button label="Salvar" icon="pi pi-check" text type="button" onClick={() => void save()} />
              </div>
            }
          >
            <div className="field">
              <label htmlFor="ej_sigla">Sigla</label>
              <InputText
                id="ej_sigla"
                value={item.sigla ?? ''}
                onChange={(e) => setItem((p) => ({ ...p, sigla: e.target.value }))}
                className={classNamesInvalid(submitted && !String(item.sigla ?? '').trim())}
                maxLength={80}
              />
            </div>
            <div className="field">
              <label htmlFor="ej_desc">Descrição</label>
              <InputText
                id="ej_desc"
                value={item.descricao ?? ''}
                onChange={(e) => setItem((p) => ({ ...p, descricao: e.target.value }))}
                className={classNamesInvalid(submitted && !String(item.descricao ?? '').trim())}
              />
            </div>
            <div className="field">
              <label htmlFor="ej_li">Limite inicial (R$)</label>
              <InputNumber
                inputId="ej_li"
                value={item.limite_inicial ?? null}
                mode="currency"
                currency="BRL"
                locale="pt-BR"
                minFractionDigits={2}
                maxFractionDigits={2}
                onValueChange={(e) => setItem((p) => ({ ...p, limite_inicial: e.value ?? 0 }))}
                className="w-full"
              />
            </div>
            <div className="field-checkbox mb-3">
              <input
                id="ej_sem_teto"
                type="checkbox"
                checked={semLimiteSuperior}
                onChange={(e) => {
                  const c = e.currentTarget.checked;
                  setSemLimiteSuperior(c);
                  if (c) {
                    setItem((p) => ({ ...p, limite_final: null }));
                  }
                }}
              />
              <label htmlFor="ej_sem_teto" className="ml-2">
                Sem limite superior (faixa aberta)
              </label>
            </div>
            {!semLimiteSuperior ? (
              <div className="field">
                <label htmlFor="ej_lf">Limite final (R$)</label>
                <InputNumber
                  inputId="ej_lf"
                  value={item.limite_final ?? null}
                  mode="currency"
                  currency="BRL"
                  locale="pt-BR"
                  minFractionDigits={2}
                  maxFractionDigits={2}
                  onValueChange={(e) => setItem((p) => ({ ...p, limite_final: e.value ?? null }))}
                  className="w-full"
                />
              </div>
            ) : null}
            <div className="field">
              <label htmlFor="ej_ano">Ano de vigência</label>
              <InputNumber
                inputId="ej_ano"
                value={item.ano_vigencia ?? null}
                useGrouping={false}
                min={1995}
                max={2100}
                onValueChange={(e) => setItem((p) => ({ ...p, ano_vigencia: e.value ?? anoAtual }))}
                className={classNamesInvalid(submitted && (!item.ano_vigencia || item.ano_vigencia < 1995))}
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
                <Button label="Não" text type="button" onClick={() => setDeleteVisible(false)} />
                <Button label="Sim" text type="button" onClick={() => void remove()} />
              </div>
            }
          >
            <span>
              Excluir <strong>{item.sigla}</strong> ({item.ano_vigencia})?
            </span>
          </Dialog>
        </div>
      </div>
    </div>
  );
};

function classNamesInvalid(invalid: boolean): string {
  return invalid ? 'p-invalid w-full' : 'w-full';
}

export default EnquadramentoJuridicoPage;
