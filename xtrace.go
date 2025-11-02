package xtracego

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"math/rand"

	"github.com/samber/lo/mutable"
	"golang.org/x/tools/go/ast/astutil"
)

type Config struct {
	TraceStmt bool
	TraceVar  bool
	TraceCall bool
	TraceCase bool
	Prefix    string
}

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func (cfg *Config) GenPrefix(seed int64) {
	r := rand.New(rand.NewSource(seed))
	v := []byte{}
	for i := 0; i < 8; i++ {
		v = append(v, alphabet[r.Intn(len(alphabet))])
	}
	cfg.Prefix = "xtracego_" + string(v)
}

func ProcessCode(cfg Config, filename string, dst io.Writer, src io.Reader) (err error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}
	src, srcBytes := buf, buf.Bytes()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.SkipObjectResolution)
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	astutil.AddNamedImport(fset, f, cfg.Prefix+"_log", "log")
	astutil.AddNamedImport(fset, f, cfg.Prefix+"_fmt", "fmt")

	s := XTrace{
		fset:      fset,
		src:       srcBytes,
		lineWidth: 80,
		prefix:    cfg.Prefix,

		funcByBody:   CollectFuncInfo(f),
		forByBody:    CollectForInfo(f),
		caseByBody:   CollectCaseInfo(f),
		ifElseByBody: CollectIfElseInfo(f),
	}

	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		switch node := c.Node().(type) {
		case *ast.GenDecl:
			switch node.Tok {
			case token.VAR:
				if _, isFile := c.Parent().(*ast.File); isFile {
					s.logFileStatement(c, node)
					s.logFileVariable(c, node)
				}
			case token.CONST:
				s.logFileStatement(c, node)
				s.logFileVariable(c, node)
			}
		case ast.Stmt:
			{
				if info, ok := s.funcByBody[node]; ok {
					s.logCall(c, info)
				}
				if info, ok := s.forByBody[node]; ok {
					s.logForVariables(c, info)
				}
				if info, ok := s.caseByBody[node]; ok {
					s.logCase(c, info)
				}
				if info, ok := s.ifElseByBody[node]; ok {
					s.logIfElse(c, info)
				}
			}

			s.tryLogLocalStatement(c, node)

			switch node := node.(type) {
			case *ast.DeclStmt:
				if decl, ok := node.Decl.(*ast.GenDecl); ok && decl.Tok == token.VAR {
					s.logLocalVariable(c, node)
				}
			case *ast.AssignStmt:
				if _, ok := c.Parent().(*ast.BlockStmt); ok {
					if node.Tok == token.ASSIGN || node.Tok == token.DEFINE {
						s.logLocalAssignment(c, node)
					}
				}
			case *ast.EmptyStmt:
			case *ast.BlockStmt:
			case *ast.ExprStmt:
			case *ast.IfStmt:
			case *ast.SwitchStmt:
			case *ast.TypeSwitchStmt:
			case *ast.CaseClause:
			case *ast.SelectStmt:
			case *ast.CommClause:
			case *ast.ForStmt:
			case *ast.RangeStmt:
			case *ast.ReturnStmt:
			case *ast.DeferStmt:
			case *ast.GoStmt:
			case *ast.BranchStmt:
			case *ast.LabeledStmt:
			case *ast.SendStmt:
			case *ast.IncDecStmt:
			}
		}

		return true
	})

	if err := printer.Fprint(dst, fset, f); err != nil {
		return fmt.Errorf("failed to print: %w", err)
	}

	return nil
}

type FuncInfo struct {
	Body     *ast.BlockStmt
	FuncDecl *ast.FuncDecl
}

func (i FuncInfo) Signature() (begin, end token.Pos) {
	return i.FuncDecl.Pos(), i.FuncDecl.Body.Pos()
}

func CollectFuncInfo(f *ast.File) (funcByBody map[ast.Stmt]*FuncInfo) {
	funcByBody = map[ast.Stmt]*FuncInfo{}
	ast.PreorderStack(f, nil, func(n ast.Node, s []ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			if node.Body != nil && len(node.Body.List) > 0 {
				funcByBody[node.Body] = &FuncInfo{
					Body:     node.Body,
					FuncDecl: node,
				}
			}
		}
		return true
	})
	return funcByBody
}

type ForInfo struct {
	Body  *ast.BlockStmt
	For   *ast.ForStmt
	Range *ast.RangeStmt
}

func (i ForInfo) Variables() (vars []*ast.Ident) {
	if r := i.Range; r != nil {
		if r.Key != nil {
			if ident, ok := r.Key.(*ast.Ident); ok && ident.Name != "_" {
				vars = append(vars, ident)
			}
		}
		if r.Value != nil {
			if ident, ok := r.Value.(*ast.Ident); ok && ident.Name != "_" {
				vars = append(vars, ident)
			}
		}
	}
	if f := i.For; f != nil {
		if assign, ok := f.Init.(*ast.AssignStmt); ok {
			for _, lhs := range assign.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok && ident.Name != "_" {
					vars = append(vars, ident)
				}
			}
		}
	}
	return vars
}

func CollectForInfo(f *ast.File) (forByBody map[ast.Stmt]*ForInfo) {
	forByBody = map[ast.Stmt]*ForInfo{}
	ast.PreorderStack(f, nil, func(n ast.Node, s []ast.Node) bool {
		switch node := n.(type) {
		case *ast.ForStmt:
			if node.Body != nil {
				forByBody[node.Body] = &ForInfo{
					Body: node.Body,
					For:  node,
				}
			}
		case *ast.RangeStmt:
			if node.Body != nil {
				forByBody[node.Body] = &ForInfo{
					Body:  node.Body,
					Range: node,
				}
			}
		}
		return true
	})
	return forByBody
}

type CaseInfo struct {
	Case *ast.CaseClause
	Comm *ast.CommClause
}

func (i CaseInfo) CaseLabel() (begin, end token.Pos) {
	if c := i.Case; c != nil {
		return c.Pos(), c.Colon
	}
	if c := i.Comm; c != nil {
		return c.Pos(), c.Colon
	}
	panic("CaseInfo must consist of one of *ast.CaseClause or *ast.CommClause.")
}

func CollectCaseInfo(f *ast.File) (caseByBody map[ast.Stmt]*CaseInfo) {
	caseByBody = map[ast.Stmt]*CaseInfo{}
	ast.PreorderStack(f, nil, func(n ast.Node, s []ast.Node) bool {
		switch node := n.(type) {
		case *ast.CaseClause:
			caseByBody[node] = &CaseInfo{Case: node}
		case *ast.CommClause:
			caseByBody[node] = &CaseInfo{Comm: node}
		}
		return true
	})
	return caseByBody
}

type IfElseType string

const (
	IfElseType_If     = "if"
	IfElseType_Else   = "else"
	IfElseType_ElseIf = "else if"
)

type IfElseInfo struct {
	Body     *ast.BlockStmt
	ElseBody *ast.BlockStmt
	Parents  []*ast.IfStmt
	IfStmt   *ast.IfStmt
}

func (i IfElseInfo) Variables() (vars []*ast.Ident) {
	for _, ifStmt := range append(append([]*ast.IfStmt{}, i.Parents...), i.IfStmt) {
		if ifStmt.Init != nil {
			if assign, ok := ifStmt.Init.(*ast.AssignStmt); ok {
				for _, lhs := range assign.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && ident.Name != "_" {
						vars = append(vars, ident)
					}
				}
			}
		}
	}
	return vars
}

func CollectIfElseInfo(f *ast.File) (ifElseByBody map[ast.Stmt]*IfElseInfo) {
	ifElseByBody = map[ast.Stmt]*IfElseInfo{}
	ast.PreorderStack(f, nil, func(n ast.Node, s []ast.Node) bool {
		switch node := n.(type) {
		case *ast.IfStmt:
			stack := []*ast.IfStmt{}
			for i := len(s) - 1; i >= 0; i-- {
				stmt, ok := s[i].(*ast.IfStmt)
				if !ok {
					break
				}
				stack = append(stack, stmt)
			}
			mutable.Reverse(stack)

			ifElseByBody[node.Body] = &IfElseInfo{
				Body:    node.Body,
				IfStmt:  node,
				Parents: stack,
			}
			if node.Else != nil {
				if blockBody, ok := node.Else.(*ast.BlockStmt); ok {
					ifElseByBody[blockBody] = &IfElseInfo{
						ElseBody: blockBody,
						IfStmt:   node,
						Parents:  stack,
					}
				}
			}
		}
		return true
	})
	return ifElseByBody
}
