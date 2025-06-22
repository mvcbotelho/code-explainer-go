# 🤖 Code Explainer

Um utilitário de linha de comando em Go que usa **modelos locais de IA com Ollama** para explicar trechos de código.

---

## 🚀 Como usar com Ollama

### 1. Instale o Ollama

Acesse [https://ollama.com](https://ollama.com) e instale a versão para seu sistema operacional.

Ou via terminal:

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

### 2. Baixe e rode o modelo `codellama`

```bash
ollama run codellama
```

> Isso abrirá um servidor local em `http://localhost:11434`

### 3. Rode o projeto

```bash
go run main.go
```

Cole o trecho de código e pressione **Ctrl+D** (Linux/macOS) ou **Ctrl+Z** (Windows) para enviar.

---

## 💻 Exemplo no terminal

```bash
$ go run main.go
Cole o trecho de código abaixo e pressione Ctrl+D (Linux/macOS) ou Ctrl+Z (Windows) para enviar:

func soma(a int, b int) int {
    return a + b
}

📘 Explicação gerada pela IA:
Esta função recebe dois números inteiros como argumentos e retorna a soma deles.

```

---

## 📦 Estrutura do projeto

```
code-explainer/
├── main.go
├── openai/
│   └── explain.go
├── go.mod
├── .gitignore
└── README.md
```

---

## 📋 Requisitos

- Go 1.20+
- Ollama instalado
- Modelo `codellama` carregado

---

## 📜 Licença

MIT License

Feito com 💡 por Marcus e 🤖 R2Dev2
