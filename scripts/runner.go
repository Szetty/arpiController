package scripts

import (
	"os/exec"
	"strconv"
	"log"
)

func RunScript(name string, port int) {
	cmd := exec.Command("python", "./scripts/" + name + ".py", "-p", strconv.Itoa(port))
	log.Println("Starting motor controller")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}