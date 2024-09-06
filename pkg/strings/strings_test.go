package strings

import "testing"

func TestIsDigit(t *testing.T) {
	tests := []struct {
		ch       byte
		expected bool
	}{
		{'5', true},
		{'a', false},
	}

	for _, tt := range tests {
		if got := IsDigit(tt.ch); got != tt.expected {
			t.Errorf("IsDigit() = %v, expected %v", got, tt.expected)
		}
	}
}

func TestIsLetter(t *testing.T) {
	tests := []struct {
		ch       byte
		expected bool
	}{
		{'a', true},
		{'5', false},
		{'_', true},
		{'$', false},
	}

	for _, tt := range tests {
		if got := IsLetter(tt.ch); got != tt.expected {
			t.Errorf("IsLetter() = %v, expected %v", got, tt.expected)
		}
	}
}

func TestIsWhiteSpace(t *testing.T) {
	tests := []struct {
		ch       byte
		expected bool
	}{
		{'\n', true},
		{' ', true},
		{'\t', true},
		{'\r', true},
		{'5', false},
		{'a', false},
	}

	for _, tt := range tests {
		if got := IsWhiteSpace(tt.ch); got != tt.expected {
			t.Errorf("IsWhiteSpace() = %v, expected %v", got, tt.expected)
		}
	}
}
