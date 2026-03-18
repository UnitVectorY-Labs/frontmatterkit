package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate YAML front matter in a Markdown file",
	Long: `Check whether the file's front matter is well-formed.

If no front matter is present, validation succeeds.
If front matter is present, it must satisfy strict YAML rules.
If the file begins with --- and the block is malformed, validation fails.`,
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

		_, err = frontmatter.ParseAndValidate(content)
		if err != nil {
			exitCode = 1
			return fmt.Errorf("validation failed: %w", err)
		}

		return nil
	},
}
