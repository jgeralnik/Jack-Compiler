package parser

import (
	"fmt"
	"os"
	"token"
)

var (
	//static []string
	argument []string
	//local []string
	className string
)

func CompileClass(tokens []token.Element, outputfile string) (err error) {
	pos := 0
	if tokens[pos].Tok != token.Keyword || tokens[pos].Value != "class" {
		panic("Attempted to compile non-class element with CompileClass")
	}

	output, err := os.Create(outputfile)
	if err != nil {
		return
	}
	defer output.Close()

	//output.WriteString("<class>\n")
	//defer output.WriteString("</class>\n")

	//output.WriteString(tokens[pos].String() + "\n") //Write 'class'
  pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write class name
	className = tokens[pos].Value
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write {
  pos++

	for ; pos < len(tokens)-1; pos++ { //Last token is }
		switch tokens[pos].Tok {
		case token.Keyword:
			switch tokens[pos].Value {
			case "static", "field":
				pos, err = compileClassVarDec(tokens, pos, output)
				if err != nil {
					return
				}
			case "constructor", "method", "function":
				pos, err = compileSubroutine(tokens, pos, output)
				if err != nil {
					return err
				}
			default:
				panic("Invalid keyword " + tokens[pos].Value + " in class")
			}
    case token.Identifier, token.IntegerConstant, token.StringConstant, token.Symbol:
      panic("Loose symbol " + tokens[pos].Value + " in class")
		}
	}

	return nil
}

//Global convention: all compile functions return pos as
//pointer to own last item

func compileClassVarDec(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	output.WriteString("<classVarDec>\n")
	defer output.WriteString("</classVarDec>\n")

	for pos = start; tokens[pos].Value != ";"; pos++ {
		output.WriteString(tokens[pos].String() + "\n")
	}
	output.WriteString(tokens[pos].String() + "\n")

	return pos, nil
}

func compileSubroutine(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<subroutineDec>\n")
	//defer output.WriteString("</subroutineDec>\n")

	pos = start
	//output.WriteString(tokens[pos].String() + "\n") //Write function/constructor/method
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write return value
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write function name
	output.WriteString(fmt.Sprintf("function %s.%s ", className, tokens[pos].Value))
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write (
	pos++

	pos, err = compileParameterList(tokens, pos, output)
	if err != nil {
		return
	}
	output.WriteString(fmt.Sprintf("%d\n", len(argument)))

	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write )
	pos++

	//output.WriteString("<subroutineBody>\n")
	//defer output.WriteString("</subroutineBody>\n")

	//output.WriteString(tokens[pos].String() + "\n") //Write {
	pos++

	for ; tokens[pos].Value == "var"; pos++ {
		pos, err = compileVarDec(tokens, pos, output)
		if err != nil {
			return
		}
	}

	pos, err = compileStatements(tokens, pos, output)
	if err != nil {
		return
	}

	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write }

	return
}

func compileParameterList(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<parameterList>\n")
	//defer output.WriteString("</parameterList>\n")
	argument = make([]string, 0)
	for pos = start; tokens[pos].Value != ")"; pos++ {
		//output.WriteString(tokens[pos].String() + "\n") //Print return value
		pos++
		//output.WriteString(tokens[pos].String() + "\n") //Print varname
		argument = append(argument, tokens[pos].Value)
	}
	pos--

	return pos, nil
}

func compileStatements(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<statements>\n")
	//defer output.WriteString("</statements>\n")

	for pos = start; tokens[pos].Value != "}"; pos++ {
		switch tokens[pos].Value {
		case "let":
			pos, err = compileLet(tokens, pos, output)
			if err != nil {
				return
			}
		case "do":
			pos, err = compileDo(tokens, pos, output)
			if err != nil {
				return
			}
		case "return":
			pos, err = compileReturn(tokens, pos, output)
			if err != nil {
				return
			}
		case "if":
			pos, err = compileIf(tokens, pos, output)
			if err != nil {
				return
			}
		case "while":
			pos, err = compileWhile(tokens, pos, output)
			if err != nil {
				return
			}

		default:
			panic(fmt.Sprintf("Invalid keyword %s in subroutine", tokens[pos].Value))
		}
	}

	pos--
	return
}
func compileVarDec(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	output.WriteString("<varDec>\n")
	defer output.WriteString("</varDec>\n")

	for pos = start; tokens[pos].Value != ";"; pos++ {
		output.WriteString(tokens[pos].String() + "\n")
	}
	output.WriteString(tokens[pos].String() + "\n")

	return pos, nil
}

func compileDo(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<doStatement>\n")
	//defer output.WriteString("</doStatement>\n")

  pos = start
  output.WriteString(tokens[pos].String() + "\n") //Write do
  pos++
  output.WriteString(tokens[pos].String() + "\n") //Write object
  pos++
  output.WriteString(tokens[pos].String() + "\n") //Write .
  pos++
  output.WriteString(tokens[pos].String() + "\n") //Write function name
  pos++
  output.WriteString(tokens[pos].String() + "\n") //Write (
  pos++

	pos, err = compileExpressionList(tokens, pos, output)
	if err != nil {
		return
	}

	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write )
	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write ;

	return pos, nil
}

func compileLet(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	output.WriteString("<letStatement>\n")
	defer output.WriteString("</letStatement>\n")

	pos = start

	output.WriteString(tokens[pos].String() + "\n") //Write let
	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write varname
	pos++

	if tokens[pos].Value == "[" {
		output.WriteString(tokens[pos].String() + "\n") //Write [
		pos++
		pos, err = compileExpression(tokens, pos, output)
		pos++
		output.WriteString(tokens[pos].String() + "\n") //Write ]
		pos++
	}

	output.WriteString(tokens[pos].String() + "\n") //Write =
	pos++

	pos, err = compileExpression(tokens, pos, output)
	if err != nil {
		return
	}

	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write ;

	return pos, nil
}

func compileWhile(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	output.WriteString("<whileStatement>\n")
	defer output.WriteString("</whileStatement>\n")

	for pos = start; tokens[pos-1].Value != "("; pos++ {
		output.WriteString(tokens[pos].String() + "\n")
	}

	pos, err = compileExpression(tokens, pos, output)
	if err != nil {
		return
	}

	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write )
	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write {
	pos++

	pos, err = compileStatements(tokens, pos, output)
	if err != nil {
		return
	}

	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write }

	return pos, nil
}

func compileReturn(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	output.WriteString("<returnStatement>\n")
	defer output.WriteString("</returnStatement>\n")

	pos = start
	output.WriteString(tokens[pos].String() + "\n") //Write return
	pos++

	if tokens[pos].Value != ";" {
		pos, err = compileExpression(tokens, pos, output)
		if err != nil {
			return
		}
		pos++
	}

	output.WriteString(tokens[pos].String() + "\n") //Write ;
	return
}

func compileIf(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	output.WriteString("<ifStatement>\n")
	defer output.WriteString("</ifStatement>\n")

	pos = start
	output.WriteString(tokens[pos].String() + "\n") //Write if
	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write (
	pos++

	pos, err = compileExpression(tokens, pos, output)
	if err != nil {
		return
	}
	pos++

	output.WriteString(tokens[pos].String() + "\n") //Write )
	pos++
	output.WriteString(tokens[pos].String() + "\n") //Write {
	pos++

	pos, err = compileStatements(tokens, pos, output)
	if err != nil {
		return
	}
	pos++

	output.WriteString(tokens[pos].String() + "\n") //Write }
	return
}

func isOp(tok token.Element) bool {
	switch tok.Value {
	case "+", "-", "*", "/", "&amp;", "|", "&lt;", "&gt;", "=":
		return true
	}
	return false
}
func compileExpression(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	output.WriteString("<expression>\n")
	defer output.WriteString("</expression>\n")

	pos, err = compileTerm(tokens, start, output)
	if err != nil {
		return
	}

	for pos++; isOp(tokens[pos]); pos++ {
		output.WriteString(tokens[pos].String() + "\n") //Write op
		pos++
		pos, err = compileTerm(tokens, pos, output)
		if err != nil {
			return
		}
	}

	pos--

	return
}

func compileTerm(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	output.WriteString("<term>\n")
	defer output.WriteString("</term>\n")

	pos = start

	switch tokens[pos].Tok {
	case token.IntegerConstant, token.StringConstant, token.Keyword:
		output.WriteString(tokens[pos].String() + "\n") //Write constant
	case token.Symbol:
		switch tokens[pos].Value {
		case "-", "~":
			output.WriteString(tokens[pos].String() + "\n") //Write unaryOp 
			pos, err = compileTerm(tokens, pos+1, output)
		case "(":
			output.WriteString(tokens[pos].String() + "\n") //Write ( 
			pos, err = compileExpression(tokens, pos+1, output)
			pos++
			output.WriteString(tokens[pos].String() + "\n") //Write ) 
		}
	case token.Identifier:
		output.WriteString(tokens[pos].String() + "\n") //Write identifier
		pos++
		switch tokens[pos].Value {
		case "(":
			output.WriteString(tokens[pos].String() + "\n") //Write ( 
			pos, err = compileExpressionList(tokens, pos+1, output)
			pos++
			output.WriteString(tokens[pos].String() + "\n") //Write ) 
		case ".":
			output.WriteString(tokens[pos].String() + "\n") //Write . 
			pos++
			output.WriteString(tokens[pos].String() + "\n") //Write identifier
			pos++
			output.WriteString(tokens[pos].String() + "\n") //Write ( 
			pos, err = compileExpressionList(tokens, pos+1, output)
			pos++
			output.WriteString(tokens[pos].String() + "\n") //Write ) 
		case "[":
			output.WriteString(tokens[pos].String() + "\n") //Write [ 
			pos, err = compileExpression(tokens, pos+1, output)
			pos++
			output.WriteString(tokens[pos].String() + "\n") //Write ] 
		default:
			pos--
		}
	default:
		panic("Unknown token type found in compileTerm")
	}

	return
}

func compileExpressionList(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	output.WriteString("<expressionList>\n")
	defer output.WriteString("</expressionList>\n")

	for pos = start; tokens[pos].Value != ")"; pos++ {
		pos, err = compileExpression(tokens, pos, output)
		if err != nil {
			return
		}

		if tokens[pos+1].Value != ")" {
			pos++
			if tokens[pos].Value != "," {
				panic("Expected , in expressionList, got " + tokens[pos].Value)
			}
			output.WriteString(tokens[pos].String() + "\n") //write ,
		}
	}
	pos--

	return
}
