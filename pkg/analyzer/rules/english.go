package rules

import (
	"fmt"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

type EnglishRule struct{}

func NewEnglishRule() *EnglishRule {
	return &EnglishRule{}
}

func (r *EnglishRule) Name() string {
	return "english"
}

func (r *EnglishRule) Check(_ *analysis.Pass, call LogCall) []analysis.Diagnostic {
	if call.MessageLit == nil || call.MessageText == "" {
		return nil
	}

	found, char := FindNonEnglish(call.MessageText)
	if !found {
		return nil
	}

	return []analysis.Diagnostic{
		{
			Pos:     call.MessageLit.Pos(),
			End:     call.MessageLit.End(),
			Message: fmt.Sprintf("log message should be in English, found non-Latin character %q", char),
		},
	}
}

func FindNonEnglish(s string) (bool, rune) {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.In(r, unicode.Latin) {
			return true, r
		}
	}

	return false, 0
}
