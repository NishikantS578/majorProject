package parser

import (
	"errors"
	"fmt"
	"majorProject/compiler/objCodeGenerator"
	"os"
	"strconv"
)

type TokenType string

const (
	LOWEST = iota
	ASSIGNMENT
	SUM
	MULTIPLICATION
)

const (
	INTEGER_LITERAL     = "INT"
	STRING_LITERAL      = "STRING"
	IDENTIFIER          = "IDENT"
	ASSIGNMENT_OPERATOR = "="
	PLUS                = "+"
	MINUS               = "-"
	ASTERISK            = "*"
	SLASH               = "/"
	OPENING_PARENTHESIS = "("
	CLOSING_PARENTHESIS = ")"
	RETURN              = "RETURN"
)

var precedences = map[TokenType]int{
	IDENTIFIER:          LOWEST,
	ASSIGNMENT_OPERATOR: ASSIGNMENT,
	INTEGER_LITERAL:     LOWEST,
	PLUS:                SUM,
	MINUS:               SUM,
	ASTERISK:            MULTIPLICATION,
	SLASH:               MULTIPLICATION,
}

type (
	prefix_parse_fn func() (objCodeGenerator.ExpressionNode, error)
	infix_parse_fn  func(
		ast objCodeGenerator.ExpressionNode,
	) (objCodeGenerator.ExpressionNode, error)
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
	p.infix_parse_fns[MINUS] = p.parse_infix_expression
	p.infix_parse_fns[ASTERISK] = p.parse_infix_expression
	p.infix_parse_fns[SLASH] = p.parse_infix_expression
	p.infix_parse_fns[ASSIGNMENT_OPERATOR] = p.parse_infix_expression

	p.prefix_parse_fns = make(map[TokenType]prefix_parse_fn)
	p.prefix_parse_fns[INTEGER_LITERAL] = p.parse_int_literal
	p.prefix_parse_fns[OPENING_PARENTHESIS] = p.parse_grouped_expr
	p.prefix_parse_fns[IDENTIFIER] = p.parse_identifier

	return p
}

func (parser *Parser) Parse_program() int {
	var currentToken Token = parser.peek()
	var stmt_node objCodeGenerator.StatementNode
	var err error

	for currentToken != (Token{}) {
		stmt_node, err = parser.parseStatement()
		if err != nil {
			return 0
		}
		parser.Ast.Statements = append(
			parser.Ast.Statements,
			stmt_node)
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

func (parser *Parser) parseStatement() (objCodeGenerator.StatementNode, error) {
	var stmt_node objCodeGenerator.StatementNode
	var currentToken Token = parser.peek()
	var err error

	switch currentToken.TypeOfToken {
	case RETURN:
	default:
		stmt_node, err = parser.parseExpression(LOWEST)
		if err != nil {
			return stmt_node, err
		}
	}

	return stmt_node, nil
}

func (parser *Parser) parseExpression(precedence int) (objCodeGenerator.ExpressionNode, error) {
	var current_token = parser.peek()
	var expr objCodeGenerator.ExpressionNode
	var err error

	var prefix = parser.prefix_parse_fns[current_token.TypeOfToken]

	if prefix == nil {
		fmt.Println("no prefix function found for", current_token.TypeOfToken)
		os.Exit(0)
		return expr, nil
	}
	expr, err = prefix()
	if err != nil {
		return expr, err
	}
	current_token = parser.peek()

	if current_token.TypeOfToken == ASSIGNMENT_OPERATOR {
	}

	for (current_token.TypeOfToken == PLUS ||
		current_token.TypeOfToken == MINUS ||
		current_token.TypeOfToken == ASTERISK ||
		current_token.TypeOfToken == SLASH) &&
		precedences[current_token.TypeOfToken] > precedence &&
		*parser.cursorPos < len(parser.tokenArr) {

		var infix = parser.infix_parse_fns[current_token.TypeOfToken]
		if infix == nil {
			fmt.Println("no infix found for ", current_token.TypeOfToken)
			return expr, errors.New("err")
		}

		expr, err = infix(expr)
		if err != nil {
			return expr, err
		}
		current_token = parser.peek()
	}

	return expr, nil
}

func (parser *Parser) parse_grouped_expr() (objCodeGenerator.ExpressionNode, error) {
	var exp_node objCodeGenerator.ExpressionNode
	var err error
	var current_token = parser.peek()

	parser.readToken()
	exp_node, err = parser.parseExpression(LOWEST)
	if err != nil {
		return exp_node, err
	}

	current_token = parser.peek()
	if current_token.TypeOfToken != CLOSING_PARENTHESIS {
		println("Expected ')'")
		return exp_node, errors.New("syntax error")
	}
	parser.readToken()

	return exp_node, nil
}

func (parser *Parser) parse_identifier() (objCodeGenerator.ExpressionNode, error) {
	var exp_node = objCodeGenerator.IdentifierNode{}
	var err error
	var current_token = parser.peek()
	exp_node.Value = current_token.Literal
	parser.readToken()

	return &exp_node, err
}

func (parser *Parser) parse_int_literal() (objCodeGenerator.ExpressionNode, error) {
	var exp_node = objCodeGenerator.IntegerLiteralNode{}
	var current_token = parser.peek()

	exp_node.Value, _ = strconv.ParseInt(current_token.Literal, 10, 64)
	parser.readToken()
	return &exp_node, nil
}

func (parser *Parser) parse_infix_expression(ast objCodeGenerator.ExpressionNode) (objCodeGenerator.ExpressionNode, error) {
	var exp_node = objCodeGenerator.InfixExpressionNode{}
	var current_token = parser.peek()
	var err error
	var current_precedence = precedences[current_token.TypeOfToken]

	if current_token.TypeOfToken == PLUS {
		exp_node.Op = objCodeGenerator.PLUS
	} else if current_token.TypeOfToken == MINUS {
		exp_node.Op = objCodeGenerator.MINUS
	} else if current_token.TypeOfToken == ASTERISK {
		exp_node.Op = objCodeGenerator.MULTIPLICATION
	} else if current_token.TypeOfToken == SLASH {
		exp_node.Op = objCodeGenerator.DIVISION
	} else if current_token.TypeOfToken == ASSIGNMENT_OPERATOR {
		exp_node.Op = objCodeGenerator.ASSIGNMENT
	} else {
		println("expected operator")
	}

	parser.readToken()
	current_token = parser.peek()

	exp_node.Left = ast
	exp_node.Right, err = parser.parseExpression(current_precedence)
	if err != nil {
		return &exp_node, err
	}
	return &exp_node, nil
}
