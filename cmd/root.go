package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/UtkarshVerma/qgmail/apicall"
	"github.com/UtkarshVerma/qgmail/auth"
)

var configFile, tokenFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qgmail",
	Short: "The power of Gmail API now in your comfy terminal.",
	Long: `qGmail is a CLI tool written in Go which lets you query info related to your
Gmail account.

qGmail uses Google's recommended authorization flow, OAuth2(with PKCE extension),
making its transactions with the API highly secure.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.ReadToken(tokenFile); err != nil {
			fmt.Println("Error: Authorization token not found. Please authorize qGmail using 'qgmail auth'.")
			os.Exit(1)
		}

		service, err := auth.NewGmailService()
		cobra.CheckErr(err)

		label, err := apicall.Label("INBOX", service)
		cobra.CheckErr(err)
		fmt.Println(label.MessagesUnread)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", `config file (default "$HOME/.config/qgmail/config.(json|toml|yaml)")`)
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token", "$HOME/.cache/qgmail/token.json", "cached token file")

	home, err := homedir.Dir()
	cobra.CheckErr(err)
	tokenFile = strings.Replace(tokenFile, "$HOME", home, 1)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// No need for error checking here as `homedir.Dir()` is already checked in `init()`
	home, _ := homedir.Dir()
	if configFile == "" {
		viper.AddConfigPath(path.Join(home, ".config", "qgmail"))
		viper.SetConfigName("config")
	} else {
		viper.SetConfigFile(configFile)
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig() // read in config, if exists
}
