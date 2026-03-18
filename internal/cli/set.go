package cli

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
var setIn string
var setInPlace bool
var setOut string

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Create or update front matter values",
	Long: `Create or update front matter values in a Markdown file.

Use --in to read from a file. If --in is omitted, input is read from stdin.
Use --out to write the updated document to a file or --in-place to overwrite the --in file.
If the document has no front matter, a new YAML block is created.
Values supplied with --set are parsed as YAML, not raw strings.`,
	Example: `  frontmatterkit set --set '.title="New Title"' --in post.md
  frontmatterkit set --set '.draft=false' --in post.md --in-place
  frontmatterkit set --from values.yaml --mode patch --in post.md --out updated.md
  cat post.md | frontmatterkit set --set '.title=From stdin'
  frontmatterkit set help`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOutputFlags(setIn, setInPlace, setOut); err != nil {
			exitCode = 2
			return err
		}

		content, err := readInput(setIn)
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
		switch setMode {
		case "overwrite":
		case "patch":
			mode = frontmatter.SetPatch
		default:
			exitCode = 2
			return fmt.Errorf("invalid --mode %q: must be overwrite or patch", setMode)
		}

		if setFrom != "" {
			if setFrom == "-" && (setIn == "" || setIn == "-") {
				exitCode = 2
				return fmt.Errorf("--from - cannot be combined with stdin input; use --in for the document")
			}
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
		if err := writeOutput(output, setIn, setInPlace, setOut); err != nil {
			exitCode = 3
			return err
		}
		return nil
	},
}

func init() {
	setCmd.Flags().StringArrayVar(&setValues, "set", nil, "Set a value: .path=yamlValue (repeatable)")
	setCmd.Flags().StringVar(&setFrom, "from", "", "Read YAML values from file or stdin (-)")
	setCmd.Flags().StringVar(&setMode, "mode", "overwrite", "Set mode: overwrite or patch")
	setCmd.Flags().StringVar(&setIn, "in", "", "Read the Markdown document from this file instead of stdin")
	setCmd.Flags().BoolVar(&setInPlace, "in-place", false, "Overwrite the source file")
	setCmd.Flags().StringVar(&setOut, "out", "", "Write the updated document to this file instead of stdout")
	setCmd.MarkFlagsMutuallyExclusive("in-place", "out")
}
