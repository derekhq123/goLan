package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/creamdog/gonfig"
	"gitlab.hpls.local/ppsc/agent-scheduler/config"
	"io"
	"io/ioutil"
	"net/http"

	"os/exec"
	"strings"
)

type QueueTypeRequestBody struct {
	QueueType string
}

// send the queue information on system
func QueryQueue(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Querying Queue...\n")
	var submitJobResponse SubmitJobResponse

	cmd := exec.Command("qstat", "-q")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	// Execute command
	printCommand(cmd)
	err1 := cmd.Run()
	printError(err1)
	printOutput(cmdOutput.Bytes())

	appId := strings.Split(string(cmdOutput.Bytes()), ".")
	fmt.Fprintf(w, "%s", appId[0])
	submitJobResponse.JobId = appId[0]

	// TODO Return
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(submitJobResponse); err != nil {
		panic(err)
	}
}

// send the info for all the jobs of userid
func QueryQueueByUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Listing Queue...\n")

	var request QueryQueueByUserRequest

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	cmd := exec.Command("qstat", "-u", request.User)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	// Execute command
	printCommand(cmd)
	err = cmd.Run() // will wait for command to return
	printError(err)
	printOutput(cmdOutput.Bytes())

	todos := Todos{
		Todo{Name: "qstat"},
	}

	json.NewEncoder(w).Encode(todos)

}

// send the queue limits for all queues
func QueryQueueLimit(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("flagvar has value ", main.flagvar)
	// get the command and parameter from config file
	command, _ := config.ConfigFile.GetString("QueryQueueLimit/command", nil)
	parameter, _ := config.ConfigFile.GetString("QueryQueueLimit/parameter", nil)

	input, _ := exec.Command(command, parameter).Output()

	// TODO Return
	w.Header().Set("Content-Type", "application/json")

	//w.WriteHeader(http.StatusCreated)

	//using table format parser and send
	if err := json.NewEncoder(w).Encode(TableParser(string(input), "QueryQueueLimit")); err != nil {

		panic(err)
	}
}

// send the summary information about PBS server
func ServerSummary(w http.ResponseWriter, r *http.Request) {

	command, _ := config.ConfigFile.GetString("ServerSummary/command", nil)
	parameter, _ := config.ConfigFile.GetString("ServerSummary/parameter", nil)

	input, _ := exec.Command(command, parameter).Output()

	// TODO Return
	w.Header().Set("Content-Type", "application/json")

	//w.WriteHeader(http.StatusCreated)
	//using table format parser and send
	if err := json.NewEncoder(w).Encode(TableParser(string(input), "qstat-B")); err != nil {

		panic(err)
	}

}

// send all information known about queueid
func QueueQueryByQueueId(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request QueueTypeRequestBody

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	// get command and parameter from config file
	command, _ := config.ConfigFile.GetString("QueueQueryByQueueId/command", nil)
	parameter, _ := config.ConfigFile.GetString("QueueQueryByQueueId/parameter", nil)

	cmdArgs := strings.Split(parameter, " ")

	if request.QueueType != "" {
		cmdArgs = append(cmdArgs, request.QueueType)
	}

	input, _ := exec.Command(command, cmdArgs...).Output()
	//using key-value parser and send
	if err := json.NewEncoder(w).Encode(KeyValueParser(string(input), "QueueQueryByQueueId")); err != nil {
		panic(err)
	}
}
