package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Serve index.html from the root folder
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// API Endpoint - Removed the "Render" specific text
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message": "Backend connection successful!"}`)
	})

	// Railway usually uses PORT, just like Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server live on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
