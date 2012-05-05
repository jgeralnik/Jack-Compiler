package main

import (
	"os"
	"path/filepath"
	"vm2hack"
)

func main() {
	if len(os.Args) > 2 {
		vm2hack.ConvertFile(os.Args[1], os.Args[2])
	} else {
		base := filepath.Base(os.Args[1][:len(os.Args[1])-len(filepath.Ext(os.Args[1]))])
		vm2hack.ConvertFile(os.Args[1], filepath.Join(filepath.Dir(os.Args[1]), base+".asm"))
	}
}
