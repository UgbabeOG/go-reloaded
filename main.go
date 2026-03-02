package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s input.txt output.txt\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// Process the text
	result := Process(string(content))

	// Write output file
	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}
}
