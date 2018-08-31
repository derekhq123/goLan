package handlers

import (
	"github.com/creamdog/gonfig"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func Contains(key string, headlist []string) bool {
	for _, v := range headlist {
		if v == key {
			return true
		}
	}
	return false
}

func TestTableParser(t *testing.T) {
	f, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	config, err := gonfig.FromJson(f)
	if err != nil {
		panic(err)
	}
	command, _ := config.GetString("qstat-q/command", nil)
	parameter, _ := config.GetString("qstat-q/parameter", nil)

	input, _ := exec.Command(command, parameter).Output()

	r := TableParser(string(input), "qstat-q")
	head, _ := config.GetString("qstat-q/head", nil)
	headList := strings.Split(head, ",")
	headListLen := len(headList)

	outputList := r["Job"]
	for _, v := range outputList {
		keySet := reflect.ValueOf(v).MapKeys()

		if len(keySet) != headListLen {
			t.Errorf("parser failed. Inner map keySet length unmatched with header list length. Got %d, expected %d", len(keySet), headListLen)
		}

		for _, n := range keySet {

			valueString := n.String()

			flag := Contains(valueString, headList)
			if flag == false {
				t.Errorf("key %s not in header list", valueString)
			}
		}
	}
}
