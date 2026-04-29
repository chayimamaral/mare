#!/bin/bash
set -e

# Configurações Genéricas
PROJECT_ID="vecx"
REGION="us-central1"
REPO="vecx-repo"
IMAGE_NAME="backend"
FULL_IMAGE_PATH="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE_NAME:latest"
LOCAL_BIN_DIR="./bin"
LOCAL_BIN_PATH="$LOCAL_BIN_DIR/vecx-backend"

echo "🚚 Trazendo o frontend para o contexto do backend..."
rm -rf ./frontend_static
if [ ! -d ../frontend/out ]; then
  echo "❌ frontend/out não encontrado."
  echo "Execute primeiro o build do frontend (ex.: cd ../frontend && ./deploy-frontend.sh)."
  exit 1
fi
cp -r ../frontend/out ./frontend_static

echo "🔨 Gerando binário local do backend..."
mkdir -p "$LOCAL_BIN_DIR"
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o "$LOCAL_BIN_PATH" ./cmd/api/main.go
echo "✅ Binário local gerado em: $LOCAL_BIN_PATH"

echo "🚀 Build e Deploy do Backend..."

# 1. Build e Push
docker build -t $FULL_IMAGE_PATH .
docker push $FULL_IMAGE_PATH

# 2. Deploy (Sem as variáveis no comando!)
# O Cloud Run vai manter as variáveis que você configurou manualmente no Console
gcloud run deploy vecx-backend \
  --image $FULL_IMAGE_PATH \
  --region $REGION \
  --allow-unauthenticated

echo "📦 Binário local disponível em: $LOCAL_BIN_PATH"
