package main

import (
	"fmt"
	"os"
)

type Generator struct {
	parseTree    NodeProg
	assemblyCode string
	stackSize    int
	vars         map[string]Variable
}

type Variable struct {
	stackLoc int
}

func (generator *Generator) initialize(parseTree NodeProg) {
	generator.parseTree = parseTree
	generator.vars = make(map[string]Variable)
}

func (generator *Generator) push(reg string) {
	generator.assemblyCode += "push " + reg + "\n"
	generator.stackSize++
}

func (generator *Generator) pop(reg string) {
	generator.assemblyCode += "pop " + reg + "\n"
	generator.stackSize--
}

func (generator *Generator) generateExpr(expr *NodeExpr) {
	if expr.intLitNode != nil {
		generator.assemblyCode += "mov rax, " + expr.intLitNode.intLit.value + "\n"
		generator.push("rax")
	} else if expr.identNode != nil {
		varLoc, ok := generator.vars[expr.identNode.ident.value]
		if !ok {
			println("Undeclared Identifier: " + expr.identNode.ident.value)
			os.Exit(0)
		}
		generator.push(fmt.Sprintf("QWORD [rsp + %d]", (generator.stackSize-varLoc.stackLoc-1)*8))
	}
}

func (generator *Generator) generateStmt(stmt NodeStmt) {
	if stmt.exitStmtNode != nil {
		generator.generateExpr(stmt.exitStmtNode.exprNode)
		generator.assemblyCode += "mov rax, 60\n"
		generator.pop("rdi")
		generator.assemblyCode += "syscall\n"
	} else if stmt.varDefStmtNode != nil {
		_, ok := generator.vars[stmt.varDefStmtNode.ident.ident.value]
		if ok {
			println("Redeclaration of variable '" + stmt.varDefStmtNode.ident.ident.value + "'")
			os.Exit(0)
		}
		var identifier = stmt.varDefStmtNode.ident.ident.value
		generator.vars[identifier] = Variable{stackLoc: generator.stackSize}
		generator.generateExpr(stmt.varDefStmtNode.exprNode)
	}
}

func (generator *Generator) generateProg() {
	generator.assemblyCode += "global _start\n_start:\n"

	for _, stmt := range generator.parseTree.stmtNodes {
		generator.generateStmt(stmt)
	}

}
