package main

import (
	"github.com/AYM1607/ccclip/internal/server/client"
)

var apiclient *client.Client

func init() {
	apiclient = client.New("https://api.ccclip.io")
}

func main() {
	rootCmd.Execute()
}
