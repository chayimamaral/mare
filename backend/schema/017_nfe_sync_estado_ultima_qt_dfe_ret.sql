-- Coluna para último qtDfeRet retornado pelo retDistNFeSC (BT SC-2021/001, regra 118).
-- Aplicar em cada schema de tenant onde existir nfe_sync_estado, por exemplo:
--   SET search_path TO nome_do_schema, public;
--   \i backend/schema/017_nfe_sync_estado_ultima_qt_dfe_ret.sql
-- Ou executar o ALTER abaixo uma vez por schema.

ALTER TABLE nfe_sync_estado
  ADD COLUMN IF NOT EXISTS ultima_qt_dfe_ret integer NOT NULL DEFAULT 0;
