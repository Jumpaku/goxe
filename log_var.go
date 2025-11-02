package xtracego

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/samber/lo/mutable"
	"golang.org/x/tools/go/ast/astutil"
)

func (s XTrace) newVariableLogStmt(name string) ast.Stmt {
	// log.Println(fmt.Sprintf(`[variable] VarName=%#v`, VarName))
	content := fmt.Sprintf("[VAR] %s=%%#v", name)
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
						&ast.Ident{Name: name},
					},
				},
			},
		},
	}
}

func (s XTrace) newVariableLogDecl(name string) *ast.GenDecl {
	//var _ = func() int {
	//	log.Println(fmt.Sprintf(`VarName: %+v path/to/source.go:123:45`, VarName))
	//	return 0
	//}()
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("_")},
				Values: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.FuncLit{
							Type: &ast.FuncType{Results: &ast.FieldList{List: []*ast.Field{{Type: ast.NewIdent("int")}}}},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									s.newVariableLogStmt(name),
									&ast.ReturnStmt{Results: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}}},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s XTrace) logFileVariable(c *astutil.Cursor, node *ast.GenDecl) {
	decls := []*ast.GenDecl{}
	for _, spec := range node.Specs {
		spec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, name := range spec.Names {
			if name.Name == "_" {
				continue
			}
			decls = append(decls, s.newVariableLogDecl(name.Name))
		}
	}
	mutable.Reverse(decls)
	for _, decl := range decls {
		c.InsertAfter(decl)
	}
}

func (s XTrace) logLocalVariable(c *astutil.Cursor, node *ast.DeclStmt) {
	decl, ok := node.Decl.(*ast.GenDecl)
	if !ok {
		return
	}
	stmts := []ast.Stmt{}
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, name := range spec.Names {
			if name.Name == "_" {
				continue
			}
			stmts = append(stmts, s.newVariableLogStmt(name.Name))
		}
	}
	mutable.Reverse(stmts)
	for _, decl := range stmts {
		c.InsertAfter(decl)
	}
}

func (s XTrace) logLocalAssignment(c *astutil.Cursor, node *ast.AssignStmt) {
	stmts := []ast.Stmt{}
	for _, lexpr := range node.Lhs {
		ident, ok := lexpr.(*ast.Ident)
		if !ok || ident.Name == "_" {
			continue
		}
		stmts = append(stmts, s.newVariableLogStmt(ident.Name))
	}
	mutable.Reverse(stmts)
	for _, decl := range stmts {
		c.InsertAfter(decl)
	}
}

func (s XTrace) logForVariables(c *astutil.Cursor, info *ForInfo) {
	body := info.Body
	vars := []ast.Stmt{}
	for _, ident := range info.Variables() {
		vars = append(vars, s.newVariableLogStmt(ident.Name))
	}
	body.List = append(vars, body.List...)
	c.Replace(body)
}
