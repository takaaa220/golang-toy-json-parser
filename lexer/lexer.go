package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenString
	TokenTrue
	TokenFalse
	TokenNull
	TokenEOF
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

type LexerError struct {
	Message string
	From    int
	To      int
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("%s, From: %d, To: %d", e.Message, e.From, e.To)
}

type Lexer struct {
	reader   *strings.Reader
	position int
}

func New(input string) *Lexer {
	return &Lexer{
		reader:   strings.NewReader(input),
		position: 0,
	}
}

func (lexer *Lexer) NextToken() (Token, error) {
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			break
		}
		lexer.position++

		switch {
		case r == '{':
			return Token{Type: TokenLeftBrace, Literal: string(r)}, nil
		case r == '}':
			return Token{Type: TokenRightBrace, Literal: string(r)}, nil
		case r == '[':
			return Token{Type: TokenLeftBracket, Literal: string(r)}, nil
		case r == ']':
			return Token{Type: TokenRightBracket, Literal: string(r)}, nil
		case r == ':':
			return Token{Type: TokenColon, Literal: string(r)}, nil
		case r == ',':
			return Token{Type: TokenComma, Literal: string(r)}, nil
		case r == '"':
			stringLiteral := ""

			start := lexer.position
			for {
				r, _, err := lexer.reader.ReadRune()
				if err != nil {
					return Token{}, &LexerError{Message: "syntax error", From: start - 1, To: lexer.position - 1}
				}
				lexer.position++

				if r == '"' {
					break
				}

				stringLiteral += string(r)
			}

			return Token{Type: TokenString, Literal: stringLiteral}, nil
		case unicode.IsDigit(r):
			numberLiteral := string(r)

			start := lexer.position
			for {
				ss, _, err := lexer.reader.ReadRune()
				if err != nil {
					break
				}
				lexer.position++

				if !unicode.IsDigit(ss) && ss != '.' {
					lexer.reader.UnreadRune()
					lexer.position--

					break
				}

				numberLiteral += string(ss)
			}

			_, err = strconv.ParseFloat(numberLiteral, 64)
			if err != nil {
				return Token{}, &LexerError{Message: fmt.Sprintf("unexpected number %s", numberLiteral), From: start - 1, To: lexer.position - 1}
			}

			return Token{Type: TokenNumber, Literal: numberLiteral}, nil
		case isLowerLetter(r):
			primitiveLiteral := string(r)

			start := lexer.position
			for {
				ss, _, err := lexer.reader.ReadRune()
				if err != nil {
					break
				}
				lexer.position++

				if !isLowerLetter(ss) {
					lexer.reader.UnreadRune()
					lexer.position--
					break
				}

				primitiveLiteral += string(ss)
			}

			// it may be better to check in parser not lexer whether literal value is valid or not
			switch primitiveLiteral {
			case "true":
				return Token{Type: TokenTrue, Literal: primitiveLiteral}, nil
			case "false":
				return Token{Type: TokenFalse, Literal: primitiveLiteral}, nil
			case "null":
				return Token{Type: TokenNull, Literal: primitiveLiteral}, nil
			default:
				return Token{}, &LexerError{Message: fmt.Sprintf("unexpected property '%s'", primitiveLiteral), From: start - 1, To: lexer.position - 1}
			}
		case unicode.IsSpace(r):
			continue
		default:
			return Token{}, &LexerError{Message: fmt.Sprintf("unexpected character '%s'", string(r)), From: lexer.position - 1, To: lexer.position - 1}
		}
	}

	return Token{Type: TokenEOF, Literal: ""}, nil
}

func isLowerLetter(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}
