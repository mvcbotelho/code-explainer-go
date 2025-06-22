package openai

import "strings"

// DetectLanguage tenta identificar a linguagem do código com base em padrões simples
func DetectLanguage(code string) string {
	code = strings.ToLower(code)

	switch {
	case strings.Contains(code, "func ") || strings.Contains(code, "package ") || strings.Contains(code, "import ("):
		return "Go"
	case strings.Contains(code, "def ") || strings.Contains(code, "print(") || strings.Contains(code, "import "):
		return "Python"
	case strings.Contains(code, "console.log") || strings.Contains(code, "function") || strings.Contains(code, "var "):
		return "JavaScript"
	case strings.Contains(code, "#include") || strings.Contains(code, "printf(") || strings.Contains(code, "int main"):
		return "C"
	case strings.Contains(code, "system.out.println") || strings.Contains(code, "public class") || strings.Contains(code, "public static"):
		return "Java"
	case strings.Contains(code, "echo ") || strings.Contains(code, "$_post") || strings.Contains(code, "<?php"):
		return "PHP"
	case strings.Contains(code, "fn ") || strings.Contains(code, "println!") || (strings.Contains(code, "let ") && strings.Contains(code, "mut")):
		return "Rust"
	case strings.Contains(code, "class ") && strings.Contains(code, "public ") && strings.Contains(code, "{") && !strings.Contains(code, "system.out"):
		return "C#"
	default:
		return "linguagem desconhecida"
	}
}
