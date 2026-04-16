stop:
	@echo "Finalizando processos do VContab..."
	@pkill -f "go run cmd/api/main.go" || true
	@pkill -f "next-dev" || true
	@pkill -f "node" || true

.PHONY: stop frontend-build webview-build webview-run

frontend-build:
	@echo "Buildando frontend (export estático em frontend/out)..."
	@cd frontend && npm run build

webview-build: frontend-build
	@echo "Compilando WebView (frontend/main.go)..."
	@cd frontend && go build -o vecontab-desktop ./main.go

webview-run: frontend-build
	@echo "Rodando WebView (gera out/ se necessário)..."
	@cd frontend && go run ./main.go