#!/bin/bash
set -e

# Raiz do repositório vecx e binário da API (make coloca em repo/bin/vecx — não confiar em $PWD)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
LOCAL_BIN_PATH="$ROOT_DIR/bin/vecx"

# Deploy em VM (override por variáveis de ambiente)
#
# Modo A — IP/host direto (ssh/scp):
#   DEPLOY_VM_HOST, DEPLOY_VM_SSH_PORT, DEPLOY_SSH_KEY_PATH, etc.
#
# Modo B — instância GCP (como: gcloud compute ssh vecontab-db --zone=us-central1-a):
#   export DEPLOY_GCE_INSTANCE=vecontab-db
#   export DEPLOY_GCE_ZONE=us-central1-a
#   opcional: DEPLOY_GCE_PROJECT, DEPLOY_GCE_IAP=1 (adiciona --tunnel-through-iap)
#
DEPLOY_VM_HOST="${DEPLOY_VM_HOST:-34.30.48.97}"
DEPLOY_VM_SSH_USER="${DEPLOY_VM_SSH_USER:-camaral}"
DEPLOY_VM_SSH_PORT="${DEPLOY_VM_SSH_PORT:-22}"
DEPLOY_GCE_INSTANCE="${DEPLOY_GCE_INSTANCE:-}"
DEPLOY_GCE_ZONE="${DEPLOY_GCE_ZONE:-us-central1-a}"
DEPLOY_GCE_PROJECT="${DEPLOY_GCE_PROJECT:-}"
DEPLOY_GCE_IAP="${DEPLOY_GCE_IAP:-}"
# Diretório na VM com permissão de escrita (sem sudo). Override: DEPLOY_VM_TARGET_DIR
DEPLOY_VM_TARGET_DIR="${DEPLOY_VM_TARGET_DIR:-/home/${DEPLOY_VM_SSH_USER}/vecx-app}"
DEPLOY_APP_PORT="${DEPLOY_APP_PORT:-8080}"
DEPLOY_SSH_KEY_PATH="${DEPLOY_SSH_KEY_PATH:-}"
DEPLOY_REMOTE_ENV_PATH="${DEPLOY_REMOTE_ENV_PATH:-$DEPLOY_VM_TARGET_DIR/.env}"
DEPLOY_REMOTE_PID_PATH="${DEPLOY_REMOTE_PID_PATH:-$DEPLOY_VM_TARGET_DIR/vecx.pid}"
DEPLOY_REMOTE_LOG_PATH="${DEPLOY_REMOTE_LOG_PATH:-$DEPLOY_VM_TARGET_DIR/vecx.log}"
DEPLOY_REMOTE_BIN_PATH="${DEPLOY_REMOTE_BIN_PATH:-$DEPLOY_VM_TARGET_DIR/vecx}"
DEPLOY_REMOTE_BIN_NEW_PATH="${DEPLOY_REMOTE_BIN_NEW_PATH:-$DEPLOY_VM_TARGET_DIR/vecx.new}"
DEPLOY_ENV_SOURCE_FILE="${DEPLOY_ENV_SOURCE_FILE:-$SCRIPT_DIR/.env}"

# ssh usa -p para porta; scp usa -P (minúsculo em scp = preserva datas e quebra o deploy)
SSH_OPTS=(-p "$DEPLOY_VM_SSH_PORT" -o StrictHostKeyChecking=accept-new)
SCP_OPTS=(-P "$DEPLOY_VM_SSH_PORT" -o StrictHostKeyChecking=accept-new)
if [ -n "$DEPLOY_SSH_KEY_PATH" ]; then
  SSH_OPTS+=(-i "$DEPLOY_SSH_KEY_PATH")
  SCP_OPTS+=(-i "$DEPLOY_SSH_KEY_PATH")
fi

GCLOUD_BASE=("--zone=$DEPLOY_GCE_ZONE")
if [ -n "$DEPLOY_GCE_PROJECT" ]; then
  GCLOUD_BASE+=("--project=$DEPLOY_GCE_PROJECT")
fi
if [ -n "$DEPLOY_GCE_IAP" ] && [ "$DEPLOY_GCE_IAP" != "0" ]; then
  GCLOUD_BASE+=("--tunnel-through-iap")
fi

GCE_TARGET="${DEPLOY_VM_SSH_USER}@${DEPLOY_GCE_INSTANCE}"

remote_mkdir() {
  if [ -n "$DEPLOY_GCE_INSTANCE" ]; then
    gcloud compute ssh "$GCE_TARGET" "${GCLOUD_BASE[@]}" \
      --command="mkdir -p '$DEPLOY_VM_TARGET_DIR'"
  else
    ssh "${SSH_OPTS[@]}" "$REMOTE" "mkdir -p '$DEPLOY_VM_TARGET_DIR'"
  fi
}

remote_scp() {
  # Uso: remote_scp ARQUIVO_LOCAL DESTINO_REMOTO (ex.: /home/user/app/vecx)
  local src="$1"
  local dest="$2"
  if [ -n "$DEPLOY_GCE_INSTANCE" ]; then
    gcloud compute scp "$src" "$GCE_TARGET:$dest" "${GCLOUD_BASE[@]}"
  else
    scp "${SCP_OPTS[@]}" "$src" "$REMOTE:$dest"
  fi
}

remote_bootstrap_vecx() {
  cat <<EOF
set -e
mkdir -p '$DEPLOY_VM_TARGET_DIR'
if [ -f '$DEPLOY_REMOTE_PID_PATH' ]; then
  OLD_PID=\$(cat '$DEPLOY_REMOTE_PID_PATH' || true)
  if [ -n "\$OLD_PID" ] && kill -0 "\$OLD_PID" 2>/dev/null; then
    kill "\$OLD_PID" || true
    sleep 1
  fi
fi
# Não fazer source do .env aqui: o vecx carrega com godotenv ao lado do binário (config.Load).
# Um .env “de aplicação” costuma ter linhas que o bash não aceita (URLs, comentários, etc.).
if [ -f '$DEPLOY_REMOTE_BIN_NEW_PATH' ]; then
  mv -f '$DEPLOY_REMOTE_BIN_NEW_PATH' '$DEPLOY_REMOTE_BIN_PATH'
fi
chmod +x '$DEPLOY_REMOTE_BIN_PATH'
export PORT='$DEPLOY_APP_PORT'
nohup '$DEPLOY_REMOTE_BIN_PATH' >> '$DEPLOY_REMOTE_LOG_PATH' 2>&1 &
echo \$! > '$DEPLOY_REMOTE_PID_PATH'
EOF
}

echo "🚚 Trazendo o frontend para o contexto do backend..."
rm -rf "$SCRIPT_DIR/frontend_static"
if [ ! -d "$ROOT_DIR/frontend/out" ]; then
  echo "❌ frontend/out não encontrado."
  echo "Execute primeiro o build do frontend (ex.: cd \"$ROOT_DIR/frontend\" && ./deploy-frontend.sh)."
  exit 1
fi
cp -a "$ROOT_DIR/frontend/out" "$SCRIPT_DIR/frontend_static"

echo "🔨 Gerando binário local do backend (via Makefile em $ROOT_DIR → $LOCAL_BIN_PATH)..."
(cd "$ROOT_DIR" && make backend-binaries-local)
echo "✅ Binário para upload: $LOCAL_BIN_PATH"

if [ ! -f "$LOCAL_BIN_PATH" ]; then
  echo "❌ Binário não encontrado em $LOCAL_BIN_PATH"
  exit 1
fi

REMOTE="$DEPLOY_VM_SSH_USER@$DEPLOY_VM_HOST"

if [ -n "$DEPLOY_GCE_INSTANCE" ]; then
  echo "🧰 Modo GCE: instância $DEPLOY_GCE_INSTANCE (zona $DEPLOY_GCE_ZONE) → $DEPLOY_VM_TARGET_DIR"
else
  echo "🧰 Preparando diretório remoto em $REMOTE:$DEPLOY_VM_TARGET_DIR ..."
fi
remote_mkdir

echo "📤 Enviando binário para VM..."
remote_scp "$LOCAL_BIN_PATH" "$DEPLOY_REMOTE_BIN_NEW_PATH"

echo "ℹ️ .env remoto preservado (deploy nunca faz upload de .env)."

echo "🚀 Reiniciando processo vecx na VM (porta $DEPLOY_APP_PORT)..."
if [ -n "$DEPLOY_GCE_INSTANCE" ]; then
  remote_bootstrap_vecx | gcloud compute ssh "$GCE_TARGET" "${GCLOUD_BASE[@]}" -- bash -s
else
  remote_bootstrap_vecx | ssh "${SSH_OPTS[@]}" "$REMOTE" bash -s
fi

if [ -n "$DEPLOY_GCE_INSTANCE" ]; then
  echo "✅ Backend deploy finalizado em $DEPLOY_GCE_INSTANCE ($DEPLOY_GCE_ZONE), porta $DEPLOY_APP_PORT"
else
  echo "✅ Backend deploy finalizado na VM $DEPLOY_VM_HOST:$DEPLOY_APP_PORT"
fi
echo "📄 Log remoto: $DEPLOY_REMOTE_LOG_PATH"
