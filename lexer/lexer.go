package lexer

import (
	"bufio"
	"strings"
	"unicode"
)

type LexerError struct {
	Message string
	From    int
	To      int
}

func (e *LexerError) Error() string {
	return e.Message
}

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenString
	TokenBoolean
	TokenNull
	TokenColon        // :
	TokenComma        // ,
	TokenLeftBrace    // {
	TokenRightBrace   // }
	TokenLeftBracket  // [
	TokenRightBracket // ]
)

type Token struct {
	Type    TokenType
	Literal string
}

func isLowerLetter(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}

func Lex(input string) ([]Token, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	tokens := []Token{}

	for {
		// i seems not to be index of input
		// TODO: fix this
		r, i, err := reader.ReadRune()
		if err != nil {
			break
		}

		switch {
		case r == '{':
			tokens = append(tokens, Token{Type: TokenLeftBrace, Literal: string(r)})
			// 次の } まで読み出す
		case r == '}':
			tokens = append(tokens, Token{Type: TokenRightBrace, Literal: string(r)})
		case r == '[':
			tokens = append(tokens, Token{Type: TokenLeftBracket, Literal: string(r)})
		case r == ']':
			tokens = append(tokens, Token{Type: TokenRightBracket, Literal: string(r)})
		case r == ':':
			tokens = append(tokens, Token{Type: TokenColon, Literal: string(r)})
		case r == ',':
			tokens = append(tokens, Token{Type: TokenComma, Literal: string(r)})
		case r == '"':
			s := ""
			ok := false

			for {
				r, i, err := reader.ReadRune()
				if err != nil {
					return nil, &LexerError{Message: "invalid string", From: i, To: i}
				}

				if r == '"' {
					ok = true
					break
				}

				s += string(r)
			}

			if !ok {
				return nil, &LexerError{Message: "syntax error", From: i, To: len(input) - 1}
			}

			tokens = append(tokens, Token{Type: TokenString, Literal: s})
		case unicode.IsDigit(r):
			s := string(r)

			for {
				ss, _, err := reader.ReadRune()
				if err != nil {
					break
				}

				if !unicode.IsDigit(ss) {
					reader.UnreadRune()

					break
				}

				s += string(ss)
			}
			tokens = append(tokens, Token{Type: TokenNumber, Literal: s})
		case isLowerLetter(r):
			s := string(r)

			for {
				ss, _, err := reader.ReadRune()
				if err != nil {
					break
				}

				if !isLowerLetter(ss) {
					reader.UnreadRune()
					break
				}

				s += string(ss)
			}

			switch s {
			case "true", "false":
				tokens = append(tokens, Token{Type: TokenBoolean, Literal: s})
			case "null":
				tokens = append(tokens, Token{Type: TokenNull, Literal: s})
			default:
				return nil, &LexerError{Message: "syntax error", From: i, To: i + len(s)}
			}
		case unicode.IsSpace(r):
			continue
		default:
			return nil, &LexerError{Message: "invalid character", From: i, To: i}
		}
	}

	return tokens, nil
}
