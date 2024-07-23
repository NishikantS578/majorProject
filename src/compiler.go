package main

import (
	"fmt"
	"os"
	"unicode"
)

type TokenType int

var STR_TOKEN_TYPE []string = []string{
	"exit",
	"integerLiteral",
	"return",
	"separator",
}

type Token struct {
	typeOfToken TokenType
	value       string
}

func tokenize(str string) []Token {
	var tokenArr []Token
	var buf string = ""
	var ch rune
	for idx := range str {
		ch = rune(str[idx])

		if unicode.IsLetter(ch) {
			for unicode.IsLetter(ch) || unicode.IsDigit(ch) {
				buf += string(ch)
				idx += 1
				ch = rune(str[idx])
			}
			idx -= 1

			if buf == "exit" {
				tokenArr = append(tokenArr, Token{typeOfToken: TokenType(0)})
			}
			buf = ""
		} else if unicode.IsDigit(ch) {
			for unicode.IsDigit(ch) {
				buf += string(ch)
				idx += 1
				ch = rune(str[idx])
			}
			idx -= 1

			tokenArr = append(tokenArr, Token{typeOfToken: TokenType(1), value: buf})
			buf = ""
		} else if ch == '(' {
			buf += string(ch)
			tokenArr = append(tokenArr, Token{typeOfToken: TokenType(3), value: buf})
			buf = ""
		} else if ch == ')' {
			buf += string(ch)
			tokenArr = append(tokenArr, Token{typeOfToken: TokenType(3), value: buf})
			buf = ""
		} else if ch == ';' {
			buf += string(ch)
			tokenArr = append(tokenArr, Token{typeOfToken: TokenType(3), value: buf})
			buf = ""
		}
	}
	return tokenArr
}

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
	tokenArr = tokenize(fileContent)

	fmt.Println("==============================Tokens==============================")
	for _, t := range tokenArr {
		fmt.Println(STR_TOKEN_TYPE[t.typeOfToken], t.value)
	}

	os.WriteFile("./a.nasm", []byte(assemblyCode), 0644)
}
