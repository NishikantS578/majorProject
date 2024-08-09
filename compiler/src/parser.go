package main

import (
	"os"
)

type NodeProg struct {
	stmtNodes []NodeStmt
}

type NodeStmt struct {
	exitStmtNode   *NodeExitStmt
	varDefStmtNode *NodeVarDefStmt
}

type NodeExitStmt struct {
	exprNode *NodeExpr
	termNode *NodeTerm
}

type NodeVarDefStmt struct {
	ident    *NodeIdent
	exprNode *NodeExpr
	termNode *NodeTerm
}

type NodeTerm struct {
	intLitNode *NodeIntLit
	identNode  *NodeIdent
}

type NodeAddExpr struct {
	lExprNode *NodeExpr
	lTermNode *NodeTerm
	rExprNode *NodeExpr
	rTermNode *NodeTerm
}

type NodeSubExpr struct {
	lExprNode *NodeExpr
	lTermNode *NodeTerm
	rExprNode *NodeExpr
	rTermNode *NodeTerm
}

type NodeMulExpr struct {
	lExprNode *NodeExpr
	lTermNode *NodeTerm
	rExprNode *NodeExpr
	rTermNode *NodeTerm
}

type NodeDivExpr struct {
	lExprNode *NodeExpr
	lTermNode *NodeTerm
	rExprNode *NodeExpr
	rTermNode *NodeTerm
}

type NodeExpr struct {
	addExprNode *NodeAddExpr
	subExprNode *NodeSubExpr
	mulExprNode *NodeMulExpr
	divExprNode *NodeDivExpr
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

func (parser *Parser) peek(ahead int) Token {
	var token Token
	if parser.index+ahead >= len(parser.tokenArr) {
		return token
	} else {
		token = parser.tokenArr[parser.index+ahead]
		return token
	}
}

func (parser *Parser) consume() Token {
	var token Token = parser.peek(0)
	parser.index += 1
	return token
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
		if exprNode == (NodeExpr{}) {
			var termNode = parser.parseTerm()
			exitStmtNode.termNode = &termNode
		} else {
			exitStmtNode.exprNode = &exprNode
		}

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
	} else {
		println("Expected exit() statement")
		os.Exit(0)
	}

	return exitStmtNode
}

func (parser *Parser) parseVarDefStmt() NodeVarDefStmt {
	var varDefStmtNode NodeVarDefStmt = NodeVarDefStmt{}
	parser.consume()
	var token Token = parser.consume()

	var identNode NodeIdent
	if token.typeOfToken == IDENTIFIER_TOKEN {
		identNode.ident = token
	} else {
		println("Unexpected Token: " + token.value + " // Expected Identifier")
		os.Exit(0)
	}
	varDefStmtNode.ident = &identNode

	token = parser.consume()
	if token.typeOfToken != ASSIGNMENT_OPERATOR_TOKEN {
		println("Unexpected Token: " + token.value + " // Expected '='")
		os.Exit(0)
	}

	var exprNode NodeExpr = parser.parseExpr()
	if exprNode == (NodeExpr{}) {
		var termNode NodeTerm = parser.parseTerm()
		varDefStmtNode.termNode = &termNode
	} else {
		varDefStmtNode.exprNode = &exprNode
	}

	token = parser.consume()
	if token.typeOfToken != SEMICOLON_TOKEN {
		println("Unexpected Token: " + token.value + " // Expected ';'")
		os.Exit(0)
	}

	return varDefStmtNode
}

func (parser *Parser) parseTerm() NodeTerm {
	var termNode NodeTerm

	var token = parser.peek(0)
	if token.typeOfToken == INTEGER_LITERAL_TOKEN {
		termNode.intLitNode = &NodeIntLit{intLit: token}
	} else if token.typeOfToken == IDENTIFIER_TOKEN {
		termNode.identNode = &NodeIdent{token}
	}
	parser.consume()

	return termNode
}

func (parser *Parser) parseAddExpr() NodeAddExpr {
	var addExprNode NodeAddExpr

	var lExprNode NodeExpr
	var lTerm NodeTerm
	var rExprNode NodeExpr
	var rTerm NodeTerm

	lTerm = parser.parseTerm()
	if lTerm == (NodeTerm{}) {
		lExprNode = parser.parseExpr()
		addExprNode.lExprNode = &lExprNode
	} else {
		addExprNode.lTermNode = &lTerm
	}

	parser.consume()

	rExprNode = parser.parseExpr()
	if rExprNode == (NodeExpr{}) {
		rTerm = parser.parseTerm()
		addExprNode.rTermNode = &rTerm
	} else {
		addExprNode.rExprNode = &rExprNode
	}

	return addExprNode
}

func (parser *Parser) parseExpr() NodeExpr {
	var exprNode NodeExpr

	var token = parser.peek(1)
	if token.typeOfToken == PLUS_TOKEN {
		var NodeAddExpr = parser.parseAddExpr()
		exprNode.addExprNode = &NodeAddExpr
	} else if token.typeOfToken == SUBTRACTION_TOKEN {
	} else if token.typeOfToken == MULTIPLICATION_TOKEN {
	} else if token.typeOfToken == DIVISION_TOKEN {
	}

	return exprNode
}

func (parser *Parser) parse() NodeProg {
	var rootNode NodeProg = NodeProg{stmtNodes: []NodeStmt{}}

	var token Token = parser.peek(0)
	var stmtNode NodeStmt
	for token != (Token{}) {
		if token.typeOfToken == EXIT_TOKEN {
			var exitStmtNode = parser.parseExitStmt()
			stmtNode = NodeStmt{exitStmtNode: &exitStmtNode}
			rootNode.stmtNodes = append(rootNode.stmtNodes, stmtNode)
		} else if token.typeOfToken == LET_TOKEN {
			var varDefStmtNode = parser.parseVarDefStmt()
			stmtNode = NodeStmt{varDefStmtNode: &varDefStmtNode}
			rootNode.stmtNodes = append(rootNode.stmtNodes, stmtNode)
		} else {
			println("Unexpected Token: " + token.value)
			os.Exit(0)
		}
		token = parser.peek(0)
	}

	parser.index = 0
	return rootNode
}
