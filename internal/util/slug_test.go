package util

import (
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "English text",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Korean text",
			input:    "안녕하세요 세계",
			expected: "안녕하세요-세계",
		},
		{
			name:     "Mixed text",
			input:    "Go 언어 튜토리얼",
			expected: "go-언어-튜토리얼",
		},
		{
			name:     "Special characters",
			input:    "Hello! World? Test.",
			expected: "hello-world-test",
		},
		{
			name:     "Multiple spaces",
			input:    "Hello    World",
			expected: "hello-world",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  Hello World  ",
			expected: "hello-world",
		},
		{
			name:     "Numbers",
			input:    "Test 123 Example",
			expected: "test-123-example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			if result != tt.expected {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsKorean(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Korean text",
			input:    "안녕하세요",
			expected: true,
		},
		{
			name:     "English text",
			input:    "Hello World",
			expected: false,
		},
		{
			name:     "Mixed text",
			input:    "Hello 안녕",
			expected: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsKorean(tt.input)
			if result != tt.expected {
				t.Errorf("IsKorean(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
