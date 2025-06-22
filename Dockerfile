# Build stage
FROM golang:1.22-alpine AS builder

# Instalar dependências necessárias para build
RUN apk add --no-cache git ca-certificates tzdata

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copiar arquivos de dependências primeiro para melhor cache
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o code-explainer main.go

# Final stage
FROM alpine:3.19

# Instalar dependências de runtime
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Criar diretório da aplicação
WORKDIR /app

# Copiar binário do stage anterior
COPY --from=builder /app/code-explainer .

# Mudar propriedade para o usuário não-root
RUN chown -R appuser:appgroup /app

# Mudar para usuário não-root
USER appuser

# Labels para metadados
LABEL maintainer="mvcbotelho"
LABEL version="1.0.0"
LABEL description="Code Explainer - Ferramenta para explicar código usando IA"
LABEL org.opencontainers.image.source="https://github.com/mvcbotelho/code-explainer"

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ps aux | grep code-explainer || exit 1

# Expor porta (se necessário no futuro)
# EXPOSE 8080

# Entrypoint
ENTRYPOINT ["./code-explainer"]
