package analyzer

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/BORONINN/go-log-linter/pkg/analyzer/rules"
)

var Analyzer = &analysis.Analyzer{
	Name:     "logcheck",
	Doc:      "checks log messages for common style and security issues",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	flagLowercase     bool
	flagEnglish       bool
	flagSpecialChars  bool
	flagSensitive     bool
	flagExtraKeywords string
)

func init() {
	Analyzer.Flags.BoolVar(&flagLowercase, "lowercase", true,
		"check that log messages start with a lowercase letter")
	Analyzer.Flags.BoolVar(&flagEnglish, "english", true,
		"check that log messages are in English")
	Analyzer.Flags.BoolVar(&flagSpecialChars, "special-chars", true,
		"check for special characters and emoji in log messages")
	Analyzer.Flags.BoolVar(&flagSensitive, "sensitive", true,
		"check for sensitive data in log messages")
	Analyzer.Flags.StringVar(&flagExtraKeywords, "sensitive-keywords", "",
		"comma-separated list of additional sensitive data keywords")
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	activeRules := buildRules()
	if len(activeRules) == 0 {
		return nil, nil
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

func buildRules() []rules.Rule {
	var active []rules.Rule

	if flagLowercase {
		active = append(active, rules.NewLowercaseRule())
	}

	if flagEnglish {
		active = append(active, rules.NewEnglishRule())
	}

	if flagSpecialChars {
		active = append(active, rules.NewSpecialCharsRule())
	}

	if flagSensitive {
		var extraKW []string
		if flagExtraKeywords != "" {
			for _, kw := range strings.Split(flagExtraKeywords, ",") {
				trimmed := strings.TrimSpace(kw)
				if trimmed != "" {
					extraKW = append(extraKW, trimmed)
				}
			}
		}

		active = append(active, rules.NewSensitiveRule(extraKW))
	}

	return active
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
