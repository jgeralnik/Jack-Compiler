package vm2hack

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	label     int = 0
	operators     = map[string]string{
		"add": "+",
		"sub": "-",
		"and": "&",
		"or":  "|",
	}
	segments = map[string]string{
		"local":    "LCL",
		"argument": "ARG",
		"this":     "THIS",
		"that":     "THAT",
	}
	bases = map[string]int{
		"temp":    5,
		"pointer": 3,
	}
	function string
)

func ConvertLine(line string, filename string) (result string) {
	words := strings.Fields(line)
	if len(words) == 0 || strings.HasPrefix(words[0], "//") {
		return
	}

	switch words[0] {
	case "push":
		result = push(words[1:], filename)
	case "pop":
		result = pop(words[1:], filename)
	case "add", "sub", "and", "or":
		result = fmt.Sprintf("@SP\nA=M-1\nD=M\n@SP\nM=M-1\nA=M-1\nM=M%sD\n", operators[words[0]])
	case "eq", "gt", "lt":
		result = fmt.Sprintf("@SP\nA=M-1\nD=M\n@SP\nM=M-1\nA=M-1\nD=M-D\n@LABEL%d\nD;J%s\n@SP\nA=M-1\nM=0\n@LABEL%d\n0;JMP\n(LABEL%d)\n@SP\nA=M-1\nM=-1\n(LABEL%d)\n", label, strings.ToUpper(words[0]), label+1, label, label+1)
		label += 2
	case "neg":
		result = "@SP\nA=M-1\nM=-M\n"
	case "not":
		result = "@SP\nA=M-1\nM=!M\n"
	case "label":
		result = fmt.Sprintf("(%s$%s)\n", function, words[1])
	case "goto":
		result = fmt.Sprintf("@%s$%s\n0;JMP\n", function, words[1])
	case "if-goto":
		result = fmt.Sprintf("@SP\nM=M-1\nA=M\nD=M\n@%s$%s\nD;JNE\n", function, words[1])
	case "function":
		result = fmt.Sprintf("(%s)\n@SP\nA=M\n", words[1])
		k, _ := strconv.Atoi(words[2])
		for i := 0; i < k; i++ {
			result += "M=0\nA=A+1\n"
		}
		result += "D=A\n@SP\nM=D\n"
	case "call":
		result = fmt.Sprintf("@LABEL%d\nD=A\n@SP\nM=M+1\nA=M-1\nM=D\n", label) //Push return address
		segments := []string{"LCL", "ARG", "THIS", "THAT"}
		for _, segment := range segments {
			result += fmt.Sprintf("@%s\nD=M\n@SP\nM=M+1\nA=M-1\nM=D\n", segment) //Push LCL,ARG,THIS,THAT
		}
		n, _ := strconv.Atoi(words[2])
		result += fmt.Sprintf("@SP\nD=M\n@%d\nD=D-A\n@ARG\nM=D\n", n+5)   //ARG = SP-n-5
		result += "@SP\nD=M\n@LCL\nM=D\n"                                 //LCL=SP
		result += fmt.Sprintf("@%s\n0;JMP\n(LABEL%d)\n", words[1], label) //Jump to function
		label += 1
	case "return":
		result = "@LCL\nD=M\n@13\nM=D\n"                   //Store frame in temp
		result += "@5\nA=D-A\nD=M\n@14\nM=D\n"             //Store RET in temp
		result += "@SP\nM=M-1\nA=M\nD=M\n@ARG\nA=M\nM=D\n" //*ARG = pop()
		result += "@ARG\nD=M+1\n@SP\nM=D\n"                //*SP = ARG+1
		segments := []string{"THAT", "THIS", "ARG", "LCL"}
		for _, segment := range segments {
			result += fmt.Sprintf("@13\nM=M-1\nA=M\nD=M\n@%s\nM=D\n", segment) //Pop THAT,THIS,ARG,LCL
		}
		result += "@14\nA=M\n0;JMP\n" //Goto ret
	default:
		panic(fmt.Sprintf("Unknown1 command %s", words[0]))
	}
	return
}

func push(words []string, filename string) (result string) {
	switch words[0] {
	case "constant":
		result = fmt.Sprintf("@%s\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", words[1])
	case "local", "argument", "this", "that":
		result = fmt.Sprintf("@%s\nD=M\n@%s\nA=D+A\nD=M\n@SP\nM=M+1\nA=M-1\nM=D\n", segments[words[0]], words[1])
	case "temp", "pointer":
		offset, err := strconv.Atoi(words[1])
		if err != nil {
			panic(fmt.Sprintf("Non number passed to %s", words[0]))
		}
		result = fmt.Sprintf("@%d\nD=M\n@SP\nM=M+1\nA=M-1\nM=D\n", offset+bases[words[0]])
	case "static":
		result = fmt.Sprintf("@%s.%s\nD=M\n@SP\nM=M+1\nA=M-1\nM=D\n", filename, words[1])
	default:
		panic(fmt.Sprintf("Invalid push: %s", words[0]))
	}
	return
}

func pop(words []string, filename string) (result string) {
	switch words[0] {
	case "local", "argument", "this", "that":
		result = fmt.Sprintf("@%s\nD=M\n@%s\nD=D+A\n@SP\nM=M-1\nA=M\nD=D+M\nA=D-M\nD=D-A\nM=D\n", segments[words[0]], words[1])
	case "temp", "pointer":
		offset, err := strconv.Atoi(words[1])
		if err != nil {
			panic(fmt.Sprintf("Non number passed to %s", words[0]))
		}
		result = fmt.Sprintf("@SP\nM=M-1\nA=M\nD=M\n@%d\nM=D\n", offset+bases[words[0]])
	case "static":
		result = fmt.Sprintf("@SP\nM=M-1\nA=M\nD=M\n@%s.%s\nM=D\n", filename, words[1])
	default:
		panic(fmt.Sprintf("Invalid pop: %s", words[0]))
	}
	return
}

func ConvertFile(inputfile string, outputfile string) error {
	file, err := os.Open(inputfile)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	output, err := os.Create(outputfile)
	if err != nil {
		return err
	}
	defer output.Close()

	filename := filepath.Base(inputfile)
	filename = filename[:len(filename)-3] //Remove ".vm" from filename

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		output.WriteString(ConvertLine(string(line), filename))
	}
	return nil
}

func walk(outputfile *os.File) filepath.WalkFunc {
	return func(inputfile string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(inputfile) == ".vm" {
			file, err := os.Open(inputfile)
			if err != nil {
				return err
			}
			defer file.Close()
			reader := bufio.NewReader(file)

			filename := filepath.Base(inputfile)
			filename = filename[:len(filename)-3] //Remove ".vm" from filename

			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
				outputfile.WriteString(ConvertLine(string(line), filename))
			}
		}
		return nil
	}
}

func ConvertDirectory(directory string, outputfile string) error {
	file, err := os.Create(outputfile)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("@256\nD=A\n@SP\nM=D\n")
	file.WriteString(ConvertLine("call Sys.init 0", "system"))
	return filepath.Walk(directory, walk(file))
}
