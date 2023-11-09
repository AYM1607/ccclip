package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getDevicesCmd)
}

var getDevicesCmd = &cobra.Command{
	Use:   "get-devices",
	Short: "Register a user with a given email and password",
	RunE: func(cmd *cobra.Command, args []string) error {
		// cc, err := configfile.EnsureAndGet()
		// if err != nil {
		// 	return err
		// }
		// if cc.DeviceId == "" {
		// 	return errors.New("your device is not registered")
		// }
		// pvk, err := configfile.LoadPrivateKey()
		// if err != nil {
		// 	return err
		// }
		// devices, err := apiclient.GetDevices(cc.DeviceId, pvk)
		// if err != nil {
		// 	return err
		// }

		// return json.NewEncoder(os.Stdout).Encode(devices)
		return nil
	},
}
