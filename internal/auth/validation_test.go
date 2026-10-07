package auth

import "testing"

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		//{name: "", input: "", want: true},
		//{name: "", input: "", want: false},
		{name: "too_short", input: "shortpwd", want: false},
		{name: "long_enough", input: "thisisaverysecurepassword", want: true},
		{name: "space", input: "this is a password with spaces", want: false},
		{name: "empty", input: "", want: false},
		{name: "contains_emoji", input: "thisisnotonlyasciiasmypassword🫪", want: false},
		{name: "11_chars", input: "01234567891", want: false},
		{name: "12_chars", input: "012345678912", want: true},
		{name: "13_chars", input: "0123456789123", want: true},
		{name: "12_chars_ascii", input: "!this!-pwd!-", want: true},
		{name: "contains_umlaut", input: "thisisaverysecurepasswordäöü", want: false},
		{name: "contains_newline", input: "thisisaverysecurepassword\n", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidPassword(tt.input)

			if got != tt.want {
				t.Errorf("isValidPassword(%q) = %t, want %t",
					tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		//{name: "", input: "", want: true},
		//{name: "", input: "", want: false},
		{name: "valid_3_chars", input: "Jui", want: true},
		{name: "valid_18_chars", input: "juiisacoolverydevv", want: true},
		{name: "too_short_2_chars", input: "jj", want: false},
		{name: "too_long_19_chars", input: "juiisacooldeveloper", want: false},
		{name: "starts_with_dot", input: ".Jui", want: false},
		{name: "starts_with_underscore", input: "_Jui", want: false},
		{name: "valid_with_special_chars", input: "Jui-_-", want: true},
		{name: "empty", input: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidUsername(tt.input)

			if got != tt.want {
				t.Errorf("isValidUsername(%q) = %t, want %t",
					tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		//{name: "", input: "", want: true},
		//{name: "", input: "", want: false},
		{name: "valid", input: "help@me.com", want: true},
		{name: "empty", input: "", want: false},
		{name: "missing_at_sign", input: "helpme.com", want: false},
		{name: "name_with_email", input: "pete <help@me.com>", want: false},
		{name: "missing_local_address", input: "@me.com", want: false},
		{name: "missing_domain", input: "help@", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidEmail(tt.input)

			if got != tt.want {
				t.Errorf("isValidEmail(%q) = %t, want %t",
					tt.input, got, tt.want)
			}
		})
	}
}
