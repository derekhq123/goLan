package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/creamdog/gonfig"
	"github.com/golang/glog"
	"gitlab.hpls.local/ppsc/agent-scheduler/config"
	"gitlab.hpls.local/ppsc/agent-scheduler/handlers"

	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	// BANNER is what is printed for help/info output.
	BANNER = `
 Server to accept REST request and make system call.
 Version: %s
`
)

var (
	port string
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.StringVar(&port, "p", "8080", "port for server to run on")
	// flag.String("log_dir", "log/", "directory to strore logs")
	// flag.Usage = func() {
	// 	fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.version))
	// 	flag.PrintDefaults()
	// }

	// get the configuration file name in commmand
	configAddress := flag.String("c", "configmockdata.json", "Config file address")
	flag.Parse()
	configName := *configAddress
	config.SetConfig(configName)

}

func handleRestCall() {

	fmt.Println("Starting application...")

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func main() {
	// c := make(chan struct{})
	// go handleRestCall()
	// //sendJobInfo()
	// go sendQueueInfo()

	// <-c
	glog.Infof("Application started.")

	router := NewRouter()
	addrs, _ := net.InterfaceAddrs()

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Printf("Port:%s, IP:%s", port, ipnet.IP.String())
				//os.Stdout.WriteString(ipnet.IP.String() + "\n")
			}
		}
	}
	glog.Infof("Serving at port %v.", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
	sendJobInfo(*configAddress)
	sendQueueInfo(*configAddress)

	glog.Flush()

}

/* run PBS qstat-a command and send the parsed output about job information
   to the URL per interval time
*/
func sendJobInfo() {

	url, _ := config.ConfigFile.GetString("url", nil)
	timeInterval, _ := config.ConfigFile.GetInt("qstat-a/interval", 0)
	command, _ := config.ConfigFile.GetString("qstat-a/command", nil)
	parameter, _ := config.ConfigFile.GetString("qstat-a/parameter", nil)

	for range time.Tick(time.Duration(timeInterval) * time.Second) {
		fmt.Println("Sending Job info...")

		input, _ := exec.Command(command, parameter).Output()

		output := handlers.TableParser(string(input), "qstat-a")

		b, err := json.Marshal(output)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(string(b))
		body := bytes.NewBuffer([]byte(b))
		res, err := http.Post(url, "application/json;charset=utf-8", body)
		if err != nil {
			panic(err)
		}
		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s", result)
	}
}

/* run PBS qstat-q command and send the parsed output about queue information
   to the URL per interval time
*/
func sendQueueInfo() {

	url, _ := config.ConfigFile.GetString("url", nil)
	timeInterval, _ := config.ConfigFile.GetInt("qstat-q/interval", 0)
	command, _ := config.ConfigFile.GetString("qstat-q/command", nil)
	parameter, _ := config.ConfigFile.GetString("qstat-q/parameter", nil)

	for range time.Tick(time.Duration(timeInterval) * time.Second) {
		fmt.Println("Sending Job info...")

		input, _ := exec.Command(command, parameter).Output()

		output := handlers.TableParser(string(input), "qstat-q")

		b, err := json.Marshal(output)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(string(b))
		body := bytes.NewBuffer([]byte(b))
		res, err := http.Post(url, "application/json;charset=utf-8", body)
		if err != nil {
			panic(err)
		}
		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s", result)
	}
}
