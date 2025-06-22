package main

import (
	"fmt"
	"os"

	"github.com/mvcbotelho/code-explainer/cmd"
)

func main() {
	// Configurar tratamento de erros
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "❌ Erro fatal: %v\n", r)
			os.Exit(1)
		}
	}()

	// Executar CLI
	cmd.Execute()
}
