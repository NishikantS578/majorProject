package main

import (
	"fmt"
	"os"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var cmdArgs []string = os.Args
	if len(cmdArgs) != 2 {
		fmt.Println("Incorrect usage.")
		fmt.Println("Correct usage: compiler <input.mnm>")
		return
	}

	var fileName string = cmdArgs[1]
	var fileBuffer []byte
	var err error

	fileBuffer, err = os.ReadFile(fileName)
	checkError(err)

	var assemblyCode = ""

	var tokenizer Tokenizer
	tokenizer.initialize(string(fileBuffer))

	var parser Parser
	parser.initialize(tokenizer.tokenize())

	var generator Generator
	generator.initialize(parser.parse())
	assemblyCode = generator.generate()

	os.WriteFile("./app.asm", []byte(assemblyCode), 0644)
	println("Compilation successfull")
}
