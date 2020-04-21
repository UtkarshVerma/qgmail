package main

import (
	"fmt"

	"github.com/utkarshverma/qgmail/apicall"
	"github.com/utkarshverma/qgmail/auth"
	"github.com/utkarshverma/qgmail/cli"
	"github.com/utkarshverma/qgmail/config"
)

var (
	oAuthConf, token      = auth.Config, &auth.Token
	configFile, tokenFile = **config.FileName, **config.TokenFile
)

func main() {
	if cli.InitMode {
		auth.GetToken(oAuthConf, token, auth.Params)
		auth.SaveToken(tokenFile, *token)
		return
	}

	err := auth.ReadToken(tokenFile, token)
	if err != nil {
		cli.InitMode = true
		fmt.Println("Authorization token not found. Fetching a new one...")
		auth.GetToken(oAuthConf, token, auth.Params)
		auth.SaveToken(tokenFile, *token)
		fmt.Printf("%+v", *token)
	}

	if !cli.InitMode {
		service, _ := auth.NewGmailService(auth.Config, *token)
		label := apicall.Label("INBOX", service)
		fmt.Println(label.MessagesUnread)
	}
}
