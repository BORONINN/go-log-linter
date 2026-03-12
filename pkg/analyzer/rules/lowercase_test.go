package rules

import "testing"

func TestStartsWithUppercase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"uppercase ASCII", "Starting server", true},
		{"single uppercase letter", "A", true},
		{"uppercase unicode", "Über alles", true},

		{"lowercase ASCII", "starting server", false},
		{"single lowercase letter", "a", false},
		{"lowercase unicode", "über alles", false},
		{"starts with digit", "123 starting", false},
		{"starts with space", " starting", false},
		{"starts with punctuation", ".starting", false},
		{"starts with bracket", "[info] starting", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StartsWithUppercase(tt.input)
			if got != tt.want {
				t.Errorf("StartsWithUppercase(%q) = %v, want %v",
					tt.input, got, tt.want)
			}
		})
	}
}
