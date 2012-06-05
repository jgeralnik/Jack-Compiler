package token

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
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

var (
	keywords = []string{"class", "constructor", "function", "method", "field",
		"static", "var", "int", "char", "boolean", "void", "true",
		"false", "null", "this", "let", "do", "if", "else", "while",
		"return"}
)

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

func Read(path string) (tokens []Element, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	var line []byte

	for {
		pos := 0

		line, _, err = reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		for ; pos < len(line); pos++ {
			start := pos
			switch {
			case line[pos] == '\t', line[pos] == ' ', line[pos] == '\n':
				continue
			case isLetter(line[start]):
				for ; isLetter(line[pos]) || isDigit(line[pos]); pos++ {
				}
				pos--
				flag := false

				//Keyword
				for _, keyword := range keywords {
					if bytes.Equal([]byte(keyword), []byte(line[start:pos])) {
						tokens = append(tokens, Element{Keyword, string(line[start:pos])})
						flag = true
						continue
					}
				}

				if flag == false {
					//Identifier
					tokens = append(tokens, Element{Identifier, string(line[start:pos])})
				}

			case isDigit(line[pos]):
				//IntegerConstant
				for ; isDigit(line[pos]); pos++ {
				}
				pos--
				tokens = append(tokens, Element{IntegerConstant, string(line[start:pos])})

			//StringConstant
			case line[start] == '"':
				var result []byte
				pos++
				for ; ; pos++ {
					if line[pos] != '"' {
						result = append(result, line[pos])
					} else {
						if line[pos-1] == '\\' {
							result = append(result, []byte("&quot;")...)
						} else {
							break
						}
					}
				}
				pos++
				tokens = append(tokens, Element{StringConstant, string(result)})

			//Symbol
			case line[pos] == '<':
				tokens = append(tokens, Element{Symbol, "&lt;"})
			case line[pos] == '>':
				tokens = append(tokens, Element{Symbol, "&gt;"})
			case line[pos] == '&':
				tokens = append(tokens, Element{Symbol, "&amp;"})
			default:
				tokens = append(tokens, Element{Symbol, string(line[pos])})
			}
		}
	}
	return
}
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
