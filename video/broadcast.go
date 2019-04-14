package video

import (
	"arpiController/video/mockImage"
	"arpiController/video/rpiCamera"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	boundary            = flag.String("boundary", "--BOUNDARY", "boundary marker")
	withRaspberryCamera = flag.Bool("withRpiCamera", false, "Specifies if raspberry camera should be used for video stream")
	data                = &threadSafeSlice{
		workers: make([]*worker, 0, 1),
	}
)

type worker struct {
	source chan []byte
	first  bool
	done   bool
}

func Broadcast() {
	log.Println("Start broadcasting")
	c := make(chan []byte)
	go broadcaster(c)
	if *withRaspberryCamera {
		go rpiCamera.GenerateImage(c, *boundary)
	} else {
		go mockImage.GenerateImage(c, *boundary)
	}
}

func StreamTo(w io.Writer, closed <-chan bool) {
	wk := &worker{
		source: make(chan []byte),
		first:  true,
	}
	fmt.Fprintf(os.Stderr, "created %p\n", wk)
	data.push(wk)
	lastTime := time.Now()
loop:
	for {
		select {
		case s, ok := <-wk.source:
			fmt.Printf("Received new image at: %v\n", time.Since(lastTime))
			lastTime = time.Now()
			if !ok {
				break loop
			}
			if !wk.first {
				w.Write([]byte("\r\n"))
			} else {
				wk.first = false
			}
			w.Write(s)
		case <-closed:
			wk.done = true
		}
	}
}

func broadcaster(ch chan []byte) {
	for {
		msg := <-ch
		data.iter(func(w *worker) bool {
			if w.done {
				fmt.Fprintf(os.Stderr, "done %p\n", w)
				close(w.source)
				return true
			} else {
				w.source <- msg
				return false
			}
		})
	}
}
