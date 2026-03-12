package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "logcheck",
	Doc:      "checks that log messages start with a lowercase letter",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var slogMethods = map[string]bool{
	"Info":  true,
	"Error": true,
	"Debug": true,
	"Warn":  true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		if !isSlogCall(pass, call) {
			return
		}

		msg, lit := extractMessage(call)
		if lit == nil {
			return
		}

		checkLowercase(pass, msg, lit)
	})

	return nil, nil
}

func isSlogCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	methodName := sel.Sel.Name
	if !slogMethods[methodName] {
		return false
	}

	if ident, ok := sel.X.(*ast.Ident); ok {
		obj := pass.TypesInfo.Uses[ident]
		if obj == nil {
			return false
		}

		pkgName, ok := obj.(*types.PkgName)
		if !ok {
			return false
		}

		return pkgName.Imported().Path() == "log/slog"
	}

	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ == nil {
		return false
	}

	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == "log/slog"
}

func extractMessage(call *ast.CallExpr) (string, *ast.BasicLit) {
	if len(call.Args) == 0 {
		return "", nil
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", nil
	}

	text, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", nil
	}

	return text, lit
}

func checkLowercase(pass *analysis.Pass, msg string, lit *ast.BasicLit) {
	if msg == "" {
		return
	}

	firstRune, _ := utf8.DecodeRuneInString(msg)
	if firstRune == utf8.RuneError {
		return
	}

	if !unicode.IsLetter(firstRune) {
		return
	}

	if !unicode.IsUpper(firstRune) {
		return
	}

	pass.Reportf(lit.Pos(), "log message should start with a lowercase letter")
}
