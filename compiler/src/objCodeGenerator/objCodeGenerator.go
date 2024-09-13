package objCodeGenerator

import (
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
	ASTERISK
	SLASH
)

type IntegerLiteralNode struct {
	Value int64
}

func (i *IntegerLiteralNode) literalString() {
}
func (i *IntegerLiteralNode) expressionNode() {
}
func (i *IntegerLiteralNode) statementNode() {
}

type InfixExpressionNode struct {
	Op    Operator
	Left  *ExpressionNode
	Right *ExpressionNode
}

func (i *InfixExpressionNode) literalString() {
}
func (i *InfixExpressionNode) expressionNode() {
}
func (i *InfixExpressionNode) statementNode() {
}

type Program struct {
	Statements []StatementNode
}

func (prog *Program) literalString() {
}

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
		objCodeGenerator.Generate(*node.Left)
		objCodeGenerator.Generate(*node.Right)
		switch node.Op {
		case PLUS:
			objCodeGenerator.emit(vm.OpAddition)
		}
	case *IntegerLiteralNode:
		var integer = &vm.Integer{Value: node.Value}
		var addr = objCodeGenerator.addConstant(integer)
		objCodeGenerator.emit(vm.OpConstant, addr)
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
