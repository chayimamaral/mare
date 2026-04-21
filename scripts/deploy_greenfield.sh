#!/usr/bin/env bash
set -euo pipefail

# Deploy greenfield completo:
# - aplica 000 (baseline)
# - aplica 010 (seed global)
# - aplica 020 (seed municipios)
# - executa validacoes de aceite
#
# Uso:
#   DB_HOST=localhost DB_PORT=5432 DB_USER=camaral DB_NAME=vecontab ./scripts/deploy_greenfield.sh
# ou:
#   ./scripts/deploy_greenfield.sh localhost 5432 camaral vecontab

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIG_DIR="$ROOT_DIR/migrations"

DB_HOST="${DB_HOST:-${1:-localhost}}"
DB_PORT="${DB_PORT:-${2:-5432}}"
DB_USER="${DB_USER:-${3:-camaral}}"
DB_NAME="${DB_NAME:-${4:-vecontab}}"

PSQL=(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1)

run_migration() {
  local file="$1"
  echo "==> Aplicando: $file"
  "${PSQL[@]}" -f "$file"
}

echo "=================================================="
echo "Deploy greenfield Vecontab"
echo "Host: $DB_HOST  Port: $DB_PORT  User: $DB_USER  DB: $DB_NAME"
echo "=================================================="

existing_public_tables="$("${PSQL[@]}" -tA -c "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';")"
if [[ "${existing_public_tables:-0}" != "0" ]]; then
  echo "ERRO: banco nao esta vazio (public possui $existing_public_tables tabelas)."
  echo "Este script e para deploy greenfield em banco vazio."
  exit 1
fi

run_migration "$MIG_DIR/20260421_000000_ef916_greenfield_schema_limpo.sql"
run_migration "$MIG_DIR/20260421_010000_ef001_seed_global_from_original.sql"
run_migration "$MIG_DIR/20260421_020000_ef001_seed_municipios_dtb2025.sql"

echo "==> Validando criterios de aceite..."

tenant_count="$("${PSQL[@]}" -tA -c "SELECT count(*) FROM public.tenant WHERE id = '00000000-0000-4000-8000-000000000001'::uuid;")"
super_count="$("${PSQL[@]}" -tA -c "SELECT count(*) FROM public.usuario WHERE tenantid='00000000-0000-4000-8000-000000000001'::uuid AND role='SUPER'::public.role;")"
schema_count="$("${PSQL[@]}" -tA -c "SELECT count(*) FROM public.tenant_schema_catalog WHERE tenant_id='00000000-0000-4000-8000-000000000001'::uuid AND schema_name='vec_sistemas';")"
estado_count="$("${PSQL[@]}" -tA -c "SELECT count(*) FROM public.estado;")"
municipio_count="$("${PSQL[@]}" -tA -c "SELECT count(*) FROM public.municipio;")"
feriado_nacional="$("${PSQL[@]}" -tA -c "SELECT count(*) FROM pg_enum e JOIN pg_type t ON t.oid=e.enumtypid WHERE t.typname='feriado' AND e.enumlabel='NACIONAL';")"

fail=0
[[ "$tenant_count" == "1" ]] || { echo "ERRO: tenant plataforma nao encontrado"; fail=1; }
[[ "$super_count" =~ ^[1-9][0-9]*$ ]] || { echo "ERRO: nenhum usuario SUPER encontrado no tenant da plataforma"; fail=1; }
[[ "$schema_count" == "1" ]] || { echo "ERRO: tenant_schema_catalog do tenant plataforma nao encontrado"; fail=1; }
[[ "$estado_count" == "27" ]] || { echo "ERRO: estado esperado=27 atual=$estado_count"; fail=1; }
[[ "$municipio_count" == "5571" ]] || { echo "ERRO: municipio esperado=5571 atual=$municipio_count"; fail=1; }
[[ "$feriado_nacional" == "1" ]] || { echo "ERRO: enum feriado sem valor NACIONAL"; fail=1; }

if [[ "$fail" -ne 0 ]]; then
  echo "=================================================="
  echo "FALHA: validacoes de aceite nao passaram."
  echo "=================================================="
  exit 1
fi

echo "=================================================="
echo "SUCESSO: deploy greenfield concluido e validado."
echo "Tenant/SUPER, estados, municipios e enum feriado OK."
echo "=================================================="

