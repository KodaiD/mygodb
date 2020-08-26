package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	ID		 uint32
	UserName []byte
	Email	 []byte
}

type Page struct {
	rows [][]byte
}

type Pager struct {
	fileDescriptor io.Reader
	Pages          []*Page
}

type Table struct {
	NumRows uint32
	Pager   *Pager
}

func main() {
	fmt.Println("Running db...")
	db := newDB(os.Stdin)
	db.run(os.Stdout)
}

type DB struct {
	scanner *bufio.Scanner
	table *Table
}

func (db *DB) run(o io.Writer) {
	for {
		printPrompt()
		if db.scanner.Scan() {
			inputBuffer := db.scanner.Text()
			if inputBuffer[0] == '.' {
				switch doMetaCommand(inputBuffer) {
				case META_COMMAND_EXIT:
					os.Exit(0)
				case META_COMMAND_SUCCESS:
					continue
				case META_COMMAND_UNRECOGNIZED_COMMAND:
					fmt.Fprintf(o, "Unrecognized command '%s' .\n", inputBuffer)
					continue
				}
			}
			statement := &Statement{}
			switch prepareStatement(inputBuffer, statement) {
			case PREPARE_SUCCESS:
				break
			case PREPARE_SYNTAXERROR:
				fmt.Fprintln(o, "Syntax error. Could not parse statement.")
				continue
			case PREPARE_UNRECOGNIZED_STATEMENT:
				fmt.Fprintf(o, "Unrecognized keyword at start of '%s'.\n", inputBuffer)
				continue
			}
			switch executeStatement(statement, db.table) {
			case EXECUTE_SUCCESS:
				fmt.Fprintln(o, "Executed.")
				break
			case EXECUTE_TABLE_FULL:
				fmt.Fprintln(o, "Error: Table full.")
				break
			}
		}
	}
}

func newDB(output io.Reader) *DB {
	db := &DB{
		scanner: bufio.NewScanner(output),
		table: newTable(),
	}
	return db
}

func newPager() *Pager {
	file, err := os.Open(FILE_NAME)
	if err != nil {
		log.Println(err)
	}
	p := &Pager{
		fileDescriptor: file,
		Pages:          make([]*Page, TABLE_MAX_PAGES),
	}
	for i := uint32(0); i < TABLE_MAX_PAGES; i++ {
		buf := make([][]byte, ROWS_PER_PAGE)
		p.Pages[i] =  &Page{rows: buf}
	}
	return p
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
		deserializeRow(table, row, i)
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
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println("---", err)
	}
	pn, rn := table.rowSlot(table.NumRows)
	table.Pager.Pages[pn].rows[rn] = data
}

func deserializeRow(table *Table, r *Row, i uint32) {
	pn, rn := table.rowSlot(i)
	data := table.Pager.Pages[pn].rows[rn]
	err := json.Unmarshal(data, r)
	if err != nil {
		fmt.Println("+++", err)
	}
}

func newTable() *Table {
	table := &Table{
		NumRows: 0,
		Pager: newPager(),
	}
	return table
}

func (r *Row) printRow() {
	fmt.Printf("(%d, %s, %s)\n", r.ID, r.UserName, r.Email)
}

func (table *Table) rowSlot(rowNum uint32) (uint32, uint32) {
	pageNum := rowNum / ROWS_PER_PAGE
	rowOffset := rowNum % ROWS_PER_PAGE
	return pageNum, rowOffset
}