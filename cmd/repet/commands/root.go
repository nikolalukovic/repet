package commands

import (
	"fmt"
	"os"

	"github.com/nikolalukovic/repet/cmd/repet/commands/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "repet",
	Short: "Repet is a simple key-value repeater/cache",
	Long:  "Repet is a simple key-value repeater/cache",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(version.VersionCmd)
}
