# ==========================================
# Estágio 1: Compilação (Builder)
# ==========================================
FROM golang:1.26.2-alpine AS builder

# Instala git e certificados (necessários para baixar dependências via HTTPS)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Otimização de Cache: Baixa as dependências antes de copiar todo o código
COPY go.mod go.sum ./
RUN go mod download

# Copia o código (inclui backend/frontend/out com stub para o go:embed).
COPY . .

# Export do Next (frontend_static) deve sobrescrever o stub — senão o COPY acima
# deixa só o index.html mínimo e o app em produção não carrega.
COPY ./frontend_static ./frontend/out

# Compila o binário estático para Linux
# -ldflags="-w -s" remove tabelas de símbolos e debug para reduzir o tamanho
#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o vecontab-backend ./cmd/api/main.go

# 1. Compila o binário do SERVIDOR (Linux)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o vecx-backend ./cmd/api/main.go

# 2. Compila o executável do CLIENTE (Windows .exe)
RUN CGO_ENABLED=0 GOOS=windows go build -ldflags="-w -s" -o vecx-client.exe ./cmd/api/main.go

# ==========================================
# Estágio 2: Execução (Final)
# ==========================================
FROM alpine:latest

# 1. Instala certificados para TLS (Postgres) e tzdata para o horário do Brasil
RUN apk --no-cache add ca-certificates tzdata

# 2. Cria um usuário comum para não rodar como root (Segurança)
# -D ignora a criação de senha
RUN adduser -D vecontabuser

# 3. Configura o fuso horário (Florianópolis/São Paulo)
ENV TZ=America/Sao_Paulo

# Define o diretório de trabalho
WORKDIR /app

# 4. Copia o binário do estágio builder e já define o dono como o novo usuário
COPY --from=builder --chown=vecontabuser:vecontabuser /app/vecontab-backend .

# 5. Muda para o usuário sem privilégios
USER vecontabuser

# O Cloud Run injetará o tráfego geralmente na porta 8080
EXPOSE 8080

# Execução do backend
# O uso de strings separadas no CMD ["cmd", "arg"] é a forma recomendada (exec form)
#CMD ["./vecontab-backend", "-port", "8080"]
CMD ["./vecontab-backend"]