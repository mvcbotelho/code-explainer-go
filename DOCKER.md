# Docker - Code Explainer

Este documento descreve as configurações Docker e como usar os containers.

## 🐳 Arquivos Docker

### Dockerfile
- **Multi-stage build** otimizado para reduzir tamanho da imagem
- **Segurança** com usuário não-root
- **Healthcheck** para monitoramento
- **Labels** para metadados
- **Cache otimizado** para builds mais rápidos

### docker-compose.yml
- **Serviço principal** para produção
- **Serviço de desenvolvimento** com hot reload
- **Volumes** para persistência e desenvolvimento
- **Healthcheck** configurado
- **Variáveis de ambiente** suportadas

### .dockerignore
- **Otimizado** para excluir arquivos desnecessários
- **Cobertura completa** de arquivos de desenvolvimento
- **Reduz tamanho** do contexto de build

## 🚀 Comandos Docker

### Build e Execução Básica
```bash
# Build da imagem
make docker

# Executar container
make docker-run

# Executar em modo desenvolvimento
make docker-run-dev
```

### Docker Compose
```bash
# Executar serviço principal
docker-compose up code-explainer

# Executar em modo desenvolvimento
docker-compose --profile dev up code-explainer-dev

# Executar em background
docker-compose up -d code-explainer

# Parar serviços
docker-compose down
```

### Comandos Avançados
```bash
# Build com tag específica
docker build -t code-explainer:v1.0.0 .

# Executar com variáveis de ambiente
docker run -e MODEL_NAME=gpt-3.5-turbo code-explainer

# Executar com volume para desenvolvimento
docker run -v $(pwd):/app code-explainer

# Inspecionar container
docker inspect code-explainer

# Ver logs
docker logs code-explainer
```

## 🔧 Configurações

### Variáveis de Ambiente
```bash
# Arquivo .env
MODEL_NAME=codellama
OLLAMA_API_URL=http://localhost:11434/api/generate
REQUEST_TIMEOUT=30
```

### Volumes
- `/app` - Diretório da aplicação
- `/go` - Cache do Go (desenvolvimento)

### Portas
- Nenhuma porta exposta (aplicação CLI)

## 🛡️ Segurança

### Usuário Não-Root
- Container roda como usuário `appuser` (UID 1001)
- Permissões mínimas necessárias
- Isolamento de processos

### Healthcheck
- Verifica se o processo está rodando
- Intervalo: 30s
- Timeout: 3s
- Retries: 3

## 📊 Monitoramento

### Healthcheck
```bash
# Verificar status do healthcheck
docker inspect --format='{{.State.Health.Status}}' code-explainer
```

### Logs
```bash
# Ver logs em tempo real
docker logs -f code-explainer

# Ver logs com timestamp
docker logs -t code-explainer
```

## 🔄 CI/CD

### Build Otimizado
```bash
# Build para produção
docker build --target production .

# Build para desenvolvimento
docker build --target builder .
```

### Push para Registry
```bash
# Tag e push
make docker-push

# Push com versão específica
docker tag code-explainer:latest mvcbotelho/code-explainer:v1.0.0
docker push mvcbotelho/code-explainer:v1.0.0
```

## 🧹 Limpeza

### Remover Containers/Imagens
```bash
# Limpar containers parados
docker container prune

# Limpar imagens não utilizadas
make docker-clean

# Limpeza completa
docker system prune -a
```

## 🐛 Troubleshooting

### Problemas Comuns

1. **Container não inicia**
   ```bash
   # Verificar logs
   docker logs code-explainer
   
   # Verificar healthcheck
   docker inspect code-explainer
   ```

2. **Permissões de arquivo**
   ```bash
   # Corrigir permissões
   docker run --user root -v $(pwd):/app code-explainer chown -R appuser:appgroup /app
   ```

3. **Problemas de rede**
   ```bash
   # Verificar conectividade
   docker exec code-explainer ping localhost
   ```

## 📈 Melhorias Implementadas

### ✅ Segurança
- [x] Usuário não-root
- [x] Imagem Alpine minimalista
- [x] Healthcheck configurado
- [x] Labels de metadados

### ✅ Performance
- [x] Multi-stage build
- [x] Cache otimizado
- [x] .dockerignore completo
- [x] Imagem final minimalista

### ✅ Desenvolvimento
- [x] Docker Compose configurado
- [x] Volumes para hot reload
- [x] Comandos Makefile melhorados
- [x] Documentação completa

### ✅ Operação
- [x] Healthcheck funcional
- [x] Logs estruturados
- [x] Variáveis de ambiente
- [x] Comandos de limpeza 