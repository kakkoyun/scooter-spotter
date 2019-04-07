package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// server will create a http.Server from the Go standard library
type server struct {
	s       *http.Server
	timeout time.Duration // duration to wait for graceful shutdown
}

// port int: Host por to listen
// handler http.Handler: Handler to invoke, http.DefaultServeMux if nil
// timeout time.Duration: Duration to wait for graceful shutdown
func newServer(port int, handler http.Handler, timeout time.Duration, logger *log.Logger) *server {
	return &server{
		s: &http.Server{
			Addr:     fmt.Sprintf(":%d", port),
			Handler:  handler,
			ErrorLog: logger,
		},
		timeout: timeout,
	}
}

// listenAndServe mirrors the function from the Go standard library, always returns a non-nil error
func (srv *server) listenAndServe() error {
	idleConnsClosed := make(chan struct{})
	srv.listenInterrupts(idleConnsClosed)

	if err := srv.s.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return fmt.Errorf("Server closed unexpectedly, %v", err)
		}
	}

	<-idleConnsClosed
	return nil
}

// listenInterrupts listens OS signals and triggers graceful shutdown
func (srv *server) listenInterrupts(idleConnsClosed chan<- struct{}) {
	// Listen for an OS interrupt signal. For example, catch CTRL-C
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-osSignals
		defer close(idleConnsClosed)

		done := make(chan struct{})
		go func() {
			// We received an interrupt signal, shutdown.
			if err := srv.s.Shutdown(context.Background()); err != nil {
				// Error from closing listeners, or context timeout.
				logger.Printf("Server graceful shutdown failed: %v", err)
			}
			close(done)
		}()

		select {
		case <-done:
			logger.Println("Server graceful shutdown completed")
			return
		case <-osSignals:
			// Another interrupt received, starting force shutdown
			logger.Println("Server force shutdown")
		case <-time.After(srv.timeout):
			logger.Println("Server graceful shutdown timed out, force shutdown")
		}

		if err := srv.s.Close(); err != nil {
			logger.Fatal("Server force shutdown failed: ", err)
		}
	}()
}
