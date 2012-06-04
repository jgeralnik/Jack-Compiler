package token

import (
	"fmt"
)

type Token int

const (
	Keyword Token = iota
	Identifier
	Symbol
	IntegerConstant
	StringConstant
)

type Element struct {
	Tok   Token
	Value string
}

func (t Token) String() string {
	switch t {
	case Keyword:
		return "keyword"
	case Identifier:
		return "identifier"
	case Symbol:
		return "symbol"
	case IntegerConstant:
		return "integerConstant"
	case StringConstant:
		return "stringConstant"
	}
	return fmt.Sprintf("UnknownToken%d", t)
}

func (e Element) String() string {
	return fmt.Sprintf("<%s>%s</%s>", e.Tok, e.Value, e.Tok)
}
