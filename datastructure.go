package main

import (
	"encoding/binary"
	"fmt"
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
	ID		 uint32
	UserName []byte
	Email	 []byte
}

// pointは書き込まれた値の最後の文字のインデックス
type Page struct {
	buf []byte
	point uint32
}

type Pager struct {
	fileDescriptor *os.File
	pages          []*Page
	point          uint32
}

type Table struct {
	NumRows uint32
	Pager   *Pager
}

func newPager() *Pager {
	file, err := os.OpenFile(FILE_NAME, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		log.Println(err)
	}
	p := &Pager{
		fileDescriptor: file,
		pages:          make([]*Page, TABLE_MAX_PAGES),
	}
	for i := uint32(0); i < TABLE_MAX_PAGES; i++ {
		buf := make([]byte, PAGE_SIZE)
		p.pages[i] =  &Page{buf: buf}
	}
	readData(p)
	return p
}

func newTable() *Table {
	table := &Table{
		NumRows: 0,
		Pager: newPager(),
	}
	return table
}

func (table *Table) closeTable()  {
	// データの永続化
	for i := uint32(0); i <= table.Pager.point; i++ {
		err := binary.Write(table.Pager.fileDescriptor, binary.BigEndian, table.Pager.pages[i].buf)
		if err != nil {
			log.Println(err)
		}
	}
	err := table.Pager.fileDescriptor.Close()
	if err != nil {
		log.Println(err)
	}
}

func readData(pager *Pager) {
	// TODO:
	// とりあえず1ページのみ読み出しにする
	err := binary.Read(pager.fileDescriptor, binary.BigEndian, pager.pages[0].buf)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(pager.pages[0].buf)
}