package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

type User struct {
	Password string
}

var (
	users = make(map[string]User)
	mu    sync.RWMutex
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request) {
		phone := r.URL.Query().Get("phone")
		pass := r.URL.Query().Get("password")

		mu.Lock()
		if _, exists := users[phone]; exists {
			http.Error(w, "User exists", http.StatusConflict)
			mu.Unlock()
			return
		}
		users[phone] = User{Password: pass}
		mu.Unlock()
		fmt.Fprint(w, `{"status": "success"}`)
	})

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
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
