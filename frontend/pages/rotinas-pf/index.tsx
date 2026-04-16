import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { DataTable, DataTableFilterMeta } from 'primereact/datatable';
import { Dialog } from 'primereact/dialog';
import { Dropdown } from 'primereact/dropdown';
import { InputNumber } from 'primereact/inputnumber';
import { InputSwitch } from 'primereact/inputswitch';
import { InputText } from 'primereact/inputtext';
import { InputTextarea } from 'primereact/inputtextarea';
import { Toast } from 'primereact/toast';
import { Toolbar } from 'primereact/toolbar';
import React, { useEffect, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import RotinaPFService from '../../services/cruds/RotinaPFService';
import PassoService from '../../services/cruds/PassoService';
import { Vec } from '../../types/types';

interface LazyTableState {
  totalRecords: number;
  first: number;
  rows: number;
  page: number;
  sortField?: string;
  sortOrder?: number;
  filters: DataTableFilterMeta;
}

const emptyRow: Vec.RotinaPFListRow = {
  id: '',
  nome: '',
  categoria: 'MENSALISTA',
  descricao: '',
  ativo: true,
  criado_em: '',
  item_count: 0,
};

const emptyItem: Vec.RotinaPFItemRow = {
  id: '',
  rotina_pf_id: '',
  ordem: 0,
  passo_id: '',
  passo_descricao: '',
  descricao: '',
  tempo_estimado: 0,
};

const categorias = [
  { label: 'Mensalista (ex.: Carnê-Leão)', value: 'MENSALISTA' },
  { label: 'Sazonal IRPF', value: 'SAZONAL_IRPF' },
  { label: 'Avulso', value: 'AVULSO' },
];

const RotinasPF = () => {
  const toast = useRef<Toast>(null);

  const [lazyState, setLazyState] = useState<LazyTableState>({
    totalRecords: 0,
    first: 0,
    rows: 20,
    page: 1,
    sortField: 'nome',
    sortOrder: 1,
    filters: { nome: { value: '', matchMode: 'contains' } },
  });

  const [dialogVisible, setDialogVisible] = useState(false);
  const [deleteVisible, setDeleteVisible] = useState(false);
  const [editing, setEditing] = useState<Vec.RotinaPFListRow>(emptyRow);
  const [itens, setItens] = useState<Vec.RotinaPFItemRow[]>([]);
  const [itensLoading, setItensLoading] = useState(false);

  const [itemDialog, setItemDialog] = useState(false);
  const [editingItem, setEditingItem] = useState<Vec.RotinaPFItemRow>(emptyItem);
  const [passosFallback, setPassosFallback] = useState<Vec.Passo[]>([]);

  const svc = RotinaPFService();

  const { data: userRole = null } = useQuery<string | null>({
    queryKey: ['user-role'],
    queryFn: async () => {
      try {
        const r = await api.get('/api/usuariorole');
        return r.data?.logado?.role ?? null;
      } catch {
        return null;
      }
    },
  });

  const podeAdmin = userRole === 'ADMIN' || userRole === 'SUPER';

  const loadList = async () => {
    const body = { ...lazyState };
    const { data } = await svc.getRotinasPFAdmin({ lazyEvent: JSON.stringify(body) });
    return {
      rows: data?.rotinas_pf ?? [],
      totalRecords: data?.totalRecords ?? 0,
    };
  };

  const { data, isFetching, refetch } = useQuery({
    queryKey: ['rotinas-pf-admin', lazyState],
    queryFn: () => loadList(),
  });

  const { data: passos = passosFallback } = useQuery<Vec.Passo[]>({
    queryKey: ['passos-lite-rotinas-pf'],
    queryFn: async () => {
      try {
        const ps = PassoService();
        const st = {
          totalRecords: 0,
          first: 0,
          rows: 400,
          page: 1,
          sortField: 'descricao',
          sortOrder: 1,
          filters: { descricao: { value: '', matchMode: 'contains' } },
        };
        const { data } = await ps.getPassos({ lazyEvent: JSON.stringify(st) });
        const lista = Array.isArray(data?.passos) ? data.passos : [];
        setPassosFallback(lista);
        return lista;
      } catch {
        setPassosFallback([]);
        return [];
      }
    },
  });

  const paginatorLeft = (
    <Button type="button" icon="pi pi-refresh" tooltip="Atualizar" className="p-button-text" onClick={() => refetch()} />
  );

  const onPage = (event: any) => {
    const nextRows = typeof event.rows === 'number' && event.rows > 0 ? event.rows : lazyState.rows;
    setLazyState((prev) => ({
      ...prev,
      first: event.first,
      rows: nextRows,
      page: event.page + 1,
      sortField: event.sortField ?? prev.sortField,
      sortOrder: event.sortOrder ?? prev.sortOrder,
    }));
  };

  const onSort = (event: any) => {
    setLazyState((prev) => ({
      ...prev,
      sortField: event.sortField ?? prev.sortField,
      sortOrder: event.sortOrder ?? prev.sortOrder,
      first: 0,
      page: 1,
    }));
  };

  const loadItens = (rotinaId: string) => {
    if (!rotinaId) {
      setItens([]);
      return;
    }
    setItensLoading(true);
    svc
      .getItens(rotinaId)
      .then(({ data }) => setItens(data.itens ?? []))
      .catch(() => {
        setItens([]);
        toast.current?.show({ severity: 'error', summary: 'Erro', detail: 'Falha ao carregar itens', life: 3500 });
      })
      .finally(() => setItensLoading(false));
  };

  const openNew = () => {
    setEditing({ ...emptyRow, categoria: 'MENSALISTA', ativo: true });
    setItens([]);
    setDialogVisible(true);
  };

  const openEdit = (row: Vec.RotinaPFListRow) => {
    setEditing({
      ...emptyRow,
      ...row,
      ativo: row.ativo !== false,
    });
    loadItens(row.id ?? '');
    setDialogVisible(true);
  };

  const saveCabecalho = () => {
    const nome = (editing.nome ?? '').trim();
    if (!nome) {
      toast.current?.show({ severity: 'warn', summary: 'Validação', detail: 'Informe o nome', life: 3000 });
      return;
    }
    const cat = (editing.categoria ?? '').trim();
    if (!cat) {
      toast.current?.show({ severity: 'warn', summary: 'Validação', detail: 'Selecione a categoria', life: 3000 });
      return;
    }
    const p = {
      nome,
      categoria: cat,
      descricao: editing.descricao ?? '',
      ativo: editing.ativo !== false,
    };
    const done = () => {
      toast.current?.show({ severity: 'success', summary: 'Salvo', detail: 'Rotina PF atualizada', life: 2500 });
      refetch();
      setDialogVisible(false);
    };
    if (editing.id) {
      svc
        .updateRotinaPF({ id: editing.id, ...p })
        .then(done)
        .catch((e) => {
          const msg = e?.response?.data?.error ?? 'Erro ao salvar';
          toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 5000 });
        });
    } else {
      svc
        .createRotinaPF(p)
        .then((res) => {
          const nid = res?.data?.rotinas_pf?.[0]?.id;
          if (nid) {
            setEditing((prev) => ({ ...prev, id: nid }));
            loadItens(nid);
          }
          toast.current?.show({ severity: 'success', summary: 'Criado', detail: 'Inclua os passos abaixo', life: 3500 });
          refetch();
        })
        .catch((e) => {
          const msg = e?.response?.data?.error ?? 'Erro ao criar';
          toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 5000 });
        });
    }
  };

  const confirmDelete = (row: Vec.RotinaPFListRow) => {
    setEditing(row);
    setDeleteVisible(true);
  };

  const doDelete = () => {
    const id = editing.id;
    if (!id) return;
    svc
      .deleteRotinaPF(id)
      .then(() => {
        toast.current?.show({ severity: 'success', summary: 'Desativada', detail: 'A rotina não aparece mais no cadastro de clientes PF', life: 4000 });
        setDeleteVisible(false);
        refetch();
      })
      .catch((e) => {
        const msg = e?.response?.data?.error ?? 'Erro';
        toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 5000 });
      });
  };

  const nextOrdem = () => {
    if (!itens.length) return 0;
    return Math.max(...itens.map((i) => Number(i.ordem) || 0)) + 1;
  };

  const openNewItem = () => {
    const rid = editing.id;
    if (!rid) {
      toast.current?.show({ severity: 'warn', summary: 'Ordem', detail: 'Salve o cabeçalho da rotina antes de incluir passos', life: 4000 });
      return;
    }
    setEditingItem({
      ...emptyItem,
      rotina_pf_id: rid,
      ordem: nextOrdem(),
      tempo_estimado: 0,
    });
    setItemDialog(true);
  };

  const openEditItem = (it: Vec.RotinaPFItemRow) => {
    setEditingItem({ ...emptyItem, ...it });
    setItemDialog(true);
  };

  const saveItem = () => {
    const rid = editing.id;
    if (!rid) return;
    const passo = (editingItem.passo_id ?? '').trim();
    const desc = (editingItem.descricao ?? '').trim();
    if (!passo && !desc) {
      toast.current?.show({ severity: 'warn', summary: 'Validação', detail: 'Preencha passo ou descrição livre', life: 3500 });
      return;
    }
    const base = {
      rotina_pf_id: rid,
      ordem: Number(editingItem.ordem) || 0,
      passo_id: passo,
      descricao: editingItem.descricao ?? '',
      tempo_estimado: Number(editingItem.tempo_estimado) || 0,
    };
    const ok = () => {
      loadItens(rid);
      setItemDialog(false);
      loadList();
    };
    if (editingItem.id) {
      svc
        .updateItem({ item_id: editingItem.id, ...base })
        .then(ok)
        .catch((e) => {
          const msg = e?.response?.data?.error ?? 'Erro ao salvar item';
          toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 5000 });
        });
    } else {
      svc
        .createItem(base)
        .then(ok)
        .catch((e) => {
          const msg = e?.response?.data?.error ?? 'Erro ao criar item';
          toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 5000 });
        });
    }
  };

  const removeItem = (it: Vec.RotinaPFItemRow) => {
    const rid = editing.id;
    if (!it.id || !rid) return;
    svc
      .deleteItem(it.id, rid)
      .then(() => {
        loadItens(rid);
        refetch();
      })
      .catch((e) => {
        const msg = e?.response?.data?.error ?? 'Erro ao excluir';
        toast.current?.show({ severity: 'error', summary: 'Erro', detail: String(msg), life: 5000 });
      });
  };

  const passoDropdownValue = (() => {
    const id = (editingItem.passo_id ?? '').trim();
    if (!id) return null;
    return passos.find((p) => (p.id ?? '').trim() === id) ?? null;
  })();

  const acoesCorpo = (row: Vec.RotinaPFListRow) => (
    <>
      <Button icon="pi pi-pencil" rounded text className="mr-1" onClick={() => openEdit(row)} tooltip="Editar" />
      {podeAdmin && (
        <Button
          icon="pi pi-ban"
          rounded
          text
          severity="danger"
          className="mr-1"
          onClick={() => confirmDelete(row)}
          tooltip="Desativar"
        />
      )}
    </>
  );

  const acoesItem = (row: Vec.RotinaPFItemRow) =>
    podeAdmin ? (
      <>
        <Button icon="pi pi-pencil" rounded text className="mr-1" onClick={() => openEditItem(row)} tooltip="Editar" />
        <Button icon="pi pi-trash" rounded text severity="danger" onClick={() => removeItem(row)} tooltip="Excluir" />
      </>
    ) : null;

  const ativoTemplate = (row: Vec.RotinaPFListRow) => (row.ativo === false ? 'Não' : 'Sim');

  const header = (
    <div className="flex flex-wrap gap-2 align-items-center justify-content-between">
      <span className="p-input-icon-left w-full sm:w-auto">
        <i className="pi pi-search" />
        <InputText
          value={(lazyState.filters?.nome as { value?: string })?.value ?? ''}
          onChange={(e) =>
            setLazyState((prev) => ({
              ...prev,
              first: 0,
              page: 1,
              filters: { ...prev.filters, nome: { value: e.target.value, matchMode: 'contains' } },
            }))
          }
          placeholder="Buscar por nome"
          className="w-full"
        />
      </span>
    </div>
  );

  return (
    <div className="grid crud-demo">
      <div className="col-12">
        <div className="card">
          <Toast ref={toast} />
          <Toolbar
            className="mb-4"
            start={
              podeAdmin ? (
                <Button label="Nova rotina PF" icon="pi pi-plus" severity="success" onClick={openNew} />
              ) : (
                <span className="text-600 text-sm">Somente administradores criam ou alteram templates. Usuários usam a lista no cadastro de clientes PF.</span>
              )
            }
          />
          <p className="text-600 text-sm mt-0 mb-3">
            Templates por tenant para clientes pessoa física (Carnê-Leão, IRPF, etc.). Itens podem referenciar passos globais ou só texto livre.
          </p>
          <DataTable
            value={data?.rows ?? []}
            lazy
            dataKey="id"
            paginator
            rows={lazyState.rows}
            totalRecords={data?.totalRecords ?? 0}
            first={lazyState.first}
            onPage={onPage}
            onSort={onSort}
            sortField={lazyState.sortField}
            sortOrder={lazyState.sortOrder === -1 ? -1 : 1}
            loading={isFetching}
            emptyMessage="Nenhuma rotina PF cadastrada."
            header={header}
            paginatorLeft={paginatorLeft}
            rowsPerPageOptions={[10, 20, 50]}
          >
            <Column field="nome" header="Nome" sortable style={{ minWidth: '14rem' }} />
            <Column field="categoria" header="Categoria" sortable style={{ minWidth: '10rem' }} />
            <Column field="ativo" header="Ativa" body={ativoTemplate} style={{ minWidth: '6rem' }} />
            <Column field="item_count" header="Itens" style={{ minWidth: '6rem' }} />
            <Column header="Ações" body={acoesCorpo} style={{ minWidth: '8rem' }} />
          </DataTable>

          <Dialog
            header={editing.id ? 'Rotina PF (edição)' : 'Rotina PF (nova)'}
            visible={dialogVisible}
            style={{ width: 'min(800px, 96vw)' }}
            onHide={() => setDialogVisible(false)}
            modal
            footer={
              <div className="flex justify-content-end gap-2">
                <Button label="Cancelar" icon="pi pi-times" text onClick={() => setDialogVisible(false)} />
                {podeAdmin && <Button label="Salvar cabeçalho" icon="pi pi-check" onClick={saveCabecalho} />}
              </div>
            }
          >
            <div className="field">
              <label htmlFor="rpf_nome">Nome</label>
              <InputText
                id="rpf_nome"
                value={editing.nome ?? ''}
                onChange={(e) => setEditing((p) => ({ ...p, nome: e.target.value }))}
                className="w-full"
                disabled={!podeAdmin}
              />
            </div>
            <div className="field">
              <label htmlFor="rpf_cat">Categoria</label>
              <Dropdown
                id="rpf_cat"
                value={editing.categoria}
                options={categorias}
                onChange={(e) => setEditing((p) => ({ ...p, categoria: e.value }))}
                className="w-full"
                disabled={!podeAdmin}
              />
            </div>
            <div className="field">
              <label htmlFor="rpf_desc">Descrição (opcional)</label>
              <InputTextarea
                id="rpf_desc"
                value={editing.descricao ?? ''}
                onChange={(e) => setEditing((p) => ({ ...p, descricao: e.target.value }))}
                rows={2}
                className="w-full"
                disabled={!podeAdmin}
              />
            </div>
            <div className="field flex align-items-center gap-2">
              <InputSwitch
                checked={editing.ativo !== false}
                onChange={(e) => setEditing((p) => ({ ...p, ativo: !!e.value }))}
                disabled={!podeAdmin}
              />
              <span>Ativa (aparece no cadastro de clientes PF)</span>
            </div>

            <h3 className="text-lg mt-4 mb-2">Passos da agenda</h3>
            {podeAdmin && (
              <Button label="Incluir passo" icon="pi pi-plus" className="mb-3" type="button" onClick={openNewItem} />
            )}
            <DataTable value={itens} loading={itensLoading} dataKey="id" emptyMessage="Nenhum passo.">
              <Column field="ordem" header="Ordem" style={{ width: '6rem' }} />
              <Column field="passo_descricao" header="Passo (catálogo)" />
              <Column field="descricao" header="Texto / complemento" />
              <Column field="tempo_estimado" header="Tempo (dias úteis)" style={{ width: '8rem' }} />
              <Column header="Ações" body={acoesItem} style={{ width: '8rem' }} />
            </DataTable>
          </Dialog>

          <Dialog
            header={editingItem.id ? 'Passo (edição)' : 'Passo (novo)'}
            visible={itemDialog}
            style={{ width: 'min(520px, 94vw)' }}
            onHide={() => setItemDialog(false)}
            modal
            footer={
              <div className="flex justify-content-end gap-2">
                <Button label="Cancelar" icon="pi pi-times" text onClick={() => setItemDialog(false)} />
                {podeAdmin && <Button label="Salvar" icon="pi pi-check" onClick={saveItem} />}
              </div>
            }
          >
            <div className="field">
              <label htmlFor="it_ordem">Ordem</label>
              <InputNumber
                id="it_ordem"
                value={editingItem.ordem}
                onValueChange={(e) => setEditingItem((p) => ({ ...p, ordem: Number(e.value) || 0 }))}
                className="w-full"
                disabled={!podeAdmin}
              />
            </div>
            <div className="field">
              <label htmlFor="it_passo">Passo do catálogo (opcional)</label>
              <Dropdown
                id="it_passo"
                value={passoDropdownValue}
                options={passos}
                onChange={(e) =>
                  setEditingItem((p) => ({
                    ...p,
                    passo_id: e.value?.id ?? '',
                    passo_descricao: e.value?.descricao ?? '',
                  }))
                }
                optionLabel="descricao"
                dataKey="id"
                filter
                showClear
                placeholder="Nenhum — usar só descrição livre"
                className="w-full"
                disabled={!podeAdmin}
              />
            </div>
            <div className="field">
              <label htmlFor="it_txt">Descrição livre (se vazia, usa a do passo)</label>
              <InputText
                id="it_txt"
                value={editingItem.descricao ?? ''}
                onChange={(e) => setEditingItem((p) => ({ ...p, descricao: e.target.value }))}
                className="w-full"
                disabled={!podeAdmin}
              />
            </div>
            <div className="field">
              <label htmlFor="it_te">Tempo estimado (dias úteis)</label>
              <InputNumber
                id="it_te"
                value={editingItem.tempo_estimado}
                onValueChange={(e) => setEditingItem((p) => ({ ...p, tempo_estimado: Number(e.value) || 0 }))}
                className="w-full"
                disabled={!podeAdmin}
              />
            </div>
          </Dialog>

          <Dialog
            header="Desativar rotina PF"
            visible={deleteVisible}
            style={{ width: '400px' }}
            modal
            onHide={() => setDeleteVisible(false)}
            footer={
              <>
                <Button label="Não" icon="pi pi-times" text onClick={() => setDeleteVisible(false)} />
                <Button label="Desativar" icon="pi pi-check" severity="danger" onClick={doDelete} />
              </>
            }
          >
            <p>
              A rotina <strong>{editing.nome}</strong> deixará de aparecer para novos vínculos. Clientes que já a usam
              permanecem com o vínculo; reative depois se precisar.
            </p>
          </Dialog>
        </div>
      </div>
    </div>
  );
};

export default RotinasPF;
