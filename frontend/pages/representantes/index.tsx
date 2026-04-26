import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { DataTable } from 'primereact/datatable';
import { Dialog } from 'primereact/dialog';
import { InputText } from 'primereact/inputtext';
import { Toast } from 'primereact/toast';
import { Toolbar } from 'primereact/toolbar';
import { classNames } from 'primereact/utils';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { AxiosError } from 'axios';
import { useEffect, useRef, useState } from 'react';
import api from '../../components/api/apiClient';

type RepresentanteRow = {
  id: string;
  nome: string;
  email_contato?: string;
  ativo: boolean;
};

type MatrizItem = {
  modulo_id: string;
  slug: string;
  habilitado: boolean;
};

type ModuloMeta = {
  id: string;
  slug: string;
  nome: string;
  ordem: number;
};

const apiErr = (err: unknown) =>
  (err as AxiosError<{ error?: string }>)?.response?.data?.error ||
  (err as Error)?.message ||
  'Operação não concluída.';

export default function RepresentantesPage() {
  const qc = useQueryClient();
  const toast = useRef<Toast>(null);

  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [draft, setDraft] = useState<RepresentanteRow>({ id: '', nome: '', email_contato: '', ativo: true });
  const [deleteTarget, setDeleteTarget] = useState<RepresentanteRow | null>(null);

  const [matrizRepId, setMatrizRepId] = useState<string | null>(null);
  const [matrizNome, setMatrizNome] = useState('');
  const [matrizLocal, setMatrizLocal] = useState<MatrizItem[] | null>(null);

  const { data: lista = [], isFetching, refetch } = useQuery({
    queryKey: ['representantes-lista'],
    queryFn: async () => {
      const { data } = await api.get<RepresentanteRow[]>('/api/representantes');
      return Array.isArray(data) ? data : [];
    },
  });

  const { data: modulosMeta = [] } = useQuery({
    queryKey: ['modulos-plataforma'],
    queryFn: async () => {
      const { data } = await api.get<ModuloMeta[]>('/api/modulos-plataforma');
      return Array.isArray(data) ? data : [];
    },
  });

  const { data: matrizApi, isFetching: matrizLoading } = useQuery({
    queryKey: ['matriz-acesso', matrizRepId],
    enabled: Boolean(matrizRepId),
    queryFn: async () => {
      const { data } = await api.get<MatrizItem[]>('/api/matriz-acesso', {
        params: { representante_id: matrizRepId },
      });
      return Array.isArray(data) ? data : [];
    },
  });

  useEffect(() => {
    if (!matrizRepId || !matrizApi) {
      return;
    }
    setMatrizLocal(matrizApi.map((r) => ({ ...r })));
  }, [matrizRepId, matrizApi]);

  const saveMatrizMutation = useMutation({
    mutationFn: async (payload: { representante_id: string; itens: MatrizItem[] }) => {
      await api.put('/api/matriz-acesso', payload);
    },
    onSuccess: () => {
      toast.current?.show({ severity: 'success', summary: 'Matriz salva', detail: 'Permissões atualizadas.', life: 3500 });
      void qc.invalidateQueries({ queryKey: ['matriz-acesso'] });
      setMatrizRepId(null);
      setMatrizLocal(null);
    },
    onError: (e) => {
      toast.current?.show({ severity: 'error', summary: 'Erro', detail: apiErr(e), life: 5000 });
    },
  });

  const openNew = () => {
    setDraft({ id: '', nome: '', email_contato: '', ativo: true });
    setSubmitted(false);
    setEditOpen(true);
  };

  const openEdit = (row: RepresentanteRow) => {
    setDraft({
      id: row.id,
      nome: row.nome,
      email_contato: row.email_contato ?? '',
      ativo: row.ativo,
    });
    setSubmitted(false);
    setEditOpen(true);
  };

  const openMatriz = (row: RepresentanteRow) => {
    setMatrizNome(row.nome);
    setMatrizRepId(row.id);
    setMatrizLocal(null);
  };

  const closeMatriz = () => {
    setMatrizRepId(null);
    setMatrizLocal(null);
    setMatrizNome('');
  };

  const toggleModulo = (moduloId: string) => {
    setMatrizLocal((prev) => {
      if (!prev) {
        return prev;
      }
      return prev.map((r) => (r.modulo_id === moduloId ? { ...r, habilitado: !r.habilitado } : r));
    });
  };

  const salvarMatriz = () => {
    if (!matrizRepId || !matrizLocal) {
      return;
    }
    saveMatrizMutation.mutate({
      representante_id: matrizRepId,
      itens: matrizLocal.map((r) => ({
        modulo_id: r.modulo_id,
        slug: r.slug,
        habilitado: r.habilitado,
      })),
    });
  };

  const salvarRepresentante = async () => {
    setSubmitted(true);
    if (!draft.nome.trim()) {
      toast.current?.show({ severity: 'warn', summary: 'Validação', detail: 'Informe o nome.', life: 3500 });
      return;
    }
    try {
      if (draft.id) {
        await api.put('/api/representantes', {
          id: draft.id,
          nome: draft.nome.trim(),
          email_contato: (draft.email_contato ?? '').trim(),
          ativo: draft.ativo,
        });
        toast.current?.show({ severity: 'success', summary: 'Atualizado', detail: 'Representante salvo.', life: 3000 });
      } else {
        await api.post('/api/representantes', {
          nome: draft.nome.trim(),
          email_contato: (draft.email_contato ?? '').trim(),
        });
        toast.current?.show({ severity: 'success', summary: 'Criado', detail: 'Representante cadastrado.', life: 3000 });
      }
      setEditOpen(false);
      await refetch();
    } catch (e) {
      toast.current?.show({ severity: 'error', summary: 'Erro', detail: apiErr(e), life: 5000 });
    }
  };

  const excluirRepresentante = async () => {
    if (!deleteTarget?.id) {
      return;
    }
    try {
      await api.delete('/api/representantes', { params: { id: deleteTarget.id } });
      toast.current?.show({ severity: 'success', summary: 'Removido', detail: 'Representante excluído.', life: 3000 });
      setDeleteOpen(false);
      setDeleteTarget(null);
      await refetch();
    } catch (e) {
      toast.current?.show({ severity: 'error', summary: 'Erro', detail: apiErr(e), life: 5000 });
    }
  };

  const colAcoes = (row: RepresentanteRow) => (
    <div className="flex flex-wrap gap-2 justify-content-end">
      <Button type="button" icon="pi pi-th-large" rounded severity="help" tooltip="Matriz de módulos" onClick={() => openMatriz(row)} />
      <Button type="button" icon="pi pi-pencil" rounded severity="success" tooltip="Editar" onClick={() => openEdit(row)} />
      <Button
        type="button"
        icon="pi pi-trash"
        rounded
        severity="warning"
        tooltip="Excluir"
        onClick={() => {
          setDeleteTarget(row);
          setDeleteOpen(true);
        }}
      />
    </div>
  );

  const colAtivo = (row: RepresentanteRow) => (row.ativo ? 'Sim' : 'Não');

  const matrizDialogVisible = Boolean(matrizRepId);

  const tituloModulo = (moduloId: string, slug: string) => {
    const meta = modulosMeta.find((m) => m.id === moduloId);
    return meta?.nome?.trim() ? meta.nome : slug.replace(/_/g, ' ');
  };

  return (
    <div className="grid crud-demo">
      <div className="col-12">
        <div className="card">
          <Toast ref={toast} />
          <Toolbar
            className="mb-4"
            left={
              <div className="my-2 flex flex-wrap align-items-center gap-2">
                <Button type="button" label="Novo representante" icon="pi pi-plus" severity="success" onClick={openNew} />
              </div>
            }
          />
          <h2 className="mt-0 mb-3 text-xl font-semibold">Representantes comerciais</h2>
          <p className="text-600 mb-4 text-sm">
            Cadastro e matriz de acesso por módulo. O que estiver marcado vale para o representante e para os tenants vinculados a ele.
          </p>
          <DataTable value={lista} loading={isFetching} dataKey="id" emptyMessage="Nenhum representante cadastrado." tableStyle={{ minWidth: '42rem' }}>
            <Column field="nome" header="Nome" sortable style={{ minWidth: '14rem' }} />
            <Column field="email_contato" header="Contato (e-mail)" style={{ minWidth: '12rem' }} />
            <Column header="Ativo" body={colAtivo} style={{ width: '6rem' }} />
            <Column header="Ações" body={colAcoes} style={{ minWidth: '12rem' }} />
          </DataTable>

          <Dialog
            visible={editOpen}
            onHide={() => setEditOpen(false)}
            header={draft.id ? 'Editar representante' : 'Novo representante'}
            modal
            className="p-fluid"
            style={{ width: 'min(32rem, 96vw)' }}
            footer={
              <>
                <Button type="button" label="Cancelar" icon="pi pi-times" text onClick={() => setEditOpen(false)} />
                <Button type="button" label="Salvar" icon="pi pi-check" text onClick={() => void salvarRepresentante()} />
              </>
            }
          >
            <div className="field">
              <label htmlFor="repNome">Nome</label>
              <InputText
                id="repNome"
                value={draft.nome}
                onChange={(e) => setDraft((p) => ({ ...p, nome: e.target.value }))}
                className={classNames({ 'p-invalid': submitted && !draft.nome.trim() })}
              />
              {submitted && !draft.nome.trim() && <small className="p-invalid">Obrigatório.</small>}
            </div>
            <div className="field">
              <label htmlFor="repEmail">E-mail de contato</label>
              <InputText
                id="repEmail"
                value={draft.email_contato ?? ''}
                onChange={(e) => setDraft((p) => ({ ...p, email_contato: e.target.value }))}
              />
            </div>
            {draft.id ? (
              <div className="field field-checkbox">
                <input
                  type="checkbox"
                  id="repAtivo"
                  checked={draft.ativo}
                  onChange={(e) => setDraft((p) => ({ ...p, ativo: e.target.checked }))}
                />
                <label htmlFor="repAtivo">Ativo</label>
              </div>
            ) : null}
          </Dialog>

          <Dialog
            visible={matrizDialogVisible}
            onHide={closeMatriz}
            header={`Matriz de módulos — ${matrizNome}`}
            modal
            className="p-fluid"
            style={{ width: 'min(40rem, 96vw)' }}
            contentStyle={{ minHeight: 'min(58vh, 38rem)' }}
            footer={
              <>
                <Button type="button" label="Cancelar" icon="pi pi-times" text onClick={closeMatriz} />
                <Button
                  type="button"
                  label="Salvar matriz"
                  icon="pi pi-check"
                  text
                  loading={saveMatrizMutation.isPending}
                  disabled={!matrizLocal?.length}
                  onClick={() => salvarMatriz()}
                />
              </>
            }
          >
            {matrizLoading ? (
              <p className="text-600 text-sm">Carregando módulos…</p>
            ) : !matrizLocal?.length ? (
              <p className="text-600 text-sm">Nenhum módulo encontrado na plataforma.</p>
            ) : (
              <ul className="list-none p-0 m-0 flex flex-column gap-3">
                {matrizLocal.map((r) => {
                  const labelId = `mod-${r.modulo_id}`;
                  const titulo = tituloModulo(r.modulo_id, r.slug);
                  return (
                    <li key={r.modulo_id} className="flex align-items-center gap-3">
                      <div className="field-checkbox m-0">
                        <input
                          type="checkbox"
                          id={labelId}
                          checked={r.habilitado}
                          onChange={() => toggleModulo(r.modulo_id)}
                        />
                        <label htmlFor={labelId} className="mb-0 cursor-pointer">
                          <span className="font-medium">{titulo}</span>
                          <span className="text-600 text-sm block">{r.slug}</span>
                        </label>
                      </div>
                    </li>
                  );
                })}
              </ul>
            )}
          </Dialog>

          <Dialog
            visible={deleteOpen}
            onHide={() => setDeleteOpen(false)}
            header="Confirmar exclusão"
            modal
            footer={
              <>
                <Button type="button" label="Não" icon="pi pi-times" text onClick={() => setDeleteOpen(false)} />
                <Button type="button" label="Sim, excluir" icon="pi pi-check" text severity="warning" onClick={() => void excluirRepresentante()} />
              </>
            }
          >
            {deleteTarget ? (
              <span>
                Excluir o representante <b>{deleteTarget.nome}</b>? Tenants vinculados podem impedir a exclusão.
              </span>
            ) : null}
          </Dialog>

          <div className="fixed" style={{ bottom: '1.5rem', right: '1.5rem', zIndex: 100 }}>
            <Button type="button" icon="pi pi-refresh" rounded tooltip="Atualizar" aria-label="Atualizar lista" onClick={() => void refetch()} />
          </div>
        </div>
      </div>
    </div>
  );
}
