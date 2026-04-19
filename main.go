package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type User struct {
	Password string `json:"password"`
}

var (
	users = make(map[string]User) // Key: phone number
	mu    sync.RWMutex
)

func main() {
	// Serve index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Sign Up Logic
	http.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request) {
		phone := r.URL.Query().Get("phone")
		pass := r.URL.Query().Get("password")

		mu.Lock()
		if _, exists := users[phone]; exists {
			http.Error(w, "User already exists", http.StatusConflict)
			mu.Unlock()
			return
		}
		users[phone] = User{Password: pass}
		mu.Unlock()
		fmt.Fprint(w, `{"status": "success"}`)
	})

	// Call Logic (Simple verification if the person exists)
	http.HandleFunc("/api/call", func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("target")
		
		mu.RLock()
		_, exists := users[target]
		mu.RUnlock()

		if !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, `{"status": "calling", "target": "%s"}`, target)
	})

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.ListenAndServe(":"+port, nil)
}
