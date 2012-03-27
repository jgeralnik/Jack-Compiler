package vm2hack

import (
	"fmt"
	"strings"
)

func ConvertLine(line string) (result string) {
	words := strings.Fields(line)
	switch words[0] {
	case "push":
		result = push(words[1:])
	default:
		panic(fmt.Sprintf("Unknown command %s", words[0]))
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

func ConvertFile(inputfile string, outputfile string) (result []byte) {
	return
}
