package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/creamdog/gonfig"
)

type Todo struct {
	Id          int `json:"id"`
	Name        string
	JobFilename string
}

type QueryQueueByUserRequest struct {
	User string
}

type Todos []Todo

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing here: %s\n", strings.Join(cmd.Args, " "))
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func StringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}

func handle_childchild(v1 bytes.Buffer, prev_key string,
	prev_value string, deli_ccs string, deli_ccsg string) string {

	var s, store_temp bytes.Buffer
	var t []string
	rnd := ""

	//The child header that belongs to this childchild
	s.WriteString(prev_value)
	s.WriteString(v1.String())

	childchildarray := strings.Split(s.String(), deli_ccs)

	for i := 0; i < len(childchildarray); i++ {

		t = strings.Split(childchildarray[i], deli_ccsg)
		temp_header := strings.TrimSpace(t[0])
		store_temp.WriteString("\"" + temp_header + "\":\"" + strings.TrimSpace(t[1]) + "\",")
	}

	// just to remove the last coma
	rnd = TrimSuffix(store_temp.String(), ",")
	final := "\"" + strings.ToLower(string(prev_key)) + "\":{" + rnd + "},"

	return final
}

//used to parse the key-value format of PBS command output to json format

func KeyValueParser(mycmd string, keyheader string) []map[string]interface{} {

	// f, _ := os.Open("configmockdata.json")

	// defer f.Close()
	// config, _ := gonfig.FromJson(f)

	parent_deli, _ := config.ConfigFile.GetString(fmt.Sprintf("%s%s", keyheader, "/parent"), ":")
	child_deli, _ := config.ConfigFile.GetString(fmt.Sprintf("%s%s", keyheader, "/child"), " = ")
	child_child_set, _ := config.ConfigFile.GetString(fmt.Sprintf("%s%s", keyheader, "/child-child/set"), ",")
	child_child_single, _ := config.ConfigFile.GetString(fmt.Sprintf("%s%s", keyheader, "/child-child/single"), " ")

	var actual_json bytes.Buffer
	var temp []string
	var save_header = ""
	var storeChildChild bytes.Buffer

	header := regexp.MustCompile(`^[^\s][a-zA-z]+`)
	child := regexp.MustCompile(`^[\s\p{Zs}]{4,7}[a-zA-Z]+`)
	childchild := regexp.MustCompile(`^([\s\p{Zs}]{1}|[\s\p{Zs}]{8})[a-zA-Z:]+`)
	childchildSituation := false
	headerExists := false
	prev_Key := ""
	prev_Value := ""

	lines, _ := StringToLines(mycmd)

	//Process output from command
	for _, line := range lines {
		//fmt.Println(line)
		if header.MatchString(line) {
			// It must be the header
			if !headerExists {
				temp = strings.Split(line, parent_deli)
				save_header = strings.ToLower(strings.TrimSpace(temp[0]))
				actual_json.WriteString("[{")
				//actual_json.WriteString("{\"test\":[{ \"jobid\":\"" + strings.ToLower(strings.Replace(line, " ", "", -1)) + "\",")
				if len(temp) > 1 {
					actual_json.WriteString("\"" + save_header + "\":\"" + strings.TrimSpace(temp[1]) + "\",")
				} else {
					actual_json.WriteString("\"" + save_header + "\":\"\",")
				}

				headerExists = true
			} else {
				//Add previous value first
				if prev_Key != "" && prev_Value != "" && !childchildSituation {
					actual_json.WriteString("\"" + strings.ToLower(prev_Key) + "\":\"" + prev_Value + "\",")
					prev_Key = ""
					prev_Value = ""
				} else if prev_Key != "" && prev_Value != "" && childchildSituation {
					actual_json.WriteString(handle_childchild(storeChildChild,
						strings.ToLower(prev_Key), prev_Value, child_child_set, child_child_single))
				}

				// A new header has been detected, so need to wrap up the previous key
				temp_use := TrimSuffix(actual_json.String(), ",")
				actual_json.Reset()
				actual_json.WriteString(temp_use)
				actual_json.WriteString("},")

				temp = strings.Split(line, parent_deli)
				save_header = strings.ToLower(strings.TrimSpace(temp[0]))
				actual_json.WriteString("{")
				if len(temp) > 1 {
					actual_json.WriteString("\"" + save_header + "\":\"" + strings.TrimSpace(temp[1]) + "\",")
				} else {
					actual_json.WriteString("\"" + save_header + "\":\"\",")
				}
			}
		} else if childchild.MatchString(line) {
			childchildSituation = true
			//Inner array detected
			//The previous value is an item, so store it
			storeChildChild.WriteString(strings.TrimSpace(line))

		} else if child.MatchString(line) {
			if childchildSituation {
				actual_json.WriteString(handle_childchild(storeChildChild,
					prev_Key, prev_Value, child_child_set, child_child_single))
				prev_Key = ""
				prev_Value = ""
				storeChildChild.Reset()
			}
			//Add previous value first
			if prev_Key != "" && prev_Value != "" && !childchildSituation {
				actual_json.WriteString("\"" + strings.ToLower(prev_Key) + "\":\"" + prev_Value + "\",")
				prev_Key = ""
				prev_Value = ""
			}
			// It must be the child
			temp = strings.Split(line, child_deli)
			prev_Key = strings.TrimSpace(temp[0])
			prev_Value = strings.TrimSpace(temp[1])

			childchildSituation = false
		}
	}

	//Add previous value first
	if prev_Key != "" && prev_Value != "" && !childchildSituation {
		actual_json.WriteString("\"" + strings.ToLower(prev_Key) + "\":\"" + prev_Value + "\",")
		prev_Key = ""
		prev_Value = ""
	} else if prev_Key != "" && prev_Value != "" && childchildSituation {
		actual_json.WriteString(handle_childchild(storeChildChild,
			strings.ToLower(prev_Key), prev_Value, child_child_set, child_child_single))
	}

	temp_use := TrimSuffix(actual_json.String(), ",")
	actual_json.Reset()
	actual_json.WriteString(temp_use)
	actual_json.WriteString("}]")

	in := []byte(actual_json.String())

	var raw []map[string]interface{}

	json.Unmarshal(in, &raw)
	//out, _ := json.MarshalIndent(raw, "", "\t")

	//return mycmd
	return raw
}

//used to parse the table format of PBS command output to json format
func TableParser(mycmd string, keyheader string) map[string][]map[string]string {
	// f, _ := os.Open("configmockdata.json")
	// defer f.Close()
	// config, _ := gonfig.FromJson(f)

	head, _ := config.ConfigFile.GetString(fmt.Sprintf("%s%s", keyheader, "/head"), nil)
	numOfSkip, _ := config.ConfigFile.GetInt(fmt.Sprintf("%s%s", keyheader, "/numOfSkip"), 0)
	lastCol, _ := config.ConfigFile.GetInt(fmt.Sprintf("%s%s", keyheader, "/lastCol"), 0)

	headList := strings.Split(head, ",")

	mapResult := make(map[string][]map[string]string)
	key := "Job"

	scanner := bufio.NewScanner(strings.NewReader(mycmd))

	index := 0
	var tmpList []map[string]string

	for scanner.Scan() {
		line := scanner.Text()
		if index == numOfSkip {
			if line != "" {
				tmpMap := make(map[string]string)
				ss := strings.Fields(line)
				var combineLastColText bytes.Buffer
				rmbLastK := 0

				for k, v := range ss {

					if (k + 1) > lastCol {
						combineLastColText.WriteString(v + " ")
					} else {
						tmpMap[headList[k]] = v
					}
					rmbLastK = k
				}

				if rmbLastK > lastCol {
					tmpMap[headList[lastCol]] = combineLastColText.String()
				}
				tmpList = append(tmpList, tmpMap)

			}
		} else {
			index += 1
		}

	}
	mapResult[key] = tmpList
	//fmt.Println(mapResult)
	//output, _ := json.Marshal(mapResult)
	//fmt.Println(string(output))
	return mapResult
}
