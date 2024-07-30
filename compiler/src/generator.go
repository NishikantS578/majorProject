package main

type Generator struct {
	parseTree NodeProg
}

func (generator *Generator) initialize(parseTree NodeProg) {
	generator.parseTree = parseTree
}

func (generator *Generator) generate() string {
	var assemblyCode string
	assemblyCode += "global _start\n_start:\n"
	assemblyCode += "	mov rax, 60\n"
	assemblyCode += "	mov rdi, " + "60" + "\n"
	assemblyCode += "	syscall\n"
	return assemblyCode
}
