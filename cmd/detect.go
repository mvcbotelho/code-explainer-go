package cmd

import (
	"fmt"
	"strings"

	"github.com/mvcbotelho/code-explainer/openai"
	"github.com/spf13/cobra"
)

var (
	detectCodeInput   string
	detectFilePath    string
	detectInteractive bool
)

// detectCmd representa o comando detect
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detecta a linguagem de programação de um código",
	Long: `Detecta automaticamente a linguagem de programação de um trecho de código.

Linguagens suportadas:
• Go, Python, JavaScript, C, Java, PHP, Rust, C#

Você pode fornecer o código de três formas:
1. Via flag --code: code-explainer detect --code "func main() {}"
2. Via arquivo: code-explainer detect --file main.go
3. Interativo: code-explainer detect (digite o código e pressione Ctrl+D)

Exemplos:
  code-explainer detect --code "print('Hello World')"
  code-explainer detect --file script.py
  code-explainer detect --code "console.log('Hello')" --verbose`,
	RunE: runDetect,
}

func init() {
	rootCmd.AddCommand(detectCmd)

	// Flags específicas do comando detect
	detectCmd.Flags().StringVarP(&detectCodeInput, "code", "c", "", "Código para detectar linguagem")
	detectCmd.Flags().StringVarP(&detectFilePath, "file", "f", "", "Arquivo contendo o código")
	detectCmd.Flags().BoolVarP(&detectInteractive, "interactive", "i", false, "Modo interativo")

	// Marcar flags como mutuamente exclusivas
	detectCmd.MarkFlagsMutuallyExclusive("code", "file", "interactive")
}

func runDetect(cmd *cobra.Command, args []string) error {
	var code string
	var err error

	// Determinar a fonte do código
	switch {
	case detectCodeInput != "":
		code = detectCodeInput
		if verbose {
			fmt.Printf("📝 Usando código fornecido via flag\n")
		}

	case detectFilePath != "":
		code, err = readFile(detectFilePath)
		if err != nil {
			return fmt.Errorf("erro ao ler arquivo %s: %w", detectFilePath, err)
		}
		if verbose {
			fmt.Printf("📁 Lendo código do arquivo: %s\n", detectFilePath)
		}

	case detectInteractive || (detectCodeInput == "" && detectFilePath == ""):
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

	// Detectar linguagem
	detectedLang := openai.DetectLanguage(code)

	// Formatar saída
	outputText := formatDetectOutput(code, detectedLang)

	// Escrever saída
	if output != "" {
		err = writeToFile(output, outputText)
		if err != nil {
			return fmt.Errorf("erro ao escrever arquivo de saída: %w", err)
		}
		if verbose {
			fmt.Printf("💾 Resultado salvo em: %s\n", output)
		}
	} else {
		fmt.Println(outputText)
	}

	return nil
}

// formatDetectOutput formata a saída da detecção
func formatDetectOutput(code, language string) string {
	var output strings.Builder

	output.WriteString("🔍 Detecção de Linguagem\n")
	output.WriteString(strings.Repeat("=", 30) + "\n\n")

	output.WriteString("💻 **Código analisado:**\n")
	output.WriteString("```\n")
	output.WriteString(code)
	output.WriteString("\n```\n\n")

	output.WriteString("🎯 **Linguagem detectada:** ")
	if language == "linguagem desconhecida" {
		output.WriteString("❓ " + language)
	} else {
		output.WriteString("✅ " + language)
	}
	output.WriteString("\n\n")

	// Adicionar informações extras se verbose
	if verbose {
		output.WriteString("📊 **Informações:**\n")
		output.WriteString(fmt.Sprintf("• Tamanho do código: %d caracteres\n", len(code)))
		output.WriteString(fmt.Sprintf("• Linhas de código: %d\n", len(strings.Split(code, "\n"))))

		// Listar linguagens suportadas
		supportedLangs := openai.GetSupportedLanguages()
		output.WriteString("• Linguagens suportadas: " + strings.Join(supportedLangs, ", ") + "\n")
	}

	return output.String()
}
