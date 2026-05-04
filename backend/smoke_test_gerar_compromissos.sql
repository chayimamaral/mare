-- =============================================================================
-- Diagnostico EF-916 / compromissos (leitura apenas — nao altera dados)
-- =============================================================================
-- Modelo atual:
--   - Negocio do escritorio: tabelas no SCHEMA do tenant (ex.: escritorio_xyz).
--   - Global em public: tenant, usuario, tenant_schema_catalog, tipoempresa_obrigacao,
--     municipio, feriados, CNAE, etc.
--   - Geracao de compromissos: API Go POST .../empresacompromissos/gerar ou worker
--     (nao depende de public.gerar_compromissos_*).
--
-- Banco novo / zerado:
--   1) Migrations: vecx/migrations (000_ordem_aplicacao.txt).
--   2) Cadastro pela aplicacao (tenant + usuario + schema provisionado).
--   3) Clientes e empresas pela aplicacao (dados no schema do tenant).
--   4) Compromissos pela UI/API ou cron do worker.
--
-- Como usar este arquivo:
--   Copie o UUID do tenant (SELECT abaixo), depois nas secoes 3–4 substitua
--   TROQUE_TENANT_UUID pelo valor.
-- =============================================================================

-- 1) Tenants e schema catalogado
SELECT t.id AS tenant_id,
       t.nome,
       t.active,
       c.schema_name
FROM public.tenant t
LEFT JOIN public.tenant_schema_catalog c ON c.tenant_id = t.id
ORDER BY t.nome;

-- 2) Usuarios (amostra)
SELECT u.id, u.email, u.role, u.tenantid, u.active
FROM public.usuario u
ORDER BY u.email
LIMIT 30;

-- 3) Conferir schema do seu tenant (substitua o UUID)
SELECT schema_name
FROM public.tenant_schema_catalog
WHERE tenant_id = 'TROQUE_TENANT_UUID'::uuid;

-- 4) Apos saber o schema_name (ex.: meu_escritorio), rode em sessao separada:
--    SET search_path TO meu_escritorio, public;
--    SELECT count(*) AS empresas FROM empresa WHERE ativo = true;
--    SELECT count(*) AS compromissos FROM empresa_compromissos;
--
--    Ou em uma linha (troque o schema):
--    SELECT count(*) FROM meu_escritorio.empresa WHERE ativo = true;

-- 5) Templates de obrigacao globais (para geracao bater com tipo_empresa do cliente)
SELECT id, descricao, tipo_empresa_id, periodicidade, ativo
FROM public.tipoempresa_obrigacao
WHERE ativo = true
ORDER BY tipo_empresa_id, descricao
LIMIT 100;

-- =============================================================================
-- Legado (issue #40): modelo antigo em public + funcoes SQL — nao usar em EF-916
-- =============================================================================
-- BEGIN;
-- SELECT public.gerar_compromissos_empresa(...);
-- ROLLBACK;
