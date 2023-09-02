package lexer

import (
	"fmt"
	"reflect"
	"testing"
)

func Lex(input string) ([]Token, error) {
	tokens := []Token{}

	lexer := NewLexer(input)
	for {
		token, err := lexer.NextToken()
		if err != nil {
			return nil, err
		}

		if token.Type == TokenEOF {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

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

func TestLexer_PeekToken1(t *testing.T) {
	tests := []struct {
		input    string
		expected Token
	}{
		{
			input: "123",
			expected: Token{
				Type:    TokenNumber,
				Literal: "123",
			},
		},
		{
			input: "true",
			expected: Token{
				Type:    TokenTrue,
				Literal: "true",
			},
		},
		{
			input: `"HELLO"`,
			expected: Token{
				Type:    TokenString,
				Literal: "HELLO",
			},
		},
		{
			input: `{"Hello": "World!"}`,
			expected: Token{
				Type:    TokenLeftBrace,
				Literal: "{",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			lexer.PeekToken()
			got, err := lexer.PeekToken()
			if (err != nil) && (tt.expected != Token{}) {
				t.Errorf("Lexer.PeekToken() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Lexer.PeekToken() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestLexer_PeekToken(t *testing.T) {
	tests := []struct {
		input    string
		current  int
		expected Token
	}{
		{
			input:   "123",
			current: 1,
			expected: Token{
				Type:    TokenEOF,
				Literal: "",
			},
		},
		{
			input:   `{"Hello": "World!"}`,
			current: 1,
			expected: Token{
				Type:    TokenString,
				Literal: "Hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s, current: %d", tt.input, tt.current), func(t *testing.T) {
			lexer := NewLexer(tt.input)
			for i := 0; i < tt.current; i++ {
				_, err := lexer.NextToken()
				if err != nil {
					t.Errorf("Lexer.PeekToken() error = %v", err)
				}
			}

			lexer.PeekToken()
			got, err := lexer.PeekToken()
			if (err != nil) && (tt.expected != Token{}) {
				t.Errorf("Lexer.PeekToken() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Lexer.PeekToken() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
