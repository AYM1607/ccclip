package main

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/AYM1607/ccclip/internal/configfile"
	"github.com/AYM1607/ccclip/pkg/crypto"
	"github.com/AYM1607/ccclip/pkg/input"
)

func init() {
	rootCmd.AddCommand(registerDeviceCommand)
}

var registerDeviceCommand = &cobra.Command{
	Use:   "register-device",
	Short: "Register a device for the given user",
	RunE: func(cmd *cobra.Command, args []string) error {
		cc, err := configfile.EnsureAndGet()
		if err != nil {
			return err
		}

		if cc.Email == "" {
			return errors.New("you don't have an account configured for thist device")
		}

		if cc.DeviceId != "" {
			return errors.New("this device is already registered")
		}

		pvk := crypto.NewPrivateKey()
		pbk := pvk.PublicKey()

		password := input.ReadPassword()
		res, err := apiclient.RegisterDevice(cc.Email, password, pbk.Bytes())
		if err != nil {
			return err
		}

		// Write the key files first, if those fail to write then we should not
		// save the device Id.
		cc.DeviceId = res.DeviceID
		err = configfile.SavePrivateKey(pvk)
		if err != nil {
			return err
		}
		err = configfile.SavePublicKey(pbk)
		if err != nil {
			return err
		}
		return configfile.Write(cc)
	},
}
