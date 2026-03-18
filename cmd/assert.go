package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/assertion"
	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
)

var assertExprs []string

var assertCmd = &cobra.Command{
	Use:   "assert [file]",
	Short: "Check front matter against declarative conditions",
	Long: `Check whether front matter satisfies one or more declarative conditions.

Multiple --assert flags are combined with logical AND.
Missing front matter is treated as an empty object.

Supported operators: exists, not exists, ==, !=, <, <=, >, >=, contains, not contains`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(assertExprs) == 0 {
			exitCode = 2
			return fmt.Errorf("at least one --assert flag is required")
		}

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

		// Evaluate all assertions
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
}
