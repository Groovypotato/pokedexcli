package main

import (
	"testing"
)


func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  This will hurt   ",
			expected: []string{"this", "will", "hurt"},
		},
		{
			input:    "Go ahead make my day   ",
			expected: []string{"go", "ahead", "make", "my", "day"},
		},
	}
	
	
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length doesn't match")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Word '%s' did not match expected word '%s' ", word, expectedWord)
			}
		}
	}
}