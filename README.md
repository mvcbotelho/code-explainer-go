# 🚀 Code Explainer

Uma ferramenta CLI inteligente para explicar código usando IA local via Ollama. Desenvolvida em Go com arquitetura modular e extensível.

## ✨ Características

- 🤖 **IA Local**: Usa Ollama com modelos como CodeLlama, sem envio de dados para nuvem
- 🔍 **Detecção Automática**: Identifica automaticamente a linguagem de programação
- 🐳 **Docker Otimizado**: Multi-stage build com imagem minimalista (~41MB)
- 🧪 **Testes Completos**: Cobertura abrangente com mocks HTTP
- 🔧 **Configurável**: Suporte a diferentes modelos e configurações via variáveis de ambiente
- 🛡️ **Seguro**: Container com usuário não-root e healthcheck
- 📚 **Documentado**: Documentação completa e exemplos

## 🚀 Instalação e Uso

### Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) instalado
- [Ollama](https://ollama.com) rodando com modelo `codellama`

### Método 1: Docker (Recomendado)

```bash
# Baixar e executar
docker run --rm -it mvcbotelho/code-explainer:latest

# Ou com modelo customizado
docker run --rm -it -e MODEL_NAME=gpt-3.5-turbo mvcbotelho/code-explainer:latest
```

### Método 2: Build Local

```bash
# Clone o repositório
git clone https://github.com/mvcbotelho/code-explainer.git
cd code-explainer

# Build da imagem
docker build -t code-explainer .

# Executar
docker run --rm -it code-explainer
```

### Método 3: Desenvolvimento Local

```bash
# Instalar Go 1.22+
go mod download
go run main.go
```

## 🐳 Docker Compose

Para desenvolvimento com hot reload:

```bash
# Serviço principal
docker-compose up code-explainer

# Modo desenvolvimento
docker-compose --profile dev up code-explainer-dev
```

## 🛠️ Comandos de Desenvolvimento

### Makefile (Linux/macOS)

```bash
# Ver todos os comandos disponíveis
make help

# Pipeline de desenvolvimento completo
make dev

# Executar testes
make test

# Build da aplicação
make build

# Build Docker
make docker

# Pipeline de release
make release
```

### Windows (PowerShell)

```bash
# Executar aplicação
go run main.go

# Executar testes
go test ./openai

# Build
go build -o code-explainer.exe

# Docker build
docker build -t code-explainer .
```

## 📝 Exemplo de Uso

```bash
$ docker run --rm -it mvcbotelho/code-explainer

Cole o trecho de código abaixo e pressione Ctrl+D (Linux/macOS) ou Ctrl+Z (Windows) para enviar:

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

📘 Explicação gerada pela IA:
Este código implementa a função de Fibonacci em Go. A função recebe um número inteiro n e retorna o n-ésimo número da sequência de Fibonacci. A implementação usa recursão: se n for 0 ou 1, retorna n; caso contrário, retorna a soma dos dois números anteriores da sequência.
```

## ⚙️ Configuração

### Variáveis de Ambiente

Crie um arquivo `.env`:

```bash
# Modelo de IA (padrão: codellama)
MODEL_NAME=codellama

# URL da API Ollama (padrão: http://localhost:11434/api/generate)
OLLAMA_API_URL=http://localhost:11434/api/generate

# Timeout para requisições (padrão: 30s)
REQUEST_TIMEOUT=30
```

### Modelos Suportados

- `codellama` (padrão)
- `gpt-3.5-turbo`
- `llama2`
- Qualquer modelo disponível no Ollama

## 🧪 Testes

```bash
# Executar todos os testes
go test ./openai

# Testes com cobertura
go test -cover ./openai

# Testes específicos
go test -run TestDetectLanguage ./openai
```

### Cobertura de Testes

- ✅ Detecção de linguagens (Go, Python, JavaScript, C, Java, PHP, Rust, C#)
- ✅ Integração com API Ollama
- ✅ Tratamento de erros HTTP
- ✅ Configurações customizadas
- ✅ Estruturas de dados Request/Response

## 🏗️ Arquitetura

```
code-explainer/
├── main.go              # Ponto de entrada da aplicação
├── openai/
│   ├── explain.go       # Integração com API Ollama
│   ├── language.go      # Detecção de linguagens
│   ├── explain_test.go  # Testes de integração
│   └── language_test.go # Testes de detecção
├── Dockerfile           # Multi-stage build otimizado
├── docker-compose.yml   # Orquestração de containers
├── Makefile            # Pipeline de desenvolvimento
└── DOCKER.md           # Documentação Docker
```

### Componentes Principais

- **DetectLanguage**: Usa expressões regulares para identificar linguagens
- **ExplainCode**: Envia código para análise via API Ollama
- **Config**: Estrutura para configurações customizáveis
- **APIError**: Tratamento específico de erros da API

## 🔧 Desenvolvimento

### Estrutura de Código

- **Modular**: Separação clara entre detecção e explicação
- **Testável**: Mocks HTTP para testes isolados
- **Extensível**: Fácil adição de novas linguagens
- **Configurável**: Injeção de dependências para URLs e modelos

### Adicionando Novas Linguagens

```go
// Em openai/language.go
patterns := []LanguagePattern{
    {
        Language: "NovaLinguagem",
        Priority: 85,
        Patterns: []*regexp.Regexp{
            regexp.MustCompile(`\bkeyword\b`),
            regexp.MustCompile(`specific_pattern`),
        },
    },
}
```

## 🐳 Docker

### Imagem Otimizada

- **Multi-stage build** para reduzir tamanho
- **Alpine 3.19** como base minimalista
- **Usuário não-root** para segurança
- **Healthcheck** para monitoramento
- **Labels** para metadados

### Comandos Docker

```bash
# Build
docker build -t code-explainer .

# Executar
docker run --rm -it code-explainer

# Desenvolvimento com volume
docker run --rm -it -v $(pwd):/app code-explainer

# Ver logs
docker logs code-explainer
```

## 📊 Monitoramento

### Healthcheck

```bash
# Verificar status
docker inspect --format='{{.State.Health.Status}}' code-explainer
```

### Logs

```bash
# Logs em tempo real
docker logs -f code-explainer

# Logs com timestamp
docker logs -t code-explainer
```

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Padrões de Commit

```
feat: nova funcionalidade
fix: correção de bug
docs: documentação
style: formatação
refactor: refatoração
test: testes
chore: tarefas de manutenção
```

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👨‍💻 Autor

**Marcus Botelho**

- GitHub: [@mvcbotelho](https://github.com/mvcbotelho)
- LinkedIn: [Marcus Botelho](https://linkedin.com/in/mvcbotelho)

---

⭐ Se este projeto te ajudou, considere dar uma estrela no repositório!

# Ler código de arquivo
code-explainer explain --file main.go

# Explicar múltiplos arquivos
code-explainer explain --dir ./src

# Salvar saída em arquivo
code-explainer explain --output explanation.md
