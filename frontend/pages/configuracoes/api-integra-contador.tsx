import { Button } from 'primereact/button';
import { InputText } from 'primereact/inputtext';
import { Panel } from 'primereact/panel';
import { Toast } from 'primereact/toast';
import { useQuery } from '@tanstack/react-query';
import { useEffect, useRef, useState } from 'react';
import api from '../../components/api/apiClient';
import { useRouteClientGuard } from '../../components/hooks/useClientGuards';

type Chaves = {
    consumer_key: string;
    consumer_secret: string;
};

export default function ApiIntegraContadorPage() {
    useRouteClientGuard();

    const toast = useRef<Toast>(null);
    const [form, setForm] = useState<Chaves>({ consumer_key: '', consumer_secret: '' });

    const { data, refetch, isFetching } = useQuery({
        queryKey: ['integra-contador-chave-autenticacao'],
        queryFn: async () => {
            const { data } = await api.get('/api/chavessuper');
            const ch = data?.chaves ?? {};
            return {
                consumer_key: ch.consumer_key ?? '',
                consumer_secret: ch.consumer_secret ?? '',
            };
        },
    });

    useEffect(() => {
        if (!data) return;
        setForm(data);
    }, [data]);

    const save = async () => {
        if (!form.consumer_key.trim() || !form.consumer_secret.trim()) {
            toast.current?.show({ severity: 'warn', summary: 'Atenção', detail: 'Preencha Consumer Key e Consumer Secret', life: 3500 });
            return;
        }
        try {
            await api.put('/api/chavessuper', form);
            toast.current?.show({ severity: 'success', summary: 'Sucesso', detail: 'Chaves salvas', life: 3000 });
            await refetch();
        } catch (e: unknown) {
            const ax = e as { response?: { data?: { error?: string; message?: string } } };
            const detail =
                ax?.response?.data?.error || ax?.response?.data?.message || 'Falha ao salvar';
            toast.current?.show({ severity: 'error', summary: 'Erro', detail, life: 4000 });
        }
    };

    const semChaves =
        data && !String(data.consumer_key).trim() && !String(data.consumer_secret).trim();

    return (
        <div className="grid">
            <div className="col-12">
                <Toast ref={toast} />
                <Panel header="Chave de Autenticação">
                    <p className="text-600 mb-4">
                        Credencial OAuth da API Integra Contador (Serpro), gravada no tenant da VEC Sistemas (o mesmo{' '}
                        <code>tenant_id</code> de todos os usuários SUPER). Não confundir com certificado digital do escritório.
                    </p>
                    <div className="field">
                        <label htmlFor="consumer_key">Consumer Key</label>
                        <InputText
                            id="consumer_key"
                            className="w-full"
                            value={form.consumer_key}
                            onChange={(e) => setForm((prev) => ({ ...prev, consumer_key: e.target.value }))}
                            autoComplete="off"
                        />
                    </div>
                    <div className="field mt-3">
                        <label htmlFor="consumer_secret">Consumer Secret</label>
                        <InputText
                            id="consumer_secret"
                            type="password"
                            className="w-full"
                            value={form.consumer_secret}
                            onChange={(e) => setForm((prev) => ({ ...prev, consumer_secret: e.target.value }))}
                            autoComplete="new-password"
                        />
                    </div>
                    <div className="mt-4 flex align-items-center gap-2 flex-wrap">
                        <Button
                            type="button"
                            label="Salvar"
                            icon="pi pi-save"
                            onClick={() => void save()}
                            disabled={isFetching}
                            loading={isFetching}
                        />
                    </div>
                    {semChaves && (
                        <small className="block mt-3 text-500">Nenhuma chave cadastrada ainda.</small>
                    )}
                </Panel>
            </div>
        </div>
    );
}
