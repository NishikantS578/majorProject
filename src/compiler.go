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
	var fileContent string
	var err error

	fileBuffer, err = os.ReadFile(fileName)
	checkError(err)
	fileContent = string(fileBuffer)

	var tokenArr []Token
	var assemblyCode = ""

	var sourceProgram Tokenizer
	sourceProgram.initialize(fileContent)
	tokenArr = tokenize(sourceProgram)

	fmt.Println("==============================Tokens==============================")
	for _, t := range tokenArr {
		fmt.Println(STR_TOKEN_TYPE[t.typeOfToken], t.value)
	}

	os.WriteFile("./a.nasm", []byte(assemblyCode), 0644)
}
