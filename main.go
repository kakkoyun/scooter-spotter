package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	port        int
	gracePeriod int // in seconds

	logger *log.Logger // default Logger
)

func init() {
	// Initialize global dependencies
	logger = log.New(os.Stdout, "http: ", log.LstdFlags)
}

func newRouter() http.Handler {
	router := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
		w.Write([]byte("\n"))
	}))

	router.Handle("/_status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
		w.Write([]byte("\n"))
	}))

	return router
}

func main() {
	// 1. Parse arguments
	flag.IntVar(&port, "port", 4000, "host port to listen")
	flag.IntVar(&gracePeriod, "grace-period", 5, "grace period to wait for connections to drain, in secods")
	flag.Parse()

	// 2. Start server
	srv := newServer(
		port,
		newRouter(),
		time.Second*time.Duration(gracePeriod),
		logger,
	)

	logger.Println("Server is starting")
	if err := srv.listenAndServe(); err != nil {
		log.Fatalln(fmt.Errorf(" %+v", err))
	}
}
