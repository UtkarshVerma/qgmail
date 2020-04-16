package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
)

type config struct {
	CredsFile string `json:"credentials"`
	TokenFile string `json:"tokenFile"`
}

// Initialize config to default values.
func newConfig() *config {
	return &config{
		CredsFile: flag.Lookup("creds").DefValue,
		TokenFile: flag.Lookup("t").DefValue,
	}
}

// To be used with flag.Visit() to override config flags at runtime
func (conf *config) updateConfig(f *flag.Flag) {
	switch f.Name {
	case "creds":
		conf.CredsFile = f.Value.String()
	case "t":
		conf.TokenFile = f.Value.String()
	}
}

func (conf *config) readConfig(configFile string) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Fatal("Error: The configuration file doesn't exist.")
	} else {
		f, _ := os.Open(configFile)
		defer f.Close()

		byteValue, _ := ioutil.ReadAll(f)

		if err = json.Unmarshal(byteValue, conf); err != nil {
			log.Fatalf("Error: Could not parse the configuration file.\n%v", err)
		}
	}

	conf.TokenFile, _ = homedir.Expand(conf.TokenFile)
	conf.CredsFile, _ = homedir.Expand(conf.CredsFile)
}
