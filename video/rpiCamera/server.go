package rpiCamera

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func RunServer(port int) {
	args := []string{
		"-o",
		fmt.Sprintf(`"output_http.so -w ./www -p %d"`, port),
		"-i",
		`"input_raspicam.so -vf -hf"`,
	}
	cmd := exec.Command("mjpg_streamer", args...)
	cmd.Env = append(os.Environ(), "LD_LIBRARY_PATH=/home/pi/mjpg-streamer/mjpg-streamer-experimental")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Command is: %+v", cmd)
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
