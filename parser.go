package main

import (
	"regexp"
	"strconv"
	"strings"
)

type Token struct {
	Value    string
	Modifier string
	ModCount int
}

// Tokenize splits text into tokens, extracting modifier information
func Tokenize(text string) []Token {
	var tokens []Token
	parts := strings.Fields(text)

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		// Check if this part starts a modifier marker
		if strings.HasPrefix(part, "(") {
			// Collect the full modifier if it spans multiple tokens
			fullModPart := part
			for i+1 < len(parts) && !strings.HasSuffix(fullModPart, ")") {
				i++
				fullModPart += " " + parts[i]
			}

			// Check if this is a complete modifier marker
			if isModifierMarker(strings.ReplaceAll(fullModPart, " ", "")) {
				if len(tokens) > 0 {
					cleanMod := strings.ReplaceAll(fullModPart, " ", "")
					mod, count := parseModifier(cleanMod)
					tokens[len(tokens)-1].Modifier = mod
					tokens[len(tokens)-1].ModCount = count
				}
			}
		} else {
			tokens = append(tokens, Token{
				Value:    part,
				Modifier: "",
				ModCount: 0,
			})
		}
	}

	return tokens
}

// isModifierMarker checks if a token is a modifier marker like (hex), (up, 2), etc.
func isModifierMarker(token string) bool {
	modifierPattern := regexp.MustCompile(`^\((hex|bin|up|low|cap)(,\s*\d+)?\)$`)
	return modifierPattern.MatchString(token)
}

// parseModifier extracts the modifier type and optional count from a modifier marker
func parseModifier(modStr string) (string, int) {
	// Remove parentheses
	modStr = strings.TrimPrefix(modStr, "(")
	modStr = strings.TrimSuffix(modStr, ")")

	parts := strings.Split(modStr, ",")
	modType := strings.TrimSpace(parts[0])

	count := 1
	if len(parts) > 1 {
		countStr := strings.TrimSpace(parts[1])
		if c, err := strconv.Atoi(countStr); err == nil {
			count = c
		}
	}

	return modType, count
}

// TokensToString converts tokens back to a string with proper spacing
func TokensToString(tokens []Token) string {
	if len(tokens) == 0 {
		return ""
	}

	var result strings.Builder
	for i, token := range tokens {
		result.WriteString(token.Value)
		if i < len(tokens)-1 {
			result.WriteString(" ")
		}
	}

	return result.String()
}
