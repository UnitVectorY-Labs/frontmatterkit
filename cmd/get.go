package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

var getFormat string
var getPath string

var getCmd = &cobra.Command{
	Use:   "get [file]",
	Short: "Extract front matter or selected values",
	Long: `Extract front matter or selected values from a Markdown file.

Default output format is YAML. Use --format json for JSON output.
If no front matter exists, behaves as if front matter were an empty object.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := ""
		if len(args) > 0 {
			filePath = args[0]
		}

		content, err := readInput(filePath)
		if err != nil {
			exitCode = 3
			return err
		}

		doc, err := frontmatter.ParseAndValidate(content)
		if err != nil {
			exitCode = 1
			return err
		}

		path, err := pathsyntax.Parse(getPath)
		if err != nil {
			exitCode = 2
			return fmt.Errorf("invalid path %q: %w", getPath, err)
		}

		node := doc.Get(path)

		if getFormat == "json" {
			return outputJSON(node)
		}
		return outputYAML(node)
	},
}

func init() {
	getCmd.Flags().StringVar(&getFormat, "format", "yaml", "Output format: yaml or json")
	getCmd.Flags().StringVar(&getPath, "path", ".", "Path to extract (jq-like syntax)")
}

func outputYAML(node *yaml.Node) error {
	if node == nil {
		fmt.Fprintln(os.Stdout, "null")
		return nil
	}
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	err := enc.Encode(node)
	enc.Close()
	return err
}

func outputJSON(node *yaml.Node) error {
	if node == nil {
		fmt.Fprintln(os.Stdout, "null")
		return nil
	}
	var val interface{}
	if err := node.Decode(&val); err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(val)
}
