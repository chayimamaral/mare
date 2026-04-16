#!/bin/bash
set -e

# 1. Carrega a URL do backend do seu arquivo local protegido
# O arquivo config_privada.env deve estar na raiz do projeto
source ../config_privada.env

# Configurações Genéricas
PROJECT_ID="vecontab"
REGION="us-central1"
REPO="vecontab-repo"
IMAGE_NAME="frontend"
FULL_IMAGE_PATH="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE_NAME:latest"

echo "🚀 Build e Deploy do Frontend..."

# 2. Build usando a variável que o 'source' carregou na memória
docker build --no-cache \
  --build-arg NEXT_PUBLIC_API_URL=$BACKEND_URL \
  -t $FULL_IMAGE_PATH .

# 3. Push e Deploy
docker push $FULL_IMAGE_PATH
gcloud run deploy vecontab-frontend \
  --image $FULL_IMAGE_PATH \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080