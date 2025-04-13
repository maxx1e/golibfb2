// internal/archive/reader.go
//
// This module opens a 7z archive and provides a way to iterate through its files.
// For this demonstration, we use the "github.com/bodgit/sevenzip" library.
// In a real implementation, additional error handling and file filtering may be added.

package archive

import (
	"github.com/bodgit/sevenzip"
)

// OpenArchive opens a 7z archive specified by the given path and returns a slice of file pointers.
func OpenArchive(archivePath string) ([]*sevenzip.File, error) {
	archive, err := sevenzip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	// The caller is responsible for closing the archive after processing.
	return archive.File, nil
}
