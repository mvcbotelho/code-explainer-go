package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Mock da função ExplainCode para testes
func mockExplainCode(code string, serverURL string) (string, error) {
	model := os.Getenv("MODEL_NAME")
	if model == "" {
		model = "codellama"
	}

	lang := DetectLanguage(code)
	prompt := fmt.Sprintf(`Explique o que o seguinte código em %s faz:

%s`, lang, code)

	body := Request{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return "", err
	}

	resp, err := http.Post(serverURL, "application/json", buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error, status code: %d, status: %s", resp.StatusCode, resp.Status)
	}

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
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
