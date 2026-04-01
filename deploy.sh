#!/bin/bash

# Aborta o script se qualquer comando falhar
set -e

echo ""
echo "--- Iniciando Deploy Global ---"
echo ""
# Executa backend
echo "📦 Processando Backend..."
(cd backend && ./deploy-backend.sh)

echo ""
# Executa frontend
echo "🎨 Processando Frontend..."
(cd frontend && ./deploy-frontend.sh)
echo ""
echo "✅ Deploy Global finalizado com sucesso!"
echo ""
