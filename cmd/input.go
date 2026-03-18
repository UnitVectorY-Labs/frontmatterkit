package cmd

import (
	"fmt"
	"io"
	"os"
)

// readInput reads content from the given file path or stdin.
// If filePath is empty or "-", reads from stdin.
func readInput(filePath string) (string, error) {
	if filePath == "" || filePath == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		return string(data), nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("reading file %s: %w", filePath, err)
	}
	return string(data), nil
}

// writeOutput writes content based on flags.
// If inPlace is true, writes to the original file.
// If outputPath is set, writes there.
// Otherwise writes to stdout.
func writeOutput(content string, filePath string, inPlace bool, outputPath string) error {
	if inPlace {
		if filePath == "" || filePath == "-" {
			return fmt.Errorf("--in-place requires a file path argument")
		}
		return os.WriteFile(filePath, []byte(content), 0644)
	}
	if outputPath != "" {
		return os.WriteFile(outputPath, []byte(content), 0644)
	}
	_, err := fmt.Fprint(os.Stdout, content)
	return err
}
