package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/salopensource/salversion/pkg/salversion"
	log "github.com/sirupsen/logrus"
)

var Version string

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		salversion.PostHandler(w, r, Version)
	}).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		salversion.GetHandler(w, r, Version)
	}).Methods("GET")

	http.Handle("/", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s", port)
	}
	go getVersion()
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error(err)
	}
}

func getVersion() {
	ctx := context.Background()
	for {
		Version, err := salversion.GetSalVersion(ctx)
		if err != nil {
			log.Error(err)
		}
		log.Info(Version)
		time.Sleep(1 * time.Hour)
	}
}
