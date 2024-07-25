package main

type Generator struct {
	parseTree NodeExit
}

func (generator *Generator) initialize(parseTree NodeExit) {
	generator.parseTree = parseTree
}

func (generator *Generator) generate() string {
	var assemblyCode string
	assemblyCode += "global _start\n_start:\n"
	assemblyCode += "	mov rax, 60\n"
	assemblyCode += "	mov rdi, " + generator.parseTree.exprNode.intLit.value + "\n"
	assemblyCode += "	syscall\n"
	return assemblyCode
}
