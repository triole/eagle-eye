package watcher

import (
	"testing"
)

func TestWatcher_rxfind(t *testing.T) {
	w := Watcher{}

	tests := []struct {
		name     string
		rx       string
		content  string
		expected string
	}{
		{
			name:     "simple string match",
			rx:       `hello`,
			content:  "say hello world",
			expected: "hello",
		},
		{
			name:     "regex pattern match",
			rx:       `\d+`,
			content:  "price is 123 dollars",
			expected: "123",
		},
		{
			name:     "no match",
			rx:       `xyz`,
			content:  "say hello world",
			expected: "",
		},
		{
			name:     "match with groups",
			rx:       `(\w+)@(\w+\.\w+)`,
			content:  "contact me at john@example.com",
			expected: "john@example.com",
		},
		{
			name:     "empty content",
			rx:       `hello`,
			content:  "",
			expected: "",
		},
		{
			name:     "empty regex",
			rx:       ``,
			content:  "hello world",
			expected: "",
		},
		{
			name:     "match at beginning",
			rx:       `^start`,
			content:  "start of the line",
			expected: "start",
		},
		{
			name:     "match at end",
			rx:       `end$`,
			content:  "at the end",
			expected: "end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := w.rxfind(tt.rx, tt.content)
			if result != tt.expected {
				t.Errorf("rxfind(%q, %q) = %q, want %q", tt.rx, tt.content, result, tt.expected)
			}
		})
	}
}
