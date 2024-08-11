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
	generator.assemblyCode += "push " + reg + "\n\n"
	generator.stackSize++
}

func (generator *Generator) pop(reg string) {
	generator.assemblyCode += "pop " + reg + "\n"
	generator.stackSize--
}

func (generator *Generator) generateTerm(termNode *NodeTerm) {
	if termNode.identNode != nil {
		varLoc, ok := generator.vars[termNode.identNode.ident.value]
		println(termNode.identNode.ident.value, varLoc.stackLoc)
		if !ok {
			println("Undeclared Identifier: " + termNode.identNode.ident.value)
			os.Exit(0)
		}
		generator.push(fmt.Sprintf("QWORD [rsp + %d]", (generator.stackSize-varLoc.stackLoc-1)*8))
	} else if termNode.intLitNode != nil {
		generator.assemblyCode += "mov rax, " + termNode.intLitNode.intLit.value + "\n"
		generator.push("rax")
	} else if termNode.exprNode != nil {
		generator.generateExpr(termNode.exprNode)
	}
}

func (generator *Generator) generateFactor(factor *NodeFactor) {
	if factor.op == 0 {
		generator.generateTerm(factor.lTermNode)
	} else if factor.rTermNode == nil {
		generator.generateTerm(factor.lTermNode)
	} else if factor.lTermNode != nil {
		generator.generateTerm(factor.lTermNode)
		generator.generateTerm(factor.rTermNode)
		if factor.op == '*' {
			generator.pop("rbx")
			generator.pop("rax")
			generator.assemblyCode += "mul rbx\n\n"
			generator.push("rax")
		} else if factor.op == '/' {
			generator.pop("rbx")
			generator.pop("rax")
			generator.assemblyCode += "div rbx\n\n"
			generator.push("rax")
		}
	} else {
		println("Expected expression")
	}
}

func (generator *Generator) generateExpr(expr *NodeExpr) {
	if expr.op == 0 {
		generator.generateFactor(expr.lFactorNode)
	} else if expr.rFactorNode == nil {
		generator.generateFactor(expr.lFactorNode)
	} else if expr.lFactorNode != nil {
		generator.generateFactor(expr.lFactorNode)
		generator.generateFactor(expr.rFactorNode)
		if expr.op == '+' {
			generator.pop("rdi")
			generator.pop("rax")
			generator.assemblyCode += "add rax, rdi\n\n"
			generator.push("rax")
		} else if expr.op == '-' {
			generator.pop("rdi")
			generator.pop("rax")
			generator.assemblyCode += "sub rax, rdi\n\n"
			generator.push("rax")
		}
	} else {
		println("Expected expression")
	}
}

func (generator *Generator) generateStmt(stmt NodeStmt) {
	if stmt.exitStmtNode != nil {
		if stmt.exitStmtNode.exprNode != nil {
			generator.generateExpr(stmt.exitStmtNode.exprNode)
		} else if stmt.exitStmtNode.termNode != nil {
			generator.generateTerm(stmt.exitStmtNode.termNode)
		}
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
		if stmt.varDefStmtNode.exprNode != nil {
			generator.generateExpr(stmt.varDefStmtNode.exprNode)
		} else if stmt.varDefStmtNode.termNode != nil {
			generator.generateTerm(stmt.varDefStmtNode.termNode)
		}
	}
}

func (generator *Generator) generateProg() {
	generator.assemblyCode += "global _start\n_start:\n"

	for _, stmt := range generator.parseTree.stmtNodes {
		generator.generateStmt(stmt)
	}

}
