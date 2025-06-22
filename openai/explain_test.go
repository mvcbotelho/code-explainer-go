package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// Mock da função ExplainCode para testes
func mockExplainCode(code string, serverURL string) (string, error) {
	config := &Config{
		APIURL:  serverURL,
		Model:   "codellama",
		Timeout: 5 * time.Second,
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

	client := &http.Client{
		Timeout: config.Timeout,
	}

	resp, err := client.Post(config.APIURL, "application/json", buf)
	if err != nil {
		return "", fmt.Errorf("erro de conexão com a API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

func TestExplainCode(t *testing.T) {
	// Teste com servidor mock
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica se é uma requisição POST
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verifica o Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Decodifica o corpo da requisição
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verifica se o modelo está correto
		if req.Model != "codellama" {
			t.Errorf("Expected model 'codellama', got %s", req.Model)
		}

		// Verifica se o stream está false
		if req.Stream {
			t.Errorf("Expected stream to be false, got %v", req.Stream)
		}

		// Verifica se o prompt contém a linguagem detectada
		if req.Prompt == "" {
			t.Errorf("Expected non-empty prompt")
		}

		// Resposta mock
		response := Response{
			Response: "Este código implementa uma função que calcula a soma de dois números.",
			Done:     true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	code := `func add(a, b int) int {
    return a + b
}`

	result, err := mockExplainCode(code, server.URL)
	if err != nil {
		t.Errorf("mockExplainCode() error = %v", err)
		return
	}

	if result == "" {
		t.Errorf("mockExplainCode() returned empty result")
	}
}

func TestExplainCodeWithCustomModel(t *testing.T) {
	// Teste com modelo customizado via variável de ambiente
	os.Setenv("MODEL_NAME", "custom-model")
	defer os.Unsetenv("MODEL_NAME")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verifica se o modelo customizado foi usado
		if req.Model != "custom-model" {
			t.Errorf("Expected model 'custom-model', got %s", req.Model)
		}

		response := Response{
			Response: "Código analisado com modelo customizado.",
			Done:     true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	code := `print("Hello World")`

	result, err := mockExplainCode(code, server.URL)
	if err != nil {
		t.Errorf("mockExplainCode() error = %v", err)
		return
	}

	if result == "" {
		t.Errorf("mockExplainCode() returned empty result")
	}
}

func TestExplainCodeHTTPError(t *testing.T) {
	// Teste com erro HTTP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	code := `func test() {}`

	_, err := mockExplainCode(code, server.URL)
	if err == nil {
		t.Errorf("mockExplainCode() should return error for HTTP 500")
	}

	// Verifica se é um APIError
	if apiErr, ok := err.(*APIError); ok {
		if apiErr.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
		}
	} else {
		t.Errorf("Expected APIError, got %T", err)
	}
}

func TestExplainCodeInvalidJSON(t *testing.T) {
	// Teste com resposta JSON inválida
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	code := `func test() {}`

	_, err := mockExplainCode(code, server.URL)
	if err == nil {
		t.Errorf("mockExplainCode() should return error for invalid JSON")
	}
}

func TestExplainCodeEmptyResponse(t *testing.T) {
	// Teste com resposta vazia
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Response: "",
			Done:     true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	code := `func test() {}`

	result, err := mockExplainCode(code, server.URL)
	if err != nil {
		t.Errorf("mockExplainCode() error = %v", err)
		return
	}

	if result != "" {
		t.Errorf("mockExplainCode() should return empty string for empty response")
	}
}

func TestConfigValidation(t *testing.T) {
	// Teste de validação de configuração
	tests := []struct {
		name   string
		config *Config
		valid  bool
	}{
		{
			name: "Config válida",
			config: &Config{
				APIURL:  "http://localhost:8080",
				Model:   "test-model",
				Timeout: 30 * time.Second,
			},
			valid: true,
		},
		{
			name:   "Config nula",
			config: nil,
			valid:  false,
		},
		{
			name: "URL vazia",
			config: &Config{
				APIURL:  "",
				Model:   "test-model",
				Timeout: 30 * time.Second,
			},
			valid: false,
		},
		{
			name: "Modelo vazio",
			config: &Config{
				APIURL:  "http://localhost:8080",
				Model:   "",
				Timeout: 30 * time.Second,
			},
			valid: false,
		},
		{
			name: "Timeout zero",
			config: &Config{
				APIURL:  "http://localhost:8080",
				Model:   "test-model",
				Timeout: 0,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if tt.valid && err != nil {
				t.Errorf("Expected valid config, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected invalid config, got no error")
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.APIURL != "http://localhost:11434/api/generate" {
		t.Errorf("Expected default API URL, got %s", config.APIURL)
	}

	if config.Model != "codellama" {
		t.Errorf("Expected default model, got %s", config.Model)
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout, got %v", config.Timeout)
	}
}

func TestRequestStruct(t *testing.T) {
	// Teste da estrutura Request
	req := Request{
		Model:  "test-model",
		Prompt: "test prompt",
		Stream: false,
	}

	// Teste de serialização JSON
	data, err := json.Marshal(req)
	if err != nil {
		t.Errorf("Failed to marshal Request: %v", err)
	}

	var decodedReq Request
	if err := json.Unmarshal(data, &decodedReq); err != nil {
		t.Errorf("Failed to unmarshal Request: %v", err)
	}

	if decodedReq.Model != req.Model {
		t.Errorf("Expected Model %s, got %s", req.Model, decodedReq.Model)
	}

	if decodedReq.Prompt != req.Prompt {
		t.Errorf("Expected Prompt %s, got %s", req.Prompt, decodedReq.Prompt)
	}

	if decodedReq.Stream != req.Stream {
		t.Errorf("Expected Stream %v, got %v", req.Stream, decodedReq.Stream)
	}
}

func TestResponseStruct(t *testing.T) {
	// Teste da estrutura Response
	resp := Response{
		Response: "test response",
		Done:     true,
	}

	// Teste de serialização JSON
	data, err := json.Marshal(resp)
	if err != nil {
		t.Errorf("Failed to marshal Response: %v", err)
	}

	var decodedResp Response
	if err := json.Unmarshal(data, &decodedResp); err != nil {
		t.Errorf("Failed to unmarshal Response: %v", err)
	}

	if decodedResp.Response != resp.Response {
		t.Errorf("Expected Response %s, got %s", resp.Response, decodedResp.Response)
	}

	if decodedResp.Done != resp.Done {
		t.Errorf("Expected Done %v, got %v", resp.Done, decodedResp.Done)
	}
}

// Teste de integração com diferentes tipos de código
func TestExplainCodeWithDifferentLanguages(t *testing.T) {
	testCases := []struct {
		name string
		code string
	}{
		{
			name: "Go code",
			code: `func main() {
    fmt.Println("Hello World")
}`,
		},
		{
			name: "Python code",
			code: `def hello():
    print("Hello World")`,
		},
		{
			name: "JavaScript code",
			code: `function hello() {
    console.log("Hello World");
}`,
		},
		{
			name: "Empty code",
			code: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := Response{
					Response: "Análise do código " + tc.name,
					Done:     true,
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			result, err := mockExplainCode(tc.code, server.URL)
			if err != nil {
				t.Errorf("mockExplainCode() error for %s: %v", tc.name, err)
				return
			}

			if result == "" {
				t.Errorf("mockExplainCode() returned empty result for %s", tc.name)
			}
		})
	}
}
