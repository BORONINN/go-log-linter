package rules

import "testing"

func TestIsEmoji(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		// Emoji that should be detected.
		{"rocket", '🚀', true},
		{"grinning face", '😀', true},
		{"fire", '🔥', true},
		{"heart", '❤', true},
		{"check mark", '✓', true},
		{"star", '★', true},
		{"sun", '☀', true},
		{"warning sign", '⚠', true},

		// Normal characters that should pass.
		{"letter a", 'a', false},
		{"letter Z", 'Z', false},
		{"digit 0", '0', false},
		{"space", ' ', false},
		{"period", '.', false},
		{"comma", ',', false},
		{"colon", ':', false},
		{"hyphen", '-', false},
		{"underscore", '_', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEmoji(tt.r)
			if got != tt.want {
				t.Errorf("IsEmoji(%q U+%04X) = %v, want %v",
					tt.r, tt.r, got, tt.want)
			}
		})
	}
}

func TestFindForbidden(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantFound bool
	}{
		// Clean messages.
		{"normal message", "server started", false},
		{"with colon", "key: value", false},
		{"with comma", "a, b, c", false},
		{"with single dot", "v1.2.3", false},
		{"with two dots", "file..bak", false},
		{"empty", "", false},
		{"with hyphen", "my-service", false},
		{"with parens", "count (total)", false},

		// Forbidden content.
		{"exclamation mark", "server started!", true},
		{"multiple exclamations", "failed!!!", true},
		{"emoji rocket", "deployed 🚀", true},
		{"emoji fire", "hot 🔥", true},
		{"ellipsis", "loading...", true},
		{"ellipsis in middle", "wait... done", true},
		{"exclamation in middle", "oh! no", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFound, _ := FindForbidden(tt.input)
			if gotFound != tt.wantFound {
				t.Errorf("FindForbidden(%q) = %v, want %v",
					tt.input, gotFound, tt.wantFound)
			}
		})
	}
}
