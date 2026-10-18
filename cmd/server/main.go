package main

import (
	"log"
	"net/http"
	"os"

	"bandwidth-hero-proxy/internal/proxy"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/index", proxy.Handler)
	http.HandleFunc("/", proxy.Handler)

	log.Printf("Starting Bandwidth Hero Proxy on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
