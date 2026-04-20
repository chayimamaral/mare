-- Desfazer a migration anterior que tratava o broadcast pelo public
DROP TABLE IF EXISTS public.broadcast_messages;

-- Criação da tabela de Caixa Postal no schema matricial (template_tenant)
CREATE TABLE IF NOT EXISTS template_tenant.caixa_postal_mensagens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    remetente_id uuid,
    remetente_nome text NOT NULL,
    tipo text NOT NULL DEFAULT 'INBOX', -- INBOX: recebido pelo tenant; OUTBOX: enviado por usuario do tenant
    is_global boolean NOT NULL DEFAULT false,
    titulo text NOT NULL,
    conteudo text NOT NULL,
    lida boolean NOT NULL DEFAULT false,
    lida_por uuid REFERENCES public.usuario(id) ON DELETE SET NULL,
    lida_em timestamptz,
    criado_em timestamptz NOT NULL DEFAULT now()
);

-- Espelhar em todos os schemas (tenants) ja vivos no banco via tenant_schema_catalog
DO $$
DECLARE
    schema_record RECORD;
BEGIN
    FOR schema_record IN
        SELECT DISTINCT schema_name FROM public.tenant_schema_catalog WHERE TRIM(schema_name) <> ''
    LOOP
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I.caixa_postal_mensagens (LIKE template_tenant.caixa_postal_mensagens INCLUDING ALL);',
            schema_record.schema_name
        );
        BEGIN
            EXECUTE format(
                'ALTER TABLE %I.caixa_postal_mensagens ADD CONSTRAINT caixa_postal_mensagens_lida_por_fkey FOREIGN KEY (lida_por) REFERENCES public.usuario(id) ON DELETE SET NULL;',
                schema_record.schema_name
            );
        EXCEPTION WHEN duplicate_object THEN
            -- constraint ja existe, ignora
        END;
    END LOOP;
END;
$$;

