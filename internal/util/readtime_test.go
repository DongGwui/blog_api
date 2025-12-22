package util

import (
	"strings"
	"testing"
)

func TestCalculateReadingTime(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "Empty content",
			content:  "",
			expected: 1,
		},
		{
			name:     "Short content",
			content:  "Hello world",
			expected: 1,
		},
		{
			name:     "200 words (1 minute)",
			content:  strings.Repeat("word ", 200),
			expected: 1,
		},
		{
			name:     "400 words (2 minutes)",
			content:  strings.Repeat("word ", 400),
			expected: 2,
		},
		{
			name:     "600 words (3 minutes)",
			content:  strings.Repeat("word ", 600),
			expected: 3,
		},
		{
			name:     "Korean content",
			content:  strings.Repeat("안녕하세요 ", 200),
			expected: 2, // Korean words count differently
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateReadingTime(tt.content)
			if result != tt.expected {
				t.Errorf("CalculateReadingTime() = %d, want %d", result, tt.expected)
			}
		})
	}
}
