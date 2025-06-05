package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "conformitea",
	Short: "Sit back and relax",
	Long:  "More details and documentation soon.",
}

func Execute() {
	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(ServeCmd())

	rootCmd.ErrOrStderr()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
