package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/creamdog/gonfig"

	"os"
	"os/exec"

	"strings"
)

// func configScenario1() {
// 	f, _ := os.Open("config.json")
// 	config, _ := gonfig.FromJson(f)
// 	head, _ := config.GetString("head", nil)
// 	numOfSkip, _ := config.GetInt("numOfSkip", 0)

// 	headList := strings.Split(head, ",")

// 	mapResult := make(map[string][]map[string]string)
// 	key := "Job"

// 	out, _ := exec.Command("cat", "qstat.out").Output()
// 	input := string(out)
// 	haha := strings.Split(input, "\n")
// 	fmt.Println(input)

// 	for k, v := range haha {
// 		fmt.Println(k)
// 		fmt.Println(v)
// 	}

// 	scanner := bufio.NewScanner(strings.NewReader(input))

// 	index := 0
// 	var tmpList []map[string]string
// 	xixi := 1
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		fmt.Println(line)
// 		fmt.Println(xixi)
// 		xixi += 1
// 		if index == numOfSkip {
// 			tmpMap := make(map[string]string)
// 			ss := strings.Fields(line)

// 			for k, v := range ss {

// 				tmpMap[headList[k]] = v
// 			}

// 			tmpList = append(tmpList, tmpMap)
// 		} else {
// 			index += 1
// 		}

// 	}
// 	mapResult[key] = tmpList
// 	//fmt.Println(mapResult)
// 	output, _ := json.Marshal(mapResult)
// 	fmt.Println(string(output))
// }
func main() {
	// configScenario1()

	out, _ := exec.Command("cat", "outputcmd/qstat-B.txt").Output()
	input := string(out)
	fmt.Println(input)
	head := "qstat-B"
	tableParser(input, head)
}

func tableParser(input string, configHead string) []byte {

	// if input == "" {
	// 	return "empty output"
	// }
	// num := len(strings.Split(input, "\n"))
	// if num <= 2 {
	// 	return input
	// }
	f, err := os.Open("config.json")
	if err != nil {
		fmt.Println("config file open wrong")
	}
	defer f.Close()
	config, err := gonfig.FromJson(f)
	if err != nil {
		fmt.Println("import json wrong")
	}
	headKey := configHead + "/head"
	numKey := configHead + "/numOfSkip"
	fmt.Println(headKey)
	fmt.Println(numKey)
	head, err := config.GetString(headKey, nil)
	if err != nil {
		fmt.Println("head content wrong")
	}
	numOfSkip, err := config.GetInt(numKey, 0)
	if err != nil {
		fmt.Println("num of skip wrong")
	}

	headList := strings.Split(head, ",")

	mapResult := make(map[string][]map[string]string)
	key := "Job"
	scanner := bufio.NewScanner(strings.NewReader(input))
	index := 0
	var tmpList []map[string]string
	for scanner.Scan() {
		line := scanner.Text()
		if index == numOfSkip {
			tmpMap := make(map[string]string)
			ss := strings.Fields(line)

			for k, v := range ss {

				tmpMap[headList[k]] = v
			}

			tmpList = append(tmpList, tmpMap)
		} else {
			index += 1
		}

	}
	mapResult[key] = tmpList
	//fmt.Println(mapResult)
	output, _ := json.Marshal(mapResult)
	fmt.Println(string(output))
	return output
}
