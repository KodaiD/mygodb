package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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

func main() {
	fmt.Println("Running db...")
	db := newDB(os.Stdin)
	db.run(os.Stdout)
}