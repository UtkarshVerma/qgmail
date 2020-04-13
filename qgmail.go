package main

import (
	"flag"
	"fmt"

	"github.com/mitchellh/go-homedir"
)

var (
	// use different config locations depending on the os
	homeDir, _ = homedir.Dir()
	auth       = newAuthParams()
	conf       *config

	// Command line flags.
	configFile   = flag.String("c", homeDir+"/.config/qgmail/config.json", "path to config file")
	credsFile    = flag.String("creds", homeDir+"/.config/qgmail/credentials.json", "path to credentials file")
	tokenFile    = flag.String("t", homeDir+"/.config/qgmail/token.json", "path to token file")
	redirectPort = flag.Int("rport", 5000, "redirect URL port for the HTTP server")
)

func init() {
	flag.Parse()

	// Read the config and credentials.
	conf = newConfig()
	conf.readConfig(*configFile)
	flag.Visit(conf.updateConfig)
}

func main() {
	// Retrieve the token.
	oauthConf := newOauthConf(*credsFile, conf)
	token, err := tokenFromFile(conf.TokenFile)
	if err != nil {
		token = getTokenFromWeb(oauthConf)
		saveToken(conf.TokenFile, token)
	}

	service, _ := newGmailService(oauthConf, token)
	labelStruct := getLabel("INBOX", service)
	fmt.Println(labelStruct.MessagesUnread)
}
