package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var exitCode int

var rootCmd = &cobra.Command{
	Use:   "frontmatterkit",
	Short: "A Unix-style CLI for managing YAML front matter in Markdown files",
	Long:  "frontmatterkit validates, queries, asserts, and minimally updates YAML front matter in Markdown files.",
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
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(unsetCmd)
	rootCmd.AddCommand(assertCmd)
}
