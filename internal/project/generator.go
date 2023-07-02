package project

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sivaprasadreddy/goquik/internal/tmpl"
)

func Generate(p *Project) error {
	genFile(p, p.ProjectName, "main.go.tmpl", "main.go")
	genFile(p, p.ProjectName, "go.mod.tmpl", "go.mod")
	return nil
}

func genFile(p *Project, filePath, templatePath, fileName string) {
	f := createFile(filePath, fileName)
	if f == nil {
		log.Printf("warn: file %s/%s %s", filePath, fileName, "already exists.")
		return
	}
	defer f.Close()

	t, err := template.ParseFS(tmpl.TemplateFS, templatePath)
	if err != nil {
		log.Fatalf("create %s error: %s", fileName, err.Error())
	}
	err = t.Execute(f, p)
	if err != nil {
		log.Fatalf("create %s error: %s", fileName, err.Error())
	}
	log.Printf("Created new %s: %s", fileName, f.Name())
}

func createFile(dirPath string, filename string) *os.File {
	filePath := filepath.Join(dirPath, filename)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create dir %s: %v", dirPath, err)
	}
	stat, _ := os.Stat(filePath)
	if stat != nil {
		return nil
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}
	return file
}
