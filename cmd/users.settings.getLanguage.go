// Code generated by go generate; DO NOT EDIT.
package cmd

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/UtkarshVerma/qgmail/auth"
	"github.com/spf13/cobra"
)

var usersSettingsGetLanguage = &cobra.Command{
	Use: "users.settings.getLanguage <userId>",
	Short: "Gets language settings.",
	Long: "userId: User's email address. The special value `me` can be used to indicate the authenticated user.\n",
	DisableFlagsInUseLine: true,
	Args: cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.ReadToken(tokenFile); err != nil {
			fmt.Println("Error: Authorization token not found. Please authorize qGmail using 'qgmail auth'.")
			os.Exit(1)
		}

		service, err := auth.NewGmailService()
		cobra.CheckErr(err)

		val, err := service.Users.Settings.GetLanguage(args[0]).Do()
		cobra.CheckErr(err)

		jsonData, _ := json.Marshal(val)
		fmt.Println(string(jsonData))
	},
}

func init() {
	requestCmd.AddCommand(usersSettingsGetLanguage)
}