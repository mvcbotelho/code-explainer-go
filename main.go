package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/mvcbotelho/code-explainer/openai"
)

func main() {
    fmt.Println("Cole o trecho de código abaixo e pressione Ctrl+D (Linux/macOS) ou Ctrl+Z (Windows) para enviar:")

    scanner := bufio.NewScanner(os.Stdin)
    var code strings.Builder
    for scanner.Scan() {
        code.WriteString(scanner.Text() + "\n")
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("Erro ao ler entrada: %v", err)
    }

    explanation, err := openai.ExplainCode(code.String())
    if err != nil {
        log.Fatalf("Erro ao explicar o código: %v", err)
    }

    fmt.Println("\n📘 Explicação gerada pela IA:")
    fmt.Println(explanation)
}
