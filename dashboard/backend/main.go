package main

import (
	"flag"
	"net/http"
	"os"

	"strconv"

	"github.com/fusion-app/fusion-app/dashboard/backend/handler"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultListenPort  int    = 8080
)

var (
	frontendDir string
	listenPort  int
)

func main() {
	flag.IntVar(&listenPort, "port", DefaultListenPort, `port this server listen to`)
	flag.Parse()

	h, err := handler.NewAPIHandler(frontendDir)
	if err != nil {
		log.Fatalf("Failed to create route handler: %v", err)
	}
	router := handler.NewRouter(h)

	http.Handle("/", router)

	p := ":" + strconv.Itoa(listenPort)
	log.Infof("fusion-app server listens on: %v", p)
	if err = http.ListenAndServe(p, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
