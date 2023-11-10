package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/AYM1607/ccclip/internal/configfile"
	"github.com/AYM1607/ccclip/pkg/input"
)

var email string

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(&email, "email", "e", "", "email will be your login identifier")

	registerCmd.MarkFlagRequired("email")
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a user with a given email and password",
	RunE: func(cmd *cobra.Command, args []string) error {
		password := input.ReadPassword()
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
