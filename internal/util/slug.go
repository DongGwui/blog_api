package util

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9가-힣]+`)
	multiHyphenRegex     = regexp.MustCompile(`-+`)
)

// GenerateSlug creates a URL-friendly slug from the given string
// Supports both English and Korean characters
func GenerateSlug(s string) string {
	// Normalize unicode
	s = norm.NFC.String(s)

	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove non-alphanumeric characters (keep Korean)
	s = nonAlphanumericRegex.ReplaceAllString(s, "-")

	// Replace multiple hyphens with single hyphen
	s = multiHyphenRegex.ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	return s
}

// IsKorean checks if the string contains Korean characters
func IsKorean(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Hangul, r) {
			return true
		}
	}
	return false
}
