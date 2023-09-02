package parser

import (
	"fmt"
	"strconv"

	"github.com/takaaa220/golang-toy-json-parser/lexer"
)

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(lexer *lexer.Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) Parse() (interface{}, error) {
	res, err := p.value()
	if err != nil {
		return nil, err
	}

	token, err := p.lexer.NextToken()
	if err != nil {
		return nil, err
	}
	if token.Type != lexer.TokenEOF {
		return nil, &ParserError{Message: fmt.Sprintf("unexpected token %v", token)}
	}

	return res, nil
}

/*
`VALUE` = `PRIMITIVE`|`OBJECT`|`ARRAY`
*/
func (p *Parser) value() (interface{}, error) {
	token, err := p.lexer.PeekToken()
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case lexer.TokenLeftBracket:
		return p.array()
	case lexer.TokenLeftBrace:
		return p.object()
	default:
		return p.primitive()
	}
}

/*
`PRIMITIVE` = `NULL`|`BOOL`|`NUM`|`STR`
*/
func (p *Parser) primitive() (interface{}, error) {
	token, err := p.lexer.NextToken()
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case lexer.TokenNumber:
		parsed, err := strconv.ParseFloat(token.Literal, 64)
		if err != nil {
			return nil, &ParserError{Message: fmt.Sprintf("invalid number %s", token.Literal)}
		}

		return parsed, nil
	case lexer.TokenString:
		return token.Literal, nil
	case lexer.TokenTrue:
		return true, nil
	case lexer.TokenFalse:
		return false, nil
	case lexer.TokenNull:
		return nil, nil
	default:
		return nil, &ParserError{Message: fmt.Sprintf("unexpected token %v", token)}
	}
}

/*
`ARRAY` = `LEFT_BRACKET` - `VALUE* - (COMMA - VALUE)*` - `RIGHT_BRACKET`
*/
func (p *Parser) array() ([]interface{}, error) {
	leftBracketToken, err := p.lexer.NextToken()
	if err != nil {
		return nil, err
	}
	if leftBracketToken.Type != lexer.TokenLeftBracket {
		return nil, &ParserError{Message: fmt.Sprintf("unexpected token %v", leftBracketToken)}
	}

	values := []interface{}{}

	for i := 0; ; i++ {
		token1, err := p.lexer.PeekToken()
		if err != nil {
			return nil, err
		}
		if token1.Type == lexer.TokenRightBracket {
			_, err := p.lexer.NextToken()
			if err != nil {
				return nil, err
			}

			return values, nil
		}

		if i > 0 {
			commaToken, err := p.lexer.NextToken()
			if err != nil {
				return nil, err
			}
			if commaToken.Type != lexer.TokenComma {
				return nil, &ParserError{Message: fmt.Sprintf("unexpected token %v", commaToken)}
			}
		}

		value, err := p.value()
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
}

/*
`OBJECT` = `LEFT_BRACE` - `PROPERTY*` - `(COMMA - PROPERTY)*` - `RIGHT_BRACE`
*/
func (p *Parser) object() (map[string]interface{}, error) {
	leftBraceToken, err := p.lexer.NextToken()
	if err != nil {
		return nil, err
	}
	if leftBraceToken.Type != lexer.TokenLeftBrace {
		return nil, &ParserError{Message: fmt.Sprintf("unexpected token %v", leftBraceToken)}
	}

	object := map[string]interface{}{}

	for i := 0; ; i++ {
		rightBraceToken, err := p.lexer.PeekToken()
		if err != nil {
			return nil, err
		}
		if rightBraceToken.Type == lexer.TokenRightBrace {
			p.lexer.NextToken()

			return object, nil
		}

		if i > 0 {
			commaToken, err := p.lexer.NextToken()
			if err != nil {
				return nil, err
			}
			if commaToken.Type != lexer.TokenComma {
				return nil, &ParserError{Message: fmt.Sprintf("unexpected token %v", commaToken)}
			}
		}

		key, value, err := p.property()
		if err != nil {
			return nil, err
		}

		object[key] = value
	}
}

/*
`PROPERTY`  = `STR` - `COLON` - `VALUE`
*/
func (p *Parser) property() (string, interface{}, error) {
	keyToken, err := p.lexer.NextToken()
	if err != nil {
		return "", nil, err
	}
	if keyToken.Type != lexer.TokenString {
		return "", nil, &ParserError{Message: fmt.Sprintf("unexpected token %v", keyToken)}
	}

	colonToken, err := p.lexer.NextToken()
	if err != nil {
		return "", nil, err
	}
	if colonToken.Type != lexer.TokenColon {
		return "", nil, &ParserError{Message: fmt.Sprintf("unexpected token %v", colonToken)}
	}

	value, err := p.value()
	if err != nil {
		return "", nil, err
	}

	return keyToken.Literal, value, nil
}
