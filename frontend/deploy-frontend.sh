#!/bin/bash
set -e

# 1. Carrega a URL do backend (config_privada.env)
source ../config_privada.env

echo "🎨 Gerando arquivos estáticos (SPA)..."

# 2. Build local que gera a pasta 'out'
# Injetamos a URL da API para que o JS saiba onde bater no Cloud Run
NEXT_PUBLIC_API_URL=$BACKEND_URL npm run build

echo "📂 Sincronizando out/ → backend/frontend/out (embed do binário Go local)..."
rm -rf ../backend/frontend/out
mkdir -p ../backend/frontend/out
cp -a out/. ../backend/frontend/out/

echo "✅ Frontend exportado em frontend/out e copiado para backend/frontend/out."