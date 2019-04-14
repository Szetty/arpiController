package rpiCamera

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

var buffersize = flag.Int("buffersize", 4096, "buffer size")

func GenerateImage(ch chan []byte, boundary string) {
	readbuffer := make([]byte, *buffersize)
	writebuffer := new(bytes.Buffer)
	for {
		n, err := os.Stdin.Read(readbuffer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gen read err %v\n", err)
			break
		}
		ProcessData(readbuffer, n, func(image []byte) {
			// header
			fmt.Fprintf(writebuffer, "%s\r\n", boundary)
			writebuffer.Write([]byte("Content-Type: image/jpeg\r\n"))
			fmt.Fprintf(writebuffer, "Content-Length: %d\r\n", len(image))
			writebuffer.Write([]byte("\r\n"))
			// image
			writebuffer.Write(image)
			// make a copy to send over channel
			cp := make([]byte, writebuffer.Len())
			copy(cp, writebuffer.Bytes())
			writebuffer.Reset()
			// send!
			ch <- cp
		})
	}
}
