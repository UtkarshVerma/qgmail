package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/mitchellh/go-homedir"
)

type config struct {
	RedirectPort int    `json:"redirectPort"`
	CredsFile    string `json:"credentials"`
	TokenFile    string `json:"tokenFile"`
}

// Initialize config to default values.
func newConfig() *config {
	rport, _ := strconv.Atoi(flag.Lookup("rport").DefValue)
	return &config{
		CredsFile:    flag.Lookup("creds").DefValue,
		TokenFile:    flag.Lookup("t").DefValue,
		RedirectPort: rport,
	}
}

// To be used with flag.Visit() to override config flags at runtime
func (conf *config) updateConfig(f *flag.Flag) {
	switch f.Name {
	case "creds":
		conf.CredsFile = f.Value.String()
	case "t":
		conf.TokenFile = f.Value.String()
	case "rport":
		rport, _ := strconv.Atoi(f.Value.String())
		conf.RedirectPort = rport
	}
}

func (conf *config) readConfig(configFile string) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Fatalf("%s: no such file or directory", configFile)
	} else {
		f, err := os.Open(configFile)
		defer f.Close()

		byteValue, _ := ioutil.ReadAll(f)

		if err = json.Unmarshal(byteValue, conf); err != nil {
			log.Fatalf("%s: %v", configFile, err)
		}
	}

	conf.TokenFile, _ = homedir.Expand(conf.TokenFile)
	conf.CredsFile, _ = homedir.Expand(conf.CredsFile)
}
