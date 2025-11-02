package xtracego

import (
	"go/ast"
	"go/token"
	"strings"
)

type XTrace struct {
	fset      *token.FileSet
	src       []byte
	lineWidth int
	prefix    string

	funcByBody   map[ast.Stmt]*FuncInfo
	forByBody    map[ast.Stmt]*ForInfo
	caseByBody   map[ast.Stmt]*CaseInfo
	ifElseByBody map[ast.Stmt]*IfElseInfo
}

func (s XTrace) fragment(pos, end token.Pos) string {
	return string(s.src[pos-1 : end-1])
}
func (s XTrace) fragmentLine(pos, end token.Pos) string {
	begin := pos - 1
	for ; begin > 0; begin-- {
		if s.src[begin-1] == '\n' || s.src[begin-1] == '\r' {
			break
		}
	}
	frag := s.fragment(begin+1, end)
	frag, _, _ = strings.Cut(frag, "\n")
	frag, _, _ = strings.Cut(frag, "\r")
	return frag
}
