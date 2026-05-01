# --- VARIÁVEIS DE CONFIGURAÇÃO ---
BIN_DIR := bin

stop:
	@echo "Finalizando processos do VContab..."
	@pkill -f "go run cmd/api/main.go" || true
	@pkill -f "next-dev" || true
	@pkill -f "node" || true

.PHONY: stop frontend-build webview-build webview-run backend-binaries-local local-agent-binaries-local local-agent-gui-binaries-local encrypt-env

# make frontend-build
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

# make webview-build
webview-build: frontend-build
	@echo "Compilando WebView (frontend/main.go)..."
	@mkdir -p $(BIN_DIR)
	@cd frontend && go build -o ../$(BIN_DIR)/vecx-desktop ./main.go

# make webview-run
webview-run: frontend-build
	@echo "Rodando WebView (gera out/ se necessário)..."
	@cd frontend && go run ./main.go

# make backend-binaries-local
backend-binaries-local: frontend-build
	@echo "Gerando binários locais em $(BIN_DIR)..."
	@mkdir -p $(BIN_DIR)/tools
	@mkdir -p .cache/go-mod .cache/go-build
	@bash -c 'set -e; \
		KEY_FILE=""; \
		if [ -f .env.senha_compilacao ]; then KEY_FILE=".env.senha_compilacao"; \
		elif [ -f backend/.env.senha_compilacao ]; then KEY_FILE="backend/.env.senha_compilacao"; \
		else echo "Arquivo de senha nao encontrado (.env.senha_compilacao ou backend/.env.senha_compilacao)"; exit 1; fi; \
		set -a; . "./$$KEY_FILE"; set +a; \
		KEY="$${VECX_MASTER_KEY:-$${VECONTAB_MASTER_KEY:-$${SENHA_COMPILACAO:-}}}"; \
		test -n "$$KEY" || { echo "VECX_MASTER_KEY ausente em $$KEY_FILE"; exit 1; }; \
		LDFLAGS="-w -s -X '\''github.com/chayimamaral/vecx/backend/pkg/masterkey.EmbeddedMasterKey=$$KEY'\''"; \
		cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$$LDFLAGS" -o ../$(BIN_DIR)/vecx ./cmd/api/main.go; \
	'
	@test -x $(BIN_DIR)/tools/garble || { \
		echo "garble não encontrado; instalando..."; \
		GOPATH="$(PWD)/.cache/go" GOBIN="$(PWD)/$(BIN_DIR)/tools" GOMODCACHE="$(PWD)/.cache/go-mod" GOCACHE="$(PWD)/.cache/go-build" go install mvdan.cc/garble@latest; \
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
		cd backend && GARBLE_CACHE="$(PWD)/.cache/garble" CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ../$(BIN_DIR)/tools/garble -literals -tiny build -ldflags="$$LDFLAGS" -o ../$(BIN_DIR)/vecx.exe ./cmd/api/main.go; \
	'
	@echo "OK: $(BIN_DIR)/vecx e $(BIN_DIR)/vecx.exe"
	@$(MAKE) local-agent-binaries-local
	@$(MAKE) local-agent-gui-binaries-local || echo "AVISO: GUI do agente nao gerada neste ambiente (dependencias graficas ausentes)."

# make local-agent-binaries-local
local-agent-binaries-local:
	@echo "Gerando binários do agente local em $(BIN_DIR)..."
	@mkdir -p $(BIN_DIR)
	@bash -c 'set -e; \
		cd agente-local && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ../$(BIN_DIR)/vecx-agent-cli ./cmd/agent/main.go; \
		command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1 || { \
			echo "Compilador MinGW nao encontrado: instale mingw64-gcc para gerar vecx-agent-cli.exe"; \
			exit 1; \
		}; \
		CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o ../$(BIN_DIR)/vecx-agent-cli.exe ./cmd/agent/main.go; \
	'
	@echo "OK: $(BIN_DIR)/vecx-agent-cli e $(BIN_DIR)/vecx-agent-cli.exe"

# make local-agent-gui-binaries-local
local-agent-gui-binaries-local:
	@echo "Gerando binários GUI do agente local em $(BIN_DIR)..."
	@mkdir -p $(BIN_DIR)
	@bash -c 'set -e; \
		if [ "$$(uname -s)" = "Linux" ]; then \
			pkg-config --exists egl wayland-egl wayland-client wayland-cursor x11 xkbcommon xkbcommon-x11 x11-xcb xcursor xfixes || { \
				echo "Dependencias GUI ausentes no Linux (pkg-config: xkbcommon-x11/x11/etc)."; \
				echo "vecx-agent sera o mesmo binario que vecx-agent-cli (somente terminal)."; \
				echo "Instale no Fedora (headers + pkg-config para Gio/CGO):"; \
				echo "sudo dnf install mesa-libEGL-devel wayland-devel libX11-devel libxcb-devel libxkbcommon-devel libxkbcommon-x11-devel libXcursor-devel libXfixes-devel"; \
				cp -f ./$(BIN_DIR)/vecx-agent-cli ./$(BIN_DIR)/vecx-agent; \
				exit 0; \
			}; \
		fi; \
		cd agente-local && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags gui -o ../$(BIN_DIR)/vecx-agent ./cmd/agent-gui/main.go; \
		command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1 || { \
			echo "Compilador MinGW nao encontrado: instale mingw64-gcc para gerar vecx-agent.exe"; \
			exit 1; \
		}; \
		CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -tags gui -o ../$(BIN_DIR)/vecx-agent.exe ./cmd/agent-gui/main.go; \
	'
	@echo "OK: $(BIN_DIR)/vecx-agent e $(BIN_DIR)/vecx-agent.exe"

# make encrypt-env ENV_FILE=bin/.env.<cliente>
encrypt-env:
	@bash -c 'set -e; \
		FILE="$${ENV_FILE:-}"; \
		if [ -z "$$FILE" ]; then \
			echo "Uso: make encrypt-env ENV_FILE=<caminho/.env.<cliente>>"; \
			echo "Ex.: make encrypt-env ENV_FILE=$(BIN_DIR)/.env.acme"; \
			exit 1; \
		fi; \
		test -f "$$FILE" || { echo "Arquivo nao encontrado: $$FILE"; exit 1; }; \
		if [ -f .env.senha_compilacao ]; then :; \
		elif [ -f backend/.env.senha_compilacao ]; then :; \
		else echo "Arquivo de senha nao encontrado (.env.senha_compilacao ou backend/.env.senha_compilacao)"; exit 1; fi; \
		cd tools/encryptor && printf "%s\n" "$(realpath "$$FILE")" | go run .; \
	'