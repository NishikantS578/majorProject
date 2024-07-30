package main

import (
	"os"
)

type NodeProg struct {
	stmtNodes []NodeStmt
}

type NodeStmt struct {
	exitStmtNode *NodeExitStmt
	varDefNode   *NodeVarDefStmt
}

type NodeExitStmt struct {
	exprNode *NodeExpr
}

type NodeVarDefStmt struct {
	ident    *NodeIdent
	exprNode *NodeExpr
}

type NodeExpr struct {
	intLitNode *NodeIntLit
	identNode  *NodeIdent
}

type NodeIntLit struct {
	intLit Token
}

type NodeIdent struct {
	ident Token
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

func (parser *Parser) parseStmt() NodeStmt {
	var stmtNode NodeStmt
	var token = parser.peek()
	if false {
	} else if token.typeOfToken == EXIT_TOKEN {
		stmtNode = NodeStmt{exitStmtNode: &NodeExitStmt{}}
	}
	return stmtNode
}

func (parser *Parser) parseExitStmt() NodeExitStmt {
	var exitStmtNode NodeExitStmt
	var token = parser.consume()
	if token.typeOfToken == EXIT_TOKEN {
		token = parser.consume()
		if token.typeOfToken != OPENING_PARENTHESIS_TOKEN {
			println("Expected opening parenthesis: '('")
			os.Exit(0)
		}

		var exprNode = parser.parseExpr()

		token = parser.consume()
		if token.typeOfToken != CLOSING_PARENTHESIS_TOKEN {
			println("Expected closing parenthesis: ')'")
			os.Exit(0)
		}

		token = parser.consume()
		if token == (Token{}) || token.typeOfToken != SEMICOLON_TOKEN {
			println("Expected semicolon: ';'")
			os.Exit(0)
		}

		exitStmtNode = NodeExitStmt{exprNode: &exprNode}
	} else {
		println("Expected exit() statement")
		os.Exit(0)
	}

	return exitStmtNode
}

func (parser *Parser) parseExpr() NodeExpr {
	var exprNode NodeExpr
	var token = parser.consume()
	if token.typeOfToken == INTEGER_LITERAL_TOKEN {
		exprNode = NodeExpr{intLitNode: &NodeIntLit{intLit: token}}
	} else if token.typeOfToken == IDENTIFIER_TOKEN {
		exprNode = NodeExpr{identNode: &NodeIdent{token}}
	}

	return exprNode
}

func (parser *Parser) parse() NodeProg {
	var rootNode NodeProg = NodeProg{stmtNodes: []NodeStmt{}}

	var token Token = parser.peek()
	for token != (Token{}) {
		if token.typeOfToken == EXIT_TOKEN {
			var stmtNode NodeStmt
			parser.parseExitStmt()
			stmtNode = NodeStmt{exitStmtNode: &NodeExitStmt{}}
			rootNode.stmtNodes = append(rootNode.stmtNodes, stmtNode)
		} else if token.typeOfToken == LET_TOKEN {
		} else {
			parser.consume()
		}
		token = parser.peek()
	}

	parser.index = 0
	return rootNode
}
