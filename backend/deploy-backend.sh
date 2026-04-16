#!/bin/bash
set -e

# Configurações Genéricas
PROJECT_ID="vecontab"
REGION="us-central1"
REPO="vecontab-repo"
IMAGE_NAME="backend"
FULL_IMAGE_PATH="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE_NAME:latest"

echo "🚀 Build e Deploy do Backend..."

# 1. Build e Push
docker build -t $FULL_IMAGE_PATH .
docker push $FULL_IMAGE_PATH

# 2. Deploy (Sem as variáveis no comando!)
# O Cloud Run vai manter as variáveis que você configurou manualmente no Console
gcloud run deploy vecontab-backend \
  --image $FULL_IMAGE_PATH \
  --region $REGION \
  --allow-unauthenticated