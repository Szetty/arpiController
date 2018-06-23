package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/hoisie/mustache"
	"log"
	"arpiController/video"
)

const (
	address = "localhost:8080"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", MainHandler)
	r.HandleFunc("/video", video.VideoHandler)
	buildApiRouter(r.PathPrefix("/api").Subrouter())
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Starting server at address " + address)
	log.Fatal(http.ListenAndServe(address, r))
}

func buildApiRouter(router *mux.Router) {
	router.HandleFunc("/updateState/{state}", func(writer http.ResponseWriter, request *http.Request) {
		state := mux.Vars(request)["state"]
		log.Println("Changing to state " + state)
		writer.WriteHeader(200)
	})
}

func MainHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(mustache.RenderFile("static/templates/main.html")))
}