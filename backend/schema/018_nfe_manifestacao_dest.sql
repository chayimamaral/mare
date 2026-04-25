-- Histórico de manifestação do destinatário (Recepção de Evento SVRS, tpEvento 2102xx). Aplicar em cada schema de tenant.
-- SET search_path TO nome_do_schema, public;

CREATE TABLE IF NOT EXISTS nfe_manifestacao_dest (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    chave_nfe char(44) NOT NULL,
    tp_evento varchar(8) NOT NULL,
    cnpj_dest varchar(14) NOT NULL,
    cstat_lote integer NOT NULL DEFAULT 0,
    x_motivo_lote text,
    cstat_evento integer NOT NULL DEFAULT 0,
    x_motivo_evento text,
    n_prot varchar(60),
    retorno_xml text,
    criado_em timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_nfe_manifest_dest_chave ON nfe_manifestacao_dest (chave_nfe);
CREATE INDEX IF NOT EXISTS idx_nfe_manifest_dest_criado ON nfe_manifestacao_dest (criado_em DESC);
