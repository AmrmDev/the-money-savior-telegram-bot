package utils

import (
	"strings"
	"unicode"
)

func FormatTitle(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))

	words := strings.Fields(text)

	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}

	return strings.Join(words, " ")
}

func NormalizeMethod(method string) string {
	m := strings.ToLower(method)

	switch m {
	case "debito", "débito":
		return "Débito"
	case "credito", "crédito":
		return "Crédito"
	case "pix":
		return "Pix"
	case "dinheiro":
		return "Dinheiro"
	default:
		return FormatTitle(method)
	}
}