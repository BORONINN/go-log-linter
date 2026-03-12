package analyzer

import (
	"go/ast"
	"go/token"
)

type LogCall struct {
	Pos token.Pos

	End token.Pos

	FuncName string

	MessageLit *ast.BasicLit

	MessageText string

	RawArg ast.Expr

	HasConcat bool

	ConcatParts []ast.Expr
}
