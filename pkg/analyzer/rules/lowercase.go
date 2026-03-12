package rules

import (
	"fmt"
	"go/token"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

type LowercaseRule struct{}

func NewLowercaseRule() *LowercaseRule {
	return &LowercaseRule{}
}

func (r *LowercaseRule) Name() string {
	return "lowercase"
}

func (r *LowercaseRule) Check(_ *analysis.Pass, call LogCall) []analysis.Diagnostic {
	if call.MessageLit == nil || call.MessageText == "" {
		return nil
	}

	if !StartsWithUppercase(call.MessageText) {
		return nil
	}

	firstRune, firstRuneSize := utf8.DecodeRuneInString(call.MessageText)
	lowerRune := unicode.ToLower(firstRune)
	lowerBytes := make([]byte, utf8.RuneLen(lowerRune))
	utf8.EncodeRune(lowerBytes, lowerRune)
	contentStart := call.MessageLit.Pos() + 1

	return []analysis.Diagnostic{
		{
			Pos:     call.MessageLit.Pos(),
			End:     call.MessageLit.End(),
			Message: "log message should start with a lowercase letter",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: fmt.Sprintf("replace %q with %q", firstRune, lowerRune),
					TextEdits: []analysis.TextEdit{
						{
							Pos:     contentStart,
							End:     contentStart + token.Pos(firstRuneSize),
							NewText: lowerBytes,
						},
					},
				},
			},
		},
	}
}

func StartsWithUppercase(s string) bool {
	if s == "" {
		return false
	}

	firstRune, _ := utf8.DecodeRuneInString(s)
	if firstRune == utf8.RuneError {
		return false
	}

	return unicode.IsLetter(firstRune) && unicode.IsUpper(firstRune)
}
