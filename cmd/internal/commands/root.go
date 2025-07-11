package commands

import (
	"fmt"
	"os"

	"conformitea/server/types"

	"github.com/spf13/cobra"
)

func Execute(cfg types.Config) {
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
