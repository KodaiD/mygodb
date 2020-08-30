package main

import (
	"bufio"
	"fmt"
	"io"
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
					db.table.closeTable()
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

func main() {
	fmt.Println("Running db...")
	db := newDB(os.Stdin)
	db.run(os.Stdout)
}