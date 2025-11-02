package xtracego

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"

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
	pattern := regexp.MustCompile(`\s+`)
	frag := s.fragment(info.Condition())
	frag = pattern.ReplaceAllString(frag, " ")
	info.Body.List = append([]ast.Stmt{s.newIfElseLogStmt(info.IfElseType, frag)}, info.Body.List...)
	c.Replace(info.Body)
}
