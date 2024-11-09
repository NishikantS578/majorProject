package lexer

import (
	"majorProject/compiler/parser"
	"unicode"
	"unicode/utf8"
)

type Tokenizer struct {
	TokenArr  []parser.Token
	src       string
	cursorPos int
	buf       string
}

func New(line string) Tokenizer {
	return Tokenizer{
		TokenArr:  []parser.Token{},
		src:       line,
		cursorPos: 0,
		buf:       "",
	}
}

func (tokenizer *Tokenizer) SetNewInput(src string) {
	tokenizer.src = src
	tokenizer.cursorPos = 0
	tokenizer.buf = ""
	tokenizer.TokenArr = []parser.Token{}
}

func (tokenizer *Tokenizer) Tokenize() {
	var currentCh, w = tokenizer.peek()

	for currentCh != 0 {
		if unicode.IsDigit(currentCh) {
			tokenizer.buf += string(currentCh)
			tokenizer.readCh(w)
			currentCh, w = tokenizer.peek()

			for unicode.IsDigit(currentCh) {
				tokenizer.buf += string(currentCh)
				tokenizer.readCh(w)
				currentCh, w = tokenizer.peek()
			}

			tokenizer.TokenArr = append(
				tokenizer.TokenArr,
				parser.Token{TypeOfToken: parser.INTEGER_LITERAL, Literal: tokenizer.buf},
			)
			tokenizer.buf = ""
		} else if unicode.IsLetter(currentCh) || currentCh == '_' {
			for unicode.IsLetter(currentCh) || currentCh == '_' || unicode.IsDigit(currentCh) {
				tokenizer.buf += string(currentCh)
				tokenizer.readCh(w)
				currentCh, w = tokenizer.peek()
			}

			if tokenizer.buf == "true" {
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.KEYWORD_TRUE, Literal: tokenizer.buf})
			} else if tokenizer.buf == "false" {
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.KEYWORD_FALSE, Literal: tokenizer.buf})
			} else if tokenizer.buf == "if" {
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.IF, Literal: tokenizer.buf})
			} else if tokenizer.buf == "else" {
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.ELSE, Literal: tokenizer.buf})
			} else if tokenizer.buf == "let" {
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.LET})
			} else {
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.IDENTIFIER, Literal: tokenizer.buf})
			}
			tokenizer.buf = ""
		} else if currentCh == '=' {
			tokenizer.readCh(w)
			currentCh, w = tokenizer.peek()
			if currentCh == '=' {
				tokenizer.readCh(w)
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.EQUAL_TO, Literal: "=="})
			} else {
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.ASSIGNMENT_OPERATOR, Literal: "="})
			}
		} else if currentCh == '!' {
			tokenizer.readCh(w)
			currentCh, w = tokenizer.peek()
			if currentCh == '=' {
				tokenizer.readCh(w)
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.NOT_EQUAL_TO, Literal: "!="})
			} else {
				tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.BOOLEAN_INVERSION, Literal: "!"})
			}
		} else if currentCh == '>' {
			tokenizer.readCh(w)
			tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.GREATER_THAN, Literal: ">"})
		} else if currentCh == '+' {
			tokenizer.readCh(w)
			tokenizer.TokenArr = append(tokenizer.TokenArr,
				parser.Token{
					TypeOfToken: parser.PLUS,
					Literal:     "+",
				},
			)
		} else if currentCh == '-' {
			tokenizer.readCh(w)
			tokenizer.TokenArr = append(tokenizer.TokenArr,
				parser.Token{
					TypeOfToken: parser.MINUS,
					Literal:     "-",
				},
			)
		} else if currentCh == '*' {
			tokenizer.readCh(w)
			tokenizer.TokenArr = append(tokenizer.TokenArr,
				parser.Token{
					TypeOfToken: parser.ASTERISK,
					Literal:     "*",
				},
			)
		} else if currentCh == '/' {
			tokenizer.readCh(w)
			tokenizer.TokenArr = append(tokenizer.TokenArr,
				parser.Token{
					TypeOfToken: parser.SLASH,
					Literal:     "/",
				},
			)
		} else if currentCh == '(' {
			tokenizer.readCh(w)
			tokenizer.TokenArr = append(tokenizer.TokenArr,
				parser.Token{
					TypeOfToken: parser.OPENING_PARENTHESIS,
					Literal:     "(",
				},
			)
		} else if currentCh == ')' {
			tokenizer.readCh(w)
			tokenizer.TokenArr = append(tokenizer.TokenArr,
				parser.Token{
					TypeOfToken: parser.CLOSING_PARENTHESIS,
					Literal:     ")",
				},
			)
		} else if currentCh == '{' {
			tokenizer.readCh(w)
			tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.OPENING_CURLY, Literal: "{"})
		} else if currentCh == '}' {
			tokenizer.readCh(w)
			tokenizer.TokenArr = append(tokenizer.TokenArr, parser.Token{TypeOfToken: parser.CLOSING_CURLY, Literal: "}"})
		} else if currentCh == '\r' || currentCh == '\n' || currentCh == ' ' {
			tokenizer.readCh(w)
		} else {
			panic("Undefined symbol:" + string(currentCh))
		}
		currentCh, w = tokenizer.peek()
	}
}

func (tokenizer *Tokenizer) peek() (rune, int) {
	var ch rune
	var w int

	if len(tokenizer.src) > tokenizer.cursorPos {
		ch, w = utf8.DecodeRuneInString(tokenizer.src[tokenizer.cursorPos:])
	}

	return ch, w
}

func (tokenizer *Tokenizer) readCh(w int) {
	tokenizer.cursorPos += w
}
