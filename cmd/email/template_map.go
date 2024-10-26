package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

var templates = loadTemplates()

func loadTemplates() map[string]*template.Template {
	tmpls := make(map[string]*template.Template)
	files, err := filepath.Glob("cmd/email/templates/**/*.html")
	if err != nil {
		log.Fatal("Failed to load email templates:", err)
	}
	for _, file := range files {
		name := strings.ReplaceAll(file[len("cmd/email/templates/"):], `\`, `/`)
		tmpls[name] = template.Must(template.ParseFiles(file))
	}
	return tmpls
}

func RenderTemplate(templatePath string, data interface{}) (string, error) {
	fmt.Println("TEMPLATES FOUND: ", len(templates), templates)
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
