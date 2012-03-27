package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func walk(root string) filepath.WalkFunc {
	count := 0
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if path == root {
				return nil
			}
			return filepath.SkipDir
		}

		if path == root+"/hello.in" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			reader := bufio.NewReader(file)

			output, err := os.Create(root + "/hello.out")
			if err != nil {
				return err
			}
			defer output.Close()

			var line []byte
			for {
				cur, isPrefix, err := reader.ReadLine()
				line = append(line, cur...)
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
				if !isPrefix {
					line = append(line, '\n')
					if bytes.Contains(line, []byte("you")) {
						os.Stdout.Write(line)
					}
					_, err = output.Write(line)
					if err != nil {
						return err
					}
					line = line[0:0]
				}
			}

			return nil
		}

		result, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(path, os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.WriteString(strconv.Itoa(count))
		if err != nil {
			return err
		}
		_, err = file.Write(result)
		if err != nil {
			return err
		}
		count++
		return nil
	}
}

func main() {
	var root string
	//fmt.Scanln(&root)
	root = "/home/joey/goProject/src/ex0/inputs"
	err := filepath.Walk(root, walk(root))
	if err != nil {
		fmt.Println(err)
	}
}
