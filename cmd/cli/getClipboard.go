package main

import (
	"errors"
	"fmt"

	"github.com/AYM1607/ccclip/internal/configfile"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getClipboardCmd)
}

var getClipboardCmd = &cobra.Command{
	Use:   "get-clipboard",
	Short: "get the currently stored clipboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		cc, err := configfile.EnsureAndGet()
		if err != nil {
			return err
		}

		if cc.DeviceId == "" {
			return errors.New("you must log in and register your device")
		}
		pvk, err := configfile.LoadPrivateKey()
		if err != nil {
			return fmt.Errorf("could not load this device's private key: %w", err)
		}

		plain, err := apiclient.GetClipboard(cc.DeviceId, pvk)
		if err != nil {
			return fmt.Errorf("could not set clipboard: %w", err)
		}

		fmt.Printf("Your current clipbard is %q\n", plain)
		return nil
	},
}
