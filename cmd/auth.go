package cmd

import (
	"fmt"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/UtkarshVerma/qgmail/auth"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authorize qGmail",
	Long: `Authorize qGmail with permissions to use Gmail API.

qGmail uses Google's recommended authorization flow, OAuth2(with PKCE extension),
making its transactions with the API highly secure.`,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(auth.GetToken())
		cobra.CheckErr(auth.SaveToken(tokenFile))
		fmt.Println("qGmail has been successfully authorized.")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVar(&auth.CredsFile, "credentials", "$HOME/.config/qgmail/credentials.json", "credentials file")

	// No need for error checking here as `homedir.Dir()` has already been checked
	home, _ := homedir.Dir()
	auth.CredsFile = strings.Replace(auth.CredsFile, "$HOME", home, 1)
	cobra.CheckErr(auth.ReadCredentials())
}
