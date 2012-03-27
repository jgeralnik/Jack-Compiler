package main

import (
	"os"
	"vm2hack"
)

func main() {
	vm2hack.ConvertFile(os.Args[1], os.Args[2])
}
