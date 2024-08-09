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

func (generator *Generator) generateTerm(termNode *NodeTerm) {
	if termNode.identNode != nil {
		varLoc, ok := generator.vars[termNode.identNode.ident.value]
		if !ok {
			println("Undeclared Identifier: " + termNode.identNode.ident.value)
			os.Exit(0)
		}
		generator.push(fmt.Sprintf("QWORD [rsp + %d]", (generator.stackSize-varLoc.stackLoc-1)*8))
	} else if termNode.intLitNode != nil {
		generator.assemblyCode += "mov rax, " + termNode.intLitNode.intLit.value + "\n"
		generator.push("rax")
	}
}

func (generator *Generator) generateAddExpr(addExprNode *NodeAddExpr) {
	if addExprNode.lExprNode != nil {
		generator.generateExpr(addExprNode.lExprNode)
	} else if addExprNode.lTermNode != nil {
		generator.generateTerm(addExprNode.lTermNode)
	}

	if addExprNode.rExprNode != nil {
		generator.generateExpr(addExprNode.rExprNode)
	} else if addExprNode.rTermNode != nil {
		generator.generateTerm(addExprNode.rTermNode)
	}
	generator.pop("rax")
	generator.pop("rbx")
	generator.assemblyCode += "add rax, rbx\n"
	generator.push("rax")
}

func (generator *Generator) generateExpr(expr *NodeExpr) {
	if expr.addExprNode != nil {
		generator.generateAddExpr(expr.addExprNode)
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
