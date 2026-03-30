# Estágio 1: Build do Next.js
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

# ESTAS LINHAS SÃO A CHAVE:
ARG NEXT_PUBLIC_API_URL=https://vecontab-backend-822119889622.us-central1.run.app
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

RUN npm run build

# Estágio 2: Runner Node.js
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=8080
ENV HOSTNAME=0.0.0.0

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 8080

CMD ["node", "server.js"]