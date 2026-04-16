import { useEffect, useMemo, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Button } from 'primereact/button';
import { Card } from 'primereact/card';
import { Dropdown } from 'primereact/dropdown';
import { InputText } from 'primereact/inputtext';
import { Toast } from 'primereact/toast';
import CatalogoServicoService, { CatalogoServico } from '../../services/cruds/CatalogoServicoService';
import RegimeTributarioService from '../../services/cruds/RegimeTributarioService';
import SerproServicoEnquadramentoService from '../../services/cruds/SerproServicoEnquadramentoService';
import TipoEmpresaService from '../../services/cruds/TipoEmpresaService';

type SelectOption = { label: string; value: string };

export default function MatrizConformidadeFiscalPage() {
  const toast = useRef<Toast>(null);
  const tipoEmpresaService = useMemo(() => TipoEmpresaService(), []);
  const regimeService = useMemo(() => RegimeTributarioService(), []);
  const catalogoService = useMemo(() => CatalogoServicoService(), []);
  const matrizService = useMemo(() => SerproServicoEnquadramentoService(), []);

  const [enquadramentoID, setEnquadramentoID] = useState('');
  const [regimeID, setRegimeID] = useState('');
  const [secaoFiltro, setSecaoFiltro] = useState('');
  const [busca, setBusca] = useState('');
  const [selecionados, setSelecionados] = useState<string[]>([]);
  const [saving, setSaving] = useState(false);

  const { data: enquadramentos = [] } = useQuery<SelectOption[]>({
    queryKey: ['matriz-conformidade-enquadramentos'],
    queryFn: async () => {
      const { data } = await tipoEmpresaService.getTiposEmpresaLite();
      const rows = Array.isArray(data?.tiposEmpresa) ? data.tiposEmpresa : [];
      return rows
        .map((r: { id?: string; descricao?: string }) => ({ value: String(r.id || '').trim(), label: String(r.descricao || '').trim() }))
        .filter((x: SelectOption) => x.value !== '');
    },
  });

  const { data: regimes = [] } = useQuery<SelectOption[]>({
    queryKey: ['matriz-conformidade-regimes'],
    queryFn: async () => {
      const { data } = await regimeService.getRegimes({
        lazyEvent: JSON.stringify({
          first: 0,
          rows: 500,
          page: 1,
          sortField: 'nome',
          sortOrder: 1,
          filters: { nome: { value: '', matchMode: 'contains' } },
        }),
      });
      const rows = Array.isArray(data?.regimes) ? data.regimes : [];
      return rows
        .map((r: { id?: string; nome?: string; codigo_crt?: number }) => ({
          value: String(r.id || '').trim(),
          label: `${String(r.nome || '').trim()}${typeof r.codigo_crt === 'number' ? ` (CRT ${r.codigo_crt})` : ''}`,
        }))
        .filter((x: SelectOption) => x.value !== '');
    },
  });

  const { data: catalogo = [] } = useQuery<CatalogoServico[]>({
    queryKey: ['matriz-conformidade-catalogo'],
    queryFn: () => catalogoService.list({ incluirInativos: false }),
  });

  const { data: selecionadosApi = [], refetch: refetchSelecionados, isFetching: loadingSelecionados } = useQuery<string[]>({
    queryKey: ['matriz-conformidade-selecionados', enquadramentoID, regimeID],
    enabled: enquadramentoID.trim() !== '' && regimeID.trim() !== '',
    queryFn: async () => {
      const { data } = await matrizService.list(enquadramentoID, regimeID);
      const ids = Array.isArray(data?.servicos_ids) ? data.servicos_ids : [];
      return ids.map((x: unknown) => String(x || '').trim()).filter((x: string) => x !== '');
    },
  });

  useEffect(() => {
    // Sincroniza somente quando o backend mudar de fato a seleção
    if (!Array.isArray(selecionadosApi)) return;
    // evita sobrescrever alterações locais enquanto o usuário marca/desmarca
    const atual = selecionados.slice().sort();
    const remoto = selecionadosApi.slice().sort();
    if (atual.length === remoto.length && atual.every((v, i) => v === remoto[i])) {
      return;
    }
    setSelecionados(selecionadosApi);
  }, [selecionadosApi, selecionados]);

  const secoes = useMemo(() => {
    const unique = Array.from(new Set(catalogo.map((s) => (s.secao || '').trim()).filter((s) => s !== '')));
    return unique.sort((a, b) => a.localeCompare(b, 'pt-BR', { sensitivity: 'base' }));
  }, [catalogo]);

  const lista = useMemo(() => {
    const termo = busca.trim().toLowerCase();
    return catalogo
      .filter((s) => (secaoFiltro ? s.secao === secaoFiltro : true))
      .filter((s) => {
        if (!termo) return true;
        return [s.codigo, s.descricao, s.id_sistema, s.id_servico, s.secao]
          .map((x) => String(x || '').toLowerCase())
          .some((x) => x.includes(termo));
      })
      .sort((a, b) => {
        const secaoCmp = a.secao.localeCompare(b.secao, 'pt-BR', { sensitivity: 'base' });
        if (secaoCmp !== 0) return secaoCmp;
        if (a.sequencial !== b.sequencial) return a.sequencial - b.sequencial;
        return a.codigo.localeCompare(b.codigo, 'pt-BR', { sensitivity: 'base' });
      });
  }, [catalogo, secaoFiltro, busca]);

  const toggleServico = (id: string, checked: boolean) => {
    setSelecionados((prev) => {
      if (checked) {
        if (prev.includes(id)) return prev;
        return [...prev, id];
      }
      return prev.filter((x) => x !== id);
    });
  };

  const salvar = async () => {
    if (!enquadramentoID || !regimeID) {
      toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Selecione enquadramento e regime.', life: 3500 });
      return;
    }
    setSaving(true);
    try {
      await matrizService.save({
        enquadramento_id: enquadramentoID,
        regime_tributario_id: regimeID,
        servicos_ids: selecionados,
      });
      toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Matriz salva com sucesso.', life: 3000 });
      refetchSelecionados();
    } catch (e: any) {
      const msg = e?.response?.data?.error || e?.message || 'Falha ao salvar matriz.';
      toast.current?.show({ severity: 'error', summary: 'Erro', detail: msg, life: 5000 });
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="grid">
      <div className="col-12">
        <Toast ref={toast} />
        <Card title="Matriz de Conformidade Fiscal (SERPRO)">
          <p className="mt-0 text-600">
            Defina os serviços disponíveis por combinação de Enquadramento Jurídico e Regime Tributário.
          </p>

          <div className="grid mb-2">
            <div className="col-12 md:col-4">
              <label className="block mb-2 font-medium">Enquadramento Jurídico</label>
              <Dropdown
                value={enquadramentoID}
                options={enquadramentos}
                optionLabel="label"
                optionValue="value"
                filter
                showClear
                className="w-full"
                placeholder="Selecione"
                onChange={(e) => setEnquadramentoID(e.value ?? '')}
              />
            </div>
            <div className="col-12 md:col-4">
              <label className="block mb-2 font-medium">Regime Tributário</label>
              <Dropdown
                value={regimeID}
                options={regimes}
                optionLabel="label"
                optionValue="value"
                filter
                showClear
                className="w-full"
                placeholder="Selecione"
                onChange={(e) => setRegimeID(e.value ?? '')}
              />
            </div>
            <div className="col-12 md:col-4 flex align-items-end">
              <Button label="Salvar Matriz" icon="pi pi-save" onClick={() => void salvar()} loading={saving} />
            </div>
          </div>

          <div className="grid mb-3">
            <div className="col-12 md:col-4">
              <label className="block mb-2 font-medium">Seção</label>
              <Dropdown
                value={secaoFiltro}
                options={secoes.map((s) => ({ label: s, value: s }))}
                optionLabel="label"
                optionValue="value"
                showClear
                className="w-full"
                placeholder="Todas as seções"
                onChange={(e) => setSecaoFiltro(e.value ?? '')}
              />
            </div>
            <div className="col-12 md:col-8">
              <label className="block mb-2 font-medium">Buscar serviço</label>
              <span className="p-input-icon-left w-full">
                <i className="pi pi-search" />
                <InputText
                  value={busca}
                  onChange={(e) => setBusca(e.target.value)}
                  placeholder="Código, descrição, idSistema, idServico..."
                  className="w-full"
                />
              </span>
            </div>
          </div>

          {(!enquadramentoID || !regimeID) ? (
            <div className="p-3 border-1 border-round border-300 surface-50 text-700">
              Selecione Enquadramento e Regime para carregar a matriz.
            </div>
          ) : (
            <div className="border-1 border-300 border-round surface-50 p-2" style={{ maxHeight: '60vh', overflowY: 'auto' }}>
              {loadingSelecionados ? (
                <div className="p-3">Carregando seleção atual...</div>
              ) : (
                <div className="grid">
                  {lista.map((s) => {
                    const checked = selecionados.includes(s.id);
                    return (
                      <div key={s.id} className="col-12 md:col-6 lg:col-4">
                        <div className="p-2 border-1 border-200 border-round surface-card h-full">
                          <div className="flex align-items-start gap-2">
                            <div className="field-checkbox m-0">
                              <input
                                id={`svc-${s.id}`}
                                type="checkbox"
                                checked={checked}
                                onChange={(e) => toggleServico(s.id, e.target.checked)}
                              />
                              <label htmlFor={`svc-${s.id}`} className="cursor-pointer line-height-3">
                                <strong>{s.codigo}</strong> - {s.descricao}
                                <br />
                                <small className="text-600">
                                  {s.secao} | {s.id_sistema}/{s.id_servico}
                                </small>
                              </label>
                            </div>
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </div>
          )}
        </Card>
      </div>
    </div>
  );
}
