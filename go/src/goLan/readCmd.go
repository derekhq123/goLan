package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"reflect"
)

func execCommand(commandName string, params string) bool {
	cmd := exec.Command(commandName, params)
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
		fmt.Println(reflect.TypeOf(line))
	}

	cmd.Wait()

	return true
}

func main() {
	para := "-l"
	execCommand("ls", para)
}
