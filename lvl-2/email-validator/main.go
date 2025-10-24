package main

import (
	"fmt"
	"strings"
)

func main() {
	// Testes
	testCases := []string{
		"user@example.com",
		"invalid.email",
		" user@example.com",
		"user@example.com ",
		"user@domain",
		"",
		"user@@example.com",
		"@example.com",
		"user@",
	}

	for _, email := range testCases {
		valid, msg := emailValidator(email)
		status := "❌"
		if valid {
			status = "✅"
		}
		fmt.Printf("%s %-25s -> %s\n", status, fmt.Sprintf("\"%s\"", email), msg)
	}
}

func emailValidator(e string) (bool, string) {
	// Verifica se está vazio
	if e == "" {
		return false, "email não pode estar vazio"
	}

	// Verifica espaços no início ou fim
	if strings.HasPrefix(e, " ") || strings.HasSuffix(e, " ") {
		return false, "email não pode começar ou terminar com espaços"
	}

	// Conta quantos @ existem
	atCount := strings.Count(e, "@")
	if atCount == 0 {
		return false, "email deve conter @"
	}
	if atCount > 1 {
		return false, "email deve conter exatamente um @"
	}

	// Separa em partes antes e depois do @
	parts := strings.Split(e, "@")
	localPart := parts[0]
	domainPart := parts[1]

	// Verifica se há conteúdo antes e depois do @
	if localPart == "" {
		return false, "email deve ter conteúdo antes do @"
	}
	if domainPart == "" {
		return false, "email deve ter conteúdo depois do @"
	}

	// Verifica se tem pelo menos um . depois do @
	if !strings.Contains(domainPart, ".") {
		return false, "domínio deve conter pelo menos um ponto"
	}

	return true, "email válido"
}
