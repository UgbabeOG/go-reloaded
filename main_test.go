package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIntegrationExample1(t *testing.T) {
	input := "it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair."
	expected := "It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair."

	result := Process(input)
	if result != expected {
		t.Errorf("Example 1 failed.\nGot:      %q\nExpected: %q", result, expected)
	}
}

func TestIntegrationExample2(t *testing.T) {
	input := "Simply add 42 (hex) and 10 (bin) and you will see the result is 68."
	expected := "Simply add 66 and 2 and you will see the result is 68."

	result := Process(input)
	if result != expected {
		t.Errorf("Example 2 failed.\nGot:      %q\nExpected: %q", result, expected)
	}
}

func TestIntegrationExample3(t *testing.T) {
	input := "There is no greater agony than bearing a untold story inside you."
	expected := "There is no greater agony than bearing an untold story inside you."

	result := Process(input)
	if result != expected {
		t.Errorf("Example 3 failed.\nGot:      %q\nExpected: %q", result, expected)
	}
}

func TestIntegrationExample4(t *testing.T) {
	input := "Punctuation tests are ... kinda boring ,what do you think ?"
	expected := "Punctuation tests are... kinda boring, what do you think?"

	result := Process(input)
	if result != expected {
		t.Errorf("Example 4 failed.\nGot:      %q\nExpected: %q", result, expected)
	}
}

func TestIntegrationWithFiles(t *testing.T) {
	// Create a temporary file for input
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "test_input.txt")
	outputPath := filepath.Join(tmpDir, "test_output.txt")

	input := "hello (up) world"
	expected := "HELLO world"

	// Write input file
	err := os.WriteFile(inputPath, []byte(input), 0644)
	if err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	// Read, process, and write
	content, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("Failed to read input file: %v", err)
	}

	result := Process(string(content))

	err = os.WriteFile(outputPath, []byte(result), 0644)
	if err != nil {
		t.Fatalf("Failed to write output file: %v", err)
	}

	// Read and verify output
	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(output) != expected {
		t.Errorf("File integration test failed.\nGot:      %q\nExpected: %q", string(output), expected)
	}
}
