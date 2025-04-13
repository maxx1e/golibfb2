package db

import (
	"database/sql"
	"log"
	"strings"

	"golibfb2/models"

	_ "github.com/mattn/go-sqlite3" // SQLite driver registration.
)

// ConnectDB opens a connection to the SQLite database file "gofb2.db".
// If the file doesn't exist, it will be created.
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./gofb2.db")
	if err != nil {
		return nil, err
	}
	// Create the schema if it doesn't already exist.
	if err := createSchema(db); err != nil {
		return nil, err
	}
	return db, nil
}

// createSchema creates the required table for storing books.
func createSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS books (
	book_id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	authors TEXT,
	genres TEXT,
	language TEXT,
	annotation TEXT,
	cover_data BLOB,
	file_name TEXT UNIQUE,
	tags TEXT,
	series TEXT,
	series_number INTEGER,
	size INTEGER,
	archive_path TEXT,
	date_added TEXT
);
`
	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("Error creating schema: %v", err)
	}
	return err
}

// InsertBook inserts a book into the database or updates it if it already exists.
// SQLite's ON CONFLICT clause is used to handle duplicates based on file_name.
func InsertBook(db *sql.DB, book *models.Book) error {
	query := `
INSERT INTO books (
	title, authors, genres, language, annotation, cover_data, file_name, tags, 
	series, series_number, size, archive_path, date_added
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(file_name) DO UPDATE SET
	title = excluded.title,
	authors = excluded.authors,
	genres = excluded.genres,
	language = excluded.language,
	annotation = excluded.annotation,
	cover_data = excluded.cover_data,
	tags = excluded.tags,
	series = excluded.series,
	series_number = excluded.series_number,
	size = excluded.size,
	archive_path = excluded.archive_path,
	date_added = excluded.date_added;
`
	_, err := db.Exec(query,
		book.Title,
		serializeStringSlice(book.Authors),
		serializeStringSlice(book.Genres),
		book.Language,
		book.Annotation,
		book.CoverData,
		book.FileName,
		serializeStringSlice(book.Tags),
		book.Series,
		book.SeriesNumber,
		book.Size,
		book.ArchivePath,
		book.DateAdded,
	)
	return err
}

// GetAllBooks retrieves all book records from the database.
func GetAllBooks(db *sql.DB) ([]models.Book, error) {
	rows, err := db.Query(`
SELECT 
	book_id, title, authors, genres, language, annotation, cover_data, file_name, 
	tags, series, series_number, size, archive_path, date_added 
FROM books
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		var authors, genres, tags string
		if err := rows.Scan(
			&b.ID, &b.Title, &authors, &genres, &b.Language, &b.Annotation, &b.CoverData,
			&b.FileName, &tags, &b.Series, &b.SeriesNumber, &b.Size, &b.ArchivePath, &b.DateAdded,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		b.Authors = deserializeStringSlice(authors)
		b.Genres = deserializeStringSlice(genres)
		b.Tags = deserializeStringSlice(tags)
		books = append(books, b)
	}
	return books, nil
}

// Helper functions to serialize and deserialize string slices for storage in the database.
func serializeStringSlice(slice []string) string {
	return strings.Join(slice, ",")
}

func deserializeStringSlice(data string) []string {
	if data == "" {
		return []string{}
	}
	return strings.Split(data, ",")
}
