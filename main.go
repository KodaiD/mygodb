package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	META_COMMAND_SUCCESS = 0
	META_COMMAND_UNRECOGNIZED_COMMAND = 1
	META_COMMAND_EXIT = 2
)

const (
	PREPARE_SUCCESS = 0
	PREPARE_UNRECOGNIZED_STATEMENT = 1
)

const (
	STATEMENT_INSERT = 0
	STATEMENT_SELECT = 1
)

type Statement struct {
	Type int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		printPrompt()
		if scanner.Scan() {
			inputBuffer := scanner.Text()
			if inputBuffer[0] == '.' {
				switch doMetaCommand(inputBuffer) {
				case META_COMMAND_EXIT:
					os.Exit(0)
				case META_COMMAND_SUCCESS:
					continue
				case META_COMMAND_UNRECOGNIZED_COMMAND:
					fmt.Printf("Unrecognized command '%s' .\n", inputBuffer)
					continue
				}
			}
			statement := &Statement{}
			switch prepareStatement(inputBuffer, statement) {
			case PREPARE_SUCCESS:
				break
			case PREPARE_UNRECOGNIZED_STATEMENT:
				fmt.Printf("Unrecognized keyword at start of '%s'.\n", inputBuffer)
				continue
			}
			executeStatement(statement)
			fmt.Println("Executed.")
		}
	}
}

func printPrompt() {
	fmt.Print("db > ")
}

func doMetaCommand(inputBuffer string) int {
	if inputBuffer == ".exit" {
		return META_COMMAND_EXIT
	} else {
		return META_COMMAND_UNRECOGNIZED_COMMAND
	}
}

func prepareStatement(inputBuffer string, statement *Statement) int {
	if inputBuffer[:6] == "insert" {
		statement.Type = STATEMENT_INSERT
		return PREPARE_SUCCESS
	}
	if inputBuffer == "select" {
		statement.Type = STATEMENT_SELECT
		return PREPARE_SUCCESS
	}
	return PREPARE_UNRECOGNIZED_STATEMENT
}

func executeStatement(statement *Statement) {
	switch statement.Type {
	case STATEMENT_INSERT:
		fmt.Println("This is where we would do an insert.")
		break
	case STATEMENT_SELECT:
		fmt.Println("This is where we would do a select.")
		break
	}
}