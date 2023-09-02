package parser

import "fmt"

type ParserError struct {
	Message string
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
