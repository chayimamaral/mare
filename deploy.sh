#!/bin/bash

# Aborta o script se qualquer comando falhar
set -e

# Captura o horário de início (segundos desde 1970 para cálculo e formato legível para exibição)
START_TIME=$(date +%s)
START_DATE=$(date +"%H:%M:%S")

echo ""
echo "--- Iniciando Deploy Global [Início: $START_DATE] ---"
echo ""

# Executa backend
echo "📦 Processando Backend..."
(cd backend && ./deploy-backend.sh)

echo ""

# Executa frontend
echo "🎨 Processando Frontend..."
(cd frontend && ./deploy-frontend.sh)

echo ""

# Captura o horário de fim
END_TIME=$(date +%s)
END_DATE=$(date +"%H:%M:%S")

# Calcula a diferença em segundos
ELAPSED=$(( END_TIME - START_TIME ))

# Formata o tempo total (Minutos e Segundos)
MINUTES=$(( ELAPSED / 60 ))
SECONDS=$(( ELAPSED % 60 ))

echo "✅ Deploy Global finalizado com sucesso!"
echo "-------------------------------------------"
echo "Início:      $START_DATE"
echo "Fim:         $END_DATE"
echo "Tempo Total: ${MINUTES}m ${SECONDS}s"
echo "-------------------------------------------"
echo ""