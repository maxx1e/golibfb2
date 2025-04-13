// internal/pipeline/pipeline.go
//
// This module orchestrates the processing of the FB2 files.
// It reads files from the archive, parses them, and writes the data into the database.
// The work is distributed across multiple worker goroutines for performance.

package pipeline

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	"database/sql"
	"golibfb2/internal/db"
	"golibfb2/internal/parser"

	"github.com/bodgit/sevenzip"
)

// ProcessArchive opens the 7z archive and processes FB2 files concurrently.
// It extracts each file, parses it into a Book, and inserts the record into the database.
func ProcessArchive(archivePath string, dbConn *sql.DB) error {
	// Open the 7z archive.
	archiveReader, err := sevenzip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer archiveReader.Close()

	// Create a channel for sending file pointers.
	fileChan := make(chan *sevenzip.File)

	// Use a WaitGroup to wait for all worker goroutines to finish.
	var wg sync.WaitGroup

	// Set the number of worker goroutines.
	numWorkers := 4 // For real use, consider runtime.NumCPU() for dynamic adjustment.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for f := range fileChan {
				// Process only files with a .fb2 extension.
				if filepath.Ext(f.Name) != ".fb2" {
					continue
				}

				// Open the file inside the archive.
				rc, err := f.Open()
				if err != nil {
					log.Printf("Worker %d: failed to open %s: %v", workerID, f.Name, err)
					continue
				}
				data, err := ioutil.ReadAll(rc)
				rc.Close() // Important: close the reader to free resources.
				if err != nil {
					log.Printf("Worker %d: failed to read %s: %v", workerID, f.Name, err)
					continue
				}

				// Parse the FB2 XML data.
				book, err := parser.ParseFB2(data)
				if err != nil {
					log.Printf("Worker %d: parse error for %s: %v", workerID, f.Name, err)
					continue
				}
				// Set the file name to avoid duplicate imports.
				book.FileName = f.Name

				// Insert the book into the database.
				if err := db.InsertBook(dbConn, book); err != nil {
					log.Printf("Worker %d: DB insert error for %s: %v", workerID, f.Name, err)
					continue
				}
				log.Printf("Worker %d: processed %s successfully.", workerID, f.Name)
			}
		}(i)
	}

	// Producer: iterate over each file in the archive and send it to the channel.
	for _, f := range archiveReader.File {
		fileChan <- f
	}
	// Close the channel to signal workers no more files will be sent.
	close(fileChan)

	// Wait until all workers finish.
	wg.Wait()
	return nil
}
