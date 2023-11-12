package main

import (
	"log"
	"os"

	"github.com/AYM1607/ccclip/internal/server"
)

func main() {
	port := os.Getenv("CCCLIP_PORT")

	if port == "" {
		port = "8080"
	}
	log.Printf("Serving on port %s", port)
	s := server.New(":" + port)
	log.Fatal(s.ListenAndServe())
}
