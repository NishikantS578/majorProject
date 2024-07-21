package main

import (
	"fmt"
	"os"
)

type TokenType int

var STR_TOKEN_TYPE [3]string = [3]string{
	"return",
	"integerLiteral",
	"semicolon",
}

type Token struct {
	typeOfToken TokenType
	value       string
}

func tokenize(str string) []Token {
	var tokenArr []Token
	return tokenArr
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// var t Token

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
	tokenArr = tokenize(fileContent)
	for _, t := range tokenArr {
		fmt.Println(STR_TOKEN_TYPE[t.typeOfToken], t.value)
	}
}
