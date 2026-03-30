# Estágio 1: Build do binário Go
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código e compila
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go 

# Estágio 2: Runner
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia o binário do estágio anterior
COPY --from=builder /app/server .
# Copia arquivos de configuração ou .env se necessário
COPY --from=builder /app/.env . 

EXPOSE 3333

CMD ["./server"]