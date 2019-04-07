package scooter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	defaultScooterSearchAPIUrl = "https://qc05n0gp78.execute-api.eu-central-1.amazonaws.com/prod/BackendGoChallenge"

	url    string
	client *http.Client
	logger *log.Logger
)

// Scooter stores meta data for a Scooter
type Scooter struct {
	ID               int  `json:"id"`
	BatteryLevel     int  `json:"battery_level"`
	AvailableForRent bool `json:"available_for_rent"`
}

func init() {
	// Initialize package dependencies
	logger = log.New(os.Stdout, "api: ", log.LstdFlags)

	if url = os.Getenv("SCOOTER_SEARCH_API_URL"); url == "" {
		url = defaultScooterSearchAPIUrl
	}

	// Customize the Transport to have larger connection pool
	defaultTransportPointer, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("DefaultTransport not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer
	defaultTransport.MaxIdleConns = 100        // TODO parameterize
	defaultTransport.MaxIdleConnsPerHost = 100 // TODO parameterize

	client = &http.Client{Transport: &defaultTransport}
}

func FindAll(max int) ([]Scooter, error) {
	logger.Printf("Find all available Scooters: %d\n", max)
	scooters := []Scooter{}
	found := make(chan Scooter)
	defer close(found)

	// Needs better algorithms find how many workers to schedule
	for i := 0; i < max*2; i++ {
		id := i + 1
		go func(id int) {
			defer func() {
				if r := recover(); r != nil {
					logger.Printf("Failed to fetch scooter with id %d, PANIC %v\n", id, r)
				}
			}()

			sc, err := fetchScooter(id)
			if err != nil {
				logger.Printf("Failed to fetch scooter with id %d, %v\n", id, err)
				return
			}

			found <- sc
		}(id)
	}

	logger.Printf("Workers are scheduled")

	for {
		select {
		case sc := <-found:
			// logger.Println("Found:", sc)
			if sc.AvailableForRent && sc.BatteryLevel > 20 {
				scooters = append(scooters, sc)
			}
			if len(scooters) >= max {
				logger.Panicln("Reached maximum allowed")
				return scooters, nil
			}
			continue
		case <-time.After(1 * time.Second): // TODO parameterize
			logger.Println("Timed out")
			return scooters, nil
		}
	}
}

func fetchScooter(id int) (Scooter, error) {
	var scooter Scooter

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return scooter, err
	}

	q := req.URL.Query()
	q.Add("id", strconv.Itoa(id))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return scooter, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&scooter); err != nil {
		return scooter, err
	}
	return scooter, nil
}
