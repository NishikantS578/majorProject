package repl

import (
	"bufio"
	"fmt"
	"io"
	"majorProject/compiler/lexer"
	"majorProject/compiler/objCodeGenerator"
	"majorProject/compiler/parser"
	"majorProject/compiler/vm"
	"os"
)

const PROMPT = ">> "

func Run(ioSrc io.Reader, ioDest io.Writer) {

	var scanner = bufio.NewScanner(ioSrc)
	var exit_code = ""
	var err error

	for {
		if ioDest == os.Stdout {
			fmt.Fprint(ioDest, PROMPT)
		}

		if !scanner.Scan() {
			break
		}

		var line = scanner.Text()
		var progLexer = lexer.New(line)
		progLexer.Tokenize()
		var progParser = parser.New(progLexer.TokenArr)
		if progParser.Parse_program() == 0 {
			continue
		}

		var objCode = objCodeGenerator.New(&progParser.Ast)
		objCode.Generate(&progParser.Ast)

		var machine = vm.New(
			objCode.InstructionList,
			objCode.ConstantPool,
		)
		err = machine.Execute()
		if err != nil {
			fmt.Println(err)
		}

		var stackTop, err = machine.StackTop()
		if err != nil {
			continue
		}

		if ioDest == os.Stdout {
			io.WriteString(ioDest, stackTop.Inspect())
			io.WriteString(ioDest, "\n")
		}
		exit_code = stackTop.Inspect()
	}
	println(exit_code)
}
