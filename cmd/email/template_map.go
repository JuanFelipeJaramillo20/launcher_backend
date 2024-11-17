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
		baseDir, err := os.Getwd()
		if err != nil {
			log.Fatal("Failed to get current working directory:", err)
		}
		tmplPath = filepath.Join(baseDir, "cmd", "email", "templates")
	}
	log.Println("Loading templates from:", tmplPath)

	files, err := filepath.Glob(filepath.Join(tmplPath, "**/*.html"))
	if err != nil {
		log.Fatal("Failed to load email templates:", err)
	}

	if len(files) == 0 {
		log.Fatal("No templates found in:", tmplPath)
	}

	tmpls := make(map[string]*template.Template)
	for _, file := range files {
		relativePath := strings.ReplaceAll(file[len(tmplPath)+1:], `\`, `/`)
		tmpls[relativePath] = template.Must(template.ParseFiles(file))
		log.Printf("Loaded template: %s", relativePath)
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
