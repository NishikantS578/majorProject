package main

import (
	"os"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	EXIT_TOKEN = iota + 1
	LET_TOKEN
	IDENTIFIER_TOKEN
	ASSIGNMENT_OPERATOR_TOKEN
	PLUS_TOKEN
	SUBTRACTION_TOKEN
	MULTIPLICATION_TOKEN
	DIVISION_TOKEN
	INTEGER_LITERAL_TOKEN
	SEMICOLON_TOKEN
	OPENING_PARENTHESIS_TOKEN
	CLOSING_PARENTHESIS_TOKEN
)

var tokenName map[TokenType]string = map[TokenType]string{
	EXIT_TOKEN:                "Exit keyword",
	LET_TOKEN:                 "let keyword",
	IDENTIFIER_TOKEN:          "Identifier",
	ASSIGNMENT_OPERATOR_TOKEN: "Assignment operator",
	PLUS_TOKEN:                "Addition operator",
	SUBTRACTION_TOKEN:         "Subtraction operator",
	MULTIPLICATION_TOKEN:      "Multiplication operator",
	DIVISION_TOKEN:            "Division operator",
	INTEGER_LITERAL_TOKEN:     "Integer literal",
	SEMICOLON_TOKEN:           "Semicolon",
	OPENING_PARENTHESIS_TOKEN: "Opening parenthesis",
	CLOSING_PARENTHESIS_TOKEN: "Closing parenthesis",
}

type Token struct {
	typeOfToken TokenType
	value       string
}

type Tokenizer struct {
	src       string
	cursorPos int
	buf       string
}

func (programTokenizer *Tokenizer) initialize(src string) {
	programTokenizer.src = src
	programTokenizer.cursorPos = 0
	programTokenizer.buf = ""
}

func (programTokenizer *Tokenizer) peek() (rune, int) {
	var char rune
	var w int
	if programTokenizer.cursorPos >= len(programTokenizer.src) {
		return char, w
	} else {
		char, w = utf8.DecodeRuneInString(programTokenizer.src[programTokenizer.cursorPos:])
		return char, w
	}
}

func (programTokenizer *Tokenizer) consume() rune {
	var char rune
	var w int
	char, w = programTokenizer.peek()
	programTokenizer.cursorPos += w
	return char
}

func (programTokenizer *Tokenizer) tokenize() []Token {
	var tokenArr []Token

	var currentRune rune
	currentRune = programTokenizer.consume()
	for currentRune != 0 {
		if unicode.IsLetter(currentRune) {
			for (unicode.IsLetter(currentRune) && !unicode.IsSpace(currentRune)) || unicode.IsDigit(currentRune) {
				programTokenizer.buf += string(currentRune)
				currentRune = programTokenizer.consume()
			}

			if programTokenizer.buf == "exit" {
				tokenArr = append(tokenArr, Token{typeOfToken: EXIT_TOKEN, value: programTokenizer.buf})
			} else if programTokenizer.buf == "let" {
				tokenArr = append(tokenArr, Token{typeOfToken: LET_TOKEN})
			} else {
				tokenArr = append(tokenArr, Token{typeOfToken: IDENTIFIER_TOKEN, value: programTokenizer.buf})
			}
			programTokenizer.buf = ""
		} else if unicode.IsDigit(currentRune) {
			for unicode.IsDigit(currentRune) {
				programTokenizer.buf += string(currentRune)
				currentRune = programTokenizer.consume()
			}

			tokenArr = append(tokenArr, Token{typeOfToken: INTEGER_LITERAL_TOKEN, value: programTokenizer.buf})
			programTokenizer.buf = ""
		} else if currentRune == '=' {
			programTokenizer.buf += string(currentRune)
			tokenArr = append(tokenArr, Token{typeOfToken: ASSIGNMENT_OPERATOR_TOKEN, value: programTokenizer.buf})
			programTokenizer.buf = ""
			currentRune = programTokenizer.consume()
		} else if currentRune == '+' {
			programTokenizer.buf += string(currentRune)
			tokenArr = append(tokenArr, Token{typeOfToken: PLUS_TOKEN, value: programTokenizer.buf})
			programTokenizer.buf = ""
			currentRune = programTokenizer.consume()
		} else if currentRune == '(' {
			programTokenizer.buf += string(currentRune)
			tokenArr = append(tokenArr, Token{typeOfToken: OPENING_PARENTHESIS_TOKEN, value: programTokenizer.buf})
			programTokenizer.buf = ""
			currentRune = programTokenizer.consume()
		} else if currentRune == ')' {
			programTokenizer.buf += string(currentRune)
			tokenArr = append(tokenArr, Token{typeOfToken: CLOSING_PARENTHESIS_TOKEN, value: programTokenizer.buf})
			programTokenizer.buf = ""
			currentRune = programTokenizer.consume()
		} else if currentRune == ';' {
			programTokenizer.buf += string(currentRune)
			tokenArr = append(tokenArr, Token{typeOfToken: SEMICOLON_TOKEN, value: programTokenizer.buf})
			programTokenizer.buf = ""
			currentRune = programTokenizer.consume()
		} else if currentRune == ' ' || currentRune == '\n' {
			currentRune = programTokenizer.consume()
		} else {
			println("Undefined symbol: '" + string(currentRune) + "'")
			os.Exit(0)
		}
	}
	programTokenizer.cursorPos = 0
	return tokenArr
}
