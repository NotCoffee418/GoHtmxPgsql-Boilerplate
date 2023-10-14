package utils

import (
	"bytes"
	"fmt"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/access"
)

func RenderedTemplateString(templatePath string, data map[string]interface{}) (string, error) {
	var html bytes.Buffer
	if err := access.Tmpl.ExecuteTemplate(&html, templatePath, data); err != nil {
		return "", err
	}

	return html.String(), nil
}

func OobSwapWrap(id, html string) string {
	return fmt.Sprintf("<div hx-swap-oob=\"true\" id=\"%s\">%s</div>", id, html)
}
