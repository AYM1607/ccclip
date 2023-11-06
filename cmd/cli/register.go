package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/AYM1607/ccclip/internal/server"
	"github.com/spf13/cobra"
)

var email string
var password string

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(&email, "email", "e", "", "email will be your login identifier")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "password will secure your account")

	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("password")
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a user with a given email and password",
	RunE: func(cmd *cobra.Command, args []string) error {
		req := server.RegisterRequest{
			Email:    email,
			Password: password,
		}
		reqJson, err := json.Marshal(req)
		if err != nil {
			return err
		}
		res, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewReader(reqJson))
		if err != nil {
			return err
		}
		resBody, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return err
		}

		log.Println(string(resBody))
		return nil
	},
}
