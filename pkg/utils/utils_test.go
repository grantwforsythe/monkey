package utils

import "testing"

func TestIsDigit(t *testing.T) {
	tests := []struct {
		name string
		ch   byte
		want bool
	}{
		{"IsDigit", '5', true},
		{"IsDigit", 'a', false},
	}
	for _, tt := range tests {
		if got := IsDigit(tt.ch); got != tt.want {
			t.Errorf("IsDigit() = %v, want %v", got, tt.want)
		}
	}
}

func TestIsLetter(t *testing.T) {
	tests := []struct {
		name string
		ch   byte
		want bool
	}{
		{"IsLetter", 'a', true},
		{"IsLetter", '5', false},
		{"IsLetter", '_', true},
		{"IsLetter", '$', false},
	}

	for _, tt := range tests {
		if got := IsLetter(tt.ch); got != tt.want {
			t.Errorf("IsLetter() = %v, want %v", got, tt.want)
		}
	}
}
