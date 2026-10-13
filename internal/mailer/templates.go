package mailer

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed templates/password_reset.html
var passwordResetHTML string

func RenderPasswordReset(resetURL string) (string, error) {
	tmpl, err := template.New("password_reset").Parse(passwordResetHTML)
	if err != nil {
		return "", err
	}
	data := struct {
		ResetURL string
	}{
		ResetURL: resetURL,
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
