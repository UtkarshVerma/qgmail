package cli

import (
	"flag"
	"log"
	"os"
)

var (
	InitCmd    = flag.NewFlagSet("init", flag.ExitOnError)
	ConfigFile = flag.String("config", "~/.config/qgmail/config.json", "Path to qGmail's configuration file.")

	TimeoutFlag          = InitCmd.Int("timeout", 1, "Timeout(in minutes) for user-consent page.")
	URLFlag              = InitCmd.Bool("show-url", false, "Displays the user consent URL.")
	PasteFlag            = InitCmd.Bool("paste", false, "Use this flag if you want to paste authorization code manually.")
	CredsFile, TokenFile *string
	InitMode             bool
)

func init() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			CredsFile = InitCmd.String("creds", "~/.config/qgmail/credentials.json", "Path to Google API client credentials.")
			TokenFile = InitCmd.String("token", "~/.config/qgmail/token.json", "Path for storing the authorization token.")
			InitCmd.Parse(os.Args[2:])
		default:
			log.Fatal("Invalid argument passed. Run 'qgmail --help` for valid options.")
		}
	}
	flag.Parse()
	if InitMode = InitCmd.Parsed(); InitMode {
		CredsFile = flag.String("creds", "~/.config/qgmail/credentials.json", "Path to Google API client credentials.")
		TokenFile = flag.String("token", "~/.config/qgmail/token.json", "Path for storing the authorization token.")
	}
}
