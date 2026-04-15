#!/bin/bash

set -e  # Para o script se qualquer comando falhar

# Configurações
PROJECT_ID="vecontab"
REGION="us-central1"
REPO="vecontab-repo"
IMAGE_NAME="backend"
FULL_IMAGE_PATH="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE_NAME:latest"

# Variáveis de Ambiente (Configuração que funcionou)
PG_URL="postgres://camaral:camaral85@34.30.48.97:5432/vecontab?sslmode=disable"
SERVER_PORT="8080"
JWT_SECRET="sua_chave_secreta_aqui"
PG_SSL_INSECURE="true"

echo ""
echo "🚀 Iniciando Deploy do Backend: $IMAGE_NAME"
echo ""

# 1. Build da Imagem
echo "📦 Gerando build da imagem Docker (Go)..."
docker build -t $FULL_IMAGE_PATH .
echo ""

# 2. Push para o Google Artifact Registry
echo "📤 Enviando imagem para o Google Cloud..."
docker push $FULL_IMAGE_PATH
echo ""

# 3. Deploy no Cloud Run (Forçando porta e variáveis para evitar cache do GCP)
echo "🌍 Atualizando serviço no Cloud Run..."
gcloud run deploy vecontab-backend \
  --image $FULL_IMAGE_PATH \
  --region $REGION \
  --allow-unauthenticated \
  --port $SERVER_PORT \
  --set-env-vars="PG_URL=$PG_URL,SERVER_PORT=$SERVER_PORT,JWT_SECRET=$JWT_SECRET,PG_SSL_INSECURE=$PG_SSL_INSECURE"
echo ""

echo "✅ Deploy do Backend finalizado e serviço online!"
echo ""