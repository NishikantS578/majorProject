class Tokenizer {
    TokenArr: Token[]
    src: string
    cursorPos: number
    buf: string

    constructor(line: string) {
        this.TokenArr = []
        this.src = line
        this.cursorPos = 0
        this.buf = ""
    }

    Tokenize() {
        let currentCh = this.peek()

        while (currentCh != "") {
            if (!Number.isNaN(parseInt(currentCh))) {
                this.buf += currentCh
                this.readCh()
                currentCh = this.peek()

                while (!Number.isNaN(parseInt(currentCh))) {
                    this.buf += currentCh
                    this.readCh()
                    currentCh = this.peek()
                }

                this.TokenArr.push(new Token(INTEGER_LITERAL, this.buf))
                this.buf = ""
            } else if (/^[a-zA-Z]$/.test(currentCh) || currentCh == "_") {
                while (/^[a-zA-Z]$/.test(currentCh) || currentCh == '_' || /^[a-zA-Z]$/.test(currentCh)) {
                    this.buf += currentCh
                    this.readCh()
                    currentCh = this.peek()
                }

                if (this.buf == "true") {
                    this.TokenArr.push(new Token(KEYWORD_TRUE, this.buf))
                } else if (this.buf == "false") {
                    this.TokenArr.push(new Token(KEYWORD_FALSE, this.buf))
                } else if (this.buf == "if") {
                    this.TokenArr.push(new Token(IF, this.buf))
                } else if (this.buf == "else") {
                    this.TokenArr.push(new Token(ELSE, this.buf))
                } else if (this.buf == "let") {
                    this.TokenArr.push(new Token(LET))
                } else {
                    this.TokenArr.push(new Token(IDENTIFIER, this.buf))
                }
                this.buf = ""
            } else if (currentCh == '=') {
                this.readCh()
                currentCh = this.peek()
                if (currentCh == '=') {
                    this.readCh()
                    this.TokenArr.push(new Token(EQUAL_TO, "=="))
                } else {
                    this.TokenArr.push(new Token(ASSIGNMENT_OPERATOR, "="))
                }
            } else if (currentCh == '!') {
                this.readCh()
                currentCh = this.peek()
                if (currentCh == '=') {
                    this.readCh()
                    this.TokenArr.push(new Token(NOT_EQUAL_TO, "!="))
                } else {
                    this.TokenArr.push(new Token(BOOLEAN_INVERSION, "!"))
                }
            } else if (currentCh == '>') {
                this.readCh()
                this.TokenArr.push(new Token(GREATER_THAN, ">"))
            } else if (currentCh == '+') {
                this.readCh()
                this.TokenArr.push(new Token(PLUS, "+",))
            } else if (currentCh == '-') {
                this.readCh()
                this.TokenArr.push(new Token(MINUS, "-",),)
            } else if (currentCh == '*') {
                this.readCh()
                this.TokenArr.push(new Token(ASTERISK, "*",),)
            } else if (currentCh == '/') {
                this.readCh()
                this.TokenArr.push(new Token(SLASH, "/",),)
            } else if (currentCh == '(') {
                this.readCh()
                this.TokenArr.push(new Token(OPENING_PARENTHESIS, "(",),)
            } else if (currentCh == ')') {
                this.readCh()
                this.TokenArr.push(new Token(CLOSING_PARENTHESIS, ")",),)
            } else if (currentCh == '{') {
                this.readCh()
                this.TokenArr.push(new Token(OPENING_CURLY, "{"))
            } else if (currentCh == '}') {
                this.readCh()
                this.TokenArr.push(new Token(CLOSING_CURLY, ")"))
            } else if (currentCh == '\r' || currentCh == '\n' || currentCh == ' ') {
                this.readCh()
            } else {
                console.log("Undefined symbol:" + currentCh)
            }
            currentCh = this.peek()
        }
    }

    peek(): string {
        let ch: string = ""

        if (this.src.length > this.cursorPos) {
            ch = this.src[this.cursorPos]
        }

        return ch
    }

    readCh() {
        this.cursorPos += 1
    }

}

type TokenType = string

const LOWEST = 0
const COMPARISON = 1
const ASSIGNMENT = 2
const SUM = 3
const MULTIPLICATION = 4
const INVERSION = 5

const INTEGER_LITERAL = "INT"
const STRING_LITERAL = "STRING"
const IDENTIFIER = "IDENT"
const KEYWORD_TRUE = "KEYWORD true"
const KEYWORD_FALSE = "KEYWORD false"
const ASSIGNMENT_OPERATOR = "="
const PLUS = "+"
const MINUS = "-"
const ASTERISK = "*"
const SLASH = "/"
const OPENING_PARENTHESIS = "("
const CLOSING_PARENTHESIS = ")"
const OPENING_CURLY = "{"
const CLOSING_CURLY = "}"
const EQUAL_TO = "=="
const NOT_EQUAL_TO = "!="
const BOOLEAN_INVERSION = "!"
const GREATER_THAN = ">"
const IF = "if"
const ELSE = "else"
const RETURN = "RETURN"
const LET = "LET KEYWORD"

const precedences = {
    IDENTIFIER: LOWEST,
    ASSIGNMENT_OPERATOR: ASSIGNMENT,
    EQUAL_TO: COMPARISON,
    NOT_EQUAL_TO: COMPARISON,
    GREATER_THAN: COMPARISON,
    INTEGER_LITERAL: LOWEST,
    PLUS: SUM,
    MINUS: SUM,
    ASTERISK: MULTIPLICATION,
    SLASH: MULTIPLICATION,
}

type prefix_parse_fn = () => ExpressionNode
type infix_parse_fn = (node: ExpressionNode) => ExpressionNode

class Token {
    TypeOfToken: TokenType
    Literal: string

    constructor(type_of_token: TokenType, literal: string = "") {
        this.TypeOfToken = type_of_token
        this.Literal = literal
    }
}

class Parser {
    Ast: StatementBlockNode
    tokenArr: [Token]
    cursorPos: number

    infix_parse_fns: Map<TokenType, infix_parse_fn>
    prefix_parse_fns: Map<TokenType, prefix_parse_fn>

    constructor(token_arr: [Token]) {
        this.Ast = new StatementBlockNode()
        this.tokenArr = token_arr
        this.cursorPos = 0

        this.infix_parse_fns = new Map<TokenType, infix_parse_fn>()
        this.infix_parse_fns.set(PLUS, this.parse_infix_expression)
        this.infix_parse_fns.set(MINUS, this.parse_infix_expression)
        this.infix_parse_fns.set(ASTERISK, this.parse_infix_expression)
        this.infix_parse_fns.set(SLASH, this.parse_infix_expression)
        this.infix_parse_fns.set(EQUAL_TO, this.parse_infix_expression)
        this.infix_parse_fns.set(NOT_EQUAL_TO, this.parse_infix_expression)
        this.infix_parse_fns.set(GREATER_THAN, this.parse_infix_expression)

        this.prefix_parse_fns = new Map<TokenType, prefix_parse_fn>()
        this.prefix_parse_fns.set(INTEGER_LITERAL, this.parse_int_literal)
        this.prefix_parse_fns.set(KEYWORD_TRUE, this.parse_bool_keyword)
        this.prefix_parse_fns.set(KEYWORD_FALSE, this.parse_bool_keyword)
        this.prefix_parse_fns.set(OPENING_PARENTHESIS, this.parse_grouped_expr)
        this.prefix_parse_fns.set(IDENTIFIER, this.parse_identifier)
        this.prefix_parse_fns.set(BOOLEAN_INVERSION, this.parse_prefix_expression)
        this.prefix_parse_fns.set(MINUS, this.parse_prefix_expression)
    }

    ParseProgram(): number {
        let currentToken: Token | null = this.peek()
        let stmt_node: StatementNode

        while (currentToken != null) {
            stmt_node = this.parseStatement()

            this.Ast.Statements.push(stmt_node)
            currentToken = this.peek()
        }
        return 1
    }

    peek(offset: number = 0): Token | null {
        let token: Token | null = null
        if (this.tokenArr.length > this.cursorPos + offset) {
            token = this.tokenArr[this.cursorPos + offset]
        }
        return token
    }

    readToken() {
        this.cursorPos += 1
    }

    parse_stmt_block_node(): StatementBlockNode {
        let currentToken: Token | null = this.peek()
        let stmt_node: StatementNode
        let stmt_block_node: StatementBlockNode = new StatementBlockNode()

        while (currentToken != null && currentToken.TypeOfToken != CLOSING_CURLY) {
            stmt_node = this.parseStatement()

            stmt_block_node.Statements.push(stmt_node)
            currentToken = this.peek()
        }
        return stmt_block_node
    }

    parseStatement(): StatementNode | null {
        let stmt_node: StatementNode | null = null
        let currentToken: Token | null = this.peek()

        if (currentToken == null) {
            return stmt_node
        }

        switch (currentToken.TypeOfToken) {
            case RETURN:
                break;
            case IF:
                stmt_node = this.parse_if_statement()
                break;
            case LET:
                stmt_node = this.parse_let_statement()
                break;
            default:
                let next_token = this.peek(1)
                if (next_token != null && next_token.TypeOfToken == ASSIGNMENT_OPERATOR) {
                    stmt_node = this.parse_assignment_statement()
                } else {
                    stmt_node = this.parseExpression(LOWEST)
                }
        }

        return stmt_node
    }

    parseExpression(precedence: number): ExpressionNode | null {
        let current_token = this.peek()
        let expr: ExpressionNode | null = null

        if (current_token == null) {
            return expr
        }

        let prefix = this.prefix_parse_fns.get(current_token.TypeOfToken)

        if (prefix == undefined) {
            console.log("no prefix function found for", current_token.TypeOfToken)
            return expr
        }

        expr = prefix()
        current_token = this.peek()

        if (current_token == null) {
            return expr
        }

        while ((current_token.TypeOfToken == PLUS ||
            current_token.TypeOfToken == MINUS ||
            current_token.TypeOfToken == ASTERISK ||
            current_token.TypeOfToken == SLASH ||
            current_token.TypeOfToken == EQUAL_TO ||
            current_token.TypeOfToken == NOT_EQUAL_TO ||
            current_token.TypeOfToken == GREATER_THAN) &&
            (precedences[current_token.TypeOfToken] > precedence) &&
            (this.cursorPos < len(this.tokenArr))) {

            infix = this.infix_parse_fns[current_token.TypeOfToken]
            if (infix == null) {
                fmt.Println("no infix found for ", current_token.TypeOfToken)
                return expr
            }

            expr = infix(expr)
            current_token = this.peek()
            if (current_token == null) {
                return expr
            }

        }

        return expr
    }

    parse_grouped_expr(): ExpressionNode {
        exp_node: ExpressionNode
        current_token = this.peek()

        this.readToken()
        exp_node = this.parseExpression(LOWEST)

        current_token = this.peek()
        if (current_token.TypeOfToken != CLOSING_PARENTHESIS) {
            console.log("Expected ')'")
            return exp_node
        }
        this.readToken()

        return exp_node
    }

    parse_identifier(): ExpressionNode {
        exp_node = new IdentifierNode()
        current_token = this.peek()
        exp_node.Value = current_token.Literal
        this.readToken()

        return exp_node
    }

    parse_int_literal(): ExpressionNode {
        exp_node = new IntegerLiteralNode()
        current_token = this.peek()

        exp_node.Value, _ = strconv.ParseInt(current_token.Literal, 10, 64)
        this.readToken()
        return exp_node
    }

    parse_bool_keyword(): ExpressionNode {
        exp_node = new KeywordBooleanNode()
        current_token = this.peek()

        if (current_token.TypeOfToken == KEYWORD_TRUE) {
            exp_node.Value = true
        } else if (current_token.TypeOfToken == KEYWORD_FALSE) {
            exp_node.Value = false
        }
        this.readToken()
        return exp_node
    }

    parse_let_statement(): LetStmtNode {
        let_stmt_node = new LetStmtNode()
        current_token = this.peek()

        if (current_token.TypeOfToken != LET) {
            return let_stmt_node
        }
        this.readToken()
        current_token = this.peek()

        let_stmt_node.Identifier = new IdentifierNode()
        let_stmt_node.Identifier.Value = current_token.Literal
        this.readToken()
        current_token = this.peek()

        if (current_token.TypeOfToken != ASSIGNMENT_OPERATOR) {
            return let_stmt_node
        }
        this.readToken()

        let_stmt_node.InitializationExpr = this.parseExpression(LOWEST)

        return let_stmt_node
    }

    parse_assignment_statement(): AssignementStmt {
        assignment_stmt_node = new AssignementStmt()
        current_token = this.peek()

        assignment_stmt_node.Identifier = new IdentifierNode()
        assignment_stmt_node.Identifier.Value = current_token.Literal
        this.readToken()
        current_token = this.peek()

        if (current_token.TypeOfToken != ASSIGNMENT_OPERATOR) {
            return assignment_stmt_node
        }
        this.readToken()

        assignment_stmt_node.InitializationExpr = this.parseExpression(LOWEST)

        return assignment_stmt_node
    }

    parse_if_statement(): IfStmtNode {
        ifstmt_node = new IfStmtNode()
        this.readToken()

        current_token = this.peek()
        if (current_token.TypeOfToken != OPENING_PARENTHESIS) {
            return ifstmt_node
        }
        this.readToken()

        exp_node: ExpressionNode
        exp_node = this.parseExpression(0)

        current_token = this.peek()
        if (current_token.TypeOfToken != CLOSING_PARENTHESIS) {
            return ifstmt_node
        }
        this.readToken()
        current_token = this.peek()
        if (current_token.TypeOfToken != OPENING_CURLY) {
            return ifstmt_node
        }
        this.readToken()

        stmt_block_node: StatementBlockNode
        stmt_block_node = this.parse_stmt_block_node()

        current_token = this.peek()
        if (current_token.TypeOfToken != CLOSING_CURLY) {
            return ifstmt_node
        }
        this.readToken()

        ifstmt_node.Condition = exp_node
        ifstmt_node.Consequence = stmt_block_node

        current_token = this.peek()
        if (current_token.TypeOfToken != ELSE) {
            return ifstmt_node
        }

        this.readToken()
        current_token = this.peek()
        if (current_token.TypeOfToken != OPENING_CURLY) {
            return ifstmt_node
        }

        this.readToken()
        stmt_block_node = this.parse_stmt_block_node()

        current_token = this.peek()
        if (current_token.TypeOfToken != CLOSING_CURLY) {
            return ifstmt_node
        }

        this.readToken()

        ifstmt_node.Alternative = stmt_block_node
        return ifstmt_node
    }

    parse_prefix_expression(): ExpressionNode {
        exp_node = new PrefixExpressionNode()
        current_token = this.peek()

        this.readToken()
        switch (current_token.TypeOfToken) {
            case MINUS:
                exp_node.Op = objCodeGenerator.MINUS
                exp_node.Child, err = this.parseExpression(LOWEST)
                break
            case BOOLEAN_INVERSION:
                exp_node.Op = objCodeGenerator.BOOLEAN_INVERSION
                exp_node.Child, err = this.parseExpression(LOWEST)
                break
            default:
                return exp_node
        }

        return exp_node
    }

    parse_infix_expression(ast: ExpressionNode): ExpressionNode {
        exp_node = new InfixExpressionNode()
        current_token = this.peek()
        current_precedence = precedences[current_token.TypeOfToken]

        if (current_token.TypeOfToken == PLUS) {
            exp_node.Op = objCodeGenerator.PLUS
        } else if ((current_token.TypeOfToken == MINUS)) {
            exp_node.Op = objCodeGenerator.MINUS
        } else if (current_token.TypeOfToken == ASTERISK) {
            exp_node.Op = objCodeGenerator.MULTIPLICATION
        } else if (current_token.TypeOfToken == SLASH) {
            exp_node.Op = objCodeGenerator.DIVISION
        } else if (current_token.TypeOfToken == ASSIGNMENT_OPERATOR) {
            exp_node.Op = objCodeGenerator.ASSIGNMENT
        } else if (current_token.TypeOfToken == EQUAL_TO) {
            exp_node.Op = objCodeGenerator.EQUAL_TO
        } else if (current_token.TypeOfToken == NOT_EQUAL_TO) {
            exp_node.Op = objCodeGenerator.NOT_EQUAL_TO
        } else if (current_token.TypeOfToken == GREATER_THAN) {
            exp_node.Op = objCodeGenerator.GREATER_THAN
        } else {
            console.log("expected operator")
        }

        this.readToken()
        current_token = this.peek()

        exp_node.Left = ast
        exp_node.Right = this.parseExpression(current_precedence)
        return exp_node
    }
}

type Operator = int

const OP_PLUS = 0
const OP_MINUS = 1
const OP_MULTIPLICATION = 2
const OP_DIVISION = 3
const OP_ASSIGNMENT = 4
const OP_EQUAL_TO = 5
const OP_NOT_EQUAL_TO = 6
const OP_BOOLEAN_INVERSION = 7
const OP_GREATER_THAN = 8

class AST_Node {
}

class StatementNode extends AST_Node {
}

class ExpressionNode extends StatementNode {
}

class StatementBlockNode extends AST_Node {
    Statements: StatementNode[]
}

class IfStmtNode extends StatementNode {
    Condition: ExpressionNode
    Consequence: StatementBlockNode
    Alternative: StatementBlockNode
}

class LetStmtNode extends StatementNode {
    Identifier: IdentifierNode
    InitializationExpr: ExpressionNode
}

class AssignementStmt extends StatementNode {
    Identifier: IdentifierNode
    InitializationExpr: ExpressionNode
}

class IntegerLiteralNode extends ExpressionNode {
    Value: int64
}

class KeywordBooleanNode extends ExpressionNode {
    Value: bool
}

class IdentifierNode extends ExpressionNode {
    Value: string
}

class PrefixExpressionNode extends ExpressionNode {
    Op: Operator
    Child: ExpressionNode
}

class InfixExpressionNode extends ExpressionNode {
    Op: Operator
    Left: ExpressionNode
    Right: ExpressionNode
}
