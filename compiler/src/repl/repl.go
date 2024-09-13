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

	for {
		if ioDest == os.Stdout {
			fmt.Fprint(ioDest, PROMPT)
		}

		if !scanner.Scan() {
			return
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
		machine.Execute()

		var stackTop, err = machine.StackTop()
		if err != nil {
			continue
		}

		io.WriteString(ioDest, stackTop.Inspect())
		io.WriteString(ioDest, "\n")
	}
}
