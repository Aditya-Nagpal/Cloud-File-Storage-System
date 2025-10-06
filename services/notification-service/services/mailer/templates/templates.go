package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
)

// TemplateCache holds parsed templates in memory
var templateCache = map[string]*template.Template{}

// LoadTemplates parses all .html files in this folder (once at startup)
func LoadTemplates() error {
	pattern := filepath.Join("services", "mailer", "templates", "*.html")
	tmpls, err := template.ParseGlob(pattern)
	if err != nil {
		return err
	}

	for _, t := range tmpls.Templates() {
		templateCache[t.Name()] = t
	}
	return nil
}

// RenderTemplate renders a template by name with provided data
func RenderTemplate(name string, data any) (string, error) {
	tmpl, ok := templateCache[name]
	if !ok {
		return "", fmt.Errorf("template not found %s", name)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
