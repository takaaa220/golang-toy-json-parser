package lexer

import "fmt"

type LexerError struct {
	Message string
	From    int
	To      int
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("%s, From: %d, To: %d", e.Message, e.From, e.To)
}
