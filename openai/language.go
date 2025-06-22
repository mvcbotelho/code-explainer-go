package openai

import (
	"regexp"
	"strings"
)

// LanguagePattern define um padrão de linguagem com expressões regulares
type LanguagePattern struct {
	Language string
	Patterns []*regexp.Regexp
	Priority int // Prioridade mais alta = mais específico
}

// languagePatterns define os padrões de detecção para cada linguagem
// Ordenados por prioridade (mais específicos primeiro)
var languagePatterns = []LanguagePattern{
	{
		Language: "Go",
		Priority: 100,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`\bpackage\s+\w+`),
			regexp.MustCompile(`import\s*\(`),
			regexp.MustCompile(`\bfunc\s+\w+\s*\(`),
			regexp.MustCompile(`\bdefer\b`),
			regexp.MustCompile(`\bgo\s+\w+`),
			regexp.MustCompile(`\bchan\b`),
		},
	},
	{
		Language: "Rust",
		Priority: 90,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`\bfn\s+\w+\s*\(`),
			regexp.MustCompile(`\blet\s+mut\b`),
			regexp.MustCompile(`println!\s*\(`),
			regexp.MustCompile(`\buse\s+`),
			regexp.MustCompile(`\bstruct\s+\w+`),
			regexp.MustCompile(`\benum\s+\w+`),
		},
	},
	{
		Language: "C#",
		Priority: 80,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`\bpublic\s+class\s+\w+`),
			regexp.MustCompile(`\bnamespace\s+\w+`),
			regexp.MustCompile(`\busing\s+System`),
			regexp.MustCompile(`\bvar\s+\w+\s*=`),
			regexp.MustCompile(`\bConsole\.WriteLine`),
		},
	},
	{
		Language: "Java",
		Priority: 70,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`\bpublic\s+class\s+\w+`),
			regexp.MustCompile(`System\.out\.println`),
			regexp.MustCompile(`\bimport\s+java\.`),
			regexp.MustCompile(`\bpublic\s+static\s+void\s+main`),
			regexp.MustCompile(`\bString\s+\w+`),
		},
	},
	{
		Language: "Python",
		Priority: 60,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`\bdef\s+\w+\s*\(`),
			regexp.MustCompile(`\bimport\s+\w+`),
			regexp.MustCompile(`\bfrom\s+\w+\s+import`),
			regexp.MustCompile(`print\s*\(`),
			regexp.MustCompile(`\bif\s+__name__\s*==\s*['"]__main__['"]`),
			regexp.MustCompile(`\bclass\s+\w+`),
		},
	},
	{
		Language: "JavaScript",
		Priority: 50,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`console\.log\s*\(`),
			regexp.MustCompile(`\bfunction\s+\w+\s*\(`),
			regexp.MustCompile(`\bvar\s+\w+`),
			regexp.MustCompile(`\blet\s+\w+`),
			regexp.MustCompile(`\bconst\s+\w+`),
			regexp.MustCompile(`\bexport\s+`),
			regexp.MustCompile(`\bimport\s+\w+\s+from`),
		},
	},
	{
		Language: "C",
		Priority: 40,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`#include\s*<`),
			regexp.MustCompile(`\bint\s+main\s*\(`),
			regexp.MustCompile(`printf\s*\(`),
			regexp.MustCompile(`\bstruct\s+\w+`),
			regexp.MustCompile(`\b#define\b`),
		},
	},
	{
		Language: "PHP",
		Priority: 30,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`<\?php`),
			regexp.MustCompile(`\becho\s+`),
			regexp.MustCompile(`\$_POST\b`),
			regexp.MustCompile(`\$_GET\b`),
			regexp.MustCompile(`\bfunction\s+\w+\s*\(`),
			regexp.MustCompile(`\bclass\s+\w+`),
		},
	},
}

// DetectLanguage tenta identificar a linguagem do código com base em padrões de expressões regulares
func DetectLanguage(code string) string {
	if code == "" {
		return "linguagem desconhecida"
	}

	// Normaliza o código para análise
	code = strings.ToLower(code)

	// Remove comentários para evitar interferência
	code = removeComments(code)

	// Procura por padrões em ordem de prioridade
	for _, lang := range languagePatterns {
		for _, pattern := range lang.Patterns {
			if pattern.MatchString(code) {
				return lang.Language
			}
		}
	}

	return "linguagem desconhecida"
}

// removeComments remove comentários comuns para melhorar a detecção
func removeComments(code string) string {
	// Remove comentários de linha única (//, #, //)
	lines := strings.Split(code, "\n")
	var cleanLines []string

	for _, line := range lines {
		// Remove comentários de linha única
		if idx := strings.Index(line, "//"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		if idx := strings.Index(line, "//"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		if strings.TrimSpace(line) != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// GetSupportedLanguages retorna a lista de linguagens suportadas
func GetSupportedLanguages() []string {
	languages := make([]string, len(languagePatterns))
	for i, lang := range languagePatterns {
		languages[i] = lang.Language
	}
	return languages
}

// AddLanguagePattern permite adicionar novos padrões de linguagem dinamicamente
func AddLanguagePattern(language string, patterns []string, priority int) error {
	var regexPatterns []*regexp.Regexp

	for _, pattern := range patterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		regexPatterns = append(regexPatterns, regex)
	}

	newPattern := LanguagePattern{
		Language: language,
		Patterns: regexPatterns,
		Priority: priority,
	}

	// Insere na posição correta baseada na prioridade
	for i, existing := range languagePatterns {
		if priority > existing.Priority {
			languagePatterns = append(languagePatterns[:i], append([]LanguagePattern{newPattern}, languagePatterns[i:]...)...)
			return nil
		}
	}

	// Se não encontrou posição, adiciona no final
	languagePatterns = append(languagePatterns, newPattern)
	return nil
}
