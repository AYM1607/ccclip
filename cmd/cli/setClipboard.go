package main

import (
	"errors"
	"fmt"

	"github.com/AYM1607/ccclip/internal/configfile"
	"github.com/spf13/cobra"
)

var clipboard string

func init() {
	rootCmd.AddCommand(setClipboardCmd)

	setClipboardCmd.Flags().StringVar(&clipboard, "clip", "", "the string to send")
	setClipboardCmd.MarkFlagRequired("clip")
}

var setClipboardCmd = &cobra.Command{
	Use:   "set-clipboard",
	Short: "set the given string as the cloud clipboard",
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

		err = apiclient.SetClipboard(clipboard, cc.DeviceId, pvk)
		if err != nil {
			return fmt.Errorf("could not set clipboard: %w", err)
		}
		return nil
	},
}
