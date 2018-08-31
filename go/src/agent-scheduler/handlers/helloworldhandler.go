package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "If my armor breaks, I fuse it back together...\n")
}

func TestJob(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Submitting Job!!!!!...\n")

	cmd := exec.Command("echo", "When the seagulls,", "follow the trawlers")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	// Execute command
	printCommand(cmd)
	err1 := cmd.Run()
	printError(err1)
	printOutput(cmdOutput.Bytes())

}
