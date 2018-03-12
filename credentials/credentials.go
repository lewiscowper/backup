package credentials

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

// Get takes in credentials from the terminal
func Get() ([]byte, error) {
	fmt.Print("Enter password for archive encryption: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}

	fmt.Print("\n")

	return bytePassword, nil
}
