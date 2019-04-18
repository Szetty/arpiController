package main

import (
	"arpiController/scripts"
	"arpiController/video"
	"flag"
	"github.com/gorilla/mux"
	"github.com/hoisie/mustache"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	address    = "0.0.0.0:8080"
	motorsPort = 8081
)

var (
	withMotors = flag.Bool("withMotors", false, "Specifies if explorer hat motors are present")
)

func init() {
	flag.Parse()
	if *withMotors {
		scripts.RunScript("motors", motorsPort)
	}
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

func VideoHandler(wr http.ResponseWriter, r *http.Request) {
	video.GetVideo(wr, r)
}

func buildApiRouter(router *mux.Router) {
	router.HandleFunc("/updateState/{state}", func(writer http.ResponseWriter, request *http.Request) {
		state := mux.Vars(request)["state"]
		log.Println("Changing to state " + state)
		writer.WriteHeader(200)
		if *withMotors {
			sendMotorRequest(state)
		}
	})
}

func sendMotorRequest(state string) {
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
}

func MainHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(mustache.RenderFile("static/templates/main.html")))
}
