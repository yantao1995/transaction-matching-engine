package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "transaction-matching-engine",
	Short: "虚拟货币交易撮合引擎",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Help()
	}
}
