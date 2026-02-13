package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forkwatch",
	Short: "Discover meaningful patches hiding in forks",
	Long:  `Forkwatch analyzes GitHub repository forks to find meaningful changes that haven't been submitted as pull requests. It groups forks by the files they modify and highlights convergence — when multiple independent forks touch the same code.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
