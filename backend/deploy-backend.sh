#!/bin/bash

set -e  # <--- CRÍTICO: Para o script se qualquer comando falhar

# Configurações
PROJECT_ID="vecontab"
REGION="us-central1"
REPO="vecontab-repo"
IMAGE_NAME="backend"
FULL_IMAGE_PATH="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE_NAME:latest"
VARS="PORT=3333,"
VARS+="PG_SSL_INSECURE=true,"
VARS+="COMPROMISSOS_WORKER_ENABLED=true,"
VARS+="COMPROMISSOS_WORKER_RUN_ON_STARTUP=true,"
VARS+="COMPROMISSOS_WORKER_CRON=0 5 1 * *,"
VARS+="COMPROMISSOS_WORKER_TIMEZONE=America/Sao_Paulo"

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

# 3. Deploy no Cloud Run

gcloud run deploy vecontab-backend \
  --image $FULL_IMAGE_PATH \
  --region us-central1 \
  --allow-unauthenticated \
  --port 3333 \
  --set-env-vars="COMPROMISSOS_WORKER_ENABLED=true,COMPROMISSOS_WORKER_CRON=0 5 1 * *,PG_SSL_INSECURE=true" \
  --set-secrets="PG_URL=PG_URL:latest,JWT_SECRET=JWT_SECRET:latest,VECONTAB_CERT_CRYPTO_KEY_HEX=VECONTAB_CERT_CRYPTO_KEY_HEX:latest"


#echo "🌍 Atualizando serviço no Cloud Run..."
#gcloud run deploy vecontab-backend \
#  --image $FULL_IMAGE_PATH \
#  --region $REGION \
#  --allow-unauthenticated
echo ""

echo "✅ Deploy do Backend finalizado!"
echo ""
