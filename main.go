package main

import (
	"bufio"
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
	UserName	[]byte
	Email		[]byte
}

type Page struct {
	rows [][]byte
}

type Table struct {
	NumRows uint32
	Pages []*Page
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	table := newTable()
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
			switch executeStatement(statement, table) {
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

func executeInsert(statement *Statement, table *Table) ExecuteResult {
	if table.NumRows >= TABLE_MAX_ROWS {
		return EXECUTE_TABLE_FULL
	}
	row := &statement.RowToInsert
	serializeRow(table, row)
	table.NumRows += 1
	return EXECUTE_SUCCESS
}

func executeSelect(statement *Statement, table *Table) ExecuteResult {
	row := &statement.RowToInsert
	for i := uint32(0); i < table.NumRows; i++ {
		deserializeRow(table, row)
		row.printRow()
	}
	return EXECUTE_SUCCESS
}

func executeStatement(statement *Statement, table *Table) ExecuteResult {
	switch statement.Type {
	case STATEMENT_INSERT:
		return executeInsert(statement, table)
	case STATEMENT_SELECT:
		return executeSelect(statement, table)
	default:
		fmt.Println("Warning...")
		return EXECUTE_SUCCESS
	}
}

func serializeRow(table *Table, r *Row) {
	pn, rn := table.rowSlot()
	buf := make([][]byte, ROWS_PER_PAGE)
	table.Pages[pn] = &Page{rows: buf}
	page := table.Pages[pn]
	bID := make([]byte, 4)
	binary.BigEndian.PutUint32(bID, r.ID)
	copy(page.rows[rn+1][:], bID)
	fmt.Println(page.rows[rn+1])
	copy(page.rows[rn+1][4:], r.UserName)
	copy(page.rows[rn+1][4+USERNAME_SIZE:], r.Email)
}

func deserializeRow(table *Table, r *Row) {
	pn, _ := table.rowSlot()
	for i := uint32(0); i <= pn; i++ {
		for j := uint32(0); j < ROWS_PER_PAGE; j++ {
			fmt.Println(table.Pages[i])
			bID := binary.BigEndian.Uint32(table.Pages[i].rows[j][:4])
			fmt.Println(bID, table.Pages[i].rows[j][4:USERNAME_SIZE], table.Pages[i].rows[j][USERNAME_SIZE:])
			r.ID = bID
			copy(r.UserName, table.Pages[i].rows[j][4:USERNAME_SIZE])
			copy(r.Email, table.Pages[i].rows[j][USERNAME_SIZE:])
		}
	}
}

func newTable() *Table {
	table := &Table{
		NumRows: 0,
		Pages: make([]*Page, TABLE_MAX_PAGES),
	}
	return table
}

func (r *Row) printRow() {
	fmt.Printf("(%d, %s, %s)\n", r.ID, r.UserName, r.Email)
}

func (table *Table) rowSlot() (uint32, uint32) {
	rowNum := table.NumRows
	pageNum := rowNum / ROWS_PER_PAGE
	rowOffset := rowNum % ROWS_PER_PAGE
	return pageNum, rowOffset
}