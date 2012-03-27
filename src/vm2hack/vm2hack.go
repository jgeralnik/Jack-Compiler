package vm2hack

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ConvertLine(line string) (result string) {
	words := strings.Fields(line)
	if len(words) == 0 || strings.HasPrefix(words[0], "//") {
		return
	}

	switch words[0] {
	case "push":
		result = push(words[1:])
	case "add":
		result = "@SP\nA=M-1\nD=M\n@SP\nM=M-1\nA=M-1\nM=D+M\n"
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
