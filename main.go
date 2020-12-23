package main

import (
	"fmt"
	"log"
	"marlinctl/commands"
	"marlinctl/util"
	"os"
	"os/exec"
	"strings"
)

func main() {

	if !commands.IsRoot() {
		fmt.Println("requires root permissions. Please run with sudo")
		return
	}
	if !commands.IsCommandAvailable("supervisorctl") {
		fmt.Println("supervisorctl not installed!!! Please install and try again")
		return
	}
	out, _ := exec.Command("sudo", "supervisorctl", "status").Output()
	if strings.Contains(string(out), "no such file") {
		// supervisor service not running
		fmt.Println("supervisor service not running!!! Please start the service and try again")
		return
	}
	if util.CheckAndUpdate() { // if update is found then new process with same arguments will be spawned
		return
	}
	err := commands.App.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
