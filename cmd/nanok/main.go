package main

import (
	"fmt"
	"os"

	"github.com/knbr13/nano-k/pkg/container"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: nanok <command> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		container.Run()
	case "child":
		container.Child()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
