package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var defaultKeywords = []string{
	"password",
	"passwd",
	"pwd",
	"token",
	"secret",
	"api_key",
	"apikey",
	"api-key",
	"private_key",
	"privatekey",
	"credential",
	"credentials",
	"auth_token",
	"access_token",
	"refresh_token",
}

type SensitiveRule struct {
	keywords []string
}

func NewSensitiveRule(customKeywords []string) *SensitiveRule {
	kw := make([]string, len(defaultKeywords))
	copy(kw, defaultKeywords)
	kw = append(kw, customKeywords...)

	return &SensitiveRule{keywords: kw}
}

func (r *SensitiveRule) Name() string {
	return "sensitive"
}

func (r *SensitiveRule) Check(_ *analysis.Pass, call LogCall) []analysis.Diagnostic {
	if call.MessageLit != nil && call.MessageText != "" {
		if kw, found := r.findKeyword(call.MessageText); found {
			return []analysis.Diagnostic{
				{
					Pos:     call.MessageLit.Pos(),
					End:     call.MessageLit.End(),
					Message: fmt.Sprintf("log message may contain sensitive data (keyword %q)", kw),
				},
			}
		}
	}

	if call.HasConcat {
		return r.checkConcat(call)
	}

	return nil
}

func (r *SensitiveRule) checkConcat(call LogCall) []analysis.Diagnostic {
	for _, part := range call.ConcatParts {
		switch p := part.(type) {
		case *ast.BasicLit:
			if p.Kind != token.STRING {
				continue
			}

			text, err := strconv.Unquote(p.Value)
			if err != nil {
				continue
			}

			if kw, found := r.findKeyword(text); found {
				return []analysis.Diagnostic{
					{
						Pos:     p.Pos(),
						End:     p.End(),
						Message: fmt.Sprintf("log message may contain sensitive data (keyword %q)", kw),
					},
				}
			}

		case *ast.Ident:
			if kw, found := r.findKeyword(p.Name); found {
				return []analysis.Diagnostic{
					{
						Pos:     call.RawArg.Pos(),
						End:     call.RawArg.End(),
						Message: fmt.Sprintf("log message may contain sensitive data (variable %q matches keyword %q)", p.Name, kw),
					},
				}
			}
		}
	}

	return nil
}

func (r *SensitiveRule) findKeyword(text string) (string, bool) {
	lower := strings.ToLower(text)

	for _, kw := range r.keywords {
		if strings.Contains(lower, kw) {
			return kw, true
		}
	}

	return "", false
}

func ContainsSensitiveKeyword(text string, keywords []string) (string, bool) {
	lower := strings.ToLower(text)

	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return kw, true
		}
	}

	return "", false
}
