package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
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

func validateOutputFlags(inputPath string, inPlace bool, outputPath string) error {
	if inPlace && outputPath != "" && outputPath != "-" {
		return fmt.Errorf("--in-place and --out cannot be used together")
	}
	if inPlace && (inputPath == "" || inputPath == "-") {
		return fmt.Errorf("--in-place requires --in to reference a file")
	}
	return nil
}

// writeOutput writes content based on flags.
// If inPlace is true, writes to the original file.
// If outputPath is set, writes there.
// Otherwise writes to stdout.
func writeOutput(content string, filePath string, inPlace bool, outputPath string) error {
	if err := validateOutputFlags(filePath, inPlace, outputPath); err != nil {
		return err
	}
	if inPlace {
		return os.WriteFile(filePath, []byte(content), 0644)
	}
	if outputPath != "" && outputPath != "-" {
		return os.WriteFile(outputPath, []byte(content), 0644)
	}
	_, err := fmt.Fprint(os.Stdout, content)
	return err
}

func renderNode(node *yaml.Node, format string) (string, error) {
	switch format {
	case "yaml":
		return renderYAML(node)
	case "json":
		return renderJSON(node)
	default:
		return "", fmt.Errorf("invalid --format %q: must be yaml or json", format)
	}
}

func renderYAML(node *yaml.Node) (string, error) {
	if node == nil {
		return "null\n", nil
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(node); err != nil {
		return "", err
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderJSON(node *yaml.Node) (string, error) {
	if node == nil {
		return "null\n", nil
	}
	var val interface{}
	if err := node.Decode(&val); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(val); err != nil {
		return "", err
	}
	return buf.String(), nil
}
