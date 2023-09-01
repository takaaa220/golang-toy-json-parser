package lexer

import (
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input   string
		want    []Token
		wantErr bool
	}{
		{
			input: "123",
			want: []Token{
				{Type: TokenNumber, Literal: "123"},
			},
			wantErr: false,
		},
		{
			input: "null",
			want: []Token{
				{Type: TokenNull, Literal: "null"},
			},
			wantErr: false,
		},
		{
			input: "true",
			want: []Token{
				{Type: TokenBoolean, Literal: "true"},
			},
			wantErr: false,
		},
		{
			input: "false",
			want: []Token{
				{Type: TokenBoolean, Literal: "false"},
			},
			wantErr: false,
		},
		{
			input: `"Hello, World!"`,
			want: []Token{
				{Type: TokenString, Literal: "Hello, World!"},
			},
			wantErr: false,
		},
		{
			input: `{"Hello": "World!", "Foo": "Bar"}`,
			want: []Token{
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
		},
		{
			input: "[1, 2, 3]",
			want: []Token{
				{Type: TokenLeftBracket, Literal: "["},
				{Type: TokenNumber, Literal: "1"},
				{Type: TokenComma, Literal: ","},
				{Type: TokenNumber, Literal: "2"},
				{Type: TokenComma, Literal: ","},
				{Type: TokenNumber, Literal: "3"},
				{Type: TokenRightBracket, Literal: "]"},
			},
			wantErr: false,
		},
		{
			input:   "aaa",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Lex(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lex() = %v, want %v", got, tt.want)
			}
		})
	}
}
