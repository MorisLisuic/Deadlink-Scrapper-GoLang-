package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Welcome to the GoLang Deadlink-Scrapper v1.")
	fmt.Println("For help run 'help'")
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter commands (type 'exit' to quit):")

	for {
		fmt.Print("> ") // Prompt for input
		scanner.Scan()  // Read user input
		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			fmt.Println("Exiting program...")
			break // Stop the loop
		}
		handleCommands(input)
	}

	fmt.Println("")
}
