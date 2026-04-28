# ==========================================
# Estágio 1: Compilação (Builder)
# ==========================================
FROM golang:1.26.2-alpine AS builder

# Adicionamos um argumento que o Docker vai receber lá de fora
ARG APP_VERSION=dev

# Instala git e certificados (necessários para baixar dependências via HTTPS)
RUN apk add --no-cache git ca-certificates
#RUN go install github.com/burrowers/garble@latest
RUN go install mvdan.cc/garble@latest

WORKDIR /app

# Otimização de Cache: Baixa as dependências antes de copiar todo o código
COPY go.mod go.sum ./
RUN go mod download

# Copia o código (inclui backend/frontend/out com stub para o go:embed).
COPY . .

# Export do Next (frontend_static) deve sobrescrever o stub — senão o COPY acima
# deixa só o index.html mínimo e o app em produção não carrega.
COPY ./frontend_static ./frontend/out

# 1. Compila o SERVIDOR (Linux) injetando a versão
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X 'main.Version=${APP_VERSION}'" \
    -o vecx-backend ./cmd/api/main.go

# 2. Compila o CLIENTE (Windows .exe) com ofuscação (garble)
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 garble -literals -tiny build \
    -ldflags="-w -s -X 'main.Version=${APP_VERSION}'" \
    -o vecx-client.exe ./cmd/api/main.go
# ==========================================
# Estágio 2: Execução (Final)
# ==========================================
FROM alpine:latest

# 1. Instala certificados para TLS (Postgres) e tzdata para o horário do Brasil
RUN apk --no-cache add ca-certificates tzdata

# 2. Cria um usuário comum para não rodar como root (Segurança)
# -D ignora a criação de senha
RUN adduser -D vecxuser

# 3. Configura o fuso horário (Florianópolis/São Paulo)
ENV TZ=America/Sao_Paulo

# Define o diretório de trabalho
WORKDIR /app

# 4. Copia o binário do estágio builder e já define o dono como o novo usuário
# 4. Copia os arquivos e já define o "dono" como o vecxuser
COPY --from=builder --chown=vecxuser:vecxuser /app/vecx-backend .
COPY --from=builder --chown=vecxuser:vecxuser /app/vecx-client.exe .
# 5. Muda para o usuário sem privilégios
USER vecxuser

# O Cloud Run injetará o tráfego geralmente na porta 8080
EXPOSE 8080

# Execução do backend
# O uso de strings separadas no CMD ["cmd", "arg"] é a forma recomendada (exec form)
#CMD ["./vecontab-backend", "-port", "8080"]
CMD ["./vecx-backend"]