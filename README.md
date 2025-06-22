# 🐳 Code Explainer - Docker Edition

Este projeto permite explicar trechos de código usando IA (modelo CodeLlama) de forma **totalmente local** via **Docker + Ollama**.

---

## 🚀 Como usar com Docker

### 1. Certifique-se de que o [Ollama](https://ollama.com) está instalado

Baixe o modelo `codellama` e inicie:

```bash
ollama run codellama
```

---

### 2. Baixe a imagem do Docker Hub

```bash
docker pull mvcbotelho/code-explainer:latest
```

> Substitua `mvcbotelho` pelo seu nome de usuário do Docker Hub, se for diferente.

---

### 3. Execute a imagem

```bash
docker run --rm -it mvcbotelho/code-explainer
```

Cole o código no terminal e pressione **Ctrl+D** (Linux/macOS) ou **Ctrl+Z** (Windows) para enviar.

---

## 💻 Exemplo

```text
$ docker run --rm -it mvcbotelho/code-explainer

Cole o trecho de código abaixo e pressione Ctrl+D (Linux/macOS) ou Ctrl+Z (Windows) para enviar:

func soma(a int, b int) int {
    return a + b
}

📘 Explicação gerada pela IA:
Esta função soma dois números inteiros e retorna o resultado.
```

---

## 📦 Requisitos

- Docker instalado
- Ollama rodando com o modelo `codellama`

---

## 🧠 Sobre o projeto

- Linguagem: Go
- Integração com IA local via Ollama
- Não envia dados para a nuvem
- Ideal para estudos, reviews e aprendizado

---

Feito com 💡 por Marcus
