package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

var supportedPackages = map[string]map[string]bool{
	"log/slog": {
		"Info":         true,
		"Error":        true,
		"Debug":        true,
		"Warn":         true,
		"InfoContext":  true,
		"ErrorContext": true,
		"DebugContext": true,
		"WarnContext":  true,
	},
	"go.uber.org/zap": {
		// *zap.Logger methods
		"Info":   true,
		"Error":  true,
		"Debug":  true,
		"Warn":   true,
		"Fatal":  true,
		"Panic":  true,
		"DPanic": true,
		// *zap.SugaredLogger methods
		"Infow":  true,
		"Errorw": true,
		"Debugw": true,
		"Warnw":  true,
		"Fatalw": true,
		"Panicw": true,
		"Infof":  true,
		"Errorf": true,
		"Debugf": true,
		"Warnf":  true,
		"Fatalf": true,
		"Panicf": true,
	},
}

func extractLogCalls(pass *analysis.Pass, insp *inspector.Inspector) []LogCall {
	var calls []LogCall

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		funcName, ok := matchLogCall(pass, call)
		if !ok {
			return
		}

		lc := buildLogCall(call, funcName)
		if lc != nil {
			calls = append(calls, *lc)
		}
	})

	return calls
}

func matchLogCall(pass *analysis.Pass, call *ast.CallExpr) (string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	methodName := sel.Sel.Name

	if ident, ok := sel.X.(*ast.Ident); ok {
		if obj := pass.TypesInfo.Uses[ident]; obj != nil {
			if pkgName, ok := obj.(*types.PkgName); ok {
				pkgPath := pkgName.Imported().Path()
				if isSupported(pkgPath, methodName) {
					return methodName, true
				}
			}
		}
	}

	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ == nil {
		return "", false
	}

	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	named, ok := typ.(*types.Named)
	if !ok {
		return "", false
	}

	obj := named.Obj()
	if obj.Pkg() == nil {
		return "", false
	}

	pkgPath := obj.Pkg().Path()
	if isSupported(pkgPath, methodName) {
		return methodName, true
	}

	return "", false
}

func isSupported(pkgPath, methodName string) bool {
	methods, ok := supportedPackages[pkgPath]
	if !ok {
		return false
	}

	return methods[methodName]
}

func buildLogCall(call *ast.CallExpr, funcName string) *LogCall {
	if len(call.Args) == 0 {
		return nil
	}

	firstArg := call.Args[0]

	lc := &LogCall{
		Pos:      call.Pos(),
		End:      call.End(),
		FuncName: funcName,
		RawArg:   firstArg,
	}

	switch arg := firstArg.(type) {
	case *ast.BasicLit:
		if arg.Kind == token.STRING {
			lc.MessageLit = arg

			text, err := strconv.Unquote(arg.Value)
			if err == nil {
				lc.MessageText = text
			}
		}

	case *ast.BinaryExpr:
		if arg.Op == token.ADD {
			lc.HasConcat = true
			lc.ConcatParts = flattenConcat(arg)

			for _, part := range lc.ConcatParts {
				lit, ok := part.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}

				text, err := strconv.Unquote(lit.Value)
				if err != nil {
					continue
				}

				lc.MessageLit = lit
				lc.MessageText = text

				break
			}
		}
	}

	return lc
}

func flattenConcat(expr ast.Expr) []ast.Expr {
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok || bin.Op != token.ADD {
		return []ast.Expr{expr}
	}

	var parts []ast.Expr
	parts = append(parts, flattenConcat(bin.X)...)
	parts = append(parts, flattenConcat(bin.Y)...)

	return parts
}
