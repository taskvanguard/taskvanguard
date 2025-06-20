package main

import (
	"fmt"
	"os"

	"xarc.dev/taskvanguard/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}