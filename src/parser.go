package main

import (
	"os"
)

type NodeExpr struct {
	intLit Token
}

type NodeExit struct {
	exprNode NodeExpr
}

type Parser struct {
	tokenArr []Token
	index    int
}

func (parser *Parser) initialize(tokenArr []Token) {
	parser.tokenArr = tokenArr
	parser.index = 0
}

func (parser *Parser) peek() Token {
	var token Token
	if parser.index+1 > len(parser.tokenArr) {
		return token
	} else {
		token = parser.tokenArr[parser.index]
		return token
	}
}

func (parser *Parser) consume() Token {
	var token Token = parser.peek()
	parser.index += 1
	return token
}

func (parser *Parser) parseExpr() NodeExpr {
	var exprNode NodeExpr
	var token = parser.consume()
	if token.typeOfToken == INTEGER_LITERAL_TOKEN {
		exprNode = NodeExpr{intLit: token}
	}
	return exprNode
}

func (parser *Parser) parse() NodeExit {
	var rootNode NodeExit

	var token Token = parser.consume()
	for token != (Token{}) {
		if token.typeOfToken == EXIT_TOKEN {
			var exprNode = parser.parseExpr()
			if exprNode != (NodeExpr{}) {
				rootNode = NodeExit{exprNode: exprNode}
			} else {
				println("Invalid Expression")
				os.Exit(0)
			}
			token = parser.consume()

			if token == (Token{}) || token.typeOfToken != SEMICOLON_TOKEN {
				println("Expected semicolon: ';'")
				os.Exit(0)
			}
			token = parser.consume()
		}
	}

	parser.index = 0
	return rootNode
}
