package scripts

import (
	"log"
	"os/exec"
	"strconv"
)

func RunScript(name string, port int) {
	cmd := exec.Command("python", "./scripts/"+name+".py", "-p", strconv.Itoa(port))
	log.Println("Starting motor controller")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	go monitorScript(cmd)
}

func monitorScript(cmd *exec.Cmd) {
	err := cmd.Wait()
	if err != nil {
		log.Fatal(err.Error())
	}
}
