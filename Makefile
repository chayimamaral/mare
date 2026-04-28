stop:
	@echo "Finalizando processos do VContab..."
	@pkill -f "go run cmd/api/main.go" || true
	@pkill -f "next-dev" || true
	@pkill -f "node" || true

.PHONY: stop frontend-build webview-build webview-run backend-binaries-local

# Mesma ideia do frontend/deploy-frontend.sh: se existir config_privada.env, injeta NEXT_PUBLIC_API_URL.
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

webview-build: frontend-build
	@echo "Compilando WebView (frontend/main.go)..."
	@cd frontend && go build -o vecontab-desktop ./main.go

webview-run: frontend-build
	@echo "Rodando WebView (gera out/ se necessário)..."
	@cd frontend && go run ./main.go

backend-binaries-local: frontend-build
	@echo "Gerando binários locais em backend/bin..."
	@mkdir -p backend/bin
	@cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s" \
		-o ./bin/vecx-backend ./cmd/api/main.go
	@mkdir -p backend/bin/tools
	@mkdir -p .cache/go-mod .cache/go-build
	@test -x backend/bin/tools/garble || { \
		echo "garble não encontrado; instalando..."; \
		GOPATH="$(PWD)/.cache/go" GOBIN="$(PWD)/backend/bin/tools" GOMODCACHE="$(PWD)/.cache/go-mod" GOCACHE="$(PWD)/.cache/go-build" go install mvdan.cc/garble@latest; \
	}
	@cd backend && GARBLE_CACHE="$(PWD)/.cache/garble" CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ../backend/bin/tools/garble -literals -tiny build \
		-ldflags="-w -s" \
		-o ./bin/vecx-client.exe ./cmd/api/main.go
	@echo "OK: backend/bin/vecx-backend e backend/bin/vecx-client.exe"