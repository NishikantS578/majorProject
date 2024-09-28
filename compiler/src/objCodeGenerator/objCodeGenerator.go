package objCodeGenerator

import (
	"fmt"
	"majorProject/compiler/vm"
)

type Node interface {
	literalString()
}

type ExpressionNode interface {
	StatementNode
	expressionNode()
}

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

type Program struct {
	Statements []StatementNode
}

func (prog *Program) literalString() {}

type StatementNode interface {
	Node
	statementNode()
}

type ObjCodeGenerator struct {
	InstructionList vm.Instructions
	ast             Node
	ConstantPool    []vm.Data
}

func New(progAst Node) ObjCodeGenerator {
	return ObjCodeGenerator{
		InstructionList: vm.Instructions([]byte{}),
		ast:             progAst,
	}
}

func (objCodeGenerator *ObjCodeGenerator) Generate(node Node) {
	if node == nil {
		return
	}
	switch node := (node).(type) {
	case *Program:
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
	}
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
	objCodeGenerator.ConstantPool = append(
		objCodeGenerator.ConstantPool,
		data,
	)

	return len(objCodeGenerator.ConstantPool) - 1
}
