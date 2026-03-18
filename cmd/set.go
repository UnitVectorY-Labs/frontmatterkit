package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

var setValues []string
var setFrom string
var setMode string
var setInPlace bool
var setOutput string

var setCmd = &cobra.Command{
	Use:   "set [file]",
	Short: "Create or update front matter values",
	Long: `Create or update front matter values in a Markdown file.

If the file has no front matter, a new YAML block is created.
Values are interpreted as YAML, not raw strings.`,
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

		mode := frontmatter.SetOverwrite
		if setMode == "patch" {
			mode = frontmatter.SetPatch
		}

		// Process --from if specified
		if setFrom != "" {
			fromContent, err := readInput(setFrom)
			if err != nil {
				exitCode = 3
				return fmt.Errorf("reading --from: %w", err)
			}
			rootPath, _ := pathsyntax.Parse(".")
			if err := doc.Set(rootPath, fromContent, mode); err != nil {
				exitCode = 1
				return fmt.Errorf("applying --from content: %w", err)
			}
		}

		// Process each --set flag
		for _, sv := range setValues {
			idx := strings.IndexByte(sv, '=')
			if idx < 1 {
				exitCode = 2
				return fmt.Errorf("invalid --set value %q: must be in format .path=value", sv)
			}
			pathStr := sv[:idx]
			value := sv[idx+1:]

			path, err := pathsyntax.Parse(pathStr)
			if err != nil {
				exitCode = 2
				return fmt.Errorf("invalid path in --set %q: %w", sv, err)
			}

			if err := doc.Set(path, value, mode); err != nil {
				exitCode = 1
				return fmt.Errorf("setting %s: %w", pathStr, err)
			}
		}

		output := doc.Render()
		return writeOutput(output, filePath, setInPlace, setOutput)
	},
}

func init() {
	setCmd.Flags().StringArrayVar(&setValues, "set", nil, "Set a value: .path=yamlValue (repeatable)")
	setCmd.Flags().StringVar(&setFrom, "from", "", "Read YAML values from file or stdin (-)")
	setCmd.Flags().StringVar(&setMode, "mode", "overwrite", "Set mode: overwrite or patch")
	setCmd.Flags().BoolVar(&setInPlace, "in-place", false, "Overwrite the source file")
	setCmd.Flags().StringVar(&setOutput, "output", "", "Write output to specified file path")
}
