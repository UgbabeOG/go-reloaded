package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Process is the main orchestration function that applies all transformations
func Process(text string) string {
	tokens := Tokenize(text)

	// Apply transformations in sequence
	tokens = applyModifiers(tokens)
	tokens = handleArticles(tokens)
	tokens = handleQuotes(tokens)
	tokens = handlePunctuation(tokens)

	return TokensToString(tokens)
}

// applyModifiers applies the modifiers (hex, bin, up, low, cap) to tokens
func applyModifiers(tokens []Token) []Token {
	for i := 0; i < len(tokens); i++ {
		token := &tokens[i]

		if token.Modifier == "" {
			continue
		}

		switch token.Modifier {
		case "hex":
			hexVal := extractNumber(token.Value)
			if hexVal != "" {
				if dec, err := strconv.ParseInt(hexVal, 16, 64); err == nil {
					token.Value = replaceNumber(token.Value, strconv.FormatInt(dec, 10))
				}
			}
		case "bin":
			binVal := extractNumber(token.Value)
			if binVal != "" {
				if dec, err := strconv.ParseInt(binVal, 2, 64); err == nil {
					token.Value = replaceNumber(token.Value, strconv.FormatInt(dec, 10))
				}
			}
		case "up":
			applyCase(tokens, i, token.ModCount, strings.ToUpper)
		case "low":
			applyCase(tokens, i, token.ModCount, strings.ToLower)
		case "cap":
			applyCase(tokens, i, token.ModCount, capitalize)
		}
	}

	return tokens
}

// extractNumber extracts the numeric part from a token (could have punctuation)
func extractNumber(token string) string {
	numPattern := regexp.MustCompile(`[0-9A-Fa-f]+`)
	match := numPattern.FindString(token)
	return match
}

// replaceNumber replaces the numeric part in a token with a new value
func replaceNumber(token, newNum string) string {
	numPattern := regexp.MustCompile(`[0-9A-Fa-f]+`)
	return numPattern.ReplaceAllString(token, newNum)
}

// capitalize returns the capitalized version of a string
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

// applyCase applies a case transformation function to the N words before and including modIndex
func applyCase(tokens []Token, modIndex int, count int, caseFn func(string) string) {
	start := modIndex - count + 1
	if start < 0 {
		start = 0
	}

	for i := start; i <= modIndex; i++ {
		// Apply case transformation to the word part only, preserving non-word characters
		tokens[i].Value = transformWordPart(tokens[i].Value, caseFn)
	}
}

// transformWordPart applies a case transformation to just the word characters in a token
func transformWordPart(token string, caseFn func(string) string) string {
	// Extract word and punctuation
	wordPattern := regexp.MustCompile(`^([a-zA-Z']+)(.*)$`)
	matches := wordPattern.FindStringSubmatch(token)

	if len(matches) == 3 {
		word := matches[1]
		suffix := matches[2]
		return caseFn(word) + suffix
	}

	return caseFn(token)
}

// handlePunctuation formats punctuation according to rules
func handlePunctuation(tokens []Token) []Token {
	var result []Token

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// Check if this token is just punctuation
		if isJustPunctuation(token.Value) {
			// Combine with previous token if possible
			if len(result) > 0 {
				// Attach punctuation to the previous token
				result[len(result)-1].Value += token.Value
			}
		} else {
			// Extract leading punctuation
			leadingPunct := ""
			j := 0
			for j < len(token.Value) && isPunctuation(rune(token.Value[j])) {
				leadingPunct += string(token.Value[j])
				j++
			}

			// If there's leading punctuation, attach it to the previous token
			if leadingPunct != "" && len(result) > 0 {
				result[len(result)-1].Value += leadingPunct
				token.Value = token.Value[j:]
			}

			if token.Value != "" {
				result = append(result, token)
			}
		}
	}

	return result
}

// isJustPunctuation checks if a token contains only punctuation
func isJustPunctuation(token string) bool {
	for _, r := range token {
		if !isPunctuation(r) {
			return false
		}
	}
	return len(token) > 0
}

// isPunctuation checks if a rune is a punctuation mark we care about
func isPunctuation(r rune) bool {
	punctSet := ".,!?:;'\""
	return strings.ContainsRune(punctSet, r)
}

// isPunctGrouped checks if a punctuation sequence is a grouped punctuation
func isPunctGrouped(punct string) bool {
	grouped := []string{"...", "!?", "?!", "...", "....", "!!!"}
	for _, g := range grouped {
		if punct == g {
			return true
		}
	}
	// Check if all characters are the same and more than 1
	if len(punct) > 1 {
		first := punct[0]
		allSame := true
		for i := 1; i < len(punct); i++ {
			if punct[i] != first {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}
	return false
}

// handleQuotes removes spaces within quote pairs
func handleQuotes(tokens []Token) []Token {
	// Find quote pairs
	quoteIndices := []int{}
	for i, token := range tokens {
		if strings.Contains(token.Value, "'") {
			quoteIndices = append(quoteIndices, i)
		}
	}

	// Process quote pairs
	for i := 0; i+1 < len(quoteIndices); i += 2 {
		openIdx := quoteIndices[i]
		closeIdx := quoteIndices[i+1]

		if openIdx < closeIdx {
			// Remove leading space from the token after opening quote
			if openIdx+1 <= closeIdx {
				tokens[openIdx].Value = strings.TrimSuffix(tokens[openIdx].Value, " ")
			}

			// Remove trailing space from the token before closing quote
			if closeIdx > openIdx {
				tokens[closeIdx].Value = strings.TrimPrefix(tokens[closeIdx].Value, " ")
			}
		}
	}

	return tokens
}

// handleArticles converts "a" to "an" before vowels or 'h'
func handleArticles(tokens []Token) []Token {
	for i := 0; i < len(tokens)-1; i++ {
		token := tokens[i]

		if strings.ToLower(token.Value) == "a" {
			nextToken := tokens[i+1]
			nextWord := strings.ToLower(nextToken.Value)

			// Skip punctuation and quotes to find the actual word
			nextWord = strings.TrimLeft(nextWord, "\"'")
			cleanWord := extractWordFromToken(nextWord)

			if len(cleanWord) > 0 {
				firstChar := rune(cleanWord[0])
				if isVowel(firstChar) || firstChar == 'h' {
					tokens[i].Value = "an"
				}
			}
		}
	}

	return tokens
}

// isVowel checks if a rune is a vowel
func isVowel(r rune) bool {
	return strings.ContainsRune("aeiouAEIOU", r)
}

// extractWordFromToken extracts just the word characters from a token
func extractWordFromToken(token string) string {
	// Remove non-word characters from the start
	for i, r := range token {
		if unicode.IsLetter(r) {
			return token[i:]
		}
	}
	return ""
}
