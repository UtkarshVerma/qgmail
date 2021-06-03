package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of qGmail",
	Long:  `All software has versions. This is qGmail's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("qGmail v1.0")
	},
}
