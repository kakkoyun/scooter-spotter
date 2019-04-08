package scooter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	defaultScooterSearchAPIUrl = "https://qc05n0gp78.execute-api.eu-central-1.amazonaws.com/prod/BackendGoChallenge"
	defaultTimeout             = 1 * time.Second

	url     string
	timeout time.Duration
	client  *http.Client
	logger  *log.Logger
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
	// TODO Better to use time.ParseDuration
	if rawTimeout := os.Getenv("SCOOTER_SEARCH_API_TIMEOUT"); rawTimeout == "" {
		timeout = defaultTimeout
	} else {
		val, err := strconv.ParseInt(rawTimeout, 10, 32)
		if err != nil {
			timeout = defaultTimeout
		}
		timeout = time.Duration(val) * time.Second
	}

	// TODO Better to use multiple clients for large max value
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
	found := make(chan Scooter, max)
	done := make(chan struct{})
	defer close(done)

	go scheduleWorkers(max, found, done)
	logger.Printf("Workers are scheduled")

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// TODO To get the wall time limit,
			// - create a channel exact time when request received and push it down here.
			logger.Println("Timed out")
			return scooters, nil
		case sc, ok := <-found:
			// logger.Println("Found:", sc)
			if !ok {
				logger.Println("Scheduler closed")
				return scooters, nil
			}
			if sc.AvailableForRent && sc.BatteryLevel > 20 {
				scooters = append(scooters, sc)
			}
			if len(scooters) == max {
				logger.Println("Reached maximum requested capacity")
				return scooters, nil
			}
			continue
		default:
			// no-op
		}
	}
}

// Helpers

func scheduleWorkers(max int, found chan<- Scooter, done <-chan struct{}) {
	// TODO Needs better algorithms find how many workers to schedule,
	// - explorer incremental approach depending results
	var wg sync.WaitGroup

	maxNumberWorkers := max * 2
	for i := 0; i < maxNumberWorkers; i++ {
		id := i + 1
		select {
		case <-done:
			// Received done, stop scheduling new workers
			logger.Println("Scheduler stopped by owner")
			break
		default:
			// Schedule a new worker
			wg.Add(1)
			go func(id int, wg *sync.WaitGroup) {
				// logger.Println("Scheduled:", id)
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
				wg.Done()
			}(id, &wg)
		}
	}

	wg.Wait()
	// Done with scheduling close channel
	close(found)
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
