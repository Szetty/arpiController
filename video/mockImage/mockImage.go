package mockImage

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"time"
)

var newImageInterval = flag.Int("interval", 5, "interval in which a new image is generated")

func GenerateImage(ch chan []byte, boundary string) {
	buffer := new(bytes.Buffer)
	imgbuffer := new(bytes.Buffer)
	m := image.NewRGBA(image.Rect(0, 0, 256, 256))
	index := 0
	for {
		generateNextBytes(index, m, buffer, imgbuffer, boundary)
		ch <- buffer.Bytes()
		index++
		buffer.Reset()
		time.Sleep(time.Duration(*newImageInterval) * time.Millisecond)
	}
}

func generateNextBytes(index int, m *image.RGBA, buffer *bytes.Buffer, imgbuffer *bytes.Buffer, boundary string) {
	// generate
	x := index % m.Bounds().Max.X
	y := index / m.Bounds().Max.X
	m.Set(x, y, color.RGBA{255, 0, 255, 255})
	jpeg.Encode(imgbuffer, m, nil)

	// output
	fmt.Fprintf(buffer, "%s\r\n", boundary)
	fmt.Fprintf(buffer, "Content-Type: image/jpeg\r\n")
	fmt.Fprintf(buffer, "Content-Length: %d\r\n", imgbuffer.Len())
	buffer.Write([]byte("\r\n"))
	imgbuffer.WriteTo(buffer)
}
