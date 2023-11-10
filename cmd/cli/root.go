package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/AYM1607/ccclip/internal/configfile"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ccclip",
	Short: "copy strings to and from your end to end encrypted cloud clipboard",
	Long:  `copy strings to and from your end to end encrypted cloud clipboard`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
