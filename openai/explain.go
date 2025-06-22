package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Response struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func ExplainCode(code string) (string, error) {
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

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", buf)
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
