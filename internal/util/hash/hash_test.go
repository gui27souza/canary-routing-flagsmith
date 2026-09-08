package hash_test

import (
	"goflagsmith/internal/util/hash"
	"testing"
)

func TestNormalizedHash(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{name: "User 5", input: "user_id_5", expected: 90},
		{name: "User 6", input: "user_id_6", expected: 71},
		{name: "User 7", input: "user_id_7", expected: 52},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := hash.NormalizedHash(tc.input)
			if output != tc.expected {
				t.Errorf("NormalizedHash() = %d; want %d", output, tc.expected)
			}
		})
	}
}

func TestNormalizedHash_Boundaries(t *testing.T) {
	inputs := []string{"", "a", "user_id_9999999999999", "UUID-very-long-string-with-special-chars-!"}

	for _, input := range inputs {
		output := hash.NormalizedHash(input)
		if output < 0 || output >= 100 {
			t.Errorf("NormalizedHash(%q) yielded %d, which is out of [0, 99] bounds", input, output)
		}
	}
}
