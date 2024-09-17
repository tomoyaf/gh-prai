package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		createPR()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "create":
		createPR()
	case "config":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gh prai config [key] [value]")
			fmt.Println("Available keys: api_key, language, template, prompt")
			os.Exit(1)
		}
		configureSettings(os.Args[2], os.Args[3])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
