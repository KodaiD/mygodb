package main

import (
	"io"
	"log"
	"os"
)

type MetaCommandResult int
type PrepareResult int
type StatementType int
type ExecuteResult int

type Statement struct {
	Type 		StatementType
	RowToInsert Row
}

type Row struct {
	Key   string
	Value []byte
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

func newTable() *Table {
	table := &Table{
		NumRows: 0,
		Pager: newPager(),
	}
	return table
}