package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// HCP Backend Server
// Responsible for data persistence (PostgreSQL/Redis) and historical data analysis
// See Task.md Task 11 for details.

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Starting HCP Backend Server on port %s...\n", port)
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
