package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var templates = loadTemplates()

func loadTemplates() map[string]*template.Template {
	tmplPath, exists := os.LookupEnv("TEMPLATE_PATH")
	if !exists {
		tmplPath = "cmd/email/templates"
	}

	files, err := filepath.Glob(filepath.Join(tmplPath, "**/*.html"))
	if err != nil {
		log.Fatal("Failed to load email templates:", err)
	}

	tmpls := make(map[string]*template.Template)
	for _, file := range files {
		name := strings.ReplaceAll(file[len(tmplPath):], `\`, `/`)
		tmpls[name] = template.Must(template.ParseFiles(file))
	}

	return tmpls
}

func RenderTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, ok := templates[templatePath]
	if !ok {
		return "", fmt.Errorf("template %s not found", templatePath)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
