package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Flags globais
	modelName string
	apiURL    string
	timeout   int
	verbose   bool
	output    string
	language  string
)

// rootCmd representa o comando base quando chamado sem subcomandos
var rootCmd = &cobra.Command{
	Use:   "code-explainer",
	Short: "Uma ferramenta CLI para explicar código usando IA local",
	Long: `Code Explainer é uma ferramenta CLI inteligente que usa IA local (Ollama) 
para explicar código de programação de forma clara e detalhada.

Características:
• Detecção automática de linguagem de programação
• Explicações em português
• Execução totalmente local (sem envio de dados para nuvem)
• Suporte a múltiplos modelos de IA
• Interface CLI intuitiva

Exemplos:
  code-explainer explain --file main.go
  code-explainer explain --code "func main() { fmt.Println('Hello') }"
  code-explainer detect --file script.py
  code-explainer list-models`,
	Version: "1.0.0",
}

// Execute adiciona todos os comandos filhos ao comando root e define flags
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Flags globais
	rootCmd.PersistentFlags().StringVarP(&modelName, "model", "m", getEnvOrDefault("MODEL_NAME", "codellama"), "Modelo de IA a ser usado")
	rootCmd.PersistentFlags().StringVarP(&apiURL, "api-url", "u", getEnvOrDefault("OLLAMA_API_URL", "http://localhost:11434/api/generate"), "URL da API Ollama")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", getEnvIntOrDefault("REQUEST_TIMEOUT", 30), "Timeout em segundos para requisições")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Modo verboso")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Arquivo de saída (padrão: stdout)")
	rootCmd.PersistentFlags().StringVarP(&language, "language", "l", "", "Forçar linguagem específica (opcional)")
}

// getEnvOrDefault retorna o valor da variável de ambiente ou o valor padrão
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault retorna o valor inteiro da variável de ambiente ou o valor padrão
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := fmt.Sscanf(value, "%d", &defaultValue); err == nil && intValue == 1 {
			return defaultValue
		}
	}
	return defaultValue
}
