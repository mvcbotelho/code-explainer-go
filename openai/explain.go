package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Config contém as configurações para a API
type Config struct {
	APIURL  string
	Model   string
	Timeout time.Duration
}

// DefaultConfig retorna uma configuração padrão
func DefaultConfig() *Config {
	return &Config{
		APIURL:  "http://localhost:11434/api/generate",
		Model:   "codellama",
		Timeout: 30 * time.Second,
	}
}

// Request representa a requisição para a API
type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Response representa a resposta da API
type Response struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// APIError representa um erro específico da API
type APIError struct {
	StatusCode int
	Message    string
	Details    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// ExplainCode envia código para análise via API com configuração customizável
func ExplainCode(code string, config *Config) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Usa variável de ambiente se disponível
	if model := os.Getenv("MODEL_NAME"); model != "" {
		config.Model = model
	}

	lang := DetectLanguage(code)
	prompt := fmt.Sprintf(`Explique o que o seguinte código em %s faz:

%s`, lang, code)

	body := Request{
		Model:  config.Model,
		Prompt: prompt,
		Stream: false,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return "", fmt.Errorf("erro ao codificar requisição: %w", err)
	}

	// Cria cliente HTTP com timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	resp, err := client.Post(config.APIURL, "application/json", buf)
	if err != nil {
		return "", fmt.Errorf("erro de conexão com a API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Tenta ler o corpo da resposta para mais detalhes
		var errorBody bytes.Buffer
		errorBody.ReadFrom(resp.Body)

		return "", &APIError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
			Details:    errorBody.String(),
		}
	}

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta da API: %w", err)
	}

	return r.Response, nil
}

// ExplainCodeWithDefaultURL é uma função de conveniência que usa a URL padrão
func ExplainCodeWithDefaultURL(code string) (string, error) {
	return ExplainCode(code, DefaultConfig())
}

// ValidateConfig valida se a configuração está correta
func ValidateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("configuração não pode ser nula")
	}

	if config.APIURL == "" {
		return fmt.Errorf("URL da API não pode estar vazia")
	}

	if config.Model == "" {
		return fmt.Errorf("modelo não pode estar vazio")
	}

	if config.Timeout <= 0 {
		return fmt.Errorf("timeout deve ser maior que zero")
	}

	return nil
}

// GetDefaultAPIURL retorna a URL padrão da API
func GetDefaultAPIURL() string {
	return DefaultConfig().APIURL
}

// GetDefaultModel retorna o modelo padrão
func GetDefaultModel() string {
	return DefaultConfig().Model
}
