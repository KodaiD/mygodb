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

const (
	COLUMN_USERNAME_SIZE uint32 = 32
	COLUMN_EMAIL_SIZE    uint32 = 255
)

const (
	ID_SIZE uint32 		= 4
	USERNAME_SIZE 		= COLUMN_USERNAME_SIZE
	EMAIL_SIZE 			= COLUMN_EMAIL_SIZE
	ID_OFFSET uint32 	= 0
	USERNAME_OFFSET 	= ID_OFFSET + ID_SIZE
	EMAIL_OFFSET	 	= USERNAME_OFFSET + USERNAME_SIZE
	ROW_SIZE 			= ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
)

const (
	PAGE_SIZE uint32 = 4096 // 4KB
	TABLE_MAX_PAGES uint32 = 100
	ROWS_PER_PAGE = PAGE_SIZE / ROW_SIZE
	TABLE_MAX_ROWS = ROWS_PER_PAGE * TABLE_MAX_PAGES
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