package xtracego

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/samber/lo/mutable"
	"golang.org/x/tools/go/ast/astutil"
)

func (s XTrace) newIfElseLogStmt(typ IfElseType, cond string) ast.Stmt {
	// log.Println(fmt.Sprintf(`[IF] if condition {`))
	// log.Println(fmt.Sprintf(`[ELSEIF] else if condition {`))
	// log.Println(fmt.Sprintf(`[ELSE] else {`))
	var content string
	switch typ {
	case IfElseType_If:
		content = fmt.Sprintf("[IF] %s", cond)
	case IfElseType_ElseIf:
		content = fmt.Sprintf("[ELSEIF] %s", cond)
	case IfElseType_Else:
		content = "[ELSE]"
	}
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: ast.NewIdent(s.prefix + "_log.Println"),
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent(s.prefix + "_fmt.Sprintf"),
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf("%q", content),
						},
					},
				},
			},
		},
	}
}

func (s XTrace) logIfElse(c *astutil.Cursor, info *IfElseInfo) {
	stmts := []ast.Stmt{}
	for i, ifStmt := range info.Parents {
		if i == 0 {
			continue
		}
		frag := s.fragmentLine(ifStmt.If)
		stmts = append(stmts, s.newStatementLogStmt(s.fset.Position(ifStmt.If), frag))
	}
	if len(info.Parents) > 0 {
		frag := s.fragmentLine(info.IfStmt.If)
		stmts = append(stmts, s.newStatementLogStmt(s.fset.Position(info.IfStmt.If), frag))
		if info.ElseBody != nil {
			frag := s.fragmentLine(info.IfStmt.Body.Rbrace)
			stmts = append(stmts, s.newStatementLogStmt(s.fset.Position(info.IfStmt.Body.Rbrace), frag))
		}
	}

	vars := info.Variables()
	shadow := map[string]bool{}
	varStmts := []ast.Stmt{}
	for i := len(vars) - 1; i >= 0; i-- {
		varStmts = append(varStmts, s.newVariableLogStmt(vars[i].Name, shadow[vars[i].Name]))
		shadow[vars[i].Name] = true
	}
	mutable.Reverse(varStmts)

	stmts = append(stmts, varStmts...)

	if info.Body != nil {
		info.Body.List = append(stmts, info.Body.List...)
		c.Replace(info.Body)
	}
	if info.ElseBody != nil {
		info.ElseBody.List = append(stmts, info.ElseBody.List...)
		c.Replace(info.ElseBody)
	}
}
