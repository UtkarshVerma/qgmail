package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/mitchellh/go-homedir"
	"github.com/utkarshverma/qgmail/cli"
)

type (
	// Config is a struct having configurations for qGmail.
	conf struct {
		FileName  *string
		CredsFile *string
		TokenFile *string
		Init      initConf
	}
	initConf struct {
		MustPaste   *bool
		MustShowURL *bool
		Timeout     *int
	}

	tmpConf struct {
		CredsFile string `json:"credentials"`
		TokenFile string `json:"token"`
		Init      struct {
			MustPaste, MustShowURL bool
			Timeout                int
		}
	}
)

var (
	config = &conf{
		FileName:  cli.ConfigFile,
		CredsFile: cli.CredsFile,
		TokenFile: cli.TokenFile,
		Init: initConf{
			MustPaste:   cli.PasteFlag,
			MustShowURL: cli.URLFlag,
			Timeout:     cli.TimeoutFlag,
		},
	}

	FileName  = &config.FileName
	CredsFile = &config.CredsFile
	TokenFile = &config.TokenFile
	Init      = &config.Init
)

func init() {
	*config.FileName, _ = homedir.Expand(*config.FileName)
	config.read(*config.FileName)
	*config.CredsFile, _ = homedir.Expand(*config.CredsFile)
	*config.TokenFile, _ = homedir.Expand(*config.TokenFile)
}

func (conf *conf) read(configPath string) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If the config file doesn't exist at default path, create it.
		if defaultPath, _ := homedir.Expand("~/.config/qgmail/config.json"); configPath == defaultPath {
			configFile, _ := json.MarshalIndent(*config, "", "\t")
			_ = ioutil.WriteFile(configPath, configFile, 0644)
		} else {
			log.Fatal(err)
		}
	} else {
		c := &tmpConf{}

		f, _ := os.Open(configPath)
		defer f.Close()

		byteValue, _ := ioutil.ReadAll(f)

		if err = json.Unmarshal(byteValue, c); err != nil {
			log.Fatal(err)
		}

		// Override configurations as specified by flags.
		flag.Visit(c.update)
		cli.InitCmd.Visit(c.update)
		conf.update(c)
	}
}

func (conf *conf) update(c *tmpConf) {
	if c.CredsFile != "" {
		conf.CredsFile = &c.CredsFile
	}
	if c.TokenFile != "" {
		conf.TokenFile = &c.TokenFile
	}
	if c.Init.MustPaste {
		conf.Init.MustPaste = &c.Init.MustPaste
	}
	if c.Init.MustShowURL {
		conf.Init.MustShowURL = &c.Init.MustShowURL
	}
	if c.Init.Timeout > 0 {
		conf.Init.Timeout = &c.Init.Timeout
	}
}

func (conf *tmpConf) update(f *flag.Flag) {
	switch f.Name {
	case "creds":
		conf.CredsFile = f.Value.String()
	case "token":
		conf.TokenFile = f.Value.String()
	case "paste":
		conf.Init.MustPaste = true
	case "show-url":
		conf.Init.MustShowURL = true
	case "timeout":
		conf.Init.Timeout, _ = strconv.Atoi(f.Value.String())
	}
}
