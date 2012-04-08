package vm2hack

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
)

func ConvertLine(line string) (result string) {
	words := strings.Fields(line)
	if len(words) == 0 || strings.HasPrefix(words[0], "//") {
		return
	}

	switch words[0] {
	case "push":
		result = push(words[1:])
	case "add", "sub", "and", "or":
		result = fmt.Sprintf("@SP\nA=M-1\nD=M\n@SP\nM=M-1\nA=M-1\nM=M%sD\n", operators[words[0]])
	case "eq", "gt", "lt":
		result = fmt.Sprintf("@SP\nA=M-1\nD=M\n@SP\nM=M-1\nA=M-1\nD=M-D\n@LABEL%d\nD;J%s\n@SP\nA=M-1\nM=0\n@LABEL%d\n0;JMP\n(LABEL%d)\n@SP\nA=M-1\nM=-1\n(LABEL%d)\n", label, strings.ToUpper(words[0]), label+1, label, label+1)
		label += 2
	case "neg":
		result = "@SP\nA=M-1\nM=-M\n"
	default:
		panic(fmt.Sprintf("Unknown1 command %s", words[0]))
	}
	return
}

func push(words []string) (result string) {
	switch words[0] {
	case "constant":
		result = fmt.Sprintf("@%s\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", words[1])
	default:
		panic("Invalid push")
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

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		output.WriteString(ConvertLine(string(line)))
	}
	return nil
}
