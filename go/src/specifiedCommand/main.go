package main

import (
	"flag"
	"fmt"

	"reflect"
	"specifiedCommand/config"
	"specifiedCommand/trash"
)

func init() {
	configAddress := flag.String("c", "config1.json", "Config file address")
	flag.Parse()
	configName := *configAddress
	config.SetConfig(configName)
}

func main() {

	fmt.Println(reflect.TypeOf(config.ConfigFile))
	number, _ := config.ConfigFile.GetInt("a", 0)
	fmt.Println(number)
	trash.TestGlobalVar()
}
