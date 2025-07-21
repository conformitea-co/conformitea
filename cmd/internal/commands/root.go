package commands

import (
	"fmt"
	"os"

	"conformitea/cmd/config"

	"github.com/spf13/cobra"
)

func Execute(config config.Config) {
	rootCmd := &cobra.Command{
		Use:   "conformitea",
		Short: "Compliance Management, Simplified.",
	}

	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(ServeCmd(config))

	rootCmd.ErrOrStderr()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong: %v\n", err)
		os.Exit(1)
	}
}
