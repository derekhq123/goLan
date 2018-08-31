package config

import (
	"github.com/creamdog/gonfig"
	"os"
)

var ConfigFile gonfig.Gonfig

//set the configuration file variable as global variable which can be accessed from other package
func SetConfig(configName string) {
	f, err := os.Open(configName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	ConfigFile, _ = gonfig.FromJson(f)

}
