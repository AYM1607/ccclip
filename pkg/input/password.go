package input

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// ReadPassword reads a single line of text from the terminal withouth echoing it out.
func ReadPassword() string {
	raw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("could not reat password from the terminal: %s", err.Error()))
	}
	return string(raw)
}
