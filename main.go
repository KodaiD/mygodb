package main

import (
	"bufio"
	"fmt"
	"os"
)

func printPrompt() {
	fmt.Print("db > ")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		printPrompt()
		if scanner.Scan() {
			inputBuffer := scanner.Text()
			if inputBuffer == ".exit" {
				os.Exit(0)
			} else {
				fmt.Printf("Unrecognized command '%s' .\n", inputBuffer)
			}
		}
	}
}