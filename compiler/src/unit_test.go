package main

import (
	"majorProject/compiler/lexer"
	"majorProject/compiler/objCodeGenerator"
	"majorProject/compiler/parser"
	"majorProject/compiler/vm"
	"testing"
)

func execute(input_string string) (string, error) {
	var err error
	var progLexer = lexer.New(input_string)
	progLexer.Tokenize()
	var progParser = parser.New(progLexer.TokenArr)
	if progParser.ParseProgram() == 0 {
		return "", err
	}

	var objCode = objCodeGenerator.New(&progParser.Ast)
	objCode.Generate(&progParser.Ast)

	var machine = vm.New(
		objCode.InstructionList,
		objCode.ConstantPool,
	)
	err = machine.Execute()
	if err != nil {
		return "", err
	}

	var stackTop vm.Data
	stackTop, err = machine.StackTop()
	if err != nil {
		return "empty stack", nil
	}

	return stackTop.Inspect(), err
}

func TestCompiler(t *testing.T) {
	var err error
	var output_actual string
	var outputs_expected = map[string]string{
		"1+2": "3",
		"4-5": "-1",
		"if(true){4+1}": "5",
		"if(4==4){3*3}": "9",
		"if(4==3){1+1}else{9*2}": "18",
		"if(4==4){1+1}else{9*2}": "2",
	}

	for input_string, output_expected := range outputs_expected {
		output_actual, err = execute(input_string)
		if err != nil {
			t.Fatal("error while executing", err)
		}
		if output_actual != output_expected {
			t.Fatal("\n\texpected: \n\t\"", output_expected, "\"\n\tgot\n\t\"", output_actual, "\"")
		}
	}
}
