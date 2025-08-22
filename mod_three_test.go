package fsm

import (
	"testing"
)

// Unit tests: basic cases
func TestModThreeFSM_Basic(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"0", 0}, // smallest input
		{"1", 1},
		{"10", 2},   // decimal 2
		{"11", 0},   // decimal 3
		{"1101", 1}, // decimal 13
		{"1110", 2}, // decimal 14
		{"1111", 0}, // decimal 15
		{"", 0},     // empty input -> treated as 0
	}

	for _, tt := range tests {
		got := modThreeFSM(tt.input)
		if got != tt.expected {
			t.Errorf("modThreeFSM(%q) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}

// Unit tests: edge cases
func TestModThreeFSM_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"000000", 0}, // all zeros
		{"111111", 0}, // 63 mod 3 == 0
		{"101010", 0}, // 42 mod 3 == 0
		{"1011", 2},   // decimal 11
	}

	for _, tt := range tests {
		got := modThreeFSM(tt.input)
		if got != tt.expected {
			t.Errorf("modThreeFSM(%q) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}
