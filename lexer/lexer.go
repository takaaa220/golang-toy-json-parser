package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	reader   *strings.Reader
	position int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		reader:   strings.NewReader(input),
		position: 0,
	}
}

func (lexer *Lexer) PeekToken() (Token, error) {
	currentPosition := lexer.position
	currentReader := *lexer.reader

	token, err := lexer.NextToken()

	lexer.position = currentPosition
	lexer.reader = &currentReader

	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func (lexer *Lexer) NextToken() (Token, error) {
	for {
		r, _, err := lexer.next()
		if err != nil {
			lexer.back()
			break
		}

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
				r, _, err := lexer.next()
				if err != nil {
					return Token{}, &LexerError{Message: "syntax error", From: start - 1, To: lexer.position - 2}
				}

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
				ss, _, err := lexer.next()
				if err != nil || !unicode.IsDigit(ss) && ss != '.' {
					lexer.back()

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
				ss, _, err := lexer.next()
				if err != nil || !isLowerLetter(ss) {
					lexer.back()
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

func (lexer *Lexer) next() (rune, int, error) {
	r, i, e := lexer.reader.ReadRune()
	lexer.position++

	return r, i, e
}

func (lexer *Lexer) back() error {
	e := lexer.reader.UnreadRune()
	lexer.position--

	return e
}

func isLowerLetter(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}
