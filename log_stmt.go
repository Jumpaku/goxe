package xtracego

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/tools/go/ast/astutil"
)

func (s XTrace) newStatementLogStmt(pos token.Position, fragment string) ast.Stmt {
	// log.Println(fmt.Sprintf(`if a == 1 { /* path/to/source.go:123:45 */`))
	content := fmt.Sprintf("%s ", fragment)
	content = strings.ReplaceAll(content, "\t", "    ")
	if len(content) < s.lineWidth {
		content += strings.Repeat(".", s.lineWidth-len(content))
	}
	content += fmt.Sprintf(" [ %s:%d:%d ]", pos.Filename, pos.Line, pos.Column)
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: ast.NewIdent(s.prefix + "_log.Println"),
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", content),
				},
			},
		},
	}
}

func (s XTrace) newStatementLogDecl(pos token.Position, fragment string) *ast.GenDecl {
	//var _ = func() int {
	//	log.Println(fmt.Sprintf(`if a == 1 { /* path/to/source.go:123:45 */`))
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
									s.newStatementLogStmt(pos, fragment),
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

func (s XTrace) logFileStatement(c *astutil.Cursor, node *ast.GenDecl) {
	for _, spec := range node.Specs {
		spec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		pos := s.fset.Position(spec.Pos())
		frag := s.fragmentLine(spec.Pos(), spec.End())
		c.InsertBefore(s.newStatementLogDecl(pos, frag))
	}
}

func (s XTrace) tryLogLocalStatement(c *astutil.Cursor, node ast.Stmt) {
	{
		insertable := false
		switch parent := c.Parent().(type) {
		case *ast.BlockStmt:
			if lo.Contains(parent.List, node) {
				insertable = true
			}
		case *ast.SwitchStmt:
		case *ast.TypeSwitchStmt:
		case *ast.SelectStmt:
		case *ast.ForStmt:
		case *ast.RangeStmt:
		case *ast.IfStmt:
		case *ast.CaseClause:
			if lo.Contains(parent.Body, node) {
				insertable = true
			}
		case *ast.CommClause:
			if lo.Contains(parent.Body, node) {
				insertable = true
			}
		case *ast.ReturnStmt:
		case *ast.DeferStmt:
		case *ast.GoStmt:
		case *ast.BranchStmt:
		case *ast.LabeledStmt:
		case *ast.SendStmt:
		case *ast.IncDecStmt:
		case *ast.ExprStmt:
		case *ast.AssignStmt:
		case *ast.EmptyStmt:
		}

		tracable := true
		switch node.(type) {
		case *ast.BlockStmt:
			tracable = false
		case *ast.SwitchStmt:
		case *ast.TypeSwitchStmt:
		case *ast.SelectStmt:
		case *ast.ForStmt:
		case *ast.RangeStmt:
		case *ast.IfStmt:
			if _, ok := c.Parent().(*ast.IfStmt); ok {
				tracable = false
			}
		case *ast.CaseClause:
			tracable = false
		case *ast.CommClause:
			tracable = false
		case *ast.ReturnStmt:
		case *ast.DeferStmt:
		case *ast.GoStmt:
		case *ast.BranchStmt:
		case *ast.LabeledStmt:
		case *ast.SendStmt:
		case *ast.IncDecStmt:
		case *ast.ExprStmt:
		case *ast.AssignStmt:
		case *ast.EmptyStmt:
			tracable = false
		}

		if !insertable || !tracable {
			return
		}
	}
	pos := s.fset.Position(node.Pos())
	frag := s.fragmentLine(node.Pos(), node.End())
	c.InsertBefore(s.newStatementLogStmt(pos, frag))
}
