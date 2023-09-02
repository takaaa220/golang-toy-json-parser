package lexer

import (
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
		err      string
	}{
		{
			input: "123",
			expected: []Token{
				{Type: TokenNumber, Literal: "123"},
			},
			err: "",
		},
		{
			input: "null",
			expected: []Token{
				{Type: TokenNull, Literal: "null"},
			},
			err: "",
		},
		{
			input: "true",
			expected: []Token{
				{Type: TokenTrue, Literal: "true"},
			},
			err: "",
		},
		{
			input: "false",
			expected: []Token{
				{Type: TokenFalse, Literal: "false"},
			},
			err: "",
		},
		{
			input: `"Hello, World!"`,
			expected: []Token{
				{Type: TokenString, Literal: "Hello, World!"},
			},
			err: "",
		},
		{
			input: `{"Hello": "World!", "Foo": "Bar"}`,
			expected: []Token{
				{Type: TokenLeftBrace, Literal: "{"},
				{Type: TokenString, Literal: "Hello"},
				{Type: TokenColon, Literal: ":"},
				{Type: TokenString, Literal: "World!"},
				{Type: TokenComma, Literal: ","},
				{Type: TokenString, Literal: "Foo"},
				{Type: TokenColon, Literal: ":"},
				{Type: TokenString, Literal: "Bar"},
				{Type: TokenRightBrace, Literal: "}"},
			},
			err: "",
		},
		{
			input: "[1, 2, 3]",
			expected: []Token{
				{Type: TokenLeftBracket, Literal: "["},
				{Type: TokenNumber, Literal: "1"},
				{Type: TokenComma, Literal: ","},
				{Type: TokenNumber, Literal: "2"},
				{Type: TokenComma, Literal: ","},
				{Type: TokenNumber, Literal: "3"},
				{Type: TokenRightBracket, Literal: "]"},
			},
			err: "",
		},
		{
			input: "0.1",
			expected: []Token{
				{Type: TokenNumber, Literal: "0.1"},
			},
		},
		{
			input:    "0.1.1",
			expected: nil,
			err:      "unexpected number 0.1.1, From: 0, To: 4",
		},
		{
			input:    `{"Hello": 1.133.1}`,
			expected: nil,
			err:      "unexpected number 1.133.1, From: 10, To: 16",
		},
		{
			input:    "aaa",
			expected: nil,
			err:      "unexpected property 'aaa', From: 0, To: 2",
		},
		{
			input:    `"Hello World!`,
			expected: nil,
			err:      "syntax error, From: 0, To: 12",
		},
		{
			input:    `{ "Hello": "World! }`,
			expected: nil,
			err:      "syntax error, From: 11, To: 19",
		},
		{
			input:    `{ "Hello": 100, "World": false, "!!!": True }`,
			expected: nil,
			err:      "unexpected character 'T', From: 39, To: 39",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Lex(tt.input)
			if err != nil {
				if tt.err == "" || err.Error() != tt.err {
					t.Errorf("Lex() error = %v", err)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Lex() = %v, want %v", got, tt.expected)
			}
		})
	}
}
