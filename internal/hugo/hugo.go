// internal/hugo/export.go
//
// This module exports book records from the database into Markdown files for Hugo.
// Each book will be saved as a content bundle (folder) containing an index.md file with YAML front matter.

package hugo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"database/sql"
	"golibfb2/internal/db"
	"golibfb2/internal/models"
)

// HugoBookFrontMatter represents the fields for the Markdown front matter.
type HugoBookFrontMatter struct {
	Title      string
	Authors    []string
	Genres     []string
	Language   string
	Annotation string
	FileName   string
	Cover      string
}

const markdownTemplate = ` + "`" + `---
title: "{{ .Title }}"
authors: [{{ range $index, $element := .Authors }}{{if $index}}, {{end}}"{{ $element }}"{{end}}]
genres: [{{ range $index, $element := .Genres }}{{if $index}}, {{end}}"{{ $element }}"{{end}}]
language: "{{ .Language }}"
annotation: "{{ .Annotation }}"
cover: "{{ .Cover }}"
file_name: "{{ .FileName }}"
draft: false
---

{{ if .Cover }}
![Cover image]({{ .Cover }})
{{ end }}

{{ .Annotation }}
` + "`" + `

// ExportToHugo exports all books from the database into a Hugo content directory.
func ExportToHugo(outputDir string, dbConn *sql.DB) error {
	books, err := db.GetAllBooks(dbConn)
	if err != nil {
		return fmt.Errorf("failed to retrieve books: %w", err)
	}

	// Define the Hugo content directory (e.g., content/books/).
	hugoBooksDir := filepath.Join(outputDir, "content", "books")
	if err := os.MkdirAll(hugoBooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create content directory: %w", err)
	}

	// Parse the Markdown template.
	tmpl, err := template.New("book").Parse(markdownTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Generate a folder and index.md file for each book.
	for _, book := range books {
		// Generate a safe folder name.
		safeTitle := safeFileName(book.Title)
		bookDir := filepath.Join(hugoBooksDir, fmt.Sprintf("%s_%d", safeTitle, book.ID))
		if err := os.MkdirAll(bookDir, 0755); err != nil {
			return fmt.Errorf("failed to create book directory: %w", err)
		}

		frontMatter := HugoBookFrontMatter{
			Title:      book.Title,
			Authors:    book.Authors, // Populate if available.
			Genres:     book.Genres,  // Populate if available.
			Language:   book.Language,
			Annotation: book.Annotation,
			FileName:   book.FileName,
			Cover:      "", // Set a cover filename if images are handled.
		}

		mdPath := filepath.Join(bookDir, "index.md")
		f, err := os.Create(mdPath)
		if err != nil {
			return fmt.Errorf("failed to create markdown file: %w", err)
		}
		// Ensure the file is closed after writing.
		if err := tmpl.Execute(f, frontMatter); err != nil {
			f.Close()
			return fmt.Errorf("failed to execute template: %w", err)
		}
		f.Close()
	}
	return nil
}

// safeFileName returns a sanitized version of the string for folder names.
func safeFileName(name string) string {
	// Replace spaces with hyphens and convert to lower case.
	safe := strings.ReplaceAll(name, " ", "-")
	return strings.ToLower(safe)
}
