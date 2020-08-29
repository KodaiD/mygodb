package main

import "fmt"

const FILE_NAME = "data.txt"

const (
	META_COMMAND_SUCCESS 				MetaCommandResult = 0
	META_COMMAND_UNRECOGNIZED_COMMAND 	MetaCommandResult = 1
	META_COMMAND_EXIT 					MetaCommandResult = 2
)

const (
	PREPARE_SUCCESS 				PrepareResult= 0
	PREPARE_UNRECOGNIZED_STATEMENT 	PrepareResult = 1
	PREPARE_SYNTAXERROR				PrepareResult = 2
)

const (
	STATEMENT_INSERT StatementType = 0
	STATEMENT_SELECT StatementType = 1
)

const (
	EXECUTE_SUCCESS 	ExecuteResult = 0
	EXECUTE_TABLE_FULL 	ExecuteResult = 1
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