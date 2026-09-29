package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: " hello world",
			expected: []string{"hello","world"},
		},
		{
			input: " hElLO wOrlD",
			expected: []string{"hello","world"},
		},
		{
			input: "    ",
			expected: []string{},
		},
		{
			input: " Hello World",
			expected: []string{"hello","world"},
		},
		{
			input: " Hello World    ",
			expected: []string{"hello","world"},
		},
		{
			input: " Hello\nWorld\tagain",
			expected: []string{"hello","world","again"},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: "PIKACHU",
			expected: []string{"pikachu"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Input: %v. Actual length does not match the expected length", c.input)
			continue
		}
		for i, _ := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Word: %v does not match expected %v", word, expectedWord)
			}
		}
	}
}
