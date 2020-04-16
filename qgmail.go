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
	initFlag   = flag.Bool("init", false, "Reconfigure qGmail.")
	configFile = flag.String("c", homeDir+"/.config/qgmail/config.json", "path to config file")
	credsFile  = flag.String("creds", homeDir+"/.config/qgmail/credentials.json", "path to credentials file")
	tokenFile  = flag.String("t", homeDir+"/.config/qgmail/token.json", "path to token file")
)

func init() {
	flag.Parse()

	// Read the config and credentials.
	conf = newConfig()
	conf.readConfig(*configFile)
	flag.Visit(conf.updateConfig)
}

func main() {
	oauthConf := newOauthConf(*credsFile, conf)

	if *initFlag {
		token := getTokenFromWeb(oauthConf)
		saveToken(conf.TokenFile, token)
		return
	}

	token, err := tokenFromFile(conf.TokenFile)
	if err != nil {
		*initFlag = true
		fmt.Print("Authorization token not present. Fetching a new one...")
		token = getTokenFromWeb(oauthConf)
		saveToken(conf.TokenFile, token)
	}

	if !*initFlag {
		service, _ := newGmailService(oauthConf, token)
		labelStruct := getLabel("INBOX", service)
		fmt.Println(labelStruct.MessagesUnread)
	}
}
