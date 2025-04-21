package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"
	"github.com/taylorskalyo/goreader/epub"
)

func readEPUB(filepath string) (*epub.Reader, error) {
	// Open the EPUB filepath
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open EPUB file: %w", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	// Read entire file into memory
	content := make([]byte, fileInfo.Size())
	_, err = io.ReadFull(file, content)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}
	bytereader := bytes.NewReader(content)
	epubReader, err := epub.NewReader(bytereader, int64(len(content)))
	if err != nil {
		return nil, fmt.Errorf("failed to create EPUB reader: %w", err)
	}
	return epubReader, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <epub_file> <output_dir>")
		os.Exit(1)
	}
	epubFile := os.Args[1]
	outputDir := os.Args[2]

	// Read epub file
	reader, err := readEPUB(epubFile)
	if err != nil {
		fmt.Printf("Failed to read EPUB file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processing: %s\n", reader.Rootfiles[0].Metadata.Title)

	// Create output directory (outputDir/bookname)
	booktitleSlug := slug.Make(reader.Rootfiles[0].Metadata.Title)
	bookTargetDir := filepath.Join(outputDir, booktitleSlug)
	err = os.MkdirAll(bookTargetDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created directory: %s\n", bookTargetDir)
}
