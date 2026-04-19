package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Serve static files (HTML/CSS/JS) from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// A simple API endpoint
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message": "Hello from the Go backend!"}`)
	})

	fmt.Println("Server starting at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
