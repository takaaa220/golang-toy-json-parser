package lexer

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
