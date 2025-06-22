package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mvcbotelho/code-explainer/openai"
	"github.com/spf13/cobra"
)

var (
	codeInput   string
	filePath    string
	interactive bool
)

// explainCmd representa o comando explain
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explica um trecho de código usando IA",
	Long: `Explica um trecho de código usando IA local (Ollama).

Você pode fornecer o código de três formas:
1. Via flag --code: code-explainer explain --code "func main() {}"
2. Via arquivo: code-explainer explain --file main.go
3. Interativo: code-explainer explain (digite o código e pressione Ctrl+D)

Exemplos:
  code-explainer explain --code "print('Hello World')"
  code-explainer explain --file main.go
  code-explainer explain --file main.go --output explanation.md
  code-explainer explain --model gpt-3.5-turbo --code "console.log('Hello')"`,
	RunE: runExplain,
}

func init() {
	rootCmd.AddCommand(explainCmd)

	// Flags específicas do comando explain
	explainCmd.Flags().StringVarP(&codeInput, "code", "c", "", "Código a ser explicado")
	explainCmd.Flags().StringVarP(&filePath, "file", "f", "", "Arquivo contendo o código")
	explainCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Modo interativo (padrão se nenhuma entrada for fornecida)")

	// Marcar flags como mutuamente exclusivas
	explainCmd.MarkFlagsMutuallyExclusive("code", "file", "interactive")
}

func runExplain(cmd *cobra.Command, args []string) error {
	var code string
	var err error

	// Determinar a fonte do código
	switch {
	case codeInput != "":
		code = codeInput
		if verbose {
			fmt.Printf("📝 Usando código fornecido via flag\n")
		}

	case filePath != "":
		code, err = readFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
		}
		if verbose {
			fmt.Printf("📁 Lendo código do arquivo: %s\n", filePath)
		}

	case interactive || (codeInput == "" && filePath == ""):
		code, err = readInteractive()
		if err != nil {
			return fmt.Errorf("erro ao ler entrada interativa: %w", err)
		}
		if verbose {
			fmt.Printf("⌨️  Usando entrada interativa\n")
		}

	default:
		return fmt.Errorf("forneça o código via --code, --file ou use modo interativo")
	}

	if code == "" {
		return fmt.Errorf("código vazio fornecido")
	}

	// Detectar linguagem se não for forçada
	detectedLang := language
	if detectedLang == "" {
		detectedLang = openai.DetectLanguage(code)
		if verbose {
			fmt.Printf("🔍 Linguagem detectada: %s\n", detectedLang)
		}
	}

	// Configurar cliente
	config := &openai.Config{
		APIURL:  apiURL,
		Model:   modelName,
		Timeout: time.Duration(timeout) * time.Second,
	}

	if verbose {
		fmt.Printf("🤖 Usando modelo: %s\n", config.Model)
		fmt.Printf("🌐 API URL: %s\n", config.APIURL)
		fmt.Printf("⏱️  Timeout: %ds\n", timeout)
		fmt.Printf("📊 Tamanho do código: %d caracteres\n", len(code))
		fmt.Println("🔄 Enviando para análise...")
	}

	// Explicar código
	explanation, err := openai.ExplainCode(code, config)
	if err != nil {
		return fmt.Errorf("erro ao explicar código: %w", err)
	}

	// Formatar saída
	outputText := formatOutput(code, detectedLang, explanation)

	// Escrever saída
	if output != "" {
		err = writeToFile(output, outputText)
		if err != nil {
			return fmt.Errorf("erro ao escrever arquivo de saída: %w", err)
		}
		if verbose {
			fmt.Printf("💾 Explicação salva em: %s\n", output)
		}
	} else {
		fmt.Println(outputText)
	}

	return nil
}

// readFile lê o conteúdo de um arquivo
func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// readInteractive lê código da entrada padrão
func readInteractive() (string, error) {
	fmt.Println("Cole o trecho de código abaixo e pressione Ctrl+D (Linux/macOS) ou Ctrl+Z (Windows) para enviar:")

	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

// formatOutput formata a saída da explicação
func formatOutput(code, language, explanation string) string {
	var output strings.Builder

	output.WriteString("📘 Explicação gerada pela IA:\n")
	output.WriteString(strings.Repeat("=", 50) + "\n\n")

	if language != "" && language != "linguagem desconhecida" {
		output.WriteString(fmt.Sprintf("🔍 **Linguagem detectada:** %s\n\n", language))
	}

	output.WriteString("💻 **Código analisado:**\n")
	output.WriteString("```\n")
	output.WriteString(code)
	output.WriteString("\n```\n\n")

	output.WriteString("🤖 **Explicação:**\n")
	output.WriteString(explanation)
	output.WriteString("\n")

	return output.String()
}

// writeToFile escreve conteúdo em um arquivo
func writeToFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
