#!/bin/bash

START_TIME=$(date +%s)
START_DATE=$(date +"%H:%M:%S")
LOCAL_DIR="/home/camaral/backups/vecontab"

echo "Limpando a pasta local dos backups"

#rm -rf "$LOCAL_DIR"/*

# No início do deploy.sh

./bkp_db.sh || { echo "Backup do banco de dados falhou, deploy cancelado"; exit 1; }

# Se o backup falhar e você quiser parar o deploy:
./backup.sh || { echo "Backup falhou, deploy cancelado"; exit 1; }

if ! systemctl is-active --quiet docker; then
    echo "❌ ERRO: O serviço Docker não está rodando no Fedora."
    echo "Execute: sudo systemctl start docker"
    exit 1
fi

#./bkp_drive.sh

# Não usamos set -e aqui para podermos fazer o resumo final mesmo com falhas
set +e

START_TIME=$(date +%s)
START_DATE=$(date +"%H:%M:%S")

# Arquivos para log completo (ajuda a diagnosticar depois no Fedora)
BACKEND_LOG=$(mktemp)
FRONTEND_LOG=$(mktemp)

echo ""
echo "--- Iniciando Deploy Global [Início: $START_DATE] ---"
echo "Log detalhado em: $BACKEND_LOG e $FRONTEND_LOG"
echo ""

# --- Executa frontend ---
echo "🎨 Processando Frontend (build estático para embed no backend)..."
if (cd frontend && ./deploy-frontend.sh > "$FRONTEND_LOG" 2>&1); then
    FRONTEND_STATUS="✅ Sucesso"
else
    FRONTEND_STATUS="❌ FALHOU"
fi
echo "Status: $FRONTEND_STATUS"

echo ""

# --- Executa backend ---
echo "📦 Processando Backend..."
# Rodamos o script e capturamos toda a saída (stdout + stderr)
if (cd backend && ./deploy-backend.sh > "$BACKEND_LOG" 2>&1); then
    BACKEND_STATUS="✅ Sucesso"
else
    BACKEND_STATUS="❌ FALHOU"
fi
# Mostra o status imediato no terminal para você não ficar no escuro
echo "Status: $BACKEND_STATUS"

echo ""

# Cálculos de tempo...
END_TIME=$(date +%s)
END_DATE=$(date +"%H:%M:%S")
ELAPSED=$(( END_TIME - START_TIME ))
MINUTES=$(( ELAPSED / 60 ))
SECONDS=$(( ELAPSED % 60 ))

# --- RESUMO FINAL ---
echo "-------------------------------------------"
echo "        RESUMO DO DEPLOY GLOBAL"
echo "-------------------------------------------"
echo "Início:       $START_DATE"
echo "Fim:          $END_DATE"
echo "Tempo Total:  ${MINUTES}m ${SECONDS}s"
echo ""

echo "Status Backend:  $BACKEND_STATUS"
if [ "$BACKEND_STATUS" == "❌ FALHOU" ]; then
    # Pega as últimas 5 linhas para garantir que veremos o erro do Docker ou GCloud
    echo "  L Últimas mensagens do log:"
    tail -n 5 "$BACKEND_LOG" | sed 's/^/    /'
fi

echo ""
echo "Status Frontend: $FRONTEND_STATUS"
if [ "$FRONTEND_STATUS" == "❌ FALHOU" ]; then
    echo "  L Últimas mensagens do log:"
    tail -n 5 "$FRONTEND_LOG" | sed 's/^/    /'
fi
echo "-------------------------------------------"

# Opcional: Apagar logs apenas se deram sucesso, manter se falharam
# rm -f "$BACKEND_LOG" "$FRONTEND_LOG"

echo "-------------------------------------------"
echo ""
echo "🧹 Faxina final: Removendo caches de build antigos (liberando espaço)..."
docker builder prune -f
echo ""
echo "-------------------------------------------"
echo ""
