package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	rootDirectory := "./"
	outputFile, err := os.Create("combined.py")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	// Use a map to track processed directories
	processedDirs := make(map[string]bool)

	err = filepath.Walk(rootDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error processing path %s: %v\n", path, err)
			return nil
		}

		// Get the directory of the current file
		dir := filepath.Dir(path)

		// Check if the directory has already been processed
		if processedDirs[dir] {
			return nil // Skip this file
		}

		// Check if it's a .py file
		if !info.IsDir() && filepath.Ext(path) == ".py" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v\n", path, err)
				return nil
			}

			fmt.Fprintf(outputFile, "### File: %s\n\n", path)
			outputFile.Write(content)
			outputFile.WriteString("\n---\n\n")

			// Mark the directory as processed after combining all .py files within it
			processedDirs[dir] = true
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path: %v", err)
	}

	fmt.Println("Files combined successfully!")
}
