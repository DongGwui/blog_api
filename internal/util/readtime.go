package util

import (
	"strings"
	"unicode"
)

const (
	// Average reading speeds
	wordsPerMinuteEnglish = 200
	charsPerMinuteKorean  = 500
)

// CalculateReadingTime estimates the reading time in minutes
// Supports both English and Korean content
func CalculateReadingTime(content string) int {
	if content == "" {
		return 1
	}

	// Strip HTML tags (simple approach)
	content = stripHTMLTags(content)

	koreanChars := 0
	englishWords := 0
	inWord := false

	for _, r := range content {
		if unicode.Is(unicode.Hangul, r) {
			koreanChars++
			inWord = false
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if !inWord {
				englishWords++
				inWord = true
			}
		} else {
			inWord = false
		}
	}

	// Calculate time for each language
	koreanMinutes := float64(koreanChars) / float64(charsPerMinuteKorean)
	englishMinutes := float64(englishWords) / float64(wordsPerMinuteEnglish)

	totalMinutes := koreanMinutes + englishMinutes

	// Minimum 1 minute
	if totalMinutes < 1 {
		return 1
	}

	// Round up
	return int(totalMinutes + 0.5)
}

func stripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false

	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}

	return result.String()
}
