package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kakkoyun/scooter-spotter/scooter"
)

// Response represents payload of service response
type Response struct {
	Data    []scooter.Scooter `json:"data,omitempty"`
	Message string            `json:"message,omitempty"`
}

func newRouter() http.Handler {
	router := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
		w.Write([]byte("\n"))
	}))

	router.Handle("/scooters", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)

		// 1. Parse parameters
		param := r.URL.Query().Get("max")
		if param == "" {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{Message: `Error: Required param "max" is missing`})
			return
		}

		max, err := strconv.ParseInt(param, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{Message: fmt.Sprintf("Error: %s is not a valid integer", param)})
			return
		}

		if max <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{Message: fmt.Sprintf("Error: %s is not a positive integer", param)})
			return
		}

		// 2. Return external service results
		result, err := scooter.FindAll(int(max))
		if err != nil {
			logger.Printf("Scooter service failed: %v", err)
			encoder.Encode(Response{Message: fmt.Sprintf("Scooter service failed")})
		}
		encoder.Encode(Response{Data: result})
	}))

	router.Handle("/_status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
		w.Write([]byte("\n"))
	}))

	return router
}
