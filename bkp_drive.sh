#!/bin/bash

START_TIME=$(date +%s)
START_DATE=$(date +"%H:%M:%S")

# Configuração de destino (centralizado para evitar erro de digitação)
REMOTE_DEST="bkp_drive:/bkp_vecx"
LOCAL_DIR="/home/camaral/backups/vecx"

echo ""
echo "--- Faxina no Google Drive: Mantendo apenas os últimos 7 ---"
# O 2>/dev/null evita o erro caso a pasta esteja vazia ou não exista
# --- Limpeza dos Backups de Banco de Dados (.sql) ---
echo "Limpando backups antigos do Banco..."
rclone lsf "$REMOTE_DEST" --format "p" 2>/dev/null | grep "bkp_vecx_db_" | sort | head -n -7 | xargs -I {} rclone delete "$REMOTE_DEST/{}"

# --- Limpeza dos Backups de Fontes (.tar.gz) ---
echo "Limpando backups antigos dos Fontes..."
rclone lsf "$REMOTE_DEST" --format "p" 2>/dev/null | grep "vecx_source_" | sort | head -n -7 | xargs -I {} rclone delete "$REMOTE_DEST/{}"
echo ""
echo "--- Iniciando cópia para o Google Drive ---"
rclone copy "$LOCAL_DIR" "$REMOTE_DEST"

# O STATUS captura se o rclone copy deu certo
IF_SUCCESS=$?

if [ $IF_SUCCESS -eq 0 ]; then
    echo "✅ Upload concluído...."
#    rm -rf "$LOCAL_DIR"/*
else
    echo "❌ ERRO no upload. Os arquivos locais foram preservados para segurança."
fi

END_TIME=$(date +%s)
END_DATE=$(date +"%H:%M:%S")
ELAPSED=$(( END_TIME - START_TIME ))
MINUTES=$(( ELAPSED / 60 ))
SECONDS=$(( ELAPSED % 60 ))

echo "-------------------------------------------"
echo "        RESUMO DO BKP para Google Drive"
echo "-------------------------------------------"
echo "Início:       $START_DATE"
echo "Fim:          $END_DATE"
echo "Tempo Total:  ${MINUTES}m ${SECONDS}s"
echo ""
