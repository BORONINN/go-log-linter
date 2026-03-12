package rules

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type LogCall struct {
	Pos         token.Pos
	End         token.Pos
	FuncName    string
	MessageLit  *ast.BasicLit
	MessageText string
	RawArg      ast.Expr
	HasConcat   bool
	ConcatParts []ast.Expr
}

type Rule interface {
	Name() string

	Check(pass *analysis.Pass, call LogCall) []analysis.Diagnostic
}
