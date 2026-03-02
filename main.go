package main

import (
	"fmt"
	"os"

	"github.com/LeartS/mb/cmd/root"
)

func main() {
	if err := root.NewCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
