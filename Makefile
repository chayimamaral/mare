stop:
	@echo "Finalizando processos do VContab..."
	@pkill -f "go run cmd/api/main.go" || true
	@pkill -f "next-dev" || true
	@pkill -f "node" || true

.PHONY: stop frontend-build webview-build webview-run backend-binaries-local local-agent-binaries-local encrypt-env

# Mesma ideia do frontend/deploy-frontend.sh: se existir config_privada.env, injeta NEXT_PUBLIC_API_URL.
#make frontend-build
frontend-build:
	@echo "Buildando frontend (export estático em frontend/out)..."
	@bash -c 'set -e; \
	  if [ -f config_privada.env ]; then \
	    set -a; . ./config_privada.env; set +a; \
	    cd frontend && NEXT_PUBLIC_API_URL="$$BACKEND_URL" npm run build; \
	  else \
	    cd frontend && npm run build; \
	  fi'
	@echo "Sincronizando out/ → backend/frontend/out (embed Go)..."
	@rm -rf backend/frontend/out && mkdir -p backend/frontend/out && cp -a frontend/out/. backend/frontend/out/

#make webview-build
webview-build: frontend-build
	@echo "Compilando WebView (frontend/main.go)..."
	@mkdir -p backend/bin
	@cd frontend && go build -o ../backend/bin/vecx-desktop ./main.go

#make webview-run
webview-run: frontend-build
	@echo "Rodando WebView (gera out/ se necessário)..."
	@cd frontend && go run ./main.go

#make backend-binaries-local
backend-binaries-local: frontend-build
	@echo "Gerando binários locais em backend/bin..."
	@mkdir -p backend/bin
	@bash -c 'set -e; \
	  KEY_FILE=""; \
	  if [ -f .env.senha_compilacao ]; then KEY_FILE=".env.senha_compilacao"; \
	  elif [ -f backend/.env.senha_compilacao ]; then KEY_FILE="backend/.env.senha_compilacao"; \
	  else echo "Arquivo de senha nao encontrado (.env.senha_compilacao ou backend/.env.senha_compilacao)"; exit 1; fi; \
	  set -a; . "./$$KEY_FILE"; set +a; \
	  KEY="$${VECX_MASTER_KEY:-$${VECONTAB_MASTER_KEY:-$${SENHA_COMPILACAO:-}}}"; \
	  test -n "$$KEY" || { echo "VECX_MASTER_KEY ausente em $$KEY_FILE"; exit 1; }; \
	  LDFLAGS="-w -s -X '\''github.com/chayimamaral/vecx/backend/pkg/masterkey.EmbeddedMasterKey=$$KEY'\''"; \
	  cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$$LDFLAGS" -o ./bin/vecx ./cmd/api/main.go; \
	'
	@mkdir -p backend/bin/tools
	@mkdir -p .cache/go-mod .cache/go-build
	@test -x backend/bin/tools/garble || { \
		echo "garble não encontrado; instalando..."; \
		GOPATH="$(PWD)/.cache/go" GOBIN="$(PWD)/backend/bin/tools" GOMODCACHE="$(PWD)/.cache/go-mod" GOCACHE="$(PWD)/.cache/go-build" go install mvdan.cc/garble@latest; \
	}
	@bash -c 'set -e; \
	  KEY_FILE=""; \
	  if [ -f .env.senha_compilacao ]; then KEY_FILE=".env.senha_compilacao"; \
	  elif [ -f backend/.env.senha_compilacao ]; then KEY_FILE="backend/.env.senha_compilacao"; \
	  else echo "Arquivo de senha nao encontrado (.env.senha_compilacao ou backend/.env.senha_compilacao)"; exit 1; fi; \
	  set -a; . "./$$KEY_FILE"; set +a; \
	  KEY="$${VECX_MASTER_KEY:-$${VECONTAB_MASTER_KEY:-$${SENHA_COMPILACAO:-}}}"; \
	  test -n "$$KEY" || { echo "VECX_MASTER_KEY ausente em $$KEY_FILE"; exit 1; }; \
	  LDFLAGS="-w -s -X '\''github.com/chayimamaral/vecx/backend/pkg/masterkey.EmbeddedMasterKey=$$KEY'\''"; \
	  cd backend && GARBLE_CACHE="$(PWD)/.cache/garble" CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ../backend/bin/tools/garble -literals -tiny build -ldflags="$$LDFLAGS" -o ./bin/vecx.exe ./cmd/api/main.go; \
	'
	@echo "OK: backend/bin/vecx e backend/bin/vecx.exe"
	@$(MAKE) local-agent-binaries-local

#make local-agent-binaries-local
local-agent-binaries-local:
	@echo "Gerando binários do agente local em backend/bin..."
	@mkdir -p backend/bin
	@bash -c 'set -e; \
	  cd agente-local && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ../backend/bin/vecx-agent ./cmd/agent/main.go; \
	  command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1 || { \
	    echo "Compilador MinGW nao encontrado: instale mingw64-gcc para gerar vecx-agent.exe"; \
	    exit 1; \
	  }; \
	  CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o ../backend/bin/vecx-agent.exe ./cmd/agent/main.go; \
	'
	@echo "OK: backend/bin/vecx-agent e backend/bin/vecx-agent.exe"


#make encrypt-env ENV_FILE=backend/bin/.env.<cliente>
encrypt-env:
	@bash -c 'set -e; \
	  FILE="$${ENV_FILE:-}"; \
	  if [ -z "$$FILE" ]; then \
	    echo "Uso: make encrypt-env ENV_FILE=<caminho/.env.<cliente>>"; \
	    echo "Ex.: make encrypt-env ENV_FILE=backend/bin/.env.acme"; \
	    exit 1; \
	  fi; \
	  test -f "$$FILE" || { echo "Arquivo nao encontrado: $$FILE"; exit 1; }; \
	  if [ -f .env.senha_compilacao ]; then :; \
	  elif [ -f backend/.env.senha_compilacao ]; then :; \
	  else echo "Arquivo de senha nao encontrado (.env.senha_compilacao ou backend/.env.senha_compilacao)"; exit 1; fi; \
	  cd tools/encryptor && printf "%s\n" "$(realpath "$$FILE")" | go run .; \
	'
