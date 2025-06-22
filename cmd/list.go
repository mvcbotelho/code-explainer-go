package cmd

import (
	"fmt"
	"strings"

	"github.com/mvcbotelho/code-explainer/openai"
	"github.com/spf13/cobra"
)

// listCmd representa o comando list
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista informações sobre modelos e linguagens suportadas",
	Long: `Lista informações sobre modelos de IA disponíveis e linguagens suportadas.

Subcomandos:
  models     - Lista modelos de IA recomendados
  languages  - Lista linguagens de programação suportadas
  config     - Mostra configuração atual

Exemplos:
  code-explainer list models
  code-explainer list languages
  code-explainer list config`,
}

// listModelsCmd lista os modelos de IA
var listModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Lista modelos de IA recomendados",
	Long: `Lista os modelos de IA recomendados para uso com o Code Explainer.

Estes modelos são compatíveis com Ollama e podem ser usados
para explicar código de programação.`,
	Run: runListModels,
}

// listLanguagesCmd lista as linguagens suportadas
var listLanguagesCmd = &cobra.Command{
	Use:   "languages",
	Short: "Lista linguagens de programação suportadas",
	Long: `Lista todas as linguagens de programação que podem ser detectadas
automaticamente pelo Code Explainer.`,
	Run: runListLanguages,
}

// listConfigCmd mostra a configuração atual
var listConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Mostra a configuração atual",
	Long: `Mostra a configuração atual do Code Explainer, incluindo
modelo padrão, URL da API e outras configurações.`,
	Run: runListConfig,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listModelsCmd)
	listCmd.AddCommand(listLanguagesCmd)
	listCmd.AddCommand(listConfigCmd)
}

func runListModels(cmd *cobra.Command, args []string) {
	fmt.Println("🤖 Modelos de IA Recomendados")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()

	models := []struct {
		Name        string
		Description string
		Size        string
		BestFor     string
	}{
		{
			Name:        "codellama",
			Description: "Modelo especializado em código, baseado no Llama 2",
			Size:        "~4GB",
			BestFor:     "Explicação de código, análise de algoritmos",
		},
		{
			Name:        "codellama:7b",
			Description: "Versão menor do CodeLlama, mais rápida",
			Size:        "~4GB",
			BestFor:     "Desenvolvimento rápido, recursos limitados",
		},
		{
			Name:        "codellama:13b",
			Description: "Versão maior do CodeLlama, mais precisa",
			Size:        "~8GB",
			BestFor:     "Análises complexas, alta precisão",
		},
		{
			Name:        "llama2",
			Description: "Modelo geral, bom para código e texto",
			Size:        "~4GB",
			BestFor:     "Uso geral, documentação",
		},
		{
			Name:        "gpt-3.5-turbo",
			Description: "Modelo OpenAI (requer API key)",
			Size:        "N/A",
			BestFor:     "Alta qualidade, uso comercial",
		},
	}

	for i, model := range models {
		fmt.Printf("%d. **%s** (%s)\n", i+1, model.Name, model.Size)
		fmt.Printf("   %s\n", model.Description)
		fmt.Printf("   Melhor para: %s\n", model.BestFor)
		fmt.Println()
	}

	fmt.Println("💡 **Dica:** Use 'ollama list' para ver modelos instalados localmente")
	fmt.Println("📥 **Instalar:** ollama pull codellama")
}

func runListLanguages(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Linguagens de Programação Suportadas")
	fmt.Println(strings.Repeat("=", 45))
	fmt.Println()

	languages := openai.GetSupportedLanguages()

	for i, lang := range languages {
		icon := getLanguageIcon(lang)
		fmt.Printf("%d. %s %s\n", i+1, icon, lang)
	}

	fmt.Println()
	fmt.Println("💡 **Dica:** A detecção é automática, mas você pode forçar uma linguagem com --language")
	fmt.Println("📝 **Exemplo:** code-explainer explain --language Python --code 'print(\"Hello\")'")
}

func runListConfig(cmd *cobra.Command, args []string) {
	fmt.Println("⚙️  Configuração Atual")
	fmt.Println(strings.Repeat("=", 25))
	fmt.Println()

	config := openai.DefaultConfig()

	fmt.Printf("🤖 **Modelo padrão:** %s\n", config.Model)
	fmt.Printf("🌐 **URL da API:** %s\n", config.APIURL)
	fmt.Printf("⏱️  **Timeout:** %v\n", config.Timeout)
	fmt.Println()

	fmt.Printf("🔧 **Configuração atual:**\n")
	fmt.Printf("   Modelo: %s\n", modelName)
	fmt.Printf("   API URL: %s\n", apiURL)
	fmt.Printf("   Timeout: %ds\n", timeout)
	fmt.Printf("   Verbose: %t\n", verbose)
	fmt.Printf("   Output: %s\n", getOutputDisplay())
	fmt.Printf("   Language: %s\n", getLanguageDisplay())
	fmt.Println()

	fmt.Println("💡 **Dicas:**")
	fmt.Println("   • Use --model para alterar o modelo")
	fmt.Println("   • Use --api-url para conectar a um servidor diferente")
	fmt.Println("   • Use --verbose para mais informações")
	fmt.Println("   • Configure variáveis de ambiente: MODEL_NAME, OLLAMA_API_URL, REQUEST_TIMEOUT")
}

// getLanguageIcon retorna um emoji para cada linguagem
func getLanguageIcon(lang string) string {
	icons := map[string]string{
		"Go":         "🐹",
		"Python":     "🐍",
		"JavaScript": "🟨",
		"C":          "🔵",
		"Java":       "☕",
		"PHP":        "🐘",
		"Rust":       "🦀",
		"C#":         "💜",
	}

	if icon, exists := icons[lang]; exists {
		return icon
	}
	return "📄"
}

// getOutputDisplay retorna a exibição do output
func getOutputDisplay() string {
	if output == "" {
		return "stdout (padrão)"
	}
	return output
}

// getLanguageDisplay retorna a exibição da linguagem
func getLanguageDisplay() string {
	if language == "" {
		return "detecção automática (padrão)"
	}
	return language
}
