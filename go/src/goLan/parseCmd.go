package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/creamdog/gonfig"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

type output struct {
	JobId         string `json:"JobId"`
	UserName      string `json:"userName"`
	Queue         string `json:"queue"`
	JobName       string `json:"jobName"`
	SessId        string `json:"sessId"`
	Nds           string `json:"nds"`
	Tsk           string `json:"tsk"`
	RequireMemory string `json:"requireMemory"`
	RequireTime   string `json:"requireTime"`
	S             string `json:"s"`
	ElapTime      string `json:"elapTime"`
}

type finalOpt struct {
	Job []output
}

func scenario1() {
	file1, err1 := os.Open("commandOpt1.txt")
	if err1 != nil {
		log.Fatal(err1)
	}
	scanner := bufio.NewScanner(file1)
	flag := false
	var tmp []output
	for scanner.Scan() {
		line := scanner.Text()
		if flag == true {
			ss := strings.Fields(line)

			res := output{
				ss[0],
				ss[1],
				ss[2],
				ss[3],
				ss[4],
				ss[5],
				ss[6],
				ss[7],
				ss[8],
				ss[9],
				ss[10],
			}
			fmt.Println(reflect.TypeOf(res))
			//fmt.Println(res)
			tmp = append(tmp, res)
			newRes, _ := json.Marshal(res)
			fmt.Println(string(newRes))
		}
		if strings.Contains(line, "-") {
			flag = true
		}

	}
	result := finalOpt{tmp}
	fmt.Println(result)
	mapp, _ := json.Marshal(result)
	fmt.Println(string(mapp))
}

//dynamic handle scenario 1
func dynamicScenario1() {
	file1, err1 := os.Open("commandOpt1.txt")
	if err1 != nil {
		log.Fatal(err1)
	}
	scanner := bufio.NewScanner(file1)
	//var outputMap map[string][]string
	var titlesList [][]string
	var titleList []string
	var lineList []string
	var indexList []int
	flag := false
	//var key string
	//wholeMap := make(map[string][]map[string]string)
	mapSlice := make([]map[string]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		ss := strings.Fields(line)

		if flag == true {
			tmpMap := make(map[string]string)
			for p, q := range ss {
				fmt.Println(q)
				tmpKey := titleList[p]
				tmpValue := q

				tmpMap[tmpKey] = tmpValue
			}
			mapSlice = append(mapSlice, tmpMap)
			fmt.Println(mapSlice)

		} else if strings.Contains(line, ":") {
			//key = strings.TrimRight(ss[0], ":")

		} else if strings.Contains(line, "-") {
			flag = true
			indexList = append(indexList, 0)
			for k, v := range line {

				if k != 0 && v == '-' && line[k-1] == ' ' {
					indexList = append(indexList, k)
				}

			}
			indexList = append(indexList, len(line))
			fmt.Println(indexList)
			if len(titlesList) == 1 {

				titles := titlesList[0]
				for _, n := range titles {
					titleList = append(titleList, n)
				}
			} else {
				maxNumber := 0

				var referenceLine string
				for m, n := range titlesList {
					if len(n) > maxNumber {
						maxNumber = len(n)

						referenceLine = lineList[m]

					}

				}
				//fmt.Println(len(referenceLine))
				indexList[len(indexList)-1] = len(referenceLine)
				for m, _ := range indexList {
					if m != (len(indexList) - 1) {
						titleLeft := indexList[m]
						titleRight := indexList[m+1]
						titleList = append(titleList, referenceLine[titleLeft:titleRight])
					}
				}
				fmt.Println(titleList)
				for m, n := range titlesList {

					if len(n) < maxNumber {
						for x, y := range indexList {
							if x == len(indexList)-1 {
								break
							}
							if lineList[m][y] != ' ' {
								leftIndex := y
								rightIndex := indexList[x+1]
								//fmt.Println(leftIndex)
								//fmt.Println(rightIndex)
								tmpString := lineList[m][leftIndex:rightIndex]
								titleList[x] = tmpString + titleList[x]
								//fmt.Println(titleList)
							}
						}
					}
				}

			}
		} else if flag == false {
			titlesList = append(titlesList, ss)
			lineList = append(lineList, line)

		}

	}

}

// handle scenario 1 using json configuration file
type Config struct {
	Head      string
	NumOfSkip int
}

func loadConfig(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func configScenario1() {
	f, _ := os.Open("config.json")
	config, _ := gonfig.FromJson(f)
	head, _ := config.GetString("head", nil)
	numOfSkip, _ := config.GetInt("numOfSkip", 0)

	headList := strings.Split(head, ",")

	mapResult := make(map[string][]map[string]string)
	key := "Job"

	out, _ := exec.Command("cat", "qstat.out").Output()
	input := string(out)
	fmt.Println(input)
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
}
func main() {
	configScenario1()

}
