package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type LexerError struct {
	Message string
	From    int
	To      int
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("%s, From: %d, To: %d", e.Message, e.From, e.To)
}

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenString
	TokenTrue
	TokenFalse
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
	reader := strings.NewReader(input)
	position := 0
	tokens := []Token{}

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		position++

		switch {
		case r == '{':
			tokens = append(tokens, Token{Type: TokenLeftBrace, Literal: string(r)})
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
			stringLiteral := ""

			p := 0
			for {
				r, _, err := reader.ReadRune()
				if err != nil {
					return nil, &LexerError{Message: "syntax error", From: position - 1, To: position - 1 + p}
				}

				p++

				if r == '"' {
					position += p
					break
				}

				stringLiteral += string(r)
			}

			tokens = append(tokens, Token{Type: TokenString, Literal: stringLiteral})
		case unicode.IsDigit(r):
			numberLiteral := string(r)

			p := 0
			for {
				ss, _, err := reader.ReadRune()
				if err != nil {
					break
				}

				p++

				if !unicode.IsDigit(ss) {
					reader.UnreadRune()
					position += p - 1

					break
				}

				numberLiteral += string(ss)
			}

			tokens = append(tokens, Token{Type: TokenNumber, Literal: numberLiteral})
		case isLowerLetter(r):
			primitiveLiteral := string(r)

			for {
				ss, _, err := reader.ReadRune()
				if err != nil {
					break
				}

				if !isLowerLetter(ss) {
					reader.UnreadRune()
					break
				}

				primitiveLiteral += string(ss)
			}

			switch primitiveLiteral {
			case "true":
				tokens = append(tokens, Token{Type: TokenTrue, Literal: primitiveLiteral})
			case "false":
				tokens = append(tokens, Token{Type: TokenFalse, Literal: primitiveLiteral})
			case "null":
				tokens = append(tokens, Token{Type: TokenNull, Literal: primitiveLiteral})
			default:
				return nil, &LexerError{Message: fmt.Sprintf("unexpected property '%s'", primitiveLiteral), From: position - 1, To: position - 1 + len(primitiveLiteral) - 1}
			}

			position = position + len(primitiveLiteral)
		case unicode.IsSpace(r):
			continue
		default:
			return nil, &LexerError{Message: fmt.Sprintf("unexpected character '%s'", string(r)), From: position - 1, To: position - 1}
		}
	}

	return tokens, nil
}
