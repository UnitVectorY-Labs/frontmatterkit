package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/assertion"
	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
)

var assertExprs []string
var assertIn string

var assertCmd = &cobra.Command{
	Use:   "assert",
	Short: "Check front matter against declarative conditions",
	Long: `Check whether front matter satisfies one or more declarative conditions.

Multiple --assert flags are combined with logical AND.
Missing front matter is treated as an empty object.
Use --in to read from a file. If --in is omitted, input is read from stdin.

Supported operators: exists, not exists, ==, !=, <, <=, >, >=, contains, not contains`,
	Example: `  frontmatterkit assert --assert '.title exists' --in post.md
  frontmatterkit assert --assert '.draft == false' --assert '.tags contains "go"' --in post.md
  cat post.md | frontmatterkit assert --assert '.title exists'
  frontmatterkit assert help`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(assertExprs) == 0 {
			exitCode = 2
			return fmt.Errorf("at least one --assert flag is required")
		}

		content, err := readInput(assertIn)
		if err != nil {
			exitCode = 3
			return err
		}

		doc, err := frontmatter.ParseAndValidate(content)
		if err != nil {
			exitCode = 1
			return err
		}

		// Parse all assertions first to fail fast on syntax errors
		assertions := make([]*assertion.Assertion, 0, len(assertExprs))
		for _, expr := range assertExprs {
			a, err := assertion.ParseAssertion(expr)
			if err != nil {
				exitCode = 2
				return fmt.Errorf("invalid assertion %q: %w", expr, err)
			}
			assertions = append(assertions, a)
		}

		var failures []string
		for i, a := range assertions {
			if err := a.Evaluate(doc); err != nil {
				failures = append(failures, fmt.Sprintf("assertion %q failed: %s", assertExprs[i], err))
			}
		}

		if len(failures) > 0 {
			exitCode = 1
			return fmt.Errorf("%s", strings.Join(failures, "\n"))
		}

		return nil
	},
}

func init() {
	assertCmd.Flags().StringArrayVar(&assertExprs, "assert", nil, "Assertion expression (repeatable)")
	assertCmd.Flags().StringVar(&assertIn, "in", "", "Read the Markdown document from this file instead of stdin")
}
