package analyzer

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"customLinterSelectel/pkg/analyzer/rules"
)

var Analyzer = &analysis.Analyzer{
	Name:     "logcheck",
	Doc:      "checks log messages for common style and security issues",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	activeRules := []rules.Rule{
		rules.NewLowercaseRule(),
		rules.NewEnglishRule(),
		rules.NewSpecialCharsRule(),
		rules.NewSensitiveRule(nil),
	}

	logCalls := extractLogCalls(pass, insp)

	for i := range logCalls {
		rc := toRulesLogCall(logCalls[i])

		for _, rule := range activeRules {
			diagnostics := rule.Check(pass, rc)

			for _, d := range diagnostics {
				pass.Report(d)
			}
		}
	}

	return nil, nil
}

func toRulesLogCall(lc LogCall) rules.LogCall {
	return rules.LogCall{
		Pos:         lc.Pos,
		End:         lc.End,
		FuncName:    lc.FuncName,
		MessageLit:  lc.MessageLit,
		MessageText: lc.MessageText,
		RawArg:      lc.RawArg,
		HasConcat:   lc.HasConcat,
		ConcatParts: lc.ConcatParts,
	}
}
