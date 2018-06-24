package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/hoisie/mustache"
	"log"
	"arpiController/video"
	"arpiController/scripts"
	"strconv"
	"io/ioutil"
)

const (
	address = "0.0.0.0:8080"
	motorsPort = 8081
)

func init() {
	video.Broadcast()
	scripts.RunScript("motors", motorsPort)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", MainHandler)
	r.HandleFunc("/video", VideoHandler)
	buildApiRouter(r.PathPrefix("/api").Subrouter())
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Starting server at address " + address)
	log.Fatal(http.ListenAndServe(address, r))
}
func VideoHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary=--BOUNDARY")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	if c, ok := writer.(http.CloseNotifier); ok {
		video.StreamTo(writer, c.CloseNotify())
	}
}

func buildApiRouter(router *mux.Router) {
	router.HandleFunc("/updateState/{state}", func(writer http.ResponseWriter, request *http.Request) {
		state := mux.Vars(request)["state"]
		log.Println("Changing to state " + state)
		writer.WriteHeader(200)
		resp, err := http.Get("http://localhost:" + strconv.Itoa(motorsPort) + "/" + state)
		if err != nil {
			log.Println(err.Error())
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
		}
		if string(respBody) != "ok" {
			log.Println(string(respBody))
		}
	})
}

func MainHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(mustache.RenderFile("static/templates/main.html")))
}