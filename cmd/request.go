package cmd

import (
	"github.com/spf13/cobra"
)

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request the Gmail API for a resource",
	Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	rootCmd.AddCommand(requestCmd)
}
