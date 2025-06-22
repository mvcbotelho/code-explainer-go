# Testes Unitários - Code Explainer

Este diretório contém os testes unitários para as funções de detecção de linguagem e explicação de código.

## Arquivos de Teste

### `language_test.go`
Testes para a função `DetectLanguage()` que identifica a linguagem de programação do código.

**Cobertura de Testes:**
- ✅ Detecção de Go (func, package, import)
- ✅ Detecção de Python (def, print, import)
- ✅ Detecção de JavaScript (console.log, function, var, let)
- ✅ Detecção de C (#include, printf, int main)
- ✅ Detecção de Java (System.out.println, public class, public static)
- ✅ Detecção de PHP (echo, $_post, <?php)
- ✅ Detecção de Rust (fn, let, println!)
- ✅ Detecção de C# (class + public + {)
- ✅ Case insensitive (maiúsculas/minúsculas)
- ✅ Prioridade de detecção
- ✅ Linguagem desconhecida

### `explain_test.go`
Testes para a função `ExplainCode()` que envia código para análise via API.

**Cobertura de Testes:**
- ✅ Requisição HTTP correta (POST, Content-Type)
- ✅ Estrutura de dados Request/Response
- ✅ Modelo padrão (codellama)
- ✅ Modelo customizado via variável de ambiente
- ✅ Tratamento de erros HTTP
- ✅ Tratamento de JSON inválido
- ✅ Resposta vazia
- ✅ Diferentes tipos de código (Go, Python, JavaScript)
- ✅ Serialização/deserialização JSON

## Como Executar os Testes

### Executar todos os testes
```bash
go test ./openai
```

### Executar testes específicos
```bash
# Apenas testes de detecção de linguagem
go test -run TestDetectLanguage ./openai

# Apenas testes de explicação de código
go test -run TestExplainCode ./openai
```

### Executar com verbose
```bash
go test -v ./openai
```

### Executar com cobertura
```bash
go test -cover ./openai
```

### Executar com relatório de cobertura detalhado
```bash
go test -coverprofile=coverage.out ./openai
go tool cover -html=coverage.out
```

## Estrutura dos Testes

### Testes de Detecção de Linguagem
- **TestDetectLanguage**: Testa todos os casos de uso com diferentes linguagens
- **TestDetectLanguageCaseInsensitive**: Verifica se a detecção é case insensitive
- **TestDetectLanguagePriority**: Testa a prioridade das detecções

### Testes de Explicação de Código
- **TestExplainCode**: Teste básico com servidor mock
- **TestExplainCodeWithCustomModel**: Teste com modelo customizado
- **TestExplainCodeHTTPError**: Teste de erro HTTP
- **TestExplainCodeInvalidJSON**: Teste de JSON inválido
- **TestExplainCodeEmptyResponse**: Teste de resposta vazia
- **TestRequestStruct**: Teste da estrutura Request
- **TestResponseStruct**: Teste da estrutura Response
- **TestExplainCodeWithDifferentLanguages**: Teste com diferentes linguagens

## Notas Importantes

1. **Mocks HTTP**: Os testes usam `httptest.NewServer` para simular respostas da API sem fazer chamadas reais
2. **Variáveis de Ambiente**: Alguns testes configuram variáveis de ambiente para testar funcionalidades específicas
3. **Limpeza**: Os testes fazem cleanup adequado com `defer server.Close()`
4. **Isolamento**: Cada teste é independente e não afeta outros testes

## Melhorias Futuras

- [ ] Adicionar testes de benchmark
- [ ] Testes de concorrência
- [ ] Testes de timeout
- [ ] Testes de rate limiting
- [ ] Testes de diferentes formatos de resposta da API 