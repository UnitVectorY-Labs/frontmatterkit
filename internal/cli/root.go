package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var exitCode int

var rootCmd = &cobra.Command{
	Use:   "frontmatterkit",
	Short: "A Unix-style CLI for managing YAML front matter in Markdown files",
	Long: `frontmatterkit validates, queries, asserts, and updates YAML front matter in Markdown files.

Use "frontmatterkit help <command>" or "frontmatterkit <command> help" for command-specific usage.`,
	Example: `  frontmatterkit validate --in post.md
  frontmatterkit get --path .title --in post.md
  frontmatterkit set --set '.draft=false' --in post.md --in-place
  cat post.md | frontmatterkit assert --assert '.title exists'`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute runs the root command and exits with the appropriate code.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		if exitCode == 0 {
			exitCode = 2
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	os.Exit(exitCode)
}

func init() {
	commands := []*cobra.Command{validateCmd, getCmd, setCmd, unsetCmd, assertCmd}
	for _, cmd := range commands {
		addHelpSubcommand(cmd)
		rootCmd.AddCommand(cmd)
	}
}

func addHelpSubcommand(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "help",
		Short: fmt.Sprintf("Show help for %s", parent.Name()),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return parent.Help()
		},
	})
}
