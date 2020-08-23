package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

/*
limitations
- support two operations: inserting a row and printing all rows
- reside only in memory (no persistence to disk)
- support a single, hard-coded table

data structure
- Store rows in blocks of memory called pages
- Each page stores as many rows as it can fit
- Rows are serialized into a compact representation with each page
- Pages are only allocated as needed
- Keep a fixed-size array of pointers to pages
*/

type MetaCommandResult int
type PrepareResult int
type StatementType int
type ExecuteResult int

type Statement struct {
	Type 		StatementType
	RowToInsert Row
}

type Row struct {
	ID			uint32
	UserName	[COLUMN_USERNAME_SIZE]byte
	Email		[COLUMN_EMAIL_SIZE]byte
}

type Page struct {
	rows [ROWS_PER_PAGE]*Row
}

type Table struct {
	NumRows uint32
	Pages [TABLE_MAX_PAGES]*Page
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	table := newTable()
	var buf bytes.Buffer
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
			case PREPARE_SYNTAXERROR:
				fmt.Println("Syntax error. Could not parse statement.")
				continue
			case PREPARE_UNRECOGNIZED_STATEMENT:
				fmt.Printf("Unrecognized keyword at start of '%s'.\n", inputBuffer)
				continue
			}
			switch executeStatement(statement, table, &buf) {
			case EXECUTE_SUCCESS:
				fmt.Println("Executed.")
				break
			case EXECUTE_TABLE_FULL:
				fmt.Println("Error: Table full.")
				break
			}
		}
	}
}

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

func prepareStatement(inputBuffer string, statement *Statement) PrepareResult {
	if inputBuffer[:6] == "insert" {
		statement.Type = STATEMENT_INSERT
		argsAssigned, err := fmt.Sscanf(inputBuffer, "insert %d %s %s", &statement.RowToInsert.ID,
									&statement.RowToInsert.UserName, &statement.RowToInsert.Email)
		if argsAssigned < 3 {
			log.Println(err)
			return PREPARE_SYNTAXERROR
		}
		return PREPARE_SUCCESS
	}
	if inputBuffer == "select" {
		statement.Type = STATEMENT_SELECT
		return PREPARE_SUCCESS
	}
	return PREPARE_UNRECOGNIZED_STATEMENT
}

func executeInsert(statement *Statement, table *Table, buf *bytes.Buffer) ExecuteResult {
	if table.NumRows >= TABLE_MAX_ROWS {
		return EXECUTE_TABLE_FULL
	}
	row := &statement.RowToInsert
	err := row.serializeRow(buf)
	if err != nil {
		log.Println(err)
	}
	table.NumRows += 1
	return EXECUTE_SUCCESS
}

func executeSelect(statement *Statement, table *Table, buf *bytes.Buffer) ExecuteResult {
	row := &statement.RowToInsert
	for i := uint32(0); i < table.NumRows; i++ {
		err := row.deserializeRow(buf)
		if err != nil {
			log.Println(err)
		}
		row.printRow()
	}
	return EXECUTE_SUCCESS
}

func executeStatement(statement *Statement, table *Table, buf *bytes.Buffer) ExecuteResult {
	switch statement.Type {
	case STATEMENT_INSERT:
		return executeInsert(statement, table, buf)
	case STATEMENT_SELECT:
		fmt.Println("This is where we would do a select.")
		return executeSelect(statement, table, buf)
	}
}

func (r *Row) serializeRow(buf *bytes.Buffer) error {
	err := binary.Write(buf, binary.BigEndian, r)
	return err
}

func (r *Row) deserializeRow(buf *bytes.Buffer) error {
	err := binary.Read(buf, binary.BigEndian, r)
	return err
}

func newTable() *Table {
	table := &Table{
		NumRows: 0,
	}
	return table
}

func (r *Row) printRow() {
	fmt.Printf("(%d, %s, %s)\n", r.ID, r.UserName, r.Email)
}

func (table *Table) rowSlot(rowNum uint32) uint32 {
	pageNum := rowNum / ROWS_PER_PAGE
	page := table.Pages[pageNum]
	if page == nil {
		page = table.Pages[pageNum]
	}
	return pageNum
}