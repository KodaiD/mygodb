package main

import "fmt"

const FILE_NAME = "data.txt"

const (
	META_COMMAND_SUCCESS 				MetaCommandResult = 1 + iota
	META_COMMAND_UNRECOGNIZED_COMMAND
	META_COMMAND_EXIT
)

const (
	PREPARE_SUCCESS 				PrepareResult= 1 + iota
	PREPARE_UNRECOGNIZED_STATEMENT
	PREPARE_SYNTAXERROR
)

const (
	STATEMENT_INSERT StatementType = 1 + iota
	STATEMENT_SELECT
)

const (
	EXECUTE_SUCCESS 	ExecuteResult = 1 + iota
	EXECUTE_TABLE_FULL
)

func printPrompt() {
	fmt.Print("db > ")
}

func doMetaCommand(inputBuffer string) MetaCommandResult {
	if inputBuffer == ".exit" {
		return META_COMMAND_EXIT
	} else {
		return META_COMMAND_UNRECOGNIZED_COMMAND
	}
}