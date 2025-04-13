# GoLibFB2

GoLibFB2 is a scalable, high-performance Go project designed to process hundreds of thousands of FB2 e-book files stored in 7z archives. The project extracts key metadata (such as title, authors, genres, language, and annotation) along with cover images, stores the data in an SQL database (SQLite by default), and exports it to a Hugo-compatible static site for quick and efficient browsing and full-text search.

> **Note:** This project uses Go's native concurrency primitives (goroutines and channels) to achieve highly-parallel file processing. It is built with robust error handling for malformed files, supports incremental updates, and is organized into modular components for easy maintenance and extension.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation and Setup](#installation-and-setup)
- [Usage](#usage)
  - [Importing FB2 Files](#importing-fb2-files)
  - [Exporting to Hugo](#exporting-to-hugo)
- [Testing](#testing)
- [Customization and Extensions](#customization-and-extensions)
- [License](#license)

## Overview

GoLibFB2 is designed to meet the following objectives:
- **Bulk Processing:** Efficiently read and process large 7z archives containing hundreds of thousands of FB2 files.
- **Metadata Extraction:** Extract book title, authors, genre(s), language, annotation, and cover image data from FB2 XML files.
- **Concurrent Processing:** Use a multi-stage processing pipeline with goroutines and channels to handle a large number of files concurrently.
- **Robust Error Handling:** Log and skip malformed or incomplete FB2 files without interrupting overall processing.
- **Database Storage:** Store the processed data in an embedded SQLite database (or optionally, PostgreSQL) using upsert semantics for incremental updates.
- **Hugo Static Site Export:** Convert stored book records into a Hugo-compatible structure, enabling fast static-site generation and client-side search.
- **Test Coverage:** Include unit tests to ensure the FB2 parsing, database operations, and pipeline processing work as expected.

## Features

- **FB2 File Processing:** Scans and extracts metadata from FB2 files contained in a 7z archive.
- **Concurrent Pipeline:** Leverages a producer/consumer model using channels to process files in parallel.
- **SQL Database Storage:** Uses SQLite (via the `go-sqlite3` driver) to store data locally without needing a separate database installation.
- **Static Site Generation:** Exports book records as Markdown files with YAML front matter for integration with Hugo.
- **Error Logging:** Robust logging to ensure that errors in individual files do not impact overall processing.
- **Unit and Integration Tests:** A test suite to validate core functionalities like XML parsing and database operations.

## Architecture

The application is divided into several modules with clear responsibilities:

- **CLI Application (`cmd/golibfb2`):**  
  The main entry point that handles command-line arguments and dispatches either the import or export command.

- **Archive Reader (`internal/archive`):**  
  Wraps the 7z archive reading (using [bodgit/sevenzip](https://github.com/bodgit/sevenzip)) to provide a stream of FB2 files.

- **FB2 Parser (`internal/parser`):**  
  Contains logic to unmarshal FB2 XML into a Go data structure, extracting key metadata and cover image references.

- **Data Models (`internal/models`):**  
  Defines core data types such as `Book` to hold metadata extracted from FB2 files.

- **Database Layer (`internal/db`):**  
  Manages database connection, schema creation, and CRUD operations. By default, it uses SQLite.

- **Hugo Exporter (`internal/hugo`):**  
  Converts database records into Hugo content pages (Markdown files with YAML front matter) suitable for generating a static site.

- **Processing Pipeline (`internal/pipeline`):**  
  Orchestrates the reading, parsing, and database insertion processes concurrently using goroutines and channels.

## Project Structure

golibfb2/ ├── cmd/ │ └── golibfb2/ │ └── main.go # CLI entry point: handles import/export commands. ├── internal/ │ ├── archive/ │ │ └── reader.go # Reads 7z archives and provides file iteration. │ ├── parser/ │ │ └── fb2.go # Parses FB2 files to extract metadata. │ ├── models/ │ │ └── book.go # Defines core data types (e.g. Book). │ ├── db/ │ │ └── db.go # Manages DB connections, schema, and queries (SQLite by default). │ ├── hugo/ │ │ └── export.go # Exports data into Hugo Markdown files. │ └── pipeline/ │ └── pipeline.go # Coordinates file processing using goroutines and channels. ├── testdata/ # Contains sample FB2 and 7z files for testing. ├── go.mod # Go module file. └── README.md # This file.

## Prerequisites

- **Go 1.18 or Later:** Download and install from [golang.org/dl](https://golang.org/dl/).
- **Visual Studio Code (Optional):** Use with the Go extension for enhanced coding and debugging support.
- **SQLite:** No installation required; the project uses an embedded SQLite database via the `go-sqlite3` driver.
- **7z Archive Files:** Ensure you have a 7z archive containing FB2 files to use for testing the import functionality.

## Installation and Setup

1. **Clone or Create the Project Directory:**
   ```bash
   mkdir golibfb2
   cd golibfb2
