package cmd

import (
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile, tokenFile string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "qgmail",
		Short: "The power of Gmail API now in your comfy terminal.",
		Long: `qGmail is a CLI tool written in Go which lets you query info related to your
Gmail account.

qGmail uses Google's recommended authorization flow, OAuth2(with PKCE extension),
making its transactions with the API highly secure.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "$HOME/.config/qgmail/config.json", "config file")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token", "$HOME/.cache/qgmail/token.json", "cached token file")

	home, err := homedir.Dir()
	cobra.CheckErr(err)
	for _, flag := range []*string{&configFile, &tokenFile} {
		*flag = strings.Replace(*flag, "$HOME", home, 1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig() // read in config, if exists
}
