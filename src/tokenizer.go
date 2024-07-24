package main

import (
	"unicode"
	"unicode/utf8"
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

type Tokenizer struct {
	src       string
	cursorPos int32
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
	if programTokenizer.cursorPos >= int32(len(programTokenizer.src)) {
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
	programTokenizer.cursorPos += int32(w)
	return char
}

func tokenize(programTokenizer Tokenizer) []Token {
	var tokenArr []Token

	var nextRune rune
	nextRune = programTokenizer.consume()
	for nextRune != 0 {
		if unicode.IsLetter(nextRune) {
			for unicode.IsLetter(nextRune) || unicode.IsDigit(nextRune) {
				programTokenizer.buf += string(nextRune)
				nextRune = programTokenizer.consume()
			}

			if programTokenizer.buf == "exit" {
				tokenArr = append(tokenArr, Token{typeOfToken: TokenType(0)})
			}
			programTokenizer.buf = ""
		} else if unicode.IsDigit(nextRune) {
			for unicode.IsDigit(nextRune) {
				programTokenizer.buf += string(nextRune)
				nextRune = programTokenizer.consume()
			}

			tokenArr = append(tokenArr, Token{typeOfToken: TokenType(1), value: programTokenizer.buf})
			programTokenizer.buf = ""
		} else if nextRune == '(' {
			programTokenizer.buf += string(nextRune)
			tokenArr = append(tokenArr, Token{typeOfToken: TokenType(3), value: programTokenizer.buf})
			programTokenizer.buf = ""
			nextRune = programTokenizer.consume()
		} else if nextRune == ')' {
			programTokenizer.buf += string(nextRune)
			tokenArr = append(tokenArr, Token{typeOfToken: TokenType(3), value: programTokenizer.buf})
			programTokenizer.buf = ""
			nextRune = programTokenizer.consume()
		} else if nextRune == ';' {
			programTokenizer.buf += string(nextRune)
			tokenArr = append(tokenArr, Token{typeOfToken: TokenType(3), value: programTokenizer.buf})
			programTokenizer.buf = ""
			nextRune = programTokenizer.consume()
		}
	}
	return tokenArr
}
