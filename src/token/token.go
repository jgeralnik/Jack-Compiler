package token

type Token int

const (
	Keyword Token = iota
	Identifier
	Symbol
	IntegerConstant
	StringConstant
)

type Element struct {
	token Token
	value string
}
