package trash

import (
	"fmt"
	"specifiedCommand/config"
)

func TestGlobalVar() {
	number, _ := config.ConfigFile.GetInt("a", 0)
	fmt.Println(number)
}
