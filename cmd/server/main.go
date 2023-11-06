package main

import (
	"log"
	"os"

	"github.com/AYM1607/ccclip/internal/config"
	"github.com/AYM1607/ccclip/internal/server"
)

func main() {
	privateKeyPath := os.Getenv("CCCLIP_PRIVATE_KEY")
	publicKeyPath := os.Getenv("CCCLIP_PUBLIC_KEY")
	if publicKeyPath == "" || privateKeyPath == "" {
		log.Fatalf("both public and privae keys must be provided")
	}

	config.Default.PrivateKeyPath = privateKeyPath
	config.Default.PublicKeyPath = publicKeyPath

	port := "8080"
	log.Printf("Serving on port %s", port)
	s := server.New(":" + port)
	log.Fatal(s.ListenAndServe())
}
