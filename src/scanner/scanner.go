package main

import (
	"os"
	"parser"
	"path/filepath"
	"token"
)

func main() {
	inputInfo, err := os.Stat(os.Args[1])
	if err != nil {
		panic(err)
	}
	if inputInfo.IsDir() {
		filepath.Walk(os.Args[1], walk)
	} else {
		tokens, _ := token.Read(os.Args[1])
		if len(os.Args) > 2 {
			parser.CompileClass(tokens, os.Args[2])
		} else {
			dir, file := filepath.Split(os.Args[1])
			base := file[:len(file)-len(filepath.Ext(file))]
			parser.CompileClass(tokens, filepath.Join(dir, base+".vm"))
		}
	}
}

func walk(inputfile string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	if filepath.Ext(inputfile) == ".jack" {
		tokens, _ := token.Read(inputfile)
		dir, file := filepath.Split(inputfile)
		base := file[:len(file)-len(filepath.Ext(file))]
		parser.CompileClass(tokens, filepath.Join(dir, base+".vm"))
	}
	return nil
}
