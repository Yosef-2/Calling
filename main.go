package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Signal struct {
	From string `json:"from"`
	Data string `json:"data"` // This will hold the WebRTC handshake
}

var (
	users   = make(map[string]string) // phone:password
	signals = make(map[string]chan Signal) // phone:channel for signals
	mu      sync.Mutex
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request) {
		phone := r.URL.Query().Get("phone")
		mu.Lock()
		users[phone] = "verified"
		signals[phone] = make(chan Signal, 10)
		mu.Unlock()
		fmt.Fprint(w, `{"status": "success"}`)
	})

	// Send a signal (Offer/Answer) to a specific phone number
	http.HandleFunc("/api/send-signal", func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("to")
		var s Signal
		json.NewDecoder(r.Body).Decode(&s)
		
		mu.Lock()
		if ch, ok := signals[target]; ok {
			ch <- s
		}
		mu.Unlock()
	})

	// Listen for incoming signals (long polling)
	http.HandleFunc("/api/get-signal", func(w http.ResponseWriter, r *http.Request) {
		phone := r.URL.Query().Get("phone")
		mu.Lock()
		ch, ok := signals[phone]
		mu.Unlock()

		if ok {
			select {
			case sig := <-ch:
				json.NewEncoder(w).Encode(sig)
			default:
				w.WriteHeader(http.StatusNoContent)
			}
		}
	})

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.ListenAndServe(":"+port, nil)
}
