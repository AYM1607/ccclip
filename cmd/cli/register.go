package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/AYM1607/ccclip/internal/configfile"
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
		err := apiclient.Register(email, password)
		if err != nil {
			return fmt.Errorf("could not register user: %w", err)
		}

		cc, err := configfile.EnsureAndGet()
		if err != nil {
			return err
		}
		cc.Email = email
		return configfile.Write(cc)
	},
}
