package objCodeGenerator

import (
	"errors"
	"fmt"
	"majorProject/compiler/vm"
)

type Operator int

const (
	PLUS = iota
	MINUS
	MULTIPLICATION
	DIVISION
	ASSIGNMENT
	EQUAL_TO
	NOT_EQUAL_TO
	BOOLEAN_INVERSION
	GREATER_THAN
)

type Node interface {
	literalString()
}

type StatementNode interface {
	Node
	statementNode()
}

type ExpressionNode interface {
	StatementNode
	expressionNode()
}

type StatementBlockNode struct {
	Statements []StatementNode
}

type IfStmtNode struct {
	Condition   ExpressionNode
	Consequence StatementBlockNode
	Alternative StatementBlockNode
}

func (i *IfStmtNode) statementNode() {}
func (i *IfStmtNode) literalString() {}

type LetStmtNode struct {
	Identifier         IdentifierNode
	InitializationExpr ExpressionNode
}

func (l *LetStmtNode) statementNode() {}
func (l *LetStmtNode) literalString() {}

type AssignementStmt struct {
	Identifier         IdentifierNode
	InitializationExpr ExpressionNode
}

func (l *AssignementStmt) statementNode() {}
func (l *AssignementStmt) literalString() {}

type IntegerLiteralNode struct {
	Value int64
}

func (i *IntegerLiteralNode) literalString()  {}
func (i *IntegerLiteralNode) expressionNode() {}
func (i *IntegerLiteralNode) statementNode()  {}

type KeywordBooleanNode struct {
	Value bool
}

func (b *KeywordBooleanNode) literalString()  {}
func (b *KeywordBooleanNode) expressionNode() {}
func (b *KeywordBooleanNode) statementNode()  {}

type IdentifierNode struct {
	Value string
}

func (i *IdentifierNode) literalString()  {}
func (i *IdentifierNode) expressionNode() {}
func (i *IdentifierNode) statementNode()  {}

type PrefixExpressionNode struct {
	Op    Operator
	Child ExpressionNode
}

func (p *PrefixExpressionNode) literalString()  {}
func (p *PrefixExpressionNode) expressionNode() {}
func (p *PrefixExpressionNode) statementNode()  {}

type InfixExpressionNode struct {
	Op    Operator
	Left  ExpressionNode
	Right ExpressionNode
}

func (i *InfixExpressionNode) literalString()  {}
func (i *InfixExpressionNode) expressionNode() {}
func (i *InfixExpressionNode) statementNode()  {}

func (prog *StatementBlockNode) literalString() {}

type ObjCodeGenerator struct {
	InstructionList vm.Instructions
	ast             Node
	ConstantPool    *[]vm.Data
	SymbolTable     *vm.SymbolTable
}

func (objCodeGenerator *ObjCodeGenerator) SetNewInput(ast Node) {
	objCodeGenerator.InstructionList = vm.Instructions([]byte{})
	objCodeGenerator.ast = ast
}

func New(progAst Node) *ObjCodeGenerator {
	return &ObjCodeGenerator{
		InstructionList: vm.Instructions([]byte{}),
		ast:             progAst,
		SymbolTable:     vm.NewSymbolTable(),
		ConstantPool:    &[]vm.Data{},
	}
}

func (objCodeGenerator *ObjCodeGenerator) Generate(node Node) error {
	if node == nil {
		return nil
	}
	switch node := (node).(type) {
	case *StatementBlockNode:
		for _, s := range node.Statements {
			objCodeGenerator.Generate(s)
		}
	case *InfixExpressionNode:
		objCodeGenerator.Generate(node.Left)
		objCodeGenerator.Generate(node.Right)
		switch node.Op {
		case PLUS:
			objCodeGenerator.emit(vm.OpAddition)
		case MINUS:
			objCodeGenerator.emit(vm.OpSubtraction)
		case MULTIPLICATION:
			objCodeGenerator.emit(vm.OpMultiplication)
		case DIVISION:
			objCodeGenerator.emit(vm.OpDivision)
		case EQUAL_TO:
			objCodeGenerator.emit(vm.OpEqual)
		case NOT_EQUAL_TO:
			objCodeGenerator.emit(vm.OpNotEqual)
		case GREATER_THAN:
			objCodeGenerator.emit(vm.OpGreaterThan)
		default:
			fmt.Println("unknown operator", node.Op)
		}
	case *PrefixExpressionNode:
		objCodeGenerator.Generate(node.Child)
		switch node.Op {
		case MINUS:
			objCodeGenerator.emit(vm.OpNegation)
		case BOOLEAN_INVERSION:
			objCodeGenerator.emit(vm.OpBooleanInversion)
		default:
			fmt.Println("expected prefix operator")
		}
	case *IntegerLiteralNode:
		var integer = &vm.Integer{Value: node.Value}
		var addr = objCodeGenerator.addConstant(integer)
		objCodeGenerator.emit(vm.OpConstant, addr)
	case *KeywordBooleanNode:
		if node.Value {
			objCodeGenerator.emit(vm.OpTrue)
		} else if !node.Value {
			objCodeGenerator.emit(vm.OpFalse)
		}
	case *IfStmtNode:
		objCodeGenerator.Generate(node.Condition)
		var jump_not_truthy_pos = len(objCodeGenerator.InstructionList)
		objCodeGenerator.emit(vm.OpJumpNotTruthy, 9999)

		objCodeGenerator.Generate(&node.Consequence)

		var jump_pos = len(objCodeGenerator.InstructionList)
		objCodeGenerator.emit(vm.OpJump, 9999)

		var after_consequence_pos = len(objCodeGenerator.InstructionList)

		objCodeGenerator.changeOperand(jump_not_truthy_pos, after_consequence_pos)

		objCodeGenerator.Generate(&node.Alternative)
		var after_alternative_pos = len(objCodeGenerator.InstructionList)

		objCodeGenerator.changeOperand(jump_pos, after_alternative_pos)
	case *LetStmtNode:
		if node.InitializationExpr == nil {
			objCodeGenerator.emit(vm.OpSetGlobal, 0)
		} else {
			var err = objCodeGenerator.Generate(node.InitializationExpr)
			if err != nil {
				return err
			}
			var symbol = objCodeGenerator.SymbolTable.Define(node.Identifier.Value)
			objCodeGenerator.emit(vm.OpSetGlobal, symbol.Index)
		}
	case *AssignementStmt:
		var err = objCodeGenerator.Generate(node.InitializationExpr)
		if err != nil {
			return err
		}
		var symbol, exists = objCodeGenerator.SymbolTable.Resolve(node.Identifier.Value)
		if !exists {
			fmt.Println("undefined symbol: ", node.Identifier.Value)
			return errors.New("undefined symbol: " + node.Identifier.Value)
		}
		objCodeGenerator.emit(vm.OpSetGlobal, symbol.Index)
	case *IdentifierNode:
		var symbol, exists = objCodeGenerator.SymbolTable.Resolve(node.Value)
		if !exists {
			fmt.Println("undefined symbol: ", node.Value)
			return errors.New("undefined symbol: " + node.Value)
		}
		objCodeGenerator.emit(vm.OpGetGlobal, symbol.Index)
	default:
		return errors.New("unknown data while generating obj code")
	}

	return nil
}

func (objCodeGenerator *ObjCodeGenerator) emit(
	opcode vm.OpCode,
	operands ...int,
) {
	var ins []byte = vm.MakeInstruction(opcode, operands...)
	objCodeGenerator.InstructionList = append(
		objCodeGenerator.InstructionList,
		ins...,
	)
}

func (objCodeGenerator *ObjCodeGenerator) addConstant(
	data vm.Data,
) int {
	*objCodeGenerator.ConstantPool = append(
		*objCodeGenerator.ConstantPool,
		data,
	)

	return len(*objCodeGenerator.ConstantPool) - 1
}

func (objCodeGenerator *ObjCodeGenerator) replaceInstruction(pos int, new_instruction []byte) {
	for i := 0; i < len(new_instruction); i++ {
		objCodeGenerator.InstructionList[pos+i] = new_instruction[i]
	}
}

func (objCodeGenerator *ObjCodeGenerator) changeOperand(opPos int, operand int) {
	var op = vm.OpCode(objCodeGenerator.InstructionList[opPos])
	var new_instruction = vm.MakeInstruction(op, operand)
	objCodeGenerator.replaceInstruction(opPos, new_instruction)
}
