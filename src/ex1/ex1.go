package main

import (
	"os"
	"path/filepath"
	"vm2hack"
)

func main() {
	inputInfo, err := os.Stat(os.Args[1])
	if err != nil {
		panic(err)
	}
	if inputInfo.IsDir() {
		if len(os.Args) > 2 {
			vm2hack.ConvertDirectory(os.Args[1], os.Args[2])
		} else {
			base := filepath.Base(os.Args[1])
			vm2hack.ConvertDirectory(os.Args[1], filepath.Join(os.Args[1], base+".asm"))
		}
	} else {
		if len(os.Args) > 2 {
			vm2hack.ConvertFile(os.Args[1], os.Args[2])
		} else {
			dir, file := filepath.Split(os.Args[1])
			base := file[:len(file)-len(filepath.Ext(file))]
			vm2hack.ConvertFile(os.Args[1], filepath.Join(dir, base+".asm"))
		}
	}
}
