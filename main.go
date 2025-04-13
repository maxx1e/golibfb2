// cmd/golibfb2/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// Import internal packages using their full module paths.
	"github.com/yourusername/golibfb2/internal/db"
	"github.com/yourusername/golibfb2/internal/hugo"
	"github.com/yourusername/golibfb2/internal/pipeline"
)

func main() {
	importCmd := flag.NewFlagSet("import", flag.ExitOnError)
	importArchive := importCmd.String("archive", "", "Path to the 7z archive containing FB2 files")

	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	exportDir := exportCmd.String("out", "", "Output directory for Hugo export")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'import' or 'export' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "import":
		importCmd.Parse(os.Args[2:])
		if *importArchive == "" {
			log.Fatal("The -archive flag is required for import")
		}
		dbConn, err := db.ConnectDB()
		if err != nil {
			log.Fatalf("Database connection error: %v", err)
		}
		defer dbConn.Close()

		if err := pipeline.ProcessArchive(*importArchive, dbConn); err != nil {
			log.Fatalf("Import failed: %v", err)
		}
		fmt.Println("Import succeeded.")

	case "export":
		exportCmd.Parse(os.Args[2:])
		if *exportDir == "" {
			log.Fatal("The -out flag is required for export")
		}
		dbConn, err := db.ConnectDB()
		if err != nil {
			log.Fatalf("Database connection error: %v", err)
		}
		defer dbConn.Close()

		if err := hugo.ExportToHugo(*exportDir, dbConn); err != nil {
			log.Fatalf("Export failed: %v", err)
		}
		fmt.Println("Export succeeded.")

	default:
		fmt.Println("Expected 'import' or 'export' subcommands")
		os.Exit(1)
	}
}
