package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Scooter stores meta data for a Scooter
type Scooter struct {
	IDd              int  `json:"id"`
	BatteryLevel     int  `json:"battery_level"`
	AvailableForRent bool `json:"available_for_rent"`
}

// Response represents payload of service response
type Response struct {
	Data    []Scooter `json:"data,omitempty"`
	Message string    `json:"message,omitempty"`
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

		max, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Response{Message: fmt.Sprintf("Error: %s is not a positive valid integer", param)})
			return
		}

		// 2. Return external service results
		encoder.Encode(Response{Message: fmt.Sprintf("Given param: %d", max)})
	}))

	router.Handle("/_status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
		w.Write([]byte("\n"))
	}))

	return router
}
