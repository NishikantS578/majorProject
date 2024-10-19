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
	COMPARISON
	ASSIGNMENT
	SUM
	MULTIPLICATION
	INVERSION
)

const (
	INTEGER_LITERAL     = "INT"
	STRING_LITERAL      = "STRING"
	IDENTIFIER          = "IDENT"
	KEYWORD_TRUE        = "KEYWORD true"
	KEYWORD_FALSE       = "KEYWORD false"
	ASSIGNMENT_OPERATOR = "="
	PLUS                = "+"
	MINUS               = "-"
	ASTERISK            = "*"
	SLASH               = "/"
	OPENING_PARENTHESIS = "("
	CLOSING_PARENTHESIS = ")"
	OPENING_CURLY       = "{"
	CLOSING_CURLY       = "}"
	EQUAL_TO            = "=="
	NOT_EQUAL_TO        = "!="
	BOOLEAN_INVERSION   = "!"
	GREATER_THAN        = ">"
	IF                  = "if"
	ELSE                = "else"
	RETURN              = "RETURN"
)

var precedences = map[TokenType]int{
	IDENTIFIER:          LOWEST,
	ASSIGNMENT_OPERATOR: ASSIGNMENT,
	EQUAL_TO:            COMPARISON,
	NOT_EQUAL_TO:        COMPARISON,
	GREATER_THAN:        COMPARISON,
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
	Ast       objCodeGenerator.StatementBlockNode
	tokenArr  []Token
	cursorPos *int

	infix_parse_fns  map[TokenType]infix_parse_fn
	prefix_parse_fns map[TokenType]prefix_parse_fn
}

func New(tokenArr []Token) Parser {
	var a int = 0
	var p = Parser{
		Ast:       objCodeGenerator.StatementBlockNode{},
		tokenArr:  tokenArr,
		cursorPos: &a,
	}

	p.infix_parse_fns = make(map[TokenType]infix_parse_fn)
	p.infix_parse_fns[PLUS] = p.parse_infix_expression
	p.infix_parse_fns[MINUS] = p.parse_infix_expression
	p.infix_parse_fns[ASTERISK] = p.parse_infix_expression
	p.infix_parse_fns[SLASH] = p.parse_infix_expression
	p.infix_parse_fns[ASSIGNMENT_OPERATOR] = p.parse_infix_expression
	p.infix_parse_fns[EQUAL_TO] = p.parse_infix_expression
	p.infix_parse_fns[NOT_EQUAL_TO] = p.parse_infix_expression
	p.infix_parse_fns[GREATER_THAN] = p.parse_infix_expression

	p.prefix_parse_fns = make(map[TokenType]prefix_parse_fn)
	p.prefix_parse_fns[INTEGER_LITERAL] = p.parse_int_literal
	p.prefix_parse_fns[KEYWORD_TRUE] = p.parse_bool_keyword
	p.prefix_parse_fns[KEYWORD_FALSE] = p.parse_bool_keyword
	p.prefix_parse_fns[OPENING_PARENTHESIS] = p.parse_grouped_expr
	p.prefix_parse_fns[IDENTIFIER] = p.parse_identifier
	p.prefix_parse_fns[BOOLEAN_INVERSION] = p.parse_prefix_expression
	p.prefix_parse_fns[MINUS] = p.parse_prefix_expression

	return p
}

func (parser *Parser) ParseProgram() int {
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

func (parser *Parser) parse_stmt_block_node() (objCodeGenerator.StatementBlockNode, error) {
	var currentToken Token = parser.peek()
	var stmt_node objCodeGenerator.StatementNode
	var stmt_block_node objCodeGenerator.StatementBlockNode
	var err error

	for currentToken.TypeOfToken != CLOSING_CURLY {
		stmt_node, err = parser.parseStatement()
		if err != nil {
			return stmt_block_node, err
		}

		stmt_block_node.Statements = append(
			stmt_block_node.Statements,
			stmt_node)
		currentToken = parser.peek()
	}
	return stmt_block_node, err
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
	case IF:
		stmt_node, err = parser.parse_if_statement()
		if err != nil {
			return stmt_node, err
		}
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

	// if current_token.TypeOfToken == ASSIGNMENT_OPERATOR {
	// }

	for (current_token.TypeOfToken == PLUS ||
		current_token.TypeOfToken == MINUS ||
		current_token.TypeOfToken == ASTERISK ||
		current_token.TypeOfToken == SLASH ||
		current_token.TypeOfToken == EQUAL_TO ||
		current_token.TypeOfToken == NOT_EQUAL_TO ||
		current_token.TypeOfToken == GREATER_THAN) &&
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

func (parser *Parser) parse_bool_keyword() (objCodeGenerator.ExpressionNode, error) {
	var exp_node = objCodeGenerator.KeywordBooleanNode{}
	var current_token = parser.peek()

	if current_token.TypeOfToken == KEYWORD_TRUE {
		exp_node.Value = true
	} else if current_token.TypeOfToken == KEYWORD_FALSE {
		exp_node.Value = false
	}
	parser.readToken()
	return &exp_node, nil
}

func (parser *Parser) parse_if_statement() (*objCodeGenerator.IfStmtNode, error) {
	var ifstmt_node = objCodeGenerator.IfStmtNode{}
	var err error
	parser.readToken()

	var current_token = parser.peek()
	if current_token.TypeOfToken != OPENING_PARENTHESIS {
		return &ifstmt_node, errors.New("expected (")
	}
	parser.readToken()

	var exp_node objCodeGenerator.ExpressionNode
	exp_node, err = parser.parseExpression(0)
	if err != nil {
		return &ifstmt_node, err
	}
	current_token = parser.peek()
	if current_token.TypeOfToken != CLOSING_PARENTHESIS {
		return &ifstmt_node, errors.New("expected )")
	}
	parser.readToken()
	current_token = parser.peek()
	if current_token.TypeOfToken != OPENING_CURLY {
		return &ifstmt_node, errors.New("expected {")
	}
	parser.readToken()

	var stmt_block_node objCodeGenerator.StatementBlockNode
	stmt_block_node, err = parser.parse_stmt_block_node()
	if err != nil {
		return &ifstmt_node, err
	}

	current_token = parser.peek()
	if current_token.TypeOfToken != CLOSING_CURLY {
		return &ifstmt_node, errors.New("expected }")
	}
	parser.readToken()

	ifstmt_node.Condition = exp_node
	ifstmt_node.Consequence = stmt_block_node

	current_token = parser.peek()
	if current_token.TypeOfToken != ELSE {
		return &ifstmt_node, err
	}

	parser.readToken()
	current_token = parser.peek()
	if current_token.TypeOfToken != OPENING_CURLY{
		return &ifstmt_node, errors.New("expected {")
	}

	parser.readToken()
	stmt_block_node, err = parser.parse_stmt_block_node()
	if err != nil{
		return &ifstmt_node, errors.New("error while parsing statement block")
	}

	current_token = parser.peek()
	if current_token.TypeOfToken != CLOSING_CURLY{
		return &ifstmt_node, errors.New("expected }")
	}

	parser.readToken()

	ifstmt_node.Alternative = stmt_block_node
	return &ifstmt_node, err
}

func (parser *Parser) parse_prefix_expression() (objCodeGenerator.ExpressionNode, error) {
	var exp_node = objCodeGenerator.PrefixExpressionNode{}
	var current_token = parser.peek()
	var err error

	parser.readToken()
	switch current_token.TypeOfToken {
	case MINUS:
		exp_node.Op = objCodeGenerator.MINUS
		exp_node.Child, err = parser.parseExpression(LOWEST)
	case BOOLEAN_INVERSION:
		exp_node.Op = objCodeGenerator.BOOLEAN_INVERSION
		exp_node.Child, err = parser.parseExpression(LOWEST)
	default:
		return &exp_node, errors.New("syntax error")
	}

	return &exp_node, err
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
	} else if current_token.TypeOfToken == EQUAL_TO {
		exp_node.Op = objCodeGenerator.EQUAL_TO
	} else if current_token.TypeOfToken == NOT_EQUAL_TO {
		exp_node.Op = objCodeGenerator.NOT_EQUAL_TO
	} else if current_token.TypeOfToken == GREATER_THAN {
		exp_node.Op = objCodeGenerator.GREATER_THAN
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
