package rules

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

type SpecialCharsRule struct{}

func NewSpecialCharsRule() *SpecialCharsRule {
	return &SpecialCharsRule{}
}

func (r *SpecialCharsRule) Name() string {
	return "special-chars"
}

var forbiddenRunes = map[rune]bool{
	'!': true,
}

var forbiddenSubstrings = []string{
	"...",
}

func (r *SpecialCharsRule) Check(_ *analysis.Pass, call LogCall) []analysis.Diagnostic {
	if call.MessageLit == nil || call.MessageText == "" {
		return nil
	}

	var diagnostics []analysis.Diagnostic

	for _, sub := range forbiddenSubstrings {
		if strings.Contains(call.MessageText, sub) {
			diagnostics = append(diagnostics, analysis.Diagnostic{
				Pos:     call.MessageLit.Pos(),
				End:     call.MessageLit.End(),
				Message: fmt.Sprintf("log message should not contain %q", sub),
			})
		}
	}

	for _, char := range call.MessageText {
		if forbiddenRunes[char] {
			diagnostics = append(diagnostics, analysis.Diagnostic{
				Pos:     call.MessageLit.Pos(),
				End:     call.MessageLit.End(),
				Message: fmt.Sprintf("log message should not contain special character %q", char),
			})

			break
		}

		if IsEmoji(char) {
			diagnostics = append(diagnostics, analysis.Diagnostic{
				Pos:     call.MessageLit.Pos(),
				End:     call.MessageLit.End(),
				Message: fmt.Sprintf("log message should not contain emoji %q", char),
			})

			break
		}
	}

	return diagnostics
}

func IsEmoji(r rune) bool {
	if unicode.Is(unicode.So, r) {
		return true
	}

	switch {
	case r >= 0x1F600 && r <= 0x1F64F:
		return true
	case r >= 0x1F300 && r <= 0x1F5FF:
		return true
	case r >= 0x1F680 && r <= 0x1F6FF:
		return true
	case r >= 0x1F900 && r <= 0x1F9FF:
		return true
	case r >= 0x2600 && r <= 0x26FF:
		return true
	case r >= 0x2700 && r <= 0x27BF:
		return true
	}

	return false
}

func FindForbidden(s string) (bool, string) {
	for _, sub := range forbiddenSubstrings {
		if strings.Contains(s, sub) {
			return true, sub
		}
	}

	for _, r := range s {
		if forbiddenRunes[r] {
			return true, string(r)
		}

		if IsEmoji(r) {
			return true, string(r)
		}
	}

	return false, ""
}
