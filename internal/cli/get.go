package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

var getFormat string
var getPath string
var getIn string
var getOut string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Extract front matter or selected values",
	Long: `Extract front matter or selected values from a Markdown file.

Use --in to read from a file. If --in is omitted, input is read from stdin.
Use --out to write the extracted value to a file instead of stdout.
If no front matter exists, the document behaves as if it had an empty object.`,
	Example: `  frontmatterkit get --path .title --in post.md
  frontmatterkit get --format json --path .tags --in post.md --out tags.json
  cat post.md | frontmatterkit get --path .author.name
  frontmatterkit get help`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := readInput(getIn)
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

		output, err := renderNode(node, getFormat)
		if err != nil {
			exitCode = 2
			return err
		}
		if err := writeOutput(output, "", false, getOut); err != nil {
			exitCode = 3
			return err
		}
		return nil
	},
}

func init() {
	getCmd.Flags().StringVar(&getFormat, "format", "yaml", "Output format: yaml or json")
	getCmd.Flags().StringVar(&getPath, "path", ".", "Path to extract (jq-like syntax)")
	getCmd.Flags().StringVar(&getIn, "in", "", "Read the Markdown document from this file instead of stdin")
	getCmd.Flags().StringVar(&getOut, "out", "", "Write the extracted value to this file instead of stdout")
}
