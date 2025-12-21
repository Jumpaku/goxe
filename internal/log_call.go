package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func (x *Xtrace) newCallLogStmt(signature string) ast.Stmt {
	// PrintlnCall("[CALL] <signature>")
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: ast.NewIdent(x.IdentifierPrintlnCall()),
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.INT,
					Value: fmt.Sprintf(`%d`, x.LineWidth),
				},
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`%q`, signature),
				},
				&ast.Ident{Name: x.IdentShowTimestamp()},
				&ast.Ident{Name: x.IdentShowGoroutine()},
			},
		},
	}
}

func (x *Xtrace) newReturnLogStmt(signature string) ast.Stmt {
	// PrintlnReturn("[RETURN] <signature>")
	return &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: ast.NewIdent(x.IdentifierPrintlnReturn()),
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.INT,
					Value: fmt.Sprintf(`%d`, x.LineWidth),
				},
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`%q`, signature),
				},
				&ast.Ident{Name: x.IdentShowTimestamp()},
				&ast.Ident{Name: x.IdentShowGoroutine()},
			},
		},
	}
}

func (x *Xtrace) logCall(c *astutil.Cursor, info *FuncInfo) {
	if !x.TraceCall {
		return
	}
	signature := strings.Join(strings.Fields(x.fragment(info.Signature())), " ")
	body := info.Body
	body.List = append(
		[]ast.Stmt{
			x.newCallLogStmt(signature),
			x.newReturnLogStmt(signature),
		},
		body.List...,
	)
	c.Replace(body)
	x.libraryRequired = true
}
