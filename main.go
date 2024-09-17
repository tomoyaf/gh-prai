package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	var createHelp bool
	createCmd.BoolVar(&createHelp, "help", false, "Show help for create command")
	createCmd.BoolVar(&createHelp, "h", false, "Show help for create command")
	createBase := createCmd.String("base", "", "Specify the base branch for the PR")

	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	var configHelp bool
	configCmd.BoolVar(&configHelp, "help", false, "Show help for config command")
	configCmd.BoolVar(&configHelp, "h", false, "Show help for config command")
	
	if len(os.Args) == 1 {
		createPR()
		os.Exit(0)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printMainHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if createHelp {
			printCreateHelp()
			os.Exit(0)
		}
		baseBranch = *createBase
		createPR()
	case "config":
		configCmd.Parse(os.Args[2:])
		if configHelp {
			printConfigHelp()
			os.Exit(0)
		}
		switch configCmd.Arg(0) {
		case "show":
			if configCmd.Arg(1) == "-h" || configCmd.Arg(1) == "--help" {
				printConfigShowHelp()
				os.Exit(0)
			}
			if configCmd.NArg() > 1 {
				fmt.Println("Error: Too many arguments for config show command")
				printConfigShowHelp()
				os.Exit(1)
			}
			showConfig()
		case "reset":
			if configCmd.Arg(1) == "-h" || configCmd.Arg(1) == "--help" {
				printConfigResetHelp()
				os.Exit(0)
			}
			if configCmd.NArg() > 1 {
				fmt.Println("Error: Too many arguments for config reset command")
				printConfigResetHelp()
				os.Exit(1)
			}
			resetConfig()
		default:
			if configCmd.NArg() < 2 {
				fmt.Println("Error: Insufficient arguments for config command")
				printConfigHelp()
				os.Exit(1)
			}
			configureSettings(configCmd.Arg(0), configCmd.Arg(1))
		}
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printMainHelp()
		os.Exit(1)
	}
}

func printMainHelp() {
	fmt.Println("Usage: gh prai [command] [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  create    Create or update a Pull Request with AI-generated title and description")
	fmt.Println("  config    Configure settings for the gh-prai extension")
	fmt.Println("\nOptions:")
	fmt.Println("  -h, --help    Show this help message")
	fmt.Println("\nRun 'gh prai <command> --help' for more information on a command.")
	fmt.Println("\nIf no command is specified, 'gh prai' will default to the 'create' command.")
}

func printCreateHelp() {
	fmt.Println("Usage: gh prai create [options]")
	fmt.Println("\nCreate or update a Pull Request with AI-generated title and description")
	fmt.Println("\nOptions:")
	fmt.Println("  --base string   Specify the base branch for the PR")
	fmt.Println("  --help, -h      Show this help message")
	fmt.Println("\nIf no options are specified, the command will use default settings.")
}

func printConfigHelp() {
	fmt.Println("Usage: gh prai config <key> <value>")
	fmt.Println("\nConfigure settings for the gh-prai extension")
	fmt.Println("\nCommands:")
	fmt.Println("  show     Show the current configuration settings")
	fmt.Println("  reset    Reset the configuration settings to default values")
	fmt.Println("\nAvailable keys:")
	fmt.Println("  api_key    Set the OpenAI API key")
	fmt.Println("  language   Set the language for PR title and description (e.g., 'en' for English, 'ja' for Japanese)")
	fmt.Println("  template   Set the template for PR description (e.g., custom template path like: './.github/pull_request_template.md', 'basic' for the basic template)")
	fmt.Println("  prompt     Set the custom prompt for AI generation")
	fmt.Println("\nOptions:")
	fmt.Println("  --help, -h     Show this help message")
}

func printConfigShowHelp() {
	fmt.Println("Usage: gh prai config show")
	fmt.Println("\nShow the current configuration settings")
}

func printConfigResetHelp() {
	fmt.Println("Usage: gh prai config reset")
	fmt.Println("\nReset the configuration settings to default values")
}
