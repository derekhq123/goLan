package config

import (
	"github.com/creamdog/gonfig"
	"os"
)

var ConfigFile gonfig.Gonfig

func SetConfig(configName string) {
	f, err := os.Open(configName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	ConfigFile, _ = gonfig.FromJson(f)

}
