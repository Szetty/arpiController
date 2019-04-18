package video

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	withRaspberryCamera   = flag.Bool("withRpiCamera", false, "Specifies if raspberry camera should be used for video stream")
	raspberryStreamerPort = flag.Int("rpiStreamerPort", 8000, "Specifies on what port should the streamer run")
)

func init() {
	flag.Parse()
}

func GetVideo(wr http.ResponseWriter, r *http.Request) {
	if *withRaspberryCamera {
		getVideoFromStreamer(wr, r)
	} else {
		wr.Write([]byte{})
		//getVideoFromStreamer(wr, r)
	}
}

func getVideoFromStreamer(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	client := &http.Client{}
	log.Printf("%+v", r)
	hostWithoutPort := strings.Split(r.Host, ":")[0]
	// hostWithoutPort := 192.168.0.127
	url := "http://" + hostWithoutPort + ":" + strconv.Itoa(*raspberryStreamerPort) + "/?action=stream"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Video request creation: %v", err)
	}
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}
	resp, err = client.Do(req)
	r.Body.Close()

	// combined for GET/POST
	if err != nil {
		log.Printf("Video error: %v", err)
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	resp.Body.Close()
}
