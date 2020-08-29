package main

import (
	"fmt"
	"os"
)

type InputBuffer struct {
	buffer       string
	bufferLength uint32
	inputLength  uint32
}

const (
	EXECUTE_SUCCESS = 0
	EXECUTE_TABLE_FULL
)
type ExecuteResult int

const (
	META_COMMAND_SUCCESS = 0
	META_COMMAND_UNRECOGNIZED_COMMAND
)
type MetaCommandResult int

const (
	PREPARE_SUCCESS = 0
	PREPARE_NEGATIVE_ID
	PREPARE_STRING_TOO_LONG
	PREPARE_SYNTAX_ERROR
	PREPARE_UNRECOGNIZED_STATEMENT
)
type PrepareResult int

const (
	STATEMENT_INSERT = 0
	STATEMENT_SELECT
)
type StatementType int

const (
	COLUMN_USERNAME_SIZE uint32 = 32
	COLUMN_EMAIL_SIZE    uint32 = 255
)

type Row struct {
	id       uint32
	username [COLUMN_USERNAME_SIZE]byte
	email    [COLUMN_EMAIL_SIZE]byte
}

type Statement struct {
	statementType StatementType
	rowToInsert   Row
}

const (
	ID_SIZE         = 32
	USERNAME_SIZE   = COLUMN_USERNAME_SIZE
	EMAIL_SIZE      = COLUMN_EMAIL_SIZE
	ID_OFFSET       = 0
	USERNAME_OFFSET = ID_OFFSET + ID_SIZE
	EMAIL_OFFSET    = USERNAME_OFFSET + USERNAME_SIZE
	ROW_SIZE        = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE

	PAGE_SIZE       = 4096
	TABLE_MAX_PAGES = 100
	ROWS_PER_PAGE   = PAGE_SIZE / ROW_SIZE
	TABLE_MAX_ROWS  = ROWS_PER_PAGE * TABLE_MAX_PAGES
)

type Pager struct {
	fileDescriptor int
	fileLength     uint32
	pages          []*Pages
}

type Table struct {
	pager   *Pager
	numRows uint32
}

func (row *Row) printRow() {
	fmt.Printf("(%d, %s, %s)\n", row.id, row.username, row.email)
}

func serializeRow(source *Row, void* destination) {
	copy(destination + ID_OFFSET, &(source->id), ID_SIZE)
	copy(destination + USERNAME_OFFSET, &(source->username), USERNAME_SIZE)
	copy(destination + EMAIL_OFFSET, &(source->email), EMAIL_SIZE)
}

func deserializeRow(void* source, destination *Row) {
	copy(&(destination->id), source + ID_OFFSET, ID_SIZE);
	copy(&(destination->username), source + USERNAME_OFFSET, USERNAME_SIZE);
	copy(&(destination->email), source + EMAIL_OFFSET, EMAIL_SIZE);
}

func (pager *Pager) getPage(pageNum uint32) {
	if pageNum > TABLE_MAX_PAGES {
		fmt.Printf("Tried to fetch page number out of bounds. %d > %d\n", pageNum, TABLE_MAX_PAGES)
		os.Exit(1)
	}
	if pager.pages[pageNum] == nil {
		page := Pager{}
		numPages := pager.fileLength / PAGE_SIZE
		if pager.fileLength % PAGE_SIZE != 0 {
			numPages += 1
		}
		if pageNum <= numPages {

		}
	}

}