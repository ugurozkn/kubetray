package main

import (
	"os"

	"github.com/ugurozkn/kubetray/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
