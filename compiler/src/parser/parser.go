package parser

import (
	"fmt"
	"majorProject/compiler/objCodeGenerator"
	"os"
	"strconv"
)

type TokenType string

const (
	LOWEST = iota
	SUM
)

const (
	INTEGER_LITERAL     = "INT"
	STRING_LITERAL      = "STRING"
	IDENTIFIER          = "IDENT"
	ASSIGNMENT_OPERATOR = "="
	PLUS                = "+"
	SUBTRACTION         = "-"
	MULTIPLICATION      = "*"
	DIVISION            = "/"
	OPENING_PARENTHESIS = "("
	CLOSING_PARENTHESIS = ")"
	RETURN              = "RETURN"
)

var precedences = map[TokenType]int{
	INTEGER_LITERAL: LOWEST,
	PLUS:            SUM,
}

type (
	prefix_parse_fn func() objCodeGenerator.ExpressionNode
	infix_parse_fn  func(
		ast objCodeGenerator.ExpressionNode,
	) objCodeGenerator.ExpressionNode
)

type Token struct {
	TypeOfToken TokenType
	Literal     string
}

type Parser struct {
	Ast       objCodeGenerator.Program
	tokenArr  []Token
	cursorPos *int

	infix_parse_fns  map[TokenType]infix_parse_fn
	prefix_parse_fns map[TokenType]prefix_parse_fn
}

func New(tokenArr []Token) Parser {
	var a int = 0
	var p = Parser{
		Ast:       objCodeGenerator.Program{},
		tokenArr:  tokenArr,
		cursorPos: &a,
	}

	p.infix_parse_fns = make(map[TokenType]infix_parse_fn)
	p.infix_parse_fns[PLUS] = p.parse_infix_expression

	p.prefix_parse_fns = make(map[TokenType]prefix_parse_fn)
	p.prefix_parse_fns[INTEGER_LITERAL] = p.parse_int_literal

	return p
}

func (parser *Parser) Parse_program() int {
	var currentToken Token = parser.peek()
	var stmt_node *objCodeGenerator.StatementNode

	for currentToken != (Token{}) {
		stmt_node = parser.parseStatement()
		if stmt_node != nil {
			parser.Ast.Statements = append(
				parser.Ast.Statements,
				*stmt_node)
		}
		currentToken = parser.peek()
	}
	return 1
}

func (parser *Parser) peek() Token {
	var token Token
	if len(parser.tokenArr) > *parser.cursorPos {
		token = parser.tokenArr[*parser.cursorPos]
	}
	return token
}

func (parser *Parser) readToken() {
	*parser.cursorPos += 1
}

func (parser *Parser) parseStatement() *objCodeGenerator.StatementNode {
	var stmt_node objCodeGenerator.StatementNode
	var currentToken Token = parser.peek()

	switch currentToken.TypeOfToken {
	case RETURN:
	default:
		stmt_node = *parser.parseExpression(LOWEST)
	}

	return &stmt_node
}

func (parser *Parser) parseExpression(precedence int) *objCodeGenerator.ExpressionNode {
	var expr objCodeGenerator.ExpressionNode
	var current_token = parser.peek()
	var prefix = parser.prefix_parse_fns[current_token.TypeOfToken]

	if prefix == nil {
		fmt.Println("no prefix function found for", current_token.TypeOfToken)
		os.Exit(0)
		return &expr
	}
	expr = prefix()
	current_token = parser.peek()

	for (current_token.TypeOfToken == PLUS ||
		current_token.TypeOfToken == SUBTRACTION ||
		current_token.TypeOfToken == MULTIPLICATION ||
		current_token.TypeOfToken == DIVISION) &&
		precedences[current_token.TypeOfToken] > precedence &&
		*parser.cursorPos < len(parser.tokenArr) {

		var infix = parser.infix_parse_fns[current_token.TypeOfToken]
		if infix == nil {
			fmt.Println("no infix found for ", current_token.TypeOfToken)
			return &expr
		}

		expr = infix(expr)
	}

	parser.readToken()
	current_token = parser.peek()

	return &expr
}

func (parser *Parser) parse_int_literal() objCodeGenerator.ExpressionNode {
	var exp_node = objCodeGenerator.IntegerLiteralNode{}
	var current_token = parser.peek()

	exp_node.Value, _ = strconv.ParseInt(current_token.Literal, 10, 64)
	parser.readToken()
	return &exp_node
}

func (parser *Parser) parse_infix_expression(ast objCodeGenerator.ExpressionNode) objCodeGenerator.ExpressionNode {
	var exp_node = objCodeGenerator.InfixExpressionNode{}
	var current_token = parser.peek()

	if current_token.TypeOfToken == PLUS {
		exp_node.Op = objCodeGenerator.PLUS
	}

	parser.readToken()
	current_token = parser.peek()

	exp_node.Left = &ast
	exp_node.Right = parser.parseExpression(precedences[current_token.TypeOfToken])
	return &exp_node
}
