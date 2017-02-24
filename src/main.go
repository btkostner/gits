// main.go
// Handles CLI entry to gits

package main

import (
	"fmt"
	"os"

	"github.com/btkostner/gits/src/cmd"
)

func main() {
	if err := cmd.MainCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
}
