package main

import (
	"encoding/base64"

	"github.com/AYM1607/ccclip/internal/server/client"
)

func b64(i []byte) string {
	return base64.StdEncoding.EncodeToString(i)
}

var apiclient *client.Client

func init() {
	apiclient = client.New("http://localhost:8080")
}

func main() {
	rootCmd.Execute()
}
