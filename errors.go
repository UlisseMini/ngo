// errors.go contains all errors for the program
package main

import (
	"errors"
)

var (
	NoHostErr error = errors.New(
		"You must specify a host to connect to. QUITTING.")
)
