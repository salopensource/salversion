package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/salopensource/salversion/pkg/salversion"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", salversion.PostHandler).Methods("POST")
	r.HandleFunc("/", salversion.GetHandler).Methods("GET")
	http.Handle("/", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error(err)
	}
}
