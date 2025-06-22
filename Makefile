# Makefile para Code Explainer com Docker

# Variáveis
APP_NAME=code-explainer
DOCKER_USER=mvcbotelho
DOCKER_TAG=latest
IMAGE=$(DOCKER_USER)/$(APP_NAME):$(DOCKER_TAG)
VERSION=dev

# Cores para output (Windows não suporta cores ANSI por padrão)
RED=
GREEN=
YELLOW=
BLUE=
NC=

.PHONY: help run build tidy test docker docker-run docker-push clean lint fmt

# Comando padrão
.DEFAULT_GOAL := help

help: ## Mostra esta mensagem de ajuda
	@echo "Code Explainer - Comandos disponíveis:"
	@echo "  run              - Executa a aplicação localmente"
	@echo "  build            - Compila a aplicação"
	@echo "  tidy             - Organiza dependências do Go"
	@echo "  test             - Executa os testes"
	@echo "  test-coverage    - Executa testes com cobertura"
	@echo "  lint             - Executa linter"
	@echo "  fmt              - Formata o código"
	@echo "  docker           - Constrói a imagem Docker"
	@echo "  docker-run       - Executa o container Docker"
	@echo "  docker-run-dev   - Executa o container em modo desenvolvimento"
	@echo "  docker-push      - Faz push da imagem para Docker Hub"
	@echo "  docker-clean     - Remove imagens Docker não utilizadas"
	@echo "  clean            - Remove arquivos de build"
	@echo "  install-deps     - Instala dependências de desenvolvimento"
	@echo "  check-env        - Verifica se o arquivo .env existe"
	@echo "  dev              - Executa pipeline de desenvolvimento completo"
	@echo "  release          - Executa pipeline de release completo"

run: ## Executa a aplicação localmente
	@echo "Executando aplicação..."
	go run main.go

build: ## Compila a aplicação
	@echo "Compilando aplicação..."
	go build -ldflags="-X main.Version=$(VERSION)" -o $(APP_NAME) main.go
	@echo "Build concluído: $(APP_NAME)"

tidy: ## Organiza dependências do Go
	@echo "Organizando dependências..."
	go mod tidy
	go mod verify

test: ## Executa os testes
	@echo "Executando testes..."
	go test -v ./...

test-coverage: ## Executa testes com cobertura
	@echo "Executando testes com cobertura..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Cobertura salva em coverage.html"

lint: ## Executa linter
	@echo "Executando linter..."
	golangci-lint run

fmt: ## Formata o código
	@echo "Formatando código..."
	go fmt ./...

docker: ## Constrói a imagem Docker
	@echo "Construindo imagem Docker..."
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "Imagem construída: $(APP_NAME):$(VERSION)"

docker-run: ## Executa o container Docker
	@echo "Executando container..."
	docker run --rm -it --env-file .env $(APP_NAME):latest

docker-run-dev: ## Executa o container em modo desenvolvimento
	@echo "Executando container em modo desenvolvimento..."
	docker run --rm -it --env-file .env -v $(PWD):/app -w /app $(APP_NAME):latest

docker-tag: ## Tag da imagem para Docker Hub
	@echo "Tagging imagem..."
	docker tag $(APP_NAME):$(VERSION) $(IMAGE)

docker-push: docker docker-tag ## Faz push da imagem para Docker Hub
	@echo "Fazendo push para Docker Hub..."
	docker push $(IMAGE)
	@echo "Push concluído: $(IMAGE)"

docker-clean: ## Remove imagens Docker não utilizadas
	@echo "Limpando imagens Docker..."
	docker image prune -f
	docker system prune -f

clean: ## Remove arquivos de build
	@echo "Limpando arquivos de build..."
	del /f /q $(APP_NAME).exe 2>nul || echo "Arquivo não encontrado"
	del /f /q coverage.out 2>nul || echo "Arquivo não encontrado"
	del /f /q coverage.html 2>nul || echo "Arquivo não encontrado"
	@echo "Limpeza concluída"

install-deps: ## Instala dependências de desenvolvimento
	@echo "Instalando dependências..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Dependências instaladas"

check-env: ## Verifica se o arquivo .env existe
	@if not exist .env ( \
		echo "Arquivo .env não encontrado!" && \
		echo "Criando arquivo .env de exemplo..." && \
		echo MODEL_NAME=codellama > .env && \
		echo "Arquivo .env criado" \
	)

dev: check-env tidy fmt lint test ## Executa pipeline de desenvolvimento completo
	@echo "Pipeline de desenvolvimento concluído!"

release: check-env tidy fmt lint test build docker ## Executa pipeline de release completo
	@echo "Release preparado!"
