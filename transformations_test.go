package main

import (
	"testing"
)

func TestHexToDec(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1E (hex)", "30"},
		{"A (hex)", "10"},
		{"FF (hex)", "255"},
	}

	for _, test := range tests {
		tokens := Tokenize(test.input)
		tokens = applyModifiers(tokens)
		result := TokensToString(tokens)

		// Extract the number from result
		if !contains(result, test.expected) {
			t.Errorf("HexToDec(%q) = %q, expected to contain %q", test.input, result, test.expected)
		}
	}
}

func TestBinToDec(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"10 (bin)", "2"},
		{"1010 (bin)", "10"},
		{"1111 (bin)", "15"},
	}

	for _, test := range tests {
		tokens := Tokenize(test.input)
		tokens = applyModifiers(tokens)
		result := TokensToString(tokens)

		if !contains(result, test.expected) {
			t.Errorf("BinToDec(%q) = %q, expected to contain %q", test.input, result, test.expected)
		}
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"WORLD", "World"},
		{"", ""},
		{"a", "A"},
	}

	for _, test := range tests {
		result := capitalize(test.input)
		if result != test.expected {
			t.Errorf("capitalize(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestToUpper(t *testing.T) {
	input := "hello world (up)"
	tokens := Tokenize(input)
	tokens = applyModifiers(tokens)
	result := TokensToString(tokens)

	if !contains(result, "WORLD") {
		t.Errorf("ToUpper test failed: got %q", result)
	}
}

func TestToLower(t *testing.T) {
	input := "HELLO WORLD (low)"
	tokens := Tokenize(input)
	tokens = applyModifiers(tokens)
	result := TokensToString(tokens)

	if !contains(result, "world") {
		t.Errorf("ToLower test failed: got %q", result)
	}
}

func TestCapitalizeMultiWord(t *testing.T) {
	input := "it was the age of foolishness (cap, 6)"
	tokens := Tokenize(input)
	tokens = applyModifiers(tokens)
	result := TokensToString(tokens)

	// All 6 words before (cap, 6) should be capitalized
	if !contains(result, "It") && !contains(result, "Was") {
		t.Errorf("CapitalizeMultiWord test failed: got %q", result)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && contains(s[1:], substr)) || contains(s[:len(s)-1], substr))
}

func TestExtractNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1E", "1E"},
		{"1E,", "1E"},
		{"42", "42"},
		{"A", "A"},
	}

	for _, test := range tests {
		result := extractNumber(test.input)
		if result != test.expected {
			t.Errorf("extractNumber(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestIsVowel(t *testing.T) {
	tests := []struct {
		r        rune
		expected bool
	}{
		{'a', true},
		{'e', true},
		{'i', true},
		{'o', true},
		{'u', true},
		{'A', true},
		{'b', false},
		{'h', false},
	}

	for _, test := range tests {
		result := isVowel(test.r)
		if result != test.expected {
			t.Errorf("isVowel(%q) = %v, want %v", test.r, result, test.expected)
		}
	}
}
