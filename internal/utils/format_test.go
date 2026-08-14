package utils

import (
	"testing"
)

func TestFormatTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single word lowercase",
			input:    "food",
			expected: "Food",
		},
		{
			name:     "multiple words lowercase",
			input:    "super market",
			expected: "Super Market",
		},
		{
			name:     "mixed case",
			input:    "UbEr RiDe",
			expected: "Uber Ride",
		},
		{
			name:     "with extra spaces",
			input:    "   food   market   ",
			expected: "Food Market",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "A",
		},
		{
			name:     "numbers in string",
			input:    "gas 123 station",
			expected: "Gas 123 Station",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTitle(tt.input)
			if result != tt.expected {
				t.Errorf("FormatTitle(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeMethod(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "debito lowercase",
			input:    "debito",
			expected: "Débito",
		},
		{
			name:     "débito with accent",
			input:    "débito",
			expected: "Débito",
		},
		{
			name:     "credito lowercase",
			input:    "credito",
			expected: "Crédito",
		},
		{
			name:     "crédito with accent",
			input:    "crédito",
			expected: "Crédito",
		},
		{
			name:     "pix",
			input:    "pix",
			expected: "Pix",
		},
		{
			name:     "PIX uppercase",
			input:    "PIX",
			expected: "Pix",
		},
		{
			name:     "dinheiro",
			input:    "dinheiro",
			expected: "Dinheiro",
		},
		{
			name:     "DINHEIRO uppercase",
			input:    "DINHEIRO",
			expected: "Dinheiro",
		},
		{
			name:     "custom method",
			input:    "transferencia",
			expected: "Transferencia",
		},
		{
			name:     "custom method with spaces",
			input:    "wire transfer",
			expected: "Wire Transfer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeMethod(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeMethod(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
