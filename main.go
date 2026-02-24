package main

import (
	"fmt"
	"os"

	"github.com/JonecoBoy/s4g/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
