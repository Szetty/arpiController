package rpiCamera

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func RunServer(port int) {
	args := fmt.Sprintf(`-o "output_http.so -w ./www -p %d" -i "input_raspicam.so -vf -hf"`, port)
	cmd := exec.Command("mjpg_streamer", args)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Println("Starting raspberry pi camera streamer")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	go monitorServer(cmd)
}

func monitorServer(cmd *exec.Cmd) {
	err := cmd.Wait()
	if err != nil {
		log.Fatalf("Streamer exited with: %v", err)
	}
}
