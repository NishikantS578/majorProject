package main

import (
	"fmt"
	"io"
	"majorProject/compiler/lexer"
	"majorProject/compiler/objCodeGenerator"
	"majorProject/compiler/parser"
	"majorProject/compiler/repl"
	"majorProject/compiler/vm"
	"os"
)

func main() {
	var args = os.Args

	var ioSrc io.Reader
	var ioDest io.Writer
	var err error

	if len(args) == 1 {
		ioSrc = os.Stdin
		ioDest = os.Stdout
		repl.Run(ioSrc, ioDest)
	} else {
		ioSrc, err = os.Open(args[1])

		if err != nil {
			fmt.Println("Could not read file: " + args[1])
			os.Exit(0)
		}

		// ioDest, err = os.Create("app.obj")

		if err != nil {
			fmt.Println("Could not create file: app.obj")
			os.Exit(0)
		}

		var exit_code = ""
		var err error

		var progLexer = lexer.New("")
		var progParser = parser.New(progLexer.TokenArr)
		var objCode = objCodeGenerator.New(&progParser.Ast)
		var machine = vm.New(
			objCode.InstructionList,
			objCode.ConstantPool,
			objCode.SymbolTable,
		)

		var line []byte
		line, err = io.ReadAll(ioSrc)
		progLexer.SetNewInput(string(line))
		progLexer.Tokenize()

		progParser.SetNewInput(progLexer.TokenArr)
		if progParser.ParseProgram() == 0 {
			os.Exit(1)
		}

		objCode.SetNewInput(&progParser.Ast)
		objCode.Generate(&progParser.Ast)

		machine.SetNewInput(objCode.InstructionList)
		err = machine.Execute()
		if err != nil {
			fmt.Println(err)
		}

		var stackTop vm.Data
		stackTop, err = machine.StackTop()
		if err != nil {
			os.Exit(1)
		}

		if ioDest == os.Stdout {
			io.WriteString(ioDest, stackTop.Inspect())
			io.WriteString(ioDest, "\n")
		}
		exit_code = stackTop.Inspect()
		println(exit_code)
	}
}
