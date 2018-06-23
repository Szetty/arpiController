package video

import (
	"github.com/dhowden/raspicam"
	"time"
	"fmt"
	"os"
)

const (
	imageLocation = "static/current.jpg"
)

var (
	cmd = raspicam.NewStill()
)

func init() {
	os.Create(imageLocation)
}

func ImageHandler() {
	for {
		f, _ := os.Open(imageLocation)
		errCh := make(chan error)
		go writeErrors(errCh)
		raspicam.Capture(cmd, f, errCh)
		time.Sleep(time.Second)
	}
}

func writeErrors(errCh chan error) {
	for x := range errCh {
		fmt.Fprintf(os.Stderr, "%v\n", x)
	}
}
