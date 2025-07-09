package commands

import (
	"fmt"
	"os"

	"conformitea/server/config"

	"github.com/spf13/cobra"
)

func Execute(cfg config.Config) {
	rootCmd := &cobra.Command{
		Use:   "conformitea",
		Short: "Compliance Management, Simplified.",
	}

	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(ServeCmd(cfg))

	rootCmd.ErrOrStderr()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong: %v\n", err)
		os.Exit(1)
	}
}
