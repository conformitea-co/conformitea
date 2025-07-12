package commands

import (
	"conformitea/server"
	"conformitea/server/types"
	"fmt"

	"github.com/spf13/cobra"
)

func ServeCmd(cfg types.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "serve",
		Short:   "Start the ConformiTea server",
		Aliases: []string{"s", "start"},
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := server.Initialize(cfg)
			if err != nil {
				return fmt.Errorf("failed to initialize server (%w)", err)
			}

			return srv.Start()
		},
	}
}
