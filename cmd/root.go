package cmd

import (
	"fmt"
	"os"

	"github.com/guimsmendes/logbookus/config"
	"github.com/guimsmendes/logbookus/internal/server"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := server.New(config.Prod, 8080)
			if err != nil {
				return fmt.Errorf("new server: %w", err)
			}

			return s.Start(cmd.Context())
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
