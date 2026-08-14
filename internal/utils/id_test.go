package utils

import (
	"regexp"
	"testing"
)

func TestGenerateExpenseID(t *testing.T) {
	// Generate multiple IDs to check for uniqueness
	ids := make(map[string]bool)
	pattern := regexp.MustCompile(`^AM\d{6}$`)

	for i := 0; i < 100; i++ {
		id := GenerateExpenseID()

		// Check format
		if !pattern.MatchString(id) {
			t.Errorf("GenerateExpenseID() = %q, expected format AM######", id)
		}

		// Check for duplicates
		if ids[id] {
			t.Errorf("GenerateExpenseID() generated duplicate ID: %s", id)
		}
		ids[id] = true
	}

	// Verify we got 100 unique IDs
	if len(ids) != 100 {
		t.Errorf("Generated only %d unique IDs out of 100", len(ids))
	}
}

func TestGenerateExpenseIDFormat(t *testing.T) {
	// Test specific format requirements
	// Generate 10 IDs
	for i := 0; i < 10; i++ {
		id := GenerateExpenseID()

		// Must start with AM
		if len(id) < 2 || id[:2] != "AM" {
			t.Errorf("GenerateExpenseID() = %q, expected to start with 'AM'", id)
		}

		// Must have exactly 8 characters
		if len(id) != 8 {
			t.Errorf("GenerateExpenseID() = %q, expected length 8, got %d", id, len(id))
		}

		// Characters after AM must be digits
		for j := 2; j < len(id); j++ {
			if id[j] < '0' || id[j] > '9' {
				t.Errorf("GenerateExpenseID() = %q, expected digits after 'AM', got %c", id, id[j])
			}
		}
	}
}
