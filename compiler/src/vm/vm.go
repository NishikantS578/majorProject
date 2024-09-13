package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Data interface {
	Type() string
	Inspect() string
}

type Integer struct {
	Value int64
}

type Instructions []byte

type OpCode byte

type Vm struct {
	sp           int
	stack        []Data
	constantPool []Data
	instructions Instructions
}

const (
	INTEGER_DATA = "INTEGER"
)

func (integer *Integer) Type() string {
	return INTEGER_DATA
}

func MakeInstruction(op OpCode, operands ...int) []byte {
	var def, ok = Definitions[op]
	if !ok {
		return []byte{}
	}

	var instructionLen int = 0
	instructionLen += 1

	for _, s := range def.OperandWidths {
		instructionLen += s
	}

	var ins = make([]byte, instructionLen)
	ins[0] = byte(op)
	var offset = 1

	for i, operand := range operands {
		var operandWidth = def.OperandWidths[i]
		switch operandWidth {
		case 1:
			ins[offset] = byte(operand)
		case 2:
			binary.BigEndian.PutUint16(ins[offset:], uint16(operand))
		}
		offset += operandWidth
	}
	return ins
}

func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

type Definition struct {
	Name          string
	OperandWidths []int
}

var Definitions = map[OpCode]*Definition{
	OpConstant: {
		Name:          "OpConstant",
		OperandWidths: []int{2},
	},
	OpAddition: {
		Name: "OpAddition",
	},
}

const (
	OpConstant OpCode = iota
	OpAddition
)

func New(ins Instructions, constPool []Data) Vm {
	vm := Vm{
		sp: 0, stack: []Data{}, constantPool: constPool,
		instructions: ins,
	}
	vm.stack = make([]Data, 2048)
	return vm
}

func (vm *Vm) Execute() {
	var ip = 0
	for ip < len(vm.instructions) {
		var op = OpCode(vm.instructions[ip])
		switch op {
		case OpConstant:
			ip++
			vm.push(
				vm.constantPool[binary.BigEndian.Uint16(
					vm.instructions[ip:ip+2],
				)],
			)
			ip++
		case OpAddition:
			n1, err := vm.pop()
			if err != nil {
				println("err")
				return
			}
			n2, err := vm.pop()
			if err != nil {
				println("err")
				return
			}
			rightValue := n1.(*Integer).Value
			leftValue := n2.(*Integer).Value
			vm.push(&Integer{Value: leftValue + rightValue})
		}
		ip++
	}
}

func (vm *Vm) StackTop() (Data, error) {
	if vm.sp == 0 {
		return &Integer{}, errors.New("stack underflow")
	}
	return vm.stack[vm.sp-1], nil
}

func (vm *Vm) pop() (Data, error) {
	d, err := vm.StackTop()
	if err == nil {
		vm.sp -= 1
	}
	return d, err
}

func (vm *Vm) push(data Data) {
	vm.stack[vm.sp] = data
	vm.sp += 1
}
