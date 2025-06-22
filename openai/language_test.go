package openai

import (
	"testing"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		// Testes para Go
		{
			name: "Go com func",
			code: `func main() {
    fmt.Println("Hello World")
}`,
			expected: "Go",
		},
		{
			name: "Go com package",
			code: `package main

import "fmt"`,
			expected: "Go",
		},
		{
			name: "Go com import",
			code: `import (
    "fmt"
    "os"
)`,
			expected: "Go",
		},
		{
			name: "Go case insensitive",
			code: `FUNC main() {
    PACKAGE main
}`,
			expected: "Go",
		},

		// Testes para Python
		{
			name: "Python com def",
			code: `def hello_world():
    print("Hello World")`,
			expected: "Python",
		},
		{
			name:     "Python com print",
			code:     `print("Hello World")`,
			expected: "Python",
		},
		{
			name: "Python com import",
			code: `import os
import sys`,
			expected: "Python",
		},

		// Testes para JavaScript
		{
			name:     "JavaScript com console.log",
			code:     `console.log("Hello World");`,
			expected: "JavaScript",
		},
		{
			name: "JavaScript com function",
			code: `function hello() {
    return "Hello World";
}`,
			expected: "JavaScript",
		},
		{
			name:     "JavaScript com var",
			code:     `var x = 10;`,
			expected: "JavaScript",
		},
		{
			name:     "JavaScript com let",
			code:     `let y = 20;`,
			expected: "JavaScript",
		},

		// Testes para C
		{
			name: "C com #include",
			code: `#include <stdio.h>
#include <stdlib.h>`,
			expected: "C",
		},
		{
			name:     "C com printf",
			code:     `printf("Hello World\n");`,
			expected: "C",
		},
		{
			name: "C com int main",
			code: `int main() {
    return 0;
}`,
			expected: "C",
		},

		// Testes para Java
		{
			name:     "Java com System.out.println",
			code:     `System.out.println("Hello World");`,
			expected: "Java",
		},
		{
			name: "Java com public class",
			code: `public class Hello {
    public static void main(String[] args) {
    }
}`,
			expected: "Java",
		},
		{
			name: "Java com public static",
			code: `public static void main(String[] args) {
    System.out.println("Hello");
}`,
			expected: "Java",
		},

		// Testes para PHP
		{
			name:     "PHP com echo",
			code:     `echo "Hello World";`,
			expected: "PHP",
		},
		{
			name:     "PHP com $_POST",
			code:     `$name = $_post['name'];`,
			expected: "PHP",
		},
		{
			name: "PHP com <?php",
			code: `<?php
echo "Hello World";
?>`,
			expected: "PHP",
		},

		// Testes para Rust
		{
			name: "Rust com fn",
			code: `fn main() {
    println!("Hello World");
}`,
			expected: "Rust",
		},
		{
			name:     "Rust com let",
			code:     `let x = 10;`,
			expected: "Rust",
		},
		{
			name:     "Rust com println!",
			code:     `println!("Hello World");`,
			expected: "Rust",
		},

		// Testes para C#
		{
			name: "C# com class e public",
			code: `public class Program {
    public static void Main() {
    }
}`,
			expected: "C#",
		},

		// Testes para linguagem desconhecida
		{
			name:     "Código vazio",
			code:     "",
			expected: "linguagem desconhecida",
		},
		{
			name:     "Texto simples",
			code:     "Hello World",
			expected: "linguagem desconhecida",
		},
		{
			name: "Código sem padrões conhecidos",
			code: `hello world
this is some text
without known patterns`,
			expected: "linguagem desconhecida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectLanguage(tt.code)
			if result != tt.expected {
				t.Errorf("DetectLanguage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDetectLanguageCaseInsensitive(t *testing.T) {
	// Teste específico para verificar se a detecção é case insensitive
	code := `FUNC MAIN() {
    PACKAGE MAIN
    IMPORT (
        "FMT"
    )
}`

	result := DetectLanguage(code)
	if result != "Go" {
		t.Errorf("DetectLanguage() deveria detectar Go independente do case, got %v", result)
	}
}

func TestDetectLanguagePriority(t *testing.T) {
	// Teste para verificar a prioridade das detecções
	// Código que contém padrões de múltiplas linguagens
	code := `func main() {
    console.log("Hello");
    print("World");
}`

	result := DetectLanguage(code)
	// Deve detectar Go primeiro (primeiro case no switch)
	if result != "Go" {
		t.Errorf("DetectLanguage() deveria priorizar Go, got %v", result)
	}
}
