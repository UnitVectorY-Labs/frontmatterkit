package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

var unsetPath string
var unsetInPlace bool
var unsetOutput string

var unsetCmd = &cobra.Command{
	Use:   "unset [file]",
	Short: "Remove a value from front matter",
	Long: `Remove a value at the specified path from front matter.

If the path does not exist, the command is a no-op.
The Markdown body and other front matter values are preserved.`,
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

		path, err := pathsyntax.Parse(unsetPath)
		if err != nil {
			exitCode = 2
			return fmt.Errorf("invalid path %q: %w", unsetPath, err)
		}

		doc.Unset(path)

		output := doc.Render()
		return writeOutput(output, filePath, unsetInPlace, unsetOutput)
	},
}

func init() {
	unsetCmd.Flags().StringVar(&unsetPath, "path", "", "Path to remove (required)")
	unsetCmd.MarkFlagRequired("path")
	unsetCmd.Flags().BoolVar(&unsetInPlace, "in-place", false, "Overwrite the source file")
	unsetCmd.Flags().StringVar(&unsetOutput, "output", "", "Write output to specified file path")
}
