package parser

import (
	"fmt"
	"os"
	"strings"
	"token"
)

var (
	//static []string
	argument  []variable
	local     []variable
	className string
	label     int = 0
)

type variable struct {
	vartype string
	name    string
}

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
	local = make([]variable, 0)
	argument = make([]variable, 0)

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
	output.WriteString(fmt.Sprintf("%d\n", len(local)))

	pos, err = compileStatements(tokens, pos, output)
	if err != nil {
		return
	}

	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write }
	output.WriteString("return\n\n")

	return
}

func compileParameterList(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<parameterList>\n")
	//defer output.WriteString("</parameterList>\n")
	for pos = start; tokens[pos].Value != ")"; pos++ {
		if tokens[pos].Value == "," {
			pos++
		}
		//output.WriteString(tokens[pos].String() + "\n") //Print type 
		vartype := tokens[pos].Value
		pos++
		//output.WriteString(tokens[pos].String() + "\n") //Print varname
		argument = append(argument, variable{vartype, tokens[pos].Value})
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
	//output.WriteString("<varDec>\n")
	//defer output.WriteString("</varDec>\n")

	pos = start

	pos++ //skip over var
	vartype := tokens[pos].Value

	for ; tokens[pos].Value != ";"; pos++ {
		pos++ //skip over vartype or comma
		local = append(local, variable{vartype, tokens[pos].Value})
	}

	return pos, nil
}

func compileDo(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<doStatement>\n")
	//defer output.WriteString("</doStatement>\n")
	var funcname string
	pos = start
	//output.WriteString(tokens[pos].String() + "\n") //Write do
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write object
	if isUpper(tokens[pos].Value) {
		//static object
		funcname = tokens[pos].Value + "." + tokens[pos+2].Value
		pos++
		//output.WriteString(tokens[pos].String() + "\n") //Write .
		pos++
		//output.WriteString(tokens[pos].String() + "\n") //Write function name
		pos++
	} else {
		fmt.Print(tokens[pos].Value)
		panic("Can't yet compile object functions!")
	}
	//output.WriteString(tokens[pos].String() + "\n") //Write (
	pos++

	pos, count, err := compileExpressionList(tokens, pos, output)
	if err != nil {
		return
	}

	pos++
	output.WriteString(fmt.Sprintf("call %s %d\n", funcname, count))
	output.WriteString("pop temp 0\n")
	//output.WriteString(tokens[pos].String() + "\n") //Write )
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write ;

	return pos, nil
}

func compileLet(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<letStatement>\n")
	//defer output.WriteString("</letStatement>\n")

	pos = start

	//output.WriteString(tokens[pos].String() + "\n") //Write let
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write varname
	myvar := ""
	for index, item := range local {
		if item.name == tokens[pos].Value {
			myvar = fmt.Sprintf("local %d", index)
		}
	}
	if myvar == "" {
		for index, item := range argument {
			if item.name == tokens[pos].Value {
				myvar = fmt.Sprintf("argument %d", index)
			}
		}
	}

	if myvar == "" {
		e := "Can't find variable " + tokens[pos].Value + "\n"
		e += "Local is:\n"
		for _, item := range local {
			e += fmt.Sprintf("%s\n", item)
		}
		e += "Argument is:\n"
		for _, item := range argument {
			e += fmt.Sprintf("%s\n", item)
		}
		panic(e)
	}

	pos++

	if tokens[pos].Value == "[" {
		panic("Arrays not yet supported in let statement!")
		output.WriteString(tokens[pos].String() + "\n") //Write [
		pos++
		pos, err = compileExpression(tokens, pos, output)
		pos++
		output.WriteString(tokens[pos].String() + "\n") //Write ]
		pos++
	}

	//output.WriteString(tokens[pos].String() + "\n") //Write =
	pos++

	pos, err = compileExpression(tokens, pos, output)
	if err != nil {
		return
	}

	output.WriteString("pop " + myvar + "\n")
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write ;

	return pos, nil
}

func compileWhile(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<whileStatement>\n")
	//defer output.WriteString("</whileStatement>\n")

	pos = start
	//output.WriteString(tokens[pos].String() + "\n") //Write while
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write (
	pos++

	mylabel := label
	label++
	output.WriteString(fmt.Sprintf("label WHILE%d\n", mylabel))
	pos, err = compileExpression(tokens, pos, output)
	if err != nil {
		return
	}

	output.WriteString(fmt.Sprintf("not\nif-goto END%d\n", mylabel))
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write )
	pos++

	//output.WriteString(tokens[pos].String() + "\n") //Write {
	pos++

	pos, err = compileStatements(tokens, pos, output)
	if err != nil {
		return
	}

	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write }

	output.WriteString(fmt.Sprintf("goto WHILE%d\n", mylabel))
	output.WriteString(fmt.Sprintf("label END%d\n", mylabel))
	return pos, nil
}

func compileReturn(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<returnStatement>\n")
	//defer output.WriteString("</returnStatement>\n")

	pos = start
	//output.WriteString(tokens[pos].String() + "\n") //Write return
	pos++

	if tokens[pos].Value != ";" {
		pos, err = compileExpression(tokens, pos, output)
		if err != nil {
			return
		}
		pos++
	} else {
		output.WriteString("push constant 0\n")
	}

	//output.WriteString(tokens[pos].String() + "\n") //Write ;
	return
}

func compileIf(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<ifStatement>\n")
	//defer output.WriteString("</ifStatement>\n")

	pos = start
	//output.WriteString(tokens[pos].String() + "\n") //Write if
	pos++
	//output.WriteString(tokens[pos].String() + "\n") //Write (
	pos++

	pos, err = compileExpression(tokens, pos, output)
	if err != nil {
		return
	}
	pos++

	//output.WriteString(tokens[pos].String() + "\n") //Write )
	pos++

	mylabel := label
	label++
	output.WriteString(fmt.Sprintf("not\nif-goto ELSE%d\n", mylabel))

	//output.WriteString(tokens[pos].String() + "\n") //Write {
	pos++

	pos, err = compileStatements(tokens, pos, output)
	if err != nil {
		return
	}
	pos++

	//output.WriteString(tokens[pos].String() + "\n") //Write }
	pos++

	output.WriteString(fmt.Sprintf("goto END%d\n", mylabel))
	output.WriteString(fmt.Sprintf("label ELSE%d\n", mylabel))

	if tokens[pos].Value == "else" {
		//output.WriteString(tokens[pos].String() + "\n") //Write else
		pos++
		//output.WriteString(tokens[pos].String() + "\n") //Write {
		pos++

		pos, err = compileStatements(tokens, pos, output)
		if err != nil {
			return
		}
		pos++

		//output.WriteString(tokens[pos].String() + "\n") //Write }
		pos++
	}

	output.WriteString(fmt.Sprintf("label END%d\n", mylabel))

	pos--
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
	//output.WriteString("<expression>\n")
	//defer output.WriteString("</expression>\n")

	pos, err = compileTerm(tokens, start, output)
	if err != nil {
		return
	}

	for pos++; isOp(tokens[pos]); pos++ {
		//output.WriteString(tokens[pos].String() + "\n") //Write op
		op := tokens[pos].Value
		pos++
		pos, err = compileTerm(tokens, pos, output)
		var ops = map[string]string{
			"+":     "add",
			"-":     "sub",
			"*":     "call Math.multiply 2",
			"/":     "call Math.divide 2",
			"&amp;": "and",
			"|":     "or",
			"&lt;":  "lt",
			"&gt;":  "gt",
			"=":     "eq",
		}

		if _, ok := ops[op]; !ok {
			panic("Illegal operator " + tokens[pos-1].Value)
		}

		output.WriteString(ops[op])
		output.WriteString("\n")

		if err != nil {
			return
		}
	}

	pos--
	return
}

func compileTerm(tokens []token.Element, start int, output *os.File) (pos int, err error) {
	//output.WriteString("<term>\n")
	//defer output.WriteString("</term>\n")

	pos = start

	switch tokens[pos].Tok {
	case token.IntegerConstant:
		//output.WriteString(tokens[pos].String() + "\n") //Write constant
		output.WriteString(fmt.Sprintf("push constant %s\n", tokens[pos].Value))
	case token.Keyword:
		switch tokens[pos].Value {
		case "true":
			output.WriteString("push constant 0\nnot\n")
		case "false":
			output.WriteString("push constant 0\n")
		default:
			panic("Invalid keyword " + tokens[pos].Value + "in terminal")
		}
	case token.StringConstant:
		panic("I don't know what to do with " + tokens[pos].Value)
	case token.Symbol:
		switch tokens[pos].Value {
		case "-":
			pos, err = compileTerm(tokens, pos+1, output)
			output.WriteString("neg\n")
		case "~":
			pos, err = compileTerm(tokens, pos+1, output)
			output.WriteString("not\n")
		case "(":
			//output.WriteString(tokens[pos].String() + "\n") //Write ( 
			pos, err = compileExpression(tokens, pos+1, output)
			pos++
			//output.WriteString(tokens[pos].String() + "\n") //Write ) 
		}
	case token.Identifier:
		//output.WriteString(tokens[pos].String() + "\n") //Write identifier
		pos++
		switch tokens[pos].Value {
		case "(":
			//Previous item was name of local function
			panic("Give me an example!")
			output.WriteString(tokens[pos].String() + "\n") //Write ( 
			pos, _, err = compileExpressionList(tokens, pos+1, output)
			pos++
			output.WriteString(tokens[pos].String() + "\n") //Write ) 
		case ".":
			//Previous item was classname
			funcname := tokens[pos-1].Value
			//output.WriteString(tokens[pos].String() + "\n") //Write . 
			pos++
			//output.WriteString(tokens[pos].String() + "\n") //Write identifier
			funcname += "." + tokens[pos].Value
			pos++
			//output.WriteString(tokens[pos].String() + "\n") //Write ( 

			var count int
			pos, count, err = compileExpressionList(tokens, pos+1, output)

			pos++
			output.WriteString(fmt.Sprintf("call %s %d\n", funcname, count))
			//output.WriteString(tokens[pos].String() + "\n") //Write ) 

		case "[":
			panic("Not implemented")
			output.WriteString(tokens[pos].String() + "\n") //Write [ 
			pos, err = compileExpression(tokens, pos+1, output)
			pos++
			output.WriteString(tokens[pos].String() + "\n") //Write ] 
		default:
			pos--
			myvar := ""
			for index, item := range local {
				if item.name == tokens[pos].Value {
					myvar = fmt.Sprintf("local %d", index)
				}
			}
			if myvar == "" {
				for index, item := range argument {
					if item.name == tokens[pos].Value {
						myvar = fmt.Sprintf("argument %d", index)
					}
				}
			}

			if myvar == "" {
				e := "Can't find variable " + tokens[pos].Value + "\n"
				e += "Local is:\n"
				for _, item := range local {
					e += fmt.Sprintf("%s\n", item)
				}
				e += "Argument is:\n"
				for _, item := range argument {
					e += fmt.Sprintf("%s\n", item)
				}
				panic(e)
			}
			output.WriteString("push " + myvar + "\n")
		}
	default:
		panic("Unknown token type found in compileTerm")
	}

	return
}

func compileExpressionList(tokens []token.Element, start int, output *os.File) (pos int, count int, err error) {
	//output.WriteString("<expressionList>\n")
	//defer output.WriteString("</expressionList>\n")
	count = 0
	for pos = start; tokens[pos].Value != ")"; pos++ {
		pos, err = compileExpression(tokens, pos, output)
		count++
		if err != nil {
			return
		}

		if tokens[pos+1].Value != ")" {
			pos++
			if tokens[pos].Value != "," {
				panic("Expected , in expressionList, got " + tokens[pos].Value)
			}
			//output.WriteString(tokens[pos].String() + "\n") //write ,
		}
	}
	pos--

	return
}

func isUpper(word string) bool {
	return word[0] == strings.ToUpper(word)[0]
}
