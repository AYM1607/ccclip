package main

import (
	"log"

	"github.com/spf13/cobra"
)

var keyset int

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
	rootCmd.PersistentFlags().IntVarP(&keyset, "keyset", "k", 0, "which key set to use, can be 1 or 2")

	rootCmd.MarkPersistentFlagRequired("keyset")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
