package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/AYM1607/ccclip/internal/configfile"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

func getClipboard() (string, error) {
	cc, err := configfile.EnsureAndGet()
	if err != nil {
		return "", err
	}

	if cc.DeviceId == "" {
		return "", errors.New("you must log in and register your device")
	}
	pvk, err := configfile.LoadPrivateKey()
	if err != nil {
		return "", fmt.Errorf("could not load this device's private key: %w", err)
	}

	plain, err := apiclient.GetClipboard(cc.DeviceId, pvk)
	if err != nil {
		return "", fmt.Errorf("could not set clipboard: %w", err)
	}

	return plain, nil
}

func setClipboard(clip []byte) error {
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

	err = apiclient.SetClipboard(clip, cc.DeviceId, pvk)
	if err != nil {
		return fmt.Errorf("could not set clipboard: %w", err)
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ccclip",
	Short: "copy strings to and from your end to end encrypted cloud clipboard",
	Long:  `copy strings to and from your end to end encrypted cloud clipboard`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd()) {
			// Nothing piped through stdin. Reading clipboard.
			clip, err := getClipboard()
			if err != nil {
				return err
			}
			_, err = os.Stdout.Write([]byte(clip))
			return err
		}

		clip, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		return setClipboard(clip)
	},
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("could not locate home directory: %s", err.Error()))
	}
	defualtConfigPath := path.Join(homeDir, ".config", "ccclip")

	rootCmd.PersistentFlags().StringVar(&configfile.Path, "config-dir", defualtConfigPath, "location of the ccclip.yaml configuration file")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
