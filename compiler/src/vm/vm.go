package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Vm struct {
	sp           int
	stack        []Data
	constantPool []Data
	instructions Instructions
}

type Data interface {
	Type() string
	Inspect() string
}

type Instructions []byte

type OpCode byte

var True = &Boolean{Value: true}
var False = &Boolean{Value: false}
var StackSize = 2048

func New(ins Instructions, constPool []Data) Vm {
	vm := Vm{
		sp: 0, stack: []Data{}, constantPool: constPool,
		instructions: ins,
	}
	vm.stack = make([]Data, StackSize)
	return vm
}

func (vm *Vm) Execute() error {
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
				return err
			}
			n2, err := vm.pop()
			if err != nil {
				return err
			}
			rightValue := n1.(*Integer).Value
			leftValue := n2.(*Integer).Value
			vm.push(&Integer{Value: leftValue + rightValue})
		case OpSubtraction:
			n1, err := vm.pop()
			if err != nil {
				return err
			}
			n2, err := vm.pop()
			if err != nil {
				return err
			}
			rightValue := n1.(*Integer).Value
			leftValue := n2.(*Integer).Value
			vm.push(&Integer{Value: leftValue - rightValue})
		case OpMultiplication:
			n1, err := vm.pop()
			if err != nil {
				return err
			}
			n2, err := vm.pop()
			if err != nil {
				return err
			}
			rightValue := n1.(*Integer).Value
			leftValue := n2.(*Integer).Value
			vm.push(&Integer{Value: leftValue * rightValue})
		case OpDivision:
			n1, err := vm.pop()
			if err != nil {
				return err
			}
			n2, err := vm.pop()
			if err != nil {
				return err
			}
			rightValue := n1.(*Integer).Value
			leftValue := n2.(*Integer).Value
			vm.push(&Integer{Value: leftValue / rightValue})
		case OpNegation:
			var n Data
			var err error
			n, err = vm.pop()
			if err != nil {
				return err
			}
			vm.push(&Integer{Value: -n.(*Integer).Value})
		case OpTrue:
			var err error
			err = vm.push(True)
			if err != nil {
				return err
			}
		case OpFalse:
			var err error
			err = vm.push(False)
			if err != nil {
				return err
			}
		case OpEqual:
			var rightValue, err = vm.pop()
			if err != nil {
				return err
			}
			var leftValue Data
			leftValue, err = vm.pop()
			if err != nil {
				return err
			}
			if leftValue.(*Integer).Value == rightValue.(*Integer).Value {
				err = vm.push(&Boolean{Value: true})
				if err != nil {
					return err
				}
			} else {
				err = vm.push(&Boolean{Value: false})
				if err != nil {
					return err
				}
			}
		case OpNotEqual:
			var rightValue, err = vm.pop()
			if err != nil {
				return err
			}
			var leftValue Data
			leftValue, err = vm.pop()
			if err != nil {
				return err
			}
			if leftValue.(*Integer).Value != rightValue.(*Integer).Value {
				err = vm.push(&Boolean{Value: true})
				if err != nil {
					return err
				}
			} else {
				err = vm.push(&Boolean{Value: false})
				if err != nil {
					return err
				}
			}
		case OpGreaterThan:
			var rightValue, err = vm.pop()
			if err != nil {
				return err
			}
			var leftValue Data
			leftValue, err = vm.pop()
			if err != nil {
				return err
			}
			if leftValue.(*Integer).Value > rightValue.(*Integer).Value {
				err = vm.push(&Boolean{Value: true})
				if err != nil {
					return err
				}
			} else {
				err = vm.push(&Boolean{Value: false})
				if err != nil {
					return err
				}
			}
		case OpBooleanInversion:
			var d Data
			var err error
			d, err = vm.pop()
			if err != nil {
				return err
			}
			switch d.(*Boolean).Value {
			case false:
				vm.push(True)
			default:
				vm.push(False)
			}
		case OpJumpNotTruthy:
			ip++
			var condition_value, err = vm.pop()
			if err != nil{
				return err
			}
			if condition_value.(*Boolean).Value == False.Value{
				ip = int(binary.BigEndian.Uint16(vm.instructions[ip:ip+2])) - 1
			} else{
				ip++
			}
		case OpJump:
			ip++
			ip = int(binary.BigEndian.Uint16(vm.instructions[ip:ip+2])) - 1
		default:
			fmt.Println("unkown instruction", op)
		}
		ip++
	}
	return nil
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

func (vm *Vm) push(data Data) error {
	if vm.sp == StackSize {
		return errors.New("stack overflow")
	}
	vm.stack[vm.sp] = data
	vm.sp += 1
	return nil
}
