package internal

import (
	"go/ast"
	"go/token"
	"strings"
)

type Xtrace struct {
	Config
	fset *token.FileSet
	src  []byte

	funcByBody   map[ast.Stmt]*FuncInfo
	forByBody    map[ast.Stmt]*ForInfo
	caseByBody   map[ast.Stmt]*CaseInfo
	ifElseByBody map[ast.Stmt]*IfElseInfo

	libraryRequired bool
}

func (x *Xtrace) fragment(pos, end token.Pos) string {
	return string(x.src[pos-1 : end-1])
}
func (x *Xtrace) fragmentLine(pos token.Pos) string {
	begin := pos - 1
	for ; begin > 0; begin-- {
		if x.src[begin-1] == '\n' || x.src[begin-1] == '\r' {
			break
		}
	}
	end := pos
	for ; end < token.Pos(len(x.src)); end++ {
		if x.src[end] == '\n' || x.src[end] == '\r' {
			break
		}
	}
	frag := x.fragment(begin+1, end+1)
	frag, _, _ = strings.Cut(frag, "\n")
	frag, _, _ = strings.Cut(frag, "\r")
	return frag
}

func (x *Xtrace) IdentShowTimestamp() string {
	if x.ShowTimestamp {
		return "true"
	}
	return "false"
}

func (x *Xtrace) IdentShowGoroutine() string {
	if x.ShowGoroutine {
		return "true"
	}
	return "false"
}
