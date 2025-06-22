# Etapa de build
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o code-explainer main.go

# Etapa final
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/code-explainer .

RUN chmod +x ./code-explainer

ENTRYPOINT ["./code-explainer"]
