// internal/parser/fb2.go
//
// This module handles parsing of FB2 files using Go's encoding/xml package.
// It defines simplified FB2 structures necessary for extracting metadata.
// In production, you might extend this to handle more elements and embedded binaries.

package parser

import (
	"encoding/xml"
	"errors"
	"strings"

	"github.com/maxx1e/golibfb2/internal/models"
)

// FictionBook represents a simplified structure of the FB2 file.
type FictionBook struct {
	XMLName     xml.Name    `xml:"FictionBook"`
	Description Description `xml:"description"`
}

type Description struct {
	TitleInfo TitleInfo `xml:"title-info"`
}

type TitleInfo struct {
	BookTitle  string     `xml:"book-title"`
	Authors    []Author   `xml:"author"`
	Genres     []string   `xml:"genre"`
	Lang       string     `xml:"lang"`
	Annotation Annotation `xml:"annotation"`
	Coverpage  Coverpage  `xml:"coverpage"`
}

type Author struct {
	FirstName  string `xml:"first-name"`
	MiddleName string `xml:"middle-name"`
	LastName   string `xml:"last-name"`
}

type Annotation struct {
	// For demonstration, we capture the inner XML.
	Text string `xml:",innerxml"`
}

type Coverpage struct {
	Image Image `xml:"image"`
}

type Image struct {
	// In FB2, cover image is referenced by an attribute (typically starting with "#").
	Href string `xml:"href,attr"`
}

// ParseFB2 unmarshals FB2 XML data into our Book model.
func ParseFB2(data []byte) (*models.Book, error) {
	var fb FictionBook
	if err := xml.Unmarshal(data, &fb); err != nil {
		return nil, err
	}

	ti := fb.Description.TitleInfo
	// Check that a mandatory field is present.
	if strings.TrimSpace(ti.BookTitle) == "" {
		return nil, errors.New("missing book title")
	}

	// Compose the full author names.
	var authors []string
	for _, a := range ti.Authors {
		// Concatenate non‑empty parts of the name.
		var fullName []string
		if a.FirstName != "" {
			fullName = append(fullName, a.FirstName)
		}
		if a.MiddleName != "" {
			fullName = append(fullName, a.MiddleName)
		}
		if a.LastName != "" {
			fullName = append(fullName, a.LastName)
		}
		authors = append(authors, strings.Join(fullName, " "))
	}

	annotation := strings.TrimSpace(ti.Annotation.Text)
	// In a complete solution, you would handle embedded cover images here.
	book := models.Book{
		Title:      ti.BookTitle,
		Authors:    authors,
		Genres:     ti.Genres,
		Language:   ti.Lang,
		Annotation: annotation,
		// CoverData and other fields can be set later.
	}
	return &book, nil
}
