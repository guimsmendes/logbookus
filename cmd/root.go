package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{Use: "logbookus"}
	rootCmd.AddCommand(commands...)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1) // Fail with non-zero exit code
	}
}

var commands = []*cobra.Command{
	{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Start the server",
		Long:    "Start the server",
		Run: func(cmd *cobra.Command, args []string) {

		},
	},
	{
		Use:     "backup",
		Aliases: []string{"b"},
		Short:   "Backup",
		Long:    "Start the server",
		Run: func(cmd *cobra.Command, args []string) {

		},
	},
}
