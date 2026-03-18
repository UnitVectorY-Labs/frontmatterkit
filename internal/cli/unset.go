package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

var unsetPath string
var unsetIn string
var unsetInPlace bool
var unsetOut string

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Remove a value from front matter",
	Long: `Remove a value at the specified path from front matter.

Use --in to read from a file. If --in is omitted, input is read from stdin.
Use --out to write the updated document to a file or --in-place to overwrite the --in file.
If the path does not exist, the command is a no-op.`,
	Example: `  frontmatterkit unset --path .draft --in post.md
  frontmatterkit unset --path .draft --in post.md --in-place
  cat post.md | frontmatterkit unset --path .author.name
  frontmatterkit unset help`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOutputFlags(unsetIn, unsetInPlace, unsetOut); err != nil {
			exitCode = 2
			return err
		}

		content, err := readInput(unsetIn)
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
		if err := writeOutput(output, unsetIn, unsetInPlace, unsetOut); err != nil {
			exitCode = 3
			return err
		}
		return nil
	},
}

func init() {
	unsetCmd.Flags().StringVar(&unsetPath, "path", "", "Path to remove (required)")
	unsetCmd.MarkFlagRequired("path")
	unsetCmd.Flags().StringVar(&unsetIn, "in", "", "Read the Markdown document from this file instead of stdin")
	unsetCmd.Flags().BoolVar(&unsetInPlace, "in-place", false, "Overwrite the source file")
	unsetCmd.Flags().StringVar(&unsetOut, "out", "", "Write the updated document to this file instead of stdout")
	unsetCmd.MarkFlagsMutuallyExclusive("in-place", "out")
}
