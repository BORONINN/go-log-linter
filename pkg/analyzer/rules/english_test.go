package rules

import "testing"

func TestFindNonEnglish(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantBool bool
		wantRune rune
	}{
		// Clean English text.
		{"plain english", "starting server", false, 0},
		{"with digits", "port 8080", false, 0},
		{"with punctuation", "key: value, ok", false, 0},
		{"empty", "", false, 0},
		{"only digits", "12345", false, 0},
		{"latin diacritics", "café résumé", false, 0},

		// Non-English text.
		{"cyrillic", "запуск сервера", true, 'з'},
		{"cyrillic mixed", "start сервер end", true, 'с'},
		{"chinese", "你好世界", true, '你'},
		{"japanese katakana", "サーバー", true, 'サ'},
		{"arabic", "مرحبا", true, 'م'},
		{"korean", "서버 시작", true, '서'},
		{"greek", "αβγ", true, 'α'},
		{"single cyrillic", "a б c", true, 'б'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBool, gotRune := FindNonEnglish(tt.input)
			if gotBool != tt.wantBool {
				t.Errorf("FindNonEnglish(%q) found = %v, want %v",
					tt.input, gotBool, tt.wantBool)
			}
			if gotRune != tt.wantRune {
				t.Errorf("FindNonEnglish(%q) rune = %q, want %q",
					tt.input, gotRune, tt.wantRune)
			}
		})
	}
}
