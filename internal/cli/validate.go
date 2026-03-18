package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate YAML front matter in a Markdown file",
	Long: `Check whether the file's front matter is well-formed.

Use --in to read from a file. If --in is omitted, input is read from stdin.
Files without front matter are valid. Files with malformed front matter exit with code 1.`,
	Example: `  frontmatterkit validate --in post.md
  cat post.md | frontmatterkit validate
  frontmatterkit validate help`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := readInput(validateIn)
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

var validateIn string

func init() {
	validateCmd.Flags().StringVar(&validateIn, "in", "", "Read the Markdown document from this file instead of stdin")
}
