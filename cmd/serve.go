package cmd

import (
	"github.com/conformitea-co/conformitea/internal/server"

	"github.com/spf13/cobra"
)

func ServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "serve",
		Short:   "Start the ConformiTea server",
		Aliases: []string{"s", "start"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.Start()
		},
	}
}
