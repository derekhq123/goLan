package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/creamdog/gonfig"
	"github.com/golang/glog"
	"gitlab.hpls.local/ppsc/agent-scheduler/config"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var currentId int
var todos Todos
var pbsscripts []PBSScript
var submitJobResponse SubmitJobResponse

func RepoCreateTodo(t Todo) Todo {
	currentId += 1
	t.Id = currentId
	todos = append(todos, t)
	return t
}

func CreateJobScript(p PBSScript) PBSScript {
	currentId += 1
	p.Id = currentId
	pbsscripts = append(pbsscripts, p)
	return p
}

type JobRequestBody struct {
	JobId string
}

type QueueJobRequestBody struct {
	UserId string
}

type NodeInfoRequestBody struct {
	JobId string
}

type SubmitJobResponse struct {
	JobId string `json:"jobid"`
}

type TerminateJobResponse struct {
	ErrCode int
	Message string
}

//type SubmitCmdResponse struct {
//	Command string      `json:"cmd"`
//	ParseType string    `json:"parsetype"`
//}

func SubmitJob(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var todo Todo
	var submitJobResponse SubmitJobResponse

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &todo); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	t := RepoCreateTodo(todo)

	cmd := exec.Command("cat", t.JobFilename)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	// Execute command
	printCommand(cmd)
	err1 := cmd.Run()
	printError(err1)
	//printOutput(cmdOutput.Bytes())

	appId := strings.Split(string(cmdOutput.Bytes()), ".")
	fmt.Fprintf(w, "%s", appId[0])
	submitJobResponse.JobId = appId[0]

	fmt.Println(submitJobResponse)

	// Return
	//w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(submitJobResponse); err != nil {
		panic(err)
	}
}

func CreateAndSubmitJobScript(w http.ResponseWriter, r *http.Request) {
	var pbsScript PBSScript
	var submitJobResponse SubmitJobResponse

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &pbsScript); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// First check if MetaID exist.
	if pbsScript.MetaID == "" {
		// w.WriteHeader(http.StatusInternalServerError)
		glog.Errorf("MetaID does not exist")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	glog.Infof("Generating PBS Script")
	p := CreateJobScript(pbsScript)

	var buffer bytes.Buffer

	buffer.WriteString("#!/bin/bash\n\n")
	glog.Infof("Inserting Job Name: " + pbsScript.Name)
	buffer.WriteString("#PBS -N " + pbsScript.Name + "\n")

	glog.Infof("Allocation of resources : select=" + pbsScript.Chunk)
	buffer.WriteString("#PBS -l select=")
	if pbsScript.Chunk != "" {
		buffer.WriteString("" + pbsScript.Chunk)
	} else {
		buffer.WriteString("1") //defaults to 1 if none specified
	}
	if pbsScript.NCPUs != "" {
		glog.Infof("                         ncpus=" + pbsScript.NCPUs)
		buffer.WriteString(":ncpus=" + pbsScript.NCPUs)
	}
	if pbsScript.Mem != "" {
		glog.Infof("                         mem=" + pbsScript.Mem)
		buffer.WriteString(":mem=" + pbsScript.Mem)
	}
	if pbsScript.MPIProcs != "" {
		glog.Infof("                         mpiprocs" + pbsScript.MPIProcs)
		buffer.WriteString(":mpiprocs=" + pbsScript.MPIProcs)
	}
	if pbsScript.OMPThreads != "" {
		glog.Infof("                         ompthreads=" + pbsScript.OMPThreads)
		buffer.WriteString(":ompthreads=" + pbsScript.OMPThreads)
	}
	buffer.WriteString("\n")

	if pbsScript.WallTime != "" {
		glog.Infof("Assigning walltime :" + pbsScript.WallTime)
		buffer.WriteString("#PBS -l walltime=" + pbsScript.WallTime + "\n")
	}

	if pbsScript.QueueName != "" {
		glog.Infof("Queue :" + pbsScript.QueueName)
		buffer.WriteString("#PBS -q " + pbsScript.QueueName + "\n")
	}

	if pbsScript.OutputFile != "" {
		glog.Infof("OutputFile :" + pbsScript.OutputFile)
		buffer.WriteString("#PBS -o " + pbsScript.OutputFile + "\n")
	}

	if pbsScript.ErrorFile != "" {
		glog.Infof("ErrorFile :" + pbsScript.ErrorFile)
		buffer.WriteString("#PBS -e " + pbsScript.ErrorFile + "\n")
	}

	if pbsScript.Description != "" {
		glog.Infof("Job Description :" + pbsScript.Description)
		buffer.WriteString("#PBS -l software=\"" + pbsScript.Description + "\"\n")
	}

	if pbsScript.Module != nil {

		glog.Infof("Loading Module.")
		buffer.WriteString("\n\n# Load Module.\n")
		buffer.WriteString("module load")
		for i := range pbsScript.Module {
			fmt.Fprintf(w, pbsScript.Module[i])
			buffer.WriteString(" " + pbsScript.Module[i])
		}

	}

	if pbsScript.WorkDirectory != "" {
		glog.Infof("Change into working directory :" + pbsScript.WorkDirectory)

		buffer.WriteString("\n\n## Change into directory where qsub command was submitted.\n")
		buffer.WriteString("cd " + pbsScript.WorkDirectory)
		buffer.WriteString("\n")

	}

	// Optional : Pre-process
	if pbsScript.PreProcess != "" {
		glog.Infof("Pre-process added.")
		buffer.WriteString("\n\n# Pre-process stage.\n")
		buffer.WriteString(pbsScript.PreProcess)
	}

	if pbsScript.CustomScript != "" {
		glog.Infof("Actual running of application." + pbsScript.CustomScript)
		buffer.WriteString("\n\n# Now run the program.\n")
		buffer.WriteString(pbsScript.CustomScript)
	}

	// Optional : Post-process
	if pbsScript.PostProcess != "" {
		glog.Infof("Post-process added.")
		buffer.WriteString("\n\n# Post-process stage.\n")
		buffer.WriteString(pbsScript.PostProcess)
	}

	// Only generate file if MetaID exist
	if pbsScript.MetaID != "" {
		var filename string
		filename = pbsScript.WorkDirectory + pbsScript.MetaID + ".pbs"
		glog.Infof("Generating PBS script file: " + filename)
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		writer := bufio.NewWriter(f)
		fmt.Fprintf(writer, "%s", buffer.String())
		writer.Flush()

		glog.Infof("Submitting PBS script to pbs scheduler.")
		cmd := exec.Command("qsub", filename)
		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput

		// Execute command
		printCommand(cmd)
		err1 := cmd.Run()
		printError(err1)
		//	printOutput(cmdOutput.Bytes())

		// FIXME : real JobId
		appId := strings.Split(string(cmdOutput.Bytes()), ".")
		glog.Infof("Retrieving JobId" + appId[0])
		fmt.Fprintf(w, "%s", appId[0])
		submitJobResponse.JobId = appId[0]

		// For Mock
		strconv.Itoa(p.Id)
		// submitJobResponse.JobId = strconv.Itoa(p.Id)

	}

	// Return
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(submitJobResponse); err != nil {
		panic(err)
	}
}

func StartJobHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Starting Job...\n")

	cmd := exec.Command("cat", "samples.pbs")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	// Execute command
	printCommand(cmd)
	err := cmd.Run() // will wait for command to return
	printError(err)
	printOutput(cmdOutput.Bytes())

	todos := Todos{
		Todo{Name: "Start Job"},
	}

	json.NewEncoder(w).Encode(todos)

}

func TerminateJobHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request JobRequestBody
	var response TerminateJobResponse

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

	if request.JobId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmd := exec.Command("qdel", request.JobId)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	cmdError := &bytes.Buffer{}
	cmd.Stderr = cmdError

	// Execute command
	printCommand(cmd)
	err = cmd.Run() // will wait for command to return
	printError(err)
	printOutput(cmdOutput.Bytes())

	if cmdOutput.Bytes() != nil {
		response.ErrCode = 1
		fmt.Fprintf(w, "%s", cmdError.Bytes())
		response.Message = string(cmdError.Bytes())
	} else {
		response.ErrCode = 0
		response.Message = ("Job Terminating.")
	}

	// Return
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}

}

func ForcePurgeJobHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Force Purge Job...\n")
	var request JobRequestBody

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

	cmd := exec.Command("qdel", "-p", request.JobId)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	// Execute command
	printCommand(cmd)
	err = cmd.Run() // will wait for command to return
	printError(err)
	printOutput(cmdOutput.Bytes())

	todos := Todos{
		Todo{Name: "Terminate job"},
	}

	json.NewEncoder(w).Encode(todos)

}

// send queued job information
func GetQueueNJobInfo(w http.ResponseWriter, r *http.Request) {
	//get command and parameter from config file
	command, _ := config.ConfigFile.GetString("GetQueueNJobInfo/command", nil)
	parameter, _ := config.ConfigFile.GetString("GetQueueNJobInfo/parameter", nil)

	input, _ := exec.Command(command, parameter).Output()

	// TODO Return
	w.Header().Set("Content-Type", "application/json")

	//w.WriteHeader(http.StatusCreated)
	//using table parser and send
	if err := json.NewEncoder(w).Encode(TableParser(string(input), "GetQueueNJobInfo")); err != nil {
		panic(err)
	}
}

// send full information about the jobid
func JobQueryByJobId(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//get command and parameter from config file
	command, _ := config.ConfigFile.GetString("JobQueryByJobId/command", nil)
	parameter, _ := config.ConfigFile.GetString("JobQueryByJobId/parameter", nil)

	input, _ := exec.Command(command, parameter).Output()

	//w.WriteHeader(http.StatusCreated)
	//using key-value parser and send
	if err := json.NewEncoder(w).Encode(KeyValueParser(string(input), "JobQueryByJobId")); err != nil {
		panic(err)
	}
}

// send info about all jobs on system
func GetJobInfo(w http.ResponseWriter, r *http.Request) {
	//get command and parameter from config file
	command, _ := config.ConfigFile.GetString("GetJobInfo/command", nil)
	parameter, _ := config.ConfigFile.GetString("GetJobInfo/parameter", nil)

	input, _ := exec.Command(command, parameter).Output()

	w.Header().Set("Content-Type", "application/json")

	//w.WriteHeader(http.StatusCreated)
	//using table parser and send
	if err := json.NewEncoder(w).Encode(TableParser(string(input), "GetJobInfo")); err != nil {
		panic(err)
	}
}

//send the info for queued job of userid
func GetQueuedJobInfoByUserId(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request QueueJobRequestBody

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

	if request.UserId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get the command and parameter from config file

	command, _ := config.ConfigFile.GetString("GetQueuedJobInfoByUserId/command", nil)
	parameter, _ := config.ConfigFile.GetString("GetQueuedJobInfoByUserId/parameter", nil)

	cmdArgs := []string{parameter, request.UserId}

	input, _ := exec.Command(command, cmdArgs...).Output()
	//using table parser and send
	if err := json.NewEncoder(w).Encode(TableParser(string(input), "GetQueuedJobInfoByUserId")); err != nil {
		panic(err)
	}
}

//send the info for all the jobs of userid
func GetJobInfoByUserId(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request QueueJobRequestBody

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

	if request.UserId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get the command and parameter from config file
	command, _ := config.ConfigFile.GetString("GetJobInfoByUserId/command", nil)
	parameter, _ := config.ConfigFile.GetString("GetJobInfoByUserId/parameter", nil)

	cmdArgs := []string{parameter, request.UserId}

	input, _ := exec.Command(command, cmdArgs...).Output()
	//using table parser and send
	if err := json.NewEncoder(w).Encode(TableParser(string(input), "GetJobInfoByUserId")); err != nil {
		panic(err)
	}
}

// send the info about all jobs with status comment
func GetJobInfoWithComment(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//get command and parameter from config file
	command, _ := config.ConfigFile.GetString("GetJobInfoWithComment/command", nil)
	parameter, _ := config.ConfigFile.GetString("GetJobInfoWithComment/parameter", nil)

	cmdArgs := strings.Split(parameter, " ")

	input, _ := exec.Command(command, cmdArgs...).Output()
	//using table parser and send
	if err := json.NewEncoder(w).Encode(TableParser(string(input), "GetJobInfoWithComment")); err != nil {
		panic(err)
	}
}

//send the node info on which jobid is running
func GetNodeInfoByJobId(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request NodeInfoRequestBody

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

	if request.JobId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f, _ := os.Open("")
	//get the command and parameter from config file
	command, _ := config.ConfigFile.GetString("GetNodeInfoByJobId/command", nil)
	parameter, _ := config.ConfigFile.GetString("GetNodeInfoByJobId/parameter", nil)

	cmdArgs := strings.Split(parameter, " ")
	cmdArgs = append(cmdArgs, request.JobId)

	input, _ := exec.Command(command, cmdArgs...).Output()
	//using table parser and send
	if err := json.NewEncoder(w).Encode(TableParser(string(input), "GetNodeInfoByJobId")); err != nil {
		panic(err)
	}
}

func GetAllPbsNodesInfo(w http.ResponseWriter, r *http.Request) {

	command, _ := config.ConfigFile.GetString("GetAllPbsNodesInfo/command", nil)
	parameter, _ := config.ConfigFile.GetString("GetAllPbsNodesInfo/parameter", nil)

	input, _ := exec.Command(command, parameter).Output()

	// TODO Return
	w.Header().Set("Content-Type", "application/json")

	//w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(KeyValueParser(string(input), "GetAllPbsNodesInfo")); err != nil {
		panic(err)
	}
}

func GetClusterInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	command, _ := config.ConfigFile.GetString("GetAllPbsNodesInfo/command", nil)
	parameter, _ := config.ConfigFile.GetString("GetAllPbsNodesInfo/parameter", nil)

	input, _ := exec.Command(command, parameter).Output()

	//w.WriteHeader(http.StatusCreated)
	result := KeyValueParser(string(input), "GetAllPbsNodesInfo")

	var main_json bytes.Buffer
	var nodesinfo bytes.Buffer
	rmb_cluster := ""
	totalnodes := 0
	busynodes := 0

	main_json.WriteString("[{")
	nodesinfo.WriteString("\"nodes\":[")
	for _, v := range result {

		if rmb_cluster == "" {
			//Fresh
			rmb_cluster = fmt.Sprintf("%v", v["cluster"])
		} else if rmb_cluster != v["cluster"] {
			//New cluster spotted
			temp_use := TrimSuffix(nodesinfo.String(), ",")
			nodesinfo.Reset()
			nodesinfo.WriteString(temp_use)
			nodesinfo.WriteString("]},{")

			main_json.WriteString("\"nodecountincluster\":" + strconv.Itoa(totalnodes) +
				",\"prctofnodeused\":" + strconv.FormatFloat(PercentOf(busynodes, totalnodes), 'f', 6, 64) + "," +
				"\"cluster\":\"" + rmb_cluster + "\"," +
				nodesinfo.String())

			//Reset to default
			rmb_cluster = fmt.Sprintf("%v", v["cluster"])
			totalnodes = 0
			busynodes = 0
			nodesinfo.Reset()
			nodesinfo.WriteString("\"nodes\":[")
		}

		totalnodes += 1

		if v["jobs"] != nil {
			busynodes += 1
		}

		res_avail_cpu, _ := strconv.ParseInt(v["resources_available.ncpus"].(string), 10, 32)
		res_assign_cpu, _ := strconv.ParseInt(v["resources_assigned.ncpus"].(string), 10, 32)
		freed_cpu := strconv.FormatInt(res_avail_cpu-res_assign_cpu, 10)

		res_avail_mem, _ := strconv.ParseInt(TrimSuffix(v["resources_available.mem"].(string), "kb"), 10, 32)
		res_assign_mem, _ := strconv.ParseInt(TrimSuffix(v["resources_assigned.mem"].(string), "kb"), 10, 32)
		freed_mem := strconv.FormatInt(res_avail_mem-res_assign_mem, 10)

		nodesinfo.WriteString("{\"name\":\"" + fmt.Sprintf("%v", v["mom"]) + "\"," +
			"\"resources_freed_ncpus\":\"" + freed_cpu + "\"," +
			"\"resources_freed_mem\":\"" + freed_mem + "kb\"},")
	}

	temp_use := TrimSuffix(nodesinfo.String(), ",")
	nodesinfo.Reset()
	nodesinfo.WriteString(temp_use)
	nodesinfo.WriteString("]}]")

	main_json.WriteString("\"nodecountincluster\":" + strconv.Itoa(totalnodes) +
		",\"prctofnodeused\":" + strconv.FormatFloat(PercentOf(busynodes, totalnodes), 'f', 6, 64) + "," +
		"\"cluster\":\"" + rmb_cluster + "\"," +
		nodesinfo.String())

	in := []byte(main_json.String())
	var raw []map[string]interface{}
	json.Unmarshal(in, &raw)
	//fmt.Println(totalnodes)
	//fmt.Println(busynodes)
	//fmt.Println(PercentOf(busynodes,totalnodes))

	if err := json.NewEncoder(w).Encode(raw); err != nil {
		panic(err)
	}
}

func PercentOf(current int, all int) float64 {
	percent := (float64(current) * float64(100)) / float64(all)
	return percent
}
