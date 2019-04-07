package main

import (
	"flag"
	"fmt"
	"log"
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
