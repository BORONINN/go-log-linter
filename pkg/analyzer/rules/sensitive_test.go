package rules

import "testing"

func TestContainsSensitiveKeyword(t *testing.T) {
	keywords := defaultKeywords

	tests := []struct {
		name     string
		input    string
		wantBool bool
		wantKW   string
	}{
		// Should detect.
		{"password in text", "user password: 123", true, "password"},
		{"token in text", "token: abc", true, "token"},
		{"api_key in text", "api_key=xyz", true, "api_key"},
		{"secret in text", "client_secret here", true, "secret"},
		{"credential", "invalid credential", true, "credential"},
		{"private_key", "loaded private_key", true, "private_key"},
		{"case insensitive", "User PASSWORD reset", true, "password"},
		{"passwd variant", "old passwd changed", true, "passwd"},
		{"access_token", "access_token expired", true, "token"},

		// Should not detect.
		{"clean message", "user authenticated", false, ""},
		{"empty", "", false, ""},
		{"server started", "server started on port 8080", false, ""},
		{"request completed", "api request completed", false, ""},
		{"normal log", "processing batch job", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKW, gotBool := ContainsSensitiveKeyword(tt.input, keywords)
			if gotBool != tt.wantBool {
				t.Errorf("ContainsSensitiveKeyword(%q) bool = %v, want %v",
					tt.input, gotBool, tt.wantBool)
			}
			if gotKW != tt.wantKW {
				t.Errorf("ContainsSensitiveKeyword(%q) keyword = %q, want %q",
					tt.input, gotKW, tt.wantKW)
			}
		})
	}
}

func TestContainsSensitiveKeywordCustom(t *testing.T) {
	custom := []string{"ssn", "credit_card"}

	tests := []struct {
		name     string
		input    string
		wantBool bool
	}{
		{"ssn match", "user ssn number", true},
		{"credit_card match", "credit_card: 1234", true},
		{"no match", "user name is john", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotBool := ContainsSensitiveKeyword(tt.input, custom)
			if gotBool != tt.wantBool {
				t.Errorf("ContainsSensitiveKeyword(%q, custom) = %v, want %v",
					tt.input, gotBool, tt.wantBool)
			}
		})
	}
}
