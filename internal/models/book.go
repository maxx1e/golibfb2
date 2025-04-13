// internal/models/book.go
//
// This file defines the core data model for a book. It holds all metadata
// extracted from an FB2 file. Additional fields (e.g., cover image data) can be added as needed.

package models

// Book represents a book's metadata extracted from an FB2 file.
type Book struct {
	// ID is populated by the database.
	ID int
	// Title of the book.
	Title string
	// Authors holds the list of authors.
	Authors []string
	// Genres holds the list of genres/categories.
	Genres []string
	// Language code (e.g., "en", "ru").
	Language string
	// Annotation or summary of the book.
	Annotation string
	// CoverData holds raw binary data for the cover image.
	CoverData []byte
	// FileName is the original filename (used to prevent re‑imports).
	FileName string
	// Tags for the better search capabilities
	Tags []string
	// Series is the name of the series the book belongs to.
	Series string
	// SeriesNumber is the number of the book in the series.
	SeriesNumber int
	// Book size in bytes.
	Size int64
	// Path is the path to the book file.
	ArchivePath string
	// DateAdded is the date when the book was added to the library.
	DateAdded string
}
