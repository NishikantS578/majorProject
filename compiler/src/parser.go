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
	exprNode   *NodeExpr
}

type NodeFactor struct {
	op        rune
	lTermNode *NodeTerm
	rTermNode *NodeTerm
}

type NodeExpr struct {
	op          rune
	lFactorNode *NodeFactor
	rFactorNode *NodeFactor
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
		println("Unexpected Token: '" + token.value + "' // Expected ';'")
		os.Exit(0)
	}

	return varDefStmtNode
}

func (parser *Parser) parseTerm() NodeTerm {
	var termNode NodeTerm

	var token = parser.peek(0)
	if token.typeOfToken == INTEGER_LITERAL_TOKEN {
		parser.consume()
		termNode.intLitNode = &NodeIntLit{intLit: token}
	} else if token.typeOfToken == IDENTIFIER_TOKEN {
		parser.consume()
		termNode.identNode = &NodeIdent{token}
	} else if token.typeOfToken == OPENING_PARENTHESIS_TOKEN {
		parser.consume()
		var exprNode = parser.parseExpr()
		termNode.exprNode = &exprNode
		if token.typeOfToken != CLOSING_PARENTHESIS_TOKEN {
			println("Expected ')'")
		}
		parser.consume()
	}

	return termNode
}

func (parser *Parser) parseFactor() NodeFactor {
	var factorNode *NodeFactor = &NodeFactor{}
	var tFactorNode *NodeFactor

	var lTermNode = parser.parseTerm()
	factorNode.lTermNode = &lTermNode
	var rTermNode NodeTerm

	var token = parser.peek(0)
	for token.typeOfToken == MULTIPLICATION_TOKEN || token.typeOfToken == DIVISION_TOKEN {
		var op rune
		if token.typeOfToken == MULTIPLICATION_TOKEN {
			op = '*'
		} else if token.typeOfToken == DIVISION_TOKEN {
			op = '/'
		}

		parser.consume()
		rTermNode = parser.parseTerm()
		factorNode.rTermNode = &rTermNode
		factorNode.op = op

		tFactorNode = factorNode
		factorNode = &NodeFactor{}
		factorNode.lTermNode = &NodeTerm{}
		factorNode.lTermNode.exprNode = &NodeExpr{}
		factorNode.lTermNode.exprNode.lFactorNode = tFactorNode

		token = parser.peek(0)
	}

	return *factorNode
}

func (parser *Parser) parseExpr() NodeExpr {
	var exprNode *NodeExpr = &NodeExpr{}
	var tExprNode *NodeExpr

	var lFactorNode = parser.parseFactor()
	exprNode.lFactorNode = &lFactorNode

	var rFactorNode NodeFactor

	var token = parser.peek(0)

	for token.typeOfToken == PLUS_TOKEN || token.typeOfToken == SUBTRACTION_TOKEN {
		var op rune
		if token.typeOfToken == PLUS_TOKEN {
			op = '+'
		} else if token.typeOfToken == SUBTRACTION_TOKEN {
			op = '-'
		}

		parser.consume()
		rFactorNode = parser.parseFactor()
		exprNode.rFactorNode = &NodeFactor{}
		*exprNode.rFactorNode = rFactorNode
		exprNode.op = op

		tExprNode = exprNode
		exprNode = &NodeExpr{}
		exprNode.lFactorNode = &NodeFactor{}
		exprNode.lFactorNode.lTermNode = &NodeTerm{}
		exprNode.lFactorNode.lTermNode.exprNode = tExprNode

		token = parser.peek(0)
	}
	return *exprNode
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
