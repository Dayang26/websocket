package websocket

import (
	"testing"
)

func TestSkipSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{" token", "token"},
		{"\ttoken", "token"},
		{"   token", "token"},
		{"\ttoken", "token"},
		{"  \ttoken", "token"},
		{"\t\t token", "token"},
		{"token", "token"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := skipSpace(tt.input)
			if actual != tt.expected {
				t.Errorf("SkipSpace(%q) = %q; want %q", tt.input, actual, tt.expected)
			}
		})
	}
}
