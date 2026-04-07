#!/bin/bash
# bkp_bd.sh - Backup Automático Vecontab (Postgres 18.3)

# Configurações (Ajuste conforme seu ambiente local)
DB_NAME="vecontab"
DB_USER="postgres"
DB_HOST="localhost"
BKP_DIR="docs/bkp_bd"
BKP_NAME="vecontab_bkp_geral.sql"
BKP_PATH="$BKP_DIR/$BKP_NAME"
TIMESTAMP=$(date +%Y%m%d_%H%M)

export PGPASSWORD='postgres'

echo "🐘 Iniciando backup automático do Postgres 18.3..."

# Garante que a pasta existe
mkdir -p "$BKP_DIR"

# Executa o pg_dump com as flags que você usava no pgAdmin:
# --no-owner: Remove comandos de OWNER TO
# --no-privileges: Remove GRANT/REVOKE
# --clean: (Opcional) Adiciona DROP TABLE antes do CREATE (ajuda na restauração)
# --if-exists: (Opcional) Evita erros se a tabela não existir no DROP
pg_dump -h "$DB_HOST" -U "$DB_USER" \
        --no-password \
        --no-owner \
        --no-privileges \
        --format=plain \
        --file="$BKP_PATH" \
        "$DB_NAME"

STATUS=$?        

# Limpa a senha da memória por segurança
unset PGPASSWORD

if [ $STATUS -eq 0 ]; then
    # Se o dump funcionou, cria a cópia de histórico
    cp "$BKP_PATH" "$BKP_DIR/vecontab_bkp_$TIMESTAMP.sql"
    echo "✅ Backup realizado com sucesso em: $BKP_PATH"
    echo "📦 Histórico: vecontab_bkp_$TIMESTAMP.sql"
    exit 0
else
    echo "❌ ERRO: Falha ao gerar o backup. Verifique a senha ou a conexão."
    exit 1
fi