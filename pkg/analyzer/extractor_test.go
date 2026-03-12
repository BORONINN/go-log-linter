package analyzer

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestFlattenConcat(t *testing.T) {
	// Build AST for: "a" + b + "c"
	// Parsed as: ("a" + b) + "c"
	litA := &ast.BasicLit{Kind: token.STRING, Value: `"a"`}
	identB := &ast.Ident{Name: "b"}
	litC := &ast.BasicLit{Kind: token.STRING, Value: `"c"`}

	expr := &ast.BinaryExpr{
		X: &ast.BinaryExpr{
			X:  litA,
			Op: token.ADD,
			Y:  identB,
		},
		Op: token.ADD,
		Y:  litC,
	}

	parts := flattenConcat(expr)

	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}

	// Verify order: litA, identB, litC.
	if parts[0] != litA {
		t.Errorf("parts[0]: expected litA, got %T", parts[0])
	}

	if parts[1] != identB {
		t.Errorf("parts[1]: expected identB, got %T", parts[1])
	}

	if parts[2] != litC {
		t.Errorf("parts[2]: expected litC, got %T", parts[2])
	}
}

func TestFlattenConcatSingleLiteral(t *testing.T) {
	lit := &ast.BasicLit{Kind: token.STRING, Value: `"hello"`}

	parts := flattenConcat(lit)

	if len(parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(parts))
	}

	if parts[0] != lit {
		t.Error("expected same literal back")
	}
}

func TestFlattenConcatDeeplyNested(t *testing.T) {
	// "a" + "b" + "c" + "d"
	// Parsed as: (("a" + "b") + "c") + "d"
	a := &ast.BasicLit{Kind: token.STRING, Value: `"a"`}
	b := &ast.BasicLit{Kind: token.STRING, Value: `"b"`}
	c := &ast.BasicLit{Kind: token.STRING, Value: `"c"`}
	d := &ast.BasicLit{Kind: token.STRING, Value: `"d"`}

	expr := &ast.BinaryExpr{
		X: &ast.BinaryExpr{
			X: &ast.BinaryExpr{
				X: a, Op: token.ADD, Y: b,
			},
			Op: token.ADD,
			Y:  c,
		},
		Op: token.ADD,
		Y:  d,
	}

	parts := flattenConcat(expr)

	if len(parts) != 4 {
		t.Fatalf("expected 4 parts, got %d", len(parts))
	}
}

func TestIsSupported(t *testing.T) {
	tests := []struct {
		pkg    string
		method string
		want   bool
	}{
		{"log/slog", "Info", true},
		{"log/slog", "Error", true},
		{"log/slog", "Debug", true},
		{"log/slog", "Warn", true},
		{"log/slog", "InfoContext", true},
		{"log/slog", "Printf", false},
		{"log/slog", "Print", false},

		{"go.uber.org/zap", "Info", true},
		{"go.uber.org/zap", "Error", true},
		{"go.uber.org/zap", "Infow", true},
		{"go.uber.org/zap", "Errorf", true},
		{"go.uber.org/zap", "Unknown", false},

		{"fmt", "Println", false},
		{"log", "Print", false},
		{"os", "Exit", false},
	}

	for _, tt := range tests {
		t.Run(tt.pkg+"."+tt.method, func(t *testing.T) {
			got := isSupported(tt.pkg, tt.method)
			if got != tt.want {
				t.Errorf("isSupported(%q, %q) = %v, want %v",
					tt.pkg, tt.method, got, tt.want)
			}
		})
	}
}
