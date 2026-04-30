#!/bin/bash
set -e

# Deploy para VM: frontend é empacotado no backend via embed.
# Variáveis aceitas (opcionais):
# - DEPLOY_VM_HOST (default 34.30.48.97)
# - DEPLOY_APP_PORT (default 8080)
DEPLOY_VM_HOST="${DEPLOY_VM_HOST:-34.30.48.97}"
DEPLOY_APP_PORT="${DEPLOY_APP_PORT:-8080}"

# 1. Carrega config privada (quando existir)
if [ -f ../config_privada.env ]; then
  # shellcheck disable=SC1091
  source ../config_privada.env
fi

echo "🎨 Gerando arquivos estáticos (SPA)..."

# 2. Build local que gera a pasta 'out'
# Deploy em VM deve apontar para a própria VM por padrão (evita herdar BACKEND_URL antigo de Cloud Run).
# Override explícito:
# - DEPLOY_FRONTEND_API_URL
API_URL="${DEPLOY_FRONTEND_API_URL:-http://$DEPLOY_VM_HOST:$DEPLOY_APP_PORT}"
NEXT_PUBLIC_API_URL="$API_URL" npm run build

echo "📂 Sincronizando out/ → backend/frontend/out (embed do binário Go local)..."
rm -rf ../backend/frontend/out
mkdir -p ../backend/frontend/out
cp -a out/. ../backend/frontend/out/

echo "✅ Frontend exportado em frontend/out e copiado para backend/frontend/out."
echo "🌐 API alvo do frontend: $API_URL"